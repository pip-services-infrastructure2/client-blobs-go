package build

import (
	"github.com/pip-services-infrastructure2/client-blobs-go/version1"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
)

type BlobsClientFactory struct {
	*cbuild.Factory
}

func NewBlobsClientFactory() *BlobsClientFactory {
	c := BlobsClientFactory{}
	c.Factory = cbuild.NewFactory()

	nullClientDescriptor := cref.NewDescriptor("service-blobs", "client", "null", "*", "1.0")
	cmdHttpClientDescriptor := cref.NewDescriptor("service-blobs", "client", "commandable-http", "*", "1.0")
	cmdGrpcClientDescriptor := cref.NewDescriptor("service-blobs", "client", "commandable-grpc", "*", "1.0")

	c.RegisterType(nullClientDescriptor, version1.NewBlobsNullClientV1)
	c.RegisterType(cmdHttpClientDescriptor, version1.NewBlobsCommandableHttpClientV1)
	c.RegisterType(cmdGrpcClientDescriptor, version1.NewBlobsCommandableGrpcClientV1)
	return &c
}
