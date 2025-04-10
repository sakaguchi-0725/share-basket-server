package usecase

type GetAccountUseCase interface {
	Execute(id string) (GetAccountOutput, error)
}

type getAccountUseCase struct{}

func (g *getAccountUseCase) Execute(id string) (GetAccountOutput, error) {
	return GetAccountOutput{
		ID:   "dummy_user_id",
		Name: "dummy user",
	}, nil
}

type GetAccountOutput struct {
	ID   string
	Name string
	// TODO: 家族情報も含める
}

func NewGetAccountUseCase() GetAccountUseCase {
	return &getAccountUseCase{}
}
