package model

var AllTags = []PipelineTag{
	TextToText,
	TextToImage,
}

type PipelineTag string

const (
	TextToText  PipelineTag = "text2text-generation"
	TextToImage PipelineTag = "text-to-image"
)

// AllTagsString returns all tags as a string slice
func AllTagsString() []string {
	var tags []string
	for _, tag := range AllTags {
		tags = append(tags, string(tag))
	}
	return tags
}
