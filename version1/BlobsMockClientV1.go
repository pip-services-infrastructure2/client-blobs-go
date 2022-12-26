package version1

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	"github.com/pip-services3-gox/pip-services3-commons-gox/data"
	"github.com/pip-services3-gox/pip-services3-commons-gox/errors"
)

type BlobsMockClientV1 struct {
	blobs       []*BlobInfoV1
	chunkSize   int64
	maxBlobSize int64
	content     map[string][]byte
}

func NewBlobsMockClientV1() *BlobsMockClientV1 {
	return &BlobsMockClientV1{
		chunkSize:   10240,
		blobs:       make([]*BlobInfoV1, 0),
		maxBlobSize: 100 * 1024,
		content:     make(map[string][]byte, 0),
	}
}

func (c *BlobsMockClientV1) matchString(value string, search string) bool {
	if value == "" && search == "" {
		return true
	}
	if value == "" || search == "" {
		return false
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(search))
}

func (c *BlobsMockClientV1) matchSearch(item *BlobInfoV1, search string) bool {
	search = strings.ToLower(search)
	return c.matchString(item.Name, search)
}

func (c *BlobsMockClientV1) composeFilter(filter *data.FilterParams) func(item *BlobInfoV1) bool {
	if filter == nil {
		filter = data.NewEmptyFilterParams()
	}

	search := filter.GetAsString("search")
	id := filter.GetAsString("id")
	name := filter.GetAsString("name")
	group := filter.GetAsString("group")
	completed, completedOk := filter.GetAsNullableBoolean("completed")
	expired, expiredOk := filter.GetAsNullableBoolean("expired")
	fromCreateTime, fromCreateTimeOK := filter.GetAsNullableDateTime("from_create_time")
	toCreateTime, toCreateTimeOk := filter.GetAsNullableDateTime("to_create_time")

	now := time.Now()

	return func(item *BlobInfoV1) bool {
		if search != "" && !c.matchSearch(item, search) {
			return false
		}
		if id != "" && id != item.Id {
			return false
		}
		if name != "" && name != item.Name {
			return false
		}
		if group != "" && group != item.Group {
			return false
		}
		if completedOk && completed != item.Completed {
			return false
		}
		if expiredOk && expired && item.ExpireTime.Unix() > now.Unix() {
			return false
		}
		if expiredOk && !expired && item.ExpireTime.Unix() <= now.Unix() {
			return false
		}
		if fromCreateTimeOK && item.CreateTime.Unix() >= fromCreateTime.Unix() {
			return false
		}
		if toCreateTimeOk && item.CreateTime.Unix() < toCreateTime.Unix() {
			return false
		}
		return true
	}
}

func (c *BlobsMockClientV1) GetBlobsByFilter(ctx context.Context, correlationId string, filter *data.FilterParams, paging *data.PagingParams) (result data.DataPage[*BlobInfoV1], err error) {
	filterFunc := c.composeFilter(filter)

	items := make([]*BlobInfoV1, 0)
	for _, v := range c.blobs {
		item := *v
		if filterFunc(&item) {
			items = append(items, &item)
		}
	}
	return *data.NewDataPage(items, len(c.blobs)), nil
}

func (c *BlobsMockClientV1) GetBlobsByIds(ctx context.Context, correlationId string, blobIds []string) (result []*BlobInfoV1, err error) {
	result = make([]*BlobInfoV1, 0)

	for _, b := range c.blobs {
		for _, id := range blobIds {
			if id == b.Id {
				buf := *b
				result = append(result, &buf)
			}
		}
	}

	return result, nil
}

func (c *BlobsMockClientV1) GetBlobById(ctx context.Context, correlationId string, blobId string) (result *BlobInfoV1, err error) {
	for _, b := range c.blobs {
		if blobId == b.Id {
			buf := *b
			result = &buf
		}
	}

	return result, nil
}

func (c *BlobsMockClientV1) CreateBlobFromUri(ctx context.Context, correlationId string, blob *BlobInfoV1,
	uri string) (result *BlobInfoV1, err error) {
	return BlobsUriProcessorV1.CreateBlobFromUri(ctx, correlationId, blob, c, uri, int(c.chunkSize))
}

func (c *BlobsMockClientV1) GetBlobUriById(ctx context.Context, correlationId string, blobId string) (result string, err error) {
	return result, nil
}

func (c *BlobsMockClientV1) CreateBlobFromData(ctx context.Context, correlationId string, blob *BlobInfoV1, buffer []byte) (result *BlobInfoV1, err error) {
	return BlobsDataProcessorV1.CreateBlobFromData(ctx, correlationId, blob, c, buffer, int(c.chunkSize))
}

func (c *BlobsMockClientV1) GetBlobDataById(ctx context.Context, correlationId string, blobId string) (result []byte, blob *BlobInfoV1, err error) {
	return BlobsDataProcessorV1.GetBlobDataById(ctx, correlationId, blobId, c, int(c.chunkSize))
}

func (c *BlobsMockClientV1) CreateBlobFromStream(ctx context.Context, correlationId string, blob *BlobInfoV1, stream io.Reader) (result *BlobInfoV1, err error) {
	return BlobsStreamProcessorV1.CreateBlobFromStream(ctx, correlationId, blob, c, stream, int(c.chunkSize))
}

func (c *BlobsMockClientV1) ReadBlobStreamById(ctx context.Context, correlationId string, blobId string, stream io.Writer) (blob *BlobInfoV1, err error) {
	return BlobsStreamProcessorV1.GetBlobStreamById(ctx, correlationId, blobId, c, stream, int(c.chunkSize))
}

func (c *BlobsMockClientV1) UpdateBlobInfo(ctx context.Context, correlationId string, blob *BlobInfoV1) (result *BlobInfoV1, err error) {
	for i, b := range c.blobs {
		if blob.Id == b.Id {
			buf := *blob
			c.blobs[i] = &buf
			return blob, nil
		}
	}

	return nil, nil
}

func (c *BlobsMockClientV1) MarkBlobsCompleted(ctx context.Context, correlationId string, blobIds []string) error {
	for _, b := range c.blobs {
		for _, id := range blobIds {
			if b.Id == id {
				b.Completed = true
			}
		}
	}

	return nil
}

func (c *BlobsMockClientV1) DeleteBlobById(ctx context.Context, correlationId string, blobId string) error {
	for i, b := range c.blobs {
		if blobId == b.Id {
			c.blobs = append(c.blobs[:i], c.blobs[i+1:]...)
			break
		}
	}

	return nil
}

func (c *BlobsMockClientV1) DeleteBlobsByIds(ctx context.Context, correlationId string, blobIds []string) error {
	for i, b := range c.blobs {
		for _, id := range blobIds {
			if b.Id == id {
				c.blobs = append(c.blobs[:i], c.blobs[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (c *BlobsMockClientV1) normilizeName(name string) string {
	if name == "" {
		return ""
	}

	name = strings.ReplaceAll(name, "\\", "/")
	pos := strings.LastIndex(name, "/")
	if pos >= 0 {
		name = name[pos+1:]
	}

	return name
}

func (c *BlobsMockClientV1) fixBlob(blob *BlobInfoV1) *BlobInfoV1 {
	if blob == nil {
		return nil
	}

	blob.CreateTime, _ = convert.DateTimeConverter.ToNullableDateTime(blob.CreateTime)
	blob.ExpireTime, _ = convert.DateTimeConverter.ToNullableDateTime(blob.ExpireTime)
	blob.Name = c.normilizeName(blob.Name)

	return blob
}

func (c *BlobsMockClientV1) BeginBlobWrite(ctx context.Context, correlationId string, blob *BlobInfoV1) (token string, err error) {
	if blob.Id == "" {
		blob.Id = data.IdGenerator.NextLong()
	}

	blob = c.fixBlob(blob)
	if blob.Size > 0 && blob.Size > c.maxBlobSize {
		return "", errors.NewBadRequestError(correlationId, "BLOB_TOO_LARGE",
			"Blob "+blob.Id+" exceeds allowed maximum size of "+strconv.FormatInt(c.maxBlobSize, 10),
		).WithDetails("blob_id", blob.Id).WithDetails("size", blob.Size).WithDetails("max_size", c.maxBlobSize)
	}

	c.blobs = append(c.blobs, blob)

	c.content[blob.Id] = make([]byte, 0)
	return blob.Id, nil

}

func (c *BlobsMockClientV1) WriteBlobChunk(ctx context.Context, correlationId string, token string, chunk []byte) (token2 string, err error) {
	if chunk == nil {
		chunk = make([]byte, 0)
	}

	id := token
	oldBuffer, ok := c.content[id]
	if !ok {
		return "", errors.NewNotFoundError(correlationId,
			"BLOB_NOT_FOUND",
			"Blob "+id+" was not found",
		).WithDetails("blob_id", id)
	}

	// Enforce maximum size
	chunkLength := 0
	if chunk != nil {
		chunkLength = len(chunk)
	}

	if c.maxBlobSize > 0 && len(oldBuffer)+chunkLength > int(c.maxBlobSize) {
		return "", errors.NewBadRequestError(
			correlationId,
			"BLOB_TOO_LARGE",
			"Blob "+id+" exceeds allowed maximum size of "+strconv.FormatInt(c.maxBlobSize, 10),
		).WithDetails("blob_id", id).WithDetails("size", len(oldBuffer)+chunkLength).WithDetails("max_size", c.maxBlobSize)
	}

	buffer := make([]byte, 0, len(chunk))
	if len(chunk) > 0 {
		buffer = chunk
	}
	c.content[id] = append(oldBuffer, buffer...)

	return token, nil

}

func (c *BlobsMockClientV1) EndBlobWrite(ctx context.Context, correlationId string, token string, chunk []byte) (blob *BlobInfoV1, err error) {
	if chunk == nil {
		chunk = make([]byte, 0)
	}

	id := token

	// Write last chunk of the blob
	_, err = c.WriteBlobChunk(ctx, correlationId, token, chunk)
	if err != nil {
		return blob, err
	}

	blob, err = c.GetBlobById(ctx, correlationId, id)
	if err != nil {
		return blob, err
	}

	if blob == nil {
		return blob, errors.NewNotFoundError(correlationId,
			"BLOB_NOT_FOUND",
			"Blob "+id+" was not found",
		).WithDetails("blob_id", id)
	}

	// Update blob info with size and create time
	buffer, bufOk := c.content[id]
	blob.CreateTime = time.Now()
	if bufOk {
		blob.Size = int64(len(buffer))
	}

	return c.UpdateBlobInfo(ctx, correlationId, blob)
}

func (c *BlobsMockClientV1) AbortBlobWrite(ctx context.Context, correlationId string, token string) error {
	id := token
	return c.DeleteBlobById(ctx, correlationId, id)
}

func (c *BlobsMockClientV1) BeginBlobRead(ctx context.Context, correlationId string, blobId string) (blob *BlobInfoV1, err error) {
	oldBuffer, bufOk := c.content[blobId]
	if !bufOk || len(oldBuffer) == 0 {
		return nil, errors.NewNotFoundError(correlationId,
			"BLOB_NOT_FOUND",
			"Blob "+blobId+" was not found",
		).WithDetails("blob_id", blobId)
	}

	return c.GetBlobById(ctx, correlationId, blobId)
}

func (c *BlobsMockClientV1) ReadBlobChunk(ctx context.Context, correlationId string, blobId string, skip int64, take int64) (chunk []byte, err error) {
	oldBuffer, bufOk := c.content[blobId]
	if !bufOk || len(oldBuffer) == 0 {
		return nil, errors.NewNotFoundError(correlationId,
			"BLOB_NOT_FOUND",
			"Blob "+blobId+" was not found",
		).WithDetails("blob_id", blobId)
	}

	if int(skip) > len(oldBuffer) {
		return make([]byte, 0), nil
	} else if int(skip+take) > len(oldBuffer) {
		return oldBuffer[skip:], nil
	} else {
		return oldBuffer[skip : skip+take], nil
	}

}

func (c *BlobsMockClientV1) EndBlobRead(ctx context.Context, correlationId string, blobId string) error {
	return nil
}
