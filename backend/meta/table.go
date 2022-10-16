package meta

import (
	"errors"
	"fmt"
	"os"

	"github.com/ramsgoli/columnar_store/config"
)

type ColMetadata struct {
	ColName [8]byte
	Type    uint8
}

type TableMetadata struct {
	TableName   [8]byte
	NumCols     uint8
	ColMetadata *[]ColMetadata
}

type Tables struct {
	Tables *[]TableMetadata
}

const HEADER_SIZE = 8
const TABLE_NAME_SIZE int64 = 8
const NUM_COLS_SIZE = 1

func GetAllTables() (*Tables, error) {
	tableMetadataPath := config.GetTableMetadataPath()
	tableMetadataFile, err := os.Open(tableMetadataPath)
	if err != nil {
		return nil, errors.New("Couldn't open table metadata file")
	}
	defer tableMetadataFile.Close()

	header := [HEADER_SIZE]byte{}
	_, err = tableMetadataFile.Read(header[:])
	if err != nil {
		return nil, errors.New("Couldn't read table metadata header")
	}

	numTables := header[0]

	var tables []TableMetadata
	// for each table
	// idea: could spawn a goroutine to read each?
	tableMetadataFile.Seek(HEADER_SIZE, 0)
	for i := 0; i < int(numTables); i++ {
		var tableName [TABLE_NAME_SIZE]byte
		_, err = tableMetadataFile.Read(tableName[:])

		var numColumnsArr [NUM_COLS_SIZE]byte
		_, err = tableMetadataFile.Read(numColumnsArr[:])
		numColumns := numColumnsArr[0]

		var columnMetadata []ColMetadata

		// for each col
		for j := 0; j < int(numColumns); j++ {
			var columnName [8]byte
			var columnType [1]byte

			_, err = tableMetadataFile.Read(columnName[:])
			_, err = tableMetadataFile.Read(columnType[:])

			c := ColMetadata{ColName: columnName, Type: columnType[0]}
			columnMetadata = append(columnMetadata, c)
		}

		t := TableMetadata{
			TableName:   tableName,
			NumCols:     numColumns,
			ColMetadata: &columnMetadata,
		}
		tables = append(tables, t)
	}

	return &Tables{Tables: &tables}, nil
}

func CreateTable(t *TableMetadata) error {
	tableMetadataPath := config.GetTableMetadataPath()
	tableMetadataFile, err := os.OpenFile(tableMetadataPath, os.O_RDWR, 0)
	if err != nil {
		return nil
	}
	defer tableMetadataFile.Close()

	// Check if file is empty.
	// If so, add the header
	stat, statErr := tableMetadataFile.Stat()
	if statErr != nil {
		return statErr
	}
	if stat.Size() == 0 {
		fmt.Println("_meta file is empty. Creating a header")
		header := [HEADER_SIZE]byte{0}
		tableMetadataFile.Write(header[:])
	}

	// seek to end of file
	tableMetadataFile.Seek(0, 2)

	// write table name
	_, err = tableMetadataFile.Write(t.TableName[:])
	if err != nil {
		return err
	}

	_, err = tableMetadataFile.Write([]byte{t.NumCols})
	if err != nil {
		return err
	}

	// write all column details
	for i := 0; i < int(t.NumCols); i++ {
		m := (*t.ColMetadata)[i]
		columnName := m.ColName
		columnType := m.Type

		_, err = tableMetadataFile.Write(columnName[:])
		if err != nil {
			return err
		}

		_, err = tableMetadataFile.Write([]byte{columnType})
		if err != nil {
			return err
		}
	}

	// update header
	var header [HEADER_SIZE]byte
	tableMetadataFile.ReadAt(header[:], 0)
	header[0] = header[0] + 1
	if _, headerWriteErr := tableMetadataFile.WriteAt(header[:], 0); headerWriteErr != nil {
		return headerWriteErr
	}

	return nil
}
