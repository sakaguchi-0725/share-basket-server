package presenter

import (
	"context"
	"net/http"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/output"
)

type signUpPresenter struct {
	w http.ResponseWriter
}

func (presenter *signUpPresenter) Render(ctx context.Context) error {
	response.NoContent(presenter.w)
	return nil
}

func NewSignUpPresenter(w http.ResponseWriter) output.SignUpOutputPort {
	return &signUpPresenter{w}
}
