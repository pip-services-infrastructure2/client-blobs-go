package version1

import (
	"context"
	"io"
	"time"

	"github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

type TBlobsStreamProcessorV1 struct{}

var BlobsStreamProcessorV1 = &TBlobsStreamProcessorV1{}

func (c *TBlobsStreamProcessorV1) CreateBlobFromStream(ctx context.Context, correlationId string, blob *BlobInfoV1,
	writer IBlobsChunkyWriterV1, stream io.Reader, chunkSize int) (*BlobInfoV1, error) {

	// Generate blob id
	if blob.Id == "" {
		blob.Id = data.IdGenerator.NextLong()
	}
	blob.CreateTime = time.Now()

	// Start writing
	token, err := writer.BeginBlobWrite(ctx, correlationId, blob)
	if err != nil {
		return nil, err
	}

	// Write in chunks
	buffer := make([]byte, chunkSize)

	for {
		size, err1 := stream.Read(buffer)

		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			return nil, err1
		}

		if size == 0 {
			continue
		}

		chunk := buffer

		if size != len(buffer) {
			chunk = buffer[0:size]
		}

		token, err = writer.WriteBlobChunk(ctx, correlationId, token, chunk)
		if err != nil {
			return nil, err
		}
	}

	// Finish writing and return blobId
	return writer.EndBlobWrite(ctx, correlationId, token, nil)
}

func (c *TBlobsStreamProcessorV1) GetBlobStreamById(ctx context.Context, correlationId string,
	blobId string, reader IBlobsChunkyReaderV1, stream io.Writer, chunkSize int) (*BlobInfoV1, error) {

	// Begin blob read
	blob, err := reader.BeginBlobRead(ctx, correlationId, blobId)

	// Read in chunks
	start := int64(0)
	for {
		buffer, err1 := reader.ReadBlobChunk(ctx, correlationId, blobId, start, int64(chunkSize))
		if err1 != nil {
			return nil, err1
		}

		// Protection against infinite loop
		if len(buffer) > 0 {
			n, err2 := stream.Write(buffer)
			if err2 != nil {
				return nil, err
			}

			start = start + int64(n)

			if len(buffer) < chunkSize {
				break
			}
		} else {
			break
		}
	}

	// Close blob read
	err = reader.EndBlobRead(ctx, correlationId, blobId)
	if err != nil {
		return nil, err
	}

	return blob, nil
}
