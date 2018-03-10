package twitter

import (
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("marija/datasources/twitter")

type Config struct {
	ConsumerKey    string
	ConsumerSecret string

	Token       string
	TokenSecret string
}
