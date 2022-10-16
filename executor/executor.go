package executor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ramsgoli/columnar_store/backend/tables"
)

func handleGet(tokens []string) {
	if t, err := tables.GetAllTables(); err == nil {
		fmt.Printf("Got %d tables\n", len(*t.Tables))
	} else {
		panic(err)
	}
}

func handleCreate(tokens []string) {
	tableName := tokens[0]
	var tableNameByteArray [8]byte
	copy(tableNameByteArray[:], []byte(tableName))
	colDetails := tokens[1:]
	numColumns := uint8(len(colDetails) / 2)

	colMetadatas := []tables.ColMetadata{}
	for i := 0; i < len(colDetails); i += 2 {
		colName := colDetails[i]

		var colType uint8
		if colTypeInt, err := strconv.Atoi(colDetails[i+1]); err != nil {
			colType = uint8(colTypeInt)
		}

		var colByteArray [8]byte
		copy(colByteArray[:], []byte(colName))
		colMetadatas = append(colMetadatas, tables.ColMetadata{ColName: colByteArray, Type: colType})
	}
	if err := tables.CreateTable(&tables.TableMetadata{
		TableName:   tableNameByteArray,
		NumCols:     numColumns,
		ColMetadata: &colMetadatas,
	}); err != nil {
		panic(err)
	}
}

func tokenize(command string) []string {
	return strings.Split(command, " ")
}

func Execute(command string) error {
	m := map[string]interface{}{
		"get":    handleGet,
		"create": handleCreate,
	}

	tokens := tokenize(command)
	first := tokens[0]
	if f, exists := m[first]; exists {
		f.(func([]string))(tokens[1:])
	} else {
		return errors.New(fmt.Sprintf("Could not recognize command %s\n", command))
	}
	return nil
}
