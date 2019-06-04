package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"office365/config"
)

// GetAccessToken 获取 access_token
func GetAccessToken(code string) (string, error) {
	requestBody := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {config.AppConfig.ClientID},
		"code":          {code},
		"redirect_uri":  {fmt.Sprintf("%v/oauth/callback", config.AppConfig.Domain)},
		"client_secret": {config.AppConfig.ClientSecret},
	}

	resp, err := http.PostForm(config.AppConfig.TokenURL, requestBody)
	if err != nil {
		return "", err
	}

	var respJSON map[string]interface{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&respJSON)

	if respJSON["error"] != nil {
		return "", fmt.Errorf("%v", respJSON["error_description"])
	}

	return fmt.Sprintf("%v", respJSON["access_token"]), nil
}

// GetSubscribedSkus 订阅的 Skus
func GetSubscribedSkus(accessToken string) ([]interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/subscribedSkus", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		return make([]interface{}, 0), err
	}

	var respJSON map[string]interface{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&respJSON)

	if respJSON["error"] != nil {
		return make([]interface{}, 0), fmt.Errorf("%v", respJSON["error"].(map[string]interface{})["code"])
	}

	return filterSubscribedSkus(respJSON), nil
}

// Sku 可用的 SKU
type Sku struct {
	SkuID        interface{}
	Total        interface{}
	Used         interface{}
	FriendlyName string
}

func filterSubscribedSkus(respJSON map[string]interface{}) []interface{} {
	subscribedSkus := respJSON["value"].([]interface{})

	var availableSkus []interface{}
	for _, sku := range subscribedSkus {
		sku := sku.(map[string]interface{})
		skuPartNum := fmt.Sprintf("%v", sku["skuPartNumber"])
		if sku["capabilityStatus"] == "Enabled" && Skus[skuPartNum] != "" {
			availableSkus = append(availableSkus, &Sku{
				SkuID:        sku["skuId"],
				Total:        sku["prepaidUnits"].(map[string]interface{})["enabled"],
				Used:         sku["consumedUnits"],
				FriendlyName: Skus[skuPartNum],
			})
		}
	}

	return availableSkus
}

