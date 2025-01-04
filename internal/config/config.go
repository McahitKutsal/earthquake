package config

type Config struct {
	Endpoint    string            `json:"endpoint"`
	Method      string            `json:"method"`
	Body        string            `json:"body"`
	Concurrency int               `json:"concurrency"`
	Requests    int               `json:"requests"`
	Headers     map[string]string `json:"headers"`
}
