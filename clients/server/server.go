package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/attic-labs/noms/clients/util"
	"github.com/attic-labs/noms/d"
	"github.com/attic-labs/noms/datas"
)

var (
	port = flag.Int("port", 8000, "")
)

func main() {
	flags := datas.NewFlags()
	flag.Parse()
	dsf, ok := flags.CreateFactory()
	if !ok {
		flag.Usage()
		return
	}

	server := datas.NewDataStoreServer(dsf, *port)

	// Shutdown server gracefully so that profile may be written
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		server.Stop()
	}()

	d.Try(func() {
		if util.MaybeStartCPUProfile() {
			defer util.StopCPUProfile()
		}
		server.Run()
	})
	dsf.Shutter()
}
