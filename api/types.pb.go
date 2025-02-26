// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ProxySessionType int32

const (
	ProxySessionType_Undefined         ProxySessionType = 0
	ProxySessionType_Client            ProxySessionType = 1
	ProxySessionType_Controller        ProxySessionType = 2
	ProxySessionType_Proxy             ProxySessionType = 3
	ProxySessionType_Actor             ProxySessionType = 4
	ProxySessionType_PassiveController ProxySessionType = 5
	ProxySessionType_DataProxy         ProxySessionType = 6
)

var ProxySessionType_name = map[int32]string{
	0: "Undefined",
	1: "Client",
	2: "Controller",
	3: "Proxy",
	4: "Actor",
	5: "PassiveController",
	6: "DataProxy",
}

var ProxySessionType_value = map[string]int32{
	"Undefined":         0,
	"Client":            1,
	"Controller":        2,
	"Proxy":             3,
	"Actor":             4,
	"PassiveController": 5,
	"DataProxy":         6,
}

func (x ProxySessionType) String() string {
	return proto.EnumName(ProxySessionType_name, int32(x))
}

func (ProxySessionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

type String struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

func (m *String) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_String.Unmarshal(m, b)
}
func (m *String) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_String.Marshal(b, m, deterministic)
}
func (m *String) XXX_Merge(src proto.Message) {
	xxx_messageInfo_String.Merge(m, src)
}
func (m *String) XXX_Size() int {
	return xxx_messageInfo_String.Size(m)
}
func (m *String) XXX_DiscardUnknown() {
	xxx_messageInfo_String.DiscardUnknown(m)
}

var xxx_messageInfo_String proto.InternalMessageInfo

func (m *String) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Bytes struct {
	Value                []byte   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bytes) Reset()         { *m = Bytes{} }
func (m *Bytes) String() string { return proto.CompactTextString(m) }
func (*Bytes) ProtoMessage()    {}
func (*Bytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}

func (m *Bytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bytes.Unmarshal(m, b)
}
func (m *Bytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bytes.Marshal(b, m, deterministic)
}
func (m *Bytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bytes.Merge(m, src)
}
func (m *Bytes) XXX_Size() int {
	return xxx_messageInfo_Bytes.Size(m)
}
func (m *Bytes) XXX_DiscardUnknown() {
	xxx_messageInfo_Bytes.DiscardUnknown(m)
}

var xxx_messageInfo_Bytes proto.InternalMessageInfo

func (m *Bytes) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type Boolean struct {
	Value                bool     `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Boolean) Reset()         { *m = Boolean{} }
func (m *Boolean) String() string { return proto.CompactTextString(m) }
func (*Boolean) ProtoMessage()    {}
func (*Boolean) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{2}
}

func (m *Boolean) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Boolean.Unmarshal(m, b)
}
func (m *Boolean) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Boolean.Marshal(b, m, deterministic)
}
func (m *Boolean) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Boolean.Merge(m, src)
}
func (m *Boolean) XXX_Size() int {
	return xxx_messageInfo_Boolean.Size(m)
}
func (m *Boolean) XXX_DiscardUnknown() {
	xxx_messageInfo_Boolean.DiscardUnknown(m)
}

var xxx_messageInfo_Boolean proto.InternalMessageInfo

func (m *Boolean) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

type ProxyClient struct {
	Type                 ProxySessionType `protobuf:"varint,1,opt,name=type,proto3,enum=info_age.bremote.ProxySessionType" json:"type,omitempty"`
	Instance             string           `protobuf:"bytes,2,opt,name=instance,proto3" json:"instance,omitempty"`
	Status               string           `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ProxyClient) Reset()         { *m = ProxyClient{} }
func (m *ProxyClient) String() string { return proto.CompactTextString(m) }
func (*ProxyClient) ProtoMessage()    {}
func (*ProxyClient) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{3}
}

func (m *ProxyClient) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyClient.Unmarshal(m, b)
}
func (m *ProxyClient) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyClient.Marshal(b, m, deterministic)
}
func (m *ProxyClient) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyClient.Merge(m, src)
}
func (m *ProxyClient) XXX_Size() int {
	return xxx_messageInfo_ProxyClient.Size(m)
}
func (m *ProxyClient) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyClient.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyClient proto.InternalMessageInfo

func (m *ProxyClient) GetType() ProxySessionType {
	if m != nil {
		return m.Type
	}
	return ProxySessionType_Undefined
}

func (m *ProxyClient) GetInstance() string {
	if m != nil {
		return m.Instance
	}
	return ""
}

func (m *ProxyClient) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterEnum("info_age.bremote.ProxySessionType", ProxySessionType_name, ProxySessionType_value)
	proto.RegisterType((*String)(nil), "info_age.bremote.String")
	proto.RegisterType((*Bytes)(nil), "info_age.bremote.Bytes")
	proto.RegisterType((*Boolean)(nil), "info_age.bremote.Boolean")
	proto.RegisterType((*ProxyClient)(nil), "info_age.bremote.ProxyClient")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor_d938547f84707355) }

var fileDescriptor_d938547f84707355 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0x4d, 0x4b, 0xf3, 0x40,
	0x10, 0xc7, 0x9f, 0xbe, 0x24, 0x4f, 0x33, 0xd5, 0xb2, 0x2e, 0x2a, 0x41, 0xf1, 0x85, 0x9c, 0x44,
	0x30, 0x07, 0x05, 0xef, 0xa6, 0x7e, 0x80, 0xd2, 0xea, 0xc5, 0x8b, 0x6c, 0xdb, 0x69, 0x59, 0x88,
	0x3b, 0x61, 0x77, 0x1a, 0xdc, 0x6f, 0x2f, 0xd9, 0x58, 0x69, 0x7b, 0xdb, 0x3f, 0xff, 0xdf, 0x0c,
	0xbf, 0x61, 0x61, 0xc8, 0xbe, 0x42, 0x97, 0x57, 0x96, 0x98, 0xa4, 0xd0, 0x66, 0x45, 0x9f, 0x6a,
	0x8d, 0xf9, 0xdc, 0xe2, 0x17, 0x31, 0x66, 0xd7, 0x10, 0xcf, 0xd8, 0x6a, 0xb3, 0x96, 0xa7, 0x10,
	0xd5, 0xaa, 0xdc, 0x60, 0xda, 0xb9, 0xed, 0xdc, 0x25, 0xd3, 0x36, 0x64, 0x57, 0x10, 0x15, 0x9e,
	0xd1, 0xed, 0xd7, 0x47, 0xdb, 0xfa, 0x06, 0xfe, 0x17, 0x44, 0x25, 0x2a, 0xb3, 0x0f, 0x0c, 0xb6,
	0x80, 0x87, 0xe1, 0xc4, 0xd2, 0xb7, 0x1f, 0x97, 0x1a, 0x0d, 0xcb, 0x67, 0xe8, 0x37, 0x3e, 0x81,
	0x19, 0x3d, 0x66, 0xf9, 0xa1, 0x4f, 0x1e, 0xe0, 0x19, 0x3a, 0xa7, 0xc9, 0xbc, 0xf9, 0x0a, 0xa7,
	0x81, 0x97, 0x17, 0x30, 0xd0, 0xc6, 0xb1, 0x32, 0x0b, 0x4c, 0xbb, 0xc1, 0xef, 0x2f, 0xcb, 0x73,
	0x88, 0x1d, 0x2b, 0xde, 0xb8, 0xb4, 0x17, 0x9a, 0xdf, 0x74, 0xef, 0x41, 0x1c, 0x6e, 0x93, 0xc7,
	0x90, 0xbc, 0x9b, 0x25, 0xae, 0xb4, 0xc1, 0xa5, 0xf8, 0x27, 0x01, 0xe2, 0x56, 0x4c, 0x74, 0xe4,
	0x08, 0x60, 0x4c, 0x86, 0x2d, 0x95, 0x25, 0x5a, 0xd1, 0x95, 0x09, 0x44, 0x61, 0x5c, 0xf4, 0x9a,
	0xe7, 0xcb, 0x82, 0xc9, 0x8a, 0xbe, 0x3c, 0x83, 0x93, 0x89, 0x72, 0x4e, 0xd7, 0xb8, 0x03, 0x47,
	0xcd, 0xde, 0x57, 0xc5, 0xaa, 0x1d, 0x88, 0x8b, 0x0c, 0x2e, 0x0d, 0x72, 0xb8, 0xee, 0x61, 0xf7,
	0x3a, 0x87, 0xb6, 0x46, 0xfb, 0xd1, 0x53, 0x95, 0x9e, 0xc7, 0xe1, 0x4b, 0x9e, 0x7e, 0x02, 0x00,
	0x00, 0xff, 0xff, 0xe1, 0xee, 0xc1, 0x50, 0xa1, 0x01, 0x00, 0x00,
}
