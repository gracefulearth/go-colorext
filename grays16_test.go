package colorext

import (
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
