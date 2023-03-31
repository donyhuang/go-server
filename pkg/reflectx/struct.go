package reflectx

import (
	"errors"
	"github.com/donyhuang/go-server/pkg/float"
	"reflect"
)

var (
	UnSupportTypeError   = errors.New("only support pointer")
	FieldTypeUnSupport   = errors.New("field only support number string bool")
	ErrorNeedPointerType = errors.New("set type need pointer struct")
	ErrorTypeNotStruct   = errors.New("type all need struct")
	ErrorTwoTypeNotEqual = errors.New("two type is not equal")
)

func ResetStructPointer(v interface{}) error {
	rValue := reflect.ValueOf(v)
	if rValue.Kind() != reflect.Pointer {
		return UnSupportTypeError
	}
	rValue = rValue.Elem()
	rValue.Set(reflect.Zero(rValue.Type()))
	return nil
}

func SetStructNotEmpty(l, r interface{}) error {
	rValue := reflect.ValueOf(r)
	lValue := reflect.ValueOf(l)
	if lValue.Kind() != reflect.Pointer {
		return ErrorNeedPointerType
	}
	lValue = lValue.Elem()
	if rValue.Kind() == reflect.Pointer {
		rValue = rValue.Elem()
	}
	if lValue.Kind() != reflect.Struct || rValue.Kind() != reflect.Struct {
		return ErrorTypeNotStruct
	}
	if lValue.Type().Name() != rValue.Type().Name() {
		return ErrorTwoTypeNotEqual
	}
	for i := 0; i < rValue.NumField(); i++ {
		if !rValue.Field(i).IsZero() && lValue.Field(i).CanSet() {
			lValue.Field(i).Set(rValue.Field(i))
		}
		lValue.Type()
	}
	return nil
}

func TruncStructField(s interface{}, fields []string, decimal int) {
	rValue := reflect.ValueOf(s)
	if rValue.Kind() != reflect.Pointer {
		return
	}
	rValue = rValue.Elem()
	for _, field := range fields {
		fieldValue := rValue.FieldByName(field)
		if fieldValue.CanFloat() {
			fieldValue.SetFloat(float.TruncFloat(fieldValue.Float(), decimal))
		}
	}
}
