package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"golang-demo/user"
	"net/http"
	"strconv"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (handler *userHandler) Store(w http.ResponseWriter, r *http.Request) {
	var input user.InputUser
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	newUser, err := handler.userService.Store(input)
	render.JSON(w, r, Response{newUser})
}

func (handler *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	pageSize := 10
	page := 1
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}
	name := r.URL.Query().Get("name")
	country := r.URL.Query().Get("country")
	users, totalCount := handler.userService.Get(name, country, page, pageSize)
	render.JSON(w, r, Response{map[string]any{"users": users, "total_count": totalCount}})
}

func (handler *userHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userByID := r.Context().Value("user").(user.User)
	render.JSON(w, r, Response{userByID})
}

func (handler *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(user.User).ID
	var input user.InputUser
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = handler.userService.Update(userId, input)
	if err != nil {
		_ = render.Render(w, r, &ErrResponse{HTTPStatusCode: 400, StatusText: "Error during update", Err: err})
		return
	}
	render.JSON(w, r, Response{map[string]string{"message": "successfully updated"}})
}

func (handler *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(user.User).ID
	err := handler.userService.Delete(id)
	if err != nil {
		_ = render.Render(w, r, &ErrResponse{HTTPStatusCode: 400, StatusText: "Error during delete", Err: err})
		return
	}
	render.JSON(w, r, Response{map[string]string{"message": "successfully deleted"}})
}

type Response struct {
	Data interface{} `json:"data"`
}

type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func (handler *userHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userByID user.User
		var err error

		if userId := chi.URLParam(r, "userId"); userId != "" {
			userByID, err = handler.userService.GetByID(userId)
		} else {
			_ = render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			_ = render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", userByID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}