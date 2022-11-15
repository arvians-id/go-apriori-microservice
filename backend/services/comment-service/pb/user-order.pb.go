// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: services/comment-service/proto/user-order.proto

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

type ListUserOrderResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserOrder []*UserOrder `protobuf:"bytes,1,rep,name=UserOrder,proto3" json:"UserOrder,omitempty"`
}

func (x *ListUserOrderResponse) Reset() {
	*x = ListUserOrderResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_comment_service_proto_user_order_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUserOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserOrderResponse) ProtoMessage() {}

func (x *ListUserOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_comment_service_proto_user_order_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUserOrderResponse.ProtoReflect.Descriptor instead.
func (*ListUserOrderResponse) Descriptor() ([]byte, []int) {
	return file_services_comment_service_proto_user_order_proto_rawDescGZIP(), []int{0}
}

func (x *ListUserOrderResponse) GetUserOrder() []*UserOrder {
	if x != nil {
		return x.UserOrder
	}
	return nil
}

type CreateUserOrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PayloadId      int64   `protobuf:"varint,1,opt,name=PayloadId,proto3" json:"PayloadId,omitempty"`
	Code           *string `protobuf:"bytes,2,opt,name=Code,proto3,oneof" json:"Code,omitempty"`
	Name           *string `protobuf:"bytes,3,opt,name=Name,proto3,oneof" json:"Name,omitempty"`
	Price          *int64  `protobuf:"varint,4,opt,name=Price,proto3,oneof" json:"Price,omitempty"`
	Image          *string `protobuf:"bytes,5,opt,name=Image,proto3,oneof" json:"Image,omitempty"`
	Quantity       *int32  `protobuf:"varint,6,opt,name=Quantity,proto3,oneof" json:"Quantity,omitempty"`
	TotalPriceItem *int64  `protobuf:"varint,7,opt,name=TotalPriceItem,proto3,oneof" json:"TotalPriceItem,omitempty"`
}

func (x *CreateUserOrderRequest) Reset() {
	*x = CreateUserOrderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_comment_service_proto_user_order_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserOrderRequest) ProtoMessage() {}

func (x *CreateUserOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_comment_service_proto_user_order_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserOrderRequest.ProtoReflect.Descriptor instead.
func (*CreateUserOrderRequest) Descriptor() ([]byte, []int) {
	return file_services_comment_service_proto_user_order_proto_rawDescGZIP(), []int{1}
}

func (x *CreateUserOrderRequest) GetPayloadId() int64 {
	if x != nil {
		return x.PayloadId
	}
	return 0
}

func (x *CreateUserOrderRequest) GetCode() string {
	if x != nil && x.Code != nil {
		return *x.Code
	}
	return ""
}

func (x *CreateUserOrderRequest) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *CreateUserOrderRequest) GetPrice() int64 {
	if x != nil && x.Price != nil {
		return *x.Price
	}
	return 0
}

func (x *CreateUserOrderRequest) GetImage() string {
	if x != nil && x.Image != nil {
		return *x.Image
	}
	return ""
}

func (x *CreateUserOrderRequest) GetQuantity() int32 {
	if x != nil && x.Quantity != nil {
		return *x.Quantity
	}
	return 0
}

func (x *CreateUserOrderRequest) GetTotalPriceItem() int64 {
	if x != nil && x.TotalPriceItem != nil {
		return *x.TotalPriceItem
	}
	return 0
}

type GetUserOrderResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserOrder *UserOrder `protobuf:"bytes,1,opt,name=UserOrder,proto3" json:"UserOrder,omitempty"`
}

func (x *GetUserOrderResponse) Reset() {
	*x = GetUserOrderResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_comment_service_proto_user_order_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserOrderResponse) ProtoMessage() {}

func (x *GetUserOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_comment_service_proto_user_order_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserOrderResponse.ProtoReflect.Descriptor instead.
func (*GetUserOrderResponse) Descriptor() ([]byte, []int) {
	return file_services_comment_service_proto_user_order_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserOrderResponse) GetUserOrder() *UserOrder {
	if x != nil {
		return x.UserOrder
	}
	return nil
}

type GetUserOrderByPayloadIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PayloadId int64 `protobuf:"varint,1,opt,name=PayloadId,proto3" json:"PayloadId,omitempty"`
}

func (x *GetUserOrderByPayloadIdRequest) Reset() {
	*x = GetUserOrderByPayloadIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_comment_service_proto_user_order_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserOrderByPayloadIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserOrderByPayloadIdRequest) ProtoMessage() {}

func (x *GetUserOrderByPayloadIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_comment_service_proto_user_order_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserOrderByPayloadIdRequest.ProtoReflect.Descriptor instead.
func (*GetUserOrderByPayloadIdRequest) Descriptor() ([]byte, []int) {
	return file_services_comment_service_proto_user_order_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserOrderByPayloadIdRequest) GetPayloadId() int64 {
	if x != nil {
		return x.PayloadId
	}
	return 0
}

type GetUserOrderByUserIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
}

func (x *GetUserOrderByUserIdRequest) Reset() {
	*x = GetUserOrderByUserIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_comment_service_proto_user_order_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserOrderByUserIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserOrderByUserIdRequest) ProtoMessage() {}

func (x *GetUserOrderByUserIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_comment_service_proto_user_order_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserOrderByUserIdRequest.ProtoReflect.Descriptor instead.
func (*GetUserOrderByUserIdRequest) Descriptor() ([]byte, []int) {
	return file_services_comment_service_proto_user_order_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserOrderByUserIdRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GetUserOrderByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *GetUserOrderByIdRequest) Reset() {
	*x = GetUserOrderByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_comment_service_proto_user_order_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserOrderByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserOrderByIdRequest) ProtoMessage() {}

func (x *GetUserOrderByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_comment_service_proto_user_order_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserOrderByIdRequest.ProtoReflect.Descriptor instead.
func (*GetUserOrderByIdRequest) Descriptor() ([]byte, []int) {
	return file_services_comment_service_proto_user_order_proto_rawDescGZIP(), []int{5}
}

func (x *GetUserOrderByIdRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_services_comment_service_proto_user_order_proto protoreflect.FileDescriptor

var file_services_comment_service_proto_user_order_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x2d, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e,
	0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x52, 0x09, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x22, 0xb2,
	0x02, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x50, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x17, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01,
	0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x50, 0x72, 0x69,
	0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x1f, 0x0a, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x04, 0x52, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x88, 0x01, 0x01,
	0x12, 0x2b, 0x0a, 0x0e, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x49, 0x74,
	0x65, 0x6d, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x48, 0x05, 0x52, 0x0e, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a,
	0x05, 0x5f, 0x43, 0x6f, 0x64, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x4e, 0x61, 0x6d, 0x65, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x50, 0x72, 0x69, 0x63, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x49,
	0x74, 0x65, 0x6d, 0x22, 0x46, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x09, 0x55,
	0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x52, 0x09, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x3e, 0x0a, 0x1e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x22, 0x35, 0x0a, 0x1b, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x22, 0x29, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64, 0x32, 0xd1, 0x02,
	0x0a, 0x10, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x59, 0x0a, 0x12, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x42, 0x79, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x12, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a,
	0x0f, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x47, 0x0a, 0x08, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x49, 0x64, 0x12, 0x1e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x06, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x1f, 0x5a, 0x1d, 0x2e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_comment_service_proto_user_order_proto_rawDescOnce sync.Once
	file_services_comment_service_proto_user_order_proto_rawDescData = file_services_comment_service_proto_user_order_proto_rawDesc
)

func file_services_comment_service_proto_user_order_proto_rawDescGZIP() []byte {
	file_services_comment_service_proto_user_order_proto_rawDescOnce.Do(func() {
		file_services_comment_service_proto_user_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_comment_service_proto_user_order_proto_rawDescData)
	})
	return file_services_comment_service_proto_user_order_proto_rawDescData
}

var file_services_comment_service_proto_user_order_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_services_comment_service_proto_user_order_proto_goTypes = []interface{}{
	(*ListUserOrderResponse)(nil),          // 0: proto.ListUserOrderResponse
	(*CreateUserOrderRequest)(nil),         // 1: proto.CreateUserOrderRequest
	(*GetUserOrderResponse)(nil),           // 2: proto.GetUserOrderResponse
	(*GetUserOrderByPayloadIdRequest)(nil), // 3: proto.GetUserOrderByPayloadIdRequest
	(*GetUserOrderByUserIdRequest)(nil),    // 4: proto.GetUserOrderByUserIdRequest
	(*GetUserOrderByIdRequest)(nil),        // 5: proto.GetUserOrderByIdRequest
	(*UserOrder)(nil),                      // 6: proto.UserOrder
}
var file_services_comment_service_proto_user_order_proto_depIdxs = []int32{
	6, // 0: proto.ListUserOrderResponse.UserOrder:type_name -> proto.UserOrder
	6, // 1: proto.GetUserOrderResponse.UserOrder:type_name -> proto.UserOrder
	3, // 2: proto.UserOrderService.FindAllByPayloadId:input_type -> proto.GetUserOrderByPayloadIdRequest
	4, // 3: proto.UserOrderService.FindAllByUserId:input_type -> proto.GetUserOrderByUserIdRequest
	5, // 4: proto.UserOrderService.FindById:input_type -> proto.GetUserOrderByIdRequest
	1, // 5: proto.UserOrderService.Create:input_type -> proto.CreateUserOrderRequest
	0, // 6: proto.UserOrderService.FindAllByPayloadId:output_type -> proto.ListUserOrderResponse
	0, // 7: proto.UserOrderService.FindAllByUserId:output_type -> proto.ListUserOrderResponse
	2, // 8: proto.UserOrderService.FindById:output_type -> proto.GetUserOrderResponse
	2, // 9: proto.UserOrderService.Create:output_type -> proto.GetUserOrderResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_services_comment_service_proto_user_order_proto_init() }
func file_services_comment_service_proto_user_order_proto_init() {
	if File_services_comment_service_proto_user_order_proto != nil {
		return
	}
	file_services_comment_service_proto_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_services_comment_service_proto_user_order_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUserOrderResponse); i {
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
		file_services_comment_service_proto_user_order_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserOrderRequest); i {
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
		file_services_comment_service_proto_user_order_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserOrderResponse); i {
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
		file_services_comment_service_proto_user_order_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserOrderByPayloadIdRequest); i {
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
		file_services_comment_service_proto_user_order_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserOrderByUserIdRequest); i {
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
		file_services_comment_service_proto_user_order_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserOrderByIdRequest); i {
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
	file_services_comment_service_proto_user_order_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_comment_service_proto_user_order_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_comment_service_proto_user_order_proto_goTypes,
		DependencyIndexes: file_services_comment_service_proto_user_order_proto_depIdxs,
		MessageInfos:      file_services_comment_service_proto_user_order_proto_msgTypes,
	}.Build()
	File_services_comment_service_proto_user_order_proto = out.File
	file_services_comment_service_proto_user_order_proto_rawDesc = nil
	file_services_comment_service_proto_user_order_proto_goTypes = nil
	file_services_comment_service_proto_user_order_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UserOrderServiceClient is the client API for UserOrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserOrderServiceClient interface {
	FindAllByPayloadId(ctx context.Context, in *GetUserOrderByPayloadIdRequest, opts ...grpc.CallOption) (*ListUserOrderResponse, error)
	FindAllByUserId(ctx context.Context, in *GetUserOrderByUserIdRequest, opts ...grpc.CallOption) (*ListUserOrderResponse, error)
	FindById(ctx context.Context, in *GetUserOrderByIdRequest, opts ...grpc.CallOption) (*GetUserOrderResponse, error)
	Create(ctx context.Context, in *CreateUserOrderRequest, opts ...grpc.CallOption) (*GetUserOrderResponse, error)
}

type userOrderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserOrderServiceClient(cc grpc.ClientConnInterface) UserOrderServiceClient {
	return &userOrderServiceClient{cc}
}

func (c *userOrderServiceClient) FindAllByPayloadId(ctx context.Context, in *GetUserOrderByPayloadIdRequest, opts ...grpc.CallOption) (*ListUserOrderResponse, error) {
	out := new(ListUserOrderResponse)
	err := c.cc.Invoke(ctx, "/proto.UserOrderService/FindAllByPayloadId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userOrderServiceClient) FindAllByUserId(ctx context.Context, in *GetUserOrderByUserIdRequest, opts ...grpc.CallOption) (*ListUserOrderResponse, error) {
	out := new(ListUserOrderResponse)
	err := c.cc.Invoke(ctx, "/proto.UserOrderService/FindAllByUserId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userOrderServiceClient) FindById(ctx context.Context, in *GetUserOrderByIdRequest, opts ...grpc.CallOption) (*GetUserOrderResponse, error) {
	out := new(GetUserOrderResponse)
	err := c.cc.Invoke(ctx, "/proto.UserOrderService/FindById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userOrderServiceClient) Create(ctx context.Context, in *CreateUserOrderRequest, opts ...grpc.CallOption) (*GetUserOrderResponse, error) {
	out := new(GetUserOrderResponse)
	err := c.cc.Invoke(ctx, "/proto.UserOrderService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserOrderServiceServer is the server API for UserOrderService service.
type UserOrderServiceServer interface {
	FindAllByPayloadId(context.Context, *GetUserOrderByPayloadIdRequest) (*ListUserOrderResponse, error)
	FindAllByUserId(context.Context, *GetUserOrderByUserIdRequest) (*ListUserOrderResponse, error)
	FindById(context.Context, *GetUserOrderByIdRequest) (*GetUserOrderResponse, error)
	Create(context.Context, *CreateUserOrderRequest) (*GetUserOrderResponse, error)
}

// UnimplementedUserOrderServiceServer can be embedded to have forward compatible implementations.
type UnimplementedUserOrderServiceServer struct {
}

func (*UnimplementedUserOrderServiceServer) FindAllByPayloadId(context.Context, *GetUserOrderByPayloadIdRequest) (*ListUserOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllByPayloadId not implemented")
}
func (*UnimplementedUserOrderServiceServer) FindAllByUserId(context.Context, *GetUserOrderByUserIdRequest) (*ListUserOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllByUserId not implemented")
}
func (*UnimplementedUserOrderServiceServer) FindById(context.Context, *GetUserOrderByIdRequest) (*GetUserOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindById not implemented")
}
func (*UnimplementedUserOrderServiceServer) Create(context.Context, *CreateUserOrderRequest) (*GetUserOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func RegisterUserOrderServiceServer(s *grpc.Server, srv UserOrderServiceServer) {
	s.RegisterService(&_UserOrderService_serviceDesc, srv)
}

func _UserOrderService_FindAllByPayloadId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserOrderByPayloadIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserOrderServiceServer).FindAllByPayloadId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserOrderService/FindAllByPayloadId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserOrderServiceServer).FindAllByPayloadId(ctx, req.(*GetUserOrderByPayloadIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserOrderService_FindAllByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserOrderByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserOrderServiceServer).FindAllByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserOrderService/FindAllByUserId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserOrderServiceServer).FindAllByUserId(ctx, req.(*GetUserOrderByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserOrderService_FindById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserOrderByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserOrderServiceServer).FindById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserOrderService/FindById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserOrderServiceServer).FindById(ctx, req.(*GetUserOrderByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserOrderService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserOrderServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserOrderService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserOrderServiceServer).Create(ctx, req.(*CreateUserOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserOrderService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.UserOrderService",
	HandlerType: (*UserOrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAllByPayloadId",
			Handler:    _UserOrderService_FindAllByPayloadId_Handler,
		},
		{
			MethodName: "FindAllByUserId",
			Handler:    _UserOrderService_FindAllByUserId_Handler,
		},
		{
			MethodName: "FindById",
			Handler:    _UserOrderService_FindById_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _UserOrderService_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/comment-service/proto/user-order.proto",
}
