package common

import (
	"math"
	"reflect"
)

func ReplaceNaNValueOfStruct(v interface{}, typeOfV reflect.Type) {
	LogP := reflect.ValueOf(v)
	if LogP.CanConvert(typeOfV) {
		LogP = LogP.Convert(typeOfV)
	}
	LogV := LogP.Elem()

	if LogV.Kind() == reflect.Struct {
		for i := 0; i < LogV.NumField(); i++ {
			field := LogV.Field(i)
			kind := field.Kind()

			if kind == reflect.Float64 {
				if field.IsValid() && field.CanSet() && math.IsNaN(field.Float()) {
					field.SetFloat(0)
				}
			}
		}
	}
}
