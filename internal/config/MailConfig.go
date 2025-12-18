package config

type MailConfig struct {
	SMTPHost string  `yaml:"smtp_host"`
	SMTPPort int     `yaml:"smtp_port"`
	Username string  `yaml:"username"`
	AuthCode string  `yaml:"auth_code"`
	From     string  `yaml:"from"`
	Timeout  float64 `yaml:"timeout"`
	Retry    uint8   `yaml:"retry"`
}

func DefaultMailConfig() *MailConfig {
	return &MailConfig{
		SMTPHost: "",
		SMTPPort: 587,
		Username: "",
		AuthCode: "",
		From:     "",
		Timeout:  5.0,
		Retry:    3,
	}
}
