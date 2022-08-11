package sbe

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type Trades struct {
	InstrumentId uint32
	TradesList   []TradesTradesList
}

type TradesTradesList struct {
	Direction     DirectionEnum
	Price         float64
	Amount        float64
	TimestampMs   uint64
	MarkPrice     float64
	IndexPrice    float64
	TradeSeq      uint64
	TradeId       uint64
	TickDirection TickDirectionEnum
	Liquidation   LiquidationEnum
	Iv            float64
	BlockTradeId  uint64
	ComboTradeId  uint64
}

func (t *Trades) Decode(_m *SbeGoMarshaller, _r io.Reader, blockLength uint16, doRangeCheck bool) error {
	if err := _m.ReadUint32(_r, &t.InstrumentId); err != nil {
		return err
	}

	if blockLength > t.SbeBlockLength() {
		_, _ = io.CopyN(ioutil.Discard, _r, int64(blockLength-t.SbeBlockLength()))
	}

	var TradesListBlockLength uint16
	if err := _m.ReadUint16(_r, &TradesListBlockLength); err != nil {
		return err
	}
	var TradesListNumInGroup uint16
	if err := _m.ReadUint16(_r, &TradesListNumInGroup); err != nil {
		return err
	}

	// Discard numGroups and numVars.
	_, _ = io.CopyN(ioutil.Discard, _r, 4)

	if cap(t.TradesList) < int(TradesListNumInGroup) {
		t.TradesList = make([]TradesTradesList, TradesListNumInGroup)
	}
	t.TradesList = t.TradesList[:TradesListNumInGroup]
	for i := range t.TradesList {
		if err := t.TradesList[i].Decode(_m, _r, uint(TradesListBlockLength)); err != nil {
			return err
		}
	}

	if doRangeCheck {
		if err := t.RangeCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Trades) RangeCheck() error {
	if t.InstrumentId < t.InstrumentIdMinValue() || t.InstrumentId > t.InstrumentIdMaxValue() {
		return fmt.Errorf("range check failed on t.InstrumentId (%v < %v > %v)", t.InstrumentIdMinValue(), t.InstrumentId, t.InstrumentIdMaxValue())
	}

	for _, prop := range t.TradesList {
		if err := prop.RangeCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (t *TradesTradesList) Decode(_m *SbeGoMarshaller, _r io.Reader, blockLength uint) error {
	if err := t.Direction.Decode(_m, _r); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.Price); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.Amount); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &t.TimestampMs); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.MarkPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.IndexPrice); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &t.TradeSeq); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &t.TradeId); err != nil {
		return err
	}

	if err := t.TickDirection.Decode(_m, _r); err != nil {
		return err
	}

	if err := t.Liquidation.Decode(_m, _r); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.Iv); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &t.BlockTradeId); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &t.ComboTradeId); err != nil {
		return err
	}

	if blockLength > t.SbeBlockLength() {
		_, _ = io.CopyN(ioutil.Discard, _r, int64(blockLength-t.SbeBlockLength()))
	}

	return nil
}

func (t *TradesTradesList) RangeCheck() error {
	if err := t.Direction.RangeCheck(); err != nil {
		return err
	}

	if t.Price < t.PriceMinValue() || t.Price > t.PriceMaxValue() {
		return fmt.Errorf("range check failed on t.Price (%v < %v > %v)", t.PriceMinValue(), t.Price, t.PriceMaxValue())
	}

	if t.Amount < t.AmountMinValue() || t.Amount > t.AmountMaxValue() {
		return fmt.Errorf("range check failed on t.Amount (%v < %v > %v)", t.AmountMinValue(), t.Amount, t.AmountMaxValue())
	}

	if t.TimestampMs < t.TimestampMsMinValue() || t.TimestampMs > t.TimestampMsMaxValue() {
		return fmt.Errorf("range check failed on t.TimestampMs (%v < %v > %v)", t.TimestampMsMinValue(), t.TimestampMs, t.TimestampMsMaxValue())
	}

	if t.MarkPrice < t.MarkPriceMinValue() || t.MarkPrice > t.MarkPriceMaxValue() {
		return fmt.Errorf("range check failed on t.MarkPrice (%v < %v > %v)", t.MarkPriceMinValue(), t.MarkPrice, t.MarkPriceMaxValue())
	}

	if t.IndexPrice < t.IndexPriceMinValue() || t.IndexPrice > t.IndexPriceMaxValue() {
		return fmt.Errorf("range check failed on t.IndexPrice (%v < %v > %v)", t.IndexPriceMinValue(), t.IndexPrice, t.IndexPriceMaxValue())
	}

	if t.TradeSeq < t.TradeSeqMinValue() || t.TradeSeq > t.TradeSeqMaxValue() {
		return fmt.Errorf("range check failed on t.TradeSeq (%v < %v > %v)", t.TradeSeqMinValue(), t.TradeSeq, t.TradeSeqMaxValue())
	}

	if t.TradeId < t.TradeIdMinValue() || t.TradeId > t.TradeIdMaxValue() {
		return fmt.Errorf("range check failed on t.TradeId (%v < %v > %v)", t.TradeIdMinValue(), t.TradeId, t.TradeIdMaxValue())
	}

	if err := t.TickDirection.RangeCheck(); err != nil {
		return err
	}
	if err := t.Liquidation.RangeCheck(); err != nil {
		return err
	}

	if t.Iv != t.IvNullValue() && (t.Iv < t.IvMinValue() || t.Iv > t.IvMaxValue()) {
		return fmt.Errorf("range check failed on t.Iv (%v < %v > %v)", t.IvMinValue(), t.Iv, t.IvMaxValue())
	}

	if t.BlockTradeId != t.BlockTradeIdNullValue() && (t.BlockTradeId < t.BlockTradeIdMinValue() || t.BlockTradeId > t.BlockTradeIdMaxValue()) {
		return fmt.Errorf("range check failed on t.BlockTradeId (%v < %v > %v)", t.BlockTradeIdMinValue(), t.BlockTradeId, t.BlockTradeIdMaxValue())
	}

	if t.ComboTradeId != t.ComboTradeIdNullValue() && (t.ComboTradeId < t.ComboTradeIdMinValue() || t.ComboTradeId > t.ComboTradeIdMaxValue()) {
		return fmt.Errorf("range check failed on t.ComboTradeId (%v < %v > %v)", t.ComboTradeIdMinValue(), t.ComboTradeId, t.ComboTradeIdMaxValue())
	}

	return nil
}

func (*Trades) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*Trades) InstrumentIdMinValue() uint32 {
	return 0
}

func (*Trades) InstrumentIdMaxValue() uint32 {
	return math.MaxUint32 - 1
}

func (*TradesTradesList) PriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*TradesTradesList) PriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*TradesTradesList) AmountMinValue() float64 {
	return -math.MaxFloat64
}

func (*TradesTradesList) AmountMaxValue() float64 {
	return math.MaxFloat64
}

func (*TradesTradesList) TimestampMsMinValue() uint64 {
	return 0
}

func (*TradesTradesList) TimestampMsMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TradesTradesList) MarkPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*TradesTradesList) MarkPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*TradesTradesList) IndexPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*TradesTradesList) IndexPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*TradesTradesList) TradeSeqMinValue() uint64 {
	return 0
}

func (*TradesTradesList) TradeSeqMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TradesTradesList) TradeIdMinValue() uint64 {
	return 0
}

func (*TradesTradesList) TradeIdMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TradesTradesList) IvMinValue() float64 {
	return -math.MaxFloat64
}

func (*TradesTradesList) IvMaxValue() float64 {
	return math.MaxFloat64
}

func (*TradesTradesList) IvNullValue() float64 {
	return math.NaN()
}

func (*TradesTradesList) BlockTradeIdMinValue() uint64 {
	return 0
}

func (*TradesTradesList) BlockTradeIdMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TradesTradesList) BlockTradeIdNullValue() uint64 {
	return math.MaxUint64
}

func (*TradesTradesList) ComboTradeIdMinValue() uint64 {
	return 0
}

func (*TradesTradesList) ComboTradeIdMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TradesTradesList) ComboTradeIdNullValue() uint64 {
	return math.MaxUint64
}

func (*TradesTradesList) SbeBlockLength() (blockLength uint) {
	return 83
}
