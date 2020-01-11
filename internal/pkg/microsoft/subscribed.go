package microsoft

import (
	"errors"

	"github.com/valyala/fastjson"
)

// Subscribed contains information about a service SKU that a company is subscribed to.
type Subscribed struct {
	GraphAPI *GraphAPI
}

// NewSubscribed returns subscribed instance
func NewSubscribed() *Subscribed {
	return &Subscribed{
		GraphAPI: NewGraphAPI(),
	}
}

// Sku is commercial subscription struct
type Sku struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ConsumedUnits int64  `json:"consumed_units"`
	PrepaidUnits  int64  `json:"prepaid_units"`
}

// ListSubscribedSkus get the list of commercial subscriptions that an organization has acquired.
func (s *Subscribed) ListSubscribedSkus(id string) (skus []Sku, err error) {
	req, err := s.GraphAPI.newRequest(id)
	if err != nil {
		return skus, err
	}

	resp, err := req.Get(s.GraphAPI.uri("/v1.0/subscribedSkus"))
	if err != nil {
		return skus, err
	}

	var parser fastjson.Parser
	skusResp, err := parser.ParseBytes(resp.Body())
	if err != nil {
		return skus, err
	}
	if skusResp.Exists("error") {
		return skus, errors.New(string(skusResp.Get("error").GetStringBytes("message")))
	}

	subscribedSkus := skusResp.GetArray("value")
	for _, sku := range subscribedSkus {
		if string(sku.GetStringBytes("capabilityStatus")) == "Enabled" {
			skus = append(skus, Sku{
				ID:            string(sku.GetStringBytes("skuId")),
				Name:          Skus[string(sku.GetStringBytes("skuPartNumber"))],
				ConsumedUnits: sku.GetInt64("consumedUnits"),
				PrepaidUnits:  sku.Get("prepaidUnits").GetInt64("enabled"),
			})
		}
	}

	return skus, nil
}
