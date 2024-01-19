/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package helpers

import (
	"encoding/json"
	"fmt"
)

func ToJSON(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("%v", obj)
	}
	return string(data)
}
