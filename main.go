package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/subosito/gotenv"
)

// README
// -Mude o nome da função de mudar() para main()
// -Execute com:
// >go run ModelsScaffold.go nome_schema nome_tabela
func main() {
	gotenv.Load()
	schema := os.Args[1]
	table := os.Args[2]
	var prefixo string

	var columns []Column

	db, err := SetupDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	db.Table("information_schema.columns").Where("table_schema = ?", schema).Where("table_name = ?", table).Find(&columns)
	if len(columns) == 0 {
		fmt.Println("Tabela Inválida")
		return
	}

	infoColumns := []ColumnInfo{}

	for i, column := range columns {
		if i == 0 {
			prefixo = strings.Split(column.Name, "_")[0]
		}
		foreign := ""
		infoColumn := ColumnInfo{}
		infoColumn.Column = column.Name

		split := strings.Split(column.Name, "_")
		for j, ss := range split {
			if j == 0 {
				if ss != prefixo {
					foreign = ss
					continue
				}
				continue
			}

			infoColumn.Field = infoColumn.Field + strings.Title(ss)

			if j == len(split)-1 {
				infoColumn.JSON = infoColumn.JSON + ss
			} else {
				infoColumn.JSON = infoColumn.JSON + ss + "_"
			}
		}

		if foreign != "" {
			infoColumn.Field = infoColumn.Field + strings.Title(foreign)
			infoColumn.JSON = infoColumn.JSON + "_" + foreign
		}

		if column.Nullable == "NO" {
			infoColumn.Validate = "required"
			if column.Type == "integer" || column.Type == "numeric" || column.Type == "bigint" {
				infoColumn.Type = "int"
			} else if column.Type == "character varying" || column.Type == "text" {
				infoColumn.Type = "string"
			} else if column.Type == "timestamp without time zone" || column.Type == "date" || column.Type == "timestamp with time zone" {
				infoColumn.Type = "*time.Time"
			} else {
				infoColumn.Type = "interface{}"
			}
		} else {
			if column.Type == "integer" || column.Type == "numeric" || column.Type == "bigint" {
				infoColumn.Type = "*int"
			} else if column.Type == "character varying" || column.Type == "text" {
				infoColumn.Type = "*string"
			} else if column.Type == "timestamp without time zone" || column.Type == "date" || column.Type == "timestamp with time zone" {
				infoColumn.Type = "*time.Time"
			} else {
				infoColumn.Type = "interface{}"
			}
		}

		if column.Default != "" {
			infoColumn.Default = column.Default
		}

		if column.MaxChar != 0 {
			if infoColumn.Validate != "" {
				infoColumn.Validate = infoColumn.Validate + ",max=" + strconv.Itoa(column.MaxChar)
			} else {
				infoColumn.Validate = "omitempty,max=" + strconv.Itoa(column.MaxChar)
			}
		}
		infoColumns = append(infoColumns, infoColumn)
	}

	splitTable := strings.Split(table, "_")
	tableName := ""

	for _, ss := range splitTable {
		tableName = tableName + strings.Title(ss)
	}

	infoTable := TableInfo{
		Schema:    schema,
		Table:     table,
		TableName: tableName,
		Columns:   infoColumns,
	}

	if err := infoTable.makeModel(); err != nil {
		fmt.Println(err.Error())
	}
	if err := infoTable.makeRequest(); err != nil {
		fmt.Println(err.Error())
	}
	if err := infoTable.makeResource(); err != nil {
		fmt.Println(err.Error())
	}
	return
}
