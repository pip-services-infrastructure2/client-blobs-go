package version1

import (
	"encoding/json"

	"github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	"github.com/pip-services3-gox/pip-services3-commons-gox/data"
	"github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	"github.com/service-infrastructure2/client-blobs-go/protos"
)

func fromError(err error) *protos.ErrorDescription {
	if err == nil {
		return nil
	}

	desc := errors.ErrorDescriptionFactory.Create(err)
	obj := &protos.ErrorDescription{
		Type:          desc.Type,
		Category:      desc.Category,
		Code:          desc.Code,
		CorrelationId: desc.CorrelationId,
		Status:        convert.StringConverter.ToString(desc.Status),
		Message:       desc.Message,
		Cause:         desc.Cause,
		StackTrace:    desc.StackTrace,
		Details:       fromMap(desc.Details),
	}

	return obj
}

func toError(obj *protos.ErrorDescription) error {
	if obj == nil || (obj.Category == "" && obj.Message == "") {
		return nil
	}

	description := &errors.ErrorDescription{
		Type:          obj.Type,
		Category:      obj.Category,
		Code:          obj.Code,
		CorrelationId: obj.CorrelationId,
		Status:        convert.IntegerConverter.ToInteger(obj.Status),
		Message:       obj.Message,
		Cause:         obj.Cause,
		StackTrace:    obj.StackTrace,
		Details:       toMap(obj.Details),
	}

	return errors.ApplicationErrorFactory.Create(description)
}

func fromMap(val map[string]any) map[string]string {
	r := map[string]string{}

	for k, v := range val {
		r[k] = convert.StringConverter.ToString(v)
	}

	return r
}

func toMap(val map[string]string) map[string]any {
	r := map[string]any{}

	for k, v := range val {
		r[k] = v
	}

	return r
}

func toJson(value any) string {
	if value == nil {
		return ""
	}

	b, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(b[:])
}

func fromJson(value string) any {
	if value == "" {
		return nil
	}

	var m any
	json.Unmarshal([]byte(value), &m)
	return m
}

func fromBlobInfo(blob *BlobInfoV1) *protos.BlobInfo {
	if blob == nil {
		return nil
	}

	obj := &protos.BlobInfo{
		Id:          blob.Id,
		Group:       blob.Group,
		Name:        blob.Name,
		Size:        blob.Size,
		ContentType: blob.ContentType,
		CreateTime:  convert.StringConverter.ToString(blob.CreateTime),
		ExpireTime:  convert.StringConverter.ToString(blob.ExpireTime),
		Completed:   blob.Completed,
	}

	return obj
}

func toBlobInfo(obj *protos.BlobInfo) *BlobInfoV1 {
	if obj == nil {
		return nil
	}

	blob := &BlobInfoV1{
		Id:          obj.Id,
		Group:       obj.Group,
		Name:        obj.Name,
		Size:        obj.Size,
		ContentType: obj.ContentType,
		CreateTime:  convert.DateTimeConverter.ToDateTime(obj.CreateTime),
		ExpireTime:  convert.DateTimeConverter.ToDateTime(obj.ExpireTime),
		Completed:   obj.Completed,
	}

	return blob
}

func fromBlobInfoPage(page data.DataPage[*BlobInfoV1]) *protos.BlobInfoPage {
	obj := &protos.BlobInfoPage{
		Total: int64(page.Total),
		Data:  make([]*protos.BlobInfo, len(page.Data)),
	}

	for i, v := range page.Data {
		blob := v
		obj.Data[i] = fromBlobInfo(blob)
	}

	return obj
}

func toBlobInfoPage(obj *protos.BlobInfoPage) data.DataPage[*BlobInfoV1] {
	if obj == nil {
		return *data.NewEmptyDataPage[*BlobInfoV1]()
	}

	blobs := make([]*BlobInfoV1, len(obj.Data))

	for i, v := range obj.Data {
		blobs[i] = toBlobInfo(v)
	}

	page := *data.NewDataPage[*BlobInfoV1](blobs, int(obj.Total))

	return page
}

func fromBlobInfos(data []*BlobInfoV1) []*protos.BlobInfo {
	if data == nil {
		return nil
	}

	obj := make([]*protos.BlobInfo, len(data))

	for i, v := range data {
		obj[i] = fromBlobInfo(v)
	}

	return obj
}

func toBlobInfos(obj []*protos.BlobInfo) []*BlobInfoV1 {
	if obj == nil {
		return nil
	}

	data := make([]*BlobInfoV1, len(obj))

	for i, v := range obj {
		data[i] = toBlobInfo(v)
	}

	return data
}
