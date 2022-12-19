package test_version1

import (
	"context"
	"os"
	"testing"

	"github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/service-infrastructure2/client-blobs-go/version1"
)

type blobsGrpcClientV1Test struct {
	client  *version1.BlobGrpcClientV1
	fixture *BlobsClientFixtureV1
}

func newBlobsGrpcClientV1Test() *blobsGrpcClientV1Test {
	return &blobsGrpcClientV1Test{}
}

func (c *blobsGrpcClientV1Test) setup(t *testing.T) {
	var GRPC_HOST = os.Getenv("GRPC_HOST")
	if GRPC_HOST == "" {
		GRPC_HOST = "localhost"
	}
	var GRPC_PORT = os.Getenv("GRPC_PORT")
	if GRPC_PORT == "" {
		GRPC_PORT = "8090"
	}

	var httpConfig = config.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", GRPC_HOST,
		"connection.port", GRPC_PORT,
	)

	c.client = version1.NewBlobGrpcClientV1()
	c.client.Configure(context.Background(), httpConfig)
	c.client.Open(context.Background(), "")

	c.fixture = NewBlobsClientFixtureV1(c.client)
}

func (c *blobsGrpcClientV1Test) teardown(t *testing.T) {
	c.client.Close(context.Background(), "")
}

func TestGrpcReadWriteChunks(t *testing.T) {
	c := newBlobsGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteChunks(t)
}

func TestGrpcReadWriteData(t *testing.T) {
	c := newBlobsGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteData(t)
}

func TestGrpcReadWriteStream(t *testing.T) {
	c := newBlobsGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteStream(t)
}

func TestGrpcWritingBlobUri(t *testing.T) {
	c := newBlobsGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestWritingBlobUri(t)
}

func TestGrpcGetUriForMissingBlob(t *testing.T) {
	c := newBlobsGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestGetUriForMissingBlob(t)
}
