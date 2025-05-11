package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"

	"github.com/XSAM/otelsql"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/zap"
)

type Connection struct {
	Conn *sql.DB
}

type ConnectionParams struct {
	fx.In

	Config *Config
	Logger log.Factory
}

func NewConnection(p ConnectionParams) *Connection {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Europe%%2FIstanbul",
		p.Config.User, p.Config.Pass, p.Config.Host, p.Config.Port, p.Config.Database)
	var db *sql.DB
	var err error

	// db, err = sql.Open("mysql", DSN)
	attrs := append(otelsql.AttributesFromDSN(DSN), semconv.DBSystemMySQL)
	db, err = otelsql.Open(
		"mysql",
		DSN,
		otelsql.WithAttributes(attrs...),
	)
	if err != nil {
		p.Logger.Bg().With(zap.Error(err)).Error("Unable to open database connection")
		return nil
	}

	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		attrs...,
	))
	if err != nil {
		p.Logger.Bg().With(zap.Error(err)).Error("Unable to register database stats metrics")
		return nil
	}

	db.SetMaxIdleConns(p.Config.MaxIdleConnections)
	db.SetMaxOpenConns(p.Config.MaxOpenConnections)
	db.SetConnMaxLifetime(time.Duration(p.Config.MaxConnectionLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(p.Config.MaxConnectionIdleTime) * time.Second)

	p.Logger.Bg().With(
		zap.String("host", p.Config.Host),
		zap.Int("port", p.Config.Port),
		zap.String("database", p.Config.Name),
	).Info("Connecting to MySQL")

	duration := time.Duration(p.Config.Timeout) * time.Second
	var ctx, cancel = context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Open doesn't open a connection. Validate DSN data:
	err = db.PingContext(ctx)
	if err != nil {
		p.Logger.Bg().Error("Could not ping database", zap.Error(err))
		return nil
	}

	return &Connection{
		Conn: db,
	}
}

func (c *Connection) GetConn() *sql.DB {
	return c.Conn
}

func (c *Connection) Start() error {
	return nil
}

func (c *Connection) Stop() error {
	return c.Conn.Close()
}
