package oauthsdk

const ContentTypeApplicationJson ContentType = "application/json; charset=UTF-8"
const ContentTypeFormURLEnconded ContentType = "application/x-www-form-urlencoded"

type ContentType string

func (c ContentType) String() string {
	return string(c)
}

type CertData struct {
	ClientCrtKey string
	ClientCrt    string
}

type Config struct {
	URL          string
	ClientID     string
	ClientSecret string
	GrantType    string
	Audience     string
	ContentType  ContentType
	//Infos about mtls
	CertData *CertData
	//Redact keys to add into bavalogs defaults
	RedactKeys []string
	//time in seconds to antecipate expiration of token
	ExpirationMarginSeconds int
}
