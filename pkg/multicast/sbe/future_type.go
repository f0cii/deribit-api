package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type FutureTypeEnum uint8
type FutureTypeValues struct {
	NotApplicable FutureTypeEnum
	Reversed      FutureTypeEnum
	Linear        FutureTypeEnum
	NullValue     FutureTypeEnum
}

var FutureType = FutureTypeValues{0, 1, 2, 255}

func (f *FutureTypeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(f)); err != nil {
		return err
	}
	return nil
}

func (f FutureTypeEnum) RangeCheck() error {
	value := reflect.ValueOf(FutureType)
	for idx := 0; idx < value.NumField(); idx++ {
		if f == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("%w on FutureType, unknown enumeration value %d", ErrRangeCheck, f)
}
