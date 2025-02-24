package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
)

type TokenService struct {
	config *configuration.Config
	token  *model.JwtToken
}

func NewTokenService(config *configuration.Config) *TokenService {
	return &TokenService{config: config, token: nil}
}
func (o *TokenService) getToken() string {
	return o.token.Token
}

func (o *TokenService) checkToken() error {
	var err error
	if o.token == nil {
		err = o.createNewToken()
		if err != nil {
			return err
		}
	} else {
		err = o.createNewToken()
	}
	return nil
}

func (o *TokenService) createNewToken() error {

	payload := strings.NewReader(
		fmt.Sprintf("grant_type=client_credentials&client_id=%s", o.config.OAuthClientId))

	req, err := http.NewRequest("POST", o.config.KeycloakUrl, payload)
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(o.config.OAuthUser, o.config.OAuthPassword)

	res, err := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("Non-OK HTTP status:", res.StatusCode)
		return errors.New("Create Token returns Non-OK HTTP status: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	tokenMap := map[string]interface{}{}
	err = json.Unmarshal(body, &tokenMap)
	if err != nil {
		return err
	}

	token := model.JwtToken{}
	token.Token = tokenMap["access_token"].(string)
	token.TokenExpiry = time.Now().Add(time.Duration(tokenMap["expires_in"].(float64)-300) * time.Second)
	token.RefreshTokenExpiry = time.Now().Add(time.Duration(tokenMap["refresh_expires_in"].(float64)-300) * time.Second)
	o.token = &token
	return err
}

func (o *TokenService) RequestData(method string, url string, payload []byte) (error, []byte, string) {
	err := o.checkToken()
	if err != nil {
		return err, nil, ""
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return err, nil, ""
	}
	req.Header.Add("content-type", "application/ld+json; version=2.0.0-dev")
	req.Header.Add("authorization", "Bearer "+o.getToken())
	res, err := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(res.Body)
	if err != nil {
		return err, nil, ""
	}
	if res.StatusCode != http.StatusOK {
		slog.Error("Non-OK HTTP status:", res.StatusCode)
		return errors.New("Non-OK HTTP status: " + res.Status), nil, ""
	}
	body, err := io.ReadAll(res.Body)
	header := res.Header.Get("location")
	if err != nil {
		return err, nil, ""
	}
	return nil, body, header
}
