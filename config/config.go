package config

import (
	"os"
	"path"
)

var dataPath string

func init() {
	d, err := os.UserHomeDir()
	if err != nil {
		panic("Couldn't obtain home dir")
	}
	dataPath = path.Join(d, "columnar", "data")
}

func GetTablesPath() string {
	return path.Join(dataPath, "tables")
}

func GetTableMetadataPath() string {
	return path.Join(dataPath, "meta", "_tables")
}
