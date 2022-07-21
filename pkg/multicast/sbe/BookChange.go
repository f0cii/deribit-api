package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type BookChangeEnum uint8
type BookChangeValues struct {
	Created   BookChangeEnum
	Changed   BookChangeEnum
	Deleted   BookChangeEnum
	NullValue BookChangeEnum
}

var BookChange = BookChangeValues{0, 1, 2, 255}

func (b *BookChangeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(b)); err != nil {
		return err
	}
	return nil
}

func (b BookChangeEnum) RangeCheck() error {
	value := reflect.ValueOf(BookChange)
	for idx := 0; idx < value.NumField(); idx++ {
		if b == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("range check failed on BookChange, unknown enumeration value %d", b)
}
