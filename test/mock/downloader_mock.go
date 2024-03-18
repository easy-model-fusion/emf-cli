package mock

import (
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
)

type MockDownloader struct {
	DownloaderModel downloader.Model
	DownloaderError error
}

func (d *MockDownloader) Execute(_ downloader.Args) (downloader.Model, error) {
	return d.DownloaderModel, d.DownloaderError
}
