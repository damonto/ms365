package microsoft

import (
	"errors"
	"net/url"

	"github.com/valyala/fastjson"
)

// User Represents an Azure AD user account.
type User struct {
	GraphAPI *GraphAPI
}

const (
	defaultPageSize = "10" // 10 records per request
)

// NewUser returns user instance
func NewUser() *User {
	return &User{
		GraphAPI: NewGraphAPI(),
	}
}

// UserResponse is the user struct
type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"mail"`
	DispalyName string `json:"display_name"`
}

// Users retrieve a list of user objects.
func (u *User) Users(id string, next string) (response map[string]interface{}, err error) {
	req, err := u.GraphAPI.newRequest(id)
	if err != nil {
		return response, err
	}

	req = req.SetQueryParam("$top", defaultPageSize).
		SetQueryParam("$select", "id,displayName,mail")
	if next != "" {
		req.SetQueryParam("$skiptoken", next)
	}

	resp, err := req.Get(u.GraphAPI.uri("/v1.0/users"))
	if err != nil {
		return response, err
	}

	var parser fastjson.Parser
	uResp, err := parser.ParseBytes(resp.Body())
	if err != nil {
		return response, err
	}
	if len(uResp.GetStringBytes("error")) > 0 {
		return response, errors.New(string(uResp.GetStringBytes("error_description")))
	}

	nextLink, err := url.Parse(string(uResp.GetStringBytes("@odata.nextLink")))
	if err != nil {
		return response, err
	}
	q, err := url.ParseQuery(nextLink.RawQuery)

	response = make(map[string]interface{}, 1)
	var respUsers = []UserResponse{}
	response["next"] = q.Get("$skiptoken")
	for _, v := range uResp.GetArray("value") {
		respUsers = append(respUsers, UserResponse{
			ID:          string(v.GetStringBytes("id")),
			Email:       string(v.GetStringBytes("mail")),
			DispalyName: string(v.GetStringBytes("displayName")),
		})
	}
	response["users"] = respUsers

	return response, nil
}

// Delete user
func (u *User) Delete(id string, uid string) error {
	req, err := u.GraphAPI.newRequest(id)
	if err != nil {
		return err
	}
	resp, err := req.Get(u.GraphAPI.uri("/v1.0/users/" + uid))
	if err != nil {
		return err
	}

	var parser fastjson.Parser
	uResp, err := parser.ParseBytes(resp.Body())
	if err != nil {
		return err
	}
	if len(uResp.GetStringBytes("error")) > 0 {
		return errors.New(string(uResp.GetStringBytes("error_description")))
	}

	return nil
}
