package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type ImageInfo struct {
	FileName  string
	FilePath  string
	ThumbPath string
	FileSize  int64
	Width     int
	Height    int
	MimeType  string
	FileType  string
}

func SaveUploadedFile(file *multipart.FileHeader, uploadDir, thumbDir string) (*ImageInfo, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileType := getFileType(ext)
	mimeType := getMimeType(ext)

	uuidName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, uuidName)
	thumbPath := filepath.Join(thumbDir, uuidName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	fileBytes := make([]byte, file.Size)
	src.Read(fileBytes)
	dst.Write(fileBytes)

	fileInfo, _ := dst.Stat()
	fileSize := fileInfo.Size()

	width, height := 0, 0
	if fileType == "image" {
		width, height, err = generateThumbnail(filePath, thumbPath, 300, 300)
		if err != nil {
			return nil, fmt.Errorf("thumbnail generation failed: %v", err)
		}
	}

	return &ImageInfo{
		FileName:  uuidName,
		FilePath:  uuidName,
		ThumbPath: uuidName,
		FileSize:  fileSize,
		Width:     width,
		Height:    height,
		MimeType:  mimeType,
		FileType:  fileType,
	}, nil
}

func generateThumbnail(srcPath, dstPath string, maxWidth, maxHeight uint) (int, int, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	thumb := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)

	out, err := os.Create(dstPath)
	if err != nil {
		return 0, 0, err
	}
	defer out.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(out, thumb, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(out, thumb)
	default:
		err = jpeg.Encode(out, thumb, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func getFileType(ext string) string {
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".bmp": true, ".webp": true, ".tiff": true,
	}
	videoExts := map[string]bool{
		".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
		".webm": true, ".flv": true, ".wmv": true,
	}

	if imageExts[ext] {
		return "image"
	}
	if videoExts[ext] {
		return "video"
	}
	return "other"
}

func getMimeType(ext string) string {
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".webp": "image/webp",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".avi":  "video/x-msvideo",
		".mkv":  "video/x-matroska",
		".webm": "video/webm",
	}
	if m, ok := mimeTypes[ext]; ok {
		return m
	}
	return "application/octet-stream"
}

func CompressImage(srcPath string, quality int) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	if strings.ToLower(format) != "jpeg" && strings.ToLower(format) != "jpg" {
		return nil
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var newImg image.Image
	if width > 1920 || height > 1080 {
		newImg = resize.Thumbnail(1920, 1080, img, resize.Lanczos3)
	} else {
		newImg = img
	}

	out, err := os.Create(srcPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, newImg, &jpeg.Options{Quality: quality})
}
