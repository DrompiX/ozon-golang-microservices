package infra

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/DrompiX/homework-3/billing/internal/domain"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) CreatePayment(ctx context.Context, p *domain.Payment) error {
	const query = `
		INSERT INTO payments (order_id, from_user, to_user, amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.pool.QueryRow(ctx, query, p.OrderId, p.FromUserId, p.ToUserId, p.Amount).Scan(&p.ID)
}

// Small check that interface is implemented correctly
var _ domain.Repository = (*postgresRepository)(nil)
