package registry

import (
	"context"
	proto "share-basket/personal-shopping/presentation/proto/gen"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Services struct {
	PersonalShopping proto.PersonalShoppingServiceServer
	Account          proto.AccountServiceServer
}

func NewServices(container *container) Services {
	return Services{
		PersonalShopping: newPersonalShoppingService(container),
		Account:          newAccountService(container),
	}
}

type personalShoppingService struct {
	container *container
	proto.UnimplementedPersonalShoppingServiceServer
}

func (service *personalShoppingService) GetAll(ctx context.Context, req *proto.GetShoppingItemsRequest) (*proto.GetShoppingItemsResponse, error) {
	handler := service.container.GetShoppingItemsHandler()
	return handler.Handle(ctx, req)
}

func newPersonalShoppingService(container *container) proto.PersonalShoppingServiceServer {
	return &personalShoppingService{container: container}
}

type accountService struct {
	container *container
	proto.UnimplementedAccountServiceServer
}

func newAccountService(container *container) proto.AccountServiceServer {
	return &accountService{container: container}
}

func (service *accountService) Get(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	handler := service.container.GetAccountHandler()
	return handler.Handle(ctx, req)
}

func (service *accountService) Create(ctx context.Context, req *proto.CreateAccountRequest) (*emptypb.Empty, error) {
	// TODO: Implements.
	return nil, nil
}
