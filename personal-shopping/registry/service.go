package registry

import (
	"context"
	proto "share-basket/personal-shopping/presentation/proto/gen"
)

type Services struct {
	PersonalShopping proto.PersonalShoppingServiceServer
}

func NewServices(container *container) Services {
	return Services{
		PersonalShopping: newPersonalShoppingService(container),
	}
}

type personalShoppingService struct {
	container *container
	proto.UnimplementedPersonalShoppingServiceServer
}

func (service *personalShoppingService) GetAll(ctx context.Context, req *proto.GetAllRequest) (*proto.GetAllResponse, error) {
	handler := service.container.GetShoppingItemsHandler()
	return handler.Handle(ctx, req)
}

func newPersonalShoppingService(container *container) proto.PersonalShoppingServiceServer {
	return &personalShoppingService{container: container}
}
