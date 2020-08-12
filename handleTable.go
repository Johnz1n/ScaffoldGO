package main

import (
	"io"
	"os"
	"strings"
)

// Column database struct
type Column struct {
	Name     string `gorm:"column:column_name"`
	Default  string `gorm:"column:column_default"`
	Nullable string `gorm:"column:is_nullable"`
	Type     string `gorm:"column:data_type"`
	MaxChar  int    `gorm:"column:character_maximum_length"`
}

// TableInfo godoc
type TableInfo struct {
	Schema    string
	Table     string
	TableName string
	Columns   []ColumnInfo
}

// ColumnInfo output struct
type ColumnInfo struct {
	Field    string
	Type     string
	JSON     string
	Column   string
	Default  string
	Validate string
}

func (t *TableInfo) makeModel() error {
	model := `package models

	import "time"
	
	// ` + t.TableName + ` ...
	type ` + t.TableName + ` struct {
		
	`
	var fill string
	for _, line := range t.Columns {
		model += line.Field + " " + line.Type + "`" + `gorm:"column:` + line.Column
		if line.Default != "" && !strings.Contains(line.Default, "nextval") {
			model += `; default:` + line.Default + `"`
		} else if strings.Contains(line.Default, "nextval") {
			model += `; primary_key:true"`
		} else {
			model += `"`
		}

		fill += "	" + strings.ToLower(string([]rune(t.TableName)[0])) + "." + line.Field + " = req." + line.Field + "\n"
		model += "`\n"
	}

	model += `}

	// TableName Seta o nome da tabela
	func (` + strings.ToLower(string([]rune(t.TableName)[0])) + ` *` + t.TableName + `) TableName() string {
		return "` + t.Schema + `.` + t.Table + `"
	}
	
	//` + t.TableName + `Fill preenche o model a partir de um request
	func (` + strings.ToLower(string([]rune(t.TableName)[0])) + ` *` + t.TableName + `) ` + t.TableName + `Fill(req requests.` + t.TableName + `) {` +
		fill +
		`}`

	file, err := os.Create(t.TableName + ".models")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(model))
	if err != nil {
		return err
	}

	return nil
}

func (t *TableInfo) makeRequest() error {
	model := `package requests

	import "time"
	
	// ` + t.Table + ` ...
	type ` + t.Table + ` struct {	

	`

	for _, line := range t.Columns {
		model += line.Field + " " + line.Type + "`" + `json:"` + line.JSON + `"`

		if line.Validate != "" {
			model += ` validate:"` + line.Validate + `"`
		}

		model += "`\n"
	}

	model += `}`

	file, err := os.Create(t.TableName + ".requests")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(model))
	if err != nil {
		return err
	}

	return nil
}

func (t *TableInfo) makeResource() error {
	model := `package resources

	import "time"
	
	// ` + t.TableName + ` ...
	type ` + t.TableName + ` struct {

	`
	var resource string
	for _, line := range t.Columns {
		model += line.Field + " " + line.Type + "`" + `json:"` + line.JSON + `"` + "`\n"
		resource += "	" + strings.ToLower(string([]rune(t.TableName)[0])) + "." + line.Field + " = mod." + line.Field + "\n"
	}

	model += `}

	//` + t.TableName + `Resource preenche um ressource apartir de um model
	func (` + strings.ToLower(string([]rune(t.TableName)[0])) + ` *` + t.TableName + `) ` + t.TableName + `Resource(mod models.` + t.TableName + `) {
		` +
		resource +
		`}`

	file, err := os.Create(t.TableName + ".resources")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(model))
	if err != nil {
		return err
	}

	return nil
}
