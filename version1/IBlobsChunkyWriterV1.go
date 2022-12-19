package version1

import "context"

type IBlobsChunkyWriterV1 interface {
	BeginBlobWrite(ctx context.Context, correlationId string, blob *BlobInfoV1) (token string, err error)

	WriteBlobChunk(ctx context.Context, correlationId string, token string, chunk []byte) (token2 string, err error)

	EndBlobWrite(ctx context.Context, correlationId string, token string, chunk []byte) (blob *BlobInfoV1, err error)

	AbortBlobWrite(ctx context.Context, correlationId string, token string) error
}
