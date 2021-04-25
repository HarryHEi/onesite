package config

import (
	"fmt"
)

var CfgRootPath = "configs"

func GetCfgPath(name string) string {
	return fmt.Sprintf("%s/%s", CfgRootPath, name)
}
