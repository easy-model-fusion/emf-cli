package downloader

import "github.com/easy-model-fusion/emf-cli/internal/utils/python"

type Downloader interface {
	Execute(downloaderArgs Args, python python.Python) (Model, error)
}
