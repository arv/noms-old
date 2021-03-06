package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/attic-labs/noms/clients/util"
	"github.com/attic-labs/noms/d"
	"github.com/attic-labs/noms/dataset"
	"github.com/attic-labs/noms/util/http/retry"
	"github.com/bradfitz/latlong"
	"github.com/garyburd/go-oauth/oauth"
)

var (
	clientIdFlag     = flag.String("client-id", "", "API keys for flickr can be created at https://www.flickr.com/services/apps/create/apply")
	clientSecretFlag = flag.String("client-secret", "", "API keys for flickr can be created at https://www.flickr.com/services/apps/create/apply")
	clientFlags      = util.NewFlags()
	ds               *dataset.Dataset
	httpClient       *http.Client
	oauthClient      oauth.Client
	tokenFlag        = flag.String("token", "", "OAuth1 token (if ommitted, flickr will attempt web auth)")
	tokenSecretFlag  = flag.String("token-secret", "", "OAuth1 token secret (if ommitted, flickr will attempt web auth)")
	user             User
)

type progressTracker struct {
	didLogin                 bool
	didGetList               bool
	numPhotos, photoProgress int
}

func (pt *progressTracker) Update() {
	progress := float32(0)
	if pt.didLogin {
		progress += 0.1
	}
	if pt.didGetList {
		progress += 0.1
	}
	if pt.numPhotos > 0 {
		remaining := 1.0 - progress
		progress += remaining * (float32(pt.photoProgress) / float32(pt.numPhotos))
	}
	clientFlags.UpdateProgress(progress)
}

type flickrAPI interface {
	Call(method string, response interface{}, args *map[string]string) error
}

type flickrCall struct {
	Stat string
}

type idAndRefOfAlbum struct {
	id  string
	ref RefOfAlbum
}

func main() {
	flag.Usage = flickrUsage
	flag.Parse()

	httpClient = util.CachingHttpClient()

	if *clientIdFlag == "" || *clientSecretFlag == "" || httpClient == nil {
		flag.Usage()
		os.Exit(1)
	}

	ds = clientFlags.CreateDataset()
	if ds == nil {
		flag.Usage()
		os.Exit(1)
	}
	defer ds.Store().Close()

	if err := clientFlags.CreateProgressFile(); err != nil {
		fmt.Println(err)
	} else {
		defer clientFlags.CloseProgressFile()
	}

	oauthClient = oauth.Client{
		TemporaryCredentialRequestURI: "https://www.flickr.com/services/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://www.flickr.com/services/oauth/authorize",
		TokenRequestURI:               "https://www.flickr.com/services/oauth/access_token",
		Credentials: oauth.Credentials{
			Token:  *clientIdFlag,
			Secret: *clientSecretFlag,
		},
	}

	var tokenCred *oauth.Credentials
	if *tokenFlag != "" && *tokenSecretFlag != "" {
		tokenCred = &oauth.Credentials{*tokenFlag, *tokenSecretFlag}
	} else {
		tokenCred = requestTokenCredentials()
	}

	api := liveFlickrAPI{tokenCred}
	if !findUser() {
		if err := requestUser(api); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	progress := progressTracker{}
	progress.didLogin = true
	progress.Update()

	user = user.SetAlbums(getAlbums(api, &progress))
	commitUser()
}

func flickrUsage() {
	essay := `The Flickr importer needs 4 items for authentication: a client ID + secret, and a token + secret:
  * For the client ID + secret, you need to apply for an API key through https://www.flickr.com/services/apps/create/apply. Specify these as -client-id and -client-secret.
  * For the token + secret, you have 2 options:
    (a) specify them as -token and -secret, or
    (b) let the Flickr importer request them for you (don't specify anything).
    The latter will take you through web auth, which on the last page will show the token and secret to use for -token and -secret from now on (unless they get revoked).`
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s\n\n", essay)
}

func findUser() bool {
	if commit, ok := ds.MaybeHead(); ok {
		if userRef, ok := commit.Value().(RefOfUser); ok {
			user = userRef.TargetValue(ds.Store())
			return true
		}
	}

	return false
}

func requestUser(api flickrAPI) error {
	response := struct {
		flickrCall
		User struct {
			Id       string `json:"id"`
			Username struct {
				Content string `json:"_content"`
			} `json:"username"`
		} `json:"user"`
	}{}

	if err := api.Call("flickr.test.login", &response, nil); err != nil {
		return err
	}

	user = user.SetId(response.User.Id).SetName(response.User.Username.Content)
	return nil
}

func requestTokenCredentials() *oauth.Credentials {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	d.Chk.NoError(err)

	callbackURL := "http://" + l.Addr().String()
	tempCred, err := oauthClient.RequestTemporaryCredentials(nil, callbackURL, url.Values{
		"perms": []string{"read"},
	})
	// If we ever hear anything from the oauth handshake, it'll be acceptance. The user declining will mean we never get called.
	d.Chk.NoError(err)

	authUrl := oauthClient.AuthorizationURL(tempCred, nil)
	fmt.Printf("Visit the following URL to authorize access to your Flickr data: %v\n", authUrl)
	tokenCred, err := awaitOAuthResponse(l, tempCred)
	d.Chk.NoError(err)
	return tokenCred
}

func getAlbum(api flickrAPI, id string, gotPhoto chan struct{}) idAndRefOfAlbum {
	response := struct {
		flickrCall
		Photoset struct {
			Id    string `json:"id"`
			Title struct {
				Content string `json:"_content"`
			} `json:"title"`
		} `json:"photoset"`
	}{}

	err := api.Call("flickr.photosets.getInfo", &response, &map[string]string{
		"photoset_id": id,
		"user_id":     user.Id(),
	})
	d.Chk.NoError(err)

	photos := getAlbumPhotos(api, id, gotPhoto)

	fmt.Printf("Photoset: %v\n", response.Photoset.Title.Content)

	album := NewAlbum().
		SetId(id).
		SetTitle(response.Photoset.Title.Content).
		SetPhotos(photos)
	// TODO: Write albums in batches.
	ref := ds.Store().WriteValue(album).(RefOfAlbum)
	return idAndRefOfAlbum{id, ref}
}

func getAlbums(api flickrAPI, progress *progressTracker) MapOfStringToRefOfAlbum {
	response := struct {
		flickrCall
		Photosets struct {
			Photoset []struct {
				Id    string `json:"id"`
				Title struct {
					Content string `json:"_content"`
				} `json:"title"`
				Photos int `json:"photos"`
			} `json:"photoset"`
		} `json:"photosets"`
	}{}

	err := api.Call("flickr.photosets.getList", &response, nil)
	d.Chk.NoError(err)

	progress.didGetList = true
	progress.Update()

	for _, ps := range response.Photosets.Photoset {
		progress.numPhotos += ps.Photos
	}

	gotPhoto := make(chan struct{}, clientFlags.Concurrency())
	go func() {
		lastUpdate := time.Now()
		for range gotPhoto {
			progress.photoProgress++
			if now := time.Now(); now.Sub(lastUpdate)/time.Millisecond >= 200 {
				progress.Update()
				lastUpdate = now
			}
		}
	}()

	out := make(chan idAndRefOfAlbum, clientFlags.Concurrency())
	for _, p := range response.Photosets.Photoset {
		p := p
		go func() {
			out <- getAlbum(api, p.Id, gotPhoto)
		}()
	}

	albums := NewMapOfStringToRefOfAlbum()
	for range response.Photosets.Photoset {
		a := <-out
		albums = albums.Set(a.id, a.ref)
	}

	close(gotPhoto)
	return albums
}

func getAlbumPhotos(api flickrAPI, id string, gotPhoto chan struct{}) SetOfRefOfRemotePhoto {
	response := struct {
		flickrCall
		Photoset struct {
			Photo []struct {
				DateTaken      string      `json:"datetaken"`
				Id             string      `json:"id"`
				Title          string      `json:"title"`
				Tags           string      `json:"tags"`
				ThumbURL       string      `json:"url_t"`
				ThumbWidth     interface{} `json:"width_t"`
				ThumbHeight    interface{} `json:"height_t"`
				SmallURL       string      `json:"url_s"`
				SmallWidth     interface{} `json:"width_s"`
				SmallHeight    interface{} `json:"height_s"`
				Latitude       interface{} `json:"latitude"`
				Longitude      interface{} `json:"longitude"`
				MediumURL      string      `json:"url_m"`
				MediumWidth    interface{} `json:"width_m"`
				MediumHeight   interface{} `json:"height_m"`
				LargeURL       string      `json:"url_l"`
				LargeWidth     interface{} `json:"width_l"`
				LargeHeight    interface{} `json:"height_l"`
				OriginalURL    string      `json:"url_o"`
				OriginalWidth  interface{} `json:"width_o"`
				OriginalHeight interface{} `json:"height_o"`
			} `json:"photo"`
		} `json:"photoset"`
	}{}

	// TODO: Implement paging. This call returns a maximum of 500 pictures in each response.
	err := api.Call("flickr.photosets.getPhotos", &response, &map[string]string{
		"photoset_id": id,
		"user_id":     user.Id(),
		"extras":      "date_taken,geo,tags,url_t,url_s,url_m,url_l,url_o",
	})
	d.Chk.NoError(err)

	store := ds.Store()
	photos := NewSetOfRefOfRemotePhoto()

	for _, p := range response.Photoset.Photo {
		photo := RemotePhotoDef{
			Id:    p.Id,
			Title: p.Title,
			Tags:  getTags(p.Tags),
		}.New()

		lat, lon := deFlickr(p.Latitude), deFlickr(p.Longitude)

		// Flickr doesn't give timezone information (in fairness, neither does EXIF), so try to figure it out from the geolocation data. This is imperfect because it won't give us daylight savings. If there is no geolocation data then assume the location is PST - it's better than GMT.
		zone := "America/Los_Angeles"
		if lat != 0.0 && lon != 0.0 {
			if z := latlong.LookupZoneName(lat, lon); z != "" {
				zone = z
			}
		}
		location, err := time.LoadLocation(zone)
		d.Chk.NoError(err)

		// DateTaken is the MySQL DATETIME format.
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", p.DateTaken, location); err == nil {
			photo = photo.SetDate(DateDef{t.Unix() * 1e3}.New())
		} else {
			fmt.Printf("Error parsing date \"%s\": %s\n", p.DateTaken, err)
		}

		sizes := NewMapOfSizeToString()
		sizes = addSize(sizes, p.ThumbURL, p.ThumbWidth, p.ThumbHeight)
		sizes = addSize(sizes, p.SmallURL, p.SmallWidth, p.SmallHeight)
		sizes = addSize(sizes, p.MediumURL, p.MediumWidth, p.MediumHeight)
		sizes = addSize(sizes, p.LargeURL, p.LargeWidth, p.LargeHeight)
		sizes = addSize(sizes, p.OriginalURL, p.OriginalWidth, p.OriginalHeight)
		photo = photo.SetSizes(sizes)

		if lat != 0.0 && lon != 0.0 {
			photo = photo.SetGeoposition(GeopositionDef{float32(lat), float32(lon)}.New())
		}

		// TODO: Write photos in batches.
		photos = photos.Insert(store.WriteValue(photo).(RefOfRemotePhoto))
		gotPhoto <- struct{}{}
	}

	return photos
}

func getTags(tagStr string) (tags SetOfStringDef) {
	tags = SetOfStringDef{}

	if tagStr == "" {
		return
	}

	for _, tag := range strings.Split(tagStr, " ") {
		tags[tag] = true
	}
	return
}

func deFlickr(argh interface{}) float64 {
	switch argh := argh.(type) {
	case float64:
		return argh
	case string:
		f64, err := strconv.ParseFloat(argh, 64)
		d.Chk.NoError(err)
		return float64(f64)
	default:
		return 0.0
	}
}

func addSize(sizes MapOfSizeToString, url string, width interface{}, height interface{}) MapOfSizeToString {
	getDim := func(v interface{}) uint32 {
		switch v := v.(type) {
		case float64:
			return uint32(v)
		case string:
			i, err := strconv.Atoi(v)
			d.Chk.NoError(err)
			return uint32(i)
		default:
			d.Chk.Fail(fmt.Sprintf("Unexpected value for image width or height: %+v", v))
			return uint32(0)
		}
	}
	if url == "" {
		return sizes
	}

	return sizes.Set(SizeDef{getDim(width), getDim(height)}.New(), url)
}

func awaitOAuthResponse(l net.Listener, tempCred *oauth.Credentials) (tokenCred *oauth.Credentials, err error) {
	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/plain")
		tokenCred, _, err = oauthClient.RequestToken(nil, tempCred, r.FormValue("oauth_verifier"))
		if err != nil {
			fmt.Fprintf(w, "%v", err)
		} else {
			d.Chk.NotNil(tokenCred)
			fmt.Fprintf(w, "Authorized: token %s token-secret %s", tokenCred.Token, tokenCred.Secret)
		}
		l.Close()
	})}
	srv.Serve(l)
	return
}

func commitUser() {
	var err error
	r := ds.Store().WriteValue(user).(RefOfUser)
	*ds, err = ds.Commit(r)
	d.Exp.NoError(err)
}

type liveFlickrAPI struct {
	tokenCred *oauth.Credentials
}

func (api liveFlickrAPI) Call(method string, response interface{}, args *map[string]string) error {
	restURL := "https://api.flickr.com/services/rest/"

	values := url.Values{
		"method":         []string{method},
		"format":         []string{"json"},
		"nojsoncallback": []string{"1"},
	}

	if args != nil {
		for k, v := range *args {
			values[k] = []string{v}
		}
	}

	res := retry.Request(restURL, func() (*http.Response, error) {
		return oauthClient.Get(nil, api.tokenCred, restURL, values)
	})

	defer res.Body.Close()
	buff, err := ioutil.ReadAll(res.Body)
	d.Chk.NoError(err)
	if err = json.Unmarshal(buff, response); err != nil {
		return err
	}

	status := reflect.ValueOf(response).Elem().FieldByName("Stat").Interface().(string)
	if status != "ok" {
		err = errors.New(fmt.Sprintf("Failed flickr API call: %v, status: %v", method, status))
	}
	return nil
}
