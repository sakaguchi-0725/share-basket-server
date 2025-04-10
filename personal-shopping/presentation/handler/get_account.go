package handler

import (
	"context"
	proto "share-basket/personal-shopping/presentation/proto/gen"
	"share-basket/personal-shopping/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetAccountHandler interface {
	Handle(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error)
}

type getAccountHandler struct {
	usecase usecase.GetAccountUseCase
}

func (handler *getAccountHandler) Handle(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	output, err := handler.usecase.Execute(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.GetAccountResponse{
		Id:   output.ID,
		Name: output.Name,
	}, nil
}

func NewGetAccountHandler(usecase usecase.GetAccountUseCase) GetAccountHandler {
	return &getAccountHandler{usecase}
}
