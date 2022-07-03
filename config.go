package db

// Config DB config
type Config struct {
	UserName     string `yaml:"userName,omitempty"`
	Password     string `yaml:"password,omitempty"`
	Host         string `yaml:"host,omitempty"`
	Port         int    `yaml:"port,omitempty"`
	DBName       string `yaml:"dbName,omitempty"`
	Charset      string `yaml:"charset,omitempty"`
	MaxIdleConns int    `yaml:"maxIdleConns,omitempty"`
	MaxOpenConns int    `yaml:"maxOpenConns,omitempty"`
	Debug        bool   `yaml:"debug,omitempty"`
}
