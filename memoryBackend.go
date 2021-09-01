package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

type memoryCell []byte

type MemoryBackend struct {
	tables map[string]*Table
}

type Record struct {
	Data     map[string]interface{}
	Metadata *TableMetadata
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer

	switch key.(type) {
	case int64:
		b := make([]byte, 8)
		uint_value := uint64(key.(int64))
		binary.BigEndian.PutUint64(b, uint_value)
		return b, nil

	default:
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(key)
		if err != nil {
			return nil, err
		}
		return buf.Bytes()[4:], nil
	}
}

func (r Record) toBytes() (bt []byte) {
	bt = make([]byte, r.Metadata.Offset)
	var pos int = 0
	for _, column := range r.Metadata.Columns {
		data := r.Data[column.Name]
		byte_arr, err := GetBytes(data)
		check(err)
		copy(bt[pos:], byte_arr)
		pos += column.Size
	}
	return
}

type Table struct {
	Metadata *TableMetadata
	Records  []Record
}

/*
func createTable() *Table {
	return &Table{
		Name:    "?tmp?",
		Columns: nil,
		Rows:    nil,
	}
}
*/

type ColumnType uint

const (
	TextType ColumnType = iota
	IntType
)

type TableMetadata struct {
	Name    string
	Columns []Column
	Offset  int
}

type Column struct {
	Name     string
	Datatype ColumnType
	Size     int
}
