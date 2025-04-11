package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	instana "github.com/instana/go-sensor"
	"go.uber.org/fx"

	"github.com/enesanbar/go-service/log"
	"go.uber.org/zap"
)

type Connection struct {
	Conn *sql.DB
}

type ConnectionParams struct {
	fx.In

	Sensor *instana.Sensor `optional:"true"`
	Config *Config
	Logger log.Factory
}

func NewConnection(p ConnectionParams) *Connection {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Europe%%2FIstanbul",
		p.Config.User, p.Config.Pass, p.Config.Host, p.Config.Port, p.Config.Name)
	var db *sql.DB
	var err error

	if p.Sensor != nil {
		instana.InstrumentSQLDriver(p.Sensor, "mysql", &mysql.MySQLDriver{})
		DSN = fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", p.Config.User, p.Config.Pass, p.Config.Host, p.Config.Name)
		if err != nil {
			p.Logger.Bg().Error("", zap.Error(err))
			return nil
		}
		location, _ := time.LoadLocation("Europe/Istanbul")
		connector, err := mysql.NewConnector(&mysql.Config{
			User:                 p.Config.User,
			Passwd:               p.Config.Pass,
			Net:                  "tcp",
			Addr:                 p.Config.Host,
			DBName:               p.Config.Name,
			AllowNativePasswords: true,
			Collation:            "utf8_general_ci",
			Loc:                  location,
			ParseTime:            true,
		})
		if err != nil {
			p.Logger.Bg().Error(err.Error(), zap.Error(err))
			return nil
		}

		db = sql.OpenDB(instana.WrapSQLConnector(p.Sensor, DSN, connector))
	} else {
		db, err = sql.Open("mysql", DSN)
	}

	if err != nil {
		p.Logger.Bg().Error("Unable to open/validate database connection", zap.Error(err))
		return nil
	}

	db.SetMaxIdleConns(p.Config.MaxIdleConnections)
	db.SetMaxOpenConns(p.Config.MaxOpenConnections)
	db.SetConnMaxLifetime(time.Duration(p.Config.MaxConnectionLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(p.Config.MaxConnectionIdletime) * time.Second)

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
	panic("implement me")
}

func (c *Connection) Stop() error {
	return c.Conn.Close()
}
