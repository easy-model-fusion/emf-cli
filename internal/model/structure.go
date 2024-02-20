package model

type Model struct {
	Name               string
	Config             Config
	PipelineTag        string
	Source             string
	ShouldBeDownloaded bool
}

type Config struct {
	Path       string
	Module     string
	Class      string
	Tokenizers []Tokenizer
}

type Tokenizer struct {
	Path  string
	Class string
}
