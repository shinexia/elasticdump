/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package mapping

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// cleanupMapping clean up some unnecessary field
func CleanUpMapping(data string) (string, error) {
	var rootMap map[string]json.RawMessage
	err := json.Unmarshal([]byte(data), &rootMap)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if len(rootMap) != 1 {
		return "", errors.Errorf("multiple indexes: %v", len(rootMap))
	}
	var indexData json.RawMessage
	for _, v := range rootMap {
		indexData = v
	}
	var dataMap map[string]json.RawMessage
	err = json.Unmarshal(indexData, &dataMap)
	if err != nil {
		return "", errors.WithStack(err)
	}
	
	//Handle possible lack of settings object in the mapping.
	if settingsData, ok := dataMap["settings"]; ok  {
		var settingsMap map[string]json.RawMessage
		
		err = json.Unmarshal(settingsData, &settingsMap)
		if err != nil {
			return "", errors.WithStack(err)
		}
		
		indexSettingData := settingsMap["index"]
		var indexSettingMap map[string]json.RawMessage
		err = json.Unmarshal(indexSettingData, &indexSettingMap)

		if err != nil {
			return "", errors.WithStack(err)
		}

		// delete .settings.index unused fields
		for _, key := range []string{"creation_date", "uuid", "version", "provided_name", "routing", "creation_date_string"} {
			delete(indexSettingMap, key)
		}

		newIndexSettingData, err := json.Marshal(indexSettingMap)

		if err != nil {
			return "", errors.WithStack(err)
		}

		settingsMap["index"] = json.RawMessage(newIndexSettingData)
		dataMap["settings"], err = json.Marshal(settingsMap)

		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	newData, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(newData), nil
}
