package sbe

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type Ticker struct {
	InstrumentId           uint32
	InstrumentState        InstrumentStateEnum
	TimestampMs            uint64
	OpenInterest           float64
	MinSellPrice           float64
	MaxBuyPrice            float64
	LastPrice              float64
	IndexPrice             float64
	MarkPrice              float64
	BestBidPrice           float64
	BestBidAmount          float64
	BestAskPrice           float64
	BestAskAmount          float64
	CurrentFunding         float64
	Funding8h              float64
	EstimatedDeliveryPrice float64
	DeliveryPrice          float64
	SettlementPrice        float64
}

func (t *Ticker) Decode(_m *SbeGoMarshaller, _r io.Reader, blockLength uint16, doRangeCheck bool) error {
	if err := _m.ReadUint32(_r, &t.InstrumentId); err != nil {
		return err
	}

	if err := t.InstrumentState.Decode(_m, _r); err != nil {
		return err
	}

	if err := _m.ReadUint64(_r, &t.TimestampMs); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.OpenInterest); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.MinSellPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.MaxBuyPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.LastPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.IndexPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.MarkPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.BestBidPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.BestBidAmount); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.BestAskPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.BestAskAmount); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.CurrentFunding); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.Funding8h); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.EstimatedDeliveryPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.DeliveryPrice); err != nil {
		return err
	}

	if err := _m.ReadFloat64(_r, &t.SettlementPrice); err != nil {
		return err
	}

	if blockLength > t.SbeBlockLength() {
		_, _ = io.CopyN(ioutil.Discard, _r, int64(blockLength-t.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := t.RangeCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Ticker) RangeCheck() error {
	if t.InstrumentId < t.InstrumentIdMinValue() || t.InstrumentId > t.InstrumentIdMaxValue() {
		return fmt.Errorf("range check failed on t.InstrumentId (%v < %v > %v)", t.InstrumentIdMinValue(), t.InstrumentId, t.InstrumentIdMaxValue())
	}

	if err := t.InstrumentState.RangeCheck(); err != nil {
		return err
	}

	if t.TimestampMs < t.TimestampMsMinValue() || t.TimestampMs > t.TimestampMsMaxValue() {
		return fmt.Errorf("range check failed on t.TimestampMs (%v < %v > %v)", t.TimestampMsMinValue(), t.TimestampMs, t.TimestampMsMaxValue())
	}

	if t.OpenInterest < t.OpenInterestMinValue() || t.OpenInterest > t.OpenInterestMaxValue() {
		return fmt.Errorf("range check failed on t.OpenInterest (%v < %v > %v)", t.OpenInterestMinValue(), t.OpenInterest, t.OpenInterestMaxValue())
	}

	if t.MinSellPrice < t.MinSellPriceMinValue() || t.MinSellPrice > t.MinSellPriceMaxValue() {
		return fmt.Errorf("range check failed on t.MinSellPrice (%v < %v > %v)", t.MinSellPriceMinValue(), t.MinSellPrice, t.MinSellPriceMaxValue())
	}

	if t.MaxBuyPrice < t.MaxBuyPriceMinValue() || t.MaxBuyPrice > t.MaxBuyPriceMaxValue() {
		return fmt.Errorf("range check failed on t.MaxBuyPrice (%v < %v > %v)", t.MaxBuyPriceMinValue(), t.MaxBuyPrice, t.MaxBuyPriceMaxValue())
	}

	if t.LastPrice != t.LastPriceNullValue() && (t.LastPrice < t.LastPriceMinValue() || t.LastPrice > t.LastPriceMaxValue()) {
		return fmt.Errorf("range check failed on t.LastPrice (%v < %v > %v)", t.LastPriceMinValue(), t.LastPrice, t.LastPriceMaxValue())
	}

	if t.IndexPrice < t.IndexPriceMinValue() || t.IndexPrice > t.IndexPriceMaxValue() {
		return fmt.Errorf("range check failed on t.IndexPrice (%v < %v > %v)", t.IndexPriceMinValue(), t.IndexPrice, t.IndexPriceMaxValue())
	}

	if t.MarkPrice < t.MarkPriceMinValue() || t.MarkPrice > t.MarkPriceMaxValue() {
		return fmt.Errorf("range check failed on t.MarkPrice (%v < %v > %v)", t.MarkPriceMinValue(), t.MarkPrice, t.MarkPriceMaxValue())
	}

	if t.BestBidPrice < t.BestBidPriceMinValue() || t.BestBidPrice > t.BestBidPriceMaxValue() {
		return fmt.Errorf("range check failed on t.BestBidPrice (%v < %v > %v)", t.BestBidPriceMinValue(), t.BestBidPrice, t.BestBidPriceMaxValue())
	}

	if t.BestBidAmount < t.BestBidAmountMinValue() || t.BestBidAmount > t.BestBidAmountMaxValue() {
		return fmt.Errorf("range check failed on t.BestBidAmount (%v < %v > %v)", t.BestBidAmountMinValue(), t.BestBidAmount, t.BestBidAmountMaxValue())
	}

	if t.BestAskPrice < t.BestAskPriceMinValue() || t.BestAskPrice > t.BestAskPriceMaxValue() {
		return fmt.Errorf("range check failed on t.BestAskPrice (%v < %v > %v)", t.BestAskPriceMinValue(), t.BestAskPrice, t.BestAskPriceMaxValue())
	}

	if t.BestAskAmount < t.BestAskAmountMinValue() || t.BestAskAmount > t.BestAskAmountMaxValue() {
		return fmt.Errorf("range check failed on t.BestAskAmount (%v < %v > %v)", t.BestAskAmountMinValue(), t.BestAskAmount, t.BestAskAmountMaxValue())
	}

	if t.CurrentFunding != t.CurrentFundingNullValue() && (t.CurrentFunding < t.CurrentFundingMinValue() || t.CurrentFunding > t.CurrentFundingMaxValue()) {
		return fmt.Errorf("range check failed on t.CurrentFunding (%v < %v > %v)", t.CurrentFundingMinValue(), t.CurrentFunding, t.CurrentFundingMaxValue())
	}

	if t.Funding8h != t.Funding8hNullValue() && (t.Funding8h < t.Funding8hMinValue() || t.Funding8h > t.Funding8hMaxValue()) {
		return fmt.Errorf("range check failed on t.Funding8h (%v < %v > %v)", t.Funding8hMinValue(), t.Funding8h, t.Funding8hMaxValue())
	}

	if t.EstimatedDeliveryPrice != t.EstimatedDeliveryPriceNullValue() && (t.EstimatedDeliveryPrice < t.EstimatedDeliveryPriceMinValue() || t.EstimatedDeliveryPrice > t.EstimatedDeliveryPriceMaxValue()) {
		return fmt.Errorf("range check failed on t.EstimatedDeliveryPrice (%v < %v > %v)", t.EstimatedDeliveryPriceMinValue(), t.EstimatedDeliveryPrice, t.EstimatedDeliveryPriceMaxValue())
	}

	if t.DeliveryPrice != t.DeliveryPriceNullValue() && (t.DeliveryPrice < t.DeliveryPriceMinValue() || t.DeliveryPrice > t.DeliveryPriceMaxValue()) {
		return fmt.Errorf("range check failed on t.DeliveryPrice (%v < %v > %v)", t.DeliveryPriceMinValue(), t.DeliveryPrice, t.DeliveryPriceMaxValue())
	}

	if t.SettlementPrice != t.SettlementPriceNullValue() && (t.SettlementPrice < t.SettlementPriceMinValue() || t.SettlementPrice > t.SettlementPriceMaxValue()) {
		return fmt.Errorf("range check failed on t.SettlementPrice (%v < %v > %v)", t.SettlementPriceMinValue(), t.SettlementPrice, t.SettlementPriceMaxValue())
	}

	return nil
}

func (*Ticker) SbeBlockLength() (blockLength uint16) {
	return 145
}

func (*Ticker) InstrumentIdMinValue() uint32 {
	return 0
}

func (*Ticker) InstrumentIdMaxValue() uint32 {
	return math.MaxUint32 - 1
}

func (*Ticker) TimestampMsMinValue() uint64 {
	return 0
}

func (*Ticker) TimestampMsMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*Ticker) OpenInterestMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) OpenInterestMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) MinSellPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) MinSellPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) MaxBuyPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) MaxBuyPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) LastPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) LastPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) LastPriceNullValue() float64 {
	return math.NaN()
}

func (*Ticker) IndexPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) IndexPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) MarkPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) MarkPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) BestBidPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) BestBidPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) BestBidAmountMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) BestBidAmountMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) BestAskPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) BestAskPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) BestAskAmountMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) BestAskAmountMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) CurrentFundingMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) CurrentFundingMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) CurrentFundingNullValue() float64 {
	return math.NaN()
}

func (*Ticker) Funding8hMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) Funding8hMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) Funding8hNullValue() float64 {
	return math.NaN()
}

func (*Ticker) EstimatedDeliveryPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) EstimatedDeliveryPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) EstimatedDeliveryPriceNullValue() float64 {
	return math.NaN()
}

func (*Ticker) DeliveryPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) DeliveryPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) DeliveryPriceNullValue() float64 {
	return math.NaN()
}

func (*Ticker) SettlementPriceMinValue() float64 {
	return -math.MaxFloat64
}

func (*Ticker) SettlementPriceMaxValue() float64 {
	return math.MaxFloat64
}

func (*Ticker) SettlementPriceNullValue() float64 {
	return math.NaN()
}
