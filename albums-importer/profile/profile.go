package profile

import (
	"fmt"
)

type Profile struct {
	URL string `yaml:"url"`
}

func LoadProfile(name string) (profile Profile, err error) {
	configMap, err := OpenConfigMap()
	if err != nil {
		err = fmt.Errorf("unable to open configuration file: %w", err)
		return
	}

	profile, ok := configMap[name]

	if !ok {
		err = fmt.Errorf("profile not found: %s", name)
		return
	}

	return
}
