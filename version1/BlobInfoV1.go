package version1

import (
	"time"
)

type BlobInfoV1 struct {
	/* Identification */
	Id    string `json:"id"`
	Group string `json:"group"`
	Name  string `json:"name"`

	/* Content */
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreateTime  time.Time `json:"create_time"`
	ExpireTime  time.Time `json:"expire_time"`
	Completed   bool      `json:"completed"`
}

func EmptyBlobInfoV1() *BlobInfoV1 {
	return &BlobInfoV1{}
}

func NewBlobInfoV1(id string, group string, name string,
	size int64, contentType string) *BlobInfoV1 {
	return &BlobInfoV1{
		Id:          id,
		Group:       group,
		Name:        name,
		Size:        size,
		ContentType: contentType,
		CreateTime:  time.Now(),
		Completed:   false,
	}
}

func NewBlobInfoV1WithExpiration(id string, group string, name string,
	size int64, contentType string, expireTime time.Time) *BlobInfoV1 {
	return &BlobInfoV1{
		Id:          id,
		Group:       group,
		Name:        name,
		Size:        size,
		ContentType: contentType,
		CreateTime:  time.Now(),
		ExpireTime:  expireTime,
		Completed:   false,
	}
}
