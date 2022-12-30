package version1

import (
	"context"
	"errors"
	"io"
	"math"
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

		if err1 == io.EOF && size == 0 {
			break
		}
		if err1 != nil && err1 != io.EOF {
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
	if err != nil {
		return nil, err
	}

	if blob == nil || blob.Size == 0 {
		return nil, errors.New("BLOB WITH ZERO SIZE")
	}

	size := blob.Size

	// Read in chunks
	skip := int64(0)
	take := int64(math.Min(float64(chunkSize), float64(size)))
	for {
		buffer, err1 := reader.ReadBlobChunk(ctx, correlationId, blobId, skip, take)
		if err1 != nil {
			return nil, err1
		}

		// Protection against infinite loop
		if len(buffer) > 0 {
			n, err2 := stream.Write(buffer)
			if err2 != nil {
				return nil, err
			}

			size -= int64(n)
			skip += int64(n)
			take = int64(math.Min(float64(chunkSize), float64(size)))

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
