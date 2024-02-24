package huggingface

var AllTags = []PipelineTag{
	TextGeneration,
	TextToImage,
}

type PipelineTag string

const (
	TextGeneration PipelineTag = "text-generation"
	TextToImage    PipelineTag = "text-to-image"
)

// AllTagsString returns all tags as a string slice
func AllTagsString() []string {
	var tags []string
	for _, tag := range AllTags {
		tags = append(tags, string(tag))
	}
	return tags
}
