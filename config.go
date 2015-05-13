/*
 * Copyright 2014 Xuyuan Pang <xuyuanp # gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package glogger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ConfigLoaderBuilder is a function which return a ConfigLoader.
type ConfigLoaderBuilder func() ConfigLoader

// ConfigLoader provide method to load config from bytes, string or a file.
type ConfigLoader interface {
	LoadConfig(m map[string]interface{}) error
}

var configLoaderBuilderRegister = NewRegister()

// RegisterConfigLoaderBuilder register the builder whith name.
func RegisterConfigLoaderBuilder(name string, builder ConfigLoaderBuilder) {
	configLoaderBuilderRegister.Register(name, builder)
}

// GetConfigLoaderBuilder return a ConfigLoaderBuilder registered with this name
func GetConfigLoaderBuilder(name string) ConfigLoaderBuilder {
	if v := configLoaderBuilderRegister.Get(name); v != nil {
		return v.(ConfigLoaderBuilder)
	}
	return nil
}

// LoadConfig parse the json format configuration.
func LoadConfig(config []byte) error {
	var configMap map[string]map[string]map[string]interface{}
	if err := json.Unmarshal(config, &configMap); err != nil {
		return err
	}

	processFunc := func(name string, conf map[string]interface{}, callback func(loader ConfigLoader)) error {
		bn, yes := conf["builder"]
		var builderName string
		if !yes {
			return fmt.Errorf("'build' field is required for section %s", name)
		}
		builderName = bn.(string)
		builder := GetConfigLoaderBuilder(builderName)
		if builder == nil {
			return fmt.Errorf("unknown builder name: %s", builderName)
		}
		loader := builder()
		if err := loader.LoadConfig(conf); err != nil {
			return err
		}
		callback(loader)
		return nil
	}

	filters, ok := configMap["filters"]
	if ok {
		for name, conf := range filters {
			if err := processFunc(name, conf, func(loader ConfigLoader) {
				filter := loader.(Filter)
				RegisterFilter(name, filter)
			}); err != nil {
				return err
			}
		}
	}
	formatters, ok := configMap["formatters"]
	if ok {
		for name, conf := range formatters {
			if builder, ok := conf["builder"]; !ok || builder.(string) == "default" {
				conf["builder"] = "github.com/Xuyuanp/glogger.DefaultFormatter"
			}
			if err := processFunc(name, conf, func(loader ConfigLoader) {
				formatter := loader.(Formatter)
				RegisterFormatter(name, formatter)
			}); err != nil {
				return err
			}
		}
	}
	handlers, ok := configMap["handlers"]
	if ok {
		for name, conf := range handlers {
			if builder, ok := conf["builder"]; !ok || builder.(string) == "default" {
				conf["builder"] = "github.com/Xuyuanp/glogger.StreamHandler"
			}
			if err := processFunc(name, conf, func(loader ConfigLoader) {
				handler := loader.(Handler)
				RegisterHandler(name, handler)
			}); err != nil {
				return err
			}
		}
	}
	loggers, ok := configMap["loggers"]
	if ok {
		for name, conf := range loggers {
			if name == "root" {
				logger := GetLogger("root")
				if err := logger.LoadConfig(conf); err != nil {
					return err
				}
				continue
			}
			logger := NewLogger()
			if err := logger.LoadConfig(conf); err != nil {
				return err
			}
			RegisterLogger(name, logger)
		}
	}
	return nil
}

// LoadConfigFromFile read file's content and call the LoadConfig method.
func LoadConfigFromFile(fileName string) error {
	var err error
	var file *os.File
	if file, err = os.Open(fileName); err != nil {
		return err
	}
	defer file.Close()
	if code, err := ioutil.ReadAll(file); err == nil {
		return LoadConfig(code)
	} else {
		return err
	}
}
