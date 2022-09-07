package sbe

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type Book struct {
	InstrumentId uint32
	TimestampMs  uint64
	PrevChangeId uint64
	ChangeId     uint64
	IsLast       YesNoEnum
	ChangesList  []BookChangesList
}

type BookChangesList struct {
	Side   BookSideEnum
	Change BookChangeEnum
	Price  float64
	Amount float64
}

func (b *Book) Decode(_m *SbeGoMarshaller, _r io.Reader, blockLength uint16, doRangeCheck bool) error {
	if err := _m.ReadUint32(_r, &b.InstrumentId); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &b.TimestampMs); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &b.PrevChangeId); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &b.ChangeId); err != nil {
		return err
	}

	if err := b.IsLast.Decode(_m, _r); err != nil {
		return err
	}

	if blockLength > b.SbeBlockLength() {
		_, _ = io.CopyN(ioutil.Discard, _r, int64(blockLength-b.SbeBlockLength()))
	}

	var ChangesListBlockLength uint16
	if err := _m.ReadUint16(_r, &ChangesListBlockLength); err != nil {
		return err
	}

	var ChangesListNumInGroup uint16
	if err := _m.ReadUint16(_r, &ChangesListNumInGroup); err != nil {
		return err
	}

	// Discard numGroups and numVars.
	_, _ = io.CopyN(ioutil.Discard, _r, 4)

	if cap(b.ChangesList) < int(ChangesListNumInGroup) {
		b.ChangesList = make([]BookChangesList, ChangesListNumInGroup)
	}
	b.ChangesList = b.ChangesList[:ChangesListNumInGroup]
	for i := range b.ChangesList {
		if err := b.ChangesList[i].Decode(_m, _r, uint(ChangesListBlockLength)); err != nil {
			return err
		}
	}

	if doRangeCheck {
		if err := b.RangeCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Book) RangeCheck() error {
	if b.InstrumentId < b.InstrumentIdMinValue() || b.InstrumentId > b.InstrumentIdMaxValue() {
		return fmt.Errorf("%w on b.InstrumentId (%v < %v > %v)", ErrRangeCheck, b.InstrumentIdMinValue(), b.InstrumentId, b.InstrumentIdMaxValue())
	}

	if b.TimestampMs < b.TimestampMsMinValue() || b.TimestampMs > b.TimestampMsMaxValue() {
		return fmt.Errorf("%w on b.TimestampMs (%v < %v > %v)", ErrRangeCheck, b.TimestampMsMinValue(), b.TimestampMs, b.TimestampMsMaxValue())
	}

	if b.PrevChangeId < b.PrevChangeIdMinValue() || b.PrevChangeId > b.PrevChangeIdMaxValue() {
		return fmt.Errorf("%w on b.PrevChangeId (%v < %v > %v)", ErrRangeCheck, b.PrevChangeIdMinValue(), b.PrevChangeId, b.PrevChangeIdMaxValue())
	}

	if b.ChangeId < b.ChangeIdMinValue() || b.ChangeId > b.ChangeIdMaxValue() {
		return fmt.Errorf("%w on b.ChangeId (%v < %v > %v)", ErrRangeCheck, b.ChangeIdMinValue(), b.ChangeId, b.ChangeIdMaxValue())
	}

	if err := b.IsLast.RangeCheck(); err != nil {
		return err
	}
	for _, prop := range b.ChangesList {
		if err := prop.RangeCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (b *BookChangesList) Decode(_m *SbeGoMarshaller, _r io.Reader, blockLength uint) error {
	if err := b.Side.Decode(_m, _r); err != nil {
		return err
	}

	if err := b.Change.Decode(_m, _r); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &b.Price); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &b.Amount); err != nil {
		return err
	}

	if blockLength > b.SbeBlockLength() {
		_, _ = io.CopyN(ioutil.Discard, _r, int64(blockLength-b.SbeBlockLength()))
	}

	return nil
}

func (b *BookChangesList) RangeCheck() error {
	if err := b.Side.RangeCheck(); err != nil {
		return err
	}

	if err := b.Change.RangeCheck(); err != nil {
		return err
	}

	if b.Price < b.PriceMinValue() || b.Price > b.PriceMaxValue() {
		return fmt.Errorf("%w on b.Price (%v < %v > %v)", ErrRangeCheck, b.PriceMinValue(), b.Price, b.PriceMaxValue())
	}

	if b.Amount < b.AmountMinValue() || b.Amount > b.AmountMaxValue() {
		return fmt.Errorf("%w on b.Amount (%v < %v > %v)", ErrRangeCheck, b.AmountMinValue(), b.Amount, b.AmountMaxValue())
	}

	return nil
}

func (*Book) SbeBlockLength() (blockLength uint16) {
	return 41 // Length of fixed fields include header
}

func (*Book) InstrumentIdMinValue() uint32 {
	return 0
}

func (*Book) InstrumentIdMaxValue() uint32 {
	return math.MaxUint32 - 1
}

func (*Book) TimestampMsMinValue() uint64 {
	return 0
}

func (*Book) TimestampMsMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*Book) PrevChangeIdMinValue() uint64 {
	return 0
}

func (*Book) PrevChangeIdMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*Book) PrevChangeIdNullValue() uint64 {
	return math.MaxUint64
}

func (*Book) ChangeIdMinValue() uint64 {
	return 0
}

func (*Book) ChangeIdMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BookChangesList) PriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*BookChangesList) PriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*BookChangesList) AmountMinValue() float64 {
	return -math.MaxFloat64
}

func (*BookChangesList) AmountMaxValue() float64 {
	return math.MaxFloat64
}

func (*BookChangesList) SbeBlockLength() (blockLength uint) {
	return 18
}
