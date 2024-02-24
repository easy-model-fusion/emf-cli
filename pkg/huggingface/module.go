package huggingface

var AllModules = []Module{
	DIFFUSERS,
	TRANSFORMERS,
}

type Module string

const (
	DIFFUSERS    Module = "diffusers"
	TRANSFORMERS Module = "transformers"
)

// AllModulesString returns all modules as a string slice
func AllModulesString() []string {
	var modules []string
	for _, module := range AllModules {
		modules = append(modules, string(module))
	}
	return modules
}