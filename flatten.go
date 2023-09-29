package flatten_json

import (
	"reflect"
	"strconv"
)

type Options struct {
	Delimiter string
	Prefix    string
}

func Flatten(src interface{}, dest map[string]interface{}, opt ...Options) {
	var options Options
	if len(opt) > 0 {
		options = Options{
			Delimiter: opt[0].Delimiter,
			Prefix:    opt[0].Prefix,
		}
	}

	switch t := src.(type) {
	case map[string][]string:
		for k, v := range t {
			for i, val := range v {
				del := options.Delimiter + strconv.Itoa(i)
				if len(opt) == 0 {
					del = "[" + strconv.Itoa(i) + "]"
				}

				dest[options.Prefix+k+del] = val
			}
		}
	case map[string]interface{}:
		for k, v := range t {
			switch reflect.ValueOf(v).Kind() {
			case reflect.Map:
				if val, ok := v.(map[string]interface{}); ok {
					Flatten(val, dest)
				}
			case reflect.Slice:
				switch child := v.(type) {
				case []interface{}:
					for i, val := range child {
						del := options.Delimiter + strconv.Itoa(i)
						if len(opt) == 0 {
							del = "[" + strconv.Itoa(i) + "]"
						}

						dest[options.Prefix+k+del] = val
					}
				case []string:
					for i, val := range child {
						del := options.Delimiter + strconv.Itoa(i)
						if len(opt) == 0 {
							del = "[" + strconv.Itoa(i) + "]"
						}

						dest[options.Prefix+k+del] = val
					}
				}
			case reflect.Struct:
				val := reflect.ValueOf(v)
				for i := 0; i < val.NumField(); i++ {
					del := options.Delimiter
					if len(opt) == 0 {
						del = "."
					}

					dest[options.Prefix+k+del+val.Type().Field(i).Name] = val.Field(i)
				}
			default:
				dest[options.Prefix+k] = v
			}
		}
	default:
		return
	}

}
