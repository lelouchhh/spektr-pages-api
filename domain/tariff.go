package domain

import (
	"context"
)

type Tariff struct {
	Id               int          `json:"ID" db:"id"`
	Price            float64      `json:"price" db:"price"`
	PeriodPerPay     string       `json:"period_per_pay" db:"period_per_pay"`
	Title            string       `json:"title" db:"title"`
	Subtitle         string       `json:"subtitle" db:"subtitle"`
	ShortDescription string       `json:"short_description" db:"short_description"`
	Types            []TariffType `json:"tariff_type"`
	City             int          `json:"city_id"`
}
type TariffType struct {
	ID          int           `json:"ID" db:"ID"`
	Type        int           `json:"type" db:"Type"`
	TypeName    string        `json:"type_name" db:"typeName"`
	Name        string        `json:"name" db:"name"`
	Description []Description `json:"description" db:"description"`
	Title       string        `json:"title" db:"Title"`
	Subtitle    string        `json:"subtitle" db:"Subtitle"`
	Icon        int           `json:"icon" db:"icon"`
	IconPath    string        `json:"icon_path" db:"icon_path"`
}
type Description struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Icon struct {
	ID   int    `json:"ID" db:"id"`
	Path string `json:"path" db:"path"`
}
type Type struct {
	ID   int    `json:"ID" db:"id"`
	Name string `json:"name" db:"name"`
}

type TariffUsecase interface {
	GetTypes(ctx context.Context) ([]Type, error)
	GetTariffTypes(ctx context.Context) ([]TariffType, error)
	GetTariffs(ctx context.Context, id int) ([]Tariff, error)
	GetIcons(ctx context.Context) ([]Icon, error)

	AddTariffType(ctx context.Context, tType TariffType) error
	RemoveTariffType(ctx context.Context, id int) error
	RemoveTariff(ctx context.Context, Id int) error
	AddTariff(ctx context.Context, tariff Tariff) error
	AddIcon(ctx context.Context, icon Icon) error
	RemoveIcon(ctx context.Context, id int) error
}

type TariffRepository interface {
	GetTypes(ctx context.Context) ([]Type, error)
	GetTariffTypes(ctx context.Context) ([]TariffType, error)
	GetTariffs(ctx context.Context, id int) ([]Tariff, error)
	GetIcons(ctx context.Context) ([]Icon, error)

	AddTariffType(ctx context.Context, tType TariffType) error
	RemoveTariffType(ctx context.Context, id int) error
	RemoveTariff(ctx context.Context, Id int) error
	AddTariff(ctx context.Context, tariff Tariff) error
	AddIcon(ctx context.Context, icon Icon) error
	RemoveIcon(ctx context.Context, id int) (string, error)
}
