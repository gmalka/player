// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: playserService.proto

package proto

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

type LoadSongResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Song []byte `protobuf:"bytes,1,opt,name=song,proto3" json:"song,omitempty"`
}

func (x *LoadSongResponse) Reset() {
	*x = LoadSongResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_playserService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadSongResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadSongResponse) ProtoMessage() {}

func (x *LoadSongResponse) ProtoReflect() protoreflect.Message {
	mi := &file_playserService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadSongResponse.ProtoReflect.Descriptor instead.
func (*LoadSongResponse) Descriptor() ([]byte, []int) {
	return file_playserService_proto_rawDescGZIP(), []int{0}
}

func (x *LoadSongResponse) GetSong() []byte {
	if x != nil {
		return x.Song
	}
	return nil
}

type SongRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *SongRequest) Reset() {
	*x = SongRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_playserService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SongRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SongRequest) ProtoMessage() {}

func (x *SongRequest) ProtoReflect() protoreflect.Message {
	mi := &file_playserService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SongRequest.ProtoReflect.Descriptor instead.
func (*SongRequest) Descriptor() ([]byte, []int) {
	return file_playserService_proto_rawDescGZIP(), []int{1}
}

func (x *SongRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type None struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *None) Reset() {
	*x = None{}
	if protoimpl.UnsafeEnabled {
		mi := &file_playserService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *None) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*None) ProtoMessage() {}

func (x *None) ProtoReflect() protoreflect.Message {
	mi := &file_playserService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use None.ProtoReflect.Descriptor instead.
func (*None) Descriptor() ([]byte, []int) {
	return file_playserService_proto_rawDescGZIP(), []int{2}
}

var File_playserService_proto protoreflect.FileDescriptor

var file_playserService_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x6c, 0x61, 0x79, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x26, 0x0a,
	0x10, 0x4c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x73, 0x6f, 0x6e, 0x67, 0x22, 0x21, 0x0a, 0x0b, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x06, 0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65,
	0x32, 0x82, 0x01, 0x0a, 0x12, 0x4d, 0x75, 0x73, 0x69, 0x63, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x08, 0x4c, 0x6f, 0x61, 0x64, 0x53,
	0x6f, 0x6e, 0x67, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x6f, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x30, 0x01, 0x12, 0x2f, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x53, 0x6f, 0x6e, 0x67, 0x73,
	0x12, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x6e, 0x65, 0x1a, 0x12, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x00, 0x30, 0x01, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_playserService_proto_rawDescOnce sync.Once
	file_playserService_proto_rawDescData = file_playserService_proto_rawDesc
)

func file_playserService_proto_rawDescGZIP() []byte {
	file_playserService_proto_rawDescOnce.Do(func() {
		file_playserService_proto_rawDescData = protoimpl.X.CompressGZIP(file_playserService_proto_rawDescData)
	})
	return file_playserService_proto_rawDescData
}

var file_playserService_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_playserService_proto_goTypes = []interface{}{
	(*LoadSongResponse)(nil), // 0: proto.LoadSongResponse
	(*SongRequest)(nil),      // 1: proto.SongRequest
	(*None)(nil),             // 2: proto.None
}
var file_playserService_proto_depIdxs = []int32{
	1, // 0: proto.MusicPlayerService.LoadSong:input_type -> proto.SongRequest
	2, // 1: proto.MusicPlayerService.GetSongs:input_type -> proto.None
	0, // 2: proto.MusicPlayerService.LoadSong:output_type -> proto.LoadSongResponse
	1, // 3: proto.MusicPlayerService.GetSongs:output_type -> proto.SongRequest
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_playserService_proto_init() }
func file_playserService_proto_init() {
	if File_playserService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_playserService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadSongResponse); i {
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
		file_playserService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SongRequest); i {
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
		file_playserService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*None); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_playserService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_playserService_proto_goTypes,
		DependencyIndexes: file_playserService_proto_depIdxs,
		MessageInfos:      file_playserService_proto_msgTypes,
	}.Build()
	File_playserService_proto = out.File
	file_playserService_proto_rawDesc = nil
	file_playserService_proto_goTypes = nil
	file_playserService_proto_depIdxs = nil
}
