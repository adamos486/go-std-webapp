package sqltest

import (
	"database/sql/driver"
	"time"
)

//AnyTime as a type allows for type checking of "AnyTime" value.
type AnyTime struct{}

//AnyString as a type allows for type checking of "AnyString".
type AnyString struct{}

//Match matches AnyTime and simply type checks then allows
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

//Match matches AnyString and simply type checks on string and allows.
func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
