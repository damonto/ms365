package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"office365/config"
	"office365/model"
	"time"
)

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

	expiresIn := time.Duration(accessToken["expires_in"].(float64))
	model.Db.Create(&model.Account{
		UserID:       me.ID,
		Email:        me.Email,
		AccessToken:  accessToken["access_token"].(string),
		RefreshToken: accessToken["refresh_token"].(string),
		ExpiresIn:    time.Now().Add(expiresIn * time.Second),
	})

	return nil
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
func GetSubscribedSkus() ([]interface{}, error) {
	accessToken := ""
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
	Total        int64
	Used         int64
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
				Total:        sku["prepaidUnits"].(map[string]interface{})["enabled"].(int64),
				Used:         sku["consumedUnits"].(int64),
				FriendlyName: Skus[skuPartNum],
			})
		}
	}

	return availableSkus
}
