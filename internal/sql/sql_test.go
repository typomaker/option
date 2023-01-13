package sql_test

import (
	"math"
	"testing"
	"time"

	sqlpkg "database/sql"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option/internal/sql"
)

func TestUnmarshal(t *testing.T) {
	t.Run("scannable", func(t *testing.T) {
		var src, dst = true, sqlpkg.NullBool{}
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, sqlpkg.NullBool{Bool: true, Valid: true}, dst)
	})
	t.Run("nil-pointer", func(t *testing.T) {
		require.Error(t, sql.Unmarshal("123", nil))
	})
	t.Run("non-pointer", func(t *testing.T) {
		var dst any
		require.Error(t, sql.Unmarshal("123", dst))
	})
	t.Run("bool to bool", func(t *testing.T) {
		var src, dst = true, false
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, true, dst)
	})
	t.Run("bool to string", func(t *testing.T) {
		var src, dst = true, ""
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, "true", dst)
	})
	t.Run("bool to int", func(t *testing.T) {
		var src, dst = true, int(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int(1), dst)
	})
	t.Run("bool to int8", func(t *testing.T) {
		var src, dst = true, int8(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int8(1), dst)
	})
	t.Run("bool to int16", func(t *testing.T) {
		var src, dst = true, int16(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int16(1), dst)
	})
	t.Run("bool to int32", func(t *testing.T) {
		var src, dst = true, int32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int32(1), dst)
	})
	t.Run("bool to int64", func(t *testing.T) {
		var src, dst = true, int64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int64(1), dst)
	})
	t.Run("bool to uint", func(t *testing.T) {
		var src, dst = true, uint(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint(1), dst)
	})
	t.Run("bool to uint8", func(t *testing.T) {
		var src, dst = true, uint8(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint8(1), dst)
	})
	t.Run("bool to uint16", func(t *testing.T) {
		var src, dst = true, uint16(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint16(1), dst)
	})
	t.Run("bool to uint32", func(t *testing.T) {
		var src, dst = true, uint32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint32(1), dst)
	})
	t.Run("bool to uint64", func(t *testing.T) {
		var src, dst = true, uint64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint64(1), dst)
	})
	t.Run("bool to float32", func(t *testing.T) {
		var src, dst = true, float32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, float32(1), dst)
	})
	t.Run("bool to float64", func(t *testing.T) {
		var src, dst = true, float64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, float64(1), dst)
	})
	t.Run("string to string", func(t *testing.T) {
		var src, dst = "123", ""
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, "123", dst)
	})
	t.Run("string to int", func(t *testing.T) {
		var src, dst = "1", int(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int(1), dst)
	})
	t.Run("string to int8", func(t *testing.T) {
		var src, dst = "1", int8(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int8(1), dst)
	})
	t.Run("string to int16", func(t *testing.T) {
		var src, dst = "1", int16(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int16(1), dst)
	})
	t.Run("string to int32", func(t *testing.T) {
		var src, dst = "1", int32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int32(1), dst)
	})
	t.Run("string to int64", func(t *testing.T) {
		var src, dst = "1", int64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int64(1), dst)
	})
	t.Run("string to uint", func(t *testing.T) {
		var src, dst = "1", uint(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint(1), dst)
	})
	t.Run("string to uint8", func(t *testing.T) {
		var src, dst = "1", uint8(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint8(1), dst)
	})
	t.Run("string to uint16", func(t *testing.T) {
		var src, dst = "1", uint16(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint16(1), dst)
	})
	t.Run("string to uint32", func(t *testing.T) {
		var src, dst = "1", uint32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint32(1), dst)
	})
	t.Run("string to uint64", func(t *testing.T) {
		var src, dst = "1", uint64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint64(1), dst)
	})
	t.Run("string to float32", func(t *testing.T) {
		var src, dst = "1.1", float32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, float32(1.1), dst)
	})
	t.Run("string to float64", func(t *testing.T) {
		var src, dst = "1.1", float64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, float64(1.1), dst)
	})
	t.Run("string to time.Duration", func(t *testing.T) {
		var src, dst = "1h", time.Duration(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, time.Hour, dst)
	})
	t.Run("int64 to int", func(t *testing.T) {
		var src, dst = int64(1), int(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int(1), dst)
	})
	t.Run("int64 to int8", func(t *testing.T) {
		var src, dst = int64(1), int8(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int8(1), dst)
	})
	t.Run("int64 to int16", func(t *testing.T) {
		var src, dst = int64(1), int16(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int16(1), dst)
	})
	t.Run("int64 to int32", func(t *testing.T) {
		var src, dst = int64(1), int32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int32(1), dst)
	})
	t.Run("int64 to int64", func(t *testing.T) {
		var src, dst = int64(1), int64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, int64(1), dst)
	})
	t.Run("int64 to uint", func(t *testing.T) {
		var src, dst = int64(1), uint(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint(1), dst)
	})
	t.Run("int64 to uint8", func(t *testing.T) {
		var src, dst = int64(1), uint8(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint8(1), dst)
	})
	t.Run("int64 to uint16", func(t *testing.T) {
		var src, dst = int64(1), uint16(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint16(1), dst)
	})
	t.Run("int64 to uint32", func(t *testing.T) {
		var src, dst = int64(1), uint32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint32(1), dst)
	})
	t.Run("int64 to uint64", func(t *testing.T) {
		var src, dst = int64(1), uint64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, uint64(1), dst)
	})
	t.Run("[]byte to []byte", func(t *testing.T) {
		var src, dst = []byte("1"), []byte{}
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, []byte("1"), dst)
	})
	t.Run("[]byte to string", func(t *testing.T) {
		var src, dst = []byte("1"), string("")
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, "1", dst)
	})
	t.Run("[]byte to encoding.BinaryUnmarshaler", func(t *testing.T) {
		var src, dst = []byte{1, 0, 0, 0, 14, 194, 139, 231, 112, 0, 0, 0, 0, 0, 0}, time.Time{}
		var err = sql.Unmarshal(src, &dst)

		require.NoError(t, err)
		require.Equal(t, "2009-11-10T23:00:00Z", dst.Format(time.RFC3339))
	})
	t.Run("time.Time to string", func(t *testing.T) {
		var src, dst = time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), string("")
		var err = sql.Unmarshal(src, &dst)

		require.NoError(t, err)
		require.Equal(t, "2009-11-10T23:00:00Z", dst)
	})
	t.Run("time.Time to time.Time", func(t *testing.T) {
		var src, dst = time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), time.Time{}
		var err = sql.Unmarshal(src, &dst)

		require.NoError(t, err)
		require.Equal(t, src, dst)
	})
	t.Run("float64 to float64", func(t *testing.T) {
		var src, dst = float64(1), float64(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, float64(1), dst)
	})
	t.Run("float64 to float32", func(t *testing.T) {
		var src, dst = float64(math.MaxFloat32), float32(0)
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, float32(math.MaxFloat32), dst)
	})
	t.Run("float64 to string", func(t *testing.T) {
		var src, dst = float64(1.111111111111), string("")
		var err = sql.Unmarshal(src, &dst)
		require.NoError(t, err)
		require.Equal(t, string("1.111111111111"), dst)
	})
}
