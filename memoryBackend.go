package mygosql

type memoryCell []byte

type MemoryBackend struct {
	tables map[string]*table
}

type table struct {
	name        string
	columns     []string
	columnTypes []ColumnType
	rows        [][]memoryCell
}

func createTable() *table {
	return &table{
		name:        "?tmp?",
		columns:     nil,
		columnTypes: nil,
		rows:        nil,
	}
}

type ColumnType uint

const (
	TextType ColumnType = iota
	IntType
)

func (c ColumnType) String() string {
	switch c {
	case TextType:
		return "TextType"
	case IntType:
		return "IntType"
	default:
		return "Error"
	}
}

type Cell interface {
	AsText() *string
	AsInt() *int32
}

type Results struct {
	Columns []ResultColumn
	Rows    [][]Cell
}

type ResultColumn struct {
	Type    ColumnType
	Name    string
	NotNull bool
}

type Index struct {
	Name       string
	Exp        string
	Type       string
	Unique     bool
	PrimaryKey bool
}

type TableMetadata struct {
	Name    string
	Columns []ResultColumn
	Indexes []Index
}
