package config

type I18n struct {
	Locales []string `yaml:"locales"`
	Default string   `yaml:"default"`
	Dir     string   `yaml:"dir"`
}

type Database struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	User    string `yaml:"user"`
	Pass    string `yaml:"pass"`
	Name    string `yaml:"name"`
	Debug   bool   `yaml:"debug"`
	SSLMode string `yaml:"ssl_mode"`
	Migrate bool   `yaml:"migrate"`
}

type Rest struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Config struct {
	DB   Database `yaml:"db"`
	I18n I18n     `yaml:"i18n"`
	Rest Rest     `yaml:"rest"`
}
