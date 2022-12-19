package version1

import (
	"context"
	"net/http"
)

type TBlobsUriProcessorV1 struct{}

var BlobsUriProcessorV1 = &TBlobsUriProcessorV1{}

func (c *TBlobsUriProcessorV1) CreateBlobFromUri(ctx context.Context, correlationId string, blob *BlobInfoV1,
	writer IBlobsChunkyWriterV1, uri string, chunkSize int) (result *BlobInfoV1, err error) {

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	return BlobsStreamProcessorV1.CreateBlobFromStream(ctx, correlationId, blob, writer, resp.Body, chunkSize)
}
