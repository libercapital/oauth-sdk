package oauthsdk

import "time"

type JWT struct {
	AccessToken *string
	Time        time.Time
}
