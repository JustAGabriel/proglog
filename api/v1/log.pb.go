// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: api/v1/log.proto

package log_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value  []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Offset uint64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_api_v1_log_proto_rawDescGZIP(), []int{0}
}

func (x *Record) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Record) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type CreateRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Record *Record `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *CreateRecordRequest) Reset() {
	*x = CreateRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_log_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRecordRequest) ProtoMessage() {}

func (x *CreateRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_log_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRecordRequest.ProtoReflect.Descriptor instead.
func (*CreateRecordRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_log_proto_rawDescGZIP(), []int{1}
}

func (x *CreateRecordRequest) GetRecord() *Record {
	if x != nil {
		return x.Record
	}
	return nil
}

type CreateRecordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset uint64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *CreateRecordResponse) Reset() {
	*x = CreateRecordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_log_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRecordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRecordResponse) ProtoMessage() {}

func (x *CreateRecordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_log_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRecordResponse.ProtoReflect.Descriptor instead.
func (*CreateRecordResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_log_proto_rawDescGZIP(), []int{2}
}

func (x *CreateRecordResponse) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type GetRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset uint64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetRecordRequest) Reset() {
	*x = GetRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_log_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecordRequest) ProtoMessage() {}

func (x *GetRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_log_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecordRequest.ProtoReflect.Descriptor instead.
func (*GetRecordRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_log_proto_rawDescGZIP(), []int{3}
}

func (x *GetRecordRequest) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type GetRecordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Record *Record `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *GetRecordResponse) Reset() {
	*x = GetRecordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_log_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRecordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecordResponse) ProtoMessage() {}

func (x *GetRecordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_log_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecordResponse.ProtoReflect.Descriptor instead.
func (*GetRecordResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_log_proto_rawDescGZIP(), []int{4}
}

func (x *GetRecordResponse) GetRecord() *Record {
	if x != nil {
		return x.Record
	}
	return nil
}

var File_api_v1_log_proto protoreflect.FileDescriptor

var file_api_v1_log_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x6c, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x22, 0x36, 0x0a, 0x06, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x22, 0x3d, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x06, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6c, 0x6f, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x22, 0x2e, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66,
	0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x22, 0x2a, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x3b, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x52, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x32, 0xa3, 0x02, 0x0a, 0x03, 0x4c,
	0x6f, 0x67, 0x12, 0x45, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x2e, 0x6c,
	0x6f, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6c, 0x6f, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x0c, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1b, 0x2e, 0x6c, 0x6f, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x03, 0x47, 0x65,
	0x74, 0x12, 0x18, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6c, 0x6f,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x18, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x19, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01,
	0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a,
	0x75, 0x73, 0x74, 0x61, 0x67, 0x61, 0x62, 0x72, 0x69, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x67,
	0x6c, 0x6f, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6c, 0x6f, 0x67, 0x5f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_log_proto_rawDescOnce sync.Once
	file_api_v1_log_proto_rawDescData = file_api_v1_log_proto_rawDesc
)

func file_api_v1_log_proto_rawDescGZIP() []byte {
	file_api_v1_log_proto_rawDescOnce.Do(func() {
		file_api_v1_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_log_proto_rawDescData)
	})
	return file_api_v1_log_proto_rawDescData
}

var file_api_v1_log_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_v1_log_proto_goTypes = []interface{}{
	(*Record)(nil),               // 0: log.v1.Record
	(*CreateRecordRequest)(nil),  // 1: log.v1.CreateRecordRequest
	(*CreateRecordResponse)(nil), // 2: log.v1.CreateRecordResponse
	(*GetRecordRequest)(nil),     // 3: log.v1.GetRecordRequest
	(*GetRecordResponse)(nil),    // 4: log.v1.GetRecordResponse
}
var file_api_v1_log_proto_depIdxs = []int32{
	0, // 0: log.v1.CreateRecordRequest.record:type_name -> log.v1.Record
	0, // 1: log.v1.GetRecordResponse.record:type_name -> log.v1.Record
	1, // 2: log.v1.Log.Create:input_type -> log.v1.CreateRecordRequest
	1, // 3: log.v1.Log.CreateStream:input_type -> log.v1.CreateRecordRequest
	3, // 4: log.v1.Log.Get:input_type -> log.v1.GetRecordRequest
	3, // 5: log.v1.Log.GetStream:input_type -> log.v1.GetRecordRequest
	2, // 6: log.v1.Log.Create:output_type -> log.v1.CreateRecordResponse
	2, // 7: log.v1.Log.CreateStream:output_type -> log.v1.CreateRecordResponse
	4, // 8: log.v1.Log.Get:output_type -> log.v1.GetRecordResponse
	4, // 9: log.v1.Log.GetStream:output_type -> log.v1.GetRecordResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_v1_log_proto_init() }
func file_api_v1_log_proto_init() {
	if File_api_v1_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_log_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_log_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRecordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_log_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRecordResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_log_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRecordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_log_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRecordResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_log_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_log_proto_goTypes,
		DependencyIndexes: file_api_v1_log_proto_depIdxs,
		MessageInfos:      file_api_v1_log_proto_msgTypes,
	}.Build()
	File_api_v1_log_proto = out.File
	file_api_v1_log_proto_rawDesc = nil
	file_api_v1_log_proto_goTypes = nil
	file_api_v1_log_proto_depIdxs = nil
}
