package config

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var conf = `
secret_key = "It's a lonely road to the top."

[server]
bind = '0.0.0.0'
port = 8000

[db]
driver_name = "mysql"
dsn = "herui:Admin@123@tcp(172.172.177.191:3306)/onesite_dev?charset=utf8mb4&parseTime=true"
max_open_conn = 10
max_idle_conn = 5

[redis]
addr = "172.172.177.191:6379"
password = ""
db = 0

[auth]
timeout = "8m03s"
`

func TestInit(t *testing.T) {
	testFile := "/tmp/oneSiteTestFile.toml"
	err := ioutil.WriteFile(testFile, []byte(conf), 0644)
	require.Nil(t, err)
	err = Init(testFile)
	require.Nil(t, err)
	err = os.Remove(testFile)
	require.Nil(t, err)
}
