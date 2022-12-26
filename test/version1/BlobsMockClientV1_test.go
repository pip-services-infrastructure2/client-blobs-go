package test_version1

import (
	"testing"

	"github.com/pip-services-infrastructure2/client-blobs-go/version1"
)

type blobsMockClientV1Test struct {
	client  *version1.BlobsMockClientV1
	fixture *BlobsClientFixtureV1
}

func newBlobsMockClientV1Test() *blobsMockClientV1Test {
	return &blobsMockClientV1Test{}
}

func (c *blobsMockClientV1Test) setup(t *testing.T) {
	c.client = version1.NewBlobsMockClientV1()
	c.fixture = NewBlobsClientFixtureV1(c.client)
}

func (c *blobsMockClientV1Test) teardown(t *testing.T) {
	c.client = nil
}

func TestMockReadWriteChunks(t *testing.T) {
	c := newBlobsMockClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteChunks(t)
}

func TestMockReadWriteData(t *testing.T) {
	c := newBlobsMockClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteData(t)
}

func TestMockReadWriteStream(t *testing.T) {
	c := newBlobsMockClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteStream(t)
}

func TestMockWritingBlobUri(t *testing.T) {
	c := newBlobsMockClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestWritingBlobUri(t)
}

func TestMockGetUriForMissingBlob(t *testing.T) {
	c := newBlobsMockClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestGetUriForMissingBlob(t)
}
