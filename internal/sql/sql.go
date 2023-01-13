package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type (
	Scanner = sql.Scanner
	Valuer  = driver.Valuer
	Value   = driver.Value
)

func Marshal(goval any) (dbval Value, err error) {
	if s, ok := goval.(Valuer); ok {
		return s.Value()
	}
	return goval, nil
}

func Unmarshal(srcval Value, dstval any) (err error) {
	if xdst, ok := dstval.(Scanner); ok {
		return xdst.Scan(srcval)
	}

	rv := reflect.ValueOf(dstval)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("sql: goval must be non-nil pointer")
	}
	var f64 float64
	var i64 int64
	var i64u uint64
	switch src := srcval.(type) {
	case bool:
		switch ref := dstval.(type) {
		case *bool:
			*ref = src
			return
		case *string:
			*ref = strconv.FormatBool(src)
			return
		case *int:
			if src {
				*ref = 1
			}
			return
		case *int8:
			if src {
				*ref = 1
			}
			return
		case *int16:
			if src {
				*ref = 1
			}
			return
		case *int32:
			if src {
				*ref = 1
			}
			return
		case *int64:
			if src {
				*ref = 1
			}
			return
		case *uint:
			if src {
				*ref = 1
			}
			return
		case *uint8:
			if src {
				*ref = 1
			}
			return
		case *uint16:
			if src {
				*ref = 1
			}
			return
		case *uint32:
			if src {
				*ref = 1
			}
			return
		case *uint64:
			if src {
				*ref = 1
			}
			return
		case *float32:
			if src {
				*ref = 1
			}
			return
		case *float64:
			if src {
				*ref = 1
			}
			return
		}
	case string:
		switch ref := any(dstval).(type) {
		case *string:
			*ref = src
			return
		case *int:
			*ref, err = strconv.Atoi(src)
			return
		case *int8:
			if i64, err = strconv.ParseInt(src, 10, 8); err == nil {
				*ref = int8(i64)
			}
			return
		case *int16:
			if i64, err = strconv.ParseInt(src, 10, 16); err == nil {
				*ref = int16(i64)
			}
			return
		case *int32:
			if i64, err = strconv.ParseInt(src, 10, 32); err == nil {
				*ref = int32(i64)
			}
			return
		case *int64:
			*ref, err = strconv.ParseInt(src, 10, 64)
			return
		case *uint:
			if i64u, err = strconv.ParseUint(src, 10, 32); err == nil {
				*ref = uint(i64u)
			}
			return
		case *uint8:
			if i64u, err = strconv.ParseUint(src, 10, 8); err == nil {
				*ref = uint8(i64u)
			}
			return
		case *uint16:
			if i64u, err = strconv.ParseUint(src, 10, 16); err == nil {
				*ref = uint16(i64u)
			}
			return
		case *uint32:
			if i64u, err = strconv.ParseUint(src, 10, 32); err == nil {
				*ref = uint32(i64u)
			}
			return
		case *uint64:
			*ref, err = strconv.ParseUint(src, 10, 64)
			return
		case *float32:
			if f64, err = strconv.ParseFloat(src, 64); err == nil {
				*ref = float32(f64)
			}
			return
		case *float64:
			*ref, err = strconv.ParseFloat(src, 64)
			return
		case *time.Duration:
			*ref, err = time.ParseDuration(src)
			return
		}
	case int64:
		switch ref := any(dstval).(type) {
		case *int:
			*ref = int(src)
			return
		case *int8:
			*ref = int8(src)
			return
		case *int16:
			*ref = int16(src)
			return
		case *int32:
			*ref = int32(src)
			return
		case *int64:
			*ref = src
			return
		case *uint:
			*ref = uint(src)
			return
		case *uint8:
			*ref = uint8(src)
			return
		case *uint16:
			*ref = uint16(src)
			return
		case *uint32:
			*ref = uint32(src)
			return
		case *uint64:
			*ref = uint64(src)
			return
		}
	case []byte:
		switch ref := any(dstval).(type) {
		case encoding.BinaryUnmarshaler:
			err = ref.UnmarshalBinary(src)
			return
		case *[]byte:
			*ref = src
			return
		case *string:
			*ref = string(src)
			return
		}
	case float64:
		switch ref := any(dstval).(type) {
		case *float64:
			*ref = src
			return
		case *float32:
			*ref = float32(src)
			return
		case *string:
			*ref = strconv.FormatFloat(src, 'f', -1, 64)
			return
		}
	case time.Time:
		switch ref := any(dstval).(type) {
		case *string:
			*ref = src.Format(time.RFC3339)
			return
		case *time.Time:
			*ref = src
			return
		}
	}
	err = fmt.Errorf("sql: unmarshal from %T to %T is not implemented", srcval, dstval)
	return
}
