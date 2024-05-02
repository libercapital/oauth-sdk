# Welcome to OAuth SDK 👋

> SDK client to connect into authenticator OAuth

## Install

```go
go github.com/libercapital/oauth-sdk
```

## Usage

```go
client, err := oauthsdk.NewClient(
Config{
  URL:          "http://localhost:1000",
  ClientID:     "CLIENT_ID",
  ClientSecret: "CLIENT_SECRET",
  GrantType:    "GRANT_TYPE",
  Audience:     "AUDIENCE",
  RedactKeys:   []string{"scope"},
  ExpirationMargin:   5,
  CertData: &CertData{
    ClientCrt:    `xpto`,
    ClientCrtKey: `xpto`,
    },
  },
)

client.GetAccessToken(context.Background())
```

## Author

👤 **Giuseppe Menti**

- Gitlab: [@giuseppe.menti@bavabank.com](https://gitlab.com/giuseppe.menti)

## Contributors

👤 **Eduardo Mello**

- Github: [@EduardoRMello](https://github.com/EduardoRMello)

## Show your support

Give a ⭐️ if this project helped you!

---

_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_
