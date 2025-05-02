package registry

import (
	"context"
	"fmt"
	"share-basket-server/core/config"
	"share-basket-server/personal/domain"
	"share-basket-server/personal/infra/persistence"
	"share-basket-server/personal/presentation/handler"
	"share-basket-server/personal/presentation/router"
	"share-basket-server/personal/usecase"

	"gorm.io/gorm"
)

type (
	repositories struct {
		userRepo             domain.UserRepository
		userService          domain.UserService
		accountRepo          domain.AccountRepository
		shoppingCategoryRepo domain.ShoppingCategoryRepository

		transaction   domain.Transaction
		authenticator domain.Authenticator
	}

	interactors struct {
		signUpInteractor                usecase.SignUpInputPort
		signUpConfirmInteractor         usecase.SignUpConfirmInputPort
		loginInteractor                 usecase.LoginInputPort
		getShoppingCategoriesInteractor usecase.GetShoppingCategoriesInputPort
	}
)

func Inject(db *gorm.DB, cfg config.AWS) (router.Handlers, error) {
	ctx := context.Background()

	repos, err := injectRepository(ctx, db, cfg)
	if err != nil {
		return router.Handlers{}, err
	}

	interactors := injectInteractor(repos)

	return router.Handlers{
		PingHandler:                  handler.MakePingHandler(),
		SignUpHandler:                handler.MakeSignUpHandler(interactors.signUpInteractor),
		SignUpConfirmHandler:         handler.MakeSignUpConfirmHandler(interactors.signUpConfirmInteractor),
		LoginHandler:                 handler.MakeLoginHandler(interactors.loginInteractor),
		GetShoppingCaterogiesHandler: handler.MakeGetShoppingCategoriesHandler(interactors.getShoppingCategoriesInteractor),
	}, nil
}

func injectRepository(ctx context.Context, db *gorm.DB, cfg config.AWS) (repositories, error) {
	authenticator, err := persistence.NewCognito(ctx, cfg)
	if err != nil {
		return repositories{}, fmt.Errorf("failed to inject authenticator: %w", err)
	}

	userRepo := persistence.NewUserPersistence(db)

	return repositories{
		userRepo:             userRepo,
		userService:          domain.NewUserService(userRepo),
		accountRepo:          persistence.NewAccountPersistence(db),
		shoppingCategoryRepo: persistence.NewShoppingCategory(db),

		authenticator: authenticator,
		transaction:   persistence.NewTransaction(db),
	}, nil
}

func injectInteractor(repos repositories) interactors {
	signUpInteractor := usecase.NewSignUpInteractor(
		repos.authenticator,
		repos.userRepo,
		repos.accountRepo,
		repos.userService,
		repos.transaction,
	)

	return interactors{
		signUpInteractor:                signUpInteractor,
		signUpConfirmInteractor:         usecase.NewSignUpConfirmInteractor(repos.authenticator),
		loginInteractor:                 usecase.NewLoginInteractor(repos.authenticator),
		getShoppingCategoriesInteractor: usecase.NewGetShoppingCategoriesInteractor(repos.shoppingCategoryRepo),
	}
}
