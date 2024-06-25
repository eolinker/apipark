package apikey

import (
	"encoding/json"
	"fmt"

	auth_driver "github.com/eolinker/apipark/module/project-authorization/auth-driver"

	"github.com/eolinker/go-common/utils"

	project_authorization_dto "github.com/eolinker/apipark/module/project-authorization/dto"
)

const (
	driver = "apikey"
)

func init() {
	auth_driver.RegisterAuthFactory(driver, auth_driver.NewFactory[Config](driver))
}

var _ auth_driver.IAuthConfig = (*Config)(nil)

type Config struct {
	Apikey string            `json:"apikey"`
	Label  map[string]string `json:"label"`
}

func (a *Config) ID() string {
	return utils.Md5(a.Apikey)
}

func (a *Config) Valid() ([]byte, error) {
	if a.Apikey == "" {
		return nil, fmt.Errorf("apikey is empty")
	}
	return json.Marshal(a)
}

func (a *Config) Detail() []project_authorization_dto.DetailItem {
	return []project_authorization_dto.DetailItem{
		{
			Key:   "Apikey",
			Value: a.Apikey,
		},
	}
}
