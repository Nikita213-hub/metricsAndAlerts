package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	host     string
	user     string
	password string
	dbname   string
	sslMode  bool
	conn     *pgx.Conn
}

func NewDatabase(host string, user string, password string, dbname string, sslMode bool) *Database {
	return &Database{
		host:     host,
		user:     user,
		password: password,
		dbname:   dbname,
		sslMode:  sslMode,
	}
}

func (db *Database) Run() error {
	ps := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		db.user, db.password, db.host, db.dbname)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, ps)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	slog.Info("Database connected successfully", "db_name", db.dbname)
	return nil
}

func (db *Database) Ping() error {
	ctx := context.Background()
	return db.conn.Ping(ctx)
}
