package utils

import "testing"

func TestDetectMimeTypeRecognisesHEIFByContent(t *testing.T) {
	box := func(brand string) []byte {
		return append([]byte{0, 0, 0, 0x1c, 'f', 't', 'y', 'p'}, []byte(brand+"\x00\x00\x00\x00")...)
	}
	cases := []struct {
		name, brand, want string
	}{
		{"IMG_0001", "heic", "image/heic"},  // iPhone, no extension
		{"photo.jpg", "heix", "image/heic"}, // wrong extension: content wins
		{"x", "mif1", "image/heif"},
	}
	for _, c := range cases {
		if got := DetectMimeType(c.name, box(c.brand)); got != c.want {
			t.Errorf("DetectMimeType(%q, %s) = %q, want %q", c.name, c.brand, got, c.want)
		}
	}
	if got := sniffHEIF([]byte("not an image at all")); got != "" {
		t.Errorf("sniffHEIF(garbage) = %q, want empty", got)
	}
}

func TestExtensionByMimeTypeDoesNotPanicOnUnknownTypes(t *testing.T) {
	if ext, err := ExtensionByMimeType("image/heic"); err != nil || ext != ".heic" {
		t.Errorf("image/heic → %q, %v", ext, err)
	}
	if _, err := ExtensionByMimeType("application/x-nothing-knows-this"); err == nil {
		t.Error("want an error for an unknown type, not a panic")
	}
}
