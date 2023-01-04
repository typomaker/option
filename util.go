package option

import "time"

type (
	Time       = Option[time.Time]
	Duration   = Option[time.Duration]
	Bool       = Option[bool]
	String     = Option[string]
	Int        = Option[int]
	Int8       = Option[int8]
	Int16      = Option[int16]
	Int32      = Option[int32]
	Int64      = Option[int64]
	Uint       = Option[uint]
	Uint8      = Option[uint8]
	Uint16     = Option[uint16]
	Uint32     = Option[uint32]
	Uint64     = Option[uint64]
	Float32    = Option[float32]
	Float64    = Option[float64]
	Complex64  = Option[complex64]
	Complex128 = Option[complex128]
)
