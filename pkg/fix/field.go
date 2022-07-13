package fix

import (
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/shopspring/decimal"
)

//func getSendingTime(msg *quickfix.Message) (time.Time, error) {
//	return msg.Header.GetTime(tag.SendingTime)
//}

func getSymbol(msg *quickfix.Message) (v string, err error) {
	var f field.SymbolField
	if err := msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func newNoMDEntriesRepeatingGroup() *quickfix.RepeatingGroup {
	return quickfix.NewRepeatingGroup(
		tag.NoMDEntries,
		quickfix.GroupTemplate{
			quickfix.GroupElement(tag.MDUpdateAction),
			quickfix.GroupElement(tag.MDEntryType),
			quickfix.GroupElement(tag.MDEntryPx),
			quickfix.GroupElement(tag.MDEntrySize),
			quickfix.GroupElement(tag.MDEntryDate),
			quickfix.GroupElement(tagDeribitTradeID),
			quickfix.GroupElement(tag.Side),
			quickfix.GroupElement(tag.Price),
			quickfix.GroupElement(tag.Text),
			quickfix.GroupElement(tag.OrderID),
			quickfix.GroupElement(tag.SecondaryOrderID),
			quickfix.GroupElement(tag.OrdStatus),
			quickfix.GroupElement(tagDeribitLabel),
			quickfix.GroupElement(tagDeribitLiquidation),
			quickfix.GroupElement(tag.TrdMatchID),
		},
	)
}

func newSnapshotNoMDEntriesRepeatingGroup() *quickfix.RepeatingGroup {
	return quickfix.NewRepeatingGroup(
		tag.NoMDEntries,
		quickfix.GroupTemplate{
			quickfix.GroupElement(tag.MDEntryType),
			quickfix.GroupElement(tag.MDEntryPx),
			quickfix.GroupElement(tag.MDEntrySize),
			quickfix.GroupElement(tag.MDEntryDate),
			quickfix.GroupElement(tagDeribitTradeID),
			quickfix.GroupElement(tag.Side),
			quickfix.GroupElement(tag.Price),
			quickfix.GroupElement(tag.Text),
			quickfix.GroupElement(tag.OrderID),
			quickfix.GroupElement(tag.SecondaryOrderID),
			quickfix.GroupElement(tag.OrdStatus),
			quickfix.GroupElement(tagDeribitLabel),
			quickfix.GroupElement(tagDeribitLiquidation),
			quickfix.GroupElement(tag.TrdMatchID),
		},
	)
}

func getMDEntries(msg *quickfix.Message) (f *quickfix.RepeatingGroup, err error) {
	if msg.IsMsgTypeOf(string(enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH)) {
		f = newSnapshotNoMDEntriesRepeatingGroup()
	} else if msg.IsMsgTypeOf(string(enum.MsgType_MARKET_DATA_INCREMENTAL_REFRESH)) {
		f = newNoMDEntriesRepeatingGroup()
	} else {
		return
	}

	err = msg.Body.GetGroup(f)
	return f, err
}

func getMDEntryType(g *quickfix.Group) (v enum.MDEntryType, err error) {
	var f field.MDEntryTypeField
	if err = g.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getMDEntryPx(g *quickfix.Group) (v decimal.Decimal, err error) {
	var f field.MDEntryPxField
	if err = g.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func hasMDEntrySize(g *quickfix.Group) bool {
	return g.Has(tag.MDEntrySize)
}

func getMDEntrySize(g *quickfix.Group) (v decimal.Decimal, err error) {
	if hasMDEntrySize(g) {
		var f field.MDEntrySizeField
		if err = g.Get(&f); err == nil {
			v = f.Value()
		}
	}
	return
}

func hasMDUpdateAction(g *quickfix.Group) bool {
	return g.Has(tag.MDUpdateAction)
}

func getMDUpdateAction(g *quickfix.Group) (string, error) {
	action, err := g.GetString(tag.MDUpdateAction)
	if err != nil {
		return "", err
	}

	switch action {
	case "0":
		return "new", nil
	case "1":
		return "change", nil
	case "2":
		return "delete", nil
	default:
		return "", nil
	}
}

func getMDEntryDate(g *quickfix.Group) (v time.Time, err error) {
	var f field.MDEntryDateField
	if err = g.Get(&f); err == nil {
		v, err = time.Parse("20060102-15:04:05.000", f.Value())
	}
	return
}

func newNoMDEntryTypesRepeatingGroup() *quickfix.RepeatingGroup {
	return quickfix.NewRepeatingGroup(
		tag.NoMDEntryTypes,
		quickfix.GroupTemplate{
			quickfix.GroupElement(tag.MDEntryType),
		},
	)
}

func newNoRelatedSymRepeatingGroup() *quickfix.RepeatingGroup {
	return quickfix.NewRepeatingGroup(
		tag.NoRelatedSym,
		quickfix.GroupTemplate{
			quickfix.GroupElement(tag.NoRelatedSym),
		},
	)
}

func getText(msg *quickfix.Message) (v string, err error) {
	var f field.TextField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getOrderStatus(msg *quickfix.Message) (v enum.OrdStatus, err error) {
	var f field.OrdStatusField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getOrderID(msg *quickfix.Message) (v string, err error) {
	var f field.OrderIDField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getOrdType(msg *quickfix.Message) (v enum.OrdType, err error) {
	var f field.OrdTypeField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getSide(msg *quickfix.Message) (v enum.Side, err error) {
	var f field.SideField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getOrderQty(msg *quickfix.Message) (v decimal.Decimal, err error) {
	var f field.OrderQtyField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getCumQty(msg *quickfix.Message) (v decimal.Decimal, err error) {
	var f field.CumQtyField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getPrice(msg *quickfix.Message) (v decimal.Decimal, err error) {
	var f field.PriceField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getAvgPx(msg *quickfix.Message) (v decimal.Decimal, err error) {
	var f field.AvgPxField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getCommission(msg *quickfix.Message) (v decimal.Decimal, err error) {
	var f field.CommissionField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getMaxShow(msg *quickfix.Message) (v decimal.Decimal, err error) {
	var f field.MaxShowField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getTransactTime(msg *quickfix.Message) (v time.Time, err error) {
	var f field.TransactTimeField
	if err = msg.Body.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getDeribitLabel(msg *quickfix.Message) (string, error) {
	return msg.Body.GetString(tagDeribitLabel)
}

func hasExecInst(msg *quickfix.Message) bool {
	return msg.Body.Has(tag.ExecInst)
}

func getExecInst(msg *quickfix.Message) (string, error) {
	if !hasExecInst(msg) {
		return "", nil
	}
	return msg.Body.GetString(tag.ExecInst)
}

func getMarkPrice(msg *quickfix.Message) (decimal.Decimal, error) {
	if !msg.Body.Has(tagMarkPrice) {
		return decimal.New(0, 0), nil
	}

	price, err := msg.Body.GetString(tagMarkPrice)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return decimal.NewFromString(price)
}

func getGroupPrice(g *quickfix.Group) (v decimal.Decimal, err error) {
	var f field.PriceField
	if err = g.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getGroupSide(g *quickfix.Group) (v enum.Side, err error) {
	var f field.SideField
	if err = g.Get(&f); err == nil {
		v = f.Value()
	}
	return
}

func getGroupTradeID(g *quickfix.Group) (string, error) {
	return g.GetString(tagDeribitTradeID)
}
