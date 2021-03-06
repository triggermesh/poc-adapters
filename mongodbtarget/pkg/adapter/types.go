/*
Copyright (c) 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mongodbtarget

// InsertPayload defines the expected data structure found at the event payload
type InsertPayload struct {
	Database   string            `json:"database"`
	Collection string            `json:"collection"`
	StrValue   string            `json:"strValue"`
	MapStrVal  map[string]string `json:"mapStrVal"`
}

// QueryPayload defines the expected data found at the "io.triggermesh.mongodb.query" payload
type QueryPayload struct {
	Database   string `json:"database"`
	Collection string `json:"collection"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}

type UpdatePayload struct {
	Database    string `json:"database"`
	Collection  string `json:"collection"`
	SearchKey   string `json:"searchKey"`
	SearchValue string `json:"searchValue"`
	UpdateKey   string `json:"updateKey"`
	UpdateValue string `json:"updateValue"`
}

type QueryResponse struct {
	Collection map[string]string `json:"collection"`
}
