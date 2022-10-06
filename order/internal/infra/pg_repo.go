package infra

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/domain"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) CreateOrder(ctx context.Context, o *domain.Order) error {
	const query = `
		INSERT INTO orders (item_id, user_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, o.ItemId, o.UserId, o.Status).Scan(&o.ID, &o.CreatedAt)
}

func (r *postgresRepository) UpdateOrder(ctx context.Context, o *domain.Order) error {
	const query = "UPDATE orders SET payment_id=$1, status=$2 WHERE id=$3"
	cmd, err := r.pool.Exec(ctx, query, o.PaymentId, o.Status, o.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		err = domain.ErrNotFound
	}
	return err
}

func (r *postgresRepository) GetOrder(ctx context.Context, orderId int64) (*domain.Order, error) {
	const query = `
		SELECT id, item_id, user_id, payment_id, status, created_at
		FROM orders
		WHERE id = $1
	`
	var o domain.Order
	err := r.pool.QueryRow(ctx, query, orderId).Scan(
		&o.ID, &o.ItemId, &o.UserId, &o.PaymentId, &o.Status, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// Small check that interface is implemented correctly
var _ domain.Repository = (*postgresRepository)(nil)
