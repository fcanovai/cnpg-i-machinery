/*
Copyright The CloudNativePG Contributors

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

package clusterstatus

import (
	"encoding/json"
	"fmt"

	"github.com/cloudnative-pg/cnpg-i/pkg/operator"
)

// NotAnObjectError is used when the passed value cannot be represented
// as a JSON object.
type NotAnObjectError struct {
	representation []byte
}

func (err NotAnObjectError) Error() string {
	return fmt.Sprintf(
		"the passed variable cannot be serialized as a JSON object: %s",
		err.representation,
	)
}

// SetClusterStatusResponseBuilder a SetStatus response builder.
type SetClusterStatusResponseBuilder struct{}

// NewSetClusterStatusResponseBuilder is an helper that creates the SetStatus endpoint responses.
func NewSetClusterStatusResponseBuilder() *SetClusterStatusResponseBuilder {
	return &SetClusterStatusResponseBuilder{}
}

// NoOpResponse this response will ensure that no changes will be done to the plugin status.
func (s SetClusterStatusResponseBuilder) NoOpResponse() *operator.SetClusterStatusResponse {
	return &operator.SetClusterStatusResponse{JsonStatus: nil}
}

// SetEmptyStatusResponse will set the plugin status to an empty object '{}'.
func (s SetClusterStatusResponseBuilder) SetEmptyStatusResponse() *operator.SetClusterStatusResponse {
	b, err := json.Marshal(map[string]string{})
	if err != nil {
		panic("JSON mashaller failed for empty map")
	}

	return &operator.SetClusterStatusResponse{JsonStatus: b}
}

// JSONStatusResponse requires a struct or map that can be translated to a JSON object,
// will set the status to the passed object.
func (s SetClusterStatusResponseBuilder) JSONStatusResponse(obj any) (*operator.SetClusterStatusResponse, error) {
	if obj == nil {
		return nil, ErrNilObject
	}

	jsonObject, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("while marshalling resource definition: %w", err)
	}

	var js map[string]interface{}
	if err := json.Unmarshal(jsonObject, &js); err != nil {
		return nil, NotAnObjectError{representation: jsonObject}
	}

	return &operator.SetClusterStatusResponse{
		JsonStatus: jsonObject,
	}, nil
}