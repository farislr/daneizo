package sqlentity

import "database/sql/driver"

type Entity interface {
	Values() []any
	Columns() []any
	StringColumns() []string
	DriverValues() []driver.Value
	MappedValues() map[string]driver.Value
}

type UpdateEntity interface {
	MappedValues() map[string]driver.Value
}
