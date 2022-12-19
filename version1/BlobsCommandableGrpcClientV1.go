package version1

import (
	"context"
	"io"

	"github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/pip-services3-gox/pip-services3-commons-gox/data"
	"github.com/pip-services3-gox/pip-services3-grpc-gox/clients"
)

type BlobsCommandableGrpcClientV1 struct {
	*clients.CommandableGrpcClient
	chunkSize int64
}

func NewBlobsCommandableGrpcClientV1() *BlobsCommandableGrpcClientV1 {
	return NewBlobsCommandableGrpcClientV1WithConfig(nil)
}

func NewBlobsCommandableGrpcClientV1WithConfig(config *cconf.ConfigParams) *BlobsCommandableGrpcClientV1 {
	c := &BlobsCommandableGrpcClientV1{
		CommandableGrpcClient: clients.NewCommandableGrpcClient("v1/blobs"),
		chunkSize:             10240,
	}

	if config != nil {
		c.Configure(context.Background(), config)
	}

	return c
}

func (c *BlobsCommandableGrpcClientV1) Configure(ctx context.Context, config *config.ConfigParams) {
	c.CommandableGrpcClient.Configure(ctx, config)
	c.chunkSize = config.GetAsLongWithDefault("options.chunk_size", c.chunkSize)
}

func (c *BlobsCommandableGrpcClientV1) GetBlobsByFilter(ctx context.Context, correlationId string, filter *data.FilterParams, paging *data.PagingParams) (result data.DataPage[*BlobInfoV1], err error) {
	params := data.NewAnyValueMapFromTuples(
		"filter", filter,
		"paging", paging,
	)

	res, err := c.CallCommand(ctx, "get_blobs_by_filter", correlationId, params)
	if err != nil {
		return *data.NewEmptyDataPage[*BlobInfoV1](), err
	}

	return clients.HandleHttpResponse[data.DataPage[*BlobInfoV1]](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) GetBlobsByIds(ctx context.Context, correlationId string, blobIds []string) (result []*BlobInfoV1, err error) {
	params := data.NewAnyValueMapFromTuples(
		"blob_ids", blobIds,
	)

	res, err := c.CallCommand(ctx, "get_blobs_by_ids", correlationId, params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[[]*BlobInfoV1](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) GetBlobById(ctx context.Context, correlationId string, blobId string) (result *BlobInfoV1, err error) {
	params := data.NewAnyValueMapFromTuples(
		"blob_id", blobId,
	)

	res, err := c.CallCommand(ctx, "get_blob_by_id", correlationId, params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*BlobInfoV1](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) CreateBlobFromUri(ctx context.Context, correlationId string, blob *BlobInfoV1,
	uri string) (result *BlobInfoV1, err error) {
	return BlobsUriProcessorV1.CreateBlobFromUri(ctx, correlationId, blob, c, uri, int(c.chunkSize))
}

func (c *BlobsCommandableGrpcClientV1) GetBlobUriById(ctx context.Context, correlationId string, blobId string) (result string, err error) {
	params := data.NewAnyValueMapFromTuples(
		"blob_id", blobId,
	)

	res, err := c.CallCommand(ctx, "get_blob_uri_by_id", correlationId, params)
	if err != nil {
		return "", err
	}

	return clients.HandleHttpResponse[string](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) CreateBlobFromData(ctx context.Context, correlationId string, blob *BlobInfoV1, buffer []byte) (result *BlobInfoV1, err error) {
	return BlobsDataProcessorV1.CreateBlobFromData(ctx, correlationId, blob, c, buffer, int(c.chunkSize))
}

func (c *BlobsCommandableGrpcClientV1) GetBlobDataById(ctx context.Context, correlationId string, blobId string) (result []byte, blob *BlobInfoV1, err error) {
	return BlobsDataProcessorV1.GetBlobDataById(ctx, correlationId, blobId, c, int(c.chunkSize))
}

func (c *BlobsCommandableGrpcClientV1) CreateBlobFromStream(ctx context.Context, correlationId string, blob *BlobInfoV1, stream io.Reader) (result *BlobInfoV1, err error) {
	return BlobsStreamProcessorV1.CreateBlobFromStream(ctx, correlationId, blob, c, stream, int(c.chunkSize))
}

func (c *BlobsCommandableGrpcClientV1) ReadBlobStreamById(ctx context.Context, correlationId string, blobId string, stream io.Writer) (blob *BlobInfoV1, err error) {
	return BlobsStreamProcessorV1.GetBlobStreamById(ctx, correlationId, blobId, c, stream, int(c.chunkSize))
}

func (c *BlobsCommandableGrpcClientV1) UpdateBlobInfo(ctx context.Context, correlationId string, blob *BlobInfoV1) (result *BlobInfoV1, err error) {
	params := data.NewAnyValueMapFromTuples(
		"blob", blob,
	)

	res, err := c.CallCommand(ctx, "update_blob_info", correlationId, params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*BlobInfoV1](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) MarkBlobsCompleted(ctx context.Context, correlationId string, blobIds []string) error {
	params := data.NewAnyValueMapFromTuples(
		"blobIds", blobIds,
	)

	_, err := c.CallCommand(ctx, "mark_blobs_completed", correlationId, params)

	return err
}

func (c *BlobsCommandableGrpcClientV1) DeleteBlobById(ctx context.Context, correlationId string, blobId string) error {
	params := data.NewAnyValueMapFromTuples(
		"blob_id", blobId,
	)

	_, err := c.CallCommand(ctx, "delete_blob_by_id", correlationId, params)

	return err
}

func (c *BlobsCommandableGrpcClientV1) DeleteBlobsByIds(ctx context.Context, correlationId string, blobIds []string) error {
	params := data.NewAnyValueMapFromTuples(
		"blob_ids", blobIds,
	)

	_, err := c.CallCommand(ctx, "delete_blobs_by_ids", correlationId, params)

	return err
}

func (c *BlobsCommandableGrpcClientV1) BeginBlobWrite(ctx context.Context, correlationId string, blob *BlobInfoV1) (token string, err error) {
	params := data.NewAnyValueMapFromTuples(
		"blob", blob,
	)

	res, err := c.CallCommand(ctx, "begin_blob_write", correlationId, params)
	if err != nil {
		return "", err
	}

	return clients.HandleHttpResponse[string](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) WriteBlobChunk(ctx context.Context, correlationId string, token string, chunk []byte) (token2 string, err error) {
	params := data.NewAnyValueMapFromTuples(
		"token", token,
		"chunk", chunk,
	)

	res, err := c.CallCommand(ctx, "write_blob_chunk", correlationId, params)
	if err != nil {
		return "", err
	}

	return clients.HandleHttpResponse[string](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) EndBlobWrite(ctx context.Context, correlationId string, token string, chunk []byte) (blob *BlobInfoV1, err error) {
	params := data.NewAnyValueMapFromTuples(
		"token", token,
		"chunk", chunk,
	)

	res, err := c.CallCommand(ctx, "end_blob_write", correlationId, params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*BlobInfoV1](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) AbortBlobWrite(ctx context.Context, correlationId string, token string) error {
	params := data.NewAnyValueMapFromTuples(
		"token", token,
	)

	_, err := c.CallCommand(ctx, "abort_blob_write", correlationId, params)

	return err
}

func (c *BlobsCommandableGrpcClientV1) BeginBlobRead(ctx context.Context, correlationId string, blobId string) (blob *BlobInfoV1, err error) {
	params := data.NewAnyValueMapFromTuples(
		"blob_id", blobId,
	)

	res, err := c.CallCommand(ctx, "begin_blob_read", correlationId, params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[*BlobInfoV1](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) ReadBlobChunk(ctx context.Context, correlationId string, blobId string, skip int64, take int64) (chunk []byte, err error) {
	params := data.NewAnyValueMapFromTuples(
		"blob_id", blobId,
		"skip", skip,
		"take", take,
	)

	res, err := c.CallCommand(ctx, "read_blob_chunk", correlationId, params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[[]byte](res, correlationId)
}

func (c *BlobsCommandableGrpcClientV1) EndBlobRead(ctx context.Context, correlationId string, blobId string) error {
	params := data.NewAnyValueMapFromTuples(
		"blob_id", blobId,
	)

	_, err := c.CallCommand(ctx, "end_blob_read", correlationId, params)

	return err
}
