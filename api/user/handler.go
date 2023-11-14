package user

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
)

type userHandler struct {
	userService Service
}

func NewUserHandler(userService Service) *userHandler {
	return &userHandler{userService}
}

// make validation errors more user friendly
func validationErrorsToStr(errors validator.ValidationErrors) string {
	var errorMessages []string
	for _, vError := range errors {
		switch fieldName := vError.Field(); fieldName {
		case "FirstName":
			errorMessages = append(errorMessages, "first_name required")
		case "LastName":
			errorMessages = append(errorMessages, "last_name required")
		case "Nickname":
			errorMessages = append(errorMessages, "nickname required")
		case "Password":
			errorMessages = append(errorMessages, "password must be ascii only,min=8,max=72")
		case "Email":
			errorMessages = append(errorMessages, "email required")
		case "Country":
			errorMessages = append(errorMessages, "country must be two-letter country code uppercase")
		}
	}
	return strings.Join(errorMessages, ";")
}

func (handler *userHandler) Store(w http.ResponseWriter, r *http.Request) {
	var input InputUser
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	validate := validator.New()
	if err = validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		_ = render.Render(w, r, &ErrResponse{HTTPStatusCode: 400, StatusText: "validation errors", Err: err, ErrorText: validationErrorsToStr(validationErrors)})
		return
	}

	created, err := handler.userService.Store(input)
	if err != nil {
		_ = render.Render(w, r, &ErrResponse{HTTPStatusCode: 400, StatusText: "error during create", Err: err, ErrorText: err.Error()})
		return
	}
	render.JSON(w, r, Response{map[string]any{"message": "successfully created", "created": created}})
}

// Get paginated users by name (should be part of firstname or lastname) and country (must be equal to)
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
	users, totalCount, err := handler.userService.Get(name, country, page, pageSize)
	if err != nil {
		_ = render.Render(w, r, &ErrResponse{HTTPStatusCode: 400, StatusText: "error during select", Err: err, ErrorText: err.Error()})
		return
	}
	render.JSON(w, r, Response{map[string]any{"users": users, "total_count": totalCount}})
}

func (handler *userHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userByID := r.Context().Value("user").(User)
	render.JSON(w, r, Response{userByID})
}

func (handler *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(User).ID
	var input InputUser
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	validate := validator.New()
	if err = validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		_ = render.Render(w, r, &ErrResponse{HTTPStatusCode: 400, StatusText: "validation errors", Err: err, ErrorText: validationErrorsToStr(validationErrors)})
		return
	}

	err = handler.userService.Update(userId, input)
	if err != nil {
		_ = render.Render(w, r, &ErrResponse{HTTPStatusCode: 400, StatusText: "error during update", Err: err, ErrorText: err.Error()})
		return
	}
	render.JSON(w, r, Response{map[string]string{"message": "successfully updated"}})
}

func (handler *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(User).ID
	err := handler.userService.Delete(id)
	if err != nil {
		_ = render.Render(w, r, &ErrResponse{HTTPStatusCode: 400, StatusText: "error during delete", Err: err, ErrorText: err.Error()})
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

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "invalid request",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "resource not found"}

// UserCtx retrieves user by id and stores it in request context
func (handler *userHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIdStr := chi.URLParam(r, "userId")
		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		userById, err := handler.userService.GetById(userId)
		if err != nil {
			_ = render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), "user", userById)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
