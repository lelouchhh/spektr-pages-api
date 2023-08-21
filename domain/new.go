package domain

import "context"

type New struct {
	Id int `json:"id,omitempty"`

	Title string `json:"title,omitempty" validate:"required"`
	Body  string `json:"body,omitempty" validate:"required"`

	Date string `json:"date,omitempty" validate:"required,date=2006-01-02"`
	//files
	Image    string `json:"image,omitempty" validate:"required"`
	Document string `json:"document,omitempty" validate:"required"`
}

type NewUsecase interface {
	GetNews(ctx context.Context) ([]New, error)
	RemoveNew(ctx context.Context, Id int) error
	AddNew(ctx context.Context, new New) error
}

type NewRepository interface {
	GetNews(ctx context.Context) ([]New, error)
	RemoveNew(ctx context.Context, Id int) error
	AddNew(ctx context.Context, new New) error
}
