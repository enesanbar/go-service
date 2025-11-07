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
	Logger log.Factory
	Conn   *sql.DB
	Config *Config
}

type ConnectionParams struct {
	fx.In

	Logger log.Factory
	Config *Config
}

func NewConnection(p ConnectionParams) *Connection {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
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

	return &Connection{
		Logger: p.Logger,
		Conn:   db,
		Config: p.Config,
	}
}

func (c *Connection) GetConn() *sql.DB {
	return c.Conn
}

func (c *Connection) Start(ctx context.Context) error {
	c.Logger.Bg().With(
		zap.String("host", c.Config.Host),
		zap.Int("port", c.Config.Port),
		zap.String("database", c.Config.Name),
	).Info("Connecting to MySQL")

	duration := time.Duration(c.Config.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	err := c.Conn.PingContext(ctx)
	if err != nil {
		c.Logger.Bg().Error("Could not ping database", zap.Error(err))
		return nil
	}
	return nil
}

func (c *Connection) Close(ctx context.Context) error {
	c.Logger.For(ctx).Info("closing MySQL connection")
	return c.Conn.Close()
}

func (c *Connection) Name() string {
	return c.Config.Name
}
