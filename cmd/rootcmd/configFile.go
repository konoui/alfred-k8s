package rootcmd

import (
	"time"

	"github.com/konoui/go-alfred"
	"github.com/spf13/viper"
)

// Config configuration
type configFile struct {
	Kubectl         Kubectl `mapstructure:"kubectl"`
	CacheTimeSecond int     `mapstructure:"cache_time_second"`
	KeyMaps         KeyMaps `mapstructure:"key_maps"`
}

// Kubectl configuration kubectl and plugin path
type Kubectl struct {
	Bin         string   `mapstructure:"kubectl_absolute_path"`
	PluginPaths []string `mapstructure:"plugin_paths"`
}

type KeyMaps struct {
	ContextKeyMap    KeyMap `mapstructure:"context_key_map"`
	NamespaceKeyMap  KeyMap `mapstructure:"namespace_key_map"`
	PodKeyMap        KeyMap `mapstructure:"pod_key_map"`
	DeploymentKeyMap KeyMap `mapstructure:"deployment_key_map"`
	ServiceKeyMap    KeyMap `mapstructure:"service_key_map"`
}

type KeyMap struct {
	Enter KeyMapKey `mapstructure:"enter"`
	Shift KeyMapKey `mapstructure:"shift"`
	Ctrl  KeyMapKey `mapstructure:"ctrl"`
	Cmd   KeyMapKey `mapstructure:"cmd"`
	Alt   KeyMapKey `mapstructure:"alt"`
}

type KeyMapKey string

const (
	CopyResourceKey    = "copy"
	CopySternKey       = "stern_copy"
	UseResourceKey     = "use"
	DeleteResourceKey  = "delete"
	CopyPortForwardKey = "port_forward_copy"
	ExecPortForwardKey = "port_forward_exec"
)

const defaulCacheValue = 70

var (
	defaultContextKeyMap = KeyMap{
		Enter: CopyResourceKey,
		Ctrl:  UseResourceKey,
		Shift: DeleteResourceKey,
	}
	defaultNamespaceKeyMap = KeyMap{
		Enter: CopyResourceKey,
		Ctrl:  UseResourceKey,
	}
	defaultPodKeyMap = KeyMap{
		Enter: CopyResourceKey,
		Ctrl:  DeleteResourceKey,
		Shift: CopySternKey,
		Alt:   CopyPortForwardKey,
	}
	defaultDeploymentKeyMap = KeyMap{
		Enter: CopyResourceKey,
		Alt:   CopyPortForwardKey,
		Shift: CopySternKey,
	}
	defaultServiceKeyMap = KeyMap{
		Enter: CopyResourceKey,
		Alt:   CopyPortForwardKey,
		Shift: CopySternKey,
	}
)

// New return alfred k8s configuration
func newConfigFile() (*configFile, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(".alfred-k8s")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/")

	v.SetDefault("kubectl.bin", "/usr/local/bin/kubectl")
	v.SetDefault("kubectl.plugin_paths", []string{"/usr/local/bin/"})
	v.SetDefault("cache_time_second", defaulCacheValue)
	v.SetDefault("key_maps.context_key_map", defaultContextKeyMap)
	v.SetDefault("key_maps.namespace_key_map", defaultNamespaceKeyMap)
	v.SetDefault("key_maps.pod_key_map", defaultPodKeyMap)
	v.SetDefault("key_maps.deployment_key_map", defaultDeploymentKeyMap)
	v.SetDefault("key_maps.service_key_map", defaultServiceKeyMap)
	if err := v.ReadInConfig(); err != nil {
		// ignore not found error. try to exec default options
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return &configFile{}, err
		}
	}

	// care no config file case here
	var c configFile
	if err := v.Unmarshal(&c); err != nil {
		return &configFile{}, err
	}

	return &c, nil
}

func (cfgFile *configFile) cacheTTL() time.Duration {
	var cacheTime time.Duration
	maxAge := cfgFile.CacheTimeSecond
	switch {
	case maxAge == 0:
		cacheTime = defaulCacheValue * time.Second
	case maxAge < 0:
		cacheTime = 0 * time.Second
	default:
		cacheTime = time.Duration(maxAge) * time.Second
	}
	return cacheTime
}

func MakeMods(km *KeyMap, modMap map[KeyMapKey]*alfred.Mod) (enterMod *alfred.Mod, mods map[alfred.ModKey]*alfred.Mod) {
	mods = make(map[alfred.ModKey]*alfred.Mod)
	enterMod = new(alfred.Mod)
	mod, ok := modMap[km.Enter]
	if ok {
		enterMod = mod
	}

	mod, ok = modMap[km.Alt]
	if ok {
		mods[alfred.ModAlt] = mod
	}

	mod, ok = modMap[km.Shift]
	if ok {
		mods[alfred.ModShift] = mod
	}

	mod, ok = modMap[km.Ctrl]
	if ok {
		mods[alfred.ModCtrl] = mod
	}

	mod, ok = modMap[km.Cmd]
	if ok {
		mods[alfred.ModCmd] = mod
	}

	return
}
