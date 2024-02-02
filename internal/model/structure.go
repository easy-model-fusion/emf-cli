package model

type Model struct {
	Name          string `json:"modelId"`
	Config        Config `json:"config"`
	PipelineTag   string `json:"pipeline_tag"`
	DirectoryPath string
	AddToBinary   bool
}

type Config struct {
	ModuleName string `json:"module_name"`
	ClassName  string `json:"class_name"`
}
