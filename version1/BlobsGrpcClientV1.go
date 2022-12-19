package version1

import (
	"context"
	"encoding/base64"
	"io"

	"github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/pip-services3-gox/pip-services3-commons-gox/data"
	"github.com/pip-services3-gox/pip-services3-grpc-gox/clients"
	"github.com/service-infrastructure2/client-blobs-go/protos"
)

type BlobGrpcClientV1 struct {
	*clients.GrpcClient
	chunkSize int
}

func NewBlobGrpcClientV1() *BlobGrpcClientV1 {
	return &BlobGrpcClientV1{
		GrpcClient: clients.NewGrpcClient("blobs_v1.Blobs"),
		chunkSize:  10240,
	}
}

func (c *BlobGrpcClientV1) Configure(ctx context.Context, config *config.ConfigParams) {
	c.GrpcClient.Configure(ctx, config)

	c.chunkSize = config.GetAsIntegerWithDefault("options.chunk_size", c.chunkSize)
}

func (c *BlobGrpcClientV1) GetBlobsByFilter(ctx context.Context, correlationId string, filter *data.FilterParams,
	paging *data.PagingParams) (result data.DataPage[*BlobInfoV1], err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.get_blobs_by_filter")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobInfoPageRequest{
		CorrelationId: correlationId,
	}
	if filter != nil {
		req.Filter = filter.Value()
	}
	if paging != nil {
		req.Paging = &protos.PagingParams{
			Skip:  paging.GetSkip(0),
			Take:  (int32)(paging.GetTake(100)),
			Total: paging.Total,
		}
	}

	reply := new(protos.BlobInfoPageReply)
	err = c.CallWithContext(ctx, "get_blobs_by_filter", correlationId, req, reply)
	if err != nil {
		return *data.NewEmptyDataPage[*BlobInfoV1](), err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return *data.NewEmptyDataPage[*BlobInfoV1](), err
	}

	result = toBlobInfoPage(reply.Page)

	return result, nil
}

func (c *BlobGrpcClientV1) GetBlobsByIds(ctx context.Context, correlationId string, blobIds []string) (result []*BlobInfoV1, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.get_blobs_by_ids")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobIdsRequest{
		CorrelationId: correlationId,
		BlobIds:       blobIds,
	}

	reply := new(protos.BlobInfoObjectsReply)
	err = c.CallWithContext(ctx, "get_blobs_by_ids", correlationId, req, reply)
	if err != nil {
		return nil, err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return nil, err
	}

	result = toBlobInfos(reply.Blobs)

	return result, nil
}

func (c *BlobGrpcClientV1) GetBlobById(ctx context.Context, correlationId string, blobId string) (result *BlobInfoV1, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.get_blob_by_id")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobIdRequest{
		CorrelationId: correlationId,
		BlobId:        blobId,
	}

	reply := new(protos.BlobInfoObjectReply)
	err = c.CallWithContext(ctx, "get_blob_by_id", correlationId, req, reply)
	if err != nil {
		return nil, err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return nil, err
	}

	result = toBlobInfo(reply.Blob)

	return result, nil
}

func (c *BlobGrpcClientV1) CreateBlobFromUri(ctx context.Context, correlationId string, blob *BlobInfoV1,
	uri string) (result *BlobInfoV1, err error) {
	return BlobsUriProcessorV1.CreateBlobFromUri(ctx, correlationId, blob, c, uri, c.chunkSize)
}

func (c *BlobGrpcClientV1) GetBlobUriById(ctx context.Context, correlationId string, blobId string) (result string, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.get_blob_uri_by_id")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobIdRequest{
		CorrelationId: correlationId,
		BlobId:        blobId,
	}

	reply := new(protos.BlobUriReply)
	err = c.CallWithContext(ctx, "get_blob_uri_by_id", correlationId, req, reply)
	if err != nil {
		return "", err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return "", err
	}

	result = reply.Uri

	return result, nil
}

func (c *BlobGrpcClientV1) CreateBlobFromData(ctx context.Context, correlationId string, blob *BlobInfoV1,
	buffer []byte) (*BlobInfoV1, error) {
	return BlobsDataProcessorV1.CreateBlobFromData(ctx, correlationId, blob, c, buffer, c.chunkSize)
}

func (c *BlobGrpcClientV1) GetBlobDataById(ctx context.Context, correlationId string,
	blobId string) ([]byte, *BlobInfoV1, error) {
	return BlobsDataProcessorV1.GetBlobDataById(ctx, correlationId, blobId, c, c.chunkSize)
}

func (c *BlobGrpcClientV1) CreateBlobFromStream(ctx context.Context, correlationId string, blob *BlobInfoV1,
	stream io.Reader) (*BlobInfoV1, error) {
	return BlobsStreamProcessorV1.CreateBlobFromStream(ctx, correlationId, blob, c, stream, c.chunkSize)
}

func (c *BlobGrpcClientV1) ReadBlobStreamById(ctx context.Context, correlationId string, blobId string,
	stream io.Writer) (*BlobInfoV1, error) {
	return BlobsStreamProcessorV1.GetBlobStreamById(ctx, correlationId, blobId, c, stream, c.chunkSize)
}

func (c *BlobGrpcClientV1) UpdateBlobInfo(ctx context.Context, correlationId string, blob *BlobInfoV1) (result *BlobInfoV1, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.update_blob_info")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobInfoObjectRequest{
		CorrelationId: correlationId,
		Blob:          fromBlobInfo(blob),
	}

	reply := new(protos.BlobInfoObjectReply)
	err = c.CallWithContext(ctx, "update_blob_info", correlationId, req, reply)
	if err != nil {
		return nil, err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return nil, err
	}

	result = toBlobInfo(reply.Blob)

	return result, nil
}

func (c *BlobGrpcClientV1) MarkBlobsCompleted(ctx context.Context, correlationId string, blobIds []string) (err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.mark_blobs_completed")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobIdsRequest{
		CorrelationId: correlationId,
		BlobIds:       blobIds,
	}

	reply := new(protos.BlobEmptyReply)
	err = c.CallWithContext(ctx, "mark_blobs_completed", correlationId, req, reply)
	if err != nil {
		return err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return err
	}

	return nil
}

func (c *BlobGrpcClientV1) DeleteBlobById(ctx context.Context, correlationId string, blobId string) (err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.delete_blob_by_id")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobIdRequest{
		CorrelationId: correlationId,
		BlobId:        blobId,
	}

	reply := new(protos.BlobEmptyReply)
	err = c.CallWithContext(ctx, "delete_blob_by_id", correlationId, req, reply)
	if err != nil {
		return err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return err
	}

	return nil
}

func (c *BlobGrpcClientV1) DeleteBlobsByIds(ctx context.Context, correlationId string, blobIds []string) (err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.delete_blobs_by_ids")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobIdsRequest{
		CorrelationId: correlationId,
		BlobIds:       blobIds,
	}

	reply := new(protos.BlobEmptyReply)
	err = c.CallWithContext(ctx, "delete_blobs_by_ids", correlationId, req, reply)
	if err != nil {
		return err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return err
	}

	return nil
}

func (c *BlobGrpcClientV1) BeginBlobRead(ctx context.Context, correlationId string, blobId string) (result *BlobInfoV1, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.begin_blob_read")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobIdRequest{
		CorrelationId: correlationId,
		BlobId:        blobId,
	}

	reply := new(protos.BlobInfoObjectReply)
	err = c.CallWithContext(ctx, "begin_blob_read", correlationId, req, reply)
	if err != nil {
		return nil, err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return nil, err
	}

	result = toBlobInfo(reply.Blob)

	return result, nil
}

func (c *BlobGrpcClientV1) ReadBlobChunk(ctx context.Context, correlationId string, blobId string, skip int64, take int64) (result []byte, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.read_blob_chunk")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobReadRequest{
		CorrelationId: correlationId,
		BlobId:        blobId,
		Skip:          skip,
		Take:          take,
	}

	reply := new(protos.BlobChunkReply)
	err = c.CallWithContext(ctx, "read_blob_chunk", correlationId, req, reply)
	if err != nil {
		return nil, err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return nil, err
	}

	result, err1 := base64.StdEncoding.DecodeString(reply.Chunk)

	return result, err1
}

func (c *BlobGrpcClientV1) EndBlobRead(ctx context.Context, correlationId string, blobId string) (err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.end_blob_read")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobIdRequest{
		CorrelationId: correlationId,
		BlobId:        blobId,
	}

	reply := new(protos.BlobEmptyReply)
	err = c.CallWithContext(ctx, "end_blob_read", correlationId, req, reply)
	if err != nil {
		return err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return err
	}

	return nil
}

func (c *BlobGrpcClientV1) BeginBlobWrite(ctx context.Context, correlationId string, blob *BlobInfoV1) (result string, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.begin_blob_write")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobInfoObjectRequest{
		CorrelationId: correlationId,
		Blob:          fromBlobInfo(blob),
	}

	reply := new(protos.BlobTokenReply)
	err = c.CallWithContext(ctx, "begin_blob_write", correlationId, req, reply)
	if err != nil {
		return "", err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return "", err
	}

	result = reply.Token

	return result, nil
}

func (c *BlobGrpcClientV1) WriteBlobChunk(ctx context.Context, correlationId string, token string, chunk []byte) (result string, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.write_blob_chunk")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobTokenWithChunkRequest{
		CorrelationId: correlationId,
		Token:         token,
		Chunk:         base64.StdEncoding.EncodeToString(chunk),
	}

	reply := new(protos.BlobTokenReply)
	err = c.CallWithContext(ctx, "write_blob_chunk", correlationId, req, reply)
	if err != nil {
		return "", err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return "", err
	}

	result = reply.Token

	return result, nil
}

func (c *BlobGrpcClientV1) EndBlobWrite(ctx context.Context, correlationId string, token string, chunk []byte) (result *BlobInfoV1, err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.end_blob_write")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobTokenWithChunkRequest{
		CorrelationId: correlationId,
		Token:         token,
		Chunk:         base64.StdEncoding.EncodeToString(chunk),
	}

	reply := new(protos.BlobInfoObjectReply)
	err = c.CallWithContext(ctx, "end_blob_write", correlationId, req, reply)
	if err != nil {
		return nil, err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return nil, err
	}

	result = toBlobInfo(reply.Blob)

	return result, nil

}

func (c *BlobGrpcClientV1) AbortBlobWrite(ctx context.Context, correlationId string, token string) (err error) {
	timing := c.Instrument(ctx, correlationId, "blobs_v1.abort_blob_write")
	defer timing.EndTiming(ctx, err)

	req := &protos.BlobTokenRequest{
		CorrelationId: correlationId,
		Token:         token,
	}

	reply := new(protos.BlobEmptyReply)
	err = c.CallWithContext(ctx, "abort_blob_write", correlationId, req, reply)
	if err != nil {
		return err
	}

	if reply.Error != nil {
		err = toError(reply.Error)
		return err
	}

	return nil
}
