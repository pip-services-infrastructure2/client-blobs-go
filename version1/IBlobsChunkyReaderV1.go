package version1

import "context"

type IBlobsChunkyReaderV1 interface {
	BeginBlobRead(ctx context.Context, correlationId string, blobId string) (blob *BlobInfoV1, err error)

	ReadBlobChunk(ctx context.Context, correlationId string, blobId string, skip int64, take int64) (chunk []byte, err error)

	EndBlobRead(ctx context.Context, correlationId string, blobId string) error
}
