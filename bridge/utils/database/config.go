package database

// Config db config
type Config struct {
	// data source name
	DSN        string `json:"dsn"`
	DriverName string `json:"driverName"`

	MaxOpenNum int `json:"maxOpenNum"`
	MaxIdleNum int `json:"maxIdleNum"`
}
