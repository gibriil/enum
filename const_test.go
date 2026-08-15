package enum

import (
	"database/sql/driver"
	"encoding"
)

var (
	_ Enum                   = shipping{}
	_ encoding.TextMarshaler = shipping{}
	_ driver.Valuer          = (*shipping)(nil)
)
