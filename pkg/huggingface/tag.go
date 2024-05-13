package huggingface

var AllTags = []PipelineTag{
	TextGeneration,
	TextToImage,
	TextToVideo,
	TextToAudio,
	TextTo3d,
	ImageToText,
	ImageToImage,
}

type PipelineTag string

const (
	TextGeneration PipelineTag = "text-generation"
	TextToImage    PipelineTag = "text-to-image"
	TextToVideo    PipelineTag = "text-to-video"
	TextTo3d       PipelineTag = "text-to-3d"
	TextToAudio    PipelineTag = "text-to-audio"
	ImageToText    PipelineTag = "image-to-text"
	ImageToImage   PipelineTag = "image-to-image"
)

// AllTagsString returns all tags as a string slice
func AllTagsString() []string {
	var tags []string
	for _, tag := range AllTags {
		tags = append(tags, string(tag))
	}
	return tags
}
