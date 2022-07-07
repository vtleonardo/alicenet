// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: proto/cobjs.proto

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

// Protobuf message implementation for struct Proposal
type Proposal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PClaims   *PClaims `protobuf:"bytes,1,opt,name=PClaims,proto3" json:"PClaims,omitempty"`
	Signature string   `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
	TxHshLst  []string `protobuf:"bytes,3,rep,name=TxHshLst,proto3" json:"TxHshLst,omitempty"`
}

func (x *Proposal) Reset() {
	*x = Proposal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Proposal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Proposal) ProtoMessage() {}

func (x *Proposal) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Proposal.ProtoReflect.Descriptor instead.
func (*Proposal) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{0}
}

func (x *Proposal) GetPClaims() *PClaims {
	if x != nil {
		return x.PClaims
	}
	return nil
}

func (x *Proposal) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *Proposal) GetTxHshLst() []string {
	if x != nil {
		return x.TxHshLst
	}
	return nil
}

// Protobuf message implementation for struct PreVoteNil
type PreVoteNil struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RCert     *RCert `protobuf:"bytes,1,opt,name=RCert,proto3" json:"RCert,omitempty"`
	Signature string `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (x *PreVoteNil) Reset() {
	*x = PreVoteNil{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreVoteNil) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreVoteNil) ProtoMessage() {}

func (x *PreVoteNil) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreVoteNil.ProtoReflect.Descriptor instead.
func (*PreVoteNil) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{1}
}

func (x *PreVoteNil) GetRCert() *RCert {
	if x != nil {
		return x.RCert
	}
	return nil
}

func (x *PreVoteNil) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

// Protobuf message implementation for struct PreCommitNil
type PreCommitNil struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RCert     *RCert `protobuf:"bytes,1,opt,name=RCert,proto3" json:"RCert,omitempty"`
	Signature string `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (x *PreCommitNil) Reset() {
	*x = PreCommitNil{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreCommitNil) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreCommitNil) ProtoMessage() {}

func (x *PreCommitNil) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreCommitNil.ProtoReflect.Descriptor instead.
func (*PreCommitNil) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{2}
}

func (x *PreCommitNil) GetRCert() *RCert {
	if x != nil {
		return x.RCert
	}
	return nil
}

func (x *PreCommitNil) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

// Protobuf message implementation for struct RCert
type RCert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RClaims  *RClaims `protobuf:"bytes,1,opt,name=RClaims,proto3" json:"RClaims,omitempty"`
	SigGroup string   `protobuf:"bytes,2,opt,name=SigGroup,proto3" json:"SigGroup,omitempty"`
}

func (x *RCert) Reset() {
	*x = RCert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RCert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RCert) ProtoMessage() {}

func (x *RCert) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RCert.ProtoReflect.Descriptor instead.
func (*RCert) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{3}
}

func (x *RCert) GetRClaims() *RClaims {
	if x != nil {
		return x.RClaims
	}
	return nil
}

func (x *RCert) GetSigGroup() string {
	if x != nil {
		return x.SigGroup
	}
	return ""
}

// Protobuf message implementation for struct NRClaims
type NRClaims struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RCert      *RCert   `protobuf:"bytes,1,opt,name=RCert,proto3" json:"RCert,omitempty"`
	RClaims    *RClaims `protobuf:"bytes,2,opt,name=RClaims,proto3" json:"RClaims,omitempty"`
	SigShare   string   `protobuf:"bytes,3,opt,name=SigShare,proto3" json:"SigShare,omitempty"`
	GroupShare string   `protobuf:"bytes,4,opt,name=GroupShare,proto3" json:"GroupShare,omitempty"`
}

func (x *NRClaims) Reset() {
	*x = NRClaims{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NRClaims) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NRClaims) ProtoMessage() {}

func (x *NRClaims) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NRClaims.ProtoReflect.Descriptor instead.
func (*NRClaims) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{4}
}

func (x *NRClaims) GetRCert() *RCert {
	if x != nil {
		return x.RCert
	}
	return nil
}

func (x *NRClaims) GetRClaims() *RClaims {
	if x != nil {
		return x.RClaims
	}
	return nil
}

func (x *NRClaims) GetSigShare() string {
	if x != nil {
		return x.SigShare
	}
	return ""
}

func (x *NRClaims) GetGroupShare() string {
	if x != nil {
		return x.GroupShare
	}
	return ""
}

// Protobuf message implementation for struct RClaims
type RClaims struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainID   uint32 `protobuf:"varint,1,opt,name=ChainID,proto3" json:"ChainID,omitempty"`
	Height    uint32 `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
	Round     uint32 `protobuf:"varint,3,opt,name=Round,proto3" json:"Round,omitempty"`
	PrevBlock string `protobuf:"bytes,4,opt,name=PrevBlock,proto3" json:"PrevBlock,omitempty"`
}

func (x *RClaims) Reset() {
	*x = RClaims{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RClaims) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RClaims) ProtoMessage() {}

func (x *RClaims) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RClaims.ProtoReflect.Descriptor instead.
func (*RClaims) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{5}
}

func (x *RClaims) GetChainID() uint32 {
	if x != nil {
		return x.ChainID
	}
	return 0
}

func (x *RClaims) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *RClaims) GetRound() uint32 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *RClaims) GetPrevBlock() string {
	if x != nil {
		return x.PrevBlock
	}
	return ""
}

// Protobuf message implementation for struct BlockHeader
type BlockHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BClaims  *BClaims `protobuf:"bytes,1,opt,name=BClaims,proto3" json:"BClaims,omitempty"`
	SigGroup string   `protobuf:"bytes,2,opt,name=SigGroup,proto3" json:"SigGroup,omitempty"`
	TxHshLst []string `protobuf:"bytes,3,rep,name=TxHshLst,proto3" json:"TxHshLst,omitempty"`
}

func (x *BlockHeader) Reset() {
	*x = BlockHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHeader) ProtoMessage() {}

func (x *BlockHeader) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHeader.ProtoReflect.Descriptor instead.
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{6}
}

func (x *BlockHeader) GetBClaims() *BClaims {
	if x != nil {
		return x.BClaims
	}
	return nil
}

func (x *BlockHeader) GetSigGroup() string {
	if x != nil {
		return x.SigGroup
	}
	return ""
}

func (x *BlockHeader) GetTxHshLst() []string {
	if x != nil {
		return x.TxHshLst
	}
	return nil
}

// Protobuf message implementation for struct BClaims
type BClaims struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainID    uint32 `protobuf:"varint,1,opt,name=ChainID,proto3" json:"ChainID,omitempty"`
	Height     uint32 `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
	TxCount    uint32 `protobuf:"varint,3,opt,name=TxCount,proto3" json:"TxCount,omitempty"`
	PrevBlock  string `protobuf:"bytes,4,opt,name=PrevBlock,proto3" json:"PrevBlock,omitempty"`
	TxRoot     string `protobuf:"bytes,5,opt,name=TxRoot,proto3" json:"TxRoot,omitempty"`
	StateRoot  string `protobuf:"bytes,6,opt,name=StateRoot,proto3" json:"StateRoot,omitempty"`
	HeaderRoot string `protobuf:"bytes,7,opt,name=HeaderRoot,proto3" json:"HeaderRoot,omitempty"`
}

func (x *BClaims) Reset() {
	*x = BClaims{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BClaims) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BClaims) ProtoMessage() {}

func (x *BClaims) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BClaims.ProtoReflect.Descriptor instead.
func (*BClaims) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{7}
}

func (x *BClaims) GetChainID() uint32 {
	if x != nil {
		return x.ChainID
	}
	return 0
}

func (x *BClaims) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *BClaims) GetTxCount() uint32 {
	if x != nil {
		return x.TxCount
	}
	return 0
}

func (x *BClaims) GetPrevBlock() string {
	if x != nil {
		return x.PrevBlock
	}
	return ""
}

func (x *BClaims) GetTxRoot() string {
	if x != nil {
		return x.TxRoot
	}
	return ""
}

func (x *BClaims) GetStateRoot() string {
	if x != nil {
		return x.StateRoot
	}
	return ""
}

func (x *BClaims) GetHeaderRoot() string {
	if x != nil {
		return x.HeaderRoot
	}
	return ""
}

// Protobuf message implementation for struct PreVote
type PreVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proposal  *Proposal `protobuf:"bytes,1,opt,name=Proposal,proto3" json:"Proposal,omitempty"`
	Signature string    `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (x *PreVote) Reset() {
	*x = PreVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreVote) ProtoMessage() {}

func (x *PreVote) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreVote.ProtoReflect.Descriptor instead.
func (*PreVote) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{8}
}

func (x *PreVote) GetProposal() *Proposal {
	if x != nil {
		return x.Proposal
	}
	return nil
}

func (x *PreVote) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

// Protobuf message implementation for struct PClaims
type PClaims struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BClaims *BClaims `protobuf:"bytes,1,opt,name=BClaims,proto3" json:"BClaims,omitempty"`
	RCert   *RCert   `protobuf:"bytes,2,opt,name=RCert,proto3" json:"RCert,omitempty"`
}

func (x *PClaims) Reset() {
	*x = PClaims{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PClaims) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PClaims) ProtoMessage() {}

func (x *PClaims) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PClaims.ProtoReflect.Descriptor instead.
func (*PClaims) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{9}
}

func (x *PClaims) GetBClaims() *BClaims {
	if x != nil {
		return x.BClaims
	}
	return nil
}

func (x *PClaims) GetRCert() *RCert {
	if x != nil {
		return x.RCert
	}
	return nil
}

// Protobuf message implementation for struct PreCommit
type PreCommit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proposal  *Proposal `protobuf:"bytes,1,opt,name=Proposal,proto3" json:"Proposal,omitempty"`
	Signature string    `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
	PreVotes  []string  `protobuf:"bytes,3,rep,name=PreVotes,proto3" json:"PreVotes,omitempty"`
}

func (x *PreCommit) Reset() {
	*x = PreCommit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreCommit) ProtoMessage() {}

func (x *PreCommit) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreCommit.ProtoReflect.Descriptor instead.
func (*PreCommit) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{10}
}

func (x *PreCommit) GetProposal() *Proposal {
	if x != nil {
		return x.Proposal
	}
	return nil
}

func (x *PreCommit) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *PreCommit) GetPreVotes() []string {
	if x != nil {
		return x.PreVotes
	}
	return nil
}

// Protobuf message implementation for struct NextHeight
type NextHeight struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NHClaims   *NHClaims `protobuf:"bytes,1,opt,name=NHClaims,proto3" json:"NHClaims,omitempty"`
	Signature  string    `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
	PreCommits []string  `protobuf:"bytes,3,rep,name=PreCommits,proto3" json:"PreCommits,omitempty"`
}

func (x *NextHeight) Reset() {
	*x = NextHeight{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NextHeight) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextHeight) ProtoMessage() {}

func (x *NextHeight) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextHeight.ProtoReflect.Descriptor instead.
func (*NextHeight) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{11}
}

func (x *NextHeight) GetNHClaims() *NHClaims {
	if x != nil {
		return x.NHClaims
	}
	return nil
}

func (x *NextHeight) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *NextHeight) GetPreCommits() []string {
	if x != nil {
		return x.PreCommits
	}
	return nil
}

// Protobuf message implementation for struct NHClaims
type NHClaims struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proposal *Proposal `protobuf:"bytes,1,opt,name=Proposal,proto3" json:"Proposal,omitempty"`
	SigShare string    `protobuf:"bytes,2,opt,name=SigShare,proto3" json:"SigShare,omitempty"`
}

func (x *NHClaims) Reset() {
	*x = NHClaims{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NHClaims) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NHClaims) ProtoMessage() {}

func (x *NHClaims) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NHClaims.ProtoReflect.Descriptor instead.
func (*NHClaims) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{12}
}

func (x *NHClaims) GetProposal() *Proposal {
	if x != nil {
		return x.Proposal
	}
	return nil
}

func (x *NHClaims) GetSigShare() string {
	if x != nil {
		return x.SigShare
	}
	return ""
}

// Protobuf message implementation for struct NextRound
type NextRound struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NRClaims  *NRClaims `protobuf:"bytes,1,opt,name=NRClaims,proto3" json:"NRClaims,omitempty"`
	Signature string    `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (x *NextRound) Reset() {
	*x = NextRound{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cobjs_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NextRound) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextRound) ProtoMessage() {}

func (x *NextRound) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cobjs_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextRound.ProtoReflect.Descriptor instead.
func (*NextRound) Descriptor() ([]byte, []int) {
	return file_proto_cobjs_proto_rawDescGZIP(), []int{13}
}

func (x *NextRound) GetNRClaims() *NRClaims {
	if x != nil {
		return x.NRClaims
	}
	return nil
}

func (x *NextRound) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

var File_proto_cobjs_proto protoreflect.FileDescriptor

var file_proto_cobjs_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x62, 0x6a, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6e, 0x0a, 0x08, 0x50, 0x72,
	0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12, 0x28, 0x0a, 0x07, 0x50, 0x43, 0x6c, 0x61, 0x69, 0x6d,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x52, 0x07, 0x50, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x54, 0x78, 0x48, 0x73, 0x68, 0x4c, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x08, 0x54, 0x78, 0x48, 0x73, 0x68, 0x4c, 0x73, 0x74, 0x22, 0x4e, 0x0a, 0x0a, 0x50, 0x72,
	0x65, 0x56, 0x6f, 0x74, 0x65, 0x4e, 0x69, 0x6c, 0x12, 0x22, 0x0a, 0x05, 0x52, 0x43, 0x65, 0x72,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x52, 0x43, 0x65, 0x72, 0x74, 0x52, 0x05, 0x52, 0x43, 0x65, 0x72, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x50, 0x0a, 0x0c, 0x50, 0x72,
	0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x4e, 0x69, 0x6c, 0x12, 0x22, 0x0a, 0x05, 0x52, 0x43,
	0x65, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x52, 0x43, 0x65, 0x72, 0x74, 0x52, 0x05, 0x52, 0x43, 0x65, 0x72, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x4d, 0x0a, 0x05,
	0x52, 0x43, 0x65, 0x72, 0x74, 0x12, 0x28, 0x0a, 0x07, 0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x52, 0x07, 0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12,
	0x1a, 0x0a, 0x08, 0x53, 0x69, 0x67, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x53, 0x69, 0x67, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x94, 0x01, 0x0a, 0x08,
	0x4e, 0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12, 0x22, 0x0a, 0x05, 0x52, 0x43, 0x65, 0x72,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x52, 0x43, 0x65, 0x72, 0x74, 0x52, 0x05, 0x52, 0x43, 0x65, 0x72, 0x74, 0x12, 0x28, 0x0a, 0x07,
	0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x52, 0x07, 0x52,
	0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x69, 0x67, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x53, 0x69, 0x67, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x68, 0x61, 0x72, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x22, 0x6f, 0x0a, 0x07, 0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x72, 0x65, 0x76, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x50, 0x72, 0x65, 0x76, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x22, 0x6f, 0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x12, 0x28, 0x0a, 0x07, 0x42, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x43, 0x6c, 0x61,
	0x69, 0x6d, 0x73, 0x52, 0x07, 0x42, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x53, 0x69, 0x67, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x53, 0x69, 0x67, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x54, 0x78, 0x48, 0x73,
	0x68, 0x4c, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x54, 0x78, 0x48, 0x73,
	0x68, 0x4c, 0x73, 0x74, 0x22, 0xc9, 0x01, 0x0a, 0x07, 0x42, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x48, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x48, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x54, 0x78, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x54, 0x78, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x50, 0x72, 0x65, 0x76, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x50, 0x72, 0x65, 0x76, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x78,
	0x52, 0x6f, 0x6f, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x54, 0x78, 0x52, 0x6f,
	0x6f, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74,
	0x12, 0x1e, 0x0a, 0x0a, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x6f, 0x6f, 0x74,
	0x22, 0x54, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x2b, 0x0a, 0x08, 0x50,
	0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x52, 0x08,
	0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x57, 0x0a, 0x07, 0x50, 0x43, 0x6c, 0x61, 0x69, 0x6d,
	0x73, 0x12, 0x28, 0x0a, 0x07, 0x42, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x43, 0x6c, 0x61, 0x69,
	0x6d, 0x73, 0x52, 0x07, 0x42, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12, 0x22, 0x0a, 0x05, 0x52,
	0x43, 0x65, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x52, 0x43, 0x65, 0x72, 0x74, 0x52, 0x05, 0x52, 0x43, 0x65, 0x72, 0x74, 0x22,
	0x72, 0x0a, 0x09, 0x50, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x2b, 0x0a, 0x08,
	0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x52,
	0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x72, 0x65, 0x56, 0x6f,
	0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x50, 0x72, 0x65, 0x56, 0x6f,
	0x74, 0x65, 0x73, 0x22, 0x77, 0x0a, 0x0a, 0x4e, 0x65, 0x78, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x12, 0x2b, 0x0a, 0x08, 0x4e, 0x48, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x48, 0x43, 0x6c,
	0x61, 0x69, 0x6d, 0x73, 0x52, 0x08, 0x4e, 0x48, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12, 0x1c,
	0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1e, 0x0a, 0x0a,
	0x50, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0a, 0x50, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x73, 0x22, 0x53, 0x0a, 0x08,
	0x4e, 0x48, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12, 0x2b, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x70,
	0x6f, 0x73, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x52, 0x08, 0x50, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x69, 0x67, 0x53, 0x68, 0x61, 0x72,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x53, 0x69, 0x67, 0x53, 0x68, 0x61, 0x72,
	0x65, 0x22, 0x56, 0x0a, 0x09, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x2b,
	0x0a, 0x08, 0x4e, 0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d,
	0x73, 0x52, 0x08, 0x4e, 0x52, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x53,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42, 0x75, 0x0a, 0x09, 0x63, 0x6f, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x0a, 0x43, 0x6f, 0x62, 0x6a, 0x73, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x61, 0x6c, 0x69, 0x63, 0x65, 0x6e,
	0x65, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0xa2, 0x02,
	0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x05, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0xca, 0x02, 0x05, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0xe2, 0x02, 0x11, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x05, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_cobjs_proto_rawDescOnce sync.Once
	file_proto_cobjs_proto_rawDescData = file_proto_cobjs_proto_rawDesc
)

func file_proto_cobjs_proto_rawDescGZIP() []byte {
	file_proto_cobjs_proto_rawDescOnce.Do(func() {
		file_proto_cobjs_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_cobjs_proto_rawDescData)
	})
	return file_proto_cobjs_proto_rawDescData
}

var file_proto_cobjs_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_proto_cobjs_proto_goTypes = []interface{}{
	(*Proposal)(nil),     // 0: proto.Proposal
	(*PreVoteNil)(nil),   // 1: proto.PreVoteNil
	(*PreCommitNil)(nil), // 2: proto.PreCommitNil
	(*RCert)(nil),        // 3: proto.RCert
	(*NRClaims)(nil),     // 4: proto.NRClaims
	(*RClaims)(nil),      // 5: proto.RClaims
	(*BlockHeader)(nil),  // 6: proto.BlockHeader
	(*BClaims)(nil),      // 7: proto.BClaims
	(*PreVote)(nil),      // 8: proto.PreVote
	(*PClaims)(nil),      // 9: proto.PClaims
	(*PreCommit)(nil),    // 10: proto.PreCommit
	(*NextHeight)(nil),   // 11: proto.NextHeight
	(*NHClaims)(nil),     // 12: proto.NHClaims
	(*NextRound)(nil),    // 13: proto.NextRound
}
var file_proto_cobjs_proto_depIdxs = []int32{
	9,  // 0: proto.Proposal.PClaims:type_name -> proto.PClaims
	3,  // 1: proto.PreVoteNil.RCert:type_name -> proto.RCert
	3,  // 2: proto.PreCommitNil.RCert:type_name -> proto.RCert
	5,  // 3: proto.RCert.RClaims:type_name -> proto.RClaims
	3,  // 4: proto.NRClaims.RCert:type_name -> proto.RCert
	5,  // 5: proto.NRClaims.RClaims:type_name -> proto.RClaims
	7,  // 6: proto.BlockHeader.BClaims:type_name -> proto.BClaims
	0,  // 7: proto.PreVote.Proposal:type_name -> proto.Proposal
	7,  // 8: proto.PClaims.BClaims:type_name -> proto.BClaims
	3,  // 9: proto.PClaims.RCert:type_name -> proto.RCert
	0,  // 10: proto.PreCommit.Proposal:type_name -> proto.Proposal
	12, // 11: proto.NextHeight.NHClaims:type_name -> proto.NHClaims
	0,  // 12: proto.NHClaims.Proposal:type_name -> proto.Proposal
	4,  // 13: proto.NextRound.NRClaims:type_name -> proto.NRClaims
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_proto_cobjs_proto_init() }
func file_proto_cobjs_proto_init() {
	if File_proto_cobjs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_cobjs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Proposal); i {
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
		file_proto_cobjs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreVoteNil); i {
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
		file_proto_cobjs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreCommitNil); i {
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
		file_proto_cobjs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RCert); i {
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
		file_proto_cobjs_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NRClaims); i {
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
		file_proto_cobjs_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RClaims); i {
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
		file_proto_cobjs_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockHeader); i {
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
		file_proto_cobjs_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BClaims); i {
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
		file_proto_cobjs_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreVote); i {
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
		file_proto_cobjs_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PClaims); i {
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
		file_proto_cobjs_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreCommit); i {
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
		file_proto_cobjs_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NextHeight); i {
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
		file_proto_cobjs_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NHClaims); i {
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
		file_proto_cobjs_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NextRound); i {
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
			RawDescriptor: file_proto_cobjs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_cobjs_proto_goTypes,
		DependencyIndexes: file_proto_cobjs_proto_depIdxs,
		MessageInfos:      file_proto_cobjs_proto_msgTypes,
	}.Build()
	File_proto_cobjs_proto = out.File
	file_proto_cobjs_proto_rawDesc = nil
	file_proto_cobjs_proto_goTypes = nil
	file_proto_cobjs_proto_depIdxs = nil
}
