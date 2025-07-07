package models

import (
	proto "github.com/misnaged/gear-go/lib/server_grpc/proto"
)

type VoucherParams struct {
	Spender       string `json:"spender"`
	Balance       string `json:"balance"`
	CodeUploading bool   `json:"code_uploading"`
	Duration      int32  `json:"duration"`
}

type EncodedVoucherIssue struct {
	EncodedCall string `json:"encoded_call"`
}

func VoucherParamsToProto(params *VoucherParams) *proto.VoucherParams {
	pb := &proto.VoucherParams{
		Spender:       params.Spender,
		Balance:       params.Balance,
		CodeUploading: params.CodeUploading,
		Duration:      params.Duration,
	}
	return pb
}
func EncodedCallFromProto(pb *proto.EncodedVoucherIssue) *EncodedVoucherIssue {
	return &EncodedVoucherIssue{
		EncodedCall: pb.EncodedCall,
	}
}
