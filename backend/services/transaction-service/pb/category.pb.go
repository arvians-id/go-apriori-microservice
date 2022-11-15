// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: services/transaction-service/proto/category.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListCategoryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Categories []*Category `protobuf:"bytes,1,rep,name=Categories,proto3" json:"Categories,omitempty"`
}

func (x *ListCategoryResponse) Reset() {
	*x = ListCategoryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_transaction_service_proto_category_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCategoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCategoryResponse) ProtoMessage() {}

func (x *ListCategoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_transaction_service_proto_category_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCategoryResponse.ProtoReflect.Descriptor instead.
func (*ListCategoryResponse) Descriptor() ([]byte, []int) {
	return file_services_transaction_service_proto_category_proto_rawDescGZIP(), []int{0}
}

func (x *ListCategoryResponse) GetCategories() []*Category {
	if x != nil {
		return x.Categories
	}
	return nil
}

type GetCategoryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Category *Category `protobuf:"bytes,1,opt,name=Category,proto3" json:"Category,omitempty"`
}

func (x *GetCategoryResponse) Reset() {
	*x = GetCategoryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_transaction_service_proto_category_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCategoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCategoryResponse) ProtoMessage() {}

func (x *GetCategoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_transaction_service_proto_category_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCategoryResponse.ProtoReflect.Descriptor instead.
func (*GetCategoryResponse) Descriptor() ([]byte, []int) {
	return file_services_transaction_service_proto_category_proto_rawDescGZIP(), []int{1}
}

func (x *GetCategoryResponse) GetCategory() *Category {
	if x != nil {
		return x.Category
	}
	return nil
}

type GetCategoryByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *GetCategoryByIdRequest) Reset() {
	*x = GetCategoryByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_transaction_service_proto_category_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCategoryByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCategoryByIdRequest) ProtoMessage() {}

func (x *GetCategoryByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_transaction_service_proto_category_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCategoryByIdRequest.ProtoReflect.Descriptor instead.
func (*GetCategoryByIdRequest) Descriptor() ([]byte, []int) {
	return file_services_transaction_service_proto_category_proto_rawDescGZIP(), []int{2}
}

func (x *GetCategoryByIdRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type CreateCategoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *CreateCategoryRequest) Reset() {
	*x = CreateCategoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_transaction_service_proto_category_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCategoryRequest) ProtoMessage() {}

func (x *CreateCategoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_transaction_service_proto_category_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCategoryRequest.ProtoReflect.Descriptor instead.
func (*CreateCategoryRequest) Descriptor() ([]byte, []int) {
	return file_services_transaction_service_proto_category_proto_rawDescGZIP(), []int{3}
}

func (x *CreateCategoryRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type UpdateCategoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdCategory int64  `protobuf:"varint,1,opt,name=IdCategory,proto3" json:"IdCategory,omitempty"`
	Name       string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *UpdateCategoryRequest) Reset() {
	*x = UpdateCategoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_transaction_service_proto_category_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCategoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCategoryRequest) ProtoMessage() {}

func (x *UpdateCategoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_transaction_service_proto_category_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCategoryRequest.ProtoReflect.Descriptor instead.
func (*UpdateCategoryRequest) Descriptor() ([]byte, []int) {
	return file_services_transaction_service_proto_category_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateCategoryRequest) GetIdCategory() int64 {
	if x != nil {
		return x.IdCategory
	}
	return 0
}

func (x *UpdateCategoryRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_services_transaction_service_proto_category_proto protoreflect.FileDescriptor

var file_services_transaction_service_proto_category_proto_rawDesc = []byte{
	0x0a, 0x31, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74,
	0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2f, 0x0a, 0x0a, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x0a, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65,
	0x73, 0x22, 0x42, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x08, 0x43, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x08, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x22, 0x28, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64, 0x22,
	0x2b, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4b, 0x0a, 0x15,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x49, 0x64, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x49, 0x64, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0xe1, 0x02, 0x0a, 0x0f, 0x43, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a,
	0x07, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a,
	0x08, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x49, 0x64, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x42, 0x79, 0x49,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x23, 0x5a,
	0x21, 0x2e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_transaction_service_proto_category_proto_rawDescOnce sync.Once
	file_services_transaction_service_proto_category_proto_rawDescData = file_services_transaction_service_proto_category_proto_rawDesc
)

func file_services_transaction_service_proto_category_proto_rawDescGZIP() []byte {
	file_services_transaction_service_proto_category_proto_rawDescOnce.Do(func() {
		file_services_transaction_service_proto_category_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_transaction_service_proto_category_proto_rawDescData)
	})
	return file_services_transaction_service_proto_category_proto_rawDescData
}

var file_services_transaction_service_proto_category_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_services_transaction_service_proto_category_proto_goTypes = []interface{}{
	(*ListCategoryResponse)(nil),   // 0: proto.ListCategoryResponse
	(*GetCategoryResponse)(nil),    // 1: proto.GetCategoryResponse
	(*GetCategoryByIdRequest)(nil), // 2: proto.GetCategoryByIdRequest
	(*CreateCategoryRequest)(nil),  // 3: proto.CreateCategoryRequest
	(*UpdateCategoryRequest)(nil),  // 4: proto.UpdateCategoryRequest
	(*Category)(nil),               // 5: proto.Category
	(*emptypb.Empty)(nil),          // 6: google.protobuf.Empty
}
var file_services_transaction_service_proto_category_proto_depIdxs = []int32{
	5, // 0: proto.ListCategoryResponse.Categories:type_name -> proto.Category
	5, // 1: proto.GetCategoryResponse.Category:type_name -> proto.Category
	6, // 2: proto.CategoryService.FindAll:input_type -> google.protobuf.Empty
	2, // 3: proto.CategoryService.FindById:input_type -> proto.GetCategoryByIdRequest
	3, // 4: proto.CategoryService.Create:input_type -> proto.CreateCategoryRequest
	4, // 5: proto.CategoryService.Update:input_type -> proto.UpdateCategoryRequest
	2, // 6: proto.CategoryService.Delete:input_type -> proto.GetCategoryByIdRequest
	0, // 7: proto.CategoryService.FindAll:output_type -> proto.ListCategoryResponse
	1, // 8: proto.CategoryService.FindById:output_type -> proto.GetCategoryResponse
	1, // 9: proto.CategoryService.Create:output_type -> proto.GetCategoryResponse
	1, // 10: proto.CategoryService.Update:output_type -> proto.GetCategoryResponse
	6, // 11: proto.CategoryService.Delete:output_type -> google.protobuf.Empty
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_services_transaction_service_proto_category_proto_init() }
func file_services_transaction_service_proto_category_proto_init() {
	if File_services_transaction_service_proto_category_proto != nil {
		return
	}
	file_services_transaction_service_proto_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_services_transaction_service_proto_category_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCategoryResponse); i {
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
		file_services_transaction_service_proto_category_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCategoryResponse); i {
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
		file_services_transaction_service_proto_category_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCategoryByIdRequest); i {
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
		file_services_transaction_service_proto_category_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCategoryRequest); i {
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
		file_services_transaction_service_proto_category_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateCategoryRequest); i {
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
			RawDescriptor: file_services_transaction_service_proto_category_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_transaction_service_proto_category_proto_goTypes,
		DependencyIndexes: file_services_transaction_service_proto_category_proto_depIdxs,
		MessageInfos:      file_services_transaction_service_proto_category_proto_msgTypes,
	}.Build()
	File_services_transaction_service_proto_category_proto = out.File
	file_services_transaction_service_proto_category_proto_rawDesc = nil
	file_services_transaction_service_proto_category_proto_goTypes = nil
	file_services_transaction_service_proto_category_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CategoryServiceClient is the client API for CategoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CategoryServiceClient interface {
	FindAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListCategoryResponse, error)
	FindById(ctx context.Context, in *GetCategoryByIdRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error)
	Create(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error)
	Update(ctx context.Context, in *UpdateCategoryRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error)
	Delete(ctx context.Context, in *GetCategoryByIdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type categoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCategoryServiceClient(cc grpc.ClientConnInterface) CategoryServiceClient {
	return &categoryServiceClient{cc}
}

func (c *categoryServiceClient) FindAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListCategoryResponse, error) {
	out := new(ListCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/FindAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) FindById(ctx context.Context, in *GetCategoryByIdRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error) {
	out := new(GetCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/FindById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) Create(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error) {
	out := new(GetCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) Update(ctx context.Context, in *UpdateCategoryRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error) {
	out := new(GetCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) Delete(ctx context.Context, in *GetCategoryByIdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CategoryServiceServer is the server API for CategoryService service.
type CategoryServiceServer interface {
	FindAll(context.Context, *emptypb.Empty) (*ListCategoryResponse, error)
	FindById(context.Context, *GetCategoryByIdRequest) (*GetCategoryResponse, error)
	Create(context.Context, *CreateCategoryRequest) (*GetCategoryResponse, error)
	Update(context.Context, *UpdateCategoryRequest) (*GetCategoryResponse, error)
	Delete(context.Context, *GetCategoryByIdRequest) (*emptypb.Empty, error)
}

// UnimplementedCategoryServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCategoryServiceServer struct {
}

func (*UnimplementedCategoryServiceServer) FindAll(context.Context, *emptypb.Empty) (*ListCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAll not implemented")
}
func (*UnimplementedCategoryServiceServer) FindById(context.Context, *GetCategoryByIdRequest) (*GetCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindById not implemented")
}
func (*UnimplementedCategoryServiceServer) Create(context.Context, *CreateCategoryRequest) (*GetCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedCategoryServiceServer) Update(context.Context, *UpdateCategoryRequest) (*GetCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedCategoryServiceServer) Delete(context.Context, *GetCategoryByIdRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func RegisterCategoryServiceServer(s *grpc.Server, srv CategoryServiceServer) {
	s.RegisterService(&_CategoryService_serviceDesc, srv)
}

func _CategoryService_FindAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).FindAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/FindAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).FindAll(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_FindById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).FindById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/FindById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).FindById(ctx, req.(*GetCategoryByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).Create(ctx, req.(*CreateCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).Update(ctx, req.(*UpdateCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).Delete(ctx, req.(*GetCategoryByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CategoryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.CategoryService",
	HandlerType: (*CategoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAll",
			Handler:    _CategoryService_FindAll_Handler,
		},
		{
			MethodName: "FindById",
			Handler:    _CategoryService_FindById_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _CategoryService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _CategoryService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _CategoryService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/transaction-service/proto/category.proto",
}
