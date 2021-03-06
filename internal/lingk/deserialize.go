package lingk

import (
	"embed"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"mime/multipart"
)

//go:embed sample
var sample embed.FS

func FromAttachment(fh *multipart.FileHeader) (*data.Data, error) {
	t, err := Unpack(fh)
	if err != nil {
		return nil, err
	}

	d, _ := t.prune()

	return d, nil
}

func Sample() (*data.Data, error) {
	t, err := unpackFs(sample, "sample/fa2021")
	if err != nil {
		return nil, err
	}

	d, _ := t.prune()
	return d, nil
}
