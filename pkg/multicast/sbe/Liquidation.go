package sbe

import (
	"fmt"
	"io"
	"reflect"
)

type LiquidationEnum uint8
type LiquidationValues struct {
	None      LiquidationEnum
	Maker     LiquidationEnum
	Taker     LiquidationEnum
	Both      LiquidationEnum
	NullValue LiquidationEnum
}

var Liquidation = LiquidationValues{0, 1, 2, 3, 255}

func (l *LiquidationEnum) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint8(_r, (*uint8)(l)); err != nil {
		return err
	}
	return nil
}

func (l LiquidationEnum) RangeCheck() error {
	value := reflect.ValueOf(Liquidation)
	for idx := 0; idx < value.NumField(); idx++ {
		if l == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on Liquidation, unknown enumeration value %d", l)
}
