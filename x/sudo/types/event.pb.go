// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: nibiru/sudo/v1/event.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	io "io"
	math "math"
	math_bits "math/bits"
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

type EventUpdateSudoers struct {
	Sudoers Sudoers `protobuf:"bytes,1,opt,name=sudoers,proto3" json:"sudoers"`
	// Action is the type of update that occured to the "sudoers"
	Action string `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
}

func (m *EventUpdateSudoers) Reset()         { *m = EventUpdateSudoers{} }
func (m *EventUpdateSudoers) String() string { return proto.CompactTextString(m) }
func (*EventUpdateSudoers) ProtoMessage()    {}
func (*EventUpdateSudoers) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e6085948b018986, []int{0}
}
func (m *EventUpdateSudoers) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventUpdateSudoers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventUpdateSudoers.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventUpdateSudoers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventUpdateSudoers.Merge(m, src)
}
func (m *EventUpdateSudoers) XXX_Size() int {
	return m.Size()
}
func (m *EventUpdateSudoers) XXX_DiscardUnknown() {
	xxx_messageInfo_EventUpdateSudoers.DiscardUnknown(m)
}

var xxx_messageInfo_EventUpdateSudoers proto.InternalMessageInfo

func (m *EventUpdateSudoers) GetSudoers() Sudoers {
	if m != nil {
		return m.Sudoers
	}
	return Sudoers{}
}

func (m *EventUpdateSudoers) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func init() {
	proto.RegisterType((*EventUpdateSudoers)(nil), "nibiru.sudo.v1.EventUpdateSudoers")
}

func init() { proto.RegisterFile("nibiru/sudo/v1/event.proto", fileDescriptor_7e6085948b018986) }

var fileDescriptor_7e6085948b018986 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xca, 0xcb, 0x4c, 0xca,
	0x2c, 0x2a, 0xd5, 0x2f, 0x2e, 0x4d, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x4f, 0x2d, 0x4b, 0xcd, 0x2b,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x83, 0xc8, 0xe9, 0x81, 0xe4, 0xf4, 0xca, 0x0c,
	0xa5, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x52, 0xfa, 0x20, 0x16, 0x44, 0x95, 0x94, 0x4c, 0x7a,
	0x7e, 0x7e, 0x7a, 0x4e, 0xaa, 0x7e, 0x62, 0x41, 0xa6, 0x7e, 0x62, 0x5e, 0x5e, 0x7e, 0x49, 0x62,
	0x49, 0x66, 0x7e, 0x5e, 0x31, 0x54, 0x16, 0xdd, 0xfc, 0xe2, 0x92, 0xc4, 0x92, 0x54, 0x88, 0x9c,
	0x52, 0x2a, 0x97, 0x90, 0x2b, 0xc8, 0xba, 0xd0, 0x82, 0x94, 0xc4, 0x92, 0xd4, 0xe0, 0xd2, 0x94,
	0xfc, 0xd4, 0xa2, 0x62, 0x21, 0x73, 0x2e, 0xf6, 0x62, 0x08, 0x53, 0x82, 0x51, 0x81, 0x51, 0x83,
	0xdb, 0x48, 0x5c, 0x0f, 0xd5, 0x1d, 0x7a, 0x50, 0x95, 0x4e, 0x2c, 0x27, 0xee, 0xc9, 0x33, 0x04,
	0xc1, 0x54, 0x0b, 0x89, 0x71, 0xb1, 0x25, 0x26, 0x83, 0xec, 0x96, 0x60, 0x52, 0x60, 0xd4, 0xe0,
	0x0c, 0x82, 0xf2, 0x9c, 0x5c, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23,
	0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a,
	0x2b, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0xdf, 0x0f, 0x6c, 0x87, 0x73,
	0x46, 0x62, 0x66, 0x9e, 0x3e, 0xd4, 0xcd, 0x15, 0x10, 0x57, 0x97, 0x54, 0x16, 0xa4, 0x16, 0x27,
	0xb1, 0x81, 0xdd, 0x6c, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xcb, 0xb3, 0x68, 0xee, 0x31, 0x01,
	0x00, 0x00,
}

func (m *EventUpdateSudoers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventUpdateSudoers) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventUpdateSudoers) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Action) > 0 {
		i -= len(m.Action)
		copy(dAtA[i:], m.Action)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Action)))
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Sudoers.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintEvent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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
func (m *EventUpdateSudoers) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Sudoers.Size()
	n += 1 + l + sovEvent(uint64(l))
	l = len(m.Action)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventUpdateSudoers) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: EventUpdateSudoers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventUpdateSudoers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sudoers", wireType)
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
			if err := m.Sudoers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Action", wireType)
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
			m.Action = string(dAtA[iNdEx:postIndex])
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
