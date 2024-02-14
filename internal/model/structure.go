package model

type Model struct {
	Name          string
	Config        Config
	PipelineTag   string
	DirectoryPath string
	AddToBinary   bool
}

type Config struct {
	ModuleName string `json:"module_name"`
	ClassName  string `json:"class_name"`
}
