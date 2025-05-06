package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

type testServer struct {
	method string
	path   string
	router *chi.Mux
}

func newTestServer(method, path string) *testServer {
	return &testServer{
		method: method,
		path:   path,
		router: chi.NewRouter(),
	}
}

func (ts *testServer) mapHandler(handler http.HandlerFunc) {
	switch ts.method {
	case http.MethodGet:
		ts.router.Get(ts.path, handler)
	case http.MethodPost:
		ts.router.Post(ts.path, handler)
	case http.MethodPut:
		ts.router.Put(ts.path, handler)
	case http.MethodDelete:
		ts.router.Delete(ts.path, handler)
	}
}

func (ts *testServer) Serve(body io.Reader, handler http.HandlerFunc) *httptest.ResponseRecorder {
	ts.mapHandler(handler)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(ts.method, ts.path, body)
	req.Header.Set("Content-Type", "application/json")

	ts.router.ServeHTTP(rec, req)

	return rec
}
