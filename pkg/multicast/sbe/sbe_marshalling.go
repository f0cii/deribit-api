// Copyright (C) 2017 MarketFactory, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// This file provides a simple bespoke marshalling layer for the
// standard binary encoding golang backend and is part of:
//
// https://github.com/real-logic/simple-binary-encoding

package sbe

import (
	"io"
	"math"
)

// Allocate via NewSbeGoMarshaller to initialize
type SbeGoMarshaller struct {
	b8 []byte // statically allocated tmp space to avoid alloc
	b1 []byte // previously created slice into b to save time
	b2 []byte // previously created slice into b to save time
	b4 []byte // previously created slice into b to save time
}

func NewSbeGoMarshaller() *SbeGoMarshaller {
	var m SbeGoMarshaller
	m.b8 = make([]byte, 8)
	m.b1 = m.b8[:1]
	m.b2 = m.b8[:2]
	m.b4 = m.b8[:4]
	return &m
}

// The "standard" MessageHeader.
//
// Most applications will use this as it's the default and optimized
// although it's possible to change it by:
// a) using a different sized BlockLength, or
// b) adding arbitrary fields
//
// If the MessageHeader is not "standard" then you can use the
// generated MessageHeader type in MessageHeader.go otherwise we
// recommend this one.
type SbeGoMessageHeader struct {
	BlockLength uint16
	TemplateId  uint16
	SchemaId    uint16
	Version     uint16
}

func (m *SbeGoMarshaller) ReadUint8(r io.Reader, v *uint8) error {
	if _, err := io.ReadFull(r, m.b1); err != nil {
		return err
	}
	*v = uint8(m.b1[0])
	return nil
}

func (m *SbeGoMarshaller) ReadUint16(r io.Reader, v *uint16) error {
	if _, err := io.ReadFull(r, m.b2); err != nil {
		return err
	}
	*v = (uint16(m.b2[0]) | uint16(m.b2[1])<<8)
	return nil
}

func (m *SbeGoMarshaller) ReadUint32(r io.Reader, v *uint32) error {
	if _, err := io.ReadFull(r, m.b4); err != nil {
		return err
	}
	*v = (uint32(m.b4[0]) | uint32(m.b4[1])<<8 |
		uint32(m.b4[2])<<16 | uint32(m.b4[3])<<24)
	return nil
}

func (m *SbeGoMarshaller) ReadUint64(r io.Reader, v *uint64) error {
	if _, err := io.ReadFull(r, m.b8); err != nil {
		return err
	}
	*v = (uint64(m.b8[0]) | uint64(m.b8[1])<<8 |
		uint64(m.b8[2])<<16 | uint64(m.b8[3])<<24 |
		uint64(m.b8[4])<<32 | uint64(m.b8[5])<<40 |
		uint64(m.b8[6])<<48 | uint64(m.b8[7])<<56)
	return nil
}

func (m *SbeGoMarshaller) ReadFloat64(r io.Reader, v *float64) error {
	if _, err := io.ReadFull(r, m.b8); err != nil {
		return err
	}
	*v = math.Float64frombits(uint64(m.b8[0]) | uint64(m.b8[1])<<8 |
		uint64(m.b8[2])<<16 | uint64(m.b8[3])<<24 |
		uint64(m.b8[4])<<32 | uint64(m.b8[5])<<40 |
		uint64(m.b8[6])<<48 | uint64(m.b8[7])<<56)
	return nil
}

func (m *SbeGoMarshaller) ReadBytes(r io.Reader, b []byte) error {
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	return nil
}
