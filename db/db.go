// Package db contains a database abstraction for storing meme
// metadata.
package db

import (
	"errors"
	"io"

	"golang.org/x/net/context"

	"bygophers.com/go/memes/caption"
	"github.com/codahale/blake2"
	"github.com/mohae/uvarint"
	"github.com/tv42/zbase32"
)

// BaseID identifies a base image.
type BaseID string

// MemeID identifies a meme.
type MemeID string

// Meme is an image with text overlaid on it.
type Meme struct {
	Base    BaseID
	Caption caption.Caption
	// TODO keep metadata like who created it
}

func (m *Meme) writeUvar(w io.Writer, n uint64) {
	// the writer is assumed to never fail (hash.Hash or bytes.Buffer)
	const maxUint64Len = 9
	buf := make([]byte, 9)
	l := uvarint.PutUvarint(buf, n)
	buf = buf[:l]
	_, _ = w.Write(buf)
}

func (m *Meme) writeString(w io.Writer, s string) {
	m.writeUvar(w, uint64(len(s)))
	_, _ = io.WriteString(w, s)
}

// ID returns an identifier that is the same for identical input.
// Fields that are considered metadata, such as creation timestamp,
// are not considered.
func (m *Meme) ID() MemeID {
	h := blake2.New(&blake2.Config{
		Size:     16,
		Salt:     []byte{0x6d, 0xc2, 0x94, 0x1b, 0x4b, 0x07, 0x0c, 0x6b},
		Personal: []byte("memes-by-gophers#meme"),
	})
	// version
	m.writeUvar(h, 1)
	m.writeString(h, string(m.Base))
	m.writeString(h, string(m.Caption.Top))
	m.writeString(h, string(m.Caption.Bottom))

	b := h.Sum(nil)
	z := zbase32.EncodeToString(b)
	return MemeID(z)
}

var (
	// ErrExists is returned from DB.Create when a meme with the same
	// ID already exists.
	ErrExists = errors.New("meme exists already")

	// ErrNotFound is returned from DB.Get when a meme is not found.
	ErrNotFound = errors.New("meme not found")
)

// DB describes the operations needed for creating and using memes.
type DB interface {
	// Create a meme. Returns ErrExists with a non-zero MemeID if the
	// meme exists already.
	Create(ctx context.Context, meme *Meme) (MemeID, error)
	// Get a meme. Returns ErrNotFound if meme does not exist.
	Get(ctx context.Context, id MemeID) (*Meme, error)
}
