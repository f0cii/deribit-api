package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type DirectionEnum uint8
type DirectionValues struct {
	Buy       DirectionEnum
	Sell      DirectionEnum
	NullValue DirectionEnum
}

var Direction = DirectionValues{0, 1, 255}

func (d DirectionEnum) String() string {
	switch d {
	case Direction.Buy:
		return "buy"
	case Direction.Sell:
		return "sell"
	default:
		return ""
	}
}

func (d *DirectionEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(d)); err != nil {
		return err
	}
	return nil
}

func (d DirectionEnum) RangeCheck() error {
	value := reflect.ValueOf(Direction)
	for idx := 0; idx < value.NumField(); idx++ {
		if d == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("%w on Direction, unknown enumeration value %d", ErrRangeCheck, d)
}
