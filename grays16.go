// Package colorext provides extended color models for use with Go's image package.
package colorext

import (
	"image"
	"image/color"
)

// GrayS16 represents a signed 16-bit grayscale color.
type GrayS16 struct {
	Y int16
}

// RGBA returns the red, green, blue and alpha components of the GrayS16 color.
// This implements the color.Color interface.
// The Y value is converted from the signed range (-32768 to 32767) to
// the unsigned range (0 to 65535) by adding 32768 to shift the range.
func (c GrayS16) RGBA() (r, g, b, a uint32) {
	// Convert signed int16 to unsigned range for RGBA output
	// Shift the range from [-32768, 32767] to [0, 65535]
	y := uint32(int32(c.Y) + 32768)
	return y, y, y, 0xffff
}

// GrayS16Model is the color model for signed 16-bit grayscale colors.
var GrayS16Model color.Model = color.ModelFunc(grayS16Model)

// grayS16Model converts any color.Color to a GrayS16.
func grayS16Model(c color.Color) color.Color {
	if _, ok := c.(GrayS16); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by the standard library.
	// Note that 19595 + 38470 + 7471 equals 65536.
	// The result y will be in the range [0, 65535].
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	// Convert from unsigned [0, 65535] to signed [-32768, 32767]
	// by subtracting 32768. Use int32 for safe intermediate calculation.
	signedY := int32(y) - 32768
	return GrayS16{int16(signedY)}
}

// GrayS16Image is an in-memory image whose At method returns GrayS16 values.
type GrayS16Image struct {
	// Pix holds the image's pixels, as signed 16-bit gray values in big-endian format.
	// The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel returns the GrayS16Image's color model.
func (p *GrayS16Image) ColorModel() color.Model {
	return GrayS16Model
}

// Bounds returns the domain for which At can return non-zero color.
func (p *GrayS16Image) Bounds() image.Rectangle {
	return p.Rect
}

// At returns the color of the pixel at (x, y).
func (p *GrayS16Image) At(x, y int) color.Color {
	return p.GrayS16At(x, y)
}

// GrayS16At returns the GrayS16 color of the pixel at (x, y).
func (p *GrayS16Image) GrayS16At(x, y int) GrayS16 {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return GrayS16{}
	}
	i := p.PixOffset(x, y)
	// Read big-endian int16
	return GrayS16{Y: int16(uint16(p.Pix[i+0])<<8 | uint16(p.Pix[i+1]))}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *GrayS16Image) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
}

// Set sets the pixel at (x, y) to a given color.
func (p *GrayS16Image) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := GrayS16Model.Convert(c).(GrayS16)
	// Write big-endian int16
	p.Pix[i+0] = uint8(uint16(c1.Y) >> 8)
	p.Pix[i+1] = uint8(uint16(c1.Y))
}

// SetGrayS16 sets the pixel at (x, y) to a given GrayS16 color.
func (p *GrayS16Image) SetGrayS16(x, y int, c GrayS16) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	// Write big-endian int16
	p.Pix[i+0] = uint8(uint16(c.Y) >> 8)
	p.Pix[i+1] = uint8(uint16(c.Y))
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *GrayS16Image) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &GrayS16Image{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &GrayS16Image{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *GrayS16Image) Opaque() bool {
	return true
}

// NewGrayS16Image returns a new GrayS16Image with the given bounds.
func NewGrayS16Image(r image.Rectangle) *GrayS16Image {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 2*w*h)
	return &GrayS16Image{
		Pix:    buf,
		Stride: 2 * w,
		Rect:   r,
	}
}
