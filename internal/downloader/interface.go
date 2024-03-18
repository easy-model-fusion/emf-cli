package downloader

type Downloader interface {
	Execute(downloaderArgs Args) (Model, error)
}
