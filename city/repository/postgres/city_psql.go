package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"spektr-pages-api/domain"
)

type psqlCityRepository struct {
	db *sqlx.DB
}

func NewCityRepository(conn *sqlx.DB) domain.CityRepository {
	return &psqlCityRepository{conn}
}

func (p *psqlCityRepository) GetCities(ctx context.Context) ([]domain.City, error) {
	var cities []domain.City
	err := p.db.SelectContext(ctx, &cities, "SELECT id, name FROM spektr.t_city")
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	return cities, nil
}

func (p *psqlCityRepository) AddCity(ctx context.Context, city domain.City) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO spektr.t_city (name) VALUES ($1)", city.Name)
	if err != nil {
		return err
	}
	return nil
}

func (p *psqlCityRepository) RemoveCity(ctx context.Context, cityID int) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM spektr.t_city WHERE id = $1", cityID)
	if err != nil {
		return err
	}
	return nil
}
func (p *psqlCityRepository) RemoveCityTariff(ctx context.Context, cityTariff domain.CityTariff) error {
	query := `DELETE FROM t_city_tariff WHERE city_id = $1 AND tariff_id = $2`
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return domain.ErrInternalServerError
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cityTariff.City, cityTariff.Tariff)
	if err != nil {
		return domain.ErrInternalServerError
	}

	return nil
}
