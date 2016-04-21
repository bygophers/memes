// Package gds contains an implementation of meme database using
// Google Cloud Datastore.
package gds

import (
	"fmt"

	"bygophers.com/go/memes/db"
	"golang.org/x/net/context"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
)

// Client is a meme database using Google Cloud Datastore.
type Client struct {
	ds *datastore.Client
}

// Open a meme database. Caller is responsible for calling Close when
// done.
func Open(ctx context.Context, projectID string) (*Client, error) {
	client, err := datastore.NewClient(ctx, projectID, cloud.WithUserAgent("memes-by-gophers"))
	if err != nil {
		return nil, err
	}
	return &Client{ds: client}, nil
}

// Close the datastore connection.
func (c *Client) Close() {
	c.ds.Close()
}

var _ db.DB = (*Client)(nil)

// Create a meme.
func (c *Client) Create(ctx context.Context, meme *db.Meme) (db.MemeID, error) {
	id := meme.ID()
	key := datastore.NewKey(ctx, "Meme", string(id), 0, nil)

	create := func(tx *datastore.Transaction) error {
		var dummy db.Meme
		switch err := tx.Get(key, &dummy); err {
		case datastore.ErrNoSuchEntity:
			// good
		case nil:
			return db.ErrExists
		default:
			return err
		}
		if _, err := tx.Put(key, meme); err != nil {
			return fmt.Errorf("datastore put: %v", err)
		}
		return nil
	}
	switch _, err := c.ds.RunInTransaction(ctx, create); err {
	case nil:
		// nothing
	case db.ErrExists:
		return id, err
	default:
		return "", err
	}
	return id, nil
}

// Get a meme.
func (c *Client) Get(ctx context.Context, id db.MemeID) (*db.Meme, error) {
	key := datastore.NewKey(ctx, "Meme", string(id), 0, nil)
	var meme db.Meme
	if err := c.ds.Get(ctx, key, &meme); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, db.ErrNotFound
		}
		return nil, fmt.Errorf("datastore get: %v", err)
	}
	return &meme, nil
}
