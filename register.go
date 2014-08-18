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

import "sync"

// register is a thread-safe map
type register struct {
	mapper map[string]interface{}
	mu     sync.RWMutex
}

// NewRegister returns a new register.
func NewRegister() *register {
	return &register{
		mapper: make(map[string]interface{}),
	}
}

// Register binds the interface and the name. If this name has been registerd, it panics.
func (r *register) Register(name string, v interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, dup := r.mapper[name]; dup {
		panic("register name: " + name + " twice")
	}
	r.mapper[name] = v
}

// Unregister unbinds the interface and the name. It returns the interface or nil
func (r *register) Unregister(name string) interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.mapper[name]; ok {
		delete(r.mapper, name)
		return v
	}
	return nil
}

// Get return an interface registerd with this name.
func (r *register) Get(name string) interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mapper[name]
}
