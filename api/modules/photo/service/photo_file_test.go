package service

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndResizeFile_PNGPreservesTransparency(t *testing.T) {
	fileHeader := newMultipartFileHeader(t, "logo.png", buildPNG(t))
	basePath := t.TempDir()

	if err := SaveAndResizeFile(fileHeader, "logo.png", basePath); err != nil {
		t.Fatalf("SaveAndResizeFile() error = %v", err)
	}

	assertPNGWithTransparentPixel(t, filepath.Join(basePath, "original", "logo.png"))
	assertPNGWithTransparentPixel(t, filepath.Join(basePath, "medium", "logo.png"))
	assertPNGWithTransparentPixel(t, filepath.Join(basePath, "thumbnail", "logo.png"))
}

func TestSaveAndResizeFile_JPEGStaysJPEG(t *testing.T) {
	fileHeader := newMultipartFileHeader(t, "logo.jpg", buildJPEG(t))
	basePath := t.TempDir()

	if err := SaveAndResizeFile(fileHeader, "logo.jpg", basePath); err != nil {
		t.Fatalf("SaveAndResizeFile() error = %v", err)
	}

	assertJPEG(t, filepath.Join(basePath, "original", "logo.jpg"))
	assertJPEG(t, filepath.Join(basePath, "medium", "logo.jpg"))
	assertJPEG(t, filepath.Join(basePath, "thumbnail", "logo.jpg"))
}

func newMultipartFileHeader(t *testing.T, filename string, body []byte) *multipart.FileHeader {
	t.Helper()

	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)

	part, err := writer.CreateFormFile("photo", filename)
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := part.Write(body); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}

	req := httptest.NewRequest("POST", "/", &reqBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err := req.ParseMultipartForm(int64(len(reqBody.Bytes()) + 1024)); err != nil {
		t.Fatalf("ParseMultipartForm() error = %v", err)
	}

	file, header, err := req.FormFile("photo")
	if err != nil {
		t.Fatalf("FormFile() error = %v", err)
	}
	file.Close()

	return header
}

func buildPNG(t *testing.T) []byte {
	t.Helper()

	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{R: 255, G: 90, B: 0, A: 255})
	img.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	img.Set(0, 1, color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	img.Set(1, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 0})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode() error = %v", err)
	}
	return buf.Bytes()
}

func buildJPEG(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 90, B: 0, A: 255})
	img.Set(1, 0, color.RGBA{R: 30, G: 30, B: 30, A: 255})
	img.Set(0, 1, color.RGBA{R: 200, G: 200, B: 200, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("jpeg.Encode() error = %v", err)
	}
	return buf.Bytes()
}

func assertPNGWithTransparentPixel(t *testing.T, path string) {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("os.Open(%q) error = %v", path, err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("png.Decode(%q) error = %v", path, err)
	}

	_, _, _, a := img.At(img.Bounds().Dx()-1, 0).RGBA()
	if a != 0 {
		t.Fatalf("expected transparent pixel in %q, got alpha=%d", path, a)
	}
}

func assertJPEG(t *testing.T, path string) {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("os.Open(%q) error = %v", path, err)
	}
	defer f.Close()

	if _, err := jpeg.Decode(f); err != nil {
		t.Fatalf("jpeg.Decode(%q) error = %v", path, err)
	}
}
