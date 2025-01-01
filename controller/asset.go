package controller

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"net/http"
	"os"
)

type AssetMeta struct {
	Path     string
	MimeType string
}

func (c *Controller) GetAssetMeta(path string) (AssetMeta, error) {
	file, err := os.Open(path)
	if err != nil {
		println("error open asset file: " + err.Error())
		return AssetMeta{}, fmt.Errorf("error open asset file: %w", err)
	}
	defer file.Close()

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return AssetMeta{}, fmt.Errorf("error detecting mimetype: %w", err)
	}

	result := AssetMeta{
		Path:     path,
		MimeType: mime.String(),
	}

	return result, nil
}

func (c *Controller) handleAsset(resp http.ResponseWriter, req *http.Request) {
	file, err := os.Open(req.URL.Path)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	resp.WriteHeader(http.StatusOK)
	io.Copy(resp, file)
}
