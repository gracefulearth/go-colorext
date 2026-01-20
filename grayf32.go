// Package colorext provides extended color models for use with Go's image package.
package colorext

import (
	"image"
	"image/color"
	"math"
)

// GrayF32 represents a 32-bit floating point grayscale color.
type GrayF32 struct {
	Y float32
}

// RGBA returns the red, green, blue and alpha components of the GrayF32 color.
// This implements the color.Color interface.
// The Y value is scaled to [0, 65535] for output, with values outside [0, 1]
// being clamped to the valid uint32 range.
func (c GrayF32) RGBA() (r, g, b, a uint32) {
	// Scale Y to [-1, 1] range, then shift to [0, 2]
	y := (c.Y / math.MaxFloat32) + 1.0
	// Scale to [0, 65535]
	y16 := min(uint32(y*32767.5), 0xffff)
	return y16, y16, y16, 0xffff
}

// GrayF32Model is the color model for 32-bit floating point grayscale colors.
var GrayF32Model color.Model = color.ModelFunc(grayF32Model)

// grayF32Model converts any color.Color to a GrayF32.
func grayF32Model(c color.Color) color.Color {
	if _, ok := c.(GrayF32); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by the standard library.
	// Note that 19595 + 38470 + 7471 equals 65536.
	// The result y will be in the range [0, 65535].
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	// Convert from [0, 65535] to [0, 1] floating point range
	return GrayF32{float32(y) / 65535.0}
}

// GrayF32Image is an in-memory image whose At method returns GrayF32 values.
type GrayF32Image struct {
	// Pix holds the image's pixels, as 32-bit floating point gray values in IEEE 754 format.
	// The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel returns the GrayF32Image's color model.
func (p *GrayF32Image) ColorModel() color.Model {
	return GrayF32Model
}

// Bounds returns the domain for which At can return non-zero color.
func (p *GrayF32Image) Bounds() image.Rectangle {
	return p.Rect
}

// At returns the color of the pixel at (x, y).
func (p *GrayF32Image) At(x, y int) color.Color {
	return p.GrayF32At(x, y)
}

// GrayF32At returns the GrayF32 color of the pixel at (x, y).
func (p *GrayF32Image) GrayF32At(x, y int) GrayF32 {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return GrayF32{}
	}
	i := p.PixOffset(x, y)
	// Read 32-bit float in big-endian format
	bits := uint32(p.Pix[i+0])<<24 | uint32(p.Pix[i+1])<<16 | uint32(p.Pix[i+2])<<8 | uint32(p.Pix[i+3])
	return GrayF32{Y: math.Float32frombits(bits)}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *GrayF32Image) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

// Set sets the pixel at (x, y) to a given color.
func (p *GrayF32Image) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := GrayF32Model.Convert(c).(GrayF32)
	// Write 32-bit float in big-endian format
	bits := math.Float32bits(c1.Y)
	p.Pix[i+0] = uint8(bits >> 24)
	p.Pix[i+1] = uint8(bits >> 16)
	p.Pix[i+2] = uint8(bits >> 8)
	p.Pix[i+3] = uint8(bits)
}

// SetGrayF32 sets the pixel at (x, y) to a given GrayF32 color.
func (p *GrayF32Image) SetGrayF32(x, y int, c GrayF32) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	// Write 32-bit float in big-endian format
	bits := math.Float32bits(c.Y)
	p.Pix[i+0] = uint8(bits >> 24)
	p.Pix[i+1] = uint8(bits >> 16)
	p.Pix[i+2] = uint8(bits >> 8)
	p.Pix[i+3] = uint8(bits)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *GrayF32Image) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &GrayF32Image{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &GrayF32Image{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque reports whether the image is fully opaque.
// GrayF32Image is always fully opaque since the GrayF32 color model has no transparency.
func (p *GrayF32Image) Opaque() bool {
	return true
}

// NewGrayF32Image returns a new GrayF32Image with the given bounds.
func NewGrayF32Image(r image.Rectangle) *GrayF32Image {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 4*w*h)
	return &GrayF32Image{
		Pix:    buf,
		Stride: 4 * w,
		Rect:   r,
	}
}
