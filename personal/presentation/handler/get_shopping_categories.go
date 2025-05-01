package handler

import (
	"net/http"
	"share-basket-server/personal/presentation/presenter"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/input"
)

func MakeGetShoppingCategoriesHandler(usecase input.GetShoppingCategoriesPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errRw := w.(*response.ErrResponseWriter)
		errRw.Err = usecase.Execute(r.Context(), presenter.NewGetShoppingCategories(w))
	}
}
