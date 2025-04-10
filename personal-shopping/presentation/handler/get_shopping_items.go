package handler

import (
	"context"
	proto "share-basket/personal-shopping/presentation/proto/gen"
	"share-basket/personal-shopping/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetShoppingItemsHandler interface {
	Handle(ctx context.Context, req *proto.GetShoppingItemsRequest) (*proto.GetShoppingItemsResponse, error)
}

type getShoppingItemsHandler struct {
	usecase usecase.GetShoppingItemsUseCase
}

func (handler *getShoppingItemsHandler) Handle(ctx context.Context, req *proto.GetShoppingItemsRequest) (*proto.GetShoppingItemsResponse, error) {
	outputs, err := handler.usecase.Execute(req.Status.String())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return handler.makeResponse(outputs), nil
}

func (handle *getShoppingItemsHandler) makeResponse(outputs []usecase.GetShoppingItemOutput) *proto.GetShoppingItemsResponse {
	items := make([]*proto.ShoppingItem, len(outputs))
	for i, output := range outputs {
		items[i] = &proto.ShoppingItem{
			Id:     output.ID,
			Name:   output.Name,
			Status: proto.Status(proto.Status_value[output.Status]),
			Category: &proto.Category{
				Id:   output.Category.ID,
				Name: output.Category.Name,
			},
		}
	}

	return &proto.GetShoppingItemsResponse{Items: items}
}

func NewGetShoppingItemsHandler(usecase usecase.GetShoppingItemsUseCase) GetShoppingItemsHandler {
	return &getShoppingItemsHandler{usecase}
}
