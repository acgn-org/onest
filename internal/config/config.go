package config

import (
	"fmt"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"os"
	"strings"
	"sync"
)

const EnvPrefix = "ONEST"
const EnvConfig = EnvPrefix + "_CONFIG"

var kFile = koanf.New(".")

func LoadConfigFile(logger logfield.LoggerWithFields) error {
	pathname := os.Getenv(EnvConfig)
	if pathname == "" {
		pathname = "config.yaml"
	}
	err := kFile.Load(file.Provider(pathname), yaml.Parser())
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warnln("config file not found, skipping...")
			return nil
		}
		return err
	}
	return nil
}

var loadConfigFileOnce = sync.OnceFunc(func() {
	logger := logfield.New(logfield.ComConfig).WithAction("load:file")
	err := LoadConfigFile(logger)
	if err != nil {
		logger.Fatalln("load config file failed:", err)
	}
})

// Load scope is used in both loading from kFile and env
func Load[T any](scope string, defaults *T) *T {
	logger := logfield.New(logfield.ComConfig).WithAction("load:" + scope)

	var conf T
	var k = koanf.New(".")

	// defaults
	if defaults != nil {
		err := k.Load(structs.Provider(defaults, "yaml"), nil)
		if err != nil {
			panic(err)
		}
	}

	// from file
	loadConfigFileOnce()
	if err := k.Merge(kFile.Cut(scope)); err != nil {
		panic(err)
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
		logger.Fatalln("load from env failed:", err)
	}

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{
		Tag: "yaml",
	}); err != nil {
		logger.Fatalln("unmarshal failed:", err)
	}
	return &conf
}
