package main

import (
	"encoding/json"
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"

	"github.com/codahale/http-handlers/logging"
)

var (
	httpAddr  = flag.String("http", ":8080", "host:port to serve production HTTP traffic on")
	adminAddr = flag.String("admin", "localhost:8081", "host:port to serve admin and debug")
)

type config struct {
}

func loadConfig(p string) (*config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var conf config
	dec := json.NewDecoder(f)
	if err := dec.Decode(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func serve(conf *config) error {
	errCh := make(chan error, 2)

	{
		mux := http.NewServeMux()
		mux.HandleFunc("/",
			func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprintf(w, "Hello, world\n")
			},
		)
		// TODO better logging, this is just a debug helper for now
		logged := logging.Wrap(mux, os.Stderr)
		logged.Start()
		defer logged.Stop()
		srv := &http.Server{
			Addr:         *httpAddr,
			Handler:      logged,
			ReadTimeout:  4 * time.Hour,
			WriteTimeout: 4 * time.Hour,
		}
		go func() {
			errCh <- srv.ListenAndServe()
		}()
	}

	{
		mux := http.NewServeMux()
		mux.HandleFunc("/",
			func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprintf(w, "TODO debug interface\n")
			},
		)
		mux.Handle("/debug/", http.DefaultServeMux)
		mux.Handle("/debug", http.DefaultServeMux)
		srv := &http.Server{
			Addr:         *adminAddr,
			Handler:      mux,
			ReadTimeout:  4 * time.Hour,
			WriteTimeout: 4 * time.Hour,
		}
		go func() {
			errCh <- srv.ListenAndServe()
		}()
	}

	return <-errCh
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s [OPTS] CONFIG\n", prog)
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	configPath := flag.Arg(0)

	conf, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := serve(conf); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
