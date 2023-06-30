package images

import "time"

type ImageStruct struct {
	Created time.Time
	ID      string
	Tags    []string
	Size    string
}

type ImageOperationStruct struct {
	Name      string `json:"name" validate:"required"`
	Operation string `json:"operation" validate:"required,oneof=pull rename remove"`
}
