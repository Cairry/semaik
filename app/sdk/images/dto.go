package images

import "time"

type ImageStruct struct {
	Created    time.Time
	ID         string
	Repository string
	RepoTags   string
	Size       int64
}

type ImageOperationStruct struct {
	Name      string `json:"name" validate:"required"`
	Operation string `json:"operation" validate:"required,oneof=start stop restart kill pause unpause rename remove"`
}
