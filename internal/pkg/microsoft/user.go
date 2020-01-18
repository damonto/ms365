package microsoft

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

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
	Email       string `json:"email"`
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
	if uResp.Exists("error") {
		return response, errors.New(string(uResp.Get("error").GetStringBytes("message")))
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
func (u *User) Delete(id string, uid string) (err error) {
	req, err := u.GraphAPI.newRequest(id)
	if err != nil {
		return err
	}
	resp, err := req.Delete(u.GraphAPI.uri("/v1.0/users/" + uid))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		var parser fastjson.Parser
		uResp, err := parser.ParseBytes(resp.Body())
		if err != nil {
			return err
		}
		if uResp.Exists("error") {
			return errors.New(string(uResp.Get("error").GetStringBytes("message")))
		}
	}

	return nil
}

// CreateUserRequest user struct
type CreateUserRequest struct {
	Name          string `json:"name" binding:"required"`
	PrincipalName string `json:"principal_name" binding:"required"`
	Domain        string `json:"domain"`
	Password      string `json:"password" binding:"required"`
	SkuID         string `json:"sku_id"`
}

// Create new user
func (u *User) Create(id string, cr CreateUserRequest) (uid string, err error) {
	req, err := u.GraphAPI.newRequest(id)
	if err != nil {
		return uid, err
	}

	if cr.Domain == "" {
		token, err := u.GraphAPI.getAccessToken(id)
		if err != nil {
			return uid, err
		}
		boom := strings.Split(token.Email, "@")
		cr.Domain = boom[1]
	}

	resp, err := req.
		SetBody(map[string]interface{}{
			"accountEnabled":    true,
			"displayName":       cr.Name,
			"mailNickname":      cr.PrincipalName,
			"userPrincipalName": cr.PrincipalName + "@" + cr.Domain,
			"usageLocation":     "US",
			"passwordProfile": map[string]interface{}{
				"forceChangePasswordNextSignIn": false,
				"password":                      cr.Password,
			},
		}).
		Post(u.GraphAPI.uri("/v1.0/users"))

	var parser fastjson.Parser
	uResp, err := parser.ParseBytes(resp.Body())
	if err != nil {
		return uid, err
	}
	if uResp.Exists("error") {
		return uid, errors.New(string(uResp.Get("error").GetStringBytes("message")))
	}

	uid = string(uResp.GetStringBytes("id"))
	if cr.SkuID != "" {
		err = u.assignLicense(id, uid, cr.SkuID)
		if err != nil {
			return uid, err
		}
	}

	return uid, nil
}

func (u *User) assignLicense(id string, uid string, skuID string) (err error) {
	req, err := u.GraphAPI.newRequest(id)
	if err != nil {
		return err
	}

	resp, err := req.
		SetBody(map[string]interface{}{
			"addLicenses": []map[string]interface{}{
				{
					"disabledPlans": []string{},
					"skuId":         skuID,
				},
			},
			"removeLicenses": []string{},
		}).
		Post(u.GraphAPI.uri("/v1.0/users/" + uid + "/assignLicense"))

	var parser fastjson.Parser
	assignResp, err := parser.ParseBytes(resp.Body())
	if err != nil {
		return err
	}
	if assignResp.Exists("error") {
		return errors.New(string(assignResp.Get("error").GetStringBytes("message")))
	}

	return nil
}

// MSUser is microsoft user object
type MSUser struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Location    string `json:"location"`
}

// Retrieve the properties and relationships of user object.
func (u *User) Retrieve(id string, uid string) (user MSUser, err error) {
	req, err := u.GraphAPI.newRequest(id)
	if err != nil {
		return user, err
	}
	resp, err := req.Get(u.GraphAPI.uri("/v1.0/users/" + uid))
	var parser fastjson.Parser
	uResp, err := parser.ParseBytes(resp.Body())
	if err != nil {
		return user, err
	}
	if uResp.Exists("error") {
		return user, errors.New(string(uResp.Get("error").GetStringBytes("message")))
	}

	return MSUser{
		ID:          string(uResp.GetStringBytes("id")),
		Email:       string(uResp.GetStringBytes("mail")),
		DisplayName: string(uResp.GetStringBytes("displayName")),
		Location:    string(uResp.GetStringBytes("officeLocation")),
	}, nil
}

// UpdatePassword reset user password
func (u *User) UpdatePassword(id string, uid string, password string) (err error) {
	req, err := u.GraphAPI.newRequest(id)
	if err != nil {
		return err
	}
	resp, err := req.
		SetBody(map[string]interface{}{
			"passwordProfile": map[string]interface{}{
				"forceChangePasswordNextSignIn": false,
				"password":                      password,
			},
		}).
		Patch(u.GraphAPI.uri("/v1.0/users/" + uid))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		var parser fastjson.Parser
		uResp, err := parser.ParseBytes(resp.Body())
		if err != nil {
			return err
		}
		if uResp.Exists("error") {
			return errors.New(string(uResp.Get("error").GetStringBytes("message")))
		}
	}

	return nil
}
