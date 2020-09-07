package natsdraw

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"time"

	nats "github.com/nats-io/nats.go"
)

// LoadTimeout is the timeout for loading the image
var LoadTimeout = 500 * time.Millisecond

// ErrEmptyName is returned if name is empty
var ErrEmptyName = errors.New("empty name")

// Image implements draw.Image over NATS
type Image struct {
	*image.RGBA
	name string
	conn *nats.EncodedConn
}

// New creates a new *Image
func New(name string, r image.Rectangle, options ...Option) (*Image, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	m := &Image{name: name, RGBA: image.NewRGBA(r)}

	for _, o := range options {
		if err := o(m); err != nil {
			return nil, err
		}
	}

	if m.conn == nil {
		if err := Connect(nats.DefaultURL)(m); err != nil {
			return nil, err
		}
	}

	if msg, err := m.conn.Conn.Request(m.subject(), nil, LoadTimeout); err == nil {
		if p, err := png.Decode(bytes.NewReader(msg.Data)); err == nil {
			m.RGBA = image.NewRGBA(p.Bounds())
			draw.Draw(m.RGBA, p.Bounds(), p, image.ZP, draw.Src)
		}
	}

	m.conn.Conn.Subscribe(m.subject(), func(msg *nats.Msg) {
		buf := new(bytes.Buffer)

		if err := png.Encode(buf, m.RGBA); err == nil {
			m.conn.Conn.Publish(msg.Reply, buf.Bytes())
		}
	})

	m.conn.Subscribe(m.setSubject(), func(p *Pixel) {
		m.RGBA.Set(p.X, p.Y, p.C)
	})

	return m, nil
}

// Option is the functional option type for *Image
type Option func(*Image) error

// Connect to the given NATS URL
func Connect(url string) Option {
	return func(m *Image) error {
		nc, err := nats.Connect(url)
		if err != nil {
			return err
		}

		ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
		if err != nil {
			return err
		}

		m.conn = ec

		return nil
	}
}

// Close the connection to NATS
func (m *Image) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

// Set the pixel at x,y to color c
func (m *Image) Set(x, y int, c color.Color) {
	if m.name == "" || m.RGBA == nil || m.conn == nil {
		return
	}

	m.RGBA.Set(x, y, c)

	s := m.setSubject()
	p := NewPixel(x, y, c)

	m.conn.Publish(s, p)
}

func (m *Image) subject() string {
	return "natsdraw." + m.name
}

func (m *Image) setSubject() string {
	return m.subject() + ".Set"
}

// Pixel contains the position and color
type Pixel struct {
	X int
	Y int
	C color.RGBA
}

// NewPixel creates a new *Pixel
func NewPixel(x, y int, c color.Color) *Pixel {
	switch c := c.(type) {
	case color.RGBA:
		return &Pixel{x, y, c}
	default:
		r, g, b, a := c.RGBA()
		return &Pixel{x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}}
	}
}

var _ draw.Image = &Image{}
