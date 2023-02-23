package aconfig

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ghodss/yaml"
	"github.com/spf13/cast"
)

type Interface struct {
	Name string `json:"name,omitempty"`
}

type Location struct {
	Address string `json:"address,omitempty"`
	Port    int    `json:"port,omitempty"`
}

type PrimaryPort struct {
	Interface Interface `json:"interface,omitempty"`
	Location  *Location `json:"location,omitempty"`
}

type MatchedPrimaryPort struct {
	Location        *Location `json:"location,omitempty"`
	ApplicationName string    `json:"application_name,omitempty"`
}

type SecondaryPort struct {
	Interface          Interface           `json:"interface,omitempty"`
	Options            KV                  `json:"options,omitempty"`
	MatchedPrimaryPort *MatchedPrimaryPort `json:"matched_primary_port,omitempty"`
}

type PortName string

type ApplicationSpec struct {
	PrimaryPorts   map[PortName]PrimaryPort   `json:"primary_ports,omitempty"`
	SecondaryPorts map[PortName]SecondaryPort `json:"secondary_ports,omitempty"`
	CustomConfig   KV                         `json:"custom_config,omitempty"`
}

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"api_version,omitempty"`
}

type ObjectMeta struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type Application struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`
	Spec       ApplicationSpec `json:"spec,omitempty"`
}

type PreProcessFunc func([]byte) ([]byte, error)

func New(configFile string, funcs ...PreProcessFunc) (*Application, error) {
	var application = &Application{}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	// Run the funcs on it
	for _, f := range funcs {
		if data, err = f(data); err != nil {
			return nil, err
		}
	}
	err = yaml.Unmarshal(data, application)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (a *Application) GetName() string {
	return a.ObjectMeta.Name
}

func (a *Application) GetAPIVersion() string {
	return a.APIVersion
}

func (a *Application) GetSecondaryPorts() map[PortName]SecondaryPort {
	return a.Spec.SecondaryPorts
}

func (a *Application) GetSecondaryPort(name PortName) *SecondaryPort {
	if a.Spec.SecondaryPorts == nil {
		return nil
	}
	if sp, exists := a.Spec.SecondaryPorts[name]; exists {
		return &sp
	}

	return nil
}

func (a *Application) GetMatchedPrimaryPort(name PortName) *MatchedPrimaryPort {
	sp := a.GetSecondaryPort(name)
	if sp == nil {
		return nil
	}

	return sp.MatchedPrimaryPort
}

func (a *Application) GetMatchedPrimaryPortLocation(name PortName) *Location {
	p := a.GetMatchedPrimaryPort(name)
	if p == nil {
		return nil
	}

	return p.Location
}

func (a *Application) GetCustomConfig() KV {
	return a.Spec.CustomConfig
}

type KV map[string]interface{}

func (kv KV) LoadTo(out interface{}) error {
	jsonStr, err := json.Marshal(kv)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(jsonStr, out); err != nil {
		return err
	}

	return nil
}

func (kv KV) get(key string) interface{} {
	return kv[key]
}

func (kv KV) Get(key string) interface{} {
	return kv.get(key)
}

func (kv KV) GetString(key string) string {
	return cast.ToString(kv.get(key))
}

func (kv KV) GetBool(key string) bool {
	return cast.ToBool(kv.get(key))
}

func (kv KV) GetDuration(key string) time.Duration {
	return cast.ToDuration(kv.get(key))
}

func (kv KV) GetFloat64(key string) float64 {
	return cast.ToFloat64(kv.get(key))
}

func (kv KV) GetInt(key string) int {
	return cast.ToInt(kv.get(key))
}

func (kv KV) GetInt32(key string) int32 {
	return cast.ToInt32(kv.get(key))
}

func (kv KV) GetInt64(key string) int64 {
	return cast.ToInt64(kv.get(key))
}

func (kv KV) GetUint(key string) uint {
	return cast.ToUint(kv.get(key))
}

func (kv KV) GetUint32(key string) uint32 {
	return cast.ToUint32(kv.get(key))
}

func (kv KV) GetUint64(key string) uint64 {
	return cast.ToUint64(kv.get(key))
}

func (kv KV) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(kv.get(key))
}

func (kv KV) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(kv.get(key))
}

func (kv KV) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(kv.get(key))
}

func (kv KV) GetStringSlice(key string) []string {
	return cast.ToStringSlice(kv.get(key))
}

func (kv KV) GetTime(key string) time.Time {
	return cast.ToTime(kv.get(key))
}
