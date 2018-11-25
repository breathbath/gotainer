package container

import (
	"errors"
	"reflect"
)

func RegisterParameters(c Container, dependenciesMaps ...interface{}) error {
	for _, currMap := range dependenciesMaps {
		reflectedMap := reflect.ValueOf(currMap)
		if reflectedMap.Kind() != reflect.Map {
			return errors.New("A map type should be provided to register parameters")
		}

		for _, mapKey := range reflectedMap.MapKeys() {
			if mapKey.Kind() != reflect.String {
				return errors.New("A map[string]interface{} should be provided to register parameters")
			}
			serviceName := mapKey.Interface().(string)

			elem := reflectedMap.MapIndex(mapKey)
			if !elem.IsValid() {
				continue
			}

			c.AddConstructor(serviceName, func(c Container) (interface{}, error) {
				return elem.Interface(), nil
			})
		}
	}
	return nil
}
