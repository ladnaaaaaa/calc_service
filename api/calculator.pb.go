package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTaskRequest) Reset() {
	*x = GetTaskRequest{}
	mi := &file_api_calculator_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTaskRequest) ProtoMessage() {}

func (x *GetTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_calculator_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetTaskRequest) Descriptor() ([]byte, []int) {
	return file_api_calculator_proto_rawDescGZIP(), []int{0}
}

type GetTaskResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        uint64                 `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Arg1          float64                `protobuf:"fixed64,2,opt,name=arg1,proto3" json:"arg1,omitempty"`
	Arg2          float64                `protobuf:"fixed64,3,opt,name=arg2,proto3" json:"arg2,omitempty"`
	Operation     string                 `protobuf:"bytes,4,opt,name=operation,proto3" json:"operation,omitempty"`
	OperationTime int64                  `protobuf:"varint,5,opt,name=operation_time,json=operationTime,proto3" json:"operation_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTaskResponse) Reset() {
	*x = GetTaskResponse{}
	mi := &file_api_calculator_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTaskResponse) ProtoMessage() {}

func (x *GetTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_calculator_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetTaskResponse) Descriptor() ([]byte, []int) {
	return file_api_calculator_proto_rawDescGZIP(), []int{1}
}

func (x *GetTaskResponse) GetTaskId() uint64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *GetTaskResponse) GetArg1() float64 {
	if x != nil {
		return x.Arg1
	}
	return 0
}

func (x *GetTaskResponse) GetArg2() float64 {
	if x != nil {
		return x.Arg2
	}
	return 0
}

func (x *GetTaskResponse) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

func (x *GetTaskResponse) GetOperationTime() int64 {
	if x != nil {
		return x.OperationTime
	}
	return 0
}

type SubmitResultRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        uint64                 `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Result        float64                `protobuf:"fixed64,2,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitResultRequest) Reset() {
	*x = SubmitResultRequest{}
	mi := &file_api_calculator_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitResultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitResultRequest) ProtoMessage() {}

func (x *SubmitResultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_calculator_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SubmitResultRequest) Descriptor() ([]byte, []int) {
	return file_api_calculator_proto_rawDescGZIP(), []int{2}
}

func (x *SubmitResultRequest) GetTaskId() uint64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *SubmitResultRequest) GetResult() float64 {
	if x != nil {
		return x.Result
	}
	return 0
}

type SubmitResultResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitResultResponse) Reset() {
	*x = SubmitResultResponse{}
	mi := &file_api_calculator_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitResultResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitResultResponse) ProtoMessage() {}

func (x *SubmitResultResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_calculator_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SubmitResultResponse) Descriptor() ([]byte, []int) {
	return file_api_calculator_proto_rawDescGZIP(), []int{3}
}

func (x *SubmitResultResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SubmitResultResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_api_calculator_proto protoreflect.FileDescriptor

const file_api_calculator_proto_rawDesc = "" +
	"\n" +
	"\x14api/calculator.proto\x12\n" +
	"calculator\"\x10\n" +
	"\x0eGetTaskRequest\"\x97\x01\n" +
	"\x0fGetTaskResponse\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\x04R\x06taskId\x12\x12\n" +
	"\x04arg1\x18\x02 \x01(\x01R\x04arg1\x12\x12\n" +
	"\x04arg2\x18\x03 \x01(\x01R\x04arg2\x12\x1c\n" +
	"\toperation\x18\x04 \x01(\tR\toperation\x12%\n" +
	"\x0eoperation_time\x18\x05 \x01(\x03R\roperationTime\"F\n" +
	"\x13SubmitResultRequest\x12\x17\n" +
	"\atask_id\x18\x01 \x01(\x04R\x06taskId\x12\x16\n" +
	"\x06result\x18\x02 \x01(\x01R\x06result\"F\n" +
	"\x14SubmitResultResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x14\n" +
	"\x05error\x18\x02 \x01(\tR\x05error2\xa7\x01\n" +
	"\n" +
	"Calculator\x12D\n" +
	"\aGetTask\x12\x1a.calculator.GetTaskRequest\x1a\x1b.calculator.GetTaskResponse\"\x00\x12S\n" +
	"\fSubmitResult\x12\x1f.calculator.SubmitResultRequest\x1a .calculator.SubmitResultResponse\"\x00B(Z&github.com/ladnaaaaaa/calc_service/apib\x06proto3"

var (
	file_api_calculator_proto_rawDescOnce sync.Once
	file_api_calculator_proto_rawDescData []byte
)

func file_api_calculator_proto_rawDescGZIP() []byte {
	file_api_calculator_proto_rawDescOnce.Do(func() {
		file_api_calculator_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_calculator_proto_rawDesc), len(file_api_calculator_proto_rawDesc)))
	})
	return file_api_calculator_proto_rawDescData
}

var file_api_calculator_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_calculator_proto_goTypes = []any{
	(*GetTaskRequest)(nil),
	(*GetTaskResponse)(nil),
	(*SubmitResultRequest)(nil),
	(*SubmitResultResponse)(nil),
}
var file_api_calculator_proto_depIdxs = []int32{
	0,
	2,
	1,
	3,
	2,
	0,
	0,
	0,
	0,
}

func init() { file_api_calculator_proto_init() }
func file_api_calculator_proto_init() {
	if File_api_calculator_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_calculator_proto_rawDesc), len(file_api_calculator_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_calculator_proto_goTypes,
		DependencyIndexes: file_api_calculator_proto_depIdxs,
		MessageInfos:      file_api_calculator_proto_msgTypes,
	}.Build()
	File_api_calculator_proto = out.File
	file_api_calculator_proto_goTypes = nil
	file_api_calculator_proto_depIdxs = nil
}
