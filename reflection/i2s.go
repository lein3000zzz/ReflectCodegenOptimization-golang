package main

import (
	"errors"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	// // todo
	// return errors.New("not implemented yet")

	outValue := reflect.ValueOf(out)

	if outValue.Kind() != reflect.Ptr {
		return errors.New("out must be a pointer")
	}

	if !outValue.CanSet() && !outValue.Elem().CanSet() {
		return errors.New("out must be settable")
	}

	outElem := outValue.Elem()

	return fillValue(data, outElem)
}

func fillValue(data interface{}, outValue reflect.Value) error {
	if data == nil {
		return nil
	}

	dataValue := reflect.ValueOf(data)
	outType := outValue.Type()

	switch outValue.Kind() {
	case reflect.Struct:

		if dataValue.Kind() != reflect.Map {
			return errors.New("expected map for struct")
		}

		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return errors.New("expected map[string]interface{} for struct")
		}

		err := fillStruct(outValue, outType, dataMap)

		return err

	case reflect.Slice:

		if dataValue.Kind() != reflect.Slice {
			return errors.New("expected slice for slice")
		}

		dataSlice, ok := data.([]interface{})
		if !ok {
			return errors.New("expected []interface{} for slice")
		}

		newSlice := reflect.MakeSlice(outType, len(dataSlice), len(dataSlice))

		for i, item := range dataSlice {
			if err := fillValue(item, newSlice.Index(i)); err != nil {
				return err
			}
		}

		outValue.Set(newSlice)

	case reflect.String:

		if dataValue.Kind() != reflect.String {
			return errors.New("expected string")
		}
		outValue.SetString(dataValue.String())

	case reflect.Int:

		switch dataValue.Kind() {
		case reflect.Float64:
			outValue.SetInt(int64(dataValue.Float()))
		case reflect.Int:
			fallthrough
		case reflect.Int64:
			outValue.SetInt(dataValue.Int())
		default:
			return errors.New("expected number for int")
		}

	case reflect.Bool:

		if dataValue.Kind() != reflect.Bool {
			return errors.New("expected bool")
		}
		outValue.SetBool(dataValue.Bool())

	default:
		return errors.New("unsupported type")
	}

	return nil
}

func fillStruct(outValue reflect.Value, outType reflect.Type, dataMap map[string]interface{}) error {
	for i := 0; i < outValue.NumField(); i++ {
		field := outValue.Field(i)
		fieldType := outType.Field(i)
		fieldName := fieldType.Name

		if !field.CanSet() {
			continue
		}

		if fieldData, exists := dataMap[fieldName]; exists {
			if err := fillValue(fieldData, field); err != nil {
				return err
			}
		}
	}
	return nil
}
