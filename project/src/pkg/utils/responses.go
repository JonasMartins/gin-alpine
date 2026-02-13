package utils

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var UnrequiredAuthRoutes = []string{"login"}

type HTTPVerbs int

const (
	GET HTTPVerbs = iota
	POST
	PUT
	DELETE
	PATCH
)

func (v HTTPVerbs) String() string {
	return [...]string{"GET", "POST", "PUT", "DELETE", "PATCH"}[v]
}

func ParseHTTPVerbs(v string) (HTTPVerbs, bool) {
	switch strings.ToLower(v) {
	case "get":
		return GET, true
	case "post":
		return POST, true
	case "put":
		return PUT, true
	case "patch":
		return PATCH, true
	case "delete":
		return DELETE, true
	}
	return 0, false
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type GetBasicDatesIntervalParams struct {
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
}

type PaginatedLimitAndPage struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type GetPaginatedBasicParams struct {
	Page    uint32  `json:"page"`
	Limit   uint32  `json:"limit"`
	Enabled *bool   `json:"enabled"`
	Term    *string `json:"term" validate:"omitempty,min=3"`
}

func NewGetPaginatedBasicParams() *GetPaginatedBasicParams {
	return &GetPaginatedBasicParams{}
}

type PaginatedResult[T any] struct {
	CurrentPage       uint32 `json:"current_page"`
	TotalPages        uint32 `json:"total_pages"`
	TotalItems        uint32 `json:"total_items"`
	TotalItemsPerPage uint32 `json:"total_items_per_page"`
	Items             *[]*T  `json:"items"`
}

func NewPaginatedResult[T any]() *PaginatedResult[T] {
	return &PaginatedResult[T]{
		Items: &[]*T{},
	}
}

func CustomErrorTranslator(err error, customMessages *map[string]string) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			key := fieldErr.Field() + "." + fieldErr.Tag()
			if msg, exists := (*customMessages)[key]; exists {
				errors[fieldErr.Field()] = msg
			} else {
				errors[fieldErr.Field()] = fieldErr.Error() // Default error message
			}
		}
	}
	return errors
}
