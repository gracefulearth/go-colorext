// Package colorext provides extended color models for use with Go's image package.
package colorext

import "image/color"

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
