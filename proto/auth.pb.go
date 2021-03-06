// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/auth.proto

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type UserRole struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *UserRole) Reset()                    { *m = UserRole{} }
func (m *UserRole) String() string            { return proto1.CompactTextString(m) }
func (*UserRole) ProtoMessage()               {}
func (*UserRole) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *UserRole) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type LoginRequest struct {
}

func (m *LoginRequest) Reset()                    { *m = LoginRequest{} }
func (m *LoginRequest) String() string            { return proto1.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()               {}
func (*LoginRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

type LoginReply struct {
	State string `protobuf:"bytes,1,opt,name=state" json:"state,omitempty"`
	Url   string `protobuf:"bytes,2,opt,name=url" json:"url,omitempty"`
	Err   string `protobuf:"bytes,3,opt,name=err" json:"err,omitempty"`
}

func (m *LoginReply) Reset()                    { *m = LoginReply{} }
func (m *LoginReply) String() string            { return proto1.CompactTextString(m) }
func (*LoginReply) ProtoMessage()               {}
func (*LoginReply) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *LoginReply) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *LoginReply) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *LoginReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type LoginUserRequest struct {
	State string `protobuf:"bytes,1,opt,name=state" json:"state,omitempty"`
	Code  string `protobuf:"bytes,2,opt,name=code" json:"code,omitempty"`
}

func (m *LoginUserRequest) Reset()                    { *m = LoginUserRequest{} }
func (m *LoginUserRequest) String() string            { return proto1.CompactTextString(m) }
func (*LoginUserRequest) ProtoMessage()               {}
func (*LoginUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *LoginUserRequest) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *LoginUserRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type LoginUserReply struct {
	Authorization *Authorization `protobuf:"bytes,1,opt,name=authorization" json:"authorization,omitempty"`
	Err           string         `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *LoginUserReply) Reset()                    { *m = LoginUserReply{} }
func (m *LoginUserReply) String() string            { return proto1.CompactTextString(m) }
func (*LoginUserReply) ProtoMessage()               {}
func (*LoginUserReply) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *LoginUserReply) GetAuthorization() *Authorization {
	if m != nil {
		return m.Authorization
	}
	return nil
}

func (m *LoginUserReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type MakeTokenRequest struct {
	Uid string `protobuf:"bytes,1,opt,name=uid" json:"uid,omitempty"`
}

func (m *MakeTokenRequest) Reset()                    { *m = MakeTokenRequest{} }
func (m *MakeTokenRequest) String() string            { return proto1.CompactTextString(m) }
func (*MakeTokenRequest) ProtoMessage()               {}
func (*MakeTokenRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *MakeTokenRequest) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

type MakeTokenReply struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	Url   string `protobuf:"bytes,2,opt,name=url" json:"url,omitempty"`
	Err   string `protobuf:"bytes,3,opt,name=err" json:"err,omitempty"`
}

func (m *MakeTokenReply) Reset()                    { *m = MakeTokenReply{} }
func (m *MakeTokenReply) String() string            { return proto1.CompactTextString(m) }
func (*MakeTokenReply) ProtoMessage()               {}
func (*MakeTokenReply) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *MakeTokenReply) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *MakeTokenReply) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *MakeTokenReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type Authorization struct {
	Uid        string `protobuf:"bytes,1,opt,name=uid" json:"uid,omitempty"`
	Username   string `protobuf:"bytes,2,opt,name=username" json:"username,omitempty"`
	Email      string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
	FullName   string `protobuf:"bytes,4,opt,name=full_name,json=fullName" json:"full_name,omitempty"`
	PictureUrl string `protobuf:"bytes,8,opt,name=picture_url,json=pictureUrl" json:"picture_url,omitempty"`
	ProviderId string `protobuf:"bytes,9,opt,name=provider_id,json=providerId" json:"provider_id,omitempty"`
	Name       string `protobuf:"bytes,10,opt,name=name" json:"name,omitempty"`
}

func (m *Authorization) Reset()                    { *m = Authorization{} }
func (m *Authorization) String() string            { return proto1.CompactTextString(m) }
func (*Authorization) ProtoMessage()               {}
func (*Authorization) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *Authorization) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *Authorization) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Authorization) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Authorization) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *Authorization) GetPictureUrl() string {
	if m != nil {
		return m.PictureUrl
	}
	return ""
}

func (m *Authorization) GetProviderId() string {
	if m != nil {
		return m.ProviderId
	}
	return ""
}

func (m *Authorization) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type LoginAsRequest struct {
	Authorization *Authorization `protobuf:"bytes,1,opt,name=authorization" json:"authorization,omitempty"`
}

func (m *LoginAsRequest) Reset()                    { *m = LoginAsRequest{} }
func (m *LoginAsRequest) String() string            { return proto1.CompactTextString(m) }
func (*LoginAsRequest) ProtoMessage()               {}
func (*LoginAsRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

func (m *LoginAsRequest) GetAuthorization() *Authorization {
	if m != nil {
		return m.Authorization
	}
	return nil
}

type LoginAsReply struct {
	Err string `protobuf:"bytes,1,opt,name=err" json:"err,omitempty"`
}

func (m *LoginAsReply) Reset()                    { *m = LoginAsReply{} }
func (m *LoginAsReply) String() string            { return proto1.CompactTextString(m) }
func (*LoginAsReply) ProtoMessage()               {}
func (*LoginAsReply) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{9} }

func (m *LoginAsReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type LogoutRequest struct {
}

func (m *LogoutRequest) Reset()                    { *m = LogoutRequest{} }
func (m *LogoutRequest) String() string            { return proto1.CompactTextString(m) }
func (*LogoutRequest) ProtoMessage()               {}
func (*LogoutRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{10} }

type LogoutReply struct {
	Err string `protobuf:"bytes,1,opt,name=err" json:"err,omitempty"`
}

func (m *LogoutReply) Reset()                    { *m = LogoutReply{} }
func (m *LogoutReply) String() string            { return proto1.CompactTextString(m) }
func (*LogoutReply) ProtoMessage()               {}
func (*LogoutReply) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{11} }

func (m *LogoutReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto1.RegisterType((*UserRole)(nil), "proto.UserRole")
	proto1.RegisterType((*LoginRequest)(nil), "proto.LoginRequest")
	proto1.RegisterType((*LoginReply)(nil), "proto.LoginReply")
	proto1.RegisterType((*LoginUserRequest)(nil), "proto.LoginUserRequest")
	proto1.RegisterType((*LoginUserReply)(nil), "proto.LoginUserReply")
	proto1.RegisterType((*MakeTokenRequest)(nil), "proto.MakeTokenRequest")
	proto1.RegisterType((*MakeTokenReply)(nil), "proto.MakeTokenReply")
	proto1.RegisterType((*Authorization)(nil), "proto.Authorization")
	proto1.RegisterType((*LoginAsRequest)(nil), "proto.LoginAsRequest")
	proto1.RegisterType((*LoginAsReply)(nil), "proto.LoginAsReply")
	proto1.RegisterType((*LogoutRequest)(nil), "proto.LogoutRequest")
	proto1.RegisterType((*LogoutReply)(nil), "proto.LogoutReply")
}

func init() { proto1.RegisterFile("proto/auth.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 368 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0xcf, 0x4b, 0x2b, 0x31,
	0x10, 0x66, 0xfb, 0xe3, 0xd1, 0x4e, 0x5f, 0xfb, 0x96, 0xd0, 0xc3, 0xf2, 0x04, 0x5b, 0x16, 0x0f,
	0x9e, 0x2a, 0xe8, 0x4d, 0xbc, 0xf4, 0x22, 0x28, 0xd5, 0xc3, 0x62, 0xaf, 0x96, 0xb5, 0x3b, 0xd6,
	0xd0, 0xb4, 0xa9, 0xd9, 0x44, 0xd0, 0x3f, 0xd0, 0xbf, 0x4b, 0x32, 0x49, 0xda, 0xad, 0x28, 0x08,
	0x9e, 0x76, 0xe6, 0xdb, 0x99, 0x6f, 0xbe, 0xf9, 0x32, 0x10, 0x6f, 0x94, 0xd4, 0xf2, 0x24, 0x37,
	0xfa, 0x69, 0x44, 0x21, 0x6b, 0xd2, 0x27, 0x3d, 0x84, 0xd6, 0xb4, 0x44, 0x95, 0x49, 0x81, 0x8c,
	0x41, 0x63, 0x9d, 0xaf, 0x30, 0x89, 0x86, 0xd1, 0x71, 0x3b, 0xa3, 0x38, 0xed, 0xc1, 0xdf, 0x89,
	0x5c, 0xf0, 0x75, 0x86, 0xcf, 0x06, 0x4b, 0x9d, 0x5e, 0x02, 0xf8, 0x7c, 0x23, 0x5e, 0x59, 0x1f,
	0x9a, 0xa5, 0xce, 0x75, 0x68, 0x71, 0x09, 0x8b, 0xa1, 0x6e, 0x94, 0x48, 0x6a, 0x84, 0xd9, 0xd0,
	0x22, 0xa8, 0x54, 0x52, 0x77, 0x08, 0x2a, 0x95, 0x5e, 0x40, 0x4c, 0x3c, 0x34, 0xdc, 0x71, 0x7f,
	0xc3, 0xc6, 0xa0, 0x31, 0x97, 0x05, 0x7a, 0x3a, 0x8a, 0xd3, 0x7b, 0xe8, 0x55, 0xba, 0xad, 0x92,
	0x73, 0xe8, 0xda, 0xe5, 0xa4, 0xe2, 0x6f, 0xb9, 0xe6, 0x72, 0x4d, 0x1c, 0x9d, 0xd3, 0xbe, 0xdb,
	0x76, 0x34, 0xae, 0xfe, 0xcb, 0xf6, 0x4b, 0x83, 0xba, 0xda, 0x4e, 0xdd, 0x11, 0xc4, 0x37, 0xf9,
	0x12, 0xef, 0xe4, 0x12, 0xc3, 0xe6, 0xb4, 0x15, 0x2f, 0xbc, 0x36, 0x1b, 0xa6, 0xd7, 0xd0, 0xab,
	0x54, 0x79, 0x3f, 0xb4, 0xcd, 0xc2, 0x06, 0x94, 0xfc, 0xc8, 0x8f, 0xf7, 0x08, 0xba, 0xe3, 0xcf,
	0xaa, 0xf6, 0xe7, 0xb1, 0xff, 0xd0, 0x32, 0x25, 0x2a, 0x7a, 0x23, 0x47, 0xb6, 0xcd, 0xed, 0x64,
	0x5c, 0xe5, 0x5c, 0x78, 0x4e, 0x97, 0xb0, 0x03, 0x68, 0x3f, 0x1a, 0x21, 0x66, 0xd4, 0xd2, 0x70,
	0x2d, 0x16, 0xb8, 0xb5, 0x2d, 0x03, 0xe8, 0x6c, 0xf8, 0x5c, 0x1b, 0x85, 0x33, 0x2b, 0xaf, 0x45,
	0xbf, 0xc1, 0x43, 0x53, 0x25, 0xa8, 0x40, 0xc9, 0x17, 0x5e, 0xa0, 0x9a, 0xf1, 0x22, 0x69, 0xfb,
	0x02, 0x0f, 0x5d, 0x15, 0xdb, 0x83, 0x81, 0xca, 0xc1, 0x4c, 0xfc, 0xd3, 0x8c, 0xcb, 0x60, 0xdc,
	0x2f, 0x9e, 0x26, 0x1d, 0xfa, 0xf3, 0xb3, 0x6c, 0xd6, 0x60, 0x6f, 0x5c, 0xb4, 0x33, 0xee, 0x1f,
	0x74, 0x27, 0x72, 0x21, 0x8d, 0x0e, 0x17, 0x3a, 0x80, 0x4e, 0x00, 0xbe, 0xec, 0x78, 0xf8, 0x43,
	0x73, 0xcf, 0x3e, 0x02, 0x00, 0x00, 0xff, 0xff, 0xa4, 0xc3, 0xc9, 0x7d, 0x14, 0x03, 0x00, 0x00,
}
