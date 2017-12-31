package natsdraw

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/nats-io/gnatsd/test"
)

func TestNew(t *testing.T) {
	s := test.RunDefaultServer()
	defer s.Shutdown()

	r := image.Rect(0, 0, 7, 7)

	m, err := New("test-new", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := m.name, "test-new"; got != want {
		t.Fatalf("m.name = %q, want %q", got, want)
	}

	if !m.Bounds().Eq(r) {
		t.Fatalf("m.Bounds() = %v, want %v", m.Bounds(), r)
	}

	red := color.RGBA{255, 0, 0, 255}

	m.Set(4, 6, red)

	time.Sleep(10 * time.Millisecond)

	if got, want := m.At(4, 6), red; got != want {
		t.Fatalf("m.At(4,6) = %v, want %v", got, want)
	}
}
