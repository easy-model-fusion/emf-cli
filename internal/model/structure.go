package model

type Model struct {
	Name          string `json:"modelId"`
	Config        Config `json:"config"`
	DirectoryPath string
	AddToBinary   bool
}

type Config struct {
	Diffusers Diffusers `json:"diffusers"`
}

type Diffusers struct {
	PipeLine string `json:"class_name"`
}
