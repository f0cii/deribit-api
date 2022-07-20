package sbe

import (
	"fmt"
	"io"
	"math"
)

type GroupSizeEncoding struct {
	BlockLength      uint16
	NumInGroup       uint16
	NumGroups        uint16
	NumVarDataFields uint16
}

func (g *GroupSizeEncoding) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint16(_r, &g.BlockLength); err != nil {
		return err
	}

	if err := _m.ReadUint16(_r, &g.NumInGroup); err != nil {
		return err
	}

	if err := _m.ReadUint16(_r, &g.NumGroups); err != nil {
		return err
	}

	if err := _m.ReadUint16(_r, &g.NumVarDataFields); err != nil {
		return err
	}

	return nil
}

func (g *GroupSizeEncoding) RangeCheck() error {
	if g.BlockLength < g.BlockLengthMinValue() || g.BlockLength > g.BlockLengthMaxValue() {
		return fmt.Errorf("Range check failed on g.BlockLength (%v < %v > %v)", g.BlockLengthMinValue(), g.BlockLength, g.BlockLengthMaxValue())
	}

	if g.NumInGroup < g.NumInGroupMinValue() || g.NumInGroup > g.NumInGroupMaxValue() {
		return fmt.Errorf("Range check failed on g.NumInGroup (%v < %v > %v)", g.NumInGroupMinValue(), g.NumInGroup, g.NumInGroupMaxValue())
	}

	if g.NumGroups < g.NumGroupsMinValue() || g.NumGroups > g.NumGroupsMaxValue() {
		return fmt.Errorf("Range check failed on g.NumGroups (%v < %v > %v)", g.NumGroupsMinValue(), g.NumGroups, g.NumGroupsMaxValue())
	}

	if g.NumVarDataFields < g.NumVarDataFieldsMinValue() || g.NumVarDataFields > g.NumVarDataFieldsMaxValue() {
		return fmt.Errorf("Range check failed on g.NumVarDataFields (%v < %v > %v)", g.NumVarDataFieldsMinValue(), g.NumVarDataFields, g.NumVarDataFieldsMaxValue())
	}

	return nil
}

func (*GroupSizeEncoding) BlockLengthMinValue() uint16 {
	return 0
}

func (*GroupSizeEncoding) BlockLengthMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*GroupSizeEncoding) NumInGroupMinValue() uint16 {
	return 0
}

func (*GroupSizeEncoding) NumInGroupMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*GroupSizeEncoding) NumGroupsMinValue() uint16 {
	return 0
}

func (*GroupSizeEncoding) NumGroupsMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*GroupSizeEncoding) NumVarDataFieldsMinValue() uint16 {
	return 0
}

func (*GroupSizeEncoding) NumVarDataFieldsMaxValue() uint16 {
	return math.MaxUint16 - 1
}
