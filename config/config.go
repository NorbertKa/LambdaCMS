package config

type Config struct {
	Port        uint16   `json:"port"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Postgre     ConfigDB `json:"postgre"`
	Redis       ConfigDB `json:"redis"`
}

type ConfigDB struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	DB       int    `json:"db"`
	Username string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c Config) Port_Int() int {
	return int(c.Port)
}

func (c ConfigDB) Port_Int() int {
	return int(c.Port)
}
