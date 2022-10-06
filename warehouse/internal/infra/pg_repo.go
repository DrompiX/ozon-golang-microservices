package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/domain"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) CreateItem(ctx context.Context, i *domain.Item) error {
	const query = `
		INSERT INTO items (seller_id, price, quantity)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	return r.pool.QueryRow(ctx, query, i.SellerId, i.Price, i.Quantity).Scan(&i.ID)
}

func (r *postgresRepository) UpdateItem(ctx context.Context, i *domain.Item) error {
	const query = "UPDATE items SET quantity=$1, price=$2 WHERE id=$3"
	cmd, err := r.pool.Exec(ctx, query, i.Quantity, i.Price, i.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		err = domain.ErrNotFound
	}
	return err
}

func (r *postgresRepository) GetItem(ctx context.Context, itemId int64) (*domain.Item, error) {
	const query = `
		SELECT id, seller_id, price, quantity
		FROM items
		WHERE id = $1
	`
	var i domain.Item
	err := r.pool.QueryRow(ctx, query, itemId).Scan(&i.ID, &i.SellerId, &i.Price, &i.Quantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("item with id %d not found", itemId)
		}
		return nil, err
	}
	return &i, nil
}

// Small check that interface is implemented correctly
var _ domain.Repository = (*postgresRepository)(nil)
