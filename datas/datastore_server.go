package datas

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/attic-labs/noms/constants"
	"github.com/attic-labs/noms/d"
	"github.com/julienschmidt/httprouter"
)

type connectionState struct {
	c  net.Conn
	cs http.ConnState
}

type dataStoreServer struct {
	dsf     Factory
	port    int
	l       *net.Listener
	csChan  chan *connectionState
	closing bool
}

func NewDataStoreServer(dsf Factory, port int) *dataStoreServer {
	return &dataStoreServer{
		dsf, port, nil, make(chan *connectionState, 16), false,
	}
}

// Run blocks while the dataStoreServer is listening. Running on a separate go routine is supported.
func (s *dataStoreServer) Run() {
	fmt.Printf("Listening on port %d...\n", s.port)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	d.Chk.NoError(err)
	s.l = &l

	router := httprouter.New()

	dsKey, rKey := "datastore", "ref"
	router.GET(fmt.Sprintf("/:%s%s:%s", dsKey, constants.RefPath, rKey), s.makeHandle(HandleRef))
	router.POST(fmt.Sprintf("/:%s%s", dsKey, constants.PostRefsPath), s.makeHandle(HandlePostRefs))
	router.POST(fmt.Sprintf("/:%s%s", dsKey, constants.GetHasPath), s.makeHandle(HandleGetHasRefs))
	router.POST(fmt.Sprintf("/:%s%s", dsKey, constants.GetRefsPath), s.makeHandle(HandleGetRefs))
	router.GET(fmt.Sprintf("/:%s%s", dsKey, constants.RootPath), s.makeHandle(HandleRootGet))
	router.POST(fmt.Sprintf("/:%s%s", dsKey, constants.RootPath), s.makeHandle(HandleRootPost))

	// Handle DEPRECATED endpoints. Remove once JS knows about particular datastores.
	unnamedDs := expectStoreCreate(s.dsf, httprouter.Params{})
	defer unnamedDs.Close()
	router.HandleMethodNotAllowed = false
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Normalize trailing slash.
		p := strings.TrimRight(r.URL.EscapedPath(), "/") + "/"
		ps := httprouter.Params{}
		if strings.HasPrefix(p, constants.RefPath) {
			pathParts := strings.Split(r.URL.EscapedPath()[1:], "/")
			if len(pathParts) > 1 {
				ps = httprouter.Params{{rKey, pathParts[len(pathParts)-1]}}
			}
			HandleRef(w, r, ps, unnamedDs)
		} else if strings.HasPrefix(p, constants.PostRefsPath) {
			HandlePostRefs(w, r, ps, unnamedDs)
		} else if strings.HasPrefix(p, constants.GetHasPath) {
			HandleGetHasRefs(w, r, ps, unnamedDs)
		} else if strings.HasPrefix(p, constants.GetRefsPath) {
			HandleGetRefs(w, r, ps, unnamedDs)
		} else if strings.HasPrefix(p, constants.RootPath) {
			if r.Method == "GET" {
				HandleRootGet(w, r, ps, unnamedDs)
			} else if r.Method == "POST" {
				HandleRootPost(w, r, ps, unnamedDs)
			} else {
				http.NotFound(w, r)
			}
		} else {
			http.NotFound(w, r)
		}
	})
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			router.ServeHTTP(w, req)
		}),
		ConnState: s.connState,
	}

	go func() {
		m := map[net.Conn]http.ConnState{}
		for connState := range s.csChan {
			switch connState.cs {
			case http.StateNew, http.StateActive, http.StateIdle:
				m[connState.c] = connState.cs
			default:
				delete(m, connState.c)
			}
		}
		for c := range m {
			c.Close()
		}
	}()

	srv.Serve(l)
}

func (s *dataStoreServer) makeHandle(hndlr Handler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		ds := expectStoreCreate(s.dsf, ps)
		defer ds.Close()
		hndlr(w, req, ps, ds)
	}
}

func expectStoreCreate(f Factory, ps httprouter.Params) DataStore {
	ds, success := f.Create(ps.ByName("datastore"))
	d.Exp.True(success, "Failed to create datastore named %s", ps.ByName("datastore"))
	return ds
}

func (s *dataStoreServer) connState(c net.Conn, cs http.ConnState) {
	if s.closing {
		d.Chk.Equal(cs, http.StateClosed)
		return
	}
	s.csChan <- &connectionState{c, cs}
}

// Will cause the dataStoreServer to stop listening and an existing call to Run() to continue.
func (s *dataStoreServer) Stop() {
	s.closing = true
	(*s.l).Close()
	close(s.csChan)
}
