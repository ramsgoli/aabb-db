package config

import (
	"os"
	"path"
)

var DataPath string

func init() {
	d, err := os.UserHomeDir()
	if err != nil {
		panic("Couldn't obtain home dir")
	}
	DataPath = path.Join(d, "columnar", "data")
}

func GetTablePath() string {
	return path.Join(DataPath, "tables")
}
