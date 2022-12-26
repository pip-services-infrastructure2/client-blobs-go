package version1

import (
	"context"
	"io"

	"github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

type BlobsNullClientV1 struct {
}

func NewBlobsNullClientV1() *BlobsNullClientV1 {
	return &BlobsNullClientV1{}
}

func (c *BlobsNullClientV1) GetBlobsByFilter(ctx context.Context, correlationId string, filter *data.FilterParams, paging *data.PagingParams) (result data.DataPage[*BlobInfoV1], err error) {
	return *data.NewEmptyDataPage[*BlobInfoV1](), nil
}

func (c *BlobsNullClientV1) GetBlobsByIds(ctx context.Context, correlationId string, blobIds []string) (result []*BlobInfoV1, err error) {
	return nil, nil
}

func (c *BlobsNullClientV1) GetBlobById(ctx context.Context, correlationId string, blobId string) (result *BlobInfoV1, err error) {
	return nil, nil
}

func (c *BlobsNullClientV1) CreateBlobFromUri(ctx context.Context, correlationId string, blob *BlobInfoV1, uri string) (result *BlobInfoV1, err error) {
	return nil, nil
}

func (c *BlobsNullClientV1) GetBlobUriById(ctx context.Context, correlationId string, blobId string) (result string, err error) {
	return result, nil
}

func (c *BlobsNullClientV1) CreateBlobFromData(ctx context.Context, correlationId string, blob *BlobInfoV1, buffer []byte) (result *BlobInfoV1, err error) {
	return nil, nil
}

func (c *BlobsNullClientV1) GetBlobDataById(ctx context.Context, correlationId string, blobId string) (result []byte, blob *BlobInfoV1, err error) {
	return make([]byte, 0), nil, nil
}

func (c *BlobsNullClientV1) CreateBlobFromStream(ctx context.Context, correlationId string, blob *BlobInfoV1, stream io.Reader) (result *BlobInfoV1, err error) {
	return nil, nil
}

func (c *BlobsNullClientV1) ReadBlobStreamById(ctx context.Context, correlationId string, blobId string, stream io.Writer) (blob *BlobInfoV1, err error) {
	return nil, nil
}

func (c *BlobsNullClientV1) UpdateBlobInfo(ctx context.Context, correlationId string, blob *BlobInfoV1) (result *BlobInfoV1, err error) {
	return blob, nil
}

func (c *BlobsNullClientV1) MarkBlobsCompleted(ctx context.Context, correlationId string, blobIds []string) error {
	return nil
}

func (c *BlobsNullClientV1) DeleteBlobById(ctx context.Context, correlationId string, blobId string) error {
	return nil
}

func (c *BlobsNullClientV1) DeleteBlobsByIds(ctx context.Context, correlationId string, blobIds []string) error {
	return nil
}

func (c *BlobsNullClientV1) BeginBlobRead(ctx context.Context, correlationId string, blobId string) (blob *BlobInfoV1, err error) {
	return nil, nil
}

func (c *BlobsNullClientV1) ReadBlobChunk(ctx context.Context, correlationId string, blobId string, skip int64, take int64) (chunk []byte, err error) {
	return nil, nil
}

func (c *BlobsNullClientV1) EndBlobRead(ctx context.Context, correlationId string, blobId string) error {
	return nil
}

func (c *BlobsNullClientV1) BeginBlobWrite(ctx context.Context, correlationId string, blob *BlobInfoV1) (token string, err error) {
	return token, nil
}

func (c *BlobsNullClientV1) WriteBlobChunk(ctx context.Context, correlationId string, token string, chunk []byte) (token2 string, err error) {
	return token, nil
}

func (c *BlobsNullClientV1) EndBlobWrite(ctx context.Context, correlationId string, token string, chunk []byte) (blob *BlobInfoV1, err error) {
	return blob, nil
}

func (c *BlobsNullClientV1) AbortBlobWrite(ctx context.Context, correlationId string, token string) error {
	return nil
}
