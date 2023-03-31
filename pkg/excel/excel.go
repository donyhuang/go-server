package excel

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"reflect"
)

const (
	Tag          = "excel"
	DefaultSheet = "Sheet1"
)

var (
	ErrorType = errors.New("only slice struct")
)

// SaveStructToBuff 结构体保存为xlsx 格式 []byte  struct tag 为excel
func SaveStructToBuff(s interface{}) ([]byte, error) {
	sTyp := reflect.TypeOf(s)
	sVal := reflect.ValueOf(s)
	if sTyp.Kind() == reflect.Pointer {
		sTyp = sTyp.Elem()
	}
	if sTyp.Kind() != reflect.Slice {
		return nil, ErrorType
	}
	sTyp = sTyp.Elem()
	if sTyp.Kind() == reflect.Pointer {
		sTyp = sTyp.Elem()
	}
	if sTyp.Kind() != reflect.Struct {
		return nil, ErrorType
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		return nil, err
	}
	fields, tags := getFieldTitleFromStruct(sTyp)
	err = f.SetSheetRow(DefaultSheet, cell, &tags)
	if err != nil {
		return nil, err
	}
	sVal = reflect.Indirect(sVal)
	for i := 0; i < sVal.Len(); i++ {
		row := reflect.Indirect(sVal.Index(i))
		cell, err = excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			return nil, nil
		}
		rows := make([]interface{}, len(fields), len(fields))
		for k, field := range fields {
			rows[k] = row.FieldByName(field).Interface()
		}
		err = f.SetSheetRow(DefaultSheet, cell, &rows)
		if err != nil {
			return nil, err
		}

	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, nil
	}
	return buf.Bytes(), nil
}

func getFieldTitleFromStruct(s reflect.Type) ([]string, []string) {
	var fields []string
	var tags []string
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get(Tag)
		if tag == "" || tag == "-" {
			continue
		}
		fields = append(fields, field.Name)
		tags = append(tags, tag)
	}
	return fields, tags
}
