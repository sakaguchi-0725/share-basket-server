package registry

import (
	"context"
	"fmt"
	"share-basket-server/core/config"
	"share-basket-server/personal/domain/service"
	"share-basket-server/personal/infra/persistence"
	"share-basket-server/personal/presentation/handler"
	"share-basket-server/personal/presentation/router"
	"share-basket-server/personal/usecase"

	"gorm.io/gorm"
)

func Inject(db *gorm.DB, cfg config.AWS) (router.Handlers, error) {
	ctx := context.Background()

	authenticator, err := persistence.NewCognito(ctx, cfg)
	if err != nil {
		return router.Handlers{}, fmt.Errorf("failed to inject: %w", err)
	}
	userRepo := persistence.NewUserPersistence(db)
	accountRepo := persistence.NewAccountPersistence(db)
	transaction := persistence.NewTransaction(db)
	userService := service.NewUserService(userRepo)

	signUpUseCase := usecase.NewSignUpUseCase(
		authenticator,
		userRepo,
		accountRepo,
		userService,
		transaction,
	)

	signUpConfirmUseCase := usecase.NewSignUpConfirmUseCase(authenticator)

	return router.Handlers{
		PingHandler:          handler.MakePingHandler(),
		SignUpHandler:        handler.MakeSignUpHandler(signUpUseCase),
		SignUpConfirmHandler: handler.MakeSignUpConfirmHandler(signUpConfirmUseCase),
	}, nil
}
