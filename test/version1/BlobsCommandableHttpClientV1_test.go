package test_version1

import (
	"context"
	"os"
	"testing"

	"github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/service-infrastructure2/client-blobs-go/version1"
)

type blobsCommandableHttpClientV1Test struct {
	client  *version1.BlobsCommandableHttpClientV1
	fixture *BlobsClientFixtureV1
}

func newBlobsCommandableHttpClientV1Test() *blobsCommandableHttpClientV1Test {
	return &blobsCommandableHttpClientV1Test{}
}

func (c *blobsCommandableHttpClientV1Test) setup(t *testing.T) {
	var HTTP_HOST = os.Getenv("HTTP_HOST")
	if HTTP_HOST == "" {
		HTTP_HOST = "localhost"
	}
	var HTTP_PORT = os.Getenv("HTTP_PORT")
	if HTTP_PORT == "" {
		HTTP_PORT = "8080"
	}

	var httpConfig = config.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", HTTP_HOST,
		"connection.port", HTTP_PORT,
	)

	c.client = version1.NewBlobsCommandableHttpClientV1()
	c.client.Configure(context.Background(), httpConfig)
	c.client.Open(context.Background(), "")

	c.fixture = NewBlobsClientFixtureV1(c.client)
}

func (c *blobsCommandableHttpClientV1Test) teardown(t *testing.T) {
	c.client.Close(context.Background(), "")
}

func TestCommandableHttpReadWriteChunks(t *testing.T) {
	c := newBlobsCommandableHttpClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteChunks(t)
}

func TestCommandableHttpReadWriteData(t *testing.T) {
	c := newBlobsCommandableHttpClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteData(t)
}

func TestCommandableHttpReadWriteStream(t *testing.T) {
	c := newBlobsCommandableHttpClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteStream(t)
}

func TestCommandableHttpWritingBlobUri(t *testing.T) {
	c := newBlobsCommandableHttpClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestWritingBlobUri(t)
}

func TestCommandableHttpGetUriForMissingBlob(t *testing.T) {
	c := newBlobsCommandableHttpClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestGetUriForMissingBlob(t)
}
