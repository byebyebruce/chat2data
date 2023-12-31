// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types.proto

package localvectordb

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Doc
type Doc struct {
	Id                   string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Content              string    `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Vector               []float32 `protobuf:"fixed32,3,rep,packed,name=vector,proto3" json:"vector,omitempty"`
	Meta                 []byte    `protobuf:"bytes,4,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Doc) Reset()         { *m = Doc{} }
func (m *Doc) String() string { return proto.CompactTextString(m) }
func (*Doc) ProtoMessage()    {}
func (*Doc) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}
func (m *Doc) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Doc.Unmarshal(m, b)
}
func (m *Doc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Doc.Marshal(b, m, deterministic)
}
func (m *Doc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Doc.Merge(m, src)
}
func (m *Doc) XXX_Size() int {
	return xxx_messageInfo_Doc.Size(m)
}
func (m *Doc) XXX_DiscardUnknown() {
	xxx_messageInfo_Doc.DiscardUnknown(m)
}

var xxx_messageInfo_Doc proto.InternalMessageInfo

func (m *Doc) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Doc) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Doc) GetVector() []float32 {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Doc) GetMeta() []byte {
	if m != nil {
		return m.Meta
	}
	return nil
}

func init() {
	proto.RegisterType((*Doc)(nil), "localvectordb.Doc")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor_d938547f84707355) }

var fileDescriptor_d938547f84707355 = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0xa9, 0x2c, 0x48,
	0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xcd, 0xc9, 0x4f, 0x4e, 0xcc, 0x29, 0x4b,
	0x4d, 0x2e, 0xc9, 0x2f, 0x4a, 0x49, 0x52, 0x8a, 0xe6, 0x62, 0x76, 0xc9, 0x4f, 0x16, 0xe2, 0xe3,
	0x62, 0xca, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x62, 0xca, 0x4c, 0x11, 0x92, 0xe0,
	0x62, 0x4f, 0xce, 0xcf, 0x2b, 0x49, 0xcd, 0x2b, 0x91, 0x60, 0x02, 0x0b, 0xc2, 0xb8, 0x42, 0x62,
	0x5c, 0x6c, 0x10, 0xcd, 0x12, 0xcc, 0x0a, 0xcc, 0x1a, 0x4c, 0x41, 0x50, 0x9e, 0x90, 0x10, 0x17,
	0x4b, 0x6e, 0x6a, 0x49, 0xa2, 0x04, 0x8b, 0x02, 0xa3, 0x06, 0x4f, 0x10, 0x98, 0xed, 0x64, 0x17,
	0x65, 0x93, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x9f, 0x54, 0x99, 0x0a,
	0x42, 0x45, 0xa5, 0xc9, 0xa9, 0xfa, 0xc9, 0x19, 0x89, 0x25, 0x46, 0x29, 0x89, 0x25, 0x89, 0xfa,
	0x28, 0xce, 0xb1, 0x46, 0xe1, 0x25, 0xb1, 0x81, 0x9d, 0x6c, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff,
	0xaa, 0x92, 0xe6, 0xd8, 0xc1, 0x00, 0x00, 0x00,
}
