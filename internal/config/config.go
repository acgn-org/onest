package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"os"
	"strings"
)

const EnvPrefix = "ONEST"
const EnvConfig = EnvPrefix + "_CONFIG"

var k = koanf.New(".")
var kFile = koanf.New(".")

func LoadConfigFile() error {
	pathname := os.Getenv(EnvConfig)
	if pathname == "" {
		pathname = "config.yaml"
	}
	return kFile.Load(file.Provider(pathname), yaml.Parser())
}

// Load should to be called after kFile is loaded
// scope is used in both loading from kFile and env
func Load[T any](scope string, defaults *T) (*T, error) {
	var conf T

	// defaults
	if defaults != nil {
		err := k.Load(structs.Provider(&conf, "yaml"), nil)
		if err != nil {
			return nil, err
		}
	}

	// from file
	if err := k.Merge(kFile.Cut(scope)); err != nil {
		return nil, err
	}

	// from env
	var prefix = fmt.Sprintf(
		"%s_%s_",
		EnvPrefix,
		strings.Join(strings.Split(strings.ToUpper(scope), "."), "_"),
	)
	if err := k.Load(env.Provider(prefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, prefix)), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	return &conf, k.Unmarshal("", &conf)
}
