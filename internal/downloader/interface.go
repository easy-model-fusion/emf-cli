package downloader

import (
	downloadermodel "github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
)

type Downloader interface {
	Execute(downloaderArgs downloadermodel.Args, python python.Python) (downloadermodel.Model, error)
}
