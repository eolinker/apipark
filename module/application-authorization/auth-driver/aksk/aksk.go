package aksk

import (
	"encoding/json"
	"fmt"

	auth_driver "github.com/eolinker/apipark/module/application-authorization/auth-driver"

	application_authorization_dto "github.com/eolinker/apipark/module/application-authorization/dto"
)

var _ auth_driver.IAuthConfig = &Config{}

const (
	driver = "aksk"
)

func init() {
	auth_driver.RegisterAuthFactory(driver, auth_driver.NewFactory[Config](driver))
}

type Config struct {
	Ak    string            `json:"ak"`
	Sk    string            `json:"sk"`
	Label map[string]string `json:"label"`
}

func (a *Config) ID() string {
	//TODO implement me
	panic("implement me")
}

func (a *Config) Valid() ([]byte, error) {
	if a.Ak == "" {
		return nil, fmt.Errorf("access key is empty")
	}
	if a.Sk == "" {
		return nil, fmt.Errorf("secret key is empty")
	}
	return json.Marshal(a)
}

func (a *Config) Detail() []application_authorization_dto.DetailItem {
	return []application_authorization_dto.DetailItem{
		{
			Key:   "Access Key",
			Value: a.Ak,
		},
		{
			Key:   "Secret Key",
			Value: a.Sk,
		},
	}
}
