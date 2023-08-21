package usecase

import (
	"context"
	"fmt"
	"os"
	"spektr-pages-api/domain"
	"strings"
	"time"
)

type TariffUsecase struct {
	tariffRepo     domain.TariffRepository
	contextTimeout time.Duration
}

func (t TariffUsecase) GetTypes(ctx context.Context) ([]domain.Type, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	types, err := t.tariffRepo.GetTypes(ctx)
	if err != nil {
		return []domain.Type{}, err
	}
	return types, nil
}

func (t TariffUsecase) GetTariffTypes(ctx context.Context) ([]domain.TariffType, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	types, err := t.tariffRepo.GetTariffTypes(ctx)
	if err != nil {
		return []domain.TariffType{}, err
	}
	return types, nil
}

func (t TariffUsecase) GetTariffs(ctx context.Context, id int) ([]domain.Tariff, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	types, err := t.tariffRepo.GetTariffs(ctx, id)
	if err != nil {
		return []domain.Tariff{}, err
	}
	return types, nil
}

func (t TariffUsecase) GetIcons(ctx context.Context) ([]domain.Icon, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	types, err := t.tariffRepo.GetIcons(ctx)
	if err != nil {
		return []domain.Icon{}, err
	}
	return types, nil
}

func (t TariffUsecase) AddTariffType(ctx context.Context, tType domain.TariffType) error {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	err := t.tariffRepo.AddTariffType(ctx, tType)
	if err != nil {
		return err
	}
	return nil
}

func (t TariffUsecase) RemoveTariffType(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	err := t.RemoveTariffType(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (t TariffUsecase) RemoveTariff(ctx context.Context, Id int) error {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	err := t.tariffRepo.RemoveTariff(ctx, Id)
	if err != nil {
		return err
	}
	return nil
}

func (t TariffUsecase) AddTariff(ctx context.Context, tariff domain.Tariff) error {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	err := t.tariffRepo.AddTariff(ctx, tariff)
	if err != nil {
		return err
	}
	return nil
}

func (t TariffUsecase) AddIcon(ctx context.Context, icon domain.Icon) error {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()
	err := t.tariffRepo.AddIcon(ctx, icon)
	if err != nil {
		return err
	}
	return nil
}

func (t TariffUsecase) RemoveIcon(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	path, err := t.tariffRepo.RemoveIcon(ctx, id)
	if err != nil {
		return err
	}
	splited := strings.Split(path, "/")
	err = os.Remove("./static/icons/" + splited[len(splited)-1])
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func NewTariffUsecase(a domain.TariffRepository, timeout time.Duration) domain.TariffUsecase {
	return &TariffUsecase{
		tariffRepo:     a,
		contextTimeout: timeout,
	}
}
