package cmd

import (
	"github.com/spf13/viper"
	"testing"
)

func TestCompareAndSync(t *testing.T) {
	s := map[string]interface{}{
		"default": "http://localhost:8080/",
	}
	viper.Set("remoteDefault", "default")
	viper.Set("remote", s)
	CompareAndSync("gg-example", "")

}
