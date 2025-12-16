package colorext

import (
	"image"
	"image/color"
	"testing"
)

func TestGrayS16_RGBA(t *testing.T) {
	tests := []struct {
		name string
		c    GrayS16
		want [4]uint32
	}{
		{
			name: "zero value",
			c:    GrayS16{Y: 0},
			want: [4]uint32{32768, 32768, 32768, 0xffff},
		},
		{
			name: "minimum value",
			c:    GrayS16{Y: -32768},
			want: [4]uint32{0, 0, 0, 0xffff},
		},
		{
			name: "maximum value",
			c:    GrayS16{Y: 32767},
			want: [4]uint32{65535, 65535, 65535, 0xffff},
		},
		{
			name: "positive value",
			c:    GrayS16{Y: 16383},
			want: [4]uint32{49151, 49151, 49151, 0xffff},
		},
		{
			name: "negative value",
			c:    GrayS16{Y: -16384},
			want: [4]uint32{16384, 16384, 16384, 0xffff},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, a := tt.c.RGBA()
			if r != tt.want[0] || g != tt.want[1] || b != tt.want[2] || a != tt.want[3] {
				t.Errorf("GrayS16{%d}.RGBA() = (%d, %d, %d, %d), want (%d, %d, %d, %d)",
					tt.c.Y, r, g, b, a, tt.want[0], tt.want[1], tt.want[2], tt.want[3])
			}
		})
	}
}

func TestGrayS16Model(t *testing.T) {
	// Test that GrayS16Model is not nil
	if GrayS16Model == nil {
		t.Fatal("GrayS16Model is nil")
	}

	// Test conversion from GrayS16 returns same value
	original := GrayS16{Y: 1000}
	converted := GrayS16Model.Convert(original)
	if grayS16, ok := converted.(GrayS16); !ok {
		t.Errorf("GrayS16Model.Convert(GrayS16) returned type %T, want GrayS16", converted)
	} else if grayS16.Y != original.Y {
		t.Errorf("GrayS16Model.Convert(GrayS16{%d}) = GrayS16{%d}, want GrayS16{%d}",
			original.Y, grayS16.Y, original.Y)
	}
}

func TestGrayS16Model_ConvertFromRGBA(t *testing.T) {
	tests := []struct {
		name  string
		input color.RGBA
		want  int16
	}{
		{
			name:  "white",
			input: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			want:  32767, // Maximum signed value
		},
		{
			name:  "black",
			input: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			want:  -32768, // Minimum signed value
		},
		{
			name:  "medium gray (128)",
			input: color.RGBA{R: 128, G: 128, B: 128, A: 255},
			want:  128, // Slightly above zero due to 8-bit to 16-bit conversion
		},
		{
			name:  "red",
			input: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			want:  -13173, // Weighted towards red (0.299)
		},
		{
			name:  "green",
			input: color.RGBA{R: 0, G: 255, B: 0, A: 255},
			want:  5701, // Weighted towards green (0.587)
		},
		{
			name:  "blue",
			input: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			want:  -25297, // Weighted towards blue (0.114)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GrayS16Model.Convert(tt.input)
			grayS16, ok := result.(GrayS16)
			if !ok {
				t.Fatalf("GrayS16Model.Convert returned type %T, want GrayS16", result)
			}
			// Allow for small rounding differences
			diff := grayS16.Y - tt.want
			if diff < -1 || diff > 1 {
				t.Errorf("GrayS16Model.Convert(%+v) = GrayS16{%d}, want approximately GrayS16{%d}",
					tt.input, grayS16.Y, tt.want)
			}
		})
	}
}

func TestGrayS16Model_ConvertFromGray16(t *testing.T) {
	tests := []struct {
		name  string
		input color.Gray16
		want  int16
	}{
		{
			name:  "maximum gray16",
			input: color.Gray16{Y: 65535},
			want:  32767,
		},
		{
			name:  "minimum gray16",
			input: color.Gray16{Y: 0},
			want:  -32768,
		},
		{
			name:  "middle gray16",
			input: color.Gray16{Y: 32768},
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GrayS16Model.Convert(tt.input)
			grayS16, ok := result.(GrayS16)
			if !ok {
				t.Fatalf("GrayS16Model.Convert returned type %T, want GrayS16", result)
			}
			if grayS16.Y != tt.want {
				t.Errorf("GrayS16Model.Convert(Gray16{%d}) = GrayS16{%d}, want GrayS16{%d}",
					tt.input.Y, grayS16.Y, tt.want)
			}
		})
	}
}

func TestGrayS16_Implements_Color(t *testing.T) {
	// Compile-time check that GrayS16 implements color.Color
	var _ color.Color = GrayS16{}
}

func TestGrayS16Model_Implements_Model(t *testing.T) {
	// Compile-time check that GrayS16Model implements color.Model
	var _ color.Model = GrayS16Model
}

func TestNewGrayS16Image(t *testing.T) {
	r := image.Rect(0, 0, 10, 10)
	img := NewGrayS16Image(r)

	if img == nil {
		t.Fatal("NewGrayS16Image returned nil")
	}

	if img.Bounds() != r {
		t.Errorf("Bounds() = %v, want %v", img.Bounds(), r)
	}

	if img.Stride != 20 {
		t.Errorf("Stride = %d, want 20", img.Stride)
	}

	expectedLen := 2 * 10 * 10
	if len(img.Pix) != expectedLen {
		t.Errorf("len(Pix) = %d, want %d", len(img.Pix), expectedLen)
	}
}

func TestGrayS16Image_Implements_Image(t *testing.T) {
	// Compile-time check that GrayS16Image implements image.Image
	var _ image.Image = &GrayS16Image{}
}

func TestGrayS16Image_ColorModel(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))
	if img.ColorModel() != GrayS16Model {
		t.Errorf("ColorModel() returned %v, want GrayS16Model", img.ColorModel())
	}
}

func TestGrayS16Image_Bounds(t *testing.T) {
	tests := []struct {
		name string
		rect image.Rectangle
	}{
		{"zero origin", image.Rect(0, 0, 10, 10)},
		{"non-zero origin", image.Rect(5, 5, 15, 15)},
		{"negative origin", image.Rect(-5, -5, 5, 5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := NewGrayS16Image(tt.rect)
			if img.Bounds() != tt.rect {
				t.Errorf("Bounds() = %v, want %v", img.Bounds(), tt.rect)
			}
		})
	}
}

func TestGrayS16Image_SetAndGet(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	tests := []struct {
		name  string
		x, y  int
		color GrayS16
	}{
		{"zero value", 0, 0, GrayS16{Y: 0}},
		{"minimum value", 1, 1, GrayS16{Y: -32768}},
		{"maximum value", 2, 2, GrayS16{Y: 32767}},
		{"positive value", 3, 3, GrayS16{Y: 16383}},
		{"negative value", 4, 4, GrayS16{Y: -16384}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img.SetGrayS16(tt.x, tt.y, tt.color)
			got := img.GrayS16At(tt.x, tt.y)
			if got.Y != tt.color.Y {
				t.Errorf("After SetGrayS16(%d, %d, GrayS16{%d}), GrayS16At(%d, %d) = GrayS16{%d}, want GrayS16{%d}",
					tt.x, tt.y, tt.color.Y, tt.x, tt.y, got.Y, tt.color.Y)
			}
		})
	}
}

func TestGrayS16Image_Set(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	// Test setting with color.Color interface
	c := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	img.Set(5, 5, c)

	got := img.At(5, 5)
	if _, ok := got.(GrayS16); !ok {
		t.Errorf("At() returned type %T, want GrayS16", got)
	}
}

func TestGrayS16Image_At(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	// Set a value and verify At returns it via color.Color interface
	expected := GrayS16{Y: 1000}
	img.SetGrayS16(5, 5, expected)

	got := img.At(5, 5)
	grayS16, ok := got.(GrayS16)
	if !ok {
		t.Fatalf("At() returned type %T, want GrayS16", got)
	}
	if grayS16.Y != expected.Y {
		t.Errorf("At(5, 5) = GrayS16{%d}, want GrayS16{%d}", grayS16.Y, expected.Y)
	}
}

func TestGrayS16Image_GrayS16At_OutOfBounds(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	tests := []struct {
		name string
		x, y int
	}{
		{"negative x", -1, 5},
		{"negative y", 5, -1},
		{"x too large", 10, 5},
		{"y too large", 5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := img.GrayS16At(tt.x, tt.y)
			if got.Y != 0 {
				t.Errorf("GrayS16At(%d, %d) = GrayS16{%d}, want GrayS16{0} for out of bounds",
					tt.x, tt.y, got.Y)
			}
		})
	}
}

func TestGrayS16Image_SetGrayS16_OutOfBounds(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	// Setting out of bounds should not panic
	img.SetGrayS16(-1, 5, GrayS16{Y: 100})
	img.SetGrayS16(5, -1, GrayS16{Y: 100})
	img.SetGrayS16(10, 5, GrayS16{Y: 100})
	img.SetGrayS16(5, 10, GrayS16{Y: 100})
}

func TestGrayS16Image_PixOffset(t *testing.T) {
	tests := []struct {
		name   string
		rect   image.Rectangle
		x, y   int
		offset int
	}{
		{"zero origin (0,0)", image.Rect(0, 0, 10, 10), 0, 0, 0},
		{"zero origin (1,0)", image.Rect(0, 0, 10, 10), 1, 0, 2},
		{"zero origin (0,1)", image.Rect(0, 0, 10, 10), 0, 1, 20},
		{"zero origin (5,5)", image.Rect(0, 0, 10, 10), 5, 5, 110},
		{"non-zero origin (5,5)", image.Rect(5, 5, 15, 15), 5, 5, 0},
		{"non-zero origin (6,5)", image.Rect(5, 5, 15, 15), 6, 5, 2},
		{"non-zero origin (5,6)", image.Rect(5, 5, 15, 15), 5, 6, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := NewGrayS16Image(tt.rect)
			got := img.PixOffset(tt.x, tt.y)
			if got != tt.offset {
				t.Errorf("PixOffset(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.offset)
			}
		})
	}
}

func TestGrayS16Image_SubImage(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	// Set some pixels in the original image
	img.SetGrayS16(5, 5, GrayS16{Y: 1000})
	img.SetGrayS16(6, 6, GrayS16{Y: 2000})

	// Create a sub-image
	sub := img.SubImage(image.Rect(5, 5, 8, 8))
	subImg, ok := sub.(*GrayS16Image)
	if !ok {
		t.Fatalf("SubImage returned type %T, want *GrayS16Image", sub)
	}

	// Verify bounds
	expectedBounds := image.Rect(5, 5, 8, 8)
	if subImg.Bounds() != expectedBounds {
		t.Errorf("SubImage bounds = %v, want %v", subImg.Bounds(), expectedBounds)
	}

	// Verify the sub-image shares pixels with original
	got := subImg.GrayS16At(5, 5)
	if got.Y != 1000 {
		t.Errorf("SubImage.GrayS16At(5, 5) = GrayS16{%d}, want GrayS16{1000}", got.Y)
	}

	// Modify the sub-image and verify it affects the original
	subImg.SetGrayS16(6, 6, GrayS16{Y: 3000})
	got = img.GrayS16At(6, 6)
	if got.Y != 3000 {
		t.Errorf("After modifying SubImage, original GrayS16At(6, 6) = GrayS16{%d}, want GrayS16{3000}", got.Y)
	}
}

func TestGrayS16Image_SubImage_Empty(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	// Create an empty sub-image
	sub := img.SubImage(image.Rect(5, 5, 5, 5))
	subImg, ok := sub.(*GrayS16Image)
	if !ok {
		t.Fatalf("SubImage returned type %T, want *GrayS16Image", sub)
	}

	if !subImg.Bounds().Empty() {
		t.Errorf("Empty SubImage bounds = %v, want empty rectangle", subImg.Bounds())
	}
}

func TestGrayS16Image_SubImage_NonIntersecting(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	// Create a non-intersecting sub-image
	sub := img.SubImage(image.Rect(20, 20, 30, 30))
	subImg := sub.(*GrayS16Image)

	if !subImg.Bounds().Empty() {
		t.Errorf("Non-intersecting SubImage bounds = %v, want empty rectangle", subImg.Bounds())
	}
}

func TestGrayS16Image_Opaque(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 10, 10))

	if !img.Opaque() {
		t.Error("Opaque() = false, want true")
	}
}

func TestGrayS16Image_NonZeroOrigin(t *testing.T) {
	// Test with non-zero origin
	img := NewGrayS16Image(image.Rect(5, 5, 15, 15))

	// Set and get a pixel
	expected := GrayS16{Y: 1234}
	img.SetGrayS16(7, 8, expected)

	got := img.GrayS16At(7, 8)
	if got.Y != expected.Y {
		t.Errorf("GrayS16At(7, 8) = GrayS16{%d}, want GrayS16{%d}", got.Y, expected.Y)
	}

	// Verify out of bounds below origin
	got = img.GrayS16At(4, 4)
	if got.Y != 0 {
		t.Errorf("GrayS16At(4, 4) = GrayS16{%d}, want GrayS16{0} (out of bounds)", got.Y)
	}
}

func TestGrayS16Image_BigEndianEncoding(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 1, 1))

	// Test big-endian encoding by setting a specific value and checking bytes
	img.SetGrayS16(0, 0, GrayS16{Y: 0x1234})

	// In big-endian, 0x1234 should be stored as [0x12, 0x34]
	if img.Pix[0] != 0x12 || img.Pix[1] != 0x34 {
		t.Errorf("Big-endian encoding: Pix = [0x%02x, 0x%02x], want [0x12, 0x34]",
			img.Pix[0], img.Pix[1])
	}

	// Test reading back
	got := img.GrayS16At(0, 0)
	if got.Y != 0x1234 {
		t.Errorf("GrayS16At(0, 0) = GrayS16{0x%04x}, want GrayS16{0x1234}", uint16(got.Y))
	}
}

func TestGrayS16Image_NegativeValueEncoding(t *testing.T) {
	img := NewGrayS16Image(image.Rect(0, 0, 1, 1))

	// Test negative value encoding
	img.SetGrayS16(0, 0, GrayS16{Y: -1})

	// -1 as int16 is 0xFFFF in two's complement
	if img.Pix[0] != 0xFF || img.Pix[1] != 0xFF {
		t.Errorf("Negative value encoding: Pix = [0x%02x, 0x%02x], want [0xFF, 0xFF]",
			img.Pix[0], img.Pix[1])
	}

	// Test reading back
	got := img.GrayS16At(0, 0)
	if got.Y != -1 {
		t.Errorf("GrayS16At(0, 0) = GrayS16{%d}, want GrayS16{-1}", got.Y)
	}
}
