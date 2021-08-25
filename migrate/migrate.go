package migrate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func CreateTable(table string, columns []Column, args Args) {
}

func ModifyTable(table string, columns []Column) {
}

func DropTable(table string) {
}

func GetSql(filePath string, runType uint8) (sqlList []string, err error) {
	fileContent := ""
	fileContent, err = readAll(filePath)
	if err != nil {
		return
	}
	myLine := Line{
		filePath: filePath,
	}
	regString := `type *?(\w*?) *?struct *?\{(?s).*?\}`
	reg := regexp.MustCompile(regString)
	list := reg.FindAllStringSubmatch(fileContent, -1)
	if len(list) == 0 {
		err = Error("struct does not exist")
		return
	}
	if len(list) > 1 {
		err = Error("struct has %d quantities, but only 1 is required", len(list))
		return
	}
	tableStruct := list[0][1]
	myLine.SetLine(list[0][0])
	funcName := ""
	if runType == RunType.Up {
		funcName = "Up"
	} else if runType == RunType.Down {
		funcName = "Down"
	} else {
		err = errors.New("Parameter 'runType' of method 'GetSql' must be in migrate.RunType")
		return
	}
	regString = `func *?\(\w *?(\w+) *?\) *?(\w+) *?\( *?\) *? \{ *?\n(?s).*?\}\n *`
	reg = regexp.MustCompile(regString)
	for _, funcV := range reg.FindAllStringSubmatch(fileContent, -1) {
		myLine.SetLine(funcV[0])
		if tableStruct != funcV[1] {
			err = Error("struct '%s' is different from '%s'", funcV[1], tableStruct)
			return
		}
		if funcV[2] != "Up" && funcV[2] != "Down" {
			err = Error("func '%s' must be Up or Down", funcV[2])
			return
		}
		if funcName != funcV[2] {
			continue
		}
		regString = `migrate.(Create|Modify|Drop)Table *?\( *?"(\w*?)"(?s).*?\)`
		reg = regexp.MustCompile(regString)
		for _, migrateV := range reg.FindAllStringSubmatch(funcV[0], -1) {
			myLine.SetLine(migrateV[0])
			switch migrateV[1] {
			case "Create":
				sql := ""
				if sql, err = getCreateTableSql(migrateV[2], migrateV[0]); err != nil {
					return
				}
				sqlList = append(sqlList, sql)
			case "Modify":
				tpSql := []string{}
				if tpSql, err = getModifyTableSql(migrateV[2], migrateV[0]); err != nil {
					return
				}
				for _, sql := range tpSql {
					sqlList = append(sqlList, sql)
				}
			case "Drop":
				sqlList = append(sqlList, getDropTableSql(migrateV[2]))
			}
		}
	}
	return
}

func getDropTableSql(table string) string {
	return fmt.Sprintf("DROP TABLE `%s`", table)
}

func getCreateTableSql(table string, column string) (sql string, err error) {
	columns := []map[string]string{}
	if columns, err = getColumns(column, true); err != nil {
		return
	}
	args := map[string]string{}
	if args, err = getArgs(column); err != nil {
		return
	}
	sql = fmt.Sprintf("CREATE TABLE `%s` (\n", table)
	keyList := []string{}
	for _, col := range columns {
		if col["KeyType"] != "" {
			keySql := ""
			if keySql, err = packageKeySql(table, col); err != nil {
				return
			}
			if keySql != "" {
				keyList = append(keyList, fmt.Sprintf("  %s,\n", keySql))
			}
		}
		if col["Type"] != "" {
			sql += fmt.Sprintf("  `%s`", col["Field"])
			pkgSql := ""
			if pkgSql, err = packageFieldSql(col); err != nil {
				return
			}
			sql += pkgSql
			sql += ",\n"
		}
	}
	for _, key := range keyList {
		sql += key
	}
	sql = strings.TrimSuffix(sql, ",\n")
	sql += "\n)"
	engine := "InnoDB"
	if args["Engine"] != "" {
		engine = args["Engine"]
	}
	sql += " ENGINE=" + engine
	charset := "utf8"
	if args["Charset"] != "" {
		charset = args["Charset"]
	}
	sql += " DEFAULT CHARSET=" + charset
	collate := ""
	if args["Collate"] != "" {
		collate = args["Collate"]
	}
	if collate != "" {
		sql += " COLLATE=" + collate
	}
	comment := ""
	if args["Comment"] != "" {
		comment = args["Comment"]
	}
	if comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", comment)
	}
	return
}

func getModifyTableSql(table string, column string) (rs []string, err error) {
	columns := []map[string]string{}
	if columns, err = getColumns(column, false); err != nil {
		return
	}
	for _, col := range columns {
		if col["AlterKeyType"] != "" && (col["Key"] != "" || col["KeyType"] != "") {
			sql := ""
			if col["AlterKeyType"] == "Add" {
				key := ""
				if key, err = packageKeySql(table, col); err != nil {
					return
				}
				sql = fmt.Sprintf("ALTER TABLE `%s` ADD %s", table, key)
			} else if col["AlterKeyType"] == "Drop" {
				if col["KeyType"] == "PRIMARY" {
					col["KeyType"] = col["KeyType"] + " KEY"
				} else if col["KeyType"] == "FOREIGN" {
					col["KeyType"] = fmt.Sprintf("CONSTRAINT `%s`", col["KeyConstraint"])
				} else {
					col["KeyType"] = fmt.Sprintf("INDEX `%s`", strings.ReplaceAll(col["Key"], ",", "_"))
				}
				sql = fmt.Sprintf("ALTER TABLE `%s` DROP %s", table, col["KeyType"])
			}
			rs = append(rs, sql)
			continue
		}
		if col["Field"] == "" || col["AlterFieldType"] == "" {
			continue
		}
		sql := fmt.Sprintf("ALTER TABLE `%s` %s COLUMN `%s`",
			table, strings.ToUpper(col["AlterFieldType"]), col["Field"])
		if col["AlterFieldType"] == "Change" {
			if col["ChangeField"] == "" {
				col["ChangeField"] = col["Field"]
			}
			sql += fmt.Sprintf(" `%s`", col["ChangeField"])
		}
		pkgSql := ""
		if pkgSql, err = packageFieldSql(col); err != nil {
			return
		}
		sql += pkgSql
		rs = append(rs, sql)
	}
	return
}

func getArgs(column string) (rs map[string]string, err error) {
	columnValue := reflect.ValueOf(Args{})
	fieldMap := map[string]interface{}{}
	for i := 0; i < columnValue.NumField(); i++ {
		fieldMap[columnValue.Type().Field(i).Name] = ""
	}
	rs = map[string]string{}
	regString := `migrate.Args *?\{((?s).*?)\}`
	reg := regexp.MustCompile(regString)
	list := reg.FindStringSubmatch(column)
	if len(list) == 0 {
		err = Error("Args does not exist")
		return
	}
	new(Line).SetLine(list[0])
	argsString := list[1]
	for _, argString := range strings.Split(argsString, ",") {
		argString = strings.Trim(argString, "\n\t ")
		if argString == "" {
			continue
		}
		argList := strings.Split(argString, ":")
		if len(argList) < 2 {
			continue
		}
		key := strings.Trim(argList[0], "\n\t ")
		if fieldMap[key] == nil {
			err = Error("Args '%s' not found", key)
			return
		}
		val := strings.Trim(argList[1], "\n\t\" ")
		rs[key] = val
	}
	return
}

func getColumns(column string, args bool) (rs []map[string]string, err error) {
	typeMap := getStructField(Type)
	keyTypeMap := getStructField(KeyType)
	keyFuncMap := getStructField(KeyFunc)
	alterFieldType := getStructField(AlterFieldType)
	alterKeyType := getStructField(AlterKeyType)
	columnValue := reflect.ValueOf(Column{})
	fieldMap := map[string]interface{}{}
	for i := 0; i < columnValue.NumField(); i++ {
		fieldMap[columnValue.Type().Field(i).Name] = ""
	}
	regString := `\[\]migrate.Column *?\{((?s).*)\}\)`
	if args {
		regString = `\[\]migrate.Column *?\{((?s).*)\} *?, *?migrate.Args`
	}
	myLine := Line{}
	reg := regexp.MustCompile(regString)
	list := reg.FindStringSubmatch(column)
	if len(list) == 0 {
		err = Error("Column does not exist")
		return
	}
	myLine.SetLine(list[0])
	regString = `\{((?s).*?)\}`
	reg = regexp.MustCompile(regString)
	for _, col := range reg.FindAllStringSubmatch(list[1], -1) {
		myLine.SetLine(col[0])
		list = []string{}
		mList := strings.Split(col[1], ":")
		for k, v := range mList {
			if k == 0 || k == len(mList)-1 {
				list = append(list, v)
			} else {
				tmpList := strings.Split(v, ",")
				list = append(list, strings.Join(tmpList[0:len(tmpList)-1], ","))
				list = append(list, strings.Join(tmpList[len(tmpList)-1:], ","))
			}
		}
		res := map[string]string{}
		for k, v := range list {
			if k%2 == 0 {
				key := strings.Trim(v, "\n\t ")
				if fieldMap[key] == nil {
					err = Error("Column '%s' not found", key)
					return
				}
				val := strings.Trim(list[k+1], "\n\t, ")
				if key != "Length" && key != "DecimalPoint" && key != "Null" &&
					key != "Unsigned" && key != "AutoIncrement" {
					val = strings.Trim(val, "\"")
				}
				if key == "Default" && val == "" {
					val = "''"
				}
				res[key] = val
			}
		}
		err = valid{
			&res,
			typeMap,
			keyTypeMap,
			keyFuncMap,
			alterFieldType,
			alterKeyType,
		}.ValidCol()
		if err != nil {
			return
		}
		rs = append(rs, res)
	}
	return
}

func packageFieldSql(col map[string]string) (sql string, err error) {
	if col["Type"] == "" {
		return
	}
	sql = " " + col["Type"]
	length := ""
	if col["Length"] != "" {
		length = col["Length"]
	}
	if col["DecimalPoint"] != "" {
		if length != "" {
			length += "," + col["DecimalPoint"]
		}
	}
	if length != "" {
		sql += fmt.Sprintf("(%s)", length)
	}
	if col["Unsigned"] == "true" {
		sql += " unsigned"
	}
	if col["Null"] == "true" {
		sql += " NULL"
	} else if col["Null"] == "false" {
		sql += " NOT NULL"
	}
	if col["Default"] != "" {
		if col["Default"] != "''" {
			col["Default"] = fmt.Sprintf("'%s'", col["Default"])
		}
		sql += fmt.Sprintf(" DEFAULT %s", col["Default"])
	}
	if col["AutoIncrement"] == "true" {
		sql += " AUTO_INCREMENT"
	}
	if col["Comment"] != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", col["Comment"])
	}
	if col["AlterFieldFirst"] == "true" {
		sql += " FIRST"
	}
	if col["AlterFieldAfter"] != "" {
		sql += fmt.Sprintf(" AFTER `%s`", col["AlterFieldAfter"])
	}
	return
}

func packageKeySql(table string, col map[string]string) (sql string, err error) {
	col["Key"] = strings.ReplaceAll(col["Key"], " ", "")
	sql = fmt.Sprintf("%s KEY", col["KeyType"])
	if col["KeyType"] == "NORMAL" || col["KeyType"] == "" {
		sql = "KEY"
	}
	if col["KeyType"] != "PRIMARY" {
		if col["Key"] == "" {
			col["Key"] = strings.ReplaceAll(col["Field"], ",", "_")
		}
		sql += fmt.Sprintf(" `%s`", col["Key"])
	}
	sql += fmt.Sprintf(" (`%s`)", strings.ReplaceAll(col["Field"], ",", "`,`"))
	if col["KeyFunc"] != "" {
		sql += " USING " + col["KeyFunc"]
	}
	if col["KeyType"] == "FOREIGN" {
		if col["KeyConstraint"] == "" {
			col["KeyConstraint"] = fmt.Sprintf("%s_%s_%s", table, col["KeyRelationTable"], col["Key"])
		}
		sql = fmt.Sprintf("CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`)", col["KeyConstraint"],
			strings.ReplaceAll(col["Key"], ",", "`,`"), col["KeyRelationTable"],
			strings.ReplaceAll(col["KeyRelationField"], ",", "`,`"))
	}
	return
}

func getStructField(dest interface{}) map[string]interface{} {
	destValue := reflect.ValueOf(dest)
	fieldMap := map[string]interface{}{}
	for i := 0; i < destValue.NumField(); i++ {
		fieldMap[destValue.Type().Field(i).Name] = ""
	}
	return fieldMap
}

func readAll(filePath string) (string, error) {
	fi, err := os.Open(filePath)
	defer fi.Close()
	if err != nil {
		return "", err
	}
	by, err := ioutil.ReadAll(fi)
	if err != nil {
		return "", err
	}
	content := string(by)
	reg := regexp.MustCompile(`//.*`)
	content = reg.ReplaceAllString(content, "")
	reg = regexp.MustCompile(`/\*(?s).*?\*/`)
	content = reg.ReplaceAllString(content, "")
	return content, nil
}

func Error(err string, args ...interface{}) error {
	err = fmt.Sprintf(err, args...)
	return fmt.Errorf("file: %s, line: %d, Error: %s", FilePath, LastLineNum, err)
}
