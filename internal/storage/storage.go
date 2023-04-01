package storage

import (
	zlog "github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var Config Data

func SaveTargets(s string) (*Data, error) {
	Config := &Data{}
	err := yaml.Unmarshal([]byte(s), &Config.Targets)
	if err != nil {
		zlog.Print("parsing failed: ", err)
		return nil, err
	}
	zlog.Print("parsing in SaveTargets OK")
	return Config, nil
}
