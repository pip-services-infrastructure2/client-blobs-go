package version1

import "context"

type TBlobsDataProcessorV1 struct{}

var BlobsDataProcessorV1 = &TBlobsDataProcessorV1{}

func (c *TBlobsDataProcessorV1) CreateBlobFromData(ctx context.Context, correlationId string, blob *BlobInfoV1,
	writer IBlobsChunkyWriterV1, data []byte, chunkSize int) (*BlobInfoV1, error) {

	buffer := data
	skip := 0
	size := len(buffer)
	token := ""

	// Start writing when first chunk comes
	token, err := writer.BeginBlobWrite(ctx, correlationId, blob)
	if err != nil {
		return nil, err
	}

	// Write chunks
	for size > chunkSize {
		take := chunkSize
		if take > len(buffer)-skip {
			take = len(buffer) - skip
		}
		chunk := buffer[skip : skip+take]

		token, err = writer.WriteBlobChunk(ctx, correlationId, token, chunk)
		if err != nil {
			writer.AbortBlobWrite(ctx, correlationId, token)
			return nil, err
		}

		skip = skip + take
		size = size - take
	}

	// End writing
	chunk := buffer[skip:]

	blob, err = writer.EndBlobWrite(ctx, correlationId, token, chunk)
	if err != nil {
		writer.AbortBlobWrite(ctx, correlationId, token)
		return nil, err
	}

	return blob, nil
}

func (c *TBlobsDataProcessorV1) GetBlobDataById(ctx context.Context, correlationId string, blobId string,
	reader IBlobsChunkyReaderV1, chunkSize int) ([]byte, *BlobInfoV1, error) {

	// Read blob, start reading
	blob, err := reader.BeginBlobRead(ctx, correlationId, blobId)
	if err != nil {
		return nil, nil, err
	}

	// Read all chunks until the end
	skip := int64(0)
	size := blob.Size
	buffer := []byte{}

	for size > 0 {
		take := int64(chunkSize)
		if take > size {
			take = size
		}

		chunk, err1 := reader.ReadBlobChunk(ctx, correlationId, blobId, skip, take)
		if err1 != nil {
			return nil, nil, err1
		}

		if chunk != nil {
			buffer = append(buffer, chunk...)
			size = size - int64(len(chunk))
			skip = skip + int64(len(chunk))
		}
	}

	// End reading
	err = reader.EndBlobRead(ctx, correlationId, blobId)
	if err != nil {
		return nil, nil, err
	}

	return buffer, blob, nil
}
