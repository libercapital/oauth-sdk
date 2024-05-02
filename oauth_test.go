package oauthsdk

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	bavahelper "gitlab.com/bavatech/architecture/software/libs/go-modules/bava-helper.git"
)

func Test_client_GetAccessToken(t *testing.T) {
	expireIn := time.Now().Add(time.Minute * 5)

	type fields struct {
		Config Config
		Mutex  *sync.Mutex
		JWT    *JWT
		http   requestClientMock
	}
	tests := []struct {
		name           string
		fields         fields
		want           *JWT
		wantErr        error
		mockBehavior   func(fields)
		assertBehavior func(*JWT)
	}{
		{
			name: "should create access token and return created token",
			fields: fields{
				Config: Config{
					URL:          "http://localhost:1000",
					ClientID:     "CLIENT_ID",
					ClientSecret: "CLIENT_SECRET",
					GrantType:    "GRANT_TYPE",
					Audience:     "AUDIENCE",
				},
				Mutex: &sync.Mutex{},
				JWT:   nil,
				http: requestClientMock{
					Responses: []mockResponse{
						{
							Url:        "http://localhost:1000",
							StatusCode: 200,
							Body:       `{"access_token":"0000000","expires_in":43199}`,
						},
					},
				},
			},
			assertBehavior: func(jwt *JWT) {
				assert.NotEmpty(t, jwt.AccessToken)
				assert.NotEmpty(t, jwt.Time)
				assert.Greater(t, jwt.Time.Unix(), time.Now().Unix())
			},
		},

		{
			name:    "error on create new request",
			wantErr: errors.New(`parse ":": missing protocol scheme`),
			fields: fields{
				Config: Config{
					URL:          ":##",
					ClientID:     "CLIENT_ID",
					ClientSecret: "CLIENT_SECRET",
					GrantType:    "GRANT_TYPE",
					Audience:     "AUDIENCE",
				},
				Mutex: &sync.Mutex{},
				JWT:   nil,
				http:  requestClientMock{},
			},
			assertBehavior: func(jwt *JWT) {
				assert.Nil(t, jwt)
			},
		},

		{
			name:    "error doing request",
			wantErr: errors.New("FAKE_ERROR"),
			fields: fields{
				Config: Config{
					URL:          "http://localhost:1000",
					ClientID:     "CLIENT_ID",
					ClientSecret: "CLIENT_SECRET",
					GrantType:    "GRANT_TYPE",
					Audience:     "AUDIENCE",
				},
				Mutex: &sync.Mutex{},
				JWT:   nil,
				http: requestClientMock{
					Error: errors.New("FAKE_ERROR"),
				},
			},
			assertBehavior: func(jwt *JWT) {
				assert.Nil(t, jwt)
			},
		},

		{
			name:    "error bad request",
			wantErr: errors.New("http request error (status=400)"),
			fields: fields{
				Config: Config{
					URL:          "http://localhost:1000",
					ClientID:     "CLIENT_ID",
					ClientSecret: "CLIENT_SECRET",
					GrantType:    "GRANT_TYPE",
					Audience:     "AUDIENCE",
				},
				Mutex: &sync.Mutex{},
				JWT:   nil,
				http: requestClientMock{
					Responses: []mockResponse{
						{
							Url:        "http://localhost:1000",
							StatusCode: 400,
							Body:       `{"access_token":"0000000","expires_in":43199}`,
						},
					},
				},
			},
			assertBehavior: func(jwt *JWT) {
				assert.Nil(t, jwt)
			},
		},

		{
			name:    "error decode response",
			wantErr: errors.New("unexpected EOF"),
			fields: fields{
				Config: Config{
					URL:          "http://localhost:1000",
					ClientID:     "CLIENT_ID",
					ClientSecret: "CLIENT_SECRET",
					GrantType:    "GRANT_TYPE",
					Audience:     "AUDIENCE",
				},
				Mutex: &sync.Mutex{},
				JWT:   nil,
				http: requestClientMock{
					Responses: []mockResponse{
						{
							Url:        "http://localhost:1000",
							StatusCode: 200,
							Body:       `{`,
						},
					},
				},
			},
			assertBehavior: func(jwt *JWT) {
				assert.Nil(t, jwt)
			},
		},

		{
			name: "success existent token",
			want: &JWT{AccessToken: bavahelper.PtrAny("abc"), Time: expireIn},
			fields: fields{
				Config: Config{
					URL:          "http://localhost:1000",
					ClientID:     "CLIENT_ID",
					ClientSecret: "CLIENT_SECRET",
					GrantType:    "GRANT_TYPE",
					Audience:     "AUDIENCE",
				},
				Mutex: &sync.Mutex{},
				JWT:   &JWT{AccessToken: bavahelper.PtrAny("abc"), Time: expireIn},
				http: requestClientMock{
					Responses: []mockResponse{
						{
							Url:        "http://localhost:1000",
							StatusCode: 200,
							Body:       `{`,
						},
					},
				},
			},
			assertBehavior: func(jwt *JWT) {
				assert.NotEmpty(t, jwt.AccessToken)
				assert.NotEmpty(t, jwt.Time)
				assert.Greater(t, jwt.Time.Unix(), time.Now().Unix())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client{
				Config: tt.fields.Config,
				Mutex:  tt.fields.Mutex,
				JWT:    tt.fields.JWT,
				http:   tt.fields.http,
			}
			got, err := c.GetAccessToken(context.TODO())
			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("client.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.assertBehavior(c.JWT)
			tt.assertBehavior(got)
		})
	}
}
