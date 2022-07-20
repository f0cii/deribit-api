package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type InstrumentKindEnum uint8
type InstrumentKindValues struct {
	Future    InstrumentKindEnum
	Option    InstrumentKindEnum
	NullValue InstrumentKindEnum
}

var InstrumentKind = InstrumentKindValues{0, 1, 255}

func (i *InstrumentKindEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(i)); err != nil {
		return err
	}
	return nil
}

func (i InstrumentKindEnum) RangeCheck() error {
	value := reflect.ValueOf(InstrumentKind)
	for idx := 0; idx < value.NumField(); idx++ {
		if i == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on InstrumentKind, unknown enumeration value %d", i)
}
