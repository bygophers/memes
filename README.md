# memes.bygophers.com is an image + caption meme server

https://memes.bygophers.com/

It also serves as an example of a "cloud app" in Go. It's a single
service for now because it's so small, but it does many things in a
cloudy way already.

Feedback is most welcome!

## How it works

### Viewing memes

Every meme (image+caption) is served from a URL like `/i/<id>`.

Good ones are also served from URLs like `/m/<word>`.

Front page serves popular and new memes, with search based on caption
text and base image description.


### Uploading base images

- User uploads a base image (without captions), and provides a name
  for it (e.g. "Curious Gopher").
- Base image is assigned an identifier (derived from its contents, to
  minimize accidental duplicate uploads; on duplicate upload, guide
  user to adding captions to existing base image).
- Base image is stored in Google Cloud Storage.
- Base images are not served from stable URLs, to avoid being used as
  a generic image host.


### Adding captions

- User can browse base images, with search.
- Page `/create/<baseimg>` allows adding a caption to the base image
  with that identifier.
- Submitted meme is assigned an identifier (derived from its contents;
  on duplicate submission, guide user to viewing the original
  submission).
- (Later: assign a shorter id too, to keep URLs short.)
- Meme metadata is stored in Google Cloud Datastore.
- Generated image is reproducible. It might be stored in Google Cloud
  Storage, or we might just cache popular ones in-memory and perhaps
  demonstrate [https://github.com/golang/groupcache](`groupcache`).
  Let's see.


## Running it

Development mode:

```console
$ go get bygophers.com/go/memes/cmd/memes-webserver
$ cat config.json
{}
$ memes-webserver config.json
```

It listens by default on `:8080` for production traffic, and
`localhost:8081` for administrative and debug uses.
