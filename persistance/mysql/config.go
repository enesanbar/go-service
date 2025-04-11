package mysql

type Config struct {
	Host                  string
	Port                  int
	User                  string
	Pass                  string
	Name                  string
	Timeout               int
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifetime int
	MaxConnectionIdletime int
}
