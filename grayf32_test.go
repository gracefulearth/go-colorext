package colorext

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestGrayF32_RGBA(t *testing.T) {
	tests := []struct {
		name string
		c    GrayF32
		want [4]uint32
	}{
		{
			name: "zero value",
			c:    GrayF32{Y: 0.0},
			want: [4]uint32{32767, 32767, 32767, 0xffff},
		},
		{
			name: "maximum value",
			c:    GrayF32{Y: math.MaxFloat32},
			want: [4]uint32{65535, 65535, 65535, 0xffff},
		},
		{
			name: "middle value",
			c:    GrayF32{Y: math.MaxFloat32 / 2},
			want: [4]uint32{49151, 49151, 49151, 0xffff},
		},
		{
			name: "quarter value",
			c:    GrayF32{Y: math.MaxFloat32 / 4},
			want: [4]uint32{40959, 40959, 40959, 0xffff},
		},
		{
			name: "three quarters value",
			c:    GrayF32{Y: 3 * (math.MaxFloat32 / 4)},
			want: [4]uint32{57343, 57343, 57343, 0xffff},
		},
		{
			name: "negative value",
			c:    GrayF32{Y: -math.MaxFloat32},
			want: [4]uint32{0, 0, 0, 0xffff},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, a := tt.c.RGBA()
			// Allow for small rounding differences
			if abs(int(r)-int(tt.want[0])) > 1 || abs(int(g)-int(tt.want[1])) > 1 ||
				abs(int(b)-int(tt.want[2])) > 1 || a != tt.want[3] {
				t.Errorf("GrayF32{%f}.RGBA() = (%d, %d, %d, %d), want approximately (%d, %d, %d, %d)",
					tt.c.Y, r, g, b, a, tt.want[0], tt.want[1], tt.want[2], tt.want[3])
			}
		})
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestGrayF32Model(t *testing.T) {
	// Test that GrayF32Model is not nil
	if GrayF32Model == nil {
		t.Fatal("GrayF32Model is nil")
	}

	// Test conversion from GrayF32 returns same value
	original := GrayF32{Y: 0.75}
	converted := GrayF32Model.Convert(original)
	if grayF32, ok := converted.(GrayF32); !ok {
		t.Errorf("GrayF32Model.Convert(GrayF32) returned type %T, want GrayF32", converted)
	} else if grayF32.Y != original.Y {
		t.Errorf("GrayF32Model.Convert(GrayF32{%f}) = GrayF32{%f}, want GrayF32{%f}",
			original.Y, grayF32.Y, original.Y)
	}
}

func TestGrayF32Model_ConvertFromRGBA(t *testing.T) {
	tests := []struct {
		name  string
		input color.RGBA
		want  float32
	}{
		{
			name:  "white",
			input: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			want:  1.0,
		},
		{
			name:  "black",
			input: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			want:  0.0,
		},
		{
			name:  "medium gray (128)",
			input: color.RGBA{R: 128, G: 128, B: 128, A: 255},
			want:  0.502, // Approximately 128/255 scaled up to 16-bit then to [0,1]
		},
		{
			name:  "red",
			input: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			want:  0.299, // Red coefficient
		},
		{
			name:  "green",
			input: color.RGBA{R: 0, G: 255, B: 0, A: 255},
			want:  0.587, // Green coefficient
		},
		{
			name:  "blue",
			input: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			want:  0.114, // Blue coefficient
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GrayF32Model.Convert(tt.input)
			grayF32, ok := result.(GrayF32)
			if !ok {
				t.Fatalf("GrayF32Model.Convert returned type %T, want GrayF32", result)
			}
			// Allow for small floating point differences
			diff := math.Abs(float64(grayF32.Y - tt.want))
			if diff > 0.001 {
				t.Errorf("GrayF32Model.Convert(%+v) = GrayF32{%f}, want approximately GrayF32{%f}",
					tt.input, grayF32.Y, tt.want)
			}
		})
	}
}

func TestGrayF32Model_ConvertFromGray16(t *testing.T) {
	tests := []struct {
		name  string
		input color.Gray16
		want  float32
	}{
		{
			name:  "maximum gray16",
			input: color.Gray16{Y: 65535},
			want:  1.0,
		},
		{
			name:  "minimum gray16",
			input: color.Gray16{Y: 0},
			want:  0.0,
		},
		{
			name:  "middle gray16",
			input: color.Gray16{Y: 32768},
			want:  32768.0 / 65535.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GrayF32Model.Convert(tt.input)
			grayF32, ok := result.(GrayF32)
			if !ok {
				t.Fatalf("GrayF32Model.Convert returned type %T, want GrayF32", result)
			}
			diff := math.Abs(float64(grayF32.Y - tt.want))
			if diff > 0.0001 {
				t.Errorf("GrayF32Model.Convert(Gray16{%d}) = GrayF32{%f}, want GrayF32{%f}",
					tt.input.Y, grayF32.Y, tt.want)
			}
		})
	}
}

func TestGrayF32_Implements_Color(t *testing.T) {
	// Compile-time check that GrayF32 implements color.Color
	var _ color.Color = GrayF32{}
}

func TestGrayF32Model_Implements_Model(t *testing.T) {
	// Compile-time check that GrayF32Model implements color.Model
	var _ color.Model = GrayF32Model
}

func TestNewGrayF32Image(t *testing.T) {
	r := image.Rect(0, 0, 10, 10)
	img := NewGrayF32Image(r)

	if img == nil {
		t.Fatal("NewGrayF32Image returned nil")
	}

	if img.Bounds() != r {
		t.Errorf("Bounds() = %v, want %v", img.Bounds(), r)
	}

	if img.Stride != 40 {
		t.Errorf("Stride = %d, want 40", img.Stride)
	}

	expectedLen := 4 * 10 * 10
	if len(img.Pix) != expectedLen {
		t.Errorf("len(Pix) = %d, want %d", len(img.Pix), expectedLen)
	}
}

func TestGrayF32Image_Implements_Image(t *testing.T) {
	// Compile-time check that GrayF32Image implements image.Image
	var _ image.Image = &GrayF32Image{}
}

func TestGrayF32Image_ColorModel(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))
	if img.ColorModel() != GrayF32Model {
		t.Errorf("ColorModel() returned %v, want GrayF32Model", img.ColorModel())
	}
}

func TestGrayF32Image_Bounds(t *testing.T) {
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
			img := NewGrayF32Image(tt.rect)
			if img.Bounds() != tt.rect {
				t.Errorf("Bounds() = %v, want %v", img.Bounds(), tt.rect)
			}
		})
	}
}

func TestGrayF32Image_SetAndGet(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

	tests := []struct {
		name  string
		x, y  int
		color GrayF32
	}{
		{"zero value", 0, 0, GrayF32{Y: 0.0}},
		{"minimum value", 1, 1, GrayF32{Y: 0.0}},
		{"maximum value", 2, 2, GrayF32{Y: 1.0}},
		{"middle value", 3, 3, GrayF32{Y: 0.5}},
		{"quarter value", 4, 4, GrayF32{Y: 0.25}},
		{"special float value", 5, 5, GrayF32{Y: 0.123456}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img.SetGrayF32(tt.x, tt.y, tt.color)
			got := img.GrayF32At(tt.x, tt.y)
			if got.Y != tt.color.Y {
				t.Errorf("After SetGrayF32(%d, %d, GrayF32{%f}), GrayF32At(%d, %d) = GrayF32{%f}, want GrayF32{%f}",
					tt.x, tt.y, tt.color.Y, tt.x, tt.y, got.Y, tt.color.Y)
			}
		})
	}
}

func TestGrayF32Image_Set(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

	// Test setting with color.Color interface
	c := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	img.Set(5, 5, c)

	got := img.At(5, 5)
	if _, ok := got.(GrayF32); !ok {
		t.Errorf("At() returned type %T, want GrayF32", got)
	}
}

func TestGrayF32Image_At(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

	// Set a value and verify At returns it via color.Color interface
	expected := GrayF32{Y: 0.75}
	img.SetGrayF32(5, 5, expected)

	got := img.At(5, 5)
	grayF32, ok := got.(GrayF32)
	if !ok {
		t.Fatalf("At() returned type %T, want GrayF32", got)
	}
	if grayF32.Y != expected.Y {
		t.Errorf("At(5, 5) = GrayF32{%f}, want GrayF32{%f}", grayF32.Y, expected.Y)
	}
}

func TestGrayF32Image_GrayF32At_OutOfBounds(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

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
			got := img.GrayF32At(tt.x, tt.y)
			if got.Y != 0.0 {
				t.Errorf("GrayF32At(%d, %d) = GrayF32{%f}, want GrayF32{0.0} for out of bounds",
					tt.x, tt.y, got.Y)
			}
		})
	}
}

func TestGrayF32Image_SetGrayF32_OutOfBounds(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

	// Setting out of bounds should not panic
	img.SetGrayF32(-1, 5, GrayF32{Y: 0.5})
	img.SetGrayF32(5, -1, GrayF32{Y: 0.5})
	img.SetGrayF32(10, 5, GrayF32{Y: 0.5})
	img.SetGrayF32(5, 10, GrayF32{Y: 0.5})
}

func TestGrayF32Image_PixOffset(t *testing.T) {
	tests := []struct {
		name   string
		rect   image.Rectangle
		x, y   int
		offset int
	}{
		{"zero origin (0,0)", image.Rect(0, 0, 10, 10), 0, 0, 0},
		{"zero origin (1,0)", image.Rect(0, 0, 10, 10), 1, 0, 4},
		{"zero origin (0,1)", image.Rect(0, 0, 10, 10), 0, 1, 40},
		{"zero origin (5,5)", image.Rect(0, 0, 10, 10), 5, 5, 220},
		{"non-zero origin (5,5)", image.Rect(5, 5, 15, 15), 5, 5, 0},
		{"non-zero origin (6,5)", image.Rect(5, 5, 15, 15), 6, 5, 4},
		{"non-zero origin (5,6)", image.Rect(5, 5, 15, 15), 5, 6, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := NewGrayF32Image(tt.rect)
			got := img.PixOffset(tt.x, tt.y)
			if got != tt.offset {
				t.Errorf("PixOffset(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.offset)
			}
		})
	}
}

func TestGrayF32Image_SubImage(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

	// Set some pixels in the original image
	img.SetGrayF32(5, 5, GrayF32{Y: 0.75})
	img.SetGrayF32(6, 6, GrayF32{Y: 0.25})

	// Create a sub-image
	sub := img.SubImage(image.Rect(5, 5, 8, 8))
	subImg, ok := sub.(*GrayF32Image)
	if !ok {
		t.Fatalf("SubImage returned type %T, want *GrayF32Image", sub)
	}

	// Verify bounds
	expectedBounds := image.Rect(5, 5, 8, 8)
	if subImg.Bounds() != expectedBounds {
		t.Errorf("SubImage bounds = %v, want %v", subImg.Bounds(), expectedBounds)
	}

	// Verify the sub-image shares pixels with original
	got := subImg.GrayF32At(5, 5)
	if got.Y != 0.75 {
		t.Errorf("SubImage.GrayF32At(5, 5) = GrayF32{%f}, want GrayF32{0.75}", got.Y)
	}

	// Modify the sub-image and verify it affects the original
	subImg.SetGrayF32(6, 6, GrayF32{Y: 0.9})
	got = img.GrayF32At(6, 6)
	if got.Y != 0.9 {
		t.Errorf("After modifying SubImage, original GrayF32At(6, 6) = GrayF32{%f}, want GrayF32{0.9}", got.Y)
	}
}

func TestGrayF32Image_SubImage_Empty(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

	// Create an empty sub-image
	sub := img.SubImage(image.Rect(5, 5, 5, 5))
	subImg, ok := sub.(*GrayF32Image)
	if !ok {
		t.Fatalf("SubImage returned type %T, want *GrayF32Image", sub)
	}

	if !subImg.Bounds().Empty() {
		t.Errorf("Empty SubImage bounds = %v, want empty rectangle", subImg.Bounds())
	}
}

func TestGrayF32Image_SubImage_NonIntersecting(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

	// Create a non-intersecting sub-image
	sub := img.SubImage(image.Rect(20, 20, 30, 30))
	subImg := sub.(*GrayF32Image)

	if !subImg.Bounds().Empty() {
		t.Errorf("Non-intersecting SubImage bounds = %v, want empty rectangle", subImg.Bounds())
	}
}

func TestGrayF32Image_Opaque(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 10, 10))

	if !img.Opaque() {
		t.Error("Opaque() = false, want true")
	}
}

func TestGrayF32Image_NonZeroOrigin(t *testing.T) {
	// Test with non-zero origin
	img := NewGrayF32Image(image.Rect(5, 5, 15, 15))

	// Set and get a pixel
	expected := GrayF32{Y: 0.123456}
	img.SetGrayF32(7, 8, expected)

	got := img.GrayF32At(7, 8)
	if got.Y != expected.Y {
		t.Errorf("GrayF32At(7, 8) = GrayF32{%f}, want GrayF32{%f}", got.Y, expected.Y)
	}

	// Verify out of bounds below origin
	got = img.GrayF32At(4, 4)
	if got.Y != 0.0 {
		t.Errorf("GrayF32At(4, 4) = GrayF32{%f}, want GrayF32{0.0} (out of bounds)", got.Y)
	}
}

func TestGrayF32Image_BigEndianFloatEncoding(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 1, 1))

	// Test encoding of 1.0 as IEEE 754 float32
	img.SetGrayF32(0, 0, GrayF32{Y: 1.0})

	// 1.0 as float32 is 0x3F800000 in IEEE 754
	expected := []uint8{0x3F, 0x80, 0x00, 0x00}
	if img.Pix[0] != expected[0] || img.Pix[1] != expected[1] ||
		img.Pix[2] != expected[2] || img.Pix[3] != expected[3] {
		t.Errorf("Float encoding: Pix = [0x%02x, 0x%02x, 0x%02x, 0x%02x], want [0x%02x, 0x%02x, 0x%02x, 0x%02x]",
			img.Pix[0], img.Pix[1], img.Pix[2], img.Pix[3],
			expected[0], expected[1], expected[2], expected[3])
	}

	// Test reading back
	got := img.GrayF32At(0, 0)
	if got.Y != 1.0 {
		t.Errorf("GrayF32At(0, 0) = GrayF32{%f}, want GrayF32{1.0}", got.Y)
	}
}

func TestGrayF32Image_SpecialFloatValues(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 4, 1))

	// Test special float values
	tests := []struct {
		name string
		x    int
		val  float32
	}{
		{"zero", 0, 0.0},
		{"negative zero", 1, float32(math.Copysign(0, -1))},
		{"infinity", 2, float32(math.Inf(1))},
		{"negative infinity", 3, float32(math.Inf(-1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img.SetGrayF32(tt.x, 0, GrayF32{Y: tt.val})
			got := img.GrayF32At(tt.x, 0)

			// Special handling for NaN and infinities
			if math.IsInf(float64(tt.val), 0) {
				if !math.IsInf(float64(got.Y), int(math.Copysign(1, float64(tt.val)))) {
					t.Errorf("GrayF32At(%d, 0) = GrayF32{%f}, want infinity with same sign as %f",
						tt.x, got.Y, tt.val)
				}
			} else if got.Y != tt.val || math.Signbit(float64(got.Y)) != math.Signbit(float64(tt.val)) {
				t.Errorf("GrayF32At(%d, 0) = GrayF32{%f}, want GrayF32{%f}",
					tt.x, got.Y, tt.val)
			}
		})
	}
}

func TestGrayF32Image_NaNHandling(t *testing.T) {
	img := NewGrayF32Image(image.Rect(0, 0, 1, 1))

	// Test NaN handling
	nanVal := float32(math.NaN())
	img.SetGrayF32(0, 0, GrayF32{Y: nanVal})
	got := img.GrayF32At(0, 0)

	if !math.IsNaN(float64(got.Y)) {
		t.Errorf("GrayF32At(0, 0) = GrayF32{%f}, want NaN", got.Y)
	}
}
