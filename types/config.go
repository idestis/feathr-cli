package types

type Config struct {
	Storage     string     `survey:"storage" yaml:"storage"`
	GenOnCreate bool       `survey:"generate_on_create" yaml:"generate_on_create"`
	SMTP        SMTPConfig `survey:"smtp" yaml:"smtp"`
}
