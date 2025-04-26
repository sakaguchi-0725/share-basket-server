package presenter

import (
	"context"
	"net/http"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/output"
)

type signUpConfirmPresenter struct {
	w http.ResponseWriter
}

func (s *signUpConfirmPresenter) Render(ctx context.Context) error {
	response.NoContent(s.w)
	return nil
}

func NewSignUpConfirmPresenter(w http.ResponseWriter) output.SignUpConfirmOutputPort {
	return &signUpConfirmPresenter{w}
}
