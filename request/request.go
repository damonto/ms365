package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"office365/config"
	"office365/model"
	"time"

	"github.com/bluele/gcache"
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

		if accessToken["error"] != nil {
			return "", fmt.Errorf("%v", accessToken["error_description"])
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
	if err != nil {
		return make([]interface{}, 0), err
	}

	var subscribedSkus map[string]interface{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&subscribedSkus)

	if subscribedSkus["error"] != nil {
		return make([]interface{}, 0), fmt.Errorf("%v", subscribedSkus["error"].(map[string]interface{})["code"])
	}

	return filterSubscribedSkus(subscribedSkus), nil
}

// Sku 可用的 SKU
type Sku struct {
	SkuID        string
	Total        float64
	Used         float64
	FriendlyName string
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
