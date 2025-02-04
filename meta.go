package errors

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

type Meta []struct {
	Key   string
	Value string
}

func (meta Meta) Len() int {
	return len(meta)
}

func (meta Meta) Less(i, j int) bool {
	return meta[i].Key < meta[j].Key
}

func (meta Meta) Swap(i, j int) {
	meta[i], meta[j] = meta[j], meta[i]
}

func WithMeta(key string, val any) Option {
	return func(o *Options) {
		if key == "" {
			return
		}
		switch v := val.(type) {
		case string:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: v})
			break
		case int:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatInt(int64(v), 10)})
			break
		case int8:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatInt(int64(v), 10)})
			break
		case int16:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatInt(int64(v), 10)})
			break
		case int32:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatInt(int64(v), 10)})
			break
		case int64:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatInt(v, 10)})
			break
		case uint:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatUint(uint64(v), 10)})
			break
		case uint16:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatUint(uint64(v), 10)})
			break
		case uint32:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatUint(uint64(v), 10)})
			break
		case uint64:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatUint(v, 10)})
			break
		case float32:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatFloat(float64(v), 'f', -1, 32)})
			break
		case float64:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatFloat(v, 'f', -1, 64)})
			break
		case bool:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: strconv.FormatBool(v)})
			break
		case byte:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: string(v)})
			break
		case []byte:
			if len(v) == 0 {
				o.Meta = append(o.Meta, struct {
					Key   string
					Value string
				}{Key: key, Value: ""})
				break
			}
			s := unsafe.String(&v[0], len(v))
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: s})
			break
		case time.Time:
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: v.Format(time.RFC3339)})
			break
		default:
			if s, ok := v.(interface{ String() string }); ok {
				o.Meta = append(o.Meta, struct {
					Key   string
					Value string
				}{Key: key, Value: s.String()})
				break
			}
			o.Meta = append(o.Meta, struct {
				Key   string
				Value string
			}{Key: key, Value: fmt.Sprintf("%v", v)})
			break
		}
	}
}
