package pvz

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/golang-avito/internal/feature/acceptance"
	"github.com/Sanchir01/golang-avito/internal/feature/product"
	"github.com/Sanchir01/golang-avito/pkg/lib/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: primaryDB}
}

func (r *Repository) GetAllPVZ(ctx context.Context, startDate, endDate time.Time, page, limit uint64) ([]*DBPVZWithReceptions, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `
    WITH product_counts AS (
        SELECT receiving_id, COUNT(*) as product_count
        FROM product
        GROUP BY receiving_id
    )
    SELECT
        p.id AS pvz_id,
        p.registration_date,
        p.city,
        a.id AS acceptance_id,
        a.created_at,
        a.pvz_id,
        a.status,
        pr.id AS product_id,
        pr.created_at AS product_created_at,
        pr.type,
        pr.receiving_id
    FROM
        pvz p
    LEFT JOIN
        acceptance a ON p.id = a.pvz_id
    LEFT JOIN
        product pr ON a.id = pr.receiving_id
    LEFT JOIN
        product_counts pc ON a.id = pc.receiving_id
    ORDER BY p.registration_date DESC
    `

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pvzMap := make(map[uuid.UUID]*DBPVZWithReceptions)
	receptionMap := make(map[uuid.UUID]map[uuid.UUID]*DBReceptionWithProducts)

	for rows.Next() {
		var (
			pvzID              uuid.UUID
			pvzReg             time.Time
			city               string
			receptionID        sql.NullString
			receptionDateTime  sql.NullTime
			receptionPVZID     sql.NullString
			status             sql.NullString
			productID          sql.NullString
			productDateTime    sql.NullTime
			productType        sql.NullString
			productReceptionID sql.NullString
		)
		err := rows.Scan(
			&pvzID, &pvzReg, &city,
			&receptionID, &receptionDateTime, &receptionPVZID, &status,
			&productID, &productDateTime, &productType, &productReceptionID,
		)
		if err != nil {
			return nil, err
		}

		pvzWithReceptions, exists := pvzMap[pvzID]
		if !exists {
			pvzWithReceptions = &DBPVZWithReceptions{
				PVZ: DBPVZ{
					ID:               pvzID,
					RegistrationDate: pvzReg,
					City:             city,
				},
			}
			pvzMap[pvzID] = pvzWithReceptions
			receptionMap[pvzID] = make(map[uuid.UUID]*DBReceptionWithProducts)
		}

		if receptionID.Valid {
			receptionsInPVZ := receptionMap[pvzID]
			receptionUUID, err := uuid.Parse(receptionID.String)
			if err != nil {
				return nil, err
			}

			reception, exists := receptionsInPVZ[receptionUUID]
			if !exists {
				pvzUUID, err := uuid.Parse(receptionPVZID.String)
				if err != nil {
					return nil, err
				}

				reception = &DBReceptionWithProducts{
					Reception: acceptance.DBAcceptance{
						ID:        receptionUUID,
						CreatedAt: receptionDateTime.Time,
						PvzId:     pvzUUID,
						Status:    status.String,
					},
				}
				receptionsInPVZ[receptionUUID] = reception
				pvzWithReceptions.Receptions = append(pvzWithReceptions.Receptions, *reception)
			}

			if productID.Valid {
				prodUUID, err := uuid.Parse(productID.String)
				if err != nil {
					return nil, err
				}

				product := product.DBProduct{
					ID:          prodUUID,
					CreatedAt:   productDateTime.Time,
					Type:        productType.String,
					ReceptionID: receptionUUID,
				}

				reception.Products = append(reception.Products, product)

				for i, r := range pvzWithReceptions.Receptions {
					if r.Reception.ID == receptionUUID {
						pvzWithReceptions.Receptions[i].Products = reception.Products
						break
					}
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]*DBPVZWithReceptions, 0, len(pvzMap))
	for _, pvz := range pvzMap {
		if len(pvz.Receptions) == 0 {
			continue
		}

		filteredReceptions := make([]DBReceptionWithProducts, 0, len(pvz.Receptions))
		for _, reception := range pvz.Receptions {
			if len(reception.Products) > 0 {
				filteredReceptions = append(filteredReceptions, reception)
			}
		}

		if len(filteredReceptions) > 0 {
			pvz.Receptions = filteredReceptions
			result = append(result, pvz)
		}
	}

	return result, nil
}
func (r *Repository) CreatePVZ(ctx context.Context, registerDate time.Time, city string, tx pgx.Tx) (*DBPVZ, error) {
	query, arg, err := sq.Insert("pvz").
		Columns("registration_date", "city").
		Values(registerDate, city).
		Suffix("RETURNING id, registration_date, city").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var pvz DBPVZ
	if err := tx.QueryRow(ctx, query, arg...).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
		if err == pgx.ErrTxCommitRollback {
			return nil, api.ErrCreatePvz
		}
		return nil, err
	}
	if err != nil {
		return nil, api.ErrorCreateQueryString
	}

	return &DBPVZ{
		ID:               pvz.ID,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}, nil
}
