package version1

import (
	"context"
	"io"

	"github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

type IBlobsClientV1 interface {
	GetBlobsByFilter(ctx context.Context, correlationId string, filter *data.FilterParams,
		paging *data.PagingParams) (result data.DataPage[*BlobInfoV1], err error)

	GetBlobsByIds(ctx context.Context, correlationId string, blobIds []string) (result []*BlobInfoV1, err error)

	GetBlobById(ctx context.Context, correlationId string, blobId string) (result *BlobInfoV1, err error)

	CreateBlobFromUri(ctx context.Context, correlationId string, blob *BlobInfoV1,
		uri string) (result *BlobInfoV1, err error)

	GetBlobUriById(ctx context.Context, correlationId string, blobId string) (result string, err error)

	CreateBlobFromData(ctx context.Context, correlationId string, blob *BlobInfoV1,
		buffer []byte) (result *BlobInfoV1, err error)

	GetBlobDataById(ctx context.Context, correlationId string,
		blobId string) (result []byte, blob *BlobInfoV1, err error)

	CreateBlobFromStream(ctx context.Context, correlationId string, blob *BlobInfoV1,
		stream io.Reader) (result *BlobInfoV1, err error)

	ReadBlobStreamById(ctx context.Context, correlationId string, blobId string,
		stream io.Writer) (blob *BlobInfoV1, err error)

	UpdateBlobInfo(ctx context.Context, correlationId string, blob *BlobInfoV1) (result *BlobInfoV1, err error)

	MarkBlobsCompleted(ctx context.Context, correlationId string, blobIds []string) error

	DeleteBlobById(ctx context.Context, correlationId string, blobId string) error

	DeleteBlobsByIds(ctx context.Context, correlationId string, blobIds []string) error
}
