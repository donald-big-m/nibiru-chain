// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: vpool/v1/event.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "github.com/regen-network/cosmos-proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ReserveSnapshotSavedEvent struct {
	Pair         string                                 `protobuf:"bytes,1,opt,name=pair,proto3" json:"pair,omitempty"`
	QuoteReserve github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=quote_reserve,json=quoteReserve,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"quote_reserve"`
	BaseReserve  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=base_reserve,json=baseReserve,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"base_reserve"`
	// MarkPrice at the end of the block.
	// (instantaneous) markPrice := quoteReserve / baseReserve
	MarkPrice      github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=mark_price,json=markPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"mark_price"`
	BlockHeight    int64                                  `protobuf:"varint,5,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	BlockTimestamp time.Time                              `protobuf:"bytes,6,opt,name=block_timestamp,json=blockTimestamp,proto3,stdtime" json:"block_timestamp"`
}

func (m *ReserveSnapshotSavedEvent) Reset()         { *m = ReserveSnapshotSavedEvent{} }
func (m *ReserveSnapshotSavedEvent) String() string { return proto.CompactTextString(m) }
func (*ReserveSnapshotSavedEvent) ProtoMessage()    {}
func (*ReserveSnapshotSavedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_faeff0bc76489252, []int{0}
}
func (m *ReserveSnapshotSavedEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReserveSnapshotSavedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReserveSnapshotSavedEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReserveSnapshotSavedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReserveSnapshotSavedEvent.Merge(m, src)
}
func (m *ReserveSnapshotSavedEvent) XXX_Size() int {
	return m.Size()
}
func (m *ReserveSnapshotSavedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ReserveSnapshotSavedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ReserveSnapshotSavedEvent proto.InternalMessageInfo

func (m *ReserveSnapshotSavedEvent) GetPair() string {
	if m != nil {
		return m.Pair
	}
	return ""
}

func (m *ReserveSnapshotSavedEvent) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *ReserveSnapshotSavedEvent) GetBlockTimestamp() time.Time {
	if m != nil {
		return m.BlockTimestamp
	}
	return time.Time{}
}

type SwapQuoteForBaseEvent struct {
	Pair        string                                 `protobuf:"bytes,1,opt,name=pair,proto3" json:"pair,omitempty"`
	QuoteAmount github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=quote_amount,json=quoteAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"quote_amount"`
	BaseAmount  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=base_amount,json=baseAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"base_amount"`
}

func (m *SwapQuoteForBaseEvent) Reset()         { *m = SwapQuoteForBaseEvent{} }
func (m *SwapQuoteForBaseEvent) String() string { return proto.CompactTextString(m) }
func (*SwapQuoteForBaseEvent) ProtoMessage()    {}
func (*SwapQuoteForBaseEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_faeff0bc76489252, []int{1}
}
func (m *SwapQuoteForBaseEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SwapQuoteForBaseEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SwapQuoteForBaseEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SwapQuoteForBaseEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwapQuoteForBaseEvent.Merge(m, src)
}
func (m *SwapQuoteForBaseEvent) XXX_Size() int {
	return m.Size()
}
func (m *SwapQuoteForBaseEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_SwapQuoteForBaseEvent.DiscardUnknown(m)
}

var xxx_messageInfo_SwapQuoteForBaseEvent proto.InternalMessageInfo

func (m *SwapQuoteForBaseEvent) GetPair() string {
	if m != nil {
		return m.Pair
	}
	return ""
}

type SwapBaseForQuoteEvent struct {
	Pair        string                                 `protobuf:"bytes,1,opt,name=pair,proto3" json:"pair,omitempty"`
	QuoteAmount github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=quote_amount,json=quoteAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"quote_amount"`
	BaseAmount  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=base_amount,json=baseAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"base_amount"`
}

func (m *SwapBaseForQuoteEvent) Reset()         { *m = SwapBaseForQuoteEvent{} }
func (m *SwapBaseForQuoteEvent) String() string { return proto.CompactTextString(m) }
func (*SwapBaseForQuoteEvent) ProtoMessage()    {}
func (*SwapBaseForQuoteEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_faeff0bc76489252, []int{2}
}
func (m *SwapBaseForQuoteEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SwapBaseForQuoteEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SwapBaseForQuoteEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SwapBaseForQuoteEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwapBaseForQuoteEvent.Merge(m, src)
}
func (m *SwapBaseForQuoteEvent) XXX_Size() int {
	return m.Size()
}
func (m *SwapBaseForQuoteEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_SwapBaseForQuoteEvent.DiscardUnknown(m)
}

var xxx_messageInfo_SwapBaseForQuoteEvent proto.InternalMessageInfo

func (m *SwapBaseForQuoteEvent) GetPair() string {
	if m != nil {
		return m.Pair
	}
	return ""
}

type MarkPriceChangedEvent struct {
	Pair      string                                 `protobuf:"bytes,1,opt,name=pair,proto3" json:"pair,omitempty"`
	Price     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=price,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"price"`
	Timestamp time.Time                              `protobuf:"bytes,3,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
}

func (m *MarkPriceChangedEvent) Reset()         { *m = MarkPriceChangedEvent{} }
func (m *MarkPriceChangedEvent) String() string { return proto.CompactTextString(m) }
func (*MarkPriceChangedEvent) ProtoMessage()    {}
func (*MarkPriceChangedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_faeff0bc76489252, []int{3}
}
func (m *MarkPriceChangedEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarkPriceChangedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarkPriceChangedEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarkPriceChangedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarkPriceChangedEvent.Merge(m, src)
}
func (m *MarkPriceChangedEvent) XXX_Size() int {
	return m.Size()
}
func (m *MarkPriceChangedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_MarkPriceChangedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_MarkPriceChangedEvent proto.InternalMessageInfo

func (m *MarkPriceChangedEvent) GetPair() string {
	if m != nil {
		return m.Pair
	}
	return ""
}

func (m *MarkPriceChangedEvent) GetTimestamp() time.Time {
	if m != nil {
		return m.Timestamp
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*ReserveSnapshotSavedEvent)(nil), "nibiru.vpool.v1.ReserveSnapshotSavedEvent")
	proto.RegisterType((*SwapQuoteForBaseEvent)(nil), "nibiru.vpool.v1.SwapQuoteForBaseEvent")
	proto.RegisterType((*SwapBaseForQuoteEvent)(nil), "nibiru.vpool.v1.SwapBaseForQuoteEvent")
	proto.RegisterType((*MarkPriceChangedEvent)(nil), "nibiru.vpool.v1.MarkPriceChangedEvent")
}

func init() { proto.RegisterFile("vpool/v1/event.proto", fileDescriptor_faeff0bc76489252) }

var fileDescriptor_faeff0bc76489252 = []byte{
	// 494 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x94, 0x4d, 0x6f, 0xd3, 0x30,
	0x18, 0xc7, 0x6b, 0xba, 0x4d, 0xd4, 0x2d, 0x4c, 0x8a, 0x36, 0x29, 0xeb, 0x21, 0x2d, 0x3b, 0xa0,
	0x4a, 0x88, 0x58, 0x85, 0x4f, 0x40, 0xf6, 0x22, 0x2e, 0x05, 0x9a, 0x72, 0xe2, 0x12, 0x39, 0x99,
	0x49, 0xa2, 0x36, 0x79, 0x82, 0xed, 0x04, 0xf8, 0x16, 0x3b, 0xf2, 0x45, 0xf8, 0x0e, 0xbb, 0xb1,
	0x23, 0xe2, 0x30, 0x50, 0xfb, 0x45, 0x90, 0xed, 0xa4, 0xe3, 0x02, 0x12, 0xb9, 0xed, 0x14, 0xdb,
	0xff, 0xe7, 0xf9, 0xf9, 0x79, 0x8b, 0xf1, 0x41, 0x55, 0x00, 0xac, 0x48, 0x35, 0x25, 0xac, 0x62,
	0xb9, 0x74, 0x0b, 0x0e, 0x12, 0xac, 0xfd, 0x3c, 0x0d, 0x53, 0x5e, 0xba, 0x5a, 0x74, 0xab, 0xe9,
	0xf0, 0x20, 0x86, 0x18, 0xb4, 0x46, 0xd4, 0xca, 0x98, 0x0d, 0x9d, 0x08, 0x44, 0x06, 0x82, 0x84,
	0x54, 0x30, 0x52, 0x4d, 0x43, 0x26, 0xe9, 0x94, 0x44, 0x90, 0xe6, 0xb5, 0x7e, 0x64, 0xf4, 0xc0,
	0x38, 0x9a, 0x4d, 0x2d, 0x8d, 0x62, 0x80, 0x78, 0xc5, 0x88, 0xde, 0x85, 0xe5, 0x7b, 0x22, 0xd3,
	0x8c, 0x09, 0x49, 0xb3, 0xc2, 0x18, 0x1c, 0x7f, 0xe9, 0xe2, 0x23, 0x9f, 0x09, 0xc6, 0x2b, 0xb6,
	0xc8, 0x69, 0x21, 0x12, 0x90, 0x0b, 0x5a, 0xb1, 0x8b, 0x33, 0x15, 0xa6, 0x65, 0xe1, 0x9d, 0x82,
	0xa6, 0xdc, 0x46, 0x63, 0x34, 0xe9, 0xf9, 0x7a, 0x6d, 0x2d, 0xf0, 0x83, 0x0f, 0x25, 0x48, 0x16,
	0x70, 0xe3, 0x66, 0xdf, 0x53, 0xa2, 0xe7, 0x5e, 0xdd, 0x8c, 0x3a, 0x3f, 0x6e, 0x46, 0x8f, 0xe3,
	0x54, 0x26, 0x65, 0xe8, 0x46, 0x90, 0xd5, 0xa1, 0xd4, 0x9f, 0xa7, 0xe2, 0x62, 0x49, 0xe4, 0xe7,
	0x82, 0x09, 0xf7, 0x94, 0x45, 0xfe, 0x40, 0x43, 0xea, 0xab, 0xad, 0x39, 0x1e, 0xa8, 0xec, 0xb6,
	0xcc, 0x6e, 0x2b, 0x66, 0x5f, 0x31, 0x1a, 0xe4, 0x0c, 0xe3, 0x8c, 0xf2, 0x65, 0x50, 0xf0, 0x34,
	0x62, 0xf6, 0x4e, 0x2b, 0x60, 0x4f, 0x11, 0xde, 0x28, 0x80, 0xf5, 0x08, 0x0f, 0xc2, 0x15, 0x44,
	0xcb, 0x20, 0x61, 0x69, 0x9c, 0x48, 0x7b, 0x77, 0x8c, 0x26, 0x5d, 0xbf, 0xaf, 0xcf, 0x5e, 0xea,
	0x23, 0x6b, 0x86, 0xf7, 0x8d, 0xc9, 0xb6, 0xc8, 0xf6, 0xde, 0x18, 0x4d, 0xfa, 0xcf, 0x86, 0xae,
	0x69, 0x83, 0xdb, 0xb4, 0xc1, 0x7d, 0xdb, 0x58, 0x78, 0xf7, 0x55, 0x48, 0x97, 0x3f, 0x47, 0xc8,
	0x7f, 0xa8, 0x9d, 0xb7, 0xca, 0xf1, 0x37, 0x84, 0x0f, 0x17, 0x1f, 0x69, 0x31, 0x57, 0x85, 0x3a,
	0x07, 0xee, 0x51, 0xc1, 0xfe, 0xde, 0x96, 0x39, 0x36, 0x15, 0x0d, 0x68, 0x06, 0x65, 0x2e, 0x5b,
	0x76, 0xa5, 0xaf, 0x19, 0x2f, 0x34, 0xc2, 0x7a, 0x8d, 0x75, 0x41, 0x1b, 0x62, 0xbb, 0x9e, 0x60,
	0x85, 0x30, 0xc0, 0x6d, 0x46, 0x2a, 0x93, 0x73, 0xe0, 0x3a, 0xb1, 0xbb, 0x9d, 0xd1, 0x57, 0x84,
	0x0f, 0x67, 0xcd, 0x8c, 0x9c, 0x24, 0x34, 0x8f, 0xff, 0xf5, 0xeb, 0x9c, 0xe2, 0x5d, 0x33, 0x8d,
	0xed, 0x52, 0x31, 0xce, 0x96, 0x87, 0x7b, 0xb7, 0x03, 0xd6, 0xfd, 0x8f, 0x01, 0xbb, 0x75, 0xf3,
	0xce, 0xae, 0xd6, 0x0e, 0xba, 0x5e, 0x3b, 0xe8, 0xd7, 0xda, 0x41, 0x97, 0x1b, 0xa7, 0x73, 0xbd,
	0x71, 0x3a, 0xdf, 0x37, 0x4e, 0xe7, 0xdd, 0x93, 0x3f, 0x82, 0x79, 0xa5, 0x9f, 0xa7, 0x93, 0x84,
	0xa6, 0x39, 0x31, 0x4f, 0x15, 0xf9, 0x44, 0xcc, 0x4b, 0xa6, 0xa3, 0x0a, 0xf7, 0xf4, 0x7d, 0xcf,
	0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0x70, 0x41, 0x2a, 0xeb, 0xdf, 0x04, 0x00, 0x00,
}

func (m *ReserveSnapshotSavedEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReserveSnapshotSavedEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReserveSnapshotSavedEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.BlockTimestamp, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.BlockTimestamp):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintEvent(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x32
	if m.BlockHeight != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.MarkPrice.Size()
		i -= size
		if _, err := m.MarkPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.BaseReserve.Size()
		i -= size
		if _, err := m.BaseReserve.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.QuoteReserve.Size()
		i -= size
		if _, err := m.QuoteReserve.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Pair) > 0 {
		i -= len(m.Pair)
		copy(dAtA[i:], m.Pair)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Pair)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SwapQuoteForBaseEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SwapQuoteForBaseEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SwapQuoteForBaseEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.BaseAmount.Size()
		i -= size
		if _, err := m.BaseAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.QuoteAmount.Size()
		i -= size
		if _, err := m.QuoteAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Pair) > 0 {
		i -= len(m.Pair)
		copy(dAtA[i:], m.Pair)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Pair)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SwapBaseForQuoteEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SwapBaseForQuoteEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SwapBaseForQuoteEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.BaseAmount.Size()
		i -= size
		if _, err := m.BaseAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.QuoteAmount.Size()
		i -= size
		if _, err := m.QuoteAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Pair) > 0 {
		i -= len(m.Pair)
		copy(dAtA[i:], m.Pair)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Pair)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MarkPriceChangedEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarkPriceChangedEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarkPriceChangedEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n2, err2 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.Timestamp, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.Timestamp):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintEvent(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x1a
	{
		size := m.Price.Size()
		i -= size
		if _, err := m.Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Pair) > 0 {
		i -= len(m.Pair)
		copy(dAtA[i:], m.Pair)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Pair)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ReserveSnapshotSavedEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Pair)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = m.QuoteReserve.Size()
	n += 1 + l + sovEvent(uint64(l))
	l = m.BaseReserve.Size()
	n += 1 + l + sovEvent(uint64(l))
	l = m.MarkPrice.Size()
	n += 1 + l + sovEvent(uint64(l))
	if m.BlockHeight != 0 {
		n += 1 + sovEvent(uint64(m.BlockHeight))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.BlockTimestamp)
	n += 1 + l + sovEvent(uint64(l))
	return n
}

func (m *SwapQuoteForBaseEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Pair)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = m.QuoteAmount.Size()
	n += 1 + l + sovEvent(uint64(l))
	l = m.BaseAmount.Size()
	n += 1 + l + sovEvent(uint64(l))
	return n
}

func (m *SwapBaseForQuoteEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Pair)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = m.QuoteAmount.Size()
	n += 1 + l + sovEvent(uint64(l))
	l = m.BaseAmount.Size()
	n += 1 + l + sovEvent(uint64(l))
	return n
}

func (m *MarkPriceChangedEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Pair)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = m.Price.Size()
	n += 1 + l + sovEvent(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.Timestamp)
	n += 1 + l + sovEvent(uint64(l))
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ReserveSnapshotSavedEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ReserveSnapshotSavedEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReserveSnapshotSavedEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pair", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pair = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QuoteReserve", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.QuoteReserve.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseReserve", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BaseReserve.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MarkPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MarkPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.BlockTimestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SwapQuoteForBaseEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SwapQuoteForBaseEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SwapQuoteForBaseEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pair", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pair = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QuoteAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.QuoteAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BaseAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SwapBaseForQuoteEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SwapBaseForQuoteEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SwapBaseForQuoteEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pair", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pair = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QuoteAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.QuoteAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BaseAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MarkPriceChangedEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MarkPriceChangedEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarkPriceChangedEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pair", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pair = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.Timestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
