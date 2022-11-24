package insert

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/ramsgoli/columnar_store/backend/meta"
	"github.com/ramsgoli/columnar_store/config"
	"github.com/ramsgoli/columnar_store/util"
)

type InsertDetails struct {
	TableName [8]byte
	Attrs     []string
}

func findTableMetadataForTable(allTables *[]meta.TableMetadata, tableName [8]byte) *meta.TableMetadata {
	for _, t := range *allTables {
		if bytes.Compare(t.TableName[:], tableName[:]) == 0 {
			return &t
		}
	}
	return nil
}

func isValidUint(num string, size int) (uint64, error) {
	return strconv.ParseUint(num, 10, size*8)
}

// check all the cols
func checkAttrs(tableMetadata *meta.TableMetadata, attrs *[]string) ([][]byte, error) {
	if int(tableMetadata.NumCols) != len(*attrs) {
		return nil, errors.New("Invalid number of columns specified")
	}

	var mappedAttrs [][]byte
	for i := 0; i < int(tableMetadata.NumCols); i++ {
		col := (*tableMetadata.ColMetadata)[i]
		var size int
		var exists bool
		if size, exists = meta.TypeToSizeMap[col.Type]; !exists {
			return nil, fmt.Errorf("Unrecognized size %d", col.Type)
		}
		finalByteSlice := make([]byte, size)
		switch col.Type {
		case 0:
			var convertedUint uint64
			var validUintErr error
			if convertedUint, validUintErr = isValidUint((*attrs)[i], size); validUintErr != nil {
				return nil, fmt.Errorf("invalid value of size %d bytes: %s", size, (*attrs)[i])
			}
			finalByteSlice[0] = byte(convertedUint)
		default:
			if len((*attrs)[i]) > size {
				return nil, errors.New("invalid varchar value")
			}
			// append each char in string to `finalByteSlice`
			for i_idx, b := range (*attrs)[i] {
				finalByteSlice[i_idx] = byte(b)
			}
		}
		mappedAttrs = append(mappedAttrs, finalByteSlice)
	}
	// no errors so far
	return mappedAttrs, nil
}

func writeCol(m *meta.TableMetadata, data []byte, i int) error {
	colPath := path.Join(config.GetTablesPath(), util.Trim(m.TableName[:]), util.Trim((*m.ColMetadata)[i].ColName[:]))
	colFileHandle, openErr := os.OpenFile(colPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer colFileHandle.Close()

	if openErr != nil {
		return openErr
	}
	_, writeErr := colFileHandle.Write(data)
	return writeErr
}

func Insert(i *InsertDetails) error {
	// check if table exists
	allTables, err := meta.GetAllTables()
	if err != nil {
		return err
	}

	tableMetadata := findTableMetadataForTable(allTables.Tables, i.TableName)
	if tableMetadata == nil {
		return errors.New(fmt.Sprintf("Table %s does not exist", i.TableName))
	}

	mappedAttrs, checkAttrsErr := checkAttrs(tableMetadata, &i.Attrs)
	if checkAttrsErr != nil {
		return err
	}

	numCols := int(tableMetadata.NumCols)
	for i := 0; i < numCols; i++ {
		writeColErr := writeCol(tableMetadata, mappedAttrs[i], i)
		if writeColErr != nil {
			return writeColErr
		}
	}

	// write all the shits
	return nil
}
