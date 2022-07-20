package sbe

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type Snapshot struct {
	InstrumentId   uint32
	TimestampMs    uint64
	ChangeId       uint64
	IsBookComplete YesNoEnum
	IsLastInBook   YesNoEnum
	LevelsList     []SnapshotLevelsList
}
type SnapshotLevelsList struct {
	Side   BookSideEnum
	Price  float64
	Amount float64
}

func (s *Snapshot) Decode(_m *SbeGoMarshaller, _r io.Reader, blockLength uint16, doRangeCheck bool) error {
	if err := _m.ReadUint32(_r, &s.InstrumentId); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &s.TimestampMs); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &s.ChangeId); err != nil {
		return err
	}

	if err := s.IsBookComplete.Decode(_m, _r); err != nil {
		return err
	}

	if err := s.IsLastInBook.Decode(_m, _r); err != nil {
		return err
	}

	if blockLength > s.SbeBlockLength() {
		_, _ = io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	var LevelsListBlockLength uint16
	if err := _m.ReadUint16(_r, &LevelsListBlockLength); err != nil {
		return err
	}
	var LevelsListNumInGroup uint16
	if err := _m.ReadUint16(_r, &LevelsListNumInGroup); err != nil {
		return err
	}

	// Discard numGroups and numVars.
	_, _ = io.CopyN(ioutil.Discard, _r, 4)

	if cap(s.LevelsList) < int(LevelsListNumInGroup) {
		s.LevelsList = make([]SnapshotLevelsList, LevelsListNumInGroup)
	}
	s.LevelsList = s.LevelsList[:LevelsListNumInGroup]
	for i := range s.LevelsList {
		if err := s.LevelsList[i].Decode(_m, _r, uint(LevelsListBlockLength)); err != nil {
			return err
		}
	}

	if doRangeCheck {
		if err := s.RangeCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Snapshot) RangeCheck() error {
	if s.InstrumentId < s.InstrumentIdMinValue() || s.InstrumentId > s.InstrumentIdMaxValue() {
		return fmt.Errorf("Range check failed on s.InstrumentId (%v < %v > %v)", s.InstrumentIdMinValue(), s.InstrumentId, s.InstrumentIdMaxValue())
	}

	if s.TimestampMs < s.TimestampMsMinValue() || s.TimestampMs > s.TimestampMsMaxValue() {
		return fmt.Errorf("Range check failed on s.TimestampMs (%v < %v > %v)", s.TimestampMsMinValue(), s.TimestampMs, s.TimestampMsMaxValue())
	}

	if s.ChangeId < s.ChangeIdMinValue() || s.ChangeId > s.ChangeIdMaxValue() {
		return fmt.Errorf("Range check failed on s.ChangeId (%v < %v > %v)", s.ChangeIdMinValue(), s.ChangeId, s.ChangeIdMaxValue())
	}

	if err := s.IsBookComplete.RangeCheck(); err != nil {
		return err
	}
	if err := s.IsLastInBook.RangeCheck(); err != nil {
		return err
	}
	for _, prop := range s.LevelsList {
		if err := prop.RangeCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapshotLevelsList) Decode(_m *SbeGoMarshaller, _r io.Reader, blockLength uint) error {
	if err := s.Side.Decode(_m, _r); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &s.Price); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &s.Amount); err != nil {
		return err
	}

	if blockLength > s.SbeBlockLength() {
		_, _ = io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	return nil
}

func (s *SnapshotLevelsList) RangeCheck() error {
	if err := s.Side.RangeCheck(); err != nil {
		return err
	}
	if s.Price < s.PriceMinValue() || s.Price > s.PriceMaxValue() {
		return fmt.Errorf("Range check failed on s.Price (%v < %v > %v)", s.PriceMinValue(), s.Price, s.PriceMaxValue())
	}

	if s.Amount < s.AmountMinValue() || s.Amount > s.AmountMaxValue() {
		return fmt.Errorf("Range check failed on s.Amount (%v < %v > %v)", s.AmountMinValue(), s.Amount, s.AmountMaxValue())
	}

	return nil
}

func (*Snapshot) SbeBlockLength() (blockLength uint16) {
	return 34
}

func (*Snapshot) InstrumentIdMinValue() uint32 {
	return 0
}

func (*Snapshot) InstrumentIdMaxValue() uint32 {
	return math.MaxUint32 - 1
}

func (*Snapshot) TimestampMsMinValue() uint64 {
	return 0
}

func (*Snapshot) TimestampMsMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*Snapshot) ChangeIdMinValue() uint64 {
	return 0
}

func (*Snapshot) ChangeIdMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*SnapshotLevelsList) PriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*SnapshotLevelsList) PriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*SnapshotLevelsList) AmountMinValue() float64 {
	return -math.MaxFloat64
}

func (*SnapshotLevelsList) AmountMaxValue() float64 {
	return math.MaxFloat64
}

func (*SnapshotLevelsList) SbeBlockLength() (blockLength uint) {
	return 17
}
