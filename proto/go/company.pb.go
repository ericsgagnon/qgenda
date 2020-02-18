// Code generated by protoc-gen-go. DO NOT EDIT.
// source: company.proto

package qgenda

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Company struct {
	Id                   string               `protobuf:"bytes,1,opt,name=id,json=test,proto3" json:"id,omitempty"`
	Name                 string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Abbreviation         string               `protobuf:"bytes,3,opt,name=abbreviation,proto3" json:"abbreviation,omitempty"`
	Createdtime          *timestamp.Timestamp `protobuf:"bytes,4,opt,name=createdtime,proto3" json:"createdtime,omitempty"`
	Location             string               `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	Phonenumber          string               `protobuf:"bytes,6,opt,name=phonenumber,proto3" json:"phonenumber,omitempty"`
	Profile              []*Profile           `protobuf:"bytes,7,rep,name=profile,proto3" json:"profile,omitempty"`
	Organization         []*Organization      `protobuf:"bytes,8,rep,name=organization,proto3" json:"organization,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Company) Reset()         { *m = Company{} }
func (m *Company) String() string { return proto.CompactTextString(m) }
func (*Company) ProtoMessage()    {}
func (*Company) Descriptor() ([]byte, []int) {
	return fileDescriptor_ade57ca5b8f3903f, []int{0}
}

func (m *Company) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Company.Unmarshal(m, b)
}
func (m *Company) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Company.Marshal(b, m, deterministic)
}
func (m *Company) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Company.Merge(m, src)
}
func (m *Company) XXX_Size() int {
	return xxx_messageInfo_Company.Size(m)
}
func (m *Company) XXX_DiscardUnknown() {
	xxx_messageInfo_Company.DiscardUnknown(m)
}

var xxx_messageInfo_Company proto.InternalMessageInfo

func (m *Company) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Company) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Company) GetAbbreviation() string {
	if m != nil {
		return m.Abbreviation
	}
	return ""
}

func (m *Company) GetCreatedtime() *timestamp.Timestamp {
	if m != nil {
		return m.Createdtime
	}
	return nil
}

func (m *Company) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *Company) GetPhonenumber() string {
	if m != nil {
		return m.Phonenumber
	}
	return ""
}

func (m *Company) GetProfile() []*Profile {
	if m != nil {
		return m.Profile
	}
	return nil
}

func (m *Company) GetOrganization() []*Organization {
	if m != nil {
		return m.Organization
	}
	return nil
}

// Profile appears to link a user role to a company...
type Profile struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Admin                bool     `protobuf:"varint,3,opt,name=admin,proto3" json:"admin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Profile) Reset()         { *m = Profile{} }
func (m *Profile) String() string { return proto.CompactTextString(m) }
func (*Profile) ProtoMessage()    {}
func (*Profile) Descriptor() ([]byte, []int) {
	return fileDescriptor_ade57ca5b8f3903f, []int{1}
}

func (m *Profile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Profile.Unmarshal(m, b)
}
func (m *Profile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Profile.Marshal(b, m, deterministic)
}
func (m *Profile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Profile.Merge(m, src)
}
func (m *Profile) XXX_Size() int {
	return xxx_messageInfo_Profile.Size(m)
}
func (m *Profile) XXX_DiscardUnknown() {
	xxx_messageInfo_Profile.DiscardUnknown(m)
}

var xxx_messageInfo_Profile proto.InternalMessageInfo

func (m *Profile) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Profile) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Profile) GetAdmin() bool {
	if m != nil {
		return m.Admin
	}
	return false
}

// Organization appears to linke multiple companies and users
type Organization struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Organization) Reset()         { *m = Organization{} }
func (m *Organization) String() string { return proto.CompactTextString(m) }
func (*Organization) ProtoMessage()    {}
func (*Organization) Descriptor() ([]byte, []int) {
	return fileDescriptor_ade57ca5b8f3903f, []int{2}
}

func (m *Organization) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Organization.Unmarshal(m, b)
}
func (m *Organization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Organization.Marshal(b, m, deterministic)
}
func (m *Organization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Organization.Merge(m, src)
}
func (m *Organization) XXX_Size() int {
	return xxx_messageInfo_Organization.Size(m)
}
func (m *Organization) XXX_DiscardUnknown() {
	xxx_messageInfo_Organization.DiscardUnknown(m)
}

var xxx_messageInfo_Organization proto.InternalMessageInfo

func (m *Organization) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Organization) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func init() {
	proto.RegisterType((*Company)(nil), "qgenda.Company")
	proto.RegisterType((*Profile)(nil), "qgenda.Profile")
	proto.RegisterType((*Organization)(nil), "qgenda.Organization")
}

func init() { proto.RegisterFile("company.proto", fileDescriptor_ade57ca5b8f3903f) }

var fileDescriptor_ade57ca5b8f3903f = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x95, 0xfe, 0x4b, 0xb9, 0x04, 0x51, 0x59, 0x1d, 0xac, 0x2c, 0x44, 0x99, 0xc2, 0xe2,
	0x4a, 0x85, 0x81, 0x81, 0x0d, 0x31, 0x83, 0x22, 0xbe, 0x80, 0x93, 0x5c, 0x83, 0x45, 0x6c, 0x07,
	0xd7, 0x45, 0x2a, 0x9f, 0x8b, 0x0f, 0x88, 0x62, 0x37, 0x55, 0xba, 0xb1, 0xf9, 0xde, 0x7b, 0x7e,
	0x77, 0xfa, 0xc1, 0x75, 0xa5, 0x65, 0xc7, 0xd5, 0x91, 0x75, 0x46, 0x5b, 0x4d, 0x16, 0x5f, 0x0d,
	0xaa, 0x9a, 0x27, 0xb7, 0x8d, 0xd6, 0x4d, 0x8b, 0x1b, 0xa7, 0x96, 0x87, 0xdd, 0xc6, 0x0a, 0x89,
	0x7b, 0xcb, 0x65, 0xe7, 0x83, 0xd9, 0xef, 0x04, 0xc2, 0x67, 0xff, 0x95, 0xac, 0x60, 0x22, 0x6a,
	0x1a, 0xa4, 0x41, 0x7e, 0x55, 0xcc, 0x2c, 0xee, 0x2d, 0x21, 0x30, 0x53, 0x5c, 0x22, 0x9d, 0x78,
	0xad, 0x7f, 0x93, 0x0c, 0x62, 0x5e, 0x96, 0x06, 0xbf, 0x05, 0xb7, 0x42, 0x2b, 0x3a, 0x75, 0xde,
	0x85, 0x46, 0x9e, 0x20, 0xaa, 0x0c, 0x72, 0x8b, 0x75, 0xbf, 0x8f, 0xce, 0xd2, 0x20, 0x8f, 0xb6,
	0x09, 0xf3, 0xc7, 0xb0, 0xe1, 0x18, 0xf6, 0x3e, 0x1c, 0x53, 0x8c, 0xe3, 0x24, 0x81, 0x65, 0xab,
	0x2b, 0xdf, 0x3e, 0x77, 0xed, 0xe7, 0x99, 0xa4, 0x10, 0x75, 0x1f, 0x5a, 0xa1, 0x3a, 0xc8, 0x12,
	0x0d, 0x5d, 0x38, 0x7b, 0x2c, 0x91, 0x3b, 0x08, 0x3b, 0xa3, 0x77, 0xa2, 0x45, 0x1a, 0xa6, 0xd3,
	0x3c, 0xda, 0xde, 0x30, 0x0f, 0x83, 0xbd, 0x79, 0xb9, 0x18, 0x7c, 0xf2, 0x08, 0xb1, 0x36, 0x0d,
	0x57, 0xe2, 0xc7, 0x2f, 0x5b, 0xba, 0xfc, 0x7a, 0xc8, 0xbf, 0x8e, 0xbc, 0xe2, 0x22, 0x99, 0xbd,
	0x40, 0x78, 0x6a, 0x3b, 0x33, 0x0a, 0x46, 0x8c, 0x56, 0x30, 0xfd, 0xc4, 0xe3, 0x09, 0x5b, 0xff,
	0x24, 0x6b, 0x98, 0xf3, 0x5a, 0x0a, 0x8f, 0x6b, 0x59, 0xf8, 0x21, 0x7b, 0x80, 0x78, 0xbc, 0xe4,
	0x7f, 0x5d, 0xe5, 0xc2, 0x01, 0xbc, 0xff, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x16, 0xfa, 0x7a, 0xd4,
	0xf4, 0x01, 0x00, 0x00,
}
