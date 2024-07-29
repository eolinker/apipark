package access

import (
	_ "embed"

	"github.com/eolinker/go-common/access"
	"gopkg.in/yaml.v3"
)

type Access = access.Access

var (
	//go:embed access.yaml
	data []byte
)

func init() {
	ts := make(map[string][]Access)
	err := yaml.Unmarshal(data, &ts)
	if err != nil {
		panic(err)
	}
	for group, asl := range ts {

		access.Add(group, asl)

	}
}
