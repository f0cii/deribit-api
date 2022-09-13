package sbe

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeBook(t *testing.T) {
	tests := []struct {
		event          []byte
		expectedOutput Book
		expectedError  error
	}{
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x96, 0x37, 0x03, 0x00,
				0x77, 0xc4, 0x15, 0x0d, 0x83, 0x01, 0x00, 0x00, 0x3c, 0x25, 0x7a, 0x7f, 0x0b, 0x00, 0x00, 0x00,
				0x3d, 0x25, 0x7a, 0x7f, 0x0b, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x60, 0x4e, 0xd3, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc0,
				0x4f, 0xed, 0x40,
			},
			Book{
				InstrumentId: 210838,
				TimestampMs:  1662371873911,
				PrevChangeId: 49383351612,
				ChangeId:     49383351613,
				IsLast:       YesNo.Yes,
				ChangesList: []BookChangesList{
					{
						Side:   BookSide.Bid,
						Change: BookChange.Changed,
						Price:  19769.5,
						Amount: 60030,
					},
				},
			},
			nil,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x48, 0x37, 0x03, 0x00,
				0xda, 0x26, 0xe6, 0x15, 0x83, 0x01, 0x00, 0x00, 0x06, 0x46, 0x6e, 0xa0, 0x06, 0x00, 0x00, 0x00,
				0x08, 0x46, 0x6e, 0xa0, 0x06, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x33, 0x33, 0x33, 0x33, 0x33, 0x6c, 0x97, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x20, 0x40,
			},
			Book{
				InstrumentId: 210760,
				TimestampMs:  1662519748314,
				PrevChangeId: 28461385222,
				ChangeId:     28461385224,
				IsLast:       YesNo.Yes,
				ChangesList: []BookChangesList{
					{
						Side:   BookSide.Ask,
						Change: BookChange.Changed,
						Price:  1499.05,
						Amount: 8,
					},
				},
			},
			nil,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x48, 0x37,
			},
			Book{},
			io.ErrUnexpectedEOF,
		},
		// Some range check error cases
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x48, 0x37, 0x03, 0x00,
				0xda, 0x26, 0xe6, 0x15, 0x83, 0x01, 0x00, 0x00, 0x06, 0x46, 0x6e, 0xa0, 0x06, 0x00, 0x00, 0x00,
				0x08, 0x46, 0x6e, 0xa0, 0x06, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x03, 0x01, 0x33, 0x33, 0x33, 0x33, 0x33, 0x6c, 0x97, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x20, 0x40,
			},
			Book{},
			ErrRangeCheck, // ChangesList - Side
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x48, 0x37, 0x03, 0x00,
				0xda, 0x26, 0xe6, 0x15, 0x83, 0x01, 0x00, 0x00, 0x06, 0x46, 0x6e, 0xa0, 0x06, 0x00, 0x00, 0x00,
				0x08, 0x46, 0x6e, 0xa0, 0x06, 0x00, 0x00, 0x00, 0x03, 0x12, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x33, 0x33, 0x33, 0x33, 0x33, 0x6c, 0x97, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x20, 0x40,
			},
			Book{},
			ErrRangeCheck, // IsLast
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x48, 0x37, 0x03, 0x00,
				0xda, 0x26, 0xe6, 0x15, 0x83, 0x01, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0x08, 0x46, 0x6e, 0xa0, 0x06, 0x00, 0x00, 0x00, 0x03, 0x12, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x33, 0x33, 0x33, 0x33, 0x33, 0x6c, 0x97, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x20, 0x40,
			},
			Book{},
			ErrRangeCheck, // PreviousChangeId
		},
	}

	marshaller := NewSbeGoMarshaller()

	for _, test := range tests {
		bufferData := bytes.NewBuffer(test.event)

		var header MessageHeader
		err := header.Decode(marshaller, bufferData)
		require.NoError(t, err)

		err = header.RangeCheck()
		require.NoError(t, err)

		var book Book
		err = book.Decode(marshaller, bufferData, header.BlockLength, true)
		require.ErrorIs(t, err, test.expectedError)

		if err == nil {
			assert.Equal(t, book, test.expectedOutput)
		}
	}
}