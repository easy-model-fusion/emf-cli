package downloadermodel

// Model represents a model returned by the downloader script.
type Model struct {
	Path      string            `json:"path"`
	Module    string            `json:"module"`
	Class     string            `json:"class"`
	Options   map[string]string `json:"options"`
	Tokenizer Tokenizer         `json:"tokenizer"`
	IsEmpty   bool
}

// Tokenizer represents a tokenizer returned by the downloader script.
type Tokenizer struct {
	Path    string            `json:"path"`
	Class   string            `json:"class"`
	Options map[string]string `json:"options"`
}

// Args represents the arguments for the script.
type Args struct {
	ModelName         string
	ModelModule       string
	ModelClass        string
	ModelOptions      []string
	TokenizerClass    string
	TokenizerOptions  []string
	SkipTokenizer     bool
	SkipModel         bool
	OnlyConfiguration bool
	DirectoryPath     string
}

// Constants related to the downloader script python arguments.
const (
	ScriptPath         = "sdk/downloader.py"
	TagPrefix          = "--"
	ModelName          = "model-name"
	ModelModule        = "model-module"
	Path               = "path"
	ModelClass         = "model-class"
	ModelOptions       = "model-options"
	TokenizerClass     = "tokenizer-class"
	TokenizerOptions   = "tokenizer-options"
	Overwrite          = "overwrite"
	Skip               = "skip"
	SkipValueModel     = "model"
	SkipValueTokenizer = "tokenizer"
	EmfClient          = "emf-client"
	OnlyConfiguration  = "only-configuration"
)
