package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/acgn-org/onest/internal/logfield"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	log "github.com/sirupsen/logrus"
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
	logger log.FieldLogger
	scope  string
	kEnv   *koanf.Koanf
	value  *atomic.Value
}

func (c *ScopedConfig[T]) Get() T {
	return c.value.Load().(T)
}

func (c *ScopedConfig[T]) Save(value T) error {
	kRaw := koanf.New(".")
	err := kRaw.Load(structs.Provider(c.value, "yaml"), nil)
	if err != nil {
		return err
	}

	kNew := koanf.New(".")
	err = kNew.Load(structs.Provider(value, "yaml"), nil)
	if err != nil {
		return err
	}

	kFileLock.Lock()
	defer kFileLock.Unlock()

	var kFileChanged bool
	for _, key := range kNew.Keys() {
		val := kNew.Get(key)
		if val != kRaw.Get(key) {
			globalKey := c.scope + "." + key

			if c.kEnv.Get(key) != nil {
				c.logger.Warnf("%s was set by environment, change is not saved to file and will not take effect on next startup", globalKey)
				continue
			}

			if kFile.Get(globalKey) != val {
				kFileChanged = true
				err := kFile.Set(c.scope+"."+key, val)
				if err != nil {
					return err
				}
			}
		}
	}

	if kFileChanged {
		data, err := kFile.Marshal(yaml.Parser())
		if err != nil {
			return err
		}

		if err := os.WriteFile(Pathname(), data, 0600); err != nil {
			return err
		}
	}

	c.value.Store(value)
	return nil
}

// LoadScoped scope is used in both loading from kFile and env
func LoadScoped[T any](scope string, defaults *T) *ScopedConfig[T] {
	logger := logfield.New(logfield.ComConfig).WithAction("load:" + scope)

	var conf T
	var k = koanf.New(".")

	// defaults
	if defaults != nil {
		err := k.Load(structs.Provider(defaults, "yaml"), nil)
		if err != nil {
			logger.Fatalln("load defaults failed:", err)
		}
	}

	// from file
	loadConfigFileOnce()
	kFileLock.RLock()
	defer kFileLock.RUnlock()
	if err := k.Merge(kFile.Cut(scope)); err != nil {
		logger.Fatalln("merge from file failed:", err)
	}

	// from env
	var kEnv = koanf.New(".")
	var prefix = fmt.Sprintf(
		"%s_%s_",
		EnvPrefix,
		strings.Join(strings.Split(strings.ToUpper(scope), "."), "_"),
	)
	if err := kEnv.Load(env.Provider(prefix, ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, prefix))
	}), nil); err != nil {
		logger.Fatalln("load from env failed:", err)
	}
	if err := k.Merge(kEnv); err != nil {
		logger.Fatalln("merge from env failed:", err)
	}

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{
		Tag: "yaml",
	}); err != nil {
		logger.Fatalln("unmarshal failed:", err)
	}

	var value atomic.Value
	value.Store(conf)
	return &ScopedConfig[T]{
		logger: logger,
		scope:  scope,
		kEnv:   kEnv,
		value:  &value,
	}
}
