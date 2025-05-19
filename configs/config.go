package configs

type Config struct {
	Port            string `env:"APP_PORT" envDefault:"8080"`
	DBUser          string `env:"POSTGRES_USER"`
	DBPass          string `env:"POSTGRES_PASSWORD"`
	DBName          string `env:"POSTGRES_DB"`
	WeatherApiToken string `env:"WEATHER_API_TOKEN"`
	SmtpUser        string `env:"SMTP_USER"`
	SmtpPass        string `env:"SMTP_PASS"`
	SmtpHost        string `env:"SMTP_HOST"`
}
