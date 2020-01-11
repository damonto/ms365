package microsoft

import (
	"errors"
	"fmt"
	"time"

	"github.com/damonto/msonline-webapi/internal/pkg/config"
	"github.com/go-resty/resty/v2"
	"github.com/valyala/fastjson"
)

// GraphAPI is the microsoft graph api instance
type GraphAPI struct {
	resty *resty.Request
}

// NewGraphAPI returns microsoft graph api instance
func NewGraphAPI() *GraphAPI {
	return &GraphAPI{
		resty: resty.New().R(),
	}
}

// GetAccessToken get an token with authorization `code`
func (ga *GraphAPI) GetAccessToken(code string) error {
	resp, err := ga.resty.
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     config.Cfg.Microsoft.ClientID,
			"client_secret": config.Cfg.Microsoft.ClientSecret,
			"code":          code,
			"redirect_uri":  fmt.Sprintf("%v/oauth/callback", config.Cfg.Microsoft.Domain),
		}).
		Post(config.Cfg.Microsoft.TokenURL)

	if err != nil {
		return err
	}

	var parser fastjson.Parser
	parsedToken, err := parser.ParseBytes(resp.Body())
	if err != nil {
		return err
	}
	if len(parsedToken.GetStringBytes("error")) > 0 {
		return errors.New(string(parsedToken.GetStringBytes("error_description")))
	}

	token := string(parsedToken.GetStringBytes("access_token"))
	user, err := ga.resty.SetAuthToken(token).Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		return err
	}

	parser = fastjson.Parser{}
	parsedUser, err := parser.ParseBytes(user.Body())
	if err != nil {
		return err
	}

	id := string(parsedUser.GetStringBytes("id"))
	NewStore().Put(id, AccessToken{
		ID:           id,
		Email:        string(parsedUser.GetStringBytes("mail")),
		AccessToken:  token,
		RefreshToken: string(parsedToken.GetStringBytes("refresh_token")),
		ExpireDate:   time.Now().Add(time.Duration(parsedToken.GetInt64("expires_in")) * time.Second),
	})

	return nil
}
