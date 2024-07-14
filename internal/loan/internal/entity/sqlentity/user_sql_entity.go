package sqlentity

import (
	"database/sql/driver"
	"time"
)

type User struct {
	ID        uint64
	Name      string
	Type      UserType
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) Columns() []any {
	return []any{
		"id",
		"name",
		"type",
		"created_at",
		"updated_at",
	}
}

func (u User) StringColumns() []string {
	vals := make([]string, len(u.Columns()))
	for i, col := range u.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (u *User) Values() []any {
	return []any{
		u.ID,
		u.Name,
		u.Type,
		u.CreatedAt,
		u.UpdatedAt,
	}
}

func (u *User) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(u.Values()))
	for i, v := range u.Values() {
		vals[i] = v
	}

	return vals
}

type UserType int

const (
	Unknown UserType = iota
	Borrower
	Investor
	Employee
)

func (ut UserType) String() string {
	return [...]string{"Unknown", "Borrower", "Investor", "Employee"}[ut]
}

func (ut UserType) Value() (driver.Value, error) {
	return ut.String(), nil
}

func (ut UserType) getMap() map[string]UserType {
	return map[string]UserType{
		"Unknown":  Unknown,
		"Borrower": Borrower,
		"Investor": Investor,
		"Employee": Employee,
	}
}

func (ut *UserType) Scan(value any) error {
	s, ok := value.(string)
	if ok {
		val := ut.getMap()[s]

		*ut = val
	}

	return nil
}
