// Package postgres ...
package postgres

import (
	"context"
	"embed"
	"log"

	"gin-alpine/src/internal/configs"
	"gin-alpine/src/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type PgRepository struct {
	DB     *pgxpool.Pool
	Config *configs.Config
}

func NewPgRepository(c *configs.Config) (*PgRepository, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(c.DBConn)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 100
	config.MinConns = 2

	// conn, err := pgx.Connect(ctx, *b.Config.DB.Conn)
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	connConfig, err := pgx.ParseConfig(c.DBConn)
	if err != nil {
		return nil, err
	}

	bridge := stdlib.OpenDB(*connConfig)
	defer func() {
		deferErr := bridge.Close()
		if deferErr != nil {
			err = deferErr
		}
	}()

	log.Println("database successfully connected")
	goose.SetBaseFS(embedMigrations)
	err = goose.Up(bridge, "migrations")
	if err != nil {
		return nil, err
	}
	return &PgRepository{DB: pool, Config: c}, nil
}
func (r *PgRepository) GetConnection(ctx context.Context) error {
	err := r.DB.Ping(ctx)
	if err != nil {
		r.DB, err = r.OpenDBConnection()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PgRepository) OpenDBConnection() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(r.Config.DBConn)
	if err != nil {
		utils.FatalResult("error connecting database", err)
	}
	config.MaxConns = 5
	config.MaxConns = 2
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (r *PgRepository) CloseConnection() {
	r.DB.Close()
}

func (r *PgRepository) DeferTx(err error, tx pgx.Tx, ctx context.Context) {
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Printf("error at db execution, rollback %v", err)
		}
	} else {
		err = tx.Commit(ctx)
		if err != nil {
			log.Printf("error at db execution, commit %v", err)
		}
	}
}
