package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bygophers.com/go/memes/db"
	"bygophers.com/go/memes/db/gds"

	"golang.org/x/net/context"
)

// TODO take from gcloud settings?
const projectID = "memes-by-gophers"

func doit(ctx context.Context, id db.MemeID) error {
	client, err := gds.Open(ctx, projectID)
	if err != nil {
		return fmt.Errorf("datastore client: %v", err)
	}
	defer client.Close()

	meme, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot get meme: %v", err)
	}
	fmt.Printf("%#v\n", meme)
	return nil
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s [OPTS] ID\n", prog)
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
	id := db.MemeID(flag.Arg(0))

	ctx := context.Background()
	if err := doit(ctx, id); err != nil {
		log.Fatal(err)
	}
}
