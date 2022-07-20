package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type BookSideEnum uint8
type BookSideValues struct {
	Ask       BookSideEnum
	Bid       BookSideEnum
	NullValue BookSideEnum
}

var BookSide = BookSideValues{0, 1, 255}

func (b *BookSideEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(b)); err != nil {
		return err
	}
	return nil
}

func (b BookSideEnum) RangeCheck() error {
	value := reflect.ValueOf(BookSide)
	for idx := 0; idx < value.NumField(); idx++ {
		if b == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on BookSide, unknown enumeration value %d", b)
}
