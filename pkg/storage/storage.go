/*
 * Copyright (c) 2020 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package storage is used to handle persistence.
package storage

// Storage interface is used to store data.
type Storage interface {
	// InitForFirstTime initialises the storage for the first time.
	InitForFirstTime(data []byte, conf map[string]string) error

	// Init initialises the storage.
	// After "InitForFirstTime" method, this function can be called to initialize the storage in the consequent uses.
	Init(conf map[string]string) error

	// Load loads the storage as a byte array.
	Load() ([]byte, error)

	// Store store data the storage.
	Store(data []byte) error
}

// Factory struct holds Storages.
type Factory struct {
	ID string
}

// Storage method returns Storage.
func (f *Factory) Storage() Storage {
	switch f.ID {
	case FileStorageID:
		return &File{}
	default:
		return &File{}
	}
}
