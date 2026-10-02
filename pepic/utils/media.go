package utils

import (
	"mime"
	"net/http"
	"path"
)

func DetectMimeType(filename string, data []byte) string {
	mimeType := http.DetectContentType(data)
	if mimeType == "application/octet-stream" {
		if heif := sniffHEIF(data); heif != "" {
			return heif
		}
		mimeType = mime.TypeByExtension(path.Ext(filename))
	}
	return mimeType
}

// sniffHEIF recognises HEIC/HEIF by content: an ISO-BMFF "ftyp" box at byte 4
// followed by a HEIF brand. http.DetectContentType does not know the format,
// and iPhone uploads do not always keep the .heic extension.
func sniffHEIF(data []byte) string {
	if len(data) < 12 || string(data[4:8]) != "ftyp" {
		return ""
	}
	switch string(data[8:12]) {
	case "heic", "heix", "heim", "heis", "hevc", "hevx":
		return "image/heic"
	case "mif1", "msf1", "heif":
		return "image/heif"
	}
	return ""
}

func FitSize(origWidth int, origHeight int, length int) (int, int) {
	if origWidth > origHeight {
		return length, (origHeight * length) / origWidth
	} else if origWidth < origHeight {
		return (origWidth * length) / origHeight, length
	} else {
		return length, length
	}
}
