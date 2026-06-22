package qrcode

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"regexp"
	"strings"

	"github.com/fogleman/gg"
	"github.com/makiuchi-d/gozxing"
	zxqrcode "github.com/makiuchi-d/gozxing/qrcode"
	skip2 "github.com/skip2/go-qrcode"

	"qrcode/internal/model"
)

var (
	urlRegex  = regexp.MustCompile(`(?i)^(https?|ftp|mailto|tel|weixin|alipays):`)
	dataURLRx = regexp.MustCompile(`^data:image/[a-zA-Z]+;base64,`)
)

func DetectType(content string) string {
	c := strings.TrimSpace(content)
	if urlRegex.MatchString(c) {
		return "url"
	}
	return "text"
}

func mapLevel(level string) skip2.RecoveryLevel {
	switch strings.ToUpper(level) {
	case "L":
		return skip2.Low
	case "Q":
		return skip2.High
	case "H":
		return skip2.Highest
	default:
		return skip2.Medium
	}
}

func parseHexColor(s string) (r, g, b, a int) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "#")
	if len(s) == 3 {
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	}
	if len(s) != 6 {
		return 14, 14, 17, 255
	}
	var v [3]int
	for i := 0; i < 3; i++ {
		var x int
		for j := 0; j < 2; j++ {
			c := s[i*2+j]
			x <<= 4
			switch {
			case c >= '0' && c <= '9':
				x += int(c - '0')
			case c >= 'a' && c <= 'f':
				x += int(c-'a') + 10
			case c >= 'A' && c <= 'F':
				x += int(c-'A') + 10
			}
		}
		v[i] = x
	}
	return v[0], v[1], v[2], 255
}

func decodeLogo(data string) (image.Image, error) {
	if data == "" {
		return nil, nil
	}
	raw := dataURLRx.ReplaceAllString(data, "")
	b, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		b, err = base64.RawStdEncoding.DecodeString(raw)
		if err != nil {
			return nil, err
		}
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	return img, err
}

func Generate(content string, style model.StyleConfig) ([]byte, error) {
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("content is empty")
	}
	style = style.WithDefaults()

	q, err := skip2.New(content, mapLevel(style.EccLevel))
	if err != nil {
		return nil, err
	}
	bitmap := q.Bitmap()
	n := len(bitmap)
	if n == 0 {
		return nil, errors.New("empty qr matrix")
	}

	margin := style.Margin
	if margin < 0 {
		margin = 0
	}
	total := n + margin*2
	if total <= 0 {
		total = n
	}
	px := int(math.Floor(float64(style.Size) / float64(total)))
	if px < 2 {
		px = 2
	}
	dim := px * total

	dc := gg.NewContext(dim, dim)
	br, bg, bb, ba := parseHexColor(style.Background)
	dc.SetRGBA(float64(br)/255, float64(bg)/255, float64(bb)/255, float64(ba)/255)
	dc.Clear()

	fr, fg, fb, fa := parseHexColor(style.Foreground)
	dc.SetRGBA(float64(fr)/255, float64(fg)/255, float64(fb)/255, float64(fa)/255)
	off := float64(margin * px)
	p := float64(px)
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			if !bitmap[y][x] {
				continue
			}
			x0 := off + float64(x)*p
			y0 := off + float64(y)*p
			switch style.CornerStyle {
			case "dot":
				dc.DrawCircle(x0+p/2, y0+p/2, p*0.42)
			case "rounded":
				dc.DrawRoundedRectangle(x0, y0, p, p, p*0.32)
			default:
				dc.DrawRectangle(x0, y0, p, p)
			}
		}
	}
	dc.Fill()

	if logoImg := mustLogo(style.Logo); logoImg != nil {
		logoSize := int(float64(dim) * 0.22)
		lx := (dim - logoSize) / 2
		ly := (dim - logoSize) / 2
		pad := int(float64(logoSize) * 0.12)
		dc.SetRGBA(1, 1, 1, 0.96)
		dc.DrawRoundedRectangle(float64(lx-pad), float64(ly-pad), float64(logoSize+pad*2), float64(logoSize+pad*2), float64(pad))
		dc.Fill()
		dc.DrawImageAnchored(resizeNearest(logoImg, logoSize, logoSize), lx, ly, 0, 0)
	}

	var buf bytes.Buffer
	if err := dc.EncodePNG(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func mustLogo(data string) image.Image {
	img, err := decodeLogo(data)
	if err != nil {
		return nil
	}
	return img
}

func resizeNearest(src image.Image, w, h int) image.Image {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	sx := float64(b.Dx()) / float64(w)
	sy := float64(b.Dy()) / float64(h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			spx := b.Min.X + int(float64(x)*sx)
			spy := b.Min.Y + int(float64(y)*sy)
			dst.Set(x, y, src.At(spx, spy))
		}
	}
	return dst
}

func Parse(img image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	reader := zxqrcode.NewQRCodeReader()
	result, err := reader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	return result.GetText(), nil
}
