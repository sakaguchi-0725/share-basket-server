package registry

import (
	"context"
	"fmt"
	"log"
	"share-basket-server/core/config"
	"share-basket-server/domain"
	"share-basket-server/infra/aws"
	"share-basket-server/infra/rdb/db"
	"share-basket-server/infra/rdb/repository"
	"share-basket-server/presentation/handler"
	"share-basket-server/presentation/server"
	"share-basket-server/presentation/validator"
	"share-basket-server/usecase"

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

func Inject(cfg config.App) (server.Handlers, error) {
	ctx := context.Background()

	db, err := db.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	repos, err := injectRepository(ctx, db, cfg.AWS)
	if err != nil {
		return server.Handlers{}, err
	}

	interactors := injectInteractor(repos)

	validator := validator.New()

	return server.Handlers{
		PingHandler:                  handler.MakePingHandler(),
		SignUpHandler:                handler.MakeSignUpHandler(interactors.signUpInteractor, validator),
		SignUpConfirmHandler:         handler.MakeSignUpConfirmHandler(interactors.signUpConfirmInteractor, validator),
		LoginHandler:                 handler.MakeLoginHandler(interactors.loginInteractor, validator),
		GetShoppingCaterogiesHandler: handler.MakeGetShoppingCategoriesHandler(interactors.getShoppingCategoriesInteractor),
	}, nil
}

func injectRepository(ctx context.Context, db *gorm.DB, cfg config.AWS) (repositories, error) {
	cognitoClient, err := aws.NewCognitoClient(ctx, cfg)
	if err != nil {
		return repositories{}, fmt.Errorf("failed to inject authenticator: %w", err)
	}

	authenticator := aws.NewCognitoPersistence(cognitoClient)
	userRepo := repository.NewUserPersistence(db)

	return repositories{
		userRepo:             userRepo,
		userService:          domain.NewUserService(userRepo),
		accountRepo:          repository.NewAccountPersistence(db),
		shoppingCategoryRepo: repository.NewShoppingCategoryPersistence(db),

		authenticator: authenticator,
		transaction:   repository.NewTransaction(db),
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
