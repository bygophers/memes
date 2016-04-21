package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bygophers.com/go/memes/caption"
	"bygophers.com/go/memes/db"
	"bygophers.com/go/memes/db/gds"

	"golang.org/x/net/context"
)

// TODO take from gcloud settings?
const projectID = "memes-by-gophers"

func doit(ctx context.Context, baseID db.BaseID, top, bottom string) error {
	client, err := gds.Open(ctx, projectID)
	if err != nil {
		return fmt.Errorf("datastore client: %v", err)
	}
	defer client.Close()

	meme := &db.Meme{
		Base: baseID,
		Caption: caption.Caption{
			Top:    top,
			Bottom: bottom,
		},
	}
	id, err := client.Create(ctx, meme)
	if err != nil {
		return fmt.Errorf("cannot create meme: %v", err)
	}
	fmt.Println(id)
	return nil
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s [OPTS] BASE_ID TOP_CAPTION BOTTOM_CAPTION\n", prog)
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(2)
	}
	baseID := db.BaseID(flag.Arg(0))
	top := flag.Arg(1)
	bottom := flag.Arg(2)

	ctx := context.Background()
	if err := doit(ctx, baseID, top, bottom); err != nil {
		log.Fatal(err)
	}
}
