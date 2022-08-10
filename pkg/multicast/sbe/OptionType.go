package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type OptionTypeEnum uint8
type OptionTypeValues struct {
	NotApplicable OptionTypeEnum
	Call          OptionTypeEnum
	Put           OptionTypeEnum
	NullValue     OptionTypeEnum
}

var OptionType = OptionTypeValues{0, 1, 2, 255}

func (o OptionTypeEnum) String() string {
	switch o {
	case OptionType.NotApplicable:
		return "not_applicable"
	case OptionType.Put:
		return "put"
	case OptionType.Call:
		return "call"
	default:
		return ""
	}
}

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
	return fmt.Errorf("range check failed on OptionType, unknown enumeration value %d", o)
}
