package microsoft

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/damonto/msonline/internal/pkg/config"
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

func (ga *GraphAPI) uri(path string) string {
	return strings.TrimRight(config.Cfg.Microsoft.Endpoint, "/") + path
}

func (ga *GraphAPI) newRequest(id string) (req *resty.Request, err error) {
	token, err := ga.getAccessToken(id)
	if time.Now().After(token.ExpireDate) {
		token, err = ga.refreshAccessToken(token)
	}

	req = resty.New().
		R().
		SetAuthToken(token.AccessToken)

	return req, err
}

func (ga *GraphAPI) refreshAccessToken(token AccessToken) (AccessToken, error) {
	resp, err := ga.resty.
		SetFormData(map[string]string{
			"grant_type":    "refresh_token",
			"client_id":     config.Cfg.Microsoft.ClientID,
			"client_secret": config.Cfg.Microsoft.ClientSecret,
			"refresh_token": token.RefreshToken,
		}).
		Post(config.Cfg.Microsoft.TokenURL)

	var parser fastjson.Parser
	parsedToken, err := parser.ParseBytes(resp.Body())
	if err != nil {
		return token, err
	}
	if len(parsedToken.GetStringBytes("error")) > 0 {
		return token, errors.New(string(parsedToken.GetStringBytes("error_description")))
	}

	// Update token
	token.AccessToken = string(parsedToken.GetStringBytes("access_token"))
	token.RefreshToken = string(parsedToken.GetStringBytes("refresh_token"))
	token.ExpireDate = time.Now().Add(time.Duration(parsedToken.GetInt64("expires_in")) * time.Second)
	err = NewStore().Put(token.ID, token)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (ga *GraphAPI) getAccessToken(id string) (accessToken AccessToken, err error) {
	accessToken, err = NewStore().Get(id)
	if err != nil {
		return accessToken, errors.New(err.Error())
	}

	return accessToken, nil
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
	user, err := ga.resty.SetAuthToken(token).Get(ga.uri("/v1.0/me"))
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
