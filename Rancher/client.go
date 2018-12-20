package rancher

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	AccessKey   string
	SecretKey   string
	Endpoint    string
	environment string
	service     string
	client      *http.Client
}

func (c Client) Environment(environment string) Client {
	c.environment = environment
	return c
}

func (c Client) Service(service string) Client {
	c.service = service
	return c
}

func (c Client) post(url string, content interface{}) error {

	var body io.Reader = nil
	if content != nil {
		jsonValue, err := json.Marshal(content)
		if err != nil {
			return errors.WithStack(err)
		}

		body = bytes.NewBuffer(jsonValue)
	}

	req, err := http.NewRequest("POST", c.Endpoint+url, body)
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Add("Authorization", c.authorization())
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		return errors.WithStack(fmt.Errorf("got Status code of %d", resp.StatusCode))
	}

	return nil
}

func (c Client) authorization() string {
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(c.AccessKey + ":" + c.SecretKey))
	return fmt.Sprintf("Basic %s", encodedAuth)
}

func (c Client) serviceAction(action string, content interface{}) error {
	if c.service == "" {
		return errors.WithStack(errors.New("Missing service"))
	}

	if c.environment == "" {
		return errors.WithStack(errors.New("Missing environment"))
	}

	return c.post(fmt.Sprintf("/projects/%s/services/%s/?action=%s", c.environment, c.service, action), content)
}

func (c Client) Upgrade(serviceUpgrade ServiceUpgrade) error {
	return c.serviceAction("upgrade", serviceUpgrade)
}

func (c Client) FinishUpgrade() error {
	return c.serviceAction("finishupgrade", nil)
}

type Service struct {
}

type LaunchConfig struct {
	ImageUuid string `json:"imageUuid"`
}

type InServiceStrategy struct {
	LaunchConfig LaunchConfig `json:"launchConfig"`
}

type ServiceUpgrade struct {
	InServiceStrategy InServiceStrategy `json:"inServiceStrategy"`
}

func (c Client) Get() (Service, error) {
	service := Service{}
	if c.service == "" {
		return service, errors.WithStack(errors.New("Missing service"))
	}

	if c.environment == "" {
		return service, errors.WithStack(errors.New("Missing environment"))
	}

	err := c.get(fmt.Sprintf("/projects/%s/services/%s", c.environment, c.service), &service)
	return service, err
}

func (c Client) get(url string, body interface{}) error {
	req, err := http.NewRequest("GET", c.Endpoint+url, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Add("Authorization", c.authorization())
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		return errors.WithStack(fmt.Errorf("got Status code of %d", resp.StatusCode))
	}

	return nil
}
