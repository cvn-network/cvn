// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ethermint/feemarket/v1/feemarket.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// Params defines the EVM module parameters
type Params struct {
	// no_base_fee forces the EIP-1559 base fee to 0 (needed for 0 price calls)
	NoBaseFee bool `protobuf:"varint,1,opt,name=no_base_fee,json=noBaseFee,proto3" json:"no_base_fee,omitempty"`
	// base_fee_change_denominator bounds the amount the base fee can change
	// between blocks.
	BaseFeeChangeDenominator uint32 `protobuf:"varint,2,opt,name=base_fee_change_denominator,json=baseFeeChangeDenominator,proto3" json:"base_fee_change_denominator,omitempty"`
	// elasticity_multiplier bounds the maximum gas limit an EIP-1559 block may
	// have.
	ElasticityMultiplier uint32 `protobuf:"varint,3,opt,name=elasticity_multiplier,json=elasticityMultiplier,proto3" json:"elasticity_multiplier,omitempty"`
	// enable_height defines at which block height the base fee calculation is
	// enabled.
	EnableHeight int64 `protobuf:"varint,5,opt,name=enable_height,json=enableHeight,proto3" json:"enable_height,omitempty"`
	// base_fee for EIP-1559 blocks.
	BaseFee github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=base_fee,json=baseFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"base_fee"`
	// min_gas_price defines the minimum gas price value for cosmos and eth
	// transactions
	MinGasPrice github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=min_gas_price,json=minGasPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"min_gas_price"`
	// min_gas_multiplier bounds the minimum gas used to be charged
	// to senders based on gas limit
	MinGasMultiplier github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=min_gas_multiplier,json=minGasMultiplier,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"min_gas_multiplier"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_4feb8b20cf98e6e1, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetNoBaseFee() bool {
	if m != nil {
		return m.NoBaseFee
	}
	return false
}

func (m *Params) GetBaseFeeChangeDenominator() uint32 {
	if m != nil {
		return m.BaseFeeChangeDenominator
	}
	return 0
}

func (m *Params) GetElasticityMultiplier() uint32 {
	if m != nil {
		return m.ElasticityMultiplier
	}
	return 0
}

func (m *Params) GetEnableHeight() int64 {
	if m != nil {
		return m.EnableHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "ethermint.feemarket.v1.Params")
}

func init() {
	proto.RegisterFile("ethermint/feemarket/v1/feemarket.proto", fileDescriptor_4feb8b20cf98e6e1)
}

var fileDescriptor_4feb8b20cf98e6e1 = []byte{
	// 396 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x41, 0x6b, 0xdb, 0x30,
	0x14, 0xc7, 0xad, 0xa5, 0x4d, 0x5d, 0x75, 0x81, 0x20, 0xba, 0x61, 0x36, 0x70, 0xcd, 0x06, 0xc5,
	0x97, 0x5a, 0x74, 0x3d, 0xef, 0x92, 0x95, 0x6d, 0x19, 0x0c, 0x8a, 0x8f, 0x63, 0x20, 0x64, 0xf7,
	0xd5, 0x16, 0xb1, 0x24, 0x23, 0xa9, 0xde, 0xfa, 0x2d, 0xf6, 0xb1, 0x7a, 0xec, 0x69, 0x8c, 0x1d,
	0xca, 0x48, 0xbe, 0xc8, 0x88, 0x9d, 0xd8, 0xb9, 0x6e, 0x27, 0x3d, 0xbd, 0xff, 0x5f, 0xbf, 0xf7,
	0x24, 0x3d, 0x7c, 0x0a, 0xae, 0x04, 0x23, 0x85, 0x72, 0xf4, 0x06, 0x40, 0x72, 0xb3, 0x00, 0x47,
	0x9b, 0xf3, 0x61, 0x93, 0xd4, 0x46, 0x3b, 0x4d, 0x9e, 0xf7, 0xbe, 0x64, 0x90, 0x9a, 0xf3, 0x17,
	0xc7, 0x85, 0x2e, 0x74, 0x6b, 0xa1, 0xeb, 0xa8, 0x73, 0xbf, 0xfa, 0x39, 0xc2, 0xe3, 0x2b, 0x6e,
	0xb8, 0xb4, 0x24, 0xc4, 0x47, 0x4a, 0xb3, 0x8c, 0x5b, 0x60, 0x37, 0x00, 0x01, 0x8a, 0x50, 0xec,
	0xa7, 0x87, 0x4a, 0xcf, 0xb8, 0x85, 0xf7, 0x00, 0xe4, 0x2d, 0x7e, 0xb9, 0x15, 0x59, 0x5e, 0x72,
	0x55, 0x00, 0xbb, 0x06, 0xa5, 0xa5, 0x50, 0xdc, 0x69, 0x13, 0x3c, 0x89, 0x50, 0x3c, 0x49, 0x83,
	0xac, 0x73, 0xbf, 0x6b, 0x0d, 0x97, 0x83, 0x4e, 0x2e, 0xf0, 0x33, 0xa8, 0xb8, 0x75, 0x22, 0x17,
	0xee, 0x8e, 0xc9, 0xdb, 0xca, 0x89, 0xba, 0x12, 0x60, 0x82, 0x51, 0x7b, 0xf0, 0x78, 0x10, 0x3f,
	0xf7, 0x1a, 0x79, 0x8d, 0x27, 0xa0, 0x78, 0x56, 0x01, 0x2b, 0x41, 0x14, 0xa5, 0x0b, 0xf6, 0x23,
	0x14, 0x8f, 0xd2, 0xa7, 0x5d, 0xf2, 0x63, 0x9b, 0x23, 0x73, 0xec, 0xf7, 0x5d, 0x8f, 0x23, 0x14,
	0x1f, 0xce, 0x92, 0xfb, 0xc7, 0x13, 0xef, 0xf7, 0xe3, 0xc9, 0x69, 0x21, 0x5c, 0x79, 0x9b, 0x25,
	0xb9, 0x96, 0x34, 0xd7, 0x56, 0x6a, 0xbb, 0x59, 0xce, 0xec, 0xf5, 0x82, 0xba, 0xbb, 0x1a, 0x6c,
	0x32, 0x57, 0x2e, 0x3d, 0xd8, 0x74, 0x4d, 0x52, 0x3c, 0x91, 0x42, 0xb1, 0x82, 0x5b, 0x56, 0x1b,
	0x91, 0x43, 0x70, 0xf0, 0xcf, 0xbc, 0x4b, 0xc8, 0xd3, 0x23, 0x29, 0xd4, 0x07, 0x6e, 0xaf, 0xd6,
	0x08, 0xf2, 0x15, 0x93, 0x2d, 0x73, 0xe7, 0xd6, 0xfe, 0x7f, 0x81, 0xa7, 0x1d, 0x78, 0x78, 0xa1,
	0x4f, 0x7b, 0xfe, 0xde, 0x74, 0x3f, 0x9d, 0x0a, 0x25, 0x9c, 0xe0, 0x55, 0xff, 0x7d, 0xb3, 0xf9,
	0xfd, 0x32, 0x44, 0x0f, 0xcb, 0x10, 0xfd, 0x59, 0x86, 0xe8, 0xc7, 0x2a, 0xf4, 0x1e, 0x56, 0xa1,
	0xf7, 0x6b, 0x15, 0x7a, 0x5f, 0xe8, 0x6e, 0xad, 0x46, 0x9d, 0x29, 0x70, 0xdf, 0xb4, 0x59, 0xac,
	0x63, 0xda, 0xbc, 0xa1, 0xdf, 0x77, 0xc6, 0xab, 0x2d, 0x9c, 0x8d, 0xdb, 0x51, 0xb9, 0xf8, 0x1b,
	0x00, 0x00, 0xff, 0xff, 0x02, 0x43, 0xb1, 0xfe, 0x82, 0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.MinGasMultiplier.Size()
		i -= size
		if _, err := m.MinGasMultiplier.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFeemarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.MinGasPrice.Size()
		i -= size
		if _, err := m.MinGasPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFeemarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.BaseFee.Size()
		i -= size
		if _, err := m.BaseFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFeemarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if m.EnableHeight != 0 {
		i = encodeVarintFeemarket(dAtA, i, uint64(m.EnableHeight))
		i--
		dAtA[i] = 0x28
	}
	if m.ElasticityMultiplier != 0 {
		i = encodeVarintFeemarket(dAtA, i, uint64(m.ElasticityMultiplier))
		i--
		dAtA[i] = 0x18
	}
	if m.BaseFeeChangeDenominator != 0 {
		i = encodeVarintFeemarket(dAtA, i, uint64(m.BaseFeeChangeDenominator))
		i--
		dAtA[i] = 0x10
	}
	if m.NoBaseFee {
		i--
		if m.NoBaseFee {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintFeemarket(dAtA []byte, offset int, v uint64) int {
	offset -= sovFeemarket(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.NoBaseFee {
		n += 2
	}
	if m.BaseFeeChangeDenominator != 0 {
		n += 1 + sovFeemarket(uint64(m.BaseFeeChangeDenominator))
	}
	if m.ElasticityMultiplier != 0 {
		n += 1 + sovFeemarket(uint64(m.ElasticityMultiplier))
	}
	if m.EnableHeight != 0 {
		n += 1 + sovFeemarket(uint64(m.EnableHeight))
	}
	l = m.BaseFee.Size()
	n += 1 + l + sovFeemarket(uint64(l))
	l = m.MinGasPrice.Size()
	n += 1 + l + sovFeemarket(uint64(l))
	l = m.MinGasMultiplier.Size()
	n += 1 + l + sovFeemarket(uint64(l))
	return n
}

func sovFeemarket(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFeemarket(x uint64) (n int) {
	return sovFeemarket(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeemarket
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NoBaseFee", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.NoBaseFee = bool(v != 0)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFeeChangeDenominator", wireType)
			}
			m.BaseFeeChangeDenominator = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BaseFeeChangeDenominator |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ElasticityMultiplier", wireType)
			}
			m.ElasticityMultiplier = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ElasticityMultiplier |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableHeight", wireType)
			}
			m.EnableHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EnableHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
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
				return ErrInvalidLengthFeemarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeemarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BaseFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinGasPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
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
				return ErrInvalidLengthFeemarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeemarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinGasPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinGasMultiplier", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
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
				return ErrInvalidLengthFeemarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeemarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinGasMultiplier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFeemarket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeemarket
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
func skipFeemarket(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFeemarket
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
					return 0, ErrIntOverflowFeemarket
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
					return 0, ErrIntOverflowFeemarket
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
				return 0, ErrInvalidLengthFeemarket
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFeemarket
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFeemarket
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFeemarket        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFeemarket          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFeemarket = fmt.Errorf("proto: unexpected end of group")
)
