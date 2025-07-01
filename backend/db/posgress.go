package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
)

type PostgresDB struct {
	DB *sql.DB
}

func zapLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment(
	//zap.AddStacktrace(zap.PanicLevel)
	)
	return logger
}

var Logger = zapLogger()

var logger = Logger

func NewPostgresSQLWrapper() Database {
	// Create the connection string

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", "localhost", 5432, "parking", "k8Fqs&N5Io", "parking", "disable") // Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Verify the connection is still alive
	err = db.Ping()
	if err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	return &PostgresDB{DB: db}
}

func (db *PostgresDB) SetData(ctx context.Context, query string, values ...interface{}) error {
	if len(values) == 0 {
		return errors.New("values cannot be empty")
	}

	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare query")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, values...)
	if err != nil {
		return errors.Wrapf(err, "failed to execute query")
	}

	return nil
}

func (db *PostgresDB) UpdateData(ctx context.Context, query string, values ...interface{}) error {
	if len(values) == 0 {
		return errors.Errorf("values cannot be empty")
	}

	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare query")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, values...)
	if err != nil {
		return errors.Wrapf(err, "failed to execute query")
	}

	return nil
}

func (db *PostgresDB) GetData(ctx context.Context, query string, values ...interface{}) (interface{}, error) {

	rows, err := db.DB.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute query")
	}
	return rows, nil
}

func (db *PostgresDB) DeleteData(ctx context.Context, query string, values ...interface{}) error {
	result, err := db.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrapf(err, "failed to execute delete query")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		// Not a hard error, but useful to know
		log.Printf("⚠️ DeleteData: no rows affected by query: %s, args: %v", query, values)
	}

	return nil
}
