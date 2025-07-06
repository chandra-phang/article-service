package configloader

type RootConfig struct {
	AppConfig     `mapstructure:"app"`
	DbConfig      `mapstructure:"db"`
	ElasticConfig `mapstructure:"elastic"`
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type DbConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	DbName     string `mapstructure:"dbname"`
	TestDbName string `mapstructure:"test_dbname"`
}

type ElasticConfig struct {
	URL string `json:"url"`
}
