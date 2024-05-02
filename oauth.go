package oauthsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/bavatech/architecture/software/libs/go-modules/bavalogs.git"
)

type requestClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint16 `json:"expires_in"`
}

func (c client) getData() io.Reader {
	if c.Config.ContentType == ContentTypeFormURLEnconded {
		data := url.Values{}

		if c.Config.ClientID != "" {
			data.Set("client_id", c.Config.ClientID)
		}

		if c.Config.ClientSecret != "" {
			data.Set("client_secret", c.Config.ClientSecret)
		}

		if c.Config.GrantType != "" {
			data.Set("grant_type", c.Config.GrantType)
		}

		if c.Config.Audience != "" {
			data.Set("audience", c.Config.Audience)
		}

		return bytes.NewBuffer([]byte(data.Encode()))
	}

	if c.Config.ContentType == ContentTypeApplicationJson {
		body := map[string]interface{}{
			"client_id":     c.Config.ClientID,
			"client_secret": c.Config.ClientSecret,
			"grant_type":    c.Config.GrantType,
			"audience":      c.Config.Audience,
		}

		bsBody, _ := json.Marshal(body)
		return bytes.NewReader(bsBody)
	}

	return nil
}

func (c *client) createAccessToken(ctx context.Context) (*JWT, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Config.URL, c.getData())
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", c.Config.ContentType.String())

	response, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request error (status=%v)", response.StatusCode)
	}

	var oauthResponse *OAuthResponse

	err = json.NewDecoder(response.Body).Decode(&oauthResponse)
	if err != nil {
		return nil, err
	}

	c.JWT = &JWT{
		Time:        time.Now().Local().Add(time.Second * time.Duration(oauthResponse.ExpiresIn)),
		AccessToken: &oauthResponse.AccessToken,
	}

	return c.JWT, nil
}

func (c *client) GetAccessToken(ctx context.Context) (*JWT, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if c.JWT != nil && c.JWT.Time.UnixNano() > time.Now().Add(time.Second*time.Duration(c.Config.ExpirationMarginSeconds)).UnixNano() {
		bavalogs.Debug(ctx).Msg("token active, returning")
		return c.JWT, nil
	}

	bavalogs.Debug(ctx).Msg("token expired, creating new token")
	return c.createAccessToken(ctx)
}
