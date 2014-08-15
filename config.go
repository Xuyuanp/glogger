/*
 * Copyright 2014 Xuyuan Pang <xuyuanp@gmail.com>
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
	"io/ioutil"
	"os"
	"sync"
)

// ConfigLoaderBuilder is a function which return a ConfigLoader.
type ConfigLoaderBuilder func() ConfigLoader

// ConfigLoader provide method to load config from bytes, string or a file.
type ConfigLoader interface {
	LoadConfig(config []byte)
	LoadConfigFromMap(m map[string]interface{})
	LoadConfigFromFile(fileName string)
}

var manager *configLoaderBuilderManager
var onceManager sync.Once

func init() {
	onceManager.Do(initManager)
}

func initManager() {
	manager = &configLoaderBuilderManager{
		mapper: make(map[string]ConfigLoaderBuilder),
	}
}

type configLoaderBuilderManager struct {
	mapper map[string]ConfigLoaderBuilder
	mu     sync.RWMutex
}

func (manager *configLoaderBuilderManager) registerConfigLoaderBuilder(name string, configLoader ConfigLoaderBuilder) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	_, dup := manager.mapper[name]
	if dup {
		panic("Duplicate ConfigLoader with name " + name)
	}
	manager.mapper[name] = configLoader
}

func (manager *configLoaderBuilderManager) getConfigLoaderBuilder(name string) ConfigLoaderBuilder {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.mapper[name]
}

// RegisterConfigLoaderBuilder register the builder whith name.
func RegisterConfigLoaderBuilder(name string, configLoader ConfigLoaderBuilder) {
	manager.registerConfigLoaderBuilder(name, configLoader)
}

// GetConfigLoaderBuilder return a ConfigLoaderBuilder registered with this name
func GetConfigLoaderBuilder(name string) ConfigLoaderBuilder {
	return manager.getConfigLoaderBuilder(name)
}

// LoadConfig parse the json format configuration.
func LoadConfig(config []byte) {
	var configMap map[string]map[string]map[string]interface{}
	err := json.Unmarshal(config, &configMap)
	if err != nil {
		panic(err)
	}

	process := func(name string, conf map[string]interface{}, callback func(loader ConfigLoader)) {
		bn, yes := conf["builder"]
		var builderName string
		if !yes {
			panic("'build' field is required for section " + name)
		}
		builderName = bn.(string)
		builder := GetConfigLoaderBuilder(builderName)
		if builder == nil {
			panic("Builder named " + builderName + " doesn't exist")
		}
		loader := builder()
		loader.LoadConfigFromMap(conf)
		callback(loader)
	}

	filters, ok := configMap["filters"]
	if ok {
		for name, conf := range filters {
			process(name, conf, func(loader ConfigLoader) {
				filter := loader.(Filter)
				RegisterFilter(name, filter)
			})
		}
	}
	formatters, ok := configMap["formatters"]
	if ok {
		for name, conf := range formatters {
			process(name, conf, func(loader ConfigLoader) {
				formatter := loader.(Formatter)
				RegisterFormatter(name, formatter)
			})
		}
	}
	handlers, ok := configMap["handlers"]
	if ok {
		for name, conf := range handlers {
			process(name, conf, func(loader ConfigLoader) {
				handler := loader.(Handler)
				handler.SetName(name)
				RegisterHandler(name, handler)
			})
		}
	}
	loggers, ok := configMap["loggers"]
	if ok {
		for name, conf := range loggers {
			logger := new(gLogger)
			logger.LoadConfigFromMap(conf)
			logger.SetName(name)
		}
	}
}

// LoadConfigFromFile read file's content and call the LoadConfig method.
func LoadConfigFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	code, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	LoadConfig(code)
}
