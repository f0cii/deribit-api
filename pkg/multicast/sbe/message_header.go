package sbe

import (
	"fmt"
	"io"
	"math"
)

type MessageHeader struct {
	BlockLength      uint16
	TemplateId       uint16
	SchemaId         uint16
	Version          uint16
	NumGroups        uint16
	NumVarDataFields uint16
}

func (m *MessageHeader) Decode(_m *SbeGoMarshaller, _r io.Reader) error {
	if err := _m.ReadUint16(_r, &m.BlockLength); err != nil {
		return err
	}

	if err := _m.ReadUint16(_r, &m.TemplateId); err != nil {
		return err
	}

	if err := _m.ReadUint16(_r, &m.SchemaId); err != nil {
		return err
	}

	if err := _m.ReadUint16(_r, &m.Version); err != nil {
		return err
	}

	if err := _m.ReadUint16(_r, &m.NumGroups); err != nil {
		return err
	}

	if err := _m.ReadUint16(_r, &m.NumVarDataFields); err != nil {
		return err
	}

	return nil
}

func (m *MessageHeader) RangeCheck() error {
	if m.BlockLength < m.BlockLengthMinValue() || m.BlockLength > m.BlockLengthMaxValue() {
		return fmt.Errorf("%w on m.BlockLength (%v < %v > %v)", ErrRangeCheck, m.BlockLengthMinValue(), m.BlockLength, m.BlockLengthMaxValue())
	}

	if m.TemplateId < m.TemplateIdMinValue() || m.TemplateId > m.TemplateIdMaxValue() {
		return fmt.Errorf("%w on m.TemplateId (%v < %v > %v)", ErrRangeCheck, m.TemplateIdMinValue(), m.TemplateId, m.TemplateIdMaxValue())
	}

	if m.SchemaId < m.SchemaIdMinValue() || m.SchemaId > m.SchemaIdMaxValue() {
		return fmt.Errorf("%w on m.SchemaId (%v < %v > %v)", ErrRangeCheck, m.SchemaIdMinValue(), m.SchemaId, m.SchemaIdMaxValue())
	}

	if m.Version < m.VersionMinValue() || m.Version > m.VersionMaxValue() {
		return fmt.Errorf("%w on m.Version (%v < %v > %v)", ErrRangeCheck, m.VersionMinValue(), m.Version, m.VersionMaxValue())
	}

	if m.NumGroups < m.NumGroupsMinValue() || m.NumGroups > m.NumGroupsMaxValue() {
		return fmt.Errorf("%w on m.NumGroups (%v < %v > %v)", ErrRangeCheck, m.NumGroupsMinValue(), m.NumGroups, m.NumGroupsMaxValue())
	}

	if m.NumVarDataFields < m.NumVarDataFieldsMinValue() || m.NumVarDataFields > m.NumVarDataFieldsMaxValue() {
		return fmt.Errorf("%w on m.NumVarDataFields (%v < %v > %v)", ErrRangeCheck, m.NumVarDataFieldsMinValue(), m.NumVarDataFields, m.NumVarDataFieldsMaxValue())
	}

	return nil
}

func (*MessageHeader) BlockLengthMinValue() uint16 {
	return 0
}

func (*MessageHeader) BlockLengthMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*MessageHeader) TemplateIdMinValue() uint16 {
	return 0
}

func (*MessageHeader) TemplateIdMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*MessageHeader) SchemaIdMinValue() uint16 {
	return 0
}

func (*MessageHeader) SchemaIdMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*MessageHeader) VersionMinValue() uint16 {
	return 0
}

func (*MessageHeader) VersionMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*MessageHeader) NumGroupsMinValue() uint16 {
	return 0
}

func (*MessageHeader) NumGroupsMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*MessageHeader) NumVarDataFieldsMinValue() uint16 {
	return 0
}

func (*MessageHeader) NumVarDataFieldsMaxValue() uint16 {
	return math.MaxUint16 - 1
}
