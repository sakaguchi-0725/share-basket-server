package handler

import (
	"context"
	proto "share-basket/personal-shopping/presentation/proto/gen"
	"share-basket/personal-shopping/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateShoppingItemHandler interface {
	Handle(ctx context.Context, req *proto.CreateShoppingItemRequest) (*proto.ShoppingItem, error)
}

type createShoppingItemHandler struct {
	usecase usecase.CreateShoppingItemUseCase
}

func (handler *createShoppingItemHandler) Handle(ctx context.Context, req *proto.CreateShoppingItemRequest) (*proto.ShoppingItem, error) {
	input := handler.makeInput(req)
	output, err := handler.usecase.Execute(input)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return handler.makeResponse(output), nil
}

func (handler *createShoppingItemHandler) makeInput(req *proto.CreateShoppingItemRequest) usecase.CreateShoppingItemInput {
	return usecase.CreateShoppingItemInput{
		Name:       req.GetName(),
		CategoryID: req.GetCategoryId(),
	}
}

func (handler *createShoppingItemHandler) makeResponse(output usecase.CreateShoppingItemOutput) *proto.ShoppingItem {
	return &proto.ShoppingItem{
		Id:     output.ID,
		Name:   output.Name,
		Status: proto.Status(proto.Status_value[output.Status]),
		Category: &proto.Category{
			Id:   output.Category.ID,
			Name: output.Category.Name,
		},
	}
}

func NewCreateShoppingItemHandler(usecase usecase.CreateShoppingItemUseCase) CreateShoppingItemHandler {
	return &createShoppingItemHandler{usecase}
}
