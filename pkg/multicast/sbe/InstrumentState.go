package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type InstrumentStateEnum uint8
type InstrumentStateValues struct {
	Created   InstrumentStateEnum
	Open      InstrumentStateEnum
	Closed    InstrumentStateEnum
	Settled   InstrumentStateEnum
	NullValue InstrumentStateEnum
}

var InstrumentState = InstrumentStateValues{0, 1, 2, 3, 255}

func (i *InstrumentStateEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(i)); err != nil {
		return err
	}
	return nil
}

func (i InstrumentStateEnum) RangeCheck() error {
	value := reflect.ValueOf(InstrumentState)
	for idx := 0; idx < value.NumField(); idx++ {
		if i == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("range check failed on InstrumentState, unknown enumeration value %d", i)
}
