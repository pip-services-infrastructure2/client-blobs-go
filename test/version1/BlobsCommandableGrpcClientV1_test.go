package test_version1

import (
	"context"
	"os"
	"testing"

	"github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/service-infrastructure2/client-blobs-go/version1"
)

type blobsCommandableGrpcClientV1Test struct {
	client  *version1.BlobsCommandableGrpcClientV1
	fixture *BlobsClientFixtureV1
}

func newBlobsCommandableGrpcClientV1Test() *blobsCommandableGrpcClientV1Test {
	return &blobsCommandableGrpcClientV1Test{}
}

func (c *blobsCommandableGrpcClientV1Test) setup(t *testing.T) {
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

	c.client = version1.NewBlobsCommandableGrpcClientV1()
	c.client.Configure(context.Background(), httpConfig)
	c.client.Open(context.Background(), "")

	c.fixture = NewBlobsClientFixtureV1(c.client)
}

func (c *blobsCommandableGrpcClientV1Test) teardown(t *testing.T) {
	c.client.Close(context.Background(), "")
}

func TestCommandableGrpcReadWriteChunks(t *testing.T) {
	c := newBlobsCommandableGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteChunks(t)
}

func TestCommandableGrpcReadWriteData(t *testing.T) {
	c := newBlobsCommandableGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteData(t)
}

func TestCommandableGrpcReadWriteStream(t *testing.T) {
	c := newBlobsCommandableGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestReadWriteStream(t)
}

func TestCommandableGrpcWritingBlobUri(t *testing.T) {
	c := newBlobsCommandableGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestWritingBlobUri(t)
}

func TestCommandableGrpcGetUriForMissingBlob(t *testing.T) {
	c := newBlobsCommandableGrpcClientV1Test()
	c.setup(t)
	defer c.teardown(t)

	c.fixture.TestGetUriForMissingBlob(t)
}
