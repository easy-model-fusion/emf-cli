package mock

import (
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
)

type MockDownloader struct {
	DownloaderModel downloader.Model
	DownloaderError error
}

func (d *MockDownloader) Execute(_ downloader.Args, _ python.Python) (downloader.Model, error) {
	return d.DownloaderModel, d.DownloaderError
}
