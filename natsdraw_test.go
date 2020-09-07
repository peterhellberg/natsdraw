package natsdraw

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/nats-io/nats-server/test"
)

func TestNew(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
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

		m.Set(3, 3, color.NRGBA{0, 255, 0, 255})

		time.Sleep(10 * time.Millisecond)

		if got, want := m.At(4, 6), red; got != want {
			t.Fatalf("m.At(4,6) = %v, want %v", got, want)
		}

		m.Close()
	})

	t.Run("empty name", func(t *testing.T) {
		if _, err := New("", image.Rectangle{}); err != ErrEmptyName {
			t.Fatalf("unexpected error, got %v", err)
		}
	})
}

func TestZeroValue(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		m := &Image{}

		m.Set(1, 1, color.Black)
	})

	t.Run("with name", func(t *testing.T) {
		m := &Image{name: "test"}

		m.Set(1, 1, color.Black)
	})
}
