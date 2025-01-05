package models

type TestRequest struct {
	Endpoint    string            `bson:"endpoint"`
	Method      string            `bson:"method"`
	Body        string            `bson:"body"`
	Concurrency int               `bson:"concurrency"`
	Requests    int               `bson:"requests"`
	Headers     map[string]string `bson:"headers"`
}
