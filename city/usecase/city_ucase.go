package usecase

import (
	"context"
	"spektr-pages-api/domain"
	"time"
)

type CityUsecase struct {
	cityRepo       domain.CityRepository
	contextTimeout time.Duration
}

func (c CityUsecase) GetCities(ctx context.Context) ([]domain.City, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	cities, err := c.cityRepo.GetCities(ctx)
	if err != nil {
		return []domain.City{}, err
	}
	return cities, nil
}

func (c CityUsecase) AddCity(ctx context.Context, city domain.City) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	// Perform any validation or business logic checks on the city data here

	err := c.cityRepo.AddCity(ctx, city)
	if err != nil {
		return err
	}
	return nil
}

func (c CityUsecase) RemoveCity(ctx context.Context, cityID int) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	err := c.cityRepo.RemoveCity(ctx, cityID)
	if err != nil {
		return err
	}
	return nil
}
func (c CityUsecase) RemoveCityTariff(ctx context.Context, cityTariff domain.CityTariff) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	err := c.cityRepo.RemoveCityTariff(ctx, cityTariff)
	if err != nil {
		return err
	}
	return nil
}
func NewCityUsecase(repo domain.CityRepository, timeout time.Duration) domain.CityUsecase {
	return &CityUsecase{
		cityRepo:       repo,
		contextTimeout: timeout,
	}
}
