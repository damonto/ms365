package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"office365/config"
	"office365/model"
	"strings"
	"time"

	"github.com/bluele/gcache"
	"github.com/prometheus/common/log"
)

var gc = gcache.New(100).LRU().Build()

// GetAccessToken 获取 access_token
func GetAccessToken(code string) error {
	requestBody := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {config.AppConfig.ClientID},
		"code":          {code},
		"redirect_uri":  {fmt.Sprintf("%v/oauth/callback", config.AppConfig.Domain)},
		"client_secret": {config.AppConfig.ClientSecret},
	}

	resp, err := http.PostForm(config.AppConfig.TokenURL, requestBody)
	if err != nil {
		return err
	}

	var accessToken map[string]interface{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&accessToken)

	if accessToken["error"] != nil {
		return fmt.Errorf("%v", accessToken["error_description"])
	}

	me, err := GetMe(accessToken["access_token"].(string))
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	var result = model.Account{}
	model.DB.Where(model.Account{Email: me.Email}).Assign(model.Account{
		UserID:       me.ID,
		Email:        me.Email,
		AccessToken:  accessToken["access_token"].(string),
		RefreshToken: accessToken["refresh_token"].(string),
		ExpiresIn:    int(accessToken["expires_in"].(float64)),
	}).FirstOrCreate(&result)

	gc.SetWithExpire(result.UserID, result.AccessToken, time.Duration(result.ExpiresIn)*time.Second)

	return nil
}

func getAccessTokenFromCache(userID string) (string, error) {
	cachedAccessToken, err := gc.Get(userID)
	if err != nil {
		var account = model.Account{}
		model.DB.First(&account, "user_id = ?", userID)
		requestBody := url.Values{
			"grant_type":    {"refresh_token"},
			"client_id":     {config.AppConfig.ClientID},
			"refresh_token": {account.RefreshToken},
			"client_secret": {config.AppConfig.ClientSecret},
		}

		resp, err := http.PostForm(config.AppConfig.TokenURL, requestBody)
		if err != nil {
			return "", err
		}

		var accessToken map[string]interface{}
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&accessToken)

		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		if accessToken["error"] != nil {
			panic(fmt.Errorf("%v", accessToken["error_description"]))
		}

		var result = model.Account{}
		model.DB.Model(&result).Where(model.Account{UserID: account.UserID}).Updates(model.Account{
			AccessToken:  accessToken["access_token"].(string),
			RefreshToken: accessToken["refresh_token"].(string),
			ExpiresIn:    int(accessToken["expires_in"].(float64)),
		})

		gc.SetWithExpire(account.UserID, result.AccessToken, time.Duration(result.ExpiresIn)*time.Second)

		return result.AccessToken, nil
	}

	return cachedAccessToken.(string), nil
}

// Me 我的个人信息
type Me struct {
	ID    string
	Email string
}

// GetMe 获取个人信息
func GetMe(accessToken string) (Me, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		return Me{}, err
	}

	var meResp map[string]interface{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&meResp)

	return Me{
		ID:    meResp["id"].(string),
		Email: meResp["mail"].(string),
	}, nil
}

// GetSubscribedSkus 订阅的 Skus
func GetSubscribedSkus(userID string) ([]interface{}, error) {
	accessToken, _ := getAccessTokenFromCache(userID)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/subscribedSkus", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	var subscribedSkus map[string]interface{}
	if err != nil {
		return make([]interface{}, 0), err
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&subscribedSkus)

	if subscribedSkus["error"] != nil {
		return make([]interface{}, 0), fmt.Errorf("%v", subscribedSkus["error"].(map[string]interface{})["code"])
	}

	return filterSubscribedSkus(subscribedSkus), nil
}

// Sku 可用的 SKU
type Sku struct {
	SkuID        string  `json:"sku_id"`
	Total        float64 `json:"total"`
	Used         float64 `json:"used"`
	FriendlyName string  `json:"friendly_name"`
}

func filterSubscribedSkus(subscribedSkus map[string]interface{}) []interface{} {
	mySkus := subscribedSkus["value"].([]interface{})
	var availableSkus []interface{}
	for _, sku := range mySkus {
		sku := sku.(map[string]interface{})
		skuPartNum := fmt.Sprintf("%v", sku["skuPartNumber"])
		if sku["capabilityStatus"] == "Enabled" && Skus[skuPartNum] != "" {
			availableSkus = append(availableSkus, &Sku{
				SkuID:        sku["skuId"].(string),
				Total:        sku["prepaidUnits"].(map[string]interface{})["enabled"].(float64),
				Used:         sku["consumedUnits"].(float64),
				FriendlyName: Skus[skuPartNum],
			})
		}
	}

	return availableSkus
}

type newUser struct {
	AccountEnabled    bool            `json:"accountEnabled"`
	DisplayName       string          `json:"displayName"`
	MailNickname      string          `json:"mailNickname"`
	UserPrincipalName string          `json:"userPrincipalName"`
	UsageLocation     string          `json:"usageLocation"`
	PasswordProfile   passwordProfile `json:"passwordProfile"`
}

type passwordProfile struct {
	ForceChangePasswordNextSignIn bool   `json:"forceChangePasswordNextSignIn"`
	Password                      string `json:"password"`
}

// CreateUser 创建新的 Office 365 用户
func CreateUser(userID string, enabled bool, nickname string, email string, password string, domain string) (map[string]interface{}, error) {
	accessToken, _ := getAccessTokenFromCache(userID)
	newUser := newUser{
		AccountEnabled:    enabled,
		DisplayName:       nickname,
		MailNickname:      strings.Trim(email, " "),
		UserPrincipalName: strings.Trim(email, " ") + "@" + domain,
		UsageLocation:     "US",
		PasswordProfile: passwordProfile{
			ForceChangePasswordNextSignIn: false,
			Password:                      password,
		},
	}
	jsonBody, _ := json.Marshal(newUser)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	var user map[string]interface{}
	if err != nil {
		return user, err
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&user)
	if user["error"] != nil {
		e := user["error"].(map[string]interface{})
		return user, fmt.Errorf("%v", e["message"])
	}

	return user, nil
}

type assignLicense struct {
	AddLicenses    []addLicenses `json:"addLicenses"`
	RemoveLicenses []string      `json:"removeLicenses"`
}
type addLicenses struct {
	DisabledPlans []string `json:"disabledPlans"`
	SkuID         string   `json:"skuId"`
}

// AssignLicense 分配给用户 License
func AssignLicense(accountID string, SkuID string, userID string) error {
	accessToken, _ := getAccessTokenFromCache(accountID)
	requestData := assignLicense{
		AddLicenses: []addLicenses{addLicenses{
			DisabledPlans: []string{},
			SkuID:         SkuID,
		}},
		RemoveLicenses: []string{},
	}

	client := &http.Client{}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(requestData)
	req, _ := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/users/"+userID+"/assignLicense", buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	var license map[string]interface{}
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&license)
	if license["error"] != nil {
		e := license["error"].(map[string]interface{})
		return fmt.Errorf("%v", e["message"])
	}

	return nil
}
