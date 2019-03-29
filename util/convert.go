package util

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

func ConvertObjToMap(obj interface{}, tag string, out map[string]interface{}) error {
	ot := reflect.TypeOf(obj).Elem()
	vt := reflect.ValueOf(obj).Elem()
	for i := 0; i < ot.NumField(); i++ {
		v := vt.Field(i)
		key := ot.Field(i).Tag.Get(tag)
		var value interface{}
		switch v.Kind() {
		case reflect.Uint8:
			value = v.Uint()
			break
		case reflect.Uint16:
			value = v.Uint()
			break
		case reflect.Uint32:
			value = v.Uint()
			break
		case reflect.Uint:
			value = v.Uint()
			break
		case reflect.Uint64:
			value = v.Uint()
			break
		case reflect.Float32:
			value = v.Float()
			break
		case reflect.Float64:
			value = v.Float()
			break
		case reflect.Bool:
			value = v.Bool()
			break
		case reflect.Int8:
			value = v.Int()
			break
		case reflect.Int16:
			value = v.Int()
			break
		case reflect.Int32:
			value = v.Int()
			break
		case reflect.Int:
			value = v.Int()
			break
		case reflect.Int64:
			value = v.Int()
			break
		case reflect.String:
			value = v.String()
			break
		default:
			return errors.New("Has Error Type:" + ot.Field(i).Name)
		}
		out[key] = value
	}
	return nil
}

func ConvertObjectToUrlValues(obj interface{}, tag string, out url.Values) error {
	ot := reflect.TypeOf(obj).Elem()
	vt := reflect.ValueOf(obj).Elem()
	for i := 0; i < ot.NumField(); i++ {
		v := vt.Field(i)
		key := ot.Field(i).Tag.Get(tag)
		value := ""
		switch v.Kind() {
		case reflect.Uint8:
			value = strconv.FormatUint(v.Uint(), 10)
			break
		case reflect.Uint16:
			value = strconv.FormatUint(v.Uint(), 10)
			break
		case reflect.Uint32:
			value = strconv.FormatUint(v.Uint(), 10)
			break
		case reflect.Uint:
			value = strconv.FormatUint(v.Uint(), 10)
			break
		case reflect.Uint64:
			value = strconv.FormatUint(v.Uint(), 10)
			break
		case reflect.Float32:
			value = strconv.FormatFloat(v.Float(), 'g', 32, 32)
			break
		case reflect.Float64:
			value = strconv.FormatFloat(v.Float(), 'g', 64, 64)
			break
		case reflect.Bool:
			value = strconv.FormatBool(v.Bool())
			break
		case reflect.Int8:
			value = strconv.FormatInt(v.Int(), 10)
			break
		case reflect.Int16:
			value = strconv.FormatInt(v.Int(), 10)
			break
		case reflect.Int32:
			value = strconv.FormatInt(v.Int(), 10)
			break
		case reflect.Int:
			value = strconv.FormatInt(v.Int(), 10)
			break
		case reflect.Int64:
			value = strconv.FormatInt(v.Int(), 10)
			break
		case reflect.String:
			value = v.String()
			break
		default:
			return errors.New("Has Error Type:" + ot.Field(i).Name)
		}
		out.Add(key, value)
	}
	return nil
}

func ConvertObjectToUrlencodedStr(obj interface{}, tag string) (string, error) {
	v := url.Values{}
	err := ConvertObjectToUrlValues(obj, tag, v)
	if nil != err {
		return "", err
	}
	return v.Encode(), err
}

func ConvertUrlValuesToObject(val url.Values, tag string, obj interface{}) error {
	ot := reflect.TypeOf(obj).Elem()
	vt := reflect.ValueOf(obj).Elem()
	for i := 0; i < ot.NumField(); i++ {
		v := vt.Field(i)
		key := ot.Field(i).Tag.Get(tag)
		value := val.Get(key)
		if "" == value {
			continue
		}
		switch v.Kind() {
		case reflect.Uint8:
			ui, err := strconv.ParseUint(value, 10, 8)
			if nil != err {
				return err
			}
			v.SetUint(ui)
			break
		case reflect.Uint16:
			ui, err := strconv.ParseUint(value, 10, 16)
			if nil != err {
				return err
			}
			v.SetUint(ui)
			break
		case reflect.Uint32:
			ui, err := strconv.ParseUint(value, 10, 32)
			if nil != err {
				return err
			}
			v.SetUint(ui)
			break
		case reflect.Uint:
			ui, err := strconv.ParseUint(value, 10, 64)
			if nil != err {
				return err
			}
			v.SetUint(ui)
			break
		case reflect.Uint64:
			ui, err := strconv.ParseUint(value, 10, 64)
			if nil != err {
				return err
			}
			v.SetUint(ui)
			break
		case reflect.Float32:
			f, err := strconv.ParseFloat(value, 32)
			if nil != err {
				return err
			}
			v.SetFloat(f)
			break
		case reflect.Float64:
			f, err := strconv.ParseFloat(value, 64)
			if nil != err {
				return err
			}
			v.SetFloat(f)
			break
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if nil != err {
				return err
			}
			v.SetBool(b)
			break
		case reflect.Int8:
			i, err := strconv.ParseInt(value, 10, 8)
			if nil != err {
				return err
			}
			v.SetInt(i)
			break
		case reflect.Int16:
			i, err := strconv.ParseInt(value, 10, 16)
			if nil != err {
				return err
			}
			v.SetInt(i)
			break
		case reflect.Int32:
			i, err := strconv.ParseInt(value, 10, 32)
			if nil != err {
				return err
			}
			v.SetInt(i)
			break
		case reflect.Int:
			i, err := strconv.ParseInt(value, 10, 64)
			if nil != err {
				return err
			}
			v.SetInt(i)
			break
		case reflect.Int64:
			i, err := strconv.ParseInt(value, 10, 64)
			if nil != err {
				return err
			}
			v.SetInt(i)
			break
		case reflect.String:
			v.SetString(value)
			break
		default:
			return errors.New("Has Error Type:" + ot.Field(i).Name)
		}
	}
	return nil
}

func ConvertUrlencodedStrToObject(str string, tag string, obj interface{}) error {
	v, err := url.ParseQuery(str)
	if nil != err {
		return err
	}
	return ConvertUrlValuesToObject(v, tag, obj)
}
