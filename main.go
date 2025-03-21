package vm6

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Vm6 struct {
	host     string
	email    string
	password string
	token    string
	client   *http.Client
}

func New(domain, email, password string) *Vm6 {
	return &Vm6{
		host:     domain,
		email:    email,
		password: password,
		client:   &http.Client{},
	}
}

func (v *Vm6) Send(method, version, service, function string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", v.host, service, version, function)

	var reqBody []byte
	var err error
	if data != nil {
		reqBody, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	for {
		req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))

		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		if v.token != "" {
			req.Header.Set("x-xsrf-token", v.token)
		}

		resp, err := v.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 503 {
			continue
		}

		if resp.StatusCode > 201 {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
		}

		return body, nil
	}
}

func (v *Vm6) Login() error {
	data := map[string]string{"email": v.email, "password": v.password}
	resp, err := v.Send("POST", "v4", "auth", "public/token", data)

	if err != nil {
		return err
	}

	var result struct {
		Confirmed bool
		ExpiresAt any `json:"expires_at"`
		Id        int `json:"id"`
		Token     string
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return err
	}

	v.token = result.Token

	return nil
}

func (v *Vm6) GetAuthKey(userID string) (string, error) {
	resp, err := v.Send("POST", "v4", "auth", fmt.Sprintf("user/%s/key", userID), nil)
	if err != nil {
		return "", err
	}
	var result map[string]string
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}
	return result["key"], nil
}

func (v *Vm6) Create(data map[string]interface{}) (map[string]interface{}, error) {
	resp, err := v.Send("POST", "v3", "vm", "host", data)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (v *Vm6) Remove(id string) error {
	_, err := v.Send("DELETE", "v3", "vm", fmt.Sprintf("host/%s", id), nil)
	return err
}

func (v *Vm6) VM(id string) (map[string]interface{}, error) {
	resp, err := v.Send("GET", "v3", "vm", fmt.Sprintf("host/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (v *Vm6) Start(id string) error {
	_, err := v.Send("POST", "v3", "vm", fmt.Sprintf("host/%s/start", id), nil)
	return err
}

func (v *Vm6) Stop(id string) error {
	_, err := v.Send("POST", "v3", "vm", fmt.Sprintf("host/%s/stop", id), map[string]bool{"force": true})
	return err
}

func (v *Vm6) Restart(id string) error {
	_, err := v.Send("POST", "v3", "vm", fmt.Sprintf("host/%s/restart", id), map[string]bool{"force": true})
	return err
}

func (v *Vm6) VncSettings(id string) (map[string]interface{}, error) {
	resp, err := v.Send("GET", "v3", "vm", fmt.Sprintf("host/%s/vnc_settings", id), nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}
