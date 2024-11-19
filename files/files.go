package files

import (
	"demo/password/output"
	"os"

	"github.com/fatih/color"
)

type JsonDataBase struct {
	filename string
}

func NewJsonDataBase(name string) *JsonDataBase {
	return &JsonDataBase{
		filename: name,
	}
}

func (db *JsonDataBase) Read() ([]byte, error) {
	data, err := os.ReadFile(db.filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *JsonDataBase) Write(content []byte) {
	file, err := os.Create(db.filename)
	if err != nil {
		output.PrintError(err)
	}
	_, err = file.Write(content)
	defer file.Close()
	if err != nil {
		output.PrintError(err)
		return
	}
	color.Green("Запись успешна")
}
