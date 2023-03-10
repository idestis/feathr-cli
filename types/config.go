package types

type Config struct {
	GenOnCreate bool       `survey:"generate_on_create" yaml:"generate_on_create"`
	GenOnUpdate bool       `survey:"generate_on_update" yaml:"generate_on_update"`
	SMTP        SMTPConfig `survey:"smtp" yaml:"smtp"`
}
