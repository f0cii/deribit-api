package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type OptionTypeEnum uint8
type OptionTypeValues struct {
	NotApplicable OptionTypeEnum
	Put           OptionTypeEnum
	Call          OptionTypeEnum
	NullValue     OptionTypeEnum
}

var OptionType = OptionTypeValues{0, 1, 2, 255}

func (o *OptionTypeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(o)); err != nil {
		return err
	}
	return nil
}

func (o OptionTypeEnum) RangeCheck() error {
	value := reflect.ValueOf(OptionType)
	for idx := 0; idx < value.NumField(); idx++ {
		if o == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on OptionType, unknown enumeration value %d", o)
}
