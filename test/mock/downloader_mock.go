package mock

import (
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
)

type MockDownloader struct {
	DownloaderModel downloadermodel.Model
	DownloaderError error
}

func (d *MockDownloader) Execute(_ downloadermodel.Args, _ python.Python) (downloadermodel.Model, error) {
	return d.DownloaderModel, d.DownloaderError
}
