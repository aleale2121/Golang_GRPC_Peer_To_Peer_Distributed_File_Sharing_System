// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: files.proto

package files

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DownloadSongRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SongId string `protobuf:"bytes,1,opt,name=SongId,proto3" json:"SongId,omitempty"`
}

func (x *DownloadSongRequest) Reset() {
	*x = DownloadSongRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadSongRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadSongRequest) ProtoMessage() {}

func (x *DownloadSongRequest) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadSongRequest.ProtoReflect.Descriptor instead.
func (*DownloadSongRequest) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{0}
}

func (x *DownloadSongRequest) GetSongId() string {
	if x != nil {
		return x.SongId
	}
	return ""
}

type DownloadSongResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkData []byte `protobuf:"bytes,1,opt,name=chunk_data,json=chunkData,proto3" json:"chunk_data,omitempty"`
}

func (x *DownloadSongResponse) Reset() {
	*x = DownloadSongResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadSongResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadSongResponse) ProtoMessage() {}

func (x *DownloadSongResponse) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadSongResponse.ProtoReflect.Descriptor instead.
func (*DownloadSongResponse) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{1}
}

func (x *DownloadSongResponse) GetChunkData() []byte {
	if x != nil {
		return x.ChunkData
	}
	return nil
}

type UploadSongRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//	*UploadSongRequest_ChunkData
	//	*UploadSongRequest_Title
	Data isUploadSongRequest_Data `protobuf_oneof:"data"`
}

func (x *UploadSongRequest) Reset() {
	*x = UploadSongRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadSongRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadSongRequest) ProtoMessage() {}

func (x *UploadSongRequest) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadSongRequest.ProtoReflect.Descriptor instead.
func (*UploadSongRequest) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{2}
}

func (m *UploadSongRequest) GetData() isUploadSongRequest_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *UploadSongRequest) GetChunkData() []byte {
	if x, ok := x.GetData().(*UploadSongRequest_ChunkData); ok {
		return x.ChunkData
	}
	return nil
}

func (x *UploadSongRequest) GetTitle() string {
	if x, ok := x.GetData().(*UploadSongRequest_Title); ok {
		return x.Title
	}
	return ""
}

type isUploadSongRequest_Data interface {
	isUploadSongRequest_Data()
}

type UploadSongRequest_ChunkData struct {
	ChunkData []byte `protobuf:"bytes,1,opt,name=chunk_data,json=chunkData,proto3,oneof"`
}

type UploadSongRequest_Title struct {
	Title string `protobuf:"bytes,2,opt,name=title,proto3,oneof"`
}

func (*UploadSongRequest_ChunkData) isUploadSongRequest_Data() {}

func (*UploadSongRequest_Title) isUploadSongRequest_Data() {}

type UploadSongResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UploadSongResponse) Reset() {
	*x = UploadSongResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadSongResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadSongResponse) ProtoMessage() {}

func (x *UploadSongResponse) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadSongResponse.ProtoReflect.Descriptor instead.
func (*UploadSongResponse) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{3}
}

func (x *UploadSongResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetSongsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetSongsRequest) Reset() {
	*x = GetSongsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSongsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSongsRequest) ProtoMessage() {}

func (x *GetSongsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSongsRequest.ProtoReflect.Descriptor instead.
func (*GetSongsRequest) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{4}
}

type GetSongsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Songs []*SongData `protobuf:"bytes,1,rep,name=songs,proto3" json:"songs,omitempty"`
}

func (x *GetSongsResponse) Reset() {
	*x = GetSongsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSongsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSongsResponse) ProtoMessage() {}

func (x *GetSongsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSongsResponse.ProtoReflect.Descriptor instead.
func (*GetSongsResponse) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{5}
}

func (x *GetSongsResponse) GetSongs() []*SongData {
	if x != nil {
		return x.Songs
	}
	return nil
}

type ConnectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *SongData `protobuf:"bytes,1,opt,name=Info,proto3" json:"Info,omitempty"`
}

func (x *ConnectRequest) Reset() {
	*x = ConnectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectRequest) ProtoMessage() {}

func (x *ConnectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectRequest.ProtoReflect.Descriptor instead.
func (*ConnectRequest) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{6}
}

func (x *ConnectRequest) GetInfo() *SongData {
	if x != nil {
		return x.Info
	}
	return nil
}

type ConnectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConnectResponse) Reset() {
	*x = ConnectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectResponse) ProtoMessage() {}

func (x *ConnectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectResponse.ProtoReflect.Descriptor instead.
func (*ConnectResponse) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{7}
}

type SongData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Port  int32    `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Songs []string `protobuf:"bytes,3,rep,name=songs,proto3" json:"songs,omitempty"`
}

func (x *SongData) Reset() {
	*x = SongData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SongData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SongData) ProtoMessage() {}

func (x *SongData) ProtoReflect() protoreflect.Message {
	mi := &file_files_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SongData.ProtoReflect.Descriptor instead.
func (*SongData) Descriptor() ([]byte, []int) {
	return file_files_proto_rawDescGZIP(), []int{8}
}

func (x *SongData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SongData) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *SongData) GetSongs() []string {
	if x != nil {
		return x.Songs
	}
	return nil
}

var File_files_proto protoreflect.FileDescriptor

var file_files_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x22, 0x2d, 0x0a, 0x13, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x53,
	0x6f, 0x6e, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x6f, 0x6e,
	0x67, 0x49, 0x64, 0x22, 0x35, 0x0a, 0x14, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x53,
	0x6f, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x22, 0x54, 0x0a, 0x11, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1f, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x16, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x24, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x11, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x6f, 0x6e,
	0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x39, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x53, 0x6f, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a,
	0x05, 0x73, 0x6f, 0x6e, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x2e, 0x53, 0x6f, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x73,
	0x6f, 0x6e, 0x67, 0x73, 0x22, 0x35, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x53, 0x6f, 0x6e,
	0x67, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x11, 0x0a, 0x0f, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x44,
	0x0a, 0x08, 0x53, 0x6f, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x6f, 0x6e, 0x67, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x73,
	0x6f, 0x6e, 0x67, 0x73, 0x32, 0x99, 0x02, 0x0a, 0x0c, 0x53, 0x6f, 0x6e, 0x67, 0x73, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x12, 0x15, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x43, 0x0a, 0x0a, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67, 0x12, 0x18, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x28, 0x01, 0x12, 0x3f, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x53, 0x6f, 0x6e, 0x67, 0x73,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x6f, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x6f, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x0c, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x53, 0x6f, 0x6e, 0x67, 0x12, 0x1a, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x44, 0x6f,
	0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1b, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f,
	0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01,
	0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_files_proto_rawDescOnce sync.Once
	file_files_proto_rawDescData = file_files_proto_rawDesc
)

func file_files_proto_rawDescGZIP() []byte {
	file_files_proto_rawDescOnce.Do(func() {
		file_files_proto_rawDescData = protoimpl.X.CompressGZIP(file_files_proto_rawDescData)
	})
	return file_files_proto_rawDescData
}

var file_files_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_files_proto_goTypes = []interface{}{
	(*DownloadSongRequest)(nil),  // 0: files.DownloadSongRequest
	(*DownloadSongResponse)(nil), // 1: files.DownloadSongResponse
	(*UploadSongRequest)(nil),    // 2: files.UploadSongRequest
	(*UploadSongResponse)(nil),   // 3: files.UploadSongResponse
	(*GetSongsRequest)(nil),      // 4: files.GetSongsRequest
	(*GetSongsResponse)(nil),     // 5: files.GetSongsResponse
	(*ConnectRequest)(nil),       // 6: files.ConnectRequest
	(*ConnectResponse)(nil),      // 7: files.ConnectResponse
	(*SongData)(nil),             // 8: files.SongData
}
var file_files_proto_depIdxs = []int32{
	8, // 0: files.GetSongsResponse.songs:type_name -> files.SongData
	8, // 1: files.ConnectRequest.Info:type_name -> files.SongData
	6, // 2: files.SongsService.Connect:input_type -> files.ConnectRequest
	2, // 3: files.SongsService.UploadSong:input_type -> files.UploadSongRequest
	4, // 4: files.SongsService.GetSongsList:input_type -> files.GetSongsRequest
	0, // 5: files.SongsService.DownloadSong:input_type -> files.DownloadSongRequest
	7, // 6: files.SongsService.Connect:output_type -> files.ConnectResponse
	3, // 7: files.SongsService.UploadSong:output_type -> files.UploadSongResponse
	5, // 8: files.SongsService.GetSongsList:output_type -> files.GetSongsResponse
	1, // 9: files.SongsService.DownloadSong:output_type -> files.DownloadSongResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_files_proto_init() }
func file_files_proto_init() {
	if File_files_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_files_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadSongRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_files_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadSongResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_files_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadSongRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_files_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadSongResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_files_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSongsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_files_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSongsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_files_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_files_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_files_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SongData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_files_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*UploadSongRequest_ChunkData)(nil),
		(*UploadSongRequest_Title)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_files_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_files_proto_goTypes,
		DependencyIndexes: file_files_proto_depIdxs,
		MessageInfos:      file_files_proto_msgTypes,
	}.Build()
	File_files_proto = out.File
	file_files_proto_rawDesc = nil
	file_files_proto_goTypes = nil
	file_files_proto_depIdxs = nil
}