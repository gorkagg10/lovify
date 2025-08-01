// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: grpc/auth-service/auth-service.proto

package service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegisterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         *string                `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	Password      *string                `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_grpc_auth_service_auth_service_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterRequest) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil && x.Password != nil {
		return *x.Password
	}
	return ""
}

type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         *string                `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	Password      *string                `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_grpc_auth_service_auth_service_proto_rawDescGZIP(), []int{1}
}

func (x *LoginRequest) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil && x.Password != nil {
		return *x.Password
	}
	return ""
}

type LoginResponse struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	SessionToken       *Token                 `protobuf:"bytes,1,opt,name=sessionToken" json:"sessionToken,omitempty"`
	CsrfToken          *Token                 `protobuf:"bytes,2,opt,name=csrfToken" json:"csrfToken,omitempty"`
	IsProfileConnected *bool                  `protobuf:"varint,3,opt,name=isProfileConnected" json:"isProfileConnected,omitempty"`
	ProfileID          *string                `protobuf:"bytes,4,opt,name=profileID" json:"profileID,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_grpc_auth_service_auth_service_proto_rawDescGZIP(), []int{2}
}

func (x *LoginResponse) GetSessionToken() *Token {
	if x != nil {
		return x.SessionToken
	}
	return nil
}

func (x *LoginResponse) GetCsrfToken() *Token {
	if x != nil {
		return x.CsrfToken
	}
	return nil
}

func (x *LoginResponse) GetIsProfileConnected() bool {
	if x != nil && x.IsProfileConnected != nil {
		return *x.IsProfileConnected
	}
	return false
}

func (x *LoginResponse) GetProfileID() string {
	if x != nil && x.ProfileID != nil {
		return *x.ProfileID
	}
	return ""
}

type Token struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Token          *string                `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	ExpirationDate *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expirationDate" json:"expirationDate,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Token) Reset() {
	*x = Token{}
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_grpc_auth_service_auth_service_proto_rawDescGZIP(), []int{3}
}

func (x *Token) GetToken() string {
	if x != nil && x.Token != nil {
		return *x.Token
	}
	return ""
}

func (x *Token) GetExpirationDate() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpirationDate
	}
	return nil
}

type AuthorizationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         *string                `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	SessionToken  *string                `protobuf:"bytes,2,opt,name=sessionToken" json:"sessionToken,omitempty"`
	CsrfToken     *string                `protobuf:"bytes,3,opt,name=csrfToken" json:"csrfToken,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthorizationRequest) Reset() {
	*x = AuthorizationRequest{}
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthorizationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizationRequest) ProtoMessage() {}

func (x *AuthorizationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_auth_service_auth_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizationRequest.ProtoReflect.Descriptor instead.
func (*AuthorizationRequest) Descriptor() ([]byte, []int) {
	return file_grpc_auth_service_auth_service_proto_rawDescGZIP(), []int{4}
}

func (x *AuthorizationRequest) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *AuthorizationRequest) GetSessionToken() string {
	if x != nil && x.SessionToken != nil {
		return *x.SessionToken
	}
	return ""
}

func (x *AuthorizationRequest) GetCsrfToken() string {
	if x != nil && x.CsrfToken != nil {
		return *x.CsrfToken
	}
	return ""
}

var File_grpc_auth_service_auth_service_proto protoreflect.FileDescriptor

const file_grpc_auth_service_auth_service_proto_rawDesc = "" +
	"\n" +
	"$grpc/auth-service/auth-service.proto\x12\x13lovify_auth_service\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"C\n" +
	"\x0fRegisterRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"@\n" +
	"\fLoginRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"\xd7\x01\n" +
	"\rLoginResponse\x12>\n" +
	"\fsessionToken\x18\x01 \x01(\v2\x1a.lovify_auth_service.TokenR\fsessionToken\x128\n" +
	"\tcsrfToken\x18\x02 \x01(\v2\x1a.lovify_auth_service.TokenR\tcsrfToken\x12.\n" +
	"\x12isProfileConnected\x18\x03 \x01(\bR\x12isProfileConnected\x12\x1c\n" +
	"\tprofileID\x18\x04 \x01(\tR\tprofileID\"a\n" +
	"\x05Token\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12B\n" +
	"\x0eexpirationDate\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x0eexpirationDate\"n\n" +
	"\x14AuthorizationRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\"\n" +
	"\fsessionToken\x18\x02 \x01(\tR\fsessionToken\x12\x1c\n" +
	"\tcsrfToken\x18\x03 \x01(\tR\tcsrfToken2\xfb\x01\n" +
	"\vAuthService\x12L\n" +
	"\fRegisterUser\x12$.lovify_auth_service.RegisterRequest\x1a\x16.google.protobuf.Empty\x12N\n" +
	"\x05Login\x12!.lovify_auth_service.LoginRequest\x1a\".lovify_auth_service.LoginResponse\x12N\n" +
	"\tAuthorize\x12).lovify_auth_service.AuthorizationRequest\x1a\x16.google.protobuf.EmptyB\x15Z\x13lovify-auth/serviceb\beditionsp\xe8\a"

var (
	file_grpc_auth_service_auth_service_proto_rawDescOnce sync.Once
	file_grpc_auth_service_auth_service_proto_rawDescData []byte
)

func file_grpc_auth_service_auth_service_proto_rawDescGZIP() []byte {
	file_grpc_auth_service_auth_service_proto_rawDescOnce.Do(func() {
		file_grpc_auth_service_auth_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_grpc_auth_service_auth_service_proto_rawDesc), len(file_grpc_auth_service_auth_service_proto_rawDesc)))
	})
	return file_grpc_auth_service_auth_service_proto_rawDescData
}

var file_grpc_auth_service_auth_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_grpc_auth_service_auth_service_proto_goTypes = []any{
	(*RegisterRequest)(nil),       // 0: lovify_auth_service.RegisterRequest
	(*LoginRequest)(nil),          // 1: lovify_auth_service.LoginRequest
	(*LoginResponse)(nil),         // 2: lovify_auth_service.LoginResponse
	(*Token)(nil),                 // 3: lovify_auth_service.Token
	(*AuthorizationRequest)(nil),  // 4: lovify_auth_service.AuthorizationRequest
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 6: google.protobuf.Empty
}
var file_grpc_auth_service_auth_service_proto_depIdxs = []int32{
	3, // 0: lovify_auth_service.LoginResponse.sessionToken:type_name -> lovify_auth_service.Token
	3, // 1: lovify_auth_service.LoginResponse.csrfToken:type_name -> lovify_auth_service.Token
	5, // 2: lovify_auth_service.Token.expirationDate:type_name -> google.protobuf.Timestamp
	0, // 3: lovify_auth_service.AuthService.RegisterUser:input_type -> lovify_auth_service.RegisterRequest
	1, // 4: lovify_auth_service.AuthService.Login:input_type -> lovify_auth_service.LoginRequest
	4, // 5: lovify_auth_service.AuthService.Authorize:input_type -> lovify_auth_service.AuthorizationRequest
	6, // 6: lovify_auth_service.AuthService.RegisterUser:output_type -> google.protobuf.Empty
	2, // 7: lovify_auth_service.AuthService.Login:output_type -> lovify_auth_service.LoginResponse
	6, // 8: lovify_auth_service.AuthService.Authorize:output_type -> google.protobuf.Empty
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_grpc_auth_service_auth_service_proto_init() }
func file_grpc_auth_service_auth_service_proto_init() {
	if File_grpc_auth_service_auth_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_grpc_auth_service_auth_service_proto_rawDesc), len(file_grpc_auth_service_auth_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_auth_service_auth_service_proto_goTypes,
		DependencyIndexes: file_grpc_auth_service_auth_service_proto_depIdxs,
		MessageInfos:      file_grpc_auth_service_auth_service_proto_msgTypes,
	}.Build()
	File_grpc_auth_service_auth_service_proto = out.File
	file_grpc_auth_service_auth_service_proto_goTypes = nil
	file_grpc_auth_service_auth_service_proto_depIdxs = nil
}
