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

import "encoding/json"

type HandlerBuilderFunc func() Handler
type FormatterBuilderFunc func() Formatter
type FilterBuilderFunc func() Filter

var handlerBuilderMap map[string]HandlerBuilderFunc
var formatterBuilderMap map[string]FormatterBuilderFunc
var filterBuilderMap map[string]FilterBuilderFunc

type ConfigLoader interface {
	LoadConfig(config []byte)
	LoadConfigFromMap(m map[string]interface{})
	LoadConfigFromFile(fileName string)
}

func RegisterFilterBuilder(name string, builder HandlerBuilderFunc) {

}

func RegisterFormatterBuilder(name string, builder HandlerBuilderFunc) {

}

func RegisterHandlerBuilder(name string, builder HandlerBuilderFunc) {

}

func GetFilterBuilder(name string) FilterBuilderFunc {
	return nil
}

func GetFormatterBuilder(name string) FormatterBuilderFunc {
	return nil
}

func GetHandlerBuilder(name string) HandlerBuilderFunc {
	return nil
}

func LoadConfig(config []byte) {
	var mapConfig map[string]map[string]map[string]interface{}
	err := json.Unmarshal(config, &mapConfig)
	if err != nil {
		panic(err)
	}

	filters, ok := mapConfig["filters"]
	if ok {
		for _, conf := range filters {
			t, yes := conf["type"]
			if !yes {
				panic("Filter config must have a 'type' field")
			}
			builder, yes := filterBuilderMap[t.(string)]
			if !yes {
				panic("No Filter builder for " + t.(string))
			}
			filter := builder()
			filter.LoadConfigFromMap(conf)
		}
	}

	formatters, ok := mapConfig["formatters"]
	if ok {
		for _, conf := range formatters {
			t, yes := conf["type"]
			if !yes {
				panic("Formatter config must have a 'type' field")
			}
			builder, yes := formatterBuilderMap[t.(string)]
			if !yes {
				panic("No Formatter builder for " + t.(string))
			}
			formatter := builder()
			formatter.LoadConfigFromMap(conf)
		}
	}

	handlers, ok := mapConfig["handlers"]
	if ok {
		for _, conf := range handlers {
			t, yes := conf["type"]
			if !yes {
				panic("Handlers config must have a 'type' field")
			}
			builder, yes := handlerBuilderMap[t.(string)]
			if !yes {
				panic("No Handler builder for " + t.(string))
			}
			handler := builder()
			handler.LoadConfigFromMap(conf)
		}
	}

	loggers, ok := mapConfig["loggers"]
	if ok {
		for name, conf := range loggers {
			logger := new(gLogger)
			logger.LoadConfigFromMap(conf)
			logger.SetName(name)
			RegisterLogger(logger)
		}
	}
}
