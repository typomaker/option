package option

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/lib/pq"
)

func (it *Option[T]) Scan(value any) (err error) {
	if value == nil {
		*it = None[T]()
		return nil
	}
	if v, ok := any(it.Get()).(sql.Scanner); ok {
		if err := v.Scan(value); err != nil {
			return err
		}
		*it = Some(v.(T))
		return nil
	}
	if v, ok := value.(T); ok {
		*it = Some(v)
		return nil
	}
	switch dst := any(&it.value).(type) {
	case *bool:
		var t sql.NullBool
		if err = t.Scan(value); err == nil {
			*it = Some(any(t.Bool).(T))
		}
		return err
	case *byte:
		var t sql.NullByte
		if err = t.Scan(value); err == nil {
			*it = Some(any(t.Byte).(T))
		}
		return err
	case *float64:
		var t sql.NullFloat64
		if err = t.Scan(value); err == nil {
			*it = Some(any(t.Float64).(T))
		}
		return err
	case *int16:
		var t sql.NullInt16
		if err = t.Scan(value); err == nil {
			*it = Some(any(t.Int16).(T))
		}
		return err
	case *int32:
		var t sql.NullInt32
		if err = t.Scan(value); err == nil {
			*it = Some(any(t.Int32).(T))
		}
		return err
	case *int64:
		var t sql.NullInt64
		if err = t.Scan(value); err == nil {
			*it = Some(any(t.Int64).(T))
		}
		return err
	case *string:
		var t sql.NullString
		if err = t.Scan(value); err == nil {
			*it = Some(any(t.String).(T))
		}
		return err
	case *time.Time:
		var t sql.NullTime
		if err = t.Scan(value); err == nil {
			*it = Some(any(t.Value).(T))
		}
		return err
	case *[]byte, *[]float64, *[]float32, *[]int64, *[]int32, *[]string, *[][]byte, *[]bool:
		if err = pq.Array(dst).Scan(value); err == nil {
			it.notZero, it.notNone = true, true
		}
		return err
	}
	var zero T
	return fmt.Errorf("can't scan %T to %T type", value, zero)
}
func (it *Option[T]) Value() (driver.Value, error) {
	if !it.IsSome() {
		return nil, nil
	}
	if v, ok := any(it.Get()).(driver.Valuer); ok {
		return v.Value()
	}
	switch src := any(it.value).(type) {
	case bool:
		var t = sql.NullBool{Bool: src, Valid: true}
		return t.Value()
	case byte:
		var t = sql.NullByte{Byte: src, Valid: true}
		return t.Value()
	case float64:
		var t = sql.NullFloat64{Float64: src, Valid: true}
		return t.Value()
	case int16:
		var t = sql.NullInt16{Int16: src, Valid: true}
		return t.Value()
	case int32:
		var t = sql.NullInt32{Int32: src, Valid: true}
		return t.Value()
	case int64:
		var t = sql.NullInt64{Int64: src, Valid: true}
		return t.Value()
	case string:
		var t = sql.NullString{String: src, Valid: true}
		return t.Value()
	case time.Time:
		var t = sql.NullTime{Time: src, Valid: true}
		return t.Value()
	case []byte, []float64, []float32, []int64, []int32, []string, [][]byte, []bool:
		return pq.Array(src).Value()
	}
	var zero T
	return nil, fmt.Errorf("can't encode value of %T", zero)
}
