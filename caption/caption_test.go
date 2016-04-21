package caption_test

import (
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"bygophers.com/go/memes/caption"
	"bygophers.com/go/memes/caption/internal/approve"
)

func load(t testing.TB, name string) image.Image {
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("loading test image: %v", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decoding test image: %v", err)
	}
	return img
}

func TestSimple(t *testing.T) {
	t.Parallel()
	base := load(t, "chipmunk.jpg")
	req := &caption.Caption{
		Top:    "chipmunks gonna",
		Bottom: "chip",
	}
	img, err := caption.Draw(base, req)
	if err != nil {
		t.Fatalf("caption: %v", err)
	}
	if err := approve.Image(img); err != nil {
		t.Fatalf("not approved: %v", err)
	}
}

func BenchmarkCaption(b *testing.B) {
	base := load(b, "chipmunk.jpg")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &caption.Caption{
			Top:    "chipmunks gonna",
			Bottom: "chip",
		}
		img, err := caption.Draw(base, req)
		if err != nil {
			b.Fatalf("caption: %v", err)
		}
		_ = img
	}
}
