// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/v1/ballot.proto

package grpc_welaw_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type VoteValue int32

const (
	VoteValue_YES         VoteValue = 0
	VoteValue_NO          VoteValue = 1
	VoteValue_PRESENT     VoteValue = 2
	VoteValue_NOT_PRESENT VoteValue = 3
)

var VoteValue_name = map[int32]string{
	0: "YES",
	1: "NO",
	2: "PRESENT",
	3: "NOT_PRESENT",
}
var VoteValue_value = map[string]int32{
	"YES":         0,
	"NO":          1,
	"PRESENT":     2,
	"NOT_PRESENT": 3,
}

func (x VoteValue) String() string {
	return proto.EnumName(VoteValue_name, int32(x))
}
func (VoteValue) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

type GetVoteOptions_RequestType int32

const (
	GetVoteOptions_BY_USER_VERSION GetVoteOptions_RequestType = 0
)

var GetVoteOptions_RequestType_name = map[int32]string{
	0: "BY_USER_VERSION",
}
var GetVoteOptions_RequestType_value = map[string]int32{
	"BY_USER_VERSION": 0,
}

func (x GetVoteOptions_RequestType) String() string {
	return proto.EnumName(GetVoteOptions_RequestType_name, int32(x))
}
func (GetVoteOptions_RequestType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor2, []int{10, 0}
}

type Vote struct {
	Uid        string                     `protobuf:"bytes,1,opt,name=uid" json:"uid,omitempty"`
	LawId      string                     `protobuf:"bytes,2,opt,name=law_id,json=lawId" json:"law_id,omitempty"`
	UserId     string                     `protobuf:"bytes,3,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	Vote       string                     `protobuf:"bytes,4,opt,name=vote" json:"vote,omitempty"`
	Comment    string                     `protobuf:"bytes,5,opt,name=comment" json:"comment,omitempty"`
	Username   string                     `protobuf:"bytes,6,opt,name=username" json:"username,omitempty"`
	Upstream   string                     `protobuf:"bytes,7,opt,name=upstream" json:"upstream,omitempty"`
	Ident      string                     `protobuf:"bytes,8,opt,name=ident" json:"ident,omitempty"`
	Branch     string                     `protobuf:"bytes,9,opt,name=branch" json:"branch,omitempty"`
	Label      string                     `protobuf:"bytes,10,opt,name=label" json:"label,omitempty"`
	CastAt     *google_protobuf.Timestamp `protobuf:"bytes,11,opt,name=cast_at,json=castAt" json:"cast_at,omitempty"`
	Version    uint32                     `protobuf:"varint,12,opt,name=version" json:"version,omitempty"`
	VersionId  string                     `protobuf:"bytes,13,opt,name=version_id,json=versionId" json:"version_id,omitempty"`
	PictureUrl string                     `protobuf:"bytes,14,opt,name=picture_url,json=pictureUrl" json:"picture_url,omitempty"`
	LawTitle   string                     `protobuf:"bytes,15,opt,name=law_title,json=lawTitle" json:"law_title,omitempty"`
	Law        *LawSet                    `protobuf:"bytes,16,opt,name=law" json:"law,omitempty"`
	User       *User                      `protobuf:"bytes,17,opt,name=user" json:"user,omitempty"`
}

func (m *Vote) Reset()                    { *m = Vote{} }
func (m *Vote) String() string            { return proto.CompactTextString(m) }
func (*Vote) ProtoMessage()               {}
func (*Vote) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *Vote) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *Vote) GetLawId() string {
	if m != nil {
		return m.LawId
	}
	return ""
}

func (m *Vote) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Vote) GetVote() string {
	if m != nil {
		return m.Vote
	}
	return ""
}

func (m *Vote) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *Vote) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Vote) GetUpstream() string {
	if m != nil {
		return m.Upstream
	}
	return ""
}

func (m *Vote) GetIdent() string {
	if m != nil {
		return m.Ident
	}
	return ""
}

func (m *Vote) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *Vote) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *Vote) GetCastAt() *google_protobuf.Timestamp {
	if m != nil {
		return m.CastAt
	}
	return nil
}

func (m *Vote) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Vote) GetVersionId() string {
	if m != nil {
		return m.VersionId
	}
	return ""
}

func (m *Vote) GetPictureUrl() string {
	if m != nil {
		return m.PictureUrl
	}
	return ""
}

func (m *Vote) GetLawTitle() string {
	if m != nil {
		return m.LawTitle
	}
	return ""
}

func (m *Vote) GetLaw() *LawSet {
	if m != nil {
		return m.Law
	}
	return nil
}

func (m *Vote) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type VoteSummary struct {
	Yays               int32 `protobuf:"varint,1,opt,name=yays" json:"yays,omitempty"`
	Nays               int32 `protobuf:"varint,2,opt,name=nays" json:"nays,omitempty"`
	Present            int32 `protobuf:"varint,3,opt,name=present" json:"present,omitempty"`
	NotPresent         int32 `protobuf:"varint,4,opt,name=not_present,json=notPresent" json:"not_present,omitempty"`
	UpstreamYays       int32 `protobuf:"varint,5,opt,name=upstream_yays,json=upstreamYays" json:"upstream_yays,omitempty"`
	UpstreamNays       int32 `protobuf:"varint,6,opt,name=upstream_nays,json=upstreamNays" json:"upstream_nays,omitempty"`
	UpstreamPresent    int32 `protobuf:"varint,7,opt,name=upstream_present,json=upstreamPresent" json:"upstream_present,omitempty"`
	UpstreamNotPresent int32 `protobuf:"varint,8,opt,name=upstream_not_present,json=upstreamNotPresent" json:"upstream_not_present,omitempty"`
}

func (m *VoteSummary) Reset()                    { *m = VoteSummary{} }
func (m *VoteSummary) String() string            { return proto.CompactTextString(m) }
func (*VoteSummary) ProtoMessage()               {}
func (*VoteSummary) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *VoteSummary) GetYays() int32 {
	if m != nil {
		return m.Yays
	}
	return 0
}

func (m *VoteSummary) GetNays() int32 {
	if m != nil {
		return m.Nays
	}
	return 0
}

func (m *VoteSummary) GetPresent() int32 {
	if m != nil {
		return m.Present
	}
	return 0
}

func (m *VoteSummary) GetNotPresent() int32 {
	if m != nil {
		return m.NotPresent
	}
	return 0
}

func (m *VoteSummary) GetUpstreamYays() int32 {
	if m != nil {
		return m.UpstreamYays
	}
	return 0
}

func (m *VoteSummary) GetUpstreamNays() int32 {
	if m != nil {
		return m.UpstreamNays
	}
	return 0
}

func (m *VoteSummary) GetUpstreamPresent() int32 {
	if m != nil {
		return m.UpstreamPresent
	}
	return 0
}

func (m *VoteSummary) GetUpstreamNotPresent() int32 {
	if m != nil {
		return m.UpstreamNotPresent
	}
	return 0
}

type CreateVotesOptions struct {
	VoteResult *VoteResult `protobuf:"bytes,1,opt,name=vote_result,json=voteResult" json:"vote_result,omitempty"`
}

func (m *CreateVotesOptions) Reset()                    { *m = CreateVotesOptions{} }
func (m *CreateVotesOptions) String() string            { return proto.CompactTextString(m) }
func (*CreateVotesOptions) ProtoMessage()               {}
func (*CreateVotesOptions) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *CreateVotesOptions) GetVoteResult() *VoteResult {
	if m != nil {
		return m.VoteResult
	}
	return nil
}

type VoteResult struct {
	UpstreamGroup string                     `protobuf:"bytes,1,opt,name=upstream_group,json=upstreamGroup" json:"upstream_group,omitempty"`
	Result        string                     `protobuf:"bytes,2,opt,name=result" json:"result,omitempty"`
	PublishedAt   *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=published_at,json=publishedAt" json:"published_at,omitempty"`
	Upstream      string                     `protobuf:"bytes,4,opt,name=upstream" json:"upstream,omitempty"`
	Ident         string                     `protobuf:"bytes,5,opt,name=ident" json:"ident,omitempty"`
}

func (m *VoteResult) Reset()                    { *m = VoteResult{} }
func (m *VoteResult) String() string            { return proto.CompactTextString(m) }
func (*VoteResult) ProtoMessage()               {}
func (*VoteResult) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *VoteResult) GetUpstreamGroup() string {
	if m != nil {
		return m.UpstreamGroup
	}
	return ""
}

func (m *VoteResult) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func (m *VoteResult) GetPublishedAt() *google_protobuf.Timestamp {
	if m != nil {
		return m.PublishedAt
	}
	return nil
}

func (m *VoteResult) GetUpstream() string {
	if m != nil {
		return m.Upstream
	}
	return ""
}

func (m *VoteResult) GetIdent() string {
	if m != nil {
		return m.Ident
	}
	return ""
}

type CreateVoteOptions struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Branch   string `protobuf:"bytes,2,opt,name=branch" json:"branch,omitempty"`
	Version  string `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
}

func (m *CreateVoteOptions) Reset()                    { *m = CreateVoteOptions{} }
func (m *CreateVoteOptions) String() string            { return proto.CompactTextString(m) }
func (*CreateVoteOptions) ProtoMessage()               {}
func (*CreateVoteOptions) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *CreateVoteOptions) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *CreateVoteOptions) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *CreateVoteOptions) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type CreateVoteRequest struct {
	Vote *Vote              `protobuf:"bytes,1,opt,name=vote" json:"vote,omitempty"`
	Opts *CreateVoteOptions `protobuf:"bytes,2,opt,name=opts" json:"opts,omitempty"`
}

func (m *CreateVoteRequest) Reset()                    { *m = CreateVoteRequest{} }
func (m *CreateVoteRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateVoteRequest) ProtoMessage()               {}
func (*CreateVoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *CreateVoteRequest) GetVote() *Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (m *CreateVoteRequest) GetOpts() *CreateVoteOptions {
	if m != nil {
		return m.Opts
	}
	return nil
}

type CreateVoteReply struct {
	Vote *Vote  `protobuf:"bytes,1,opt,name=vote" json:"vote,omitempty"`
	Err  string `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *CreateVoteReply) Reset()                    { *m = CreateVoteReply{} }
func (m *CreateVoteReply) String() string            { return proto.CompactTextString(m) }
func (*CreateVoteReply) ProtoMessage()               {}
func (*CreateVoteReply) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6} }

func (m *CreateVoteReply) GetVote() *Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (m *CreateVoteReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type DeleteVoteOptions struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Branch   string `protobuf:"bytes,2,opt,name=branch" json:"branch,omitempty"`
	Version  string `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
}

func (m *DeleteVoteOptions) Reset()                    { *m = DeleteVoteOptions{} }
func (m *DeleteVoteOptions) String() string            { return proto.CompactTextString(m) }
func (*DeleteVoteOptions) ProtoMessage()               {}
func (*DeleteVoteOptions) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{7} }

func (m *DeleteVoteOptions) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *DeleteVoteOptions) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *DeleteVoteOptions) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type DeleteVoteRequest struct {
	Upstream string             `protobuf:"bytes,1,opt,name=upstream" json:"upstream,omitempty"`
	Ident    string             `protobuf:"bytes,2,opt,name=ident" json:"ident,omitempty"`
	Opts     *DeleteVoteOptions `protobuf:"bytes,3,opt,name=opts" json:"opts,omitempty"`
}

func (m *DeleteVoteRequest) Reset()                    { *m = DeleteVoteRequest{} }
func (m *DeleteVoteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteVoteRequest) ProtoMessage()               {}
func (*DeleteVoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{8} }

func (m *DeleteVoteRequest) GetUpstream() string {
	if m != nil {
		return m.Upstream
	}
	return ""
}

func (m *DeleteVoteRequest) GetIdent() string {
	if m != nil {
		return m.Ident
	}
	return ""
}

func (m *DeleteVoteRequest) GetOpts() *DeleteVoteOptions {
	if m != nil {
		return m.Opts
	}
	return nil
}

type DeleteVoteReply struct {
	Err string `protobuf:"bytes,1,opt,name=err" json:"err,omitempty"`
}

func (m *DeleteVoteReply) Reset()                    { *m = DeleteVoteReply{} }
func (m *DeleteVoteReply) String() string            { return proto.CompactTextString(m) }
func (*DeleteVoteReply) ProtoMessage()               {}
func (*DeleteVoteReply) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{9} }

func (m *DeleteVoteReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type GetVoteOptions struct {
	ReqType  GetVoteOptions_RequestType `protobuf:"varint,1,opt,name=req_type,json=reqType,enum=grpc.welaw.v1.GetVoteOptions_RequestType" json:"req_type,omitempty"`
	Upstream string                     `protobuf:"bytes,2,opt,name=upstream" json:"upstream,omitempty"`
	Ident    string                     `protobuf:"bytes,3,opt,name=ident" json:"ident,omitempty"`
	Username string                     `protobuf:"bytes,4,opt,name=username" json:"username,omitempty"`
	Branch   string                     `protobuf:"bytes,5,opt,name=branch" json:"branch,omitempty"`
	Version  string                     `protobuf:"bytes,6,opt,name=version" json:"version,omitempty"`
	Quiet    bool                       `protobuf:"varint,7,opt,name=quiet" json:"quiet,omitempty"`
}

func (m *GetVoteOptions) Reset()                    { *m = GetVoteOptions{} }
func (m *GetVoteOptions) String() string            { return proto.CompactTextString(m) }
func (*GetVoteOptions) ProtoMessage()               {}
func (*GetVoteOptions) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{10} }

func (m *GetVoteOptions) GetReqType() GetVoteOptions_RequestType {
	if m != nil {
		return m.ReqType
	}
	return GetVoteOptions_BY_USER_VERSION
}

func (m *GetVoteOptions) GetUpstream() string {
	if m != nil {
		return m.Upstream
	}
	return ""
}

func (m *GetVoteOptions) GetIdent() string {
	if m != nil {
		return m.Ident
	}
	return ""
}

func (m *GetVoteOptions) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *GetVoteOptions) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *GetVoteOptions) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *GetVoteOptions) GetQuiet() bool {
	if m != nil {
		return m.Quiet
	}
	return false
}

type GetVoteRequest struct {
	Opts *GetVoteOptions `protobuf:"bytes,1,opt,name=opts" json:"opts,omitempty"`
}

func (m *GetVoteRequest) Reset()                    { *m = GetVoteRequest{} }
func (m *GetVoteRequest) String() string            { return proto.CompactTextString(m) }
func (*GetVoteRequest) ProtoMessage()               {}
func (*GetVoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{11} }

func (m *GetVoteRequest) GetOpts() *GetVoteOptions {
	if m != nil {
		return m.Opts
	}
	return nil
}

type GetVoteReply struct {
	Vote *Vote  `protobuf:"bytes,1,opt,name=vote" json:"vote,omitempty"`
	Err  string `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *GetVoteReply) Reset()                    { *m = GetVoteReply{} }
func (m *GetVoteReply) String() string            { return proto.CompactTextString(m) }
func (*GetVoteReply) ProtoMessage()               {}
func (*GetVoteReply) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{12} }

func (m *GetVoteReply) GetVote() *Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (m *GetVoteReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type ListVotesOptions struct {
	Category string `protobuf:"bytes,1,opt,name=category" json:"category,omitempty"`
	Upstream string `protobuf:"bytes,2,opt,name=upstream" json:"upstream,omitempty"`
	Ident    string `protobuf:"bytes,3,opt,name=ident" json:"ident,omitempty"`
	Branch   string `protobuf:"bytes,4,opt,name=branch" json:"branch,omitempty"`
	Version  uint32 `protobuf:"varint,5,opt,name=version" json:"version,omitempty"`
	Username string `protobuf:"bytes,6,opt,name=username" json:"username,omitempty"`
	PageSize uint32 `protobuf:"varint,7,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
	PageNum  uint32 `protobuf:"varint,8,opt,name=page_num,json=pageNum" json:"page_num,omitempty"`
}

func (m *ListVotesOptions) Reset()                    { *m = ListVotesOptions{} }
func (m *ListVotesOptions) String() string            { return proto.CompactTextString(m) }
func (*ListVotesOptions) ProtoMessage()               {}
func (*ListVotesOptions) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{13} }

func (m *ListVotesOptions) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

func (m *ListVotesOptions) GetUpstream() string {
	if m != nil {
		return m.Upstream
	}
	return ""
}

func (m *ListVotesOptions) GetIdent() string {
	if m != nil {
		return m.Ident
	}
	return ""
}

func (m *ListVotesOptions) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *ListVotesOptions) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ListVotesOptions) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *ListVotesOptions) GetPageSize() uint32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListVotesOptions) GetPageNum() uint32 {
	if m != nil {
		return m.PageNum
	}
	return 0
}

type ListVotesRequest struct {
	Opts *ListVotesOptions `protobuf:"bytes,1,opt,name=opts" json:"opts,omitempty"`
}

func (m *ListVotesRequest) Reset()                    { *m = ListVotesRequest{} }
func (m *ListVotesRequest) String() string            { return proto.CompactTextString(m) }
func (*ListVotesRequest) ProtoMessage()               {}
func (*ListVotesRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{14} }

func (m *ListVotesRequest) GetOpts() *ListVotesOptions {
	if m != nil {
		return m.Opts
	}
	return nil
}

type ListVotesReply struct {
	Votes []*Vote `protobuf:"bytes,1,rep,name=votes" json:"votes,omitempty"`
	Total int32   `protobuf:"varint,2,opt,name=total" json:"total,omitempty"`
	Err   string  `protobuf:"bytes,3,opt,name=err" json:"err,omitempty"`
}

func (m *ListVotesReply) Reset()                    { *m = ListVotesReply{} }
func (m *ListVotesReply) String() string            { return proto.CompactTextString(m) }
func (*ListVotesReply) ProtoMessage()               {}
func (*ListVotesReply) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{15} }

func (m *ListVotesReply) GetVotes() []*Vote {
	if m != nil {
		return m.Votes
	}
	return nil
}

func (m *ListVotesReply) GetTotal() int32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *ListVotesReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type UpdateVoteOptions struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Branch   string `protobuf:"bytes,2,opt,name=branch" json:"branch,omitempty"`
	Version  string `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
}

func (m *UpdateVoteOptions) Reset()                    { *m = UpdateVoteOptions{} }
func (m *UpdateVoteOptions) String() string            { return proto.CompactTextString(m) }
func (*UpdateVoteOptions) ProtoMessage()               {}
func (*UpdateVoteOptions) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{16} }

func (m *UpdateVoteOptions) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UpdateVoteOptions) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *UpdateVoteOptions) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type UpdateVoteRequest struct {
	Vote *Vote              `protobuf:"bytes,1,opt,name=vote" json:"vote,omitempty"`
	Opts *UpdateVoteOptions `protobuf:"bytes,2,opt,name=opts" json:"opts,omitempty"`
}

func (m *UpdateVoteRequest) Reset()                    { *m = UpdateVoteRequest{} }
func (m *UpdateVoteRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateVoteRequest) ProtoMessage()               {}
func (*UpdateVoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{17} }

func (m *UpdateVoteRequest) GetVote() *Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (m *UpdateVoteRequest) GetOpts() *UpdateVoteOptions {
	if m != nil {
		return m.Opts
	}
	return nil
}

type UpdateVoteReply struct {
	Vote *Vote  `protobuf:"bytes,1,opt,name=vote" json:"vote,omitempty"`
	Err  string `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *UpdateVoteReply) Reset()                    { *m = UpdateVoteReply{} }
func (m *UpdateVoteReply) String() string            { return proto.CompactTextString(m) }
func (*UpdateVoteReply) ProtoMessage()               {}
func (*UpdateVoteReply) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{18} }

func (m *UpdateVoteReply) GetVote() *Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (m *UpdateVoteReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto.RegisterType((*Vote)(nil), "grpc.welaw.v1.Vote")
	proto.RegisterType((*VoteSummary)(nil), "grpc.welaw.v1.VoteSummary")
	proto.RegisterType((*CreateVotesOptions)(nil), "grpc.welaw.v1.CreateVotesOptions")
	proto.RegisterType((*VoteResult)(nil), "grpc.welaw.v1.VoteResult")
	proto.RegisterType((*CreateVoteOptions)(nil), "grpc.welaw.v1.CreateVoteOptions")
	proto.RegisterType((*CreateVoteRequest)(nil), "grpc.welaw.v1.CreateVoteRequest")
	proto.RegisterType((*CreateVoteReply)(nil), "grpc.welaw.v1.CreateVoteReply")
	proto.RegisterType((*DeleteVoteOptions)(nil), "grpc.welaw.v1.DeleteVoteOptions")
	proto.RegisterType((*DeleteVoteRequest)(nil), "grpc.welaw.v1.DeleteVoteRequest")
	proto.RegisterType((*DeleteVoteReply)(nil), "grpc.welaw.v1.DeleteVoteReply")
	proto.RegisterType((*GetVoteOptions)(nil), "grpc.welaw.v1.GetVoteOptions")
	proto.RegisterType((*GetVoteRequest)(nil), "grpc.welaw.v1.GetVoteRequest")
	proto.RegisterType((*GetVoteReply)(nil), "grpc.welaw.v1.GetVoteReply")
	proto.RegisterType((*ListVotesOptions)(nil), "grpc.welaw.v1.ListVotesOptions")
	proto.RegisterType((*ListVotesRequest)(nil), "grpc.welaw.v1.ListVotesRequest")
	proto.RegisterType((*ListVotesReply)(nil), "grpc.welaw.v1.ListVotesReply")
	proto.RegisterType((*UpdateVoteOptions)(nil), "grpc.welaw.v1.UpdateVoteOptions")
	proto.RegisterType((*UpdateVoteRequest)(nil), "grpc.welaw.v1.UpdateVoteRequest")
	proto.RegisterType((*UpdateVoteReply)(nil), "grpc.welaw.v1.UpdateVoteReply")
	proto.RegisterEnum("grpc.welaw.v1.VoteValue", VoteValue_name, VoteValue_value)
	proto.RegisterEnum("grpc.welaw.v1.GetVoteOptions_RequestType", GetVoteOptions_RequestType_name, GetVoteOptions_RequestType_value)
}

func init() { proto.RegisterFile("api/v1/ballot.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 1043 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x96, 0xcd, 0x6e, 0xdb, 0x46,
	0x10, 0xc7, 0x43, 0x51, 0x9f, 0x43, 0xcb, 0xa2, 0xd7, 0x49, 0xcb, 0x38, 0x08, 0x6c, 0x30, 0x28,
	0x62, 0xf7, 0x20, 0xd7, 0x76, 0x4f, 0x01, 0x7a, 0x48, 0x13, 0xc3, 0x30, 0x60, 0xc8, 0x06, 0x65,
	0x1b, 0xf0, 0x89, 0x58, 0x49, 0x5b, 0x85, 0xc0, 0x52, 0xa4, 0x97, 0x4b, 0x0b, 0x4a, 0x5f, 0xa2,
	0xb7, 0xbe, 0x4c, 0x9f, 0xa7, 0xe7, 0x3e, 0x42, 0x31, 0x4b, 0x2e, 0x45, 0xca, 0x51, 0x5a, 0x34,
	0xcd, 0x6d, 0x67, 0xf6, 0xaf, 0x9d, 0x9d, 0xdf, 0xce, 0x0c, 0x05, 0xdb, 0x34, 0x0e, 0x0e, 0x1f,
	0x8e, 0x0e, 0x47, 0x94, 0xf3, 0x48, 0xf6, 0x63, 0x11, 0xc9, 0x88, 0x74, 0xa7, 0x22, 0x1e, 0xf7,
	0xe7, 0x8c, 0xd3, 0x79, 0xff, 0xe1, 0x68, 0x67, 0x77, 0x1a, 0x45, 0x53, 0xce, 0x0e, 0xd5, 0xe6,
	0x28, 0xfd, 0xe5, 0x50, 0x06, 0x21, 0x4b, 0x24, 0x0d, 0xe3, 0x4c, 0xbf, 0x63, 0xe7, 0x87, 0xa0,
	0x3e, 0xf3, 0x6c, 0xe5, 0x9e, 0x34, 0x61, 0x22, 0x73, 0xb9, 0x7f, 0x99, 0x50, 0xbf, 0x8d, 0x24,
	0x23, 0x36, 0x98, 0x69, 0x30, 0x71, 0x8c, 0x3d, 0x63, 0xbf, 0xe3, 0xe1, 0x92, 0x3c, 0x83, 0x26,
	0xa7, 0x73, 0x3f, 0x98, 0x38, 0x35, 0xe5, 0x6c, 0x70, 0x3a, 0x3f, 0x9f, 0x90, 0x6f, 0xa1, 0x85,
	0xbf, 0x47, 0xbf, 0xa9, 0xfc, 0x4d, 0x34, 0xcf, 0x27, 0x84, 0x40, 0xfd, 0x21, 0x92, 0xcc, 0xa9,
	0x2b, 0xaf, 0x5a, 0x13, 0x07, 0x5a, 0xe3, 0x28, 0x0c, 0xd9, 0x4c, 0x3a, 0x0d, 0xe5, 0xd6, 0x26,
	0xd9, 0x81, 0x36, 0xfe, 0x6e, 0x46, 0x43, 0xe6, 0x34, 0xd5, 0x56, 0x61, 0xab, 0xbd, 0x38, 0x91,
	0x82, 0xd1, 0xd0, 0x69, 0xe5, 0x7b, 0xb9, 0x4d, 0x9e, 0x42, 0x23, 0x98, 0xe0, 0x79, 0xed, 0xec,
	0x52, 0xca, 0x20, 0xdf, 0x40, 0x73, 0x24, 0xe8, 0x6c, 0xfc, 0xc1, 0xe9, 0x64, 0x77, 0xca, 0x2c,
	0x54, 0x73, 0x3a, 0x62, 0xdc, 0x01, 0x9d, 0xc2, 0x88, 0x71, 0x72, 0x02, 0xad, 0x31, 0x4d, 0xa4,
	0x4f, 0xa5, 0x63, 0xed, 0x19, 0xfb, 0xd6, 0xf1, 0x4e, 0x3f, 0x83, 0xd9, 0xd7, 0x30, 0xfb, 0xd7,
	0x1a, 0xa6, 0xd7, 0x44, 0xe9, 0x5b, 0x89, 0xa9, 0x3c, 0x30, 0x91, 0x04, 0xd1, 0xcc, 0xd9, 0xd8,
	0x33, 0xf6, 0xbb, 0x9e, 0x36, 0xc9, 0x4b, 0x80, 0x7c, 0x89, 0x50, 0xba, 0x2a, 0x52, 0x27, 0xf7,
	0x9c, 0x4f, 0xc8, 0x2e, 0x58, 0x71, 0x30, 0x96, 0xa9, 0x60, 0x7e, 0x2a, 0xb8, 0xb3, 0xa9, 0xf6,
	0x21, 0x77, 0xdd, 0x08, 0x4e, 0x5e, 0x40, 0x07, 0x41, 0xcb, 0x40, 0x72, 0xe6, 0xf4, 0xb2, 0x7c,
	0x39, 0x9d, 0x5f, 0xa3, 0x4d, 0x5e, 0x83, 0xc9, 0xe9, 0xdc, 0xb1, 0xd5, 0x3d, 0x9f, 0xf5, 0x2b,
	0x35, 0xd0, 0xbf, 0xa0, 0xf3, 0x21, 0x93, 0x1e, 0x2a, 0xc8, 0x6b, 0xa8, 0x23, 0x40, 0x67, 0x4b,
	0x29, 0xb7, 0x57, 0x94, 0x37, 0x09, 0x13, 0x9e, 0x12, 0xb8, 0xbf, 0xd7, 0xc0, 0xc2, 0x27, 0x1f,
	0xa6, 0x61, 0x48, 0xc5, 0x02, 0xdf, 0x6d, 0x41, 0x17, 0x89, 0x7a, 0xfa, 0x86, 0xa7, 0xd6, 0xe8,
	0x9b, 0xa1, 0xaf, 0x96, 0xf9, 0x70, 0x8d, 0x00, 0x62, 0xc1, 0x12, 0x64, 0x6f, 0x2a, 0xb7, 0x36,
	0x31, 0xc3, 0x59, 0x24, 0x7d, 0xbd, 0x5b, 0x57, 0xbb, 0x30, 0x8b, 0xe4, 0x55, 0x2e, 0x78, 0x05,
	0x5d, 0xfd, 0x80, 0xbe, 0x8a, 0xd5, 0x50, 0x92, 0x0d, 0xed, 0xbc, 0xc3, 0xf3, 0xcb, 0x22, 0x15,
	0xbc, 0x59, 0x15, 0x0d, 0x50, 0x74, 0x00, 0x76, 0x21, 0xd2, 0xf1, 0x5a, 0x4a, 0xd7, 0xd3, 0x7e,
	0x1d, 0xf4, 0x07, 0x78, 0xba, 0x3c, 0xaf, 0x74, 0xbd, 0xb6, 0x92, 0x93, 0xe2, 0xd8, 0xe2, 0x9a,
	0xee, 0x15, 0x90, 0x77, 0x82, 0x51, 0xc9, 0x10, 0x4f, 0x72, 0x19, 0xcb, 0x20, 0x9a, 0x25, 0xe4,
	0x0d, 0x58, 0x58, 0xcb, 0xbe, 0x60, 0x49, 0xca, 0xa5, 0xc2, 0x64, 0x1d, 0x3f, 0x5f, 0xe1, 0x8b,
	0xbf, 0xf0, 0x94, 0xc0, 0x83, 0x87, 0x62, 0xed, 0xfe, 0x61, 0x00, 0x2c, 0xb7, 0xc8, 0x77, 0xb0,
	0x59, 0x5c, 0x69, 0x2a, 0xa2, 0x34, 0xce, 0xfb, 0xad, 0x48, 0xfc, 0x0c, 0x9d, 0x58, 0xcd, 0x79,
	0xb0, 0xac, 0xf3, 0x72, 0x8b, 0xfc, 0x04, 0x1b, 0x71, 0x3a, 0xe2, 0x41, 0xf2, 0x81, 0x4d, 0xb0,
	0x78, 0xcd, 0x7f, 0x2c, 0x5e, 0xab, 0xd0, 0xbf, 0x95, 0x95, 0xb6, 0xaa, 0xaf, 0x6b, 0xab, 0x46,
	0xa9, 0xad, 0x5c, 0x0a, 0x5b, 0x4b, 0x20, 0x9a, 0x47, 0xb9, 0x73, 0x8d, 0x95, 0xce, 0x5d, 0xf6,
	0x61, 0xad, 0xd2, 0x87, 0xa5, 0xe6, 0xc9, 0x86, 0x86, 0x36, 0x5d, 0x51, 0x0e, 0xe1, 0xb1, 0xfb,
	0x94, 0x25, 0x12, 0x6b, 0x59, 0x8d, 0x12, 0xe3, 0x93, 0xb5, 0xac, 0x94, 0xd9, 0x7c, 0xf9, 0x11,
	0xea, 0x51, 0x2c, 0xb3, 0x3a, 0xb5, 0x8e, 0xf7, 0x56, 0x84, 0x8f, 0xee, 0xee, 0x29, 0xb5, 0x7b,
	0x01, 0xbd, 0x72, 0xcc, 0x98, 0x2f, 0xfe, 0x7d, 0x44, 0x1b, 0x4c, 0x26, 0x44, 0x9e, 0x1e, 0x2e,
	0x11, 0xd2, 0x7b, 0xc6, 0xd9, 0xd7, 0x84, 0xf4, 0x6b, 0x39, 0x84, 0x86, 0x54, 0x7e, 0x4e, 0x63,
	0xdd, 0x73, 0xd6, 0xca, 0x53, 0x52, 0xd3, 0x32, 0x3f, 0x49, 0xeb, 0x51, 0x12, 0x39, 0xad, 0x57,
	0xd0, 0x2b, 0x07, 0x47, 0x5a, 0x39, 0x04, 0x63, 0x09, 0xe1, 0xb7, 0x1a, 0x6c, 0x9e, 0x31, 0x59,
	0x46, 0xf0, 0x1e, 0xda, 0x82, 0xdd, 0xfb, 0x72, 0x11, 0x67, 0x08, 0x36, 0x8f, 0x0f, 0x56, 0x22,
	0x56, 0x7f, 0xd0, 0xcf, 0x13, 0xbb, 0x5e, 0xc4, 0xcc, 0x6b, 0x09, 0x76, 0x8f, 0x8b, 0x4a, 0x96,
	0xb5, 0x75, 0x59, 0x9a, 0xe5, 0x2c, 0xcb, 0xe8, 0xeb, 0x6b, 0xd1, 0x37, 0xd6, 0xa1, 0x6f, 0x56,
	0xd0, 0x63, 0x8c, 0xfb, 0x34, 0x60, 0xd9, 0x94, 0x69, 0x7b, 0x99, 0xe1, 0xba, 0x60, 0x95, 0x6e,
	0x4b, 0xb6, 0xa1, 0xf7, 0xf3, 0x9d, 0x7f, 0x33, 0x3c, 0xf5, 0xfc, 0xdb, 0x53, 0x6f, 0x78, 0x7e,
	0x39, 0xb0, 0x9f, 0xb8, 0xef, 0x0a, 0x22, 0xfa, 0xc5, 0x8e, 0x72, 0xfe, 0x59, 0x91, 0xbd, 0xfc,
	0x2c, 0x8d, 0x1c, 0xfe, 0x39, 0x6c, 0x14, 0x87, 0x7c, 0x61, 0x9d, 0xfe, 0x69, 0x80, 0x7d, 0x11,
	0x24, 0xb2, 0x32, 0xdc, 0x76, 0xa0, 0x3d, 0xa6, 0x92, 0x4d, 0x23, 0xb1, 0xd0, 0x45, 0xa4, 0xed,
	0xff, 0x80, 0x7e, 0x89, 0xb7, 0xbe, 0x0e, 0x6f, 0xa3, 0xfa, 0xed, 0xfc, 0xdc, 0xdf, 0x80, 0x17,
	0xd0, 0x89, 0xe9, 0x94, 0xf9, 0x49, 0xf0, 0x91, 0x29, 0xfc, 0x5d, 0xaf, 0x8d, 0x8e, 0x61, 0xf0,
	0x91, 0x91, 0xe7, 0xa0, 0xd6, 0xfe, 0x2c, 0x0d, 0xd5, 0x44, 0xef, 0x7a, 0x2d, 0xb4, 0x07, 0x69,
	0xe8, 0x9e, 0x95, 0xf2, 0xd4, 0xe8, 0x4f, 0x2a, 0xe8, 0x77, 0x57, 0xbf, 0xa3, 0x2b, 0x58, 0x72,
	0xf8, 0x63, 0xd8, 0x2c, 0x1d, 0x84, 0xf8, 0x0f, 0xa0, 0x81, 0x74, 0xf1, 0x1c, 0x73, 0x1d, 0xff,
	0x4c, 0x81, 0x84, 0x64, 0x24, 0x29, 0xcf, 0xbf, 0xa1, 0x99, 0xa1, 0x9f, 0xc5, 0xac, 0x8c, 0x8f,
	0x9b, 0x78, 0xf2, 0xb5, 0x67, 0xec, 0x32, 0xc4, 0xff, 0x3c, 0x63, 0x1f, 0xdd, 0x7d, 0x39, 0x63,
	0xcb, 0x31, 0xbf, 0xac, 0x76, 0xbf, 0x7f, 0x03, 0x1d, 0xdc, 0xbf, 0xa5, 0x3c, 0x65, 0xa4, 0x05,
	0xe6, 0xdd, 0xe9, 0xd0, 0x7e, 0x42, 0x9a, 0x50, 0x1b, 0x5c, 0xda, 0x06, 0xb1, 0xa0, 0x75, 0xe5,
	0x9d, 0x0e, 0x4f, 0x07, 0xd7, 0x76, 0x8d, 0xf4, 0xc0, 0x1a, 0x5c, 0x5e, 0xfb, 0xda, 0x61, 0x8e,
	0x9a, 0xea, 0xbb, 0x78, 0xf2, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x52, 0x64, 0xeb, 0xf8, 0x55,
	0x0b, 0x00, 0x00,
}