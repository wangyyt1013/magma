// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cwf/protos/mconfig/mconfigs.proto

package mconfig

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

//-----------------------------------------------------------------------------
// Health configs
//-----------------------------------------------------------------------------
type CwfGatewayHealthConfig struct {
	// cpu utilization threshold
	CpuUtilThresholdPct float32 `protobuf:"fixed32,1,opt,name=cpu_util_threshold_pct,json=cpuUtilThresholdPct,proto3" json:"cpu_util_threshold_pct,omitempty"`
	// mem utilization threshold
	MemUtilThresholdPct float32 `protobuf:"fixed32,2,opt,name=mem_util_threshold_pct,json=memUtilThresholdPct,proto3" json:"mem_util_threshold_pct,omitempty"`
	// interval between probes
	GreProbeInterval uint32 `protobuf:"varint,3,opt,name=gre_probe_interval,json=greProbeInterval,proto3" json:"gre_probe_interval,omitempty"`
	// packets sent for each icmp probe
	IcmpProbePktCount uint32 `protobuf:"varint,4,opt,name=icmp_probe_pkt_count,json=icmpProbePktCount,proto3" json:"icmp_probe_pkt_count,omitempty"`
	// gre peers to probe
	GrePeers []*CwfGatewayHealthConfigGrePeer `protobuf:"bytes,5,rep,name=gre_peers,json=grePeers,proto3" json:"gre_peers,omitempty"`
	// virtual IP used by AP/WLC to connect to HA cluster
	ClusterVirtualIp     string   `protobuf:"bytes,6,opt,name=cluster_virtual_ip,json=clusterVirtualIp,proto3" json:"cluster_virtual_ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CwfGatewayHealthConfig) Reset()         { *m = CwfGatewayHealthConfig{} }
func (m *CwfGatewayHealthConfig) String() string { return proto.CompactTextString(m) }
func (*CwfGatewayHealthConfig) ProtoMessage()    {}
func (*CwfGatewayHealthConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab79e679bf56b47d, []int{0}
}

func (m *CwfGatewayHealthConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CwfGatewayHealthConfig.Unmarshal(m, b)
}
func (m *CwfGatewayHealthConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CwfGatewayHealthConfig.Marshal(b, m, deterministic)
}
func (m *CwfGatewayHealthConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CwfGatewayHealthConfig.Merge(m, src)
}
func (m *CwfGatewayHealthConfig) XXX_Size() int {
	return xxx_messageInfo_CwfGatewayHealthConfig.Size(m)
}
func (m *CwfGatewayHealthConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_CwfGatewayHealthConfig.DiscardUnknown(m)
}

var xxx_messageInfo_CwfGatewayHealthConfig proto.InternalMessageInfo

func (m *CwfGatewayHealthConfig) GetCpuUtilThresholdPct() float32 {
	if m != nil {
		return m.CpuUtilThresholdPct
	}
	return 0
}

func (m *CwfGatewayHealthConfig) GetMemUtilThresholdPct() float32 {
	if m != nil {
		return m.MemUtilThresholdPct
	}
	return 0
}

func (m *CwfGatewayHealthConfig) GetGreProbeInterval() uint32 {
	if m != nil {
		return m.GreProbeInterval
	}
	return 0
}

func (m *CwfGatewayHealthConfig) GetIcmpProbePktCount() uint32 {
	if m != nil {
		return m.IcmpProbePktCount
	}
	return 0
}

func (m *CwfGatewayHealthConfig) GetGrePeers() []*CwfGatewayHealthConfigGrePeer {
	if m != nil {
		return m.GrePeers
	}
	return nil
}

func (m *CwfGatewayHealthConfig) GetClusterVirtualIp() string {
	if m != nil {
		return m.ClusterVirtualIp
	}
	return ""
}

type CwfGatewayHealthConfigGrePeer struct {
	Ip                   string   `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CwfGatewayHealthConfigGrePeer) Reset()         { *m = CwfGatewayHealthConfigGrePeer{} }
func (m *CwfGatewayHealthConfigGrePeer) String() string { return proto.CompactTextString(m) }
func (*CwfGatewayHealthConfigGrePeer) ProtoMessage()    {}
func (*CwfGatewayHealthConfigGrePeer) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab79e679bf56b47d, []int{0, 0}
}

func (m *CwfGatewayHealthConfigGrePeer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CwfGatewayHealthConfigGrePeer.Unmarshal(m, b)
}
func (m *CwfGatewayHealthConfigGrePeer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CwfGatewayHealthConfigGrePeer.Marshal(b, m, deterministic)
}
func (m *CwfGatewayHealthConfigGrePeer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CwfGatewayHealthConfigGrePeer.Merge(m, src)
}
func (m *CwfGatewayHealthConfigGrePeer) XXX_Size() int {
	return xxx_messageInfo_CwfGatewayHealthConfigGrePeer.Size(m)
}
func (m *CwfGatewayHealthConfigGrePeer) XXX_DiscardUnknown() {
	xxx_messageInfo_CwfGatewayHealthConfigGrePeer.DiscardUnknown(m)
}

var xxx_messageInfo_CwfGatewayHealthConfigGrePeer proto.InternalMessageInfo

func (m *CwfGatewayHealthConfigGrePeer) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func init() {
	proto.RegisterType((*CwfGatewayHealthConfig)(nil), "magma.mconfig.CwfGatewayHealthConfig")
	proto.RegisterType((*CwfGatewayHealthConfigGrePeer)(nil), "magma.mconfig.CwfGatewayHealthConfig.grePeer")
}

func init() { proto.RegisterFile("cwf/protos/mconfig/mconfigs.proto", fileDescriptor_ab79e679bf56b47d) }

var fileDescriptor_ab79e679bf56b47d = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x31, 0x4f, 0xeb, 0x30,
	0x10, 0xc7, 0x95, 0xf4, 0xbd, 0xbe, 0x57, 0xa3, 0xa2, 0x12, 0x50, 0x15, 0x98, 0x52, 0x58, 0x32,
	0x40, 0x22, 0xd1, 0x6f, 0x40, 0x07, 0x28, 0x53, 0x15, 0x01, 0x03, 0x8b, 0xe5, 0xba, 0xd7, 0xd4,
	0xaa, 0x5d, 0x5b, 0xce, 0xb9, 0x15, 0x5f, 0x9c, 0x19, 0xd9, 0x4d, 0x07, 0xa0, 0x53, 0xe2, 0xfb,
	0xdd, 0xcf, 0x3e, 0xfd, 0x8f, 0x8c, 0xf8, 0x6e, 0x59, 0x1a, 0xab, 0x51, 0x37, 0xa5, 0xe2, 0x7a,
	0xb3, 0x14, 0xf5, 0xe1, 0xdb, 0x14, 0xa1, 0x9e, 0xf4, 0x15, 0xab, 0x15, 0x2b, 0xda, 0xea, 0xf5,
	0x67, 0x4c, 0x86, 0x93, 0xdd, 0xf2, 0x91, 0x21, 0xec, 0xd8, 0xc7, 0x13, 0x30, 0x89, 0xab, 0x49,
	0x40, 0xc9, 0x98, 0x0c, 0xb9, 0x71, 0xd4, 0xa1, 0x90, 0x14, 0x57, 0x16, 0x9a, 0x95, 0x96, 0x0b,
	0x6a, 0x38, 0xa6, 0x51, 0x16, 0xe5, 0x71, 0x75, 0xce, 0x8d, 0x7b, 0x45, 0x21, 0x5f, 0x0e, 0x6c,
	0xc6, 0xd1, 0x4b, 0x0a, 0xd4, 0x31, 0x29, 0xde, 0x4b, 0x0a, 0xd4, 0x2f, 0xe9, 0x96, 0x24, 0xb5,
	0x05, 0x6a, 0xac, 0x9e, 0x03, 0x15, 0x1b, 0x04, 0xbb, 0x65, 0x32, 0xed, 0x64, 0x51, 0xde, 0xaf,
	0x06, 0xb5, 0x85, 0x99, 0x07, 0xd3, 0xb6, 0x9e, 0x94, 0xe4, 0x42, 0x70, 0x65, 0xda, 0x76, 0xb3,
	0x46, 0xca, 0xb5, 0xdb, 0x60, 0xfa, 0x27, 0xf4, 0x9f, 0x79, 0x16, 0x84, 0xd9, 0x1a, 0x27, 0x1e,
	0x24, 0xcf, 0xa4, 0x17, 0xae, 0x07, 0xb0, 0x4d, 0xfa, 0x37, 0xeb, 0xe4, 0x27, 0xf7, 0x77, 0xc5,
	0xb7, 0x18, 0x8a, 0xe3, 0x11, 0x14, 0xfe, 0x6d, 0x00, 0x5b, 0xfd, 0x6f, 0x7f, 0x1a, 0x3f, 0x2a,
	0x97, 0xae, 0x41, 0xb0, 0x74, 0x2b, 0x2c, 0x3a, 0x26, 0xa9, 0x30, 0x69, 0x37, 0x8b, 0xf2, 0x5e,
	0x35, 0x68, 0xc9, 0xdb, 0x1e, 0x4c, 0xcd, 0xd5, 0x25, 0xf9, 0xd7, 0x9a, 0xc9, 0x29, 0x89, 0x85,
	0x09, 0xc9, 0xf5, 0xaa, 0x58, 0x98, 0x87, 0x9b, 0xf7, 0x51, 0x18, 0xa1, 0xf4, 0x2b, 0xe3, 0x52,
	0xbb, 0x45, 0x59, 0xeb, 0x1f, 0xbb, 0x9b, 0x77, 0xc3, 0x79, 0xfc, 0x15, 0x00, 0x00, 0xff, 0xff,
	0x7e, 0xc1, 0xde, 0x59, 0xd8, 0x01, 0x00, 0x00,
}
