package config

type Config struct {
	FilePath string
	URL      string
	Method   string
	Body     string
	Delay    int
}

func New() *Config {
	return &Config{
		Method: "POST",
		Delay:  0,
	}
}
