package sbe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookChangeString(t *testing.T) {
	tests := []struct {
		bookChange BookChangeEnum
		expected   string
	}{
		{
			BookChange.Created,
			"new",
		},
		{
			BookChange.Changed,
			"change",
		},
		{
			BookChange.Deleted,
			"delete",
		},
		{
			BookChange.NullValue,
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.bookChange.String(), test.expected)
	}
}

func TestDirectionString(t *testing.T) {
	tests := []struct {
		direction DirectionEnum
		expected  string
	}{
		{
			Direction.Buy,
			"buy",
		},
		{
			Direction.Sell,
			"sell",
		},
		{
			Direction.NullValue,
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.direction.String(), test.expected)
	}
}

func TestInstrumentKindString(t *testing.T) {
	tests := []struct {
		kind     InstrumentKindEnum
		expected string
	}{
		{
			InstrumentKind.Future,
			"future",
		},
		{
			InstrumentKind.Option,
			"option",
		},
		{
			InstrumentKind.NullValue,
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.kind.String(), test.expected)
	}
}

func TestInstrumentStateString(t *testing.T) {
	tests := []struct {
		kind     InstrumentStateEnum
		expected string
	}{
		{
			InstrumentState.Created,
			"created",
		},
		{
			InstrumentState.Open,
			"open",
		},
		{
			InstrumentState.Closed,
			"closed",
		},
		{
			InstrumentState.Settled,
			"settled",
		},
		{
			InstrumentState.NullValue,
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.kind.String(), test.expected)
	}
}

func TestLiquidationString(t *testing.T) {
	tests := []struct {
		kind     LiquidationEnum
		expected string
	}{
		{
			Liquidation.None,
			"none",
		},
		{
			Liquidation.Maker,
			"maker",
		},
		{
			Liquidation.Taker,
			"taker",
		},
		{
			Liquidation.Both,
			"both",
		},
		{
			Liquidation.NullValue,
			"none",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.kind.String(), test.expected)
	}
}

func TestOptionTypeString(t *testing.T) {
	tests := []struct {
		optionType OptionTypeEnum
		expected   string
	}{
		{
			OptionType.NotApplicable,
			"not_applicable",
		},
		{
			OptionType.Call,
			"call",
		},
		{
			OptionType.Put,
			"put",
		},
		{
			OptionType.NullValue,
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.optionType.String(), test.expected)
	}
}

func TestPeriodString(t *testing.T) {
	tests := []struct {
		period   PeriodEnum
		expected string
	}{
		{
			Period.Perpetual,
			"perpetual",
		},
		{
			Period.Minute,
			"minute",
		},
		{
			Period.Hour,
			"hour",
		},
		{
			Period.Day,
			"day",
		},
		{
			Period.Week,
			"week",
		},
		{
			Period.Month,
			"month",
		},
		{
			Period.Year,
			"year",
		},
		{
			Period.NullValue,
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.period.String(), test.expected)
	}
}
