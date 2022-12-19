package test_version1

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pip-services3-gox/pip-services3-commons-gox/data"

	"github.com/service-infrastructure2/client-blobs-go/version1"
	"github.com/stretchr/testify/assert"
)

type BlobsClientFixtureV1 struct {
	Client version1.IBlobsClientV1
}

var BLOB_ID1 = data.IdGenerator.NextLong()
var BLOB_ID2 = data.IdGenerator.NextLong()
var BLOB1 = version1.NewBlobInfoV1(BLOB_ID1, "test", "test1.dat", 100, "text/plain")
var BLOB2 = version1.NewBlobInfoV1(BLOB_ID2, "test", "test2.dat", 100, "text/plain")

func NewBlobsClientFixtureV1(client version1.IBlobsClientV1) *BlobsClientFixtureV1 {
	return &BlobsClientFixtureV1{
		Client: client,
	}
}

func (c *BlobsClientFixtureV1) clear() {
	page, _ := c.Client.GetBlobsByFilter(context.Background(), "", nil, nil)

	for _, v := range page.Data {
		blob := v
		c.Client.DeleteBlobById(context.Background(), "", blob.Id)
	}
}

func (c *BlobsClientFixtureV1) TestReadWriteChunks(t *testing.T) {
	c.clear()
	defer c.clear()

	blobId := data.IdGenerator.NextLong()

	// Start writing blob
	blob := version1.NewBlobInfoV1(
		blobId, "test", "file-"+blobId+".dat", 6, "application/binary",
	)

	w := c.Client.(version1.IBlobsChunkyWriterV1)
	tok, err := w.BeginBlobWrite(context.Background(), "", blob)

	assert.Nil(t, err)
	assert.NotEqual(t, "", tok)
	token := tok

	// Write blob
	data := []byte{1, 2, 3}

	tok, err = w.WriteBlobChunk(context.Background(), "", token, data)

	assert.Nil(t, err)
	assert.NotEqual(t, "", tok)
	token = tok

	// Finish writing blob
	data = []byte{4, 5, 6}

	blob, err = w.EndBlobWrite(context.Background(), "", token, data)

	assert.Nil(t, err)
	assert.NotNil(t, blob)

	// Start reading
	r := c.Client.(version1.IBlobsChunkyReaderV1)
	blob, err = r.BeginBlobRead(context.Background(), "", blobId)

	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, int64(6), blob.Size)

	// Read first chunk
	data, err = r.ReadBlobChunk(context.Background(), "", blobId, 0, 3)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, uint8(1), data[0])
	assert.Equal(t, uint8(2), data[1])
	assert.Equal(t, uint8(3), data[2])

	// Get blobs
	page, err1 := c.Client.GetBlobsByFilter(context.Background(), "", nil, nil)
	assert.Nil(t, err1)

	assert.NotNil(t, page)
	assert.Equal(t, 1, len(page.Data))

	// Delete blob
	err = c.Client.DeleteBlobsByIds(context.Background(), "", []string{blobId})
	assert.Nil(t, err)

	// Try to get deleted blob
	blob, err = c.Client.GetBlobById(context.Background(), "", blobId)

	assert.Nil(t, err)
	assert.Nil(t, blob)
}

func (c *BlobsClientFixtureV1) TestReadWriteData(t *testing.T) {
	c.clear()
	defer c.clear()

	blobId := data.IdGenerator.NextLong()

	// Create blob
	blob := version1.NewBlobInfoV1(
		blobId, "test", "file-"+blobId+".dat", 6, "application/binary",
	)
	data := []byte{1, 2, 3, 4, 5, 6}

	blob1, err := c.Client.CreateBlobFromData(context.Background(), "", blob, data)

	assert.Nil(t, err)
	assert.NotNil(t, blob1)
	assert.Equal(t, int64(6), blob1.Size)

	// Get blob info
	blob, err = c.Client.GetBlobById(context.Background(), "", blobId)

	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, int64(6), blob.Size)

	// Read blob
	data, blob, err = c.Client.GetBlobDataById(context.Background(), "", blobId)

	assert.Nil(t, err)

	assert.Equal(t, int64(6), blob.Size)
	assert.Equal(t, 6, len(data))
	assert.Equal(t, uint8(1), data[0])
	assert.Equal(t, uint8(2), data[1])
	assert.Equal(t, uint8(3), data[2])
	assert.Equal(t, uint8(4), data[3])
	assert.Equal(t, uint8(5), data[4])
	assert.Equal(t, uint8(6), data[5])
}

func (c *BlobsClientFixtureV1) TestReadWriteStream(t *testing.T) {
	c.clear()
	defer c.clear()

	blobId := data.IdGenerator.NextLong()
	sample, _ := ioutil.ReadFile("../../data/file.txt")

	blob := version1.NewBlobInfoV1(
		blobId, "test", "../../data/file.txt", 0, "text/plain",
	)

	rs, err := os.Open("../../data/file.txt")
	assert.Nil(t, err)

	blob, err = c.Client.CreateBlobFromStream(context.Background(), "", blob, rs)

	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, "file.txt", blob.Name)
	assert.Equal(t, len(sample), int(blob.Size))

	// Get blob info
	blob, err = c.Client.GetBlobById(context.Background(), "", blobId)

	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, "file.txt", blob.Name)
	assert.Equal(t, len(sample), int(blob.Size))

	// Read blob
	os.Remove("../../data/file.tmp")

	ws, _ := os.Create("../../data/file.tmp")

	blob, err = c.Client.ReadBlobStreamById(context.Background(), "", blobId, ws)

	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, "file.txt", blob.Name)
	assert.Equal(t, len(sample), int(blob.Size))

	sample1, _ := ioutil.ReadFile("../../data/file.tmp")
	os.Remove("../../data/file.tmp")
	assert.Equal(t, sample, sample1)
}

func (c *BlobsClientFixtureV1) TestWritingBlobUri(t *testing.T) {
	c.clear()
	defer c.clear()

	// Writing blob
	blobId := data.IdGenerator.NextLong()
	blob := version1.NewBlobInfoV1(
		blobId, "test", "blob-"+blobId+".dat", 0, "text/plain",
	)

	blob1, err := c.Client.CreateBlobFromUri(context.Background(), "", blob, "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png")

	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.Equal(t, blob.Name, blob1.Name)
	assert.Equal(t, blob.Group, blob1.Group)
	assert.Equal(t, blob.ContentType, blob1.ContentType)
	assert.True(t, blob1.Size > 0)

	// Reading blob
	data, _, err1 := c.Client.GetBlobDataById(context.Background(), "", blobId)

	assert.Nil(t, err1)
	assert.True(t, len(data) > 0)

	c.Client.DeleteBlobsByIds(context.Background(), "", []string{blobId})
}

func (c *BlobsClientFixtureV1) TestGetUriForMissingBlob(t *testing.T) {
	c.clear()
	defer c.clear()

	uri, err := c.Client.GetBlobUriById(context.Background(), "", "123")

	assert.Equal(t, "", uri)
	assert.Nil(t, err)
}
