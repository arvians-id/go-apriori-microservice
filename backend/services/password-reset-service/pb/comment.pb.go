// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: services/password-reset-service/proto/comment.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ListRatingFromCommentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RatingFromComments []*RatingFromComment `protobuf:"bytes,1,rep,name=RatingFromComments,proto3" json:"RatingFromComments,omitempty"`
}

func (x *ListRatingFromCommentResponse) Reset() {
	*x = ListRatingFromCommentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRatingFromCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRatingFromCommentResponse) ProtoMessage() {}

func (x *ListRatingFromCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRatingFromCommentResponse.ProtoReflect.Descriptor instead.
func (*ListRatingFromCommentResponse) Descriptor() ([]byte, []int) {
	return file_services_password_reset_service_proto_comment_proto_rawDescGZIP(), []int{0}
}

func (x *ListRatingFromCommentResponse) GetRatingFromComments() []*RatingFromComment {
	if x != nil {
		return x.RatingFromComments
	}
	return nil
}

type ListCommentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Comment []*Comment `protobuf:"bytes,1,rep,name=comment,proto3" json:"comment,omitempty"`
}

func (x *ListCommentResponse) Reset() {
	*x = ListCommentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCommentResponse) ProtoMessage() {}

func (x *ListCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCommentResponse.ProtoReflect.Descriptor instead.
func (*ListCommentResponse) Descriptor() ([]byte, []int) {
	return file_services_password_reset_service_proto_comment_proto_rawDescGZIP(), []int{1}
}

func (x *ListCommentResponse) GetComment() []*Comment {
	if x != nil {
		return x.Comment
	}
	return nil
}

type GetCommentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Comment *Comment `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
}

func (x *GetCommentResponse) Reset() {
	*x = GetCommentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentResponse) ProtoMessage() {}

func (x *GetCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentResponse.ProtoReflect.Descriptor instead.
func (*GetCommentResponse) Descriptor() ([]byte, []int) {
	return file_services_password_reset_service_proto_comment_proto_rawDescGZIP(), []int{2}
}

func (x *GetCommentResponse) GetComment() *Comment {
	if x != nil {
		return x.Comment
	}
	return nil
}

type GetCommentByProductCodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductCode string `protobuf:"bytes,1,opt,name=ProductCode,proto3" json:"ProductCode,omitempty"`
}

func (x *GetCommentByProductCodeRequest) Reset() {
	*x = GetCommentByProductCodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCommentByProductCodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentByProductCodeRequest) ProtoMessage() {}

func (x *GetCommentByProductCodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentByProductCodeRequest.ProtoReflect.Descriptor instead.
func (*GetCommentByProductCodeRequest) Descriptor() ([]byte, []int) {
	return file_services_password_reset_service_proto_comment_proto_rawDescGZIP(), []int{3}
}

func (x *GetCommentByProductCodeRequest) GetProductCode() string {
	if x != nil {
		return x.ProductCode
	}
	return ""
}

type GetCommentByFiltersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductCode string `protobuf:"bytes,1,opt,name=ProductCode,proto3" json:"ProductCode,omitempty"`
	Rating      string `protobuf:"bytes,2,opt,name=Rating,proto3" json:"Rating,omitempty"`
	Tag         string `protobuf:"bytes,3,opt,name=Tag,proto3" json:"Tag,omitempty"`
}

func (x *GetCommentByFiltersRequest) Reset() {
	*x = GetCommentByFiltersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCommentByFiltersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentByFiltersRequest) ProtoMessage() {}

func (x *GetCommentByFiltersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentByFiltersRequest.ProtoReflect.Descriptor instead.
func (*GetCommentByFiltersRequest) Descriptor() ([]byte, []int) {
	return file_services_password_reset_service_proto_comment_proto_rawDescGZIP(), []int{4}
}

func (x *GetCommentByFiltersRequest) GetProductCode() string {
	if x != nil {
		return x.ProductCode
	}
	return ""
}

func (x *GetCommentByFiltersRequest) GetRating() string {
	if x != nil {
		return x.Rating
	}
	return ""
}

func (x *GetCommentByFiltersRequest) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

type GetCommentByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *GetCommentByIdRequest) Reset() {
	*x = GetCommentByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCommentByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentByIdRequest) ProtoMessage() {}

func (x *GetCommentByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentByIdRequest.ProtoReflect.Descriptor instead.
func (*GetCommentByIdRequest) Descriptor() ([]byte, []int) {
	return file_services_password_reset_service_proto_comment_proto_rawDescGZIP(), []int{5}
}

func (x *GetCommentByIdRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetCommentByUserOrderIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserOrderId int64 `protobuf:"varint,1,opt,name=UserOrderId,proto3" json:"UserOrderId,omitempty"`
}

func (x *GetCommentByUserOrderIdRequest) Reset() {
	*x = GetCommentByUserOrderIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCommentByUserOrderIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentByUserOrderIdRequest) ProtoMessage() {}

func (x *GetCommentByUserOrderIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentByUserOrderIdRequest.ProtoReflect.Descriptor instead.
func (*GetCommentByUserOrderIdRequest) Descriptor() ([]byte, []int) {
	return file_services_password_reset_service_proto_comment_proto_rawDescGZIP(), []int{6}
}

func (x *GetCommentByUserOrderIdRequest) GetUserOrderId() int64 {
	if x != nil {
		return x.UserOrderId
	}
	return 0
}

type CreateCommentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserOrderId int64  `protobuf:"varint,1,opt,name=UserOrderId,proto3" json:"UserOrderId,omitempty"`
	ProductCode string `protobuf:"bytes,2,opt,name=ProductCode,proto3" json:"ProductCode,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=Description,proto3" json:"Description,omitempty"`
	Tag         string `protobuf:"bytes,4,opt,name=Tag,proto3" json:"Tag,omitempty"`
	Rating      int32  `protobuf:"varint,5,opt,name=Rating,proto3" json:"Rating,omitempty"`
}

func (x *CreateCommentRequest) Reset() {
	*x = CreateCommentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCommentRequest) ProtoMessage() {}

func (x *CreateCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_password_reset_service_proto_comment_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCommentRequest.ProtoReflect.Descriptor instead.
func (*CreateCommentRequest) Descriptor() ([]byte, []int) {
	return file_services_password_reset_service_proto_comment_proto_rawDescGZIP(), []int{7}
}

func (x *CreateCommentRequest) GetUserOrderId() int64 {
	if x != nil {
		return x.UserOrderId
	}
	return 0
}

func (x *CreateCommentRequest) GetProductCode() string {
	if x != nil {
		return x.ProductCode
	}
	return ""
}

func (x *CreateCommentRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateCommentRequest) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *CreateCommentRequest) GetRating() int32 {
	if x != nil {
		return x.Rating
	}
	return 0
}

var File_services_password_reset_service_proto_comment_proto protoreflect.FileDescriptor

var file_services_password_reset_service_proto_comment_proto_rawDesc = []byte{
	0x0a, 0x33, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x32, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2d,
	0x72, 0x65, 0x73, 0x65, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x69, 0x0a, 0x1d, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x46, 0x72,
	0x6f, 0x6d, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x48, 0x0a, 0x12, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x46, 0x72, 0x6f, 0x6d, 0x43,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x46, 0x72, 0x6f, 0x6d,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x12, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x46,
	0x72, 0x6f, 0x6d, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x3f, 0x0a, 0x13, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x3e, 0x0a, 0x12,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x42, 0x0a, 0x1e,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20,
	0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65,
	0x22, 0x68, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x79,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20,
	0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x54, 0x61, 0x67, 0x22, 0x27, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x49, 0x64, 0x22, 0x42, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x55, 0x73, 0x65, 0x72,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0xa6, 0x01, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x20, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x43, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x54, 0x61, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x61, 0x74, 0x69,
	0x6e, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x32, 0xb0, 0x03, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x69, 0x0a, 0x1a, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x61,
	0x74, 0x69, 0x6e, 0x67, 0x42, 0x79, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x46, 0x72, 0x6f, 0x6d, 0x43,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55,
	0x0a, 0x14, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x42, 0x79, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x46, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x08, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x49,
	0x64, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x11, 0x46, 0x69,
	0x6e, 0x64, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x40, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x26, 0x5a, 0x24, 0x2e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x74,
	0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_services_password_reset_service_proto_comment_proto_rawDescOnce sync.Once
	file_services_password_reset_service_proto_comment_proto_rawDescData = file_services_password_reset_service_proto_comment_proto_rawDesc
)

func file_services_password_reset_service_proto_comment_proto_rawDescGZIP() []byte {
	file_services_password_reset_service_proto_comment_proto_rawDescOnce.Do(func() {
		file_services_password_reset_service_proto_comment_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_password_reset_service_proto_comment_proto_rawDescData)
	})
	return file_services_password_reset_service_proto_comment_proto_rawDescData
}

var file_services_password_reset_service_proto_comment_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_services_password_reset_service_proto_comment_proto_goTypes = []interface{}{
	(*ListRatingFromCommentResponse)(nil),  // 0: proto.ListRatingFromCommentResponse
	(*ListCommentResponse)(nil),            // 1: proto.ListCommentResponse
	(*GetCommentResponse)(nil),             // 2: proto.GetCommentResponse
	(*GetCommentByProductCodeRequest)(nil), // 3: proto.GetCommentByProductCodeRequest
	(*GetCommentByFiltersRequest)(nil),     // 4: proto.GetCommentByFiltersRequest
	(*GetCommentByIdRequest)(nil),          // 5: proto.GetCommentByIdRequest
	(*GetCommentByUserOrderIdRequest)(nil), // 6: proto.GetCommentByUserOrderIdRequest
	(*CreateCommentRequest)(nil),           // 7: proto.CreateCommentRequest
	(*RatingFromComment)(nil),              // 8: proto.RatingFromComment
	(*Comment)(nil),                        // 9: proto.Comment
}
var file_services_password_reset_service_proto_comment_proto_depIdxs = []int32{
	8, // 0: proto.ListRatingFromCommentResponse.RatingFromComments:type_name -> proto.RatingFromComment
	9, // 1: proto.ListCommentResponse.comment:type_name -> proto.Comment
	9, // 2: proto.GetCommentResponse.comment:type_name -> proto.Comment
	3, // 3: proto.CommentService.FindAllRatingByProductCode:input_type -> proto.GetCommentByProductCodeRequest
	4, // 4: proto.CommentService.FindAllByProductCode:input_type -> proto.GetCommentByFiltersRequest
	5, // 5: proto.CommentService.FindById:input_type -> proto.GetCommentByIdRequest
	6, // 6: proto.CommentService.FindByUserOrderId:input_type -> proto.GetCommentByUserOrderIdRequest
	7, // 7: proto.CommentService.Create:input_type -> proto.CreateCommentRequest
	0, // 8: proto.CommentService.FindAllRatingByProductCode:output_type -> proto.ListRatingFromCommentResponse
	1, // 9: proto.CommentService.FindAllByProductCode:output_type -> proto.ListCommentResponse
	2, // 10: proto.CommentService.FindById:output_type -> proto.GetCommentResponse
	2, // 11: proto.CommentService.FindByUserOrderId:output_type -> proto.GetCommentResponse
	2, // 12: proto.CommentService.Create:output_type -> proto.GetCommentResponse
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_services_password_reset_service_proto_comment_proto_init() }
func file_services_password_reset_service_proto_comment_proto_init() {
	if File_services_password_reset_service_proto_comment_proto != nil {
		return
	}
	file_services_password_reset_service_proto_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_services_password_reset_service_proto_comment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRatingFromCommentResponse); i {
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
		file_services_password_reset_service_proto_comment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCommentResponse); i {
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
		file_services_password_reset_service_proto_comment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCommentResponse); i {
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
		file_services_password_reset_service_proto_comment_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCommentByProductCodeRequest); i {
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
		file_services_password_reset_service_proto_comment_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCommentByFiltersRequest); i {
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
		file_services_password_reset_service_proto_comment_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCommentByIdRequest); i {
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
		file_services_password_reset_service_proto_comment_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCommentByUserOrderIdRequest); i {
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
		file_services_password_reset_service_proto_comment_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCommentRequest); i {
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
			RawDescriptor: file_services_password_reset_service_proto_comment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_password_reset_service_proto_comment_proto_goTypes,
		DependencyIndexes: file_services_password_reset_service_proto_comment_proto_depIdxs,
		MessageInfos:      file_services_password_reset_service_proto_comment_proto_msgTypes,
	}.Build()
	File_services_password_reset_service_proto_comment_proto = out.File
	file_services_password_reset_service_proto_comment_proto_rawDesc = nil
	file_services_password_reset_service_proto_comment_proto_goTypes = nil
	file_services_password_reset_service_proto_comment_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CommentServiceClient is the client API for CommentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CommentServiceClient interface {
	FindAllRatingByProductCode(ctx context.Context, in *GetCommentByProductCodeRequest, opts ...grpc.CallOption) (*ListRatingFromCommentResponse, error)
	FindAllByProductCode(ctx context.Context, in *GetCommentByFiltersRequest, opts ...grpc.CallOption) (*ListCommentResponse, error)
	FindById(ctx context.Context, in *GetCommentByIdRequest, opts ...grpc.CallOption) (*GetCommentResponse, error)
	FindByUserOrderId(ctx context.Context, in *GetCommentByUserOrderIdRequest, opts ...grpc.CallOption) (*GetCommentResponse, error)
	Create(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*GetCommentResponse, error)
}

type commentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommentServiceClient(cc grpc.ClientConnInterface) CommentServiceClient {
	return &commentServiceClient{cc}
}

func (c *commentServiceClient) FindAllRatingByProductCode(ctx context.Context, in *GetCommentByProductCodeRequest, opts ...grpc.CallOption) (*ListRatingFromCommentResponse, error) {
	out := new(ListRatingFromCommentResponse)
	err := c.cc.Invoke(ctx, "/proto.CommentService/FindAllRatingByProductCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) FindAllByProductCode(ctx context.Context, in *GetCommentByFiltersRequest, opts ...grpc.CallOption) (*ListCommentResponse, error) {
	out := new(ListCommentResponse)
	err := c.cc.Invoke(ctx, "/proto.CommentService/FindAllByProductCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) FindById(ctx context.Context, in *GetCommentByIdRequest, opts ...grpc.CallOption) (*GetCommentResponse, error) {
	out := new(GetCommentResponse)
	err := c.cc.Invoke(ctx, "/proto.CommentService/FindById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) FindByUserOrderId(ctx context.Context, in *GetCommentByUserOrderIdRequest, opts ...grpc.CallOption) (*GetCommentResponse, error) {
	out := new(GetCommentResponse)
	err := c.cc.Invoke(ctx, "/proto.CommentService/FindByUserOrderId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) Create(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*GetCommentResponse, error) {
	out := new(GetCommentResponse)
	err := c.cc.Invoke(ctx, "/proto.CommentService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommentServiceServer is the server API for CommentService service.
type CommentServiceServer interface {
	FindAllRatingByProductCode(context.Context, *GetCommentByProductCodeRequest) (*ListRatingFromCommentResponse, error)
	FindAllByProductCode(context.Context, *GetCommentByFiltersRequest) (*ListCommentResponse, error)
	FindById(context.Context, *GetCommentByIdRequest) (*GetCommentResponse, error)
	FindByUserOrderId(context.Context, *GetCommentByUserOrderIdRequest) (*GetCommentResponse, error)
	Create(context.Context, *CreateCommentRequest) (*GetCommentResponse, error)
}

// UnimplementedCommentServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCommentServiceServer struct {
}

func (*UnimplementedCommentServiceServer) FindAllRatingByProductCode(context.Context, *GetCommentByProductCodeRequest) (*ListRatingFromCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllRatingByProductCode not implemented")
}
func (*UnimplementedCommentServiceServer) FindAllByProductCode(context.Context, *GetCommentByFiltersRequest) (*ListCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllByProductCode not implemented")
}
func (*UnimplementedCommentServiceServer) FindById(context.Context, *GetCommentByIdRequest) (*GetCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindById not implemented")
}
func (*UnimplementedCommentServiceServer) FindByUserOrderId(context.Context, *GetCommentByUserOrderIdRequest) (*GetCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByUserOrderId not implemented")
}
func (*UnimplementedCommentServiceServer) Create(context.Context, *CreateCommentRequest) (*GetCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func RegisterCommentServiceServer(s *grpc.Server, srv CommentServiceServer) {
	s.RegisterService(&_CommentService_serviceDesc, srv)
}

func _CommentService_FindAllRatingByProductCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentByProductCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).FindAllRatingByProductCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CommentService/FindAllRatingByProductCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).FindAllRatingByProductCode(ctx, req.(*GetCommentByProductCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_FindAllByProductCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentByFiltersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).FindAllByProductCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CommentService/FindAllByProductCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).FindAllByProductCode(ctx, req.(*GetCommentByFiltersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_FindById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).FindById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CommentService/FindById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).FindById(ctx, req.(*GetCommentByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_FindByUserOrderId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentByUserOrderIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).FindByUserOrderId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CommentService/FindByUserOrderId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).FindByUserOrderId(ctx, req.(*GetCommentByUserOrderIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CommentService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).Create(ctx, req.(*CreateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CommentService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.CommentService",
	HandlerType: (*CommentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAllRatingByProductCode",
			Handler:    _CommentService_FindAllRatingByProductCode_Handler,
		},
		{
			MethodName: "FindAllByProductCode",
			Handler:    _CommentService_FindAllByProductCode_Handler,
		},
		{
			MethodName: "FindById",
			Handler:    _CommentService_FindById_Handler,
		},
		{
			MethodName: "FindByUserOrderId",
			Handler:    _CommentService_FindByUserOrderId_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _CommentService_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/password-reset-service/proto/comment.proto",
}
