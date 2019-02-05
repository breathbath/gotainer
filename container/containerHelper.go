package container

import (
	"errors"
	"reflect"
)

//RegisterParameters adds scalar parameter values as dependencies
func RegisterParameters(c Container, dependenciesMaps ...interface{}) error {
	errs := []error{}
	for _, currMap := range dependenciesMaps {
		reflectedMap := reflect.ValueOf(currMap)
		if reflectedMap.Kind() != reflect.Map {
			errs = append(errs, errors.New("A map type should be provided to register parameters"))
			continue
		}

		for _, mapKey := range reflectedMap.MapKeys() {
			if mapKey.Kind() != reflect.String {
				errs = append(
					errs,
					errors.New("A map[string]interface{} should be provided to register parameters"),
				)
				continue
			}
			serviceName := mapKey.Interface().(string)

			elem := reflectedMap.MapIndex(mapKey)
			if !elem.IsValid() {
				continue
			}

			err := c.AddConstructor(serviceName, func(c Container) (interface{}, error) {
				return elem.Interface(), nil
			})
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return mergeErrors(errs)
}
