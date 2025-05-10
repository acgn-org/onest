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
var kFileLock sync.RWMutex

func Pathname() string {
	pathname := os.Getenv(EnvConfig)
	if pathname == "" {
		pathname = "config.yaml"
	}
	return pathname
}

func LoadConfigFile(logger logfield.LoggerWithFields) error {
	kFileLock.Lock()
	defer kFileLock.Unlock()

	err := kFile.Load(file.Provider(Pathname()), yaml.Parser())
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

type ScopedConfig[T any] struct {
	lock  sync.RWMutex
	scope string
	value T
}

func (c *ScopedConfig[T]) Get() T {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.value
}

func (c *ScopedConfig[T]) Save(value T) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := Save(c.scope, value)
	if err != nil {
		return err
	}
	c.value = value
	return nil
}

// Load scope is used in both loading from kFile and env
func Load[T any](scope string, defaults *T) *ScopedConfig[T] {
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
	kFileLock.RLock()
	defer kFileLock.RUnlock()
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
	return &ScopedConfig[T]{
		scope: scope,
		value: conf,
	}
}

func Save(scope string, value any) error {
	var k = koanf.New(".")
	if err := k.Load(structs.Provider(value, "yaml"), nil); err != nil {
		return err
	}

	kFileLock.Lock()
	defer kFileLock.Unlock()

	if err := kFile.MergeAt(k, scope); err != nil {
		return err
	}

	data, err := kFile.Marshal(yaml.Parser())
	if err != nil {
		return err
	}

	return os.WriteFile(Pathname(), data, 0600)
}
