package uuid

import (
	"database/sql/driver"
	"github.com/satori/go.uuid"
)

type UUID struct {
	uuid.UUID
}

var NilUUID = UUID{UUID: uuid.Nil}

func NewUUID() UUID {
	return UUID{UUID: uuid.NewV4()}
}

func (u UUID) String() string {
	return u.UUID.String()
}

// FromString returns UUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (UUID, error) {
	u, err := uuid.FromString(input)
	if err != nil {
		return NilUUID, err
	}

	return UUID{UUID: u}, nil
}

type NullUUID struct {
	UUID  UUID
	Valid bool
	uuid  uuid.NullUUID
}

func (u *NullUUID) Scan(src interface{}) error {
	return u.uuid.Scan(src)
}

func (u NullUUID) Value() (driver.Value, error) {
	return u.uuid.Value()
}
