// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// P2PClient is the client API for P2P service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type P2PClient interface {
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	GetBlockHeaders(ctx context.Context, in *GetBlockHeadersRequest, opts ...grpc.CallOption) (*GetBlockHeadersResponse, error)
	GetMinedTxs(ctx context.Context, in *GetMinedTxsRequest, opts ...grpc.CallOption) (*GetMinedTxsResponse, error)
	GetPendingTxs(ctx context.Context, in *GetPendingTxsRequest, opts ...grpc.CallOption) (*GetPendingTxsResponse, error)
	GetSnapShotNode(ctx context.Context, in *GetSnapShotNodeRequest, opts ...grpc.CallOption) (*GetSnapShotNodeResponse, error)
	GetSnapShotStateData(ctx context.Context, in *GetSnapShotStateDataRequest, opts ...grpc.CallOption) (*GetSnapShotStateDataResponse, error)
	GetSnapShotHdrNode(ctx context.Context, in *GetSnapShotHdrNodeRequest, opts ...grpc.CallOption) (*GetSnapShotHdrNodeResponse, error)
	GossipTransaction(ctx context.Context, in *GossipTransactionMessage, opts ...grpc.CallOption) (*GossipTransactionAck, error)
	GossipProposal(ctx context.Context, in *GossipProposalMessage, opts ...grpc.CallOption) (*GossipProposalAck, error)
	GossipPreVote(ctx context.Context, in *GossipPreVoteMessage, opts ...grpc.CallOption) (*GossipPreVoteAck, error)
	GossipPreVoteNil(ctx context.Context, in *GossipPreVoteNilMessage, opts ...grpc.CallOption) (*GossipPreVoteNilAck, error)
	GossipPreCommit(ctx context.Context, in *GossipPreCommitMessage, opts ...grpc.CallOption) (*GossipPreCommitAck, error)
	GossipPreCommitNil(ctx context.Context, in *GossipPreCommitNilMessage, opts ...grpc.CallOption) (*GossipPreCommitNilAck, error)
	GossipNextRound(ctx context.Context, in *GossipNextRoundMessage, opts ...grpc.CallOption) (*GossipNextRoundAck, error)
	GossipNextHeight(ctx context.Context, in *GossipNextHeightMessage, opts ...grpc.CallOption) (*GossipNextHeightAck, error)
	GossipBlockHeader(ctx context.Context, in *GossipBlockHeaderMessage, opts ...grpc.CallOption) (*GossipBlockHeaderAck, error)
	GetPeers(ctx context.Context, in *GetPeersRequest, opts ...grpc.CallOption) (*GetPeersResponse, error)
}

type p2PClient struct {
	cc grpc.ClientConnInterface
}

func NewP2PClient(cc grpc.ClientConnInterface) P2PClient {
	return &p2PClient{cc}
}

func (c *p2PClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/proto.P2P/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GetBlockHeaders(ctx context.Context, in *GetBlockHeadersRequest, opts ...grpc.CallOption) (*GetBlockHeadersResponse, error) {
	out := new(GetBlockHeadersResponse)
	err := c.cc.Invoke(ctx, "/proto.P2P/GetBlockHeaders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GetMinedTxs(ctx context.Context, in *GetMinedTxsRequest, opts ...grpc.CallOption) (*GetMinedTxsResponse, error) {
	out := new(GetMinedTxsResponse)
	err := c.cc.Invoke(ctx, "/proto.P2P/GetMinedTxs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GetPendingTxs(ctx context.Context, in *GetPendingTxsRequest, opts ...grpc.CallOption) (*GetPendingTxsResponse, error) {
	out := new(GetPendingTxsResponse)
	err := c.cc.Invoke(ctx, "/proto.P2P/GetPendingTxs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GetSnapShotNode(ctx context.Context, in *GetSnapShotNodeRequest, opts ...grpc.CallOption) (*GetSnapShotNodeResponse, error) {
	out := new(GetSnapShotNodeResponse)
	err := c.cc.Invoke(ctx, "/proto.P2P/GetSnapShotNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GetSnapShotStateData(ctx context.Context, in *GetSnapShotStateDataRequest, opts ...grpc.CallOption) (*GetSnapShotStateDataResponse, error) {
	out := new(GetSnapShotStateDataResponse)
	err := c.cc.Invoke(ctx, "/proto.P2P/GetSnapShotStateData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GetSnapShotHdrNode(ctx context.Context, in *GetSnapShotHdrNodeRequest, opts ...grpc.CallOption) (*GetSnapShotHdrNodeResponse, error) {
	out := new(GetSnapShotHdrNodeResponse)
	err := c.cc.Invoke(ctx, "/proto.P2P/GetSnapShotHdrNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipTransaction(ctx context.Context, in *GossipTransactionMessage, opts ...grpc.CallOption) (*GossipTransactionAck, error) {
	out := new(GossipTransactionAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipProposal(ctx context.Context, in *GossipProposalMessage, opts ...grpc.CallOption) (*GossipProposalAck, error) {
	out := new(GossipProposalAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipProposal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipPreVote(ctx context.Context, in *GossipPreVoteMessage, opts ...grpc.CallOption) (*GossipPreVoteAck, error) {
	out := new(GossipPreVoteAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipPreVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipPreVoteNil(ctx context.Context, in *GossipPreVoteNilMessage, opts ...grpc.CallOption) (*GossipPreVoteNilAck, error) {
	out := new(GossipPreVoteNilAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipPreVoteNil", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipPreCommit(ctx context.Context, in *GossipPreCommitMessage, opts ...grpc.CallOption) (*GossipPreCommitAck, error) {
	out := new(GossipPreCommitAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipPreCommit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipPreCommitNil(ctx context.Context, in *GossipPreCommitNilMessage, opts ...grpc.CallOption) (*GossipPreCommitNilAck, error) {
	out := new(GossipPreCommitNilAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipPreCommitNil", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipNextRound(ctx context.Context, in *GossipNextRoundMessage, opts ...grpc.CallOption) (*GossipNextRoundAck, error) {
	out := new(GossipNextRoundAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipNextRound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipNextHeight(ctx context.Context, in *GossipNextHeightMessage, opts ...grpc.CallOption) (*GossipNextHeightAck, error) {
	out := new(GossipNextHeightAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipNextHeight", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GossipBlockHeader(ctx context.Context, in *GossipBlockHeaderMessage, opts ...grpc.CallOption) (*GossipBlockHeaderAck, error) {
	out := new(GossipBlockHeaderAck)
	err := c.cc.Invoke(ctx, "/proto.P2P/GossipBlockHeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) GetPeers(ctx context.Context, in *GetPeersRequest, opts ...grpc.CallOption) (*GetPeersResponse, error) {
	out := new(GetPeersResponse)
	err := c.cc.Invoke(ctx, "/proto.P2P/GetPeers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// P2PServer is the server API for P2P service.
// All implementations should embed UnimplementedP2PServer
// for forward compatibility
type P2PServer interface {
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	GetBlockHeaders(context.Context, *GetBlockHeadersRequest) (*GetBlockHeadersResponse, error)
	GetMinedTxs(context.Context, *GetMinedTxsRequest) (*GetMinedTxsResponse, error)
	GetPendingTxs(context.Context, *GetPendingTxsRequest) (*GetPendingTxsResponse, error)
	GetSnapShotNode(context.Context, *GetSnapShotNodeRequest) (*GetSnapShotNodeResponse, error)
	GetSnapShotStateData(context.Context, *GetSnapShotStateDataRequest) (*GetSnapShotStateDataResponse, error)
	GetSnapShotHdrNode(context.Context, *GetSnapShotHdrNodeRequest) (*GetSnapShotHdrNodeResponse, error)
	GossipTransaction(context.Context, *GossipTransactionMessage) (*GossipTransactionAck, error)
	GossipProposal(context.Context, *GossipProposalMessage) (*GossipProposalAck, error)
	GossipPreVote(context.Context, *GossipPreVoteMessage) (*GossipPreVoteAck, error)
	GossipPreVoteNil(context.Context, *GossipPreVoteNilMessage) (*GossipPreVoteNilAck, error)
	GossipPreCommit(context.Context, *GossipPreCommitMessage) (*GossipPreCommitAck, error)
	GossipPreCommitNil(context.Context, *GossipPreCommitNilMessage) (*GossipPreCommitNilAck, error)
	GossipNextRound(context.Context, *GossipNextRoundMessage) (*GossipNextRoundAck, error)
	GossipNextHeight(context.Context, *GossipNextHeightMessage) (*GossipNextHeightAck, error)
	GossipBlockHeader(context.Context, *GossipBlockHeaderMessage) (*GossipBlockHeaderAck, error)
	GetPeers(context.Context, *GetPeersRequest) (*GetPeersResponse, error)
}

// UnimplementedP2PServer should be embedded to have forward compatible implementations.
type UnimplementedP2PServer struct {
}

func (UnimplementedP2PServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedP2PServer) GetBlockHeaders(context.Context, *GetBlockHeadersRequest) (*GetBlockHeadersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockHeaders not implemented")
}
func (UnimplementedP2PServer) GetMinedTxs(context.Context, *GetMinedTxsRequest) (*GetMinedTxsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMinedTxs not implemented")
}
func (UnimplementedP2PServer) GetPendingTxs(context.Context, *GetPendingTxsRequest) (*GetPendingTxsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPendingTxs not implemented")
}
func (UnimplementedP2PServer) GetSnapShotNode(context.Context, *GetSnapShotNodeRequest) (*GetSnapShotNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSnapShotNode not implemented")
}
func (UnimplementedP2PServer) GetSnapShotStateData(context.Context, *GetSnapShotStateDataRequest) (*GetSnapShotStateDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSnapShotStateData not implemented")
}
func (UnimplementedP2PServer) GetSnapShotHdrNode(context.Context, *GetSnapShotHdrNodeRequest) (*GetSnapShotHdrNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSnapShotHdrNode not implemented")
}
func (UnimplementedP2PServer) GossipTransaction(context.Context, *GossipTransactionMessage) (*GossipTransactionAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipTransaction not implemented")
}
func (UnimplementedP2PServer) GossipProposal(context.Context, *GossipProposalMessage) (*GossipProposalAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipProposal not implemented")
}
func (UnimplementedP2PServer) GossipPreVote(context.Context, *GossipPreVoteMessage) (*GossipPreVoteAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipPreVote not implemented")
}
func (UnimplementedP2PServer) GossipPreVoteNil(context.Context, *GossipPreVoteNilMessage) (*GossipPreVoteNilAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipPreVoteNil not implemented")
}
func (UnimplementedP2PServer) GossipPreCommit(context.Context, *GossipPreCommitMessage) (*GossipPreCommitAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipPreCommit not implemented")
}
func (UnimplementedP2PServer) GossipPreCommitNil(context.Context, *GossipPreCommitNilMessage) (*GossipPreCommitNilAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipPreCommitNil not implemented")
}
func (UnimplementedP2PServer) GossipNextRound(context.Context, *GossipNextRoundMessage) (*GossipNextRoundAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipNextRound not implemented")
}
func (UnimplementedP2PServer) GossipNextHeight(context.Context, *GossipNextHeightMessage) (*GossipNextHeightAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipNextHeight not implemented")
}
func (UnimplementedP2PServer) GossipBlockHeader(context.Context, *GossipBlockHeaderMessage) (*GossipBlockHeaderAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GossipBlockHeader not implemented")
}
func (UnimplementedP2PServer) GetPeers(context.Context, *GetPeersRequest) (*GetPeersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeers not implemented")
}

// UnsafeP2PServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to P2PServer will
// result in compilation errors.
type UnsafeP2PServer interface {
	mustEmbedUnimplementedP2PServer()
}

func RegisterP2PServer(s grpc.ServiceRegistrar, srv P2PServer) {
	s.RegisterService(&P2P_ServiceDesc, srv)
}

func _P2P_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GetBlockHeaders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockHeadersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GetBlockHeaders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GetBlockHeaders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GetBlockHeaders(ctx, req.(*GetBlockHeadersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GetMinedTxs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMinedTxsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GetMinedTxs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GetMinedTxs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GetMinedTxs(ctx, req.(*GetMinedTxsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GetPendingTxs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPendingTxsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GetPendingTxs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GetPendingTxs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GetPendingTxs(ctx, req.(*GetPendingTxsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GetSnapShotNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSnapShotNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GetSnapShotNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GetSnapShotNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GetSnapShotNode(ctx, req.(*GetSnapShotNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GetSnapShotStateData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSnapShotStateDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GetSnapShotStateData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GetSnapShotStateData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GetSnapShotStateData(ctx, req.(*GetSnapShotStateDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GetSnapShotHdrNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSnapShotHdrNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GetSnapShotHdrNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GetSnapShotHdrNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GetSnapShotHdrNode(ctx, req.(*GetSnapShotHdrNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipTransactionMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipTransaction(ctx, req.(*GossipTransactionMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipProposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipProposalMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipProposal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipProposal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipProposal(ctx, req.(*GossipProposalMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipPreVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipPreVoteMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipPreVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipPreVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipPreVote(ctx, req.(*GossipPreVoteMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipPreVoteNil_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipPreVoteNilMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipPreVoteNil(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipPreVoteNil",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipPreVoteNil(ctx, req.(*GossipPreVoteNilMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipPreCommit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipPreCommitMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipPreCommit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipPreCommit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipPreCommit(ctx, req.(*GossipPreCommitMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipPreCommitNil_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipPreCommitNilMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipPreCommitNil(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipPreCommitNil",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipPreCommitNil(ctx, req.(*GossipPreCommitNilMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipNextRound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipNextRoundMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipNextRound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipNextRound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipNextRound(ctx, req.(*GossipNextRoundMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipNextHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipNextHeightMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipNextHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipNextHeight",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipNextHeight(ctx, req.(*GossipNextHeightMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GossipBlockHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipBlockHeaderMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GossipBlockHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GossipBlockHeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GossipBlockHeader(ctx, req.(*GossipBlockHeaderMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_GetPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPeersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).GetPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2P/GetPeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).GetPeers(ctx, req.(*GetPeersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// P2P_ServiceDesc is the grpc.ServiceDesc for P2P service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var P2P_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.P2P",
	HandlerType: (*P2PServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _P2P_Status_Handler,
		},
		{
			MethodName: "GetBlockHeaders",
			Handler:    _P2P_GetBlockHeaders_Handler,
		},
		{
			MethodName: "GetMinedTxs",
			Handler:    _P2P_GetMinedTxs_Handler,
		},
		{
			MethodName: "GetPendingTxs",
			Handler:    _P2P_GetPendingTxs_Handler,
		},
		{
			MethodName: "GetSnapShotNode",
			Handler:    _P2P_GetSnapShotNode_Handler,
		},
		{
			MethodName: "GetSnapShotStateData",
			Handler:    _P2P_GetSnapShotStateData_Handler,
		},
		{
			MethodName: "GetSnapShotHdrNode",
			Handler:    _P2P_GetSnapShotHdrNode_Handler,
		},
		{
			MethodName: "GossipTransaction",
			Handler:    _P2P_GossipTransaction_Handler,
		},
		{
			MethodName: "GossipProposal",
			Handler:    _P2P_GossipProposal_Handler,
		},
		{
			MethodName: "GossipPreVote",
			Handler:    _P2P_GossipPreVote_Handler,
		},
		{
			MethodName: "GossipPreVoteNil",
			Handler:    _P2P_GossipPreVoteNil_Handler,
		},
		{
			MethodName: "GossipPreCommit",
			Handler:    _P2P_GossipPreCommit_Handler,
		},
		{
			MethodName: "GossipPreCommitNil",
			Handler:    _P2P_GossipPreCommitNil_Handler,
		},
		{
			MethodName: "GossipNextRound",
			Handler:    _P2P_GossipNextRound_Handler,
		},
		{
			MethodName: "GossipNextHeight",
			Handler:    _P2P_GossipNextHeight_Handler,
		},
		{
			MethodName: "GossipBlockHeader",
			Handler:    _P2P_GossipBlockHeader_Handler,
		},
		{
			MethodName: "GetPeers",
			Handler:    _P2P_GetPeers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/p2p.proto",
}

// P2PDiscoveryClient is the client API for P2PDiscovery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type P2PDiscoveryClient interface {
	GetPeers(ctx context.Context, in *GetPeersRequest, opts ...grpc.CallOption) (*GetPeersResponse, error)
}

type p2PDiscoveryClient struct {
	cc grpc.ClientConnInterface
}

func NewP2PDiscoveryClient(cc grpc.ClientConnInterface) P2PDiscoveryClient {
	return &p2PDiscoveryClient{cc}
}

func (c *p2PDiscoveryClient) GetPeers(ctx context.Context, in *GetPeersRequest, opts ...grpc.CallOption) (*GetPeersResponse, error) {
	out := new(GetPeersResponse)
	err := c.cc.Invoke(ctx, "/proto.P2PDiscovery/GetPeers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// P2PDiscoveryServer is the server API for P2PDiscovery service.
// All implementations should embed UnimplementedP2PDiscoveryServer
// for forward compatibility
type P2PDiscoveryServer interface {
	GetPeers(context.Context, *GetPeersRequest) (*GetPeersResponse, error)
}

// UnimplementedP2PDiscoveryServer should be embedded to have forward compatible implementations.
type UnimplementedP2PDiscoveryServer struct {
}

func (UnimplementedP2PDiscoveryServer) GetPeers(context.Context, *GetPeersRequest) (*GetPeersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeers not implemented")
}

// UnsafeP2PDiscoveryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to P2PDiscoveryServer will
// result in compilation errors.
type UnsafeP2PDiscoveryServer interface {
	mustEmbedUnimplementedP2PDiscoveryServer()
}

func RegisterP2PDiscoveryServer(s grpc.ServiceRegistrar, srv P2PDiscoveryServer) {
	s.RegisterService(&P2PDiscovery_ServiceDesc, srv)
}

func _P2PDiscovery_GetPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPeersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PDiscoveryServer).GetPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.P2PDiscovery/GetPeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PDiscoveryServer).GetPeers(ctx, req.(*GetPeersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// P2PDiscovery_ServiceDesc is the grpc.ServiceDesc for P2PDiscovery service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var P2PDiscovery_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.P2PDiscovery",
	HandlerType: (*P2PDiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPeers",
			Handler:    _P2PDiscovery_GetPeers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/p2p.proto",
}