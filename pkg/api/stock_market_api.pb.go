// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/stock_market_api.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/empty"
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

type Period int32

const (
	Period_HOUR Period = 0
	Period_DAY  Period = 1
	Period_WEEK Period = 2
	Period_ALL  Period = 3
)

var Period_name = map[int32]string{
	0: "HOUR",
	1: "DAY",
	2: "WEEK",
	3: "ALL",
}

var Period_value = map[string]int32{
	"HOUR": 0,
	"DAY":  1,
	"WEEK": 2,
	"ALL":  3,
}

func (x Period) String() string {
	return proto.EnumName(Period_name, int32(x))
}

func (Period) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{0}
}

type UserId struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserId) Reset()         { *m = UserId{} }
func (m *UserId) String() string { return proto.CompactTextString(m) }
func (*UserId) ProtoMessage()    {}
func (*UserId) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{0}
}

func (m *UserId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserId.Unmarshal(m, b)
}
func (m *UserId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserId.Marshal(b, m, deterministic)
}
func (m *UserId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserId.Merge(m, src)
}
func (m *UserId) XXX_Size() int {
	return xxx_messageInfo_UserId.Size(m)
}
func (m *UserId) XXX_DiscardUnknown() {
	xxx_messageInfo_UserId.DiscardUnknown(m)
}

var xxx_messageInfo_UserId proto.InternalMessageInfo

func (m *UserId) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type StockRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Amount               int32    `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	UserId               *UserId  `protobuf:"bytes,3,opt,name=userId,proto3" json:"userId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StockRequest) Reset()         { *m = StockRequest{} }
func (m *StockRequest) String() string { return proto.CompactTextString(m) }
func (*StockRequest) ProtoMessage()    {}
func (*StockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{1}
}

func (m *StockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StockRequest.Unmarshal(m, b)
}
func (m *StockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StockRequest.Marshal(b, m, deterministic)
}
func (m *StockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StockRequest.Merge(m, src)
}
func (m *StockRequest) XXX_Size() int {
	return xxx_messageInfo_StockRequest.Size(m)
}
func (m *StockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StockRequest proto.InternalMessageInfo

func (m *StockRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StockRequest) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *StockRequest) GetUserId() *UserId {
	if m != nil {
		return m.UserId
	}
	return nil
}

type Stock struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Amount               int32    `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stock) Reset()         { *m = Stock{} }
func (m *Stock) String() string { return proto.CompactTextString(m) }
func (*Stock) ProtoMessage()    {}
func (*Stock) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{2}
}

func (m *Stock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stock.Unmarshal(m, b)
}
func (m *Stock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stock.Marshal(b, m, deterministic)
}
func (m *Stock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stock.Merge(m, src)
}
func (m *Stock) XXX_Size() int {
	return xxx_messageInfo_Stock.Size(m)
}
func (m *Stock) XXX_DiscardUnknown() {
	xxx_messageInfo_Stock.DiscardUnknown(m)
}

var xxx_messageInfo_Stock proto.InternalMessageInfo

func (m *Stock) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Stock) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type StockChanges struct {
	Stock                *Stock   `protobuf:"bytes,1,opt,name=stock,proto3" json:"stock,omitempty"`
	OldPrice             float64  `protobuf:"fixed64,2,opt,name=oldPrice,proto3" json:"oldPrice,omitempty"`
	CurrentPrice         float64  `protobuf:"fixed64,3,opt,name=currentPrice,proto3" json:"currentPrice,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StockChanges) Reset()         { *m = StockChanges{} }
func (m *StockChanges) String() string { return proto.CompactTextString(m) }
func (*StockChanges) ProtoMessage()    {}
func (*StockChanges) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{3}
}

func (m *StockChanges) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StockChanges.Unmarshal(m, b)
}
func (m *StockChanges) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StockChanges.Marshal(b, m, deterministic)
}
func (m *StockChanges) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StockChanges.Merge(m, src)
}
func (m *StockChanges) XXX_Size() int {
	return xxx_messageInfo_StockChanges.Size(m)
}
func (m *StockChanges) XXX_DiscardUnknown() {
	xxx_messageInfo_StockChanges.DiscardUnknown(m)
}

var xxx_messageInfo_StockChanges proto.InternalMessageInfo

func (m *StockChanges) GetStock() *Stock {
	if m != nil {
		return m.Stock
	}
	return nil
}

func (m *StockChanges) GetOldPrice() float64 {
	if m != nil {
		return m.OldPrice
	}
	return 0
}

func (m *StockChanges) GetCurrentPrice() float64 {
	if m != nil {
		return m.CurrentPrice
	}
	return 0
}

type GetStocksResponse struct {
	Stocks               []*Stock `protobuf:"bytes,1,rep,name=stocks,proto3" json:"stocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetStocksResponse) Reset()         { *m = GetStocksResponse{} }
func (m *GetStocksResponse) String() string { return proto.CompactTextString(m) }
func (*GetStocksResponse) ProtoMessage()    {}
func (*GetStocksResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{4}
}

func (m *GetStocksResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetStocksResponse.Unmarshal(m, b)
}
func (m *GetStocksResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetStocksResponse.Marshal(b, m, deterministic)
}
func (m *GetStocksResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetStocksResponse.Merge(m, src)
}
func (m *GetStocksResponse) XXX_Size() int {
	return xxx_messageInfo_GetStocksResponse.Size(m)
}
func (m *GetStocksResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetStocksResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetStocksResponse proto.InternalMessageInfo

func (m *GetStocksResponse) GetStocks() []*Stock {
	if m != nil {
		return m.Stocks
	}
	return nil
}

type GetPortfolioChangesRequest struct {
	UserId               *UserId  `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Period               Period   `protobuf:"varint,2,opt,name=period,proto3,enum=api.Period" json:"period,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPortfolioChangesRequest) Reset()         { *m = GetPortfolioChangesRequest{} }
func (m *GetPortfolioChangesRequest) String() string { return proto.CompactTextString(m) }
func (*GetPortfolioChangesRequest) ProtoMessage()    {}
func (*GetPortfolioChangesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{5}
}

func (m *GetPortfolioChangesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPortfolioChangesRequest.Unmarshal(m, b)
}
func (m *GetPortfolioChangesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPortfolioChangesRequest.Marshal(b, m, deterministic)
}
func (m *GetPortfolioChangesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPortfolioChangesRequest.Merge(m, src)
}
func (m *GetPortfolioChangesRequest) XXX_Size() int {
	return xxx_messageInfo_GetPortfolioChangesRequest.Size(m)
}
func (m *GetPortfolioChangesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPortfolioChangesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPortfolioChangesRequest proto.InternalMessageInfo

func (m *GetPortfolioChangesRequest) GetUserId() *UserId {
	if m != nil {
		return m.UserId
	}
	return nil
}

func (m *GetPortfolioChangesRequest) GetPeriod() Period {
	if m != nil {
		return m.Period
	}
	return Period_HOUR
}

type GetPortfolioChangesResponse struct {
	Stocks               []*StockChanges `protobuf:"bytes,1,rep,name=stocks,proto3" json:"stocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *GetPortfolioChangesResponse) Reset()         { *m = GetPortfolioChangesResponse{} }
func (m *GetPortfolioChangesResponse) String() string { return proto.CompactTextString(m) }
func (*GetPortfolioChangesResponse) ProtoMessage()    {}
func (*GetPortfolioChangesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{6}
}

func (m *GetPortfolioChangesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPortfolioChangesResponse.Unmarshal(m, b)
}
func (m *GetPortfolioChangesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPortfolioChangesResponse.Marshal(b, m, deterministic)
}
func (m *GetPortfolioChangesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPortfolioChangesResponse.Merge(m, src)
}
func (m *GetPortfolioChangesResponse) XXX_Size() int {
	return xxx_messageInfo_GetPortfolioChangesResponse.Size(m)
}
func (m *GetPortfolioChangesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPortfolioChangesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPortfolioChangesResponse proto.InternalMessageInfo

func (m *GetPortfolioChangesResponse) GetStocks() []*StockChanges {
	if m != nil {
		return m.Stocks
	}
	return nil
}

type AddNotificationRequest struct {
	UserId               *UserId  `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	StockName            string   `protobuf:"bytes,2,opt,name=stockName,proto3" json:"stockName,omitempty"`
	Threshold            float64  `protobuf:"fixed64,3,opt,name=threshold,proto3" json:"threshold,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddNotificationRequest) Reset()         { *m = AddNotificationRequest{} }
func (m *AddNotificationRequest) String() string { return proto.CompactTextString(m) }
func (*AddNotificationRequest) ProtoMessage()    {}
func (*AddNotificationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ae39547e7469637, []int{7}
}

func (m *AddNotificationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddNotificationRequest.Unmarshal(m, b)
}
func (m *AddNotificationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddNotificationRequest.Marshal(b, m, deterministic)
}
func (m *AddNotificationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddNotificationRequest.Merge(m, src)
}
func (m *AddNotificationRequest) XXX_Size() int {
	return xxx_messageInfo_AddNotificationRequest.Size(m)
}
func (m *AddNotificationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddNotificationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddNotificationRequest proto.InternalMessageInfo

func (m *AddNotificationRequest) GetUserId() *UserId {
	if m != nil {
		return m.UserId
	}
	return nil
}

func (m *AddNotificationRequest) GetStockName() string {
	if m != nil {
		return m.StockName
	}
	return ""
}

func (m *AddNotificationRequest) GetThreshold() float64 {
	if m != nil {
		return m.Threshold
	}
	return 0
}

func init() {
	proto.RegisterEnum("api.Period", Period_name, Period_value)
	proto.RegisterType((*UserId)(nil), "api.UserId")
	proto.RegisterType((*StockRequest)(nil), "api.StockRequest")
	proto.RegisterType((*Stock)(nil), "api.Stock")
	proto.RegisterType((*StockChanges)(nil), "api.StockChanges")
	proto.RegisterType((*GetStocksResponse)(nil), "api.GetStocksResponse")
	proto.RegisterType((*GetPortfolioChangesRequest)(nil), "api.GetPortfolioChangesRequest")
	proto.RegisterType((*GetPortfolioChangesResponse)(nil), "api.GetPortfolioChangesResponse")
	proto.RegisterType((*AddNotificationRequest)(nil), "api.AddNotificationRequest")
}

func init() { proto.RegisterFile("api/stock_market_api.proto", fileDescriptor_0ae39547e7469637) }

var fileDescriptor_0ae39547e7469637 = []byte{
	// 561 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x4f, 0x6f, 0xda, 0x4e,
	0x10, 0xc5, 0x38, 0xf1, 0x2f, 0x0c, 0x51, 0x7e, 0x64, 0x2b, 0x21, 0x64, 0x2a, 0x15, 0xb9, 0x97,
	0xb4, 0x52, 0x6c, 0x09, 0x0e, 0x91, 0x7a, 0xa3, 0x29, 0x0a, 0x55, 0x53, 0x8a, 0x1c, 0x45, 0x55,
	0x73, 0x41, 0xc6, 0x1e, 0x60, 0x85, 0xed, 0x75, 0xd7, 0x0b, 0x15, 0xfd, 0x24, 0xfd, 0xb8, 0x95,
	0xc7, 0xe6, 0x4f, 0x29, 0x1c, 0xd2, 0xdb, 0xee, 0xbc, 0x99, 0x37, 0x6f, 0x66, 0x1e, 0x98, 0x5e,
	0xc2, 0x9d, 0x54, 0x09, 0x7f, 0x3e, 0x8a, 0x3c, 0x39, 0x47, 0x35, 0xf2, 0x12, 0x6e, 0x27, 0x52,
	0x28, 0xc1, 0x74, 0x2f, 0xe1, 0x66, 0x73, 0x2a, 0xc4, 0x34, 0x44, 0x87, 0x42, 0xe3, 0xc5, 0xc4,
	0xc1, 0x28, 0x51, 0xab, 0x3c, 0xc3, 0x6a, 0x80, 0xf1, 0x98, 0xa2, 0xfc, 0x18, 0xb0, 0x0b, 0x28,
	0xf3, 0xa0, 0xa1, 0xb5, 0xb4, 0xab, 0x53, 0xb7, 0xcc, 0x03, 0x6b, 0x04, 0xe7, 0x0f, 0x19, 0xab,
	0x8b, 0xdf, 0x17, 0x98, 0x2a, 0xc6, 0xe0, 0x24, 0xf6, 0x22, 0xa4, 0x8c, 0x8a, 0x4b, 0x6f, 0x56,
	0x07, 0xc3, 0x8b, 0xc4, 0x22, 0x56, 0x8d, 0x32, 0xd5, 0x15, 0x3f, 0xf6, 0x1a, 0x8c, 0x05, 0xb1,
	0x36, 0xf4, 0x96, 0x76, 0x55, 0x6d, 0x57, 0xed, 0x4c, 0x53, 0xde, 0xc8, 0x2d, 0x20, 0xab, 0x03,
	0xa7, 0xd4, 0xe0, 0x39, 0xcc, 0x56, 0x52, 0xa8, 0xba, 0x9d, 0x79, 0xf1, 0x14, 0x53, 0xd6, 0x82,
	0x53, 0x9a, 0x9d, 0x8a, 0xab, 0x6d, 0xa0, 0x46, 0xb9, 0xee, 0x1c, 0x60, 0x26, 0x9c, 0x89, 0x30,
	0x18, 0x4a, 0xee, 0x23, 0x71, 0x69, 0xee, 0xe6, 0xcf, 0x2c, 0x38, 0xf7, 0x17, 0x52, 0x62, 0xac,
	0x72, 0x5c, 0x27, 0xfc, 0x8f, 0x98, 0x75, 0x03, 0x97, 0x77, 0xa8, 0x88, 0x32, 0x75, 0x31, 0x4d,
	0x44, 0x9c, 0x66, 0x85, 0x06, 0xb1, 0xa7, 0x0d, 0xad, 0xa5, 0xef, 0xf5, 0x2d, 0x10, 0x6b, 0x02,
	0xe6, 0x1d, 0xaa, 0xa1, 0x90, 0x6a, 0x22, 0x42, 0x2e, 0x0a, 0xc5, 0xeb, 0x75, 0x6e, 0x57, 0xa4,
	0x1d, 0x5d, 0x51, 0x96, 0x94, 0xa0, 0xe4, 0x22, 0x20, 0xe5, 0x17, 0x45, 0xd2, 0x90, 0x42, 0x6e,
	0x01, 0x59, 0x7d, 0x68, 0x1e, 0xec, 0x53, 0x48, 0x7d, 0xb3, 0x27, 0xf5, 0x72, 0x2b, 0x75, 0x9d,
	0xba, 0x56, 0xbc, 0x82, 0x7a, 0x37, 0x08, 0x06, 0x42, 0xf1, 0x09, 0xf7, 0x3d, 0xc5, 0x45, 0xfc,
	0x2c, 0xb5, 0x2f, 0xa1, 0x42, 0x44, 0x83, 0xec, 0x98, 0x65, 0x3a, 0xe6, 0x36, 0x90, 0xa1, 0x6a,
	0x26, 0x31, 0x9d, 0x89, 0x30, 0x28, 0x16, 0xbd, 0x0d, 0xbc, 0xb5, 0xc1, 0xc8, 0xc7, 0x62, 0x67,
	0x70, 0xd2, 0xff, 0xf2, 0xe8, 0xd6, 0x4a, 0xec, 0x3f, 0xd0, 0x3f, 0x74, 0xbf, 0xd5, 0xb4, 0x2c,
	0xf4, 0xb5, 0xd7, 0xfb, 0x54, 0x2b, 0x67, 0xa1, 0xee, 0xfd, 0x7d, 0x4d, 0x6f, 0xff, 0xd2, 0x81,
	0xd1, 0x0c, 0x9f, 0xc9, 0xf3, 0x0f, 0x28, 0x97, 0xd9, 0x41, 0x3b, 0x00, 0xb7, 0x12, 0x3d, 0x85,
	0x99, 0x34, 0x56, 0xb7, 0x73, 0xeb, 0xdb, 0x6b, 0xeb, 0xdb, 0xbd, 0xcc, 0xfa, 0xe6, 0xae, 0x7a,
	0xab, 0xc4, 0xda, 0x50, 0xd9, 0x5c, 0x98, 0xed, 0x62, 0x66, 0x9d, 0x3e, 0x7f, 0x9d, 0xdf, 0x2a,
	0xb1, 0x1b, 0x38, 0xeb, 0x06, 0x41, 0xee, 0xdf, 0x9d, 0x8d, 0x16, 0xfb, 0x32, 0x8f, 0x74, 0xb6,
	0x4a, 0xec, 0x1d, 0x54, 0x5d, 0x8c, 0xc4, 0x12, 0xff, 0xa1, 0xf6, 0x09, 0x5e, 0x1c, 0xb8, 0x34,
	0x7b, 0xb5, 0x56, 0x79, 0xc4, 0x6b, 0x66, 0xeb, 0x78, 0xc2, 0x66, 0xa0, 0x3e, 0xfc, 0xbf, 0x77,
	0x7b, 0xd6, 0xa4, 0xb2, 0xc3, 0x8e, 0x38, 0xae, 0xf2, 0xbd, 0xf3, 0x74, 0x3d, 0xe5, 0x2a, 0xf4,
	0xc6, 0xb6, 0xf8, 0x29, 0x62, 0x3b, 0xc0, 0xa5, 0xe3, 0xcf, 0x78, 0x18, 0xae, 0x06, 0xdc, 0x9f,
	0x3b, 0x33, 0x11, 0xe1, 0x0f, 0x21, 0xe7, 0xd7, 0x6d, 0x27, 0x99, 0x4f, 0x1d, 0x2f, 0xe1, 0x63,
	0x83, 0x28, 0x3a, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x77, 0x05, 0x29, 0xcc, 0xca, 0x04, 0x00,
	0x00,
}
