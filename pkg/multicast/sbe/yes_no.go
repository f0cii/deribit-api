package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type YesNoEnum uint8
type YesNoValues struct {
	No        YesNoEnum
	Yes       YesNoEnum
	NullValue YesNoEnum
}

var YesNo = YesNoValues{0, 1, 255}

func (y *YesNoEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(y)); err != nil {
		return err
	}
	return nil
}

func (y YesNoEnum) RangeCheck() error {
	value := reflect.ValueOf(YesNo)
	for idx := 0; idx < value.NumField(); idx++ {
		if y == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("%w on YesNo, unknown enumeration value %d", ErrRangeCheck, y)
}
