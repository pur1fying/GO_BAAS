package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/pur1fying/GO_BAAS/internal/global_info"
	"github.com/pur1fying/GO_BAAS/internal/logger"
	"gopkg.in/yaml.v2"
)

var Config *GOBAASConfig

type GOBAASConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Mail     MailConfig     `yaml:"mail"`
}

func GenerateDefaultConfig() *GOBAASConfig {
	return &GOBAASConfig{
		Server:   *DefaultServerConfig(),
		Database: *DefaultDatabaseConfig(),
		Mail:     *DefaultMailConfig(),
	}
}

func Load(path string) error {
	if path == "" {
		path = global_info.GO_BAAS_DEFAULT_CONFIG_PATH
	}
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := createDefaultConfigFile(); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var oldConfigMap map[string]interface{}
	if err := yaml.Unmarshal(data, &oldConfigMap); err != nil {
		return err
	}

	logger.SubTitle("Config Load")
	logger.BAASInfo("Path:", path)
	updated := update(oldConfigMap)

	data, err = yaml.Marshal(updated)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &Config); err != nil {
		return err
	}

	return Save()
}

func Save() error {
	logger.SubTitle("Config Save")
	data, _ := yaml.Marshal(Config)
	return os.WriteFile(global_info.GO_BAAS_DEFAULT_CONFIG_PATH, data, 0644)
}

func createDefaultConfigFile() error {
	logger.SubTitle("Write Default Config")
	logger.BAASInfo("Path:", global_info.GO_BAAS_DEFAULT_CONFIG_PATH)
	cfg := GenerateDefaultConfig()
	data, _ := yaml.Marshal(cfg)
	err := os.MkdirAll(global_info.GO_BAAS_CONFIG_DIR, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(global_info.GO_BAAS_DEFAULT_CONFIG_PATH, data, 0644)
}

func update(old map[string]interface{}) map[string]interface{} {
	logger.SubTitle("Config Key Update")
	defaultMap := structToMap(GenerateDefaultConfig())
	return mergeMaps(old, defaultMap, "")
}

func structToMap(s interface{}) map[string]interface{} {
	data, _ := yaml.Marshal(s)
	var m map[string]interface{}
	_ = yaml.Unmarshal(data, &m)
	return m
}

func mergeMaps(old, defaults map[string]interface{}, path string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, defaultValue := range defaults {
		fullPath := key
		if path != "" {
			fullPath = path + "." + key
		}

		oldValue, exists := old[key]

		switch defaultVal := defaultValue.(type) {
		case map[string]interface{}:
			if exists {
				switch oldVal := oldValue.(type) {
				case map[string]interface{}:
					result[key] = mergeMaps(oldVal, defaultVal, fullPath)
				case map[interface{}]interface{}:
					result[key] = mergeMaps(convertMap(oldVal), defaultVal, fullPath)
				default:
					result[key] = defaultVal
					logValue(false, fullPath, oldValue)
					logMap(true, fullPath, defaultVal)
				}
			} else {
				result[key] = defaultVal
				logMap(true, fullPath, defaultVal)
			}

		case map[interface{}]interface{}:
			convertedDefault := convertMap(defaultVal)
			if exists {
				switch oldVal := oldValue.(type) {
				case map[string]interface{}:
					result[key] = mergeMaps(oldVal, convertedDefault, fullPath)
				case map[interface{}]interface{}:
					result[key] = mergeMaps(convertMap(oldVal), convertedDefault, fullPath)
				default:
					result[key] = convertedDefault
					logValue(false, fullPath, oldValue)
					logMap(true, fullPath, convertedDefault)
				}
			} else {
				result[key] = convertedDefault
				logMap(true, fullPath, convertedDefault)
			}

		default:
			if !exists {
				result[key] = defaultValue
				logValue(true, fullPath, defaultValue)
			}
		}
	}

	for key, oldValue := range old {
		if _, exists := defaults[key]; !exists {
			fullPath := key
			if path != "" {
				fullPath = path + "." + key
			}

			switch val := oldValue.(type) {
			case map[string]interface{}:
				logMap(false, fullPath, val)
			case map[interface{}]interface{}:
				convertedVal := convertMap(val)
				logMap(false, fullPath, convertedVal)
			default:
				logValue(false, fullPath, oldValue)
			}
		} else {
			result[key] = oldValue
		}
	}
	return result
}

func convertMap(input map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range input {
		if strKey, ok := k.(string); ok {
			if nestedMap, ok := v.(map[interface{}]interface{}); ok {
				result[strKey] = convertMap(nestedMap)
			} else if nestedMapStr, ok := v.(map[string]interface{}); ok {
				result[strKey] = nestedMapStr
			} else {
				result[strKey] = v
			}
		}
	}
	return result
}

func logValue(isAdd bool, path string, value interface{}) {
	st := "Add    :"
	if !isAdd {
		st = "Del    :"
	}
	v := fmt.Sprintf("%v", value)
	if len(v) == 0 {
		v = "\"\""
	}
	logger.BAASInfo(st, path, "=", v)
}

func logMap(isAdd bool, path string, value interface{}) {
	st := "Add Map:"
	if !isAdd {
		st = "Del Map:"
	}
	switch v := value.(type) {
	case map[string]interface{}:
		logger.BAASInfo(st, path)
		for key, val := range v {
			subPath := path + "." + key
			logMap(isAdd, subPath, val)
		}
	default:
		logValue(isAdd, path, value)
	}
}
