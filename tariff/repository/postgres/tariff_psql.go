package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
	"spektr-pages-api/domain"
)

type psqlTariffRepository struct {
	db *sqlx.DB
}

func NewTariffRepository(conn *sqlx.DB) domain.TariffRepository {
	return &psqlTariffRepository{conn}
}

func jsonToString(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(b)
}

func (p *psqlTariffRepository) GetTypes(ctx context.Context) ([]domain.Type, error) {
	var types []domain.Type
	err := p.db.SelectContext(ctx, &types, "SELECT id, name FROM spektr.t_type")
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	return types, nil
}

func (p *psqlTariffRepository) GetTariffTypes(ctx context.Context) ([]domain.TariffType, error) {
	var types []domain.TariffType
	query := `
		SELECT
			tt.id,
			tt.name,
			tt.description,
			tt.title,
			tt.subtitle,
			tt.icon,
			tt.type,
			ot.name AS typeName,
			ic.path
		FROM
			spektr.t_tariff_type tt
		JOIN
			t_type ot ON tt.type = ot.id
		JOIN
			t_icon ic ON ic.id = tt.icon`
	rows, err := p.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var tt domain.TariffType
		var descriptionJSON []byte
		if err := rows.Scan(
			&tt.ID, &tt.Name, &descriptionJSON, &tt.Title, &tt.Subtitle, &tt.Icon,
			&tt.Type, &tt.TypeName, &tt.IconPath,
		); err != nil {
			return nil, domain.ErrInternalServerError
		}
		var descriptions []domain.Description
		if err := json.Unmarshal(descriptionJSON, &descriptions); err != nil {
			return nil, domain.ErrInternalServerError
		}
		tt.Description = descriptions
		types = append(types, tt)
	}
	if err := rows.Err(); err != nil {
		return nil, domain.ErrInternalServerError
	}
	return types, nil
}

func (p *psqlTariffRepository) getTariffTypesById(ctx context.Context, id int) ([]domain.TariffType, error) {
	var types []domain.TariffType
	query := `
		SELECT
			tt.id,
			tt.name,
			tt.description,
			tt.title,
			tt.subtitle,
			tt.icon,
			tt.type,
			ot.name AS typeName,
			ic.path
		FROM
			spektr.t_tariff_type tt
		JOIN
			t_type ot ON tt.type = ot.id
		JOIN
			t_icon ic ON ic.id = tt.icon
		JOIN
			spektr.t_tariff_type_tariff tttt ON tt.id = tttt.tariff_type_id
		WHERE
			tttt.tariff_id = $1`
	rows, err := p.db.QueryxContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tt domain.TariffType
		var descriptionJSON []byte
		if err := rows.Scan(
			&tt.ID, &tt.Name, &descriptionJSON, &tt.Title, &tt.Subtitle, &tt.Icon,
			&tt.Type, &tt.TypeName, &tt.IconPath,
		); err != nil {
			return nil, err
		}
		var descriptions []domain.Description
		if err := json.Unmarshal(descriptionJSON, &descriptions); err != nil {
			return nil, err
		}
		tt.Description = descriptions
		types = append(types, tt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return types, nil
}

func (p *psqlTariffRepository) GetTariffs(ctx context.Context, id int) ([]domain.Tariff, error) {
	var tariffs []domain.Tariff
	err := p.db.SelectContext(ctx, &tariffs, "SELECT id, price, period_per_pay, title, subtitle, short_description FROM spektr.t_tariff join t_city_tariff on t_tariff.id = t_city_tariff.tariff_id where city_id = $1;", id)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	var eg errgroup.Group
	for i := range tariffs {
		v := &tariffs[i]
		eg.Go(func() error {
			tariffTypes, err := p.getTariffTypesById(ctx, v.Id)
			if err != nil {
				return err
			}
			v.Types = tariffTypes
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, domain.ErrInternalServerError
	}
	fmt.Println()
	return tariffs, nil
}

func (p *psqlTariffRepository) GetIcons(ctx context.Context) ([]domain.Icon, error) {
	var icons []domain.Icon
	err := p.db.SelectContext(ctx, &icons, "SELECT id, path FROM spektr.t_icon")
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	return icons, nil
}

func (p *psqlTariffRepository) AddTariffType(ctx context.Context, tType domain.TariffType) error {
	query := `INSERT INTO t_tariff_type (name, description, title, subtitle, icon, type) VALUES ($1, $2, $3, $4, $5, $6);`
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return domain.ErrInternalServerError
	}
	defer stmt.Close() // Close the statement after use
	_, err = stmt.ExecContext(ctx, tType.Name, jsonToString(tType.Description), tType.Title, tType.Subtitle, tType.Icon, tType.Type)
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (p *psqlTariffRepository) RemoveTariffType(ctx context.Context, id int) error {
	query := `DELETE FROM spektr.t_tariff_type WHERE id = $1`
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return domain.ErrInternalServerError
	}
	defer stmt.Close() // Close the statement after use
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (p *psqlTariffRepository) RemoveTariff(ctx context.Context, id int) error {
	query := `DELETE FROM t_tariff WHERE id = $1`
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return domain.ErrInternalServerError
	}
	defer stmt.Close() // Close the statement after use
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (p *psqlTariffRepository) AddTariff(ctx context.Context, tariff domain.Tariff) error {
	query := `INSERT INTO spektr.t_tariff (price, period_per_pay, title, subtitle, short_description) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return domain.ErrInternalServerError
	}
	defer stmt.Close() // Close the statement after use
	var insertedID int
	err = stmt.QueryRowContext(ctx, tariff.Price, tariff.PeriodPerPay, tariff.Title, tariff.Subtitle, tariff.ShortDescription).Scan(&insertedID)
	if err != nil {
		return domain.ErrInternalServerError
	}
	query = `INSERT INTO spektr.t_tariff_type_tariff (tariff_id, tariff_type_id) VALUES ($1, $2)`
	var g errgroup.Group
	for _, v := range tariff.Types {
		v := v
		g.Go(func() error {
			stmt, err := p.db.PrepareContext(ctx, query)
			if err != nil {
				return domain.ErrInternalServerError
			}
			defer stmt.Close() // Close the statement after use
			_, err = stmt.ExecContext(ctx, insertedID, v.ID)
			if err != nil {
				return domain.ErrInternalServerError
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	cityTariffQuery := `INSERT INTO spektr.t_city_tariff (city_id, tariff_id) VALUES ($1, $2)`
	stmt, err = p.db.PrepareContext(ctx, cityTariffQuery)
	if err != nil {
		return domain.ErrInternalServerError
	}
	defer stmt.Close() // Close the statement after use
	_, err = stmt.ExecContext(ctx, tariff.City, insertedID)
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (p *psqlTariffRepository) AddIcon(ctx context.Context, icon domain.Icon) error {
	query := `INSERT INTO spektr.t_icon (path) VALUES ($1)`
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return domain.ErrInternalServerError
	}
	defer stmt.Close() // Close the statement after use
	_, err = stmt.ExecContext(ctx, icon.Path)
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (p *psqlTariffRepository) RemoveIcon(ctx context.Context, id int) (string, error) {
	var path string
	err := p.db.GetContext(ctx, &path, "SELECT path FROM spektr.t_icon WHERE id =$1", id)
	if err != nil {
		return "", domain.ErrInternalServerError
	}
	query := `DELETE FROM t_icon WHERE id = $1`
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return "", domain.ErrInternalServerError
	}
	defer stmt.Close() // Close the statement after use
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return "", domain.ErrInternalServerError
	}
	return path, nil
}
