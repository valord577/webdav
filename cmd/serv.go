package cmd

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/valord577/clix"
	log "github.com/valord577/webdav/logger"
	"github.com/valord577/webdav/rt"
	"github.com/valord577/webdav/serve"
)

// @author valor.

var servCmd = &clix.Command{
	Name:    "serv",
	Summary: "Startup webdav server.",

	Run: runServ,
}

var cfg string

func init() {
	servCmd.FlagStringVar(&cfg, "c", "", "Declare the configuration file path")
}

func runServ(_ *clix.Command, _ []string) error {

	if cfg == "" {
		msg := "no configuration file declared"
		log.Errorf(msg)
		return errors.New(msg)
	}
	err := rt.ReadInFile(cfg)
	if err != nil {
		msg := "error to read cfg file: " + err.Error()
		log.Errorf(msg)
		return errors.New(msg)
	}
	log.Infof("activated cfg file: %s", cfg)

	serv := serve.WebDAVServ()

	go func() {
		var err error = nil

		if rt.UseTLS() {
			serv.TLSConfig = &tls.Config{
				MinVersion: rt.MinVersionTLS(),
			}
			err = serv.ListenAndServeTLS(rt.CerFilepath(), rt.CerFilepath())
		} else {
			err = serv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Block and listen for signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	s := <-sig
	log.Infof("starting to shutdown by signal: %d", s)

	err = serv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
	log.Infof("shutdown successfully")
	return err
}
