package model

type Model struct {
	Name          string `json:"name"`
	PipeLine      string `json:"pipeline"`
	DirectoryPath string `json:"directorypath"`
	AddToBinary   bool   `json:"addtobinary"`
}
