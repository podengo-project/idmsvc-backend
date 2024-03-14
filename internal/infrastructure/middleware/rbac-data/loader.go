package rbac_data

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type RbacFile struct {
	Version string  `yaml:"version"`
	Prefix  string  `yaml:"prefix"`
	Data    RBACMap `yaml:"data"`
}

const rbacFileVersion = "1.0"

func RBACMapLoad(data []byte) (string, RBACMap) {
	var dataParsed RbacFile
	if err := yaml.Unmarshal(data, &dataParsed); err != nil {
		panic(errors.New(err.Error()))
	}
	if dataParsed.Version != rbacFileVersion {
		panic(fmt.Errorf("it was expected rbac map config data version '%s', but '%s' was found", rbacFileVersion, dataParsed.Version))
	}
	return dataParsed.Prefix, dataParsed.Data
}
