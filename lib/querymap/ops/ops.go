package ops

import (
	"fmt"
	"reflect"
)

func floatify(val interface{}) (float64, error) {
	floatType := reflect.TypeOf(float64(0))
	v := reflect.ValueOf(val)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func GT(val, check interface{}) (bool, error) {
	a, err := floatify(val)
	if err != nil {
		return false, err
	}
	b, err := floatify(check)
	if err != nil {
		return false, err
	}
	return a > b, nil
}

func LT(val, check interface{}) (bool, error) {
	a, err := floatify(val)
	if err != nil {
		return false, err
	}
	b, err := floatify(check)
	if err != nil {
		return false, err
	}
	return a < b, nil
}

func GTE(val, check interface{}) (bool, error) {
	a, err := floatify(val)
	if err != nil {
		return false, err
	}
	b, err := floatify(check)
	if err != nil {
		return false, err
	}
	return a >= b, nil
}

func LTE(val, check interface{}) (bool, error) {
	a, err := floatify(val)
	if err != nil {
		return false, err
	}
	b, err := floatify(check)
	if err != nil {
		return false, err
	}
	return a <= b, nil
}

func EQ(val, check interface{}) (bool, error) {
	a, err := floatify(val)
	if err != nil {
		return false, err
	}
	b, err := floatify(check)
	if err != nil {
		return false, err
	}
	return a == b, nil
}
