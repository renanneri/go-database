package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	tableFromCSV("data.csv")
}

func tableFromCSV(s string) {
	f, err := os.Open(s)
	check(err)

	scanner := bufio.NewScanner(f)

	table := Table{}
	metadata := TableMetadata{}
	metadata.Name = "Penalties"
	table.Metadata = &metadata
	// Reading Header
	scanner.Scan()
	line := scanner.Text()
	metadata.Columns, metadata.Offset = parseHeader(line)

	for scanner.Scan() {
		line := scanner.Text()
		record := stringToRecord(line, table)
		table.Records = append(table.Records, record)

	}
	table.persiste()
	table.flush()
	table.recovery()
}

func stringToRecord(s string, table Table) (record Record) {
	splited := strings.Split(s, ";")
	record.Data = make(map[string]interface{})
	record.Metadata = table.Metadata
	for i, value := range table.Metadata.Columns {
		switch value.Datatype {
		case IntType:
			intValue, err := strconv.ParseInt(splited[i], 10, 64)
			check(err)
			record.Data[value.Name] = intValue
		default:
			record.Data[value.Name] = splited[i]
		}
	}

	return
}

func parseHeader(s string) (columns []Column, offset int) {
	splited := strings.Split(s, ";")

	var size int
	var err error
	for _, v := range splited {
		column_info := strings.Split(v, "|")
		column := Column{}
		column.Name = column_info[0]

		switch column_info[1] {
		case "int":
			column.Datatype = IntType
			size = 8
		case "string":
			column.Datatype = TextType
			size, err = strconv.Atoi(column_info[2])
			check(err)
		default:
			column.Datatype = TextType
			size, err = strconv.Atoi(column_info[2])
			check(err)

		}
		column.Size = size
		offset += size
		columns = append(columns, column)
	}

	return
}

func (table Table) persiste() {
	output := table.Metadata.Name + ".db"
	var file, err = os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for _, record := range table.Records {
		bytes := record.toBytes()
		file.Write(bytes)
	}
}

func (table Table) recovery() {
	output := table.Metadata.Name + ".db"
	file, err := os.OpenFile(output, os.O_RDONLY, 0644)
	check(err)
	stat, err := file.Stat()
	check(err)
	size := stat.Size()
	offset := int64(table.Metadata.Offset)
	var i int64
	byte_arr := make([]byte, offset)

	for i = 0; i < size; i += offset {
		file.Seek(i, 0)
		_, err := file.Read(byte_arr)
		check(err)
		record := table.bytesToRecord(byte_arr)
		table.Records = append(table.Records, record)
		fmt.Println(record)
	}
}

func (table Table) bytesToRecord(bt []byte) (record Record) {
	record.Data = make(map[string]interface{})
	record.Metadata = table.Metadata
	offset := 0
	for _, Column := range table.Metadata.Columns {
		switch Column.Datatype {
		case IntType:
			record.Data[Column.Name] = binary.BigEndian.Uint64(bt[offset : offset+Column.Size])
		default:
			record.Data[Column.Name] = string(bt[offset : offset+Column.Size])
		}
		offset += Column.Size
	}
	return
}

func (table Table) flush() {
	table.Records = make([]Record, 10)
}
