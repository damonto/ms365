package microsoft

import (
	"fmt"
	"time"

	"github.com/damonto/office365/internal/pkg/config"
	"github.com/go-resty/resty/v2"
	"github.com/valyala/fastjson"
)

// GraphAPI is the microsoft graph api instance
type GraphAPI struct {
	resty    *resty.Request
	fastjson fastjson.Parser
}

// NewGraphAPI returns microsoft graph api instance
func NewGraphAPI() *GraphAPI {
	return &GraphAPI{
		resty: resty.New().R(),
	}
}

// AccessToken is the microsoft graph api access token sturct
type AccessToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpireDate   time.Time `json:"expire_date"`
}

// GetAccessToken get an token with authorization	 `code`
func (ga *GraphAPI) GetAccessToken(code string) error {
	resp, err := ga.resty.
		SetBody(map[string]string{
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

	parsed, err := ga.fastjson.ParseBytes(resp.Body())
	if err != nil {
		return err
	}
	if len(parsed.GetStringBytes("error")) >= 0 {
		return fmt.Errorf("%v", parsed.GetStringBytes("error_description"))
	}

	accessToken := AccessToken{
		AccessToken:  string(parsed.GetStringBytes("access_token")),
		RefreshToken: string(parsed.GetStringBytes("refresh_token")),
		ExpireDate:   time.Now().Add(time.Duration(parsed.GetInt64("expires_in")) * time.Second),
	}

	fmt.Println(accessToken)

	return nil
}
