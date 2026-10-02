package entity

import (
	"strings"

	"github.com/vas3k/pepic/pepic/config"
)

type ProcessingFile struct {
	Filename string
	Mime     string
	Path     string
	Bytes    []byte
	Size     int64
}

func (p *ProcessingFile) Url() string {
	return config.App.Global.BaseUrl + p.Filename
}

func (p *ProcessingFile) IsGIF() bool {
	return p.Mime == "image/gif"
}

// IsHEIF: iPhone photos. Browsers other than Safari cannot show them.
func (p *ProcessingFile) IsHEIF() bool {
	return p.Mime == "image/heic" || p.Mime == "image/heif"
}

func (p *ProcessingFile) IsImage() bool {
	return strings.HasPrefix(p.Mime, "image/")
}

func (p *ProcessingFile) IsVideo() bool {
	return strings.HasPrefix(p.Mime, "video/")
}
