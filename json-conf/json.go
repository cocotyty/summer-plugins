package json_conf

import (
	"encoding/json"
	"errors"
	"github.com/cocotyty/summer"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type jsonObject map[string]interface{}

func getValue(from interface{}, k string, path string) (interface{}, error) {
	switch value := from.(type) {
	case map[string]interface{}:
		return value[k], nil
	case []interface{}:
		index, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}
		if index < len(value) {
			return value[index], nil
		}
		return nil, errors.New("out of range")
	default:
		return nil, errors.New("path error :" + path + ",not in json")
	}
}
func (jo jsonObject) find(path string) (object interface{}, err error) {
	object = map[string]interface{}(jo)
	paths := strings.Split(path, ".")
	for _, v := range paths {
		object, err = getValue(object, v, path)
		if err != nil {
			return nil, err
		}
	}
	return object, nil
}
func LoadJSON(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return LoadJSONBytes(data)
}
func LoadJSONBytes(data []byte) error {
	v := &jsonObject{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	summer.PluginRegister(SummerJSON{jsonObject: v}, summer.BeforeInit)
	return nil
}

type SummerJSON struct {
	jsonObject *jsonObject
}

// look up the value which field wanted
func (sj SummerJSON) Look(Holder *summer.Holder, path string, sf *reflect.StructField) reflect.Value {
	v, err := sj.jsonObject.find(path)
	if err != nil {
		panic(err)
	}
	if sf.Type.Kind() == reflect.Slice {
		if sf.Type.Elem().Kind() == reflect.String {
			strs := []string{}
			if list, ok := v.([]interface{}); ok {
				for _, elm := range list {
					str, ok := elm.(string)
					if !ok {
						panic("JSON is Wrong! @" + path)
					}
					strs = append(strs, str)
				}
			}
			return reflect.ValueOf(strs)
		}
		if sf.Type.Elem().Kind() == reflect.Int {
			ints := []int{}
			if list, ok := v.([]interface{}); ok {
				for _, elm := range list {
					i, ok := elm.(int)
					if !ok {
						panic("JSON is Wrong! @" + path)
					}
					ints = append(ints, i)
				}
			}
			return reflect.ValueOf(ints)
		}
	}
	if reflect.TypeOf(v).Kind() == reflect.Float64 {
		switch sf.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num := v.(float64)
			return reflect.ValueOf(int(num))
		}
	}

	return reflect.ValueOf(v)
}

// tell  summer the plugin prefix
func (sj SummerJSON) Prefix() string {
	return "json"
}

// zIndex represent the sequence of plugins
func (sj SummerJSON) ZIndex() int {
	return 0
}
