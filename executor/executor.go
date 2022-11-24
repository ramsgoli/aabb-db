package executor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ramsgoli/columnar_store/backend/insert"
	"github.com/ramsgoli/columnar_store/backend/meta"
)

func handleDescribe(tokens []string) {
	if t, err := meta.GetAllTables(); err == nil {
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

	colMetadatas := []meta.ColMetadata{}
	for i := 0; i < len(colDetails); i += 2 {
		colName := colDetails[i]

		var colType uint8
		if colTypeInt, err := strconv.Atoi(colDetails[i+1]); err == nil {
			colType = uint8(colTypeInt)
		}

		var colByteArray [8]byte
		copy(colByteArray[:], []byte(colName))
		colMetadatas = append(colMetadatas, meta.ColMetadata{ColName: colByteArray, Type: colType})
	}
	if err := meta.CreateTable(&meta.TableMetadata{
		TableName:   tableNameByteArray,
		NumCols:     numColumns,
		ColMetadata: &colMetadatas,
	}); err != nil {
		panic(err)
	}
}

func handleInsert(tokens []string) {
	tableName := tokens[0]
	attributes := tokens[1:]

	var tableNameBytes [8]byte
	copy(tableNameBytes[:], []byte(tableName))
	if err := insert.Insert(&insert.InsertDetails{TableName: tableNameBytes, Attrs: attributes}); err != nil {
		panic(err)
	}
}

func tokenize(command string) []string {
	return strings.Split(command, " ")
}

func Execute(command string) error {
	m := map[string]interface{}{
		"\\d":    handleDescribe,
		"create": handleCreate,
		"insert": handleInsert,
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
