package migrate

import (
	"strconv"
	"strings"
)

type valid struct {
	col            *map[string]string
	typeMap        map[string]interface{}
	keyTypeMap     map[string]interface{}
	keyFuncMap     map[string]interface{}
	alterFieldType map[string]interface{}
	alterKeyType   map[string]interface{}
}

// 验证Column
func (v valid) ValidCol() error {
	if err := v.validKeyCol(); err != nil {
		return err
	}
	if err := v.validFieldCol(); err != nil {
		return err
	}
	if err := v.validLogicCol(); err != nil {
		return err
	}
	return nil
}

// 验证Column的键
func (v valid) validKeyCol() error {
	col := v.col
	keyTypeMap := v.keyTypeMap
	keyFuncMap := v.keyFuncMap
	alterKeyType := v.alterKeyType
	newCol := *col
	// 验证键类型是否合法
	if newCol["KeyType"] != "" {
		if strings.Index(newCol["KeyType"], "migrate.KeyType.") == -1 {
			return Error("Column 'KeyType' does not exist in migrate.KeyType")
		}
		newCol["KeyType"] = strings.TrimPrefix(newCol["KeyType"], "migrate.KeyType.")
		if keyTypeMap[newCol["KeyType"]] == nil {
			return Error("Column 'KeyType' does not exist in migrate.KeyType")
		}
		newCol["KeyType"] = strings.ToUpper(newCol["KeyType"])
	}
	// 验证键方法是否合法
	if newCol["KeyFunc"] != "" {
		if strings.Index(newCol["KeyFunc"], "migrate.KeyFunc.") == -1 {
			return Error("Column 'KeyFunc' does not exist in migrate.KeyFunc")
		}
		newCol["KeyFunc"] = strings.TrimPrefix(newCol["KeyFunc"], "migrate.KeyFunc.")
		if keyFuncMap[newCol["KeyFunc"]] == nil {
			return Error("Column 'KeyFunc' does not exist in migrate.KeyFunc")
		}
		newCol["KeyFunc"] = strings.ToUpper(newCol["KeyFunc"])
	}
	// 验证修改键类型是否合法
	if newCol["AlterKeyType"] != "" {
		if strings.Index(newCol["AlterKeyType"], "migrate.AlterKeyType.") == -1 {
			return Error("Column 'AlterKeyType' does not exist in migrate.AlterKeyType")
		}
		newCol["AlterKeyType"] = strings.TrimPrefix(newCol["AlterKeyType"], "migrate.AlterKeyType.")
		if alterKeyType[newCol["AlterKeyType"]] == nil {
			return Error("Column 'AlterKeyType' does not exist in migrate.AlterKeyType")
		}
	}
	return nil
}

// 验证Column的字段
func (v valid) validFieldCol() error {
	col := v.col
	typeMap := v.typeMap
	alterFieldType := v.alterFieldType
	newCol := *col
	// 验证字段类型是否合法
	if newCol["Type"] != "" {
		if strings.Index(newCol["Type"], "migrate.Type.") == -1 {
			return Error("Column 'Type' does not exist in migrate.Type")
		}
		newCol["Type"] = strings.TrimPrefix(newCol["Type"], "migrate.Type.")
		if typeMap[newCol["Type"]] == nil {
			return Error("Column 'Type' does not exist in migrate.Type")
		}
		newCol["Type"] = strings.ToLower(newCol["Type"])
	}
	// 验证修改字段类型是否合法
	if newCol["AlterFieldType"] != "" {
		if strings.Index(newCol["AlterFieldType"], "migrate.AlterFieldType.") == -1 {
			return Error("Column 'AlterFieldType' does not exist in migrate.AlterFieldType")
		}
		newCol["AlterFieldType"] = strings.TrimPrefix(newCol["AlterFieldType"], "migrate.AlterFieldType.")
		if alterFieldType[newCol["AlterFieldType"]] == nil {
			return Error("Column 'AlterFieldType' does not exist in migrate.AlterFieldType")
		}
	}
	// 验证长度为数字
	if newCol["Length"] != "" {
		if _, err := strconv.ParseInt(newCol["Length"], 10, 64); err != nil {
			return Error("Column 'Length' must be an integer")
		}
	}
	// 验证小数点
	if newCol["DecimalPoint"] != "" {
		if newCol["Length"] == "" {
			return Error("Column 'DecimalPoint' must have 'Length'")
		}
		if _, err := strconv.ParseInt(newCol["DecimalPoint"], 10, 64); err != nil {
			return Error("Column 'DecimalPoint' must be an integer")
		}
	}
	// 验证Null为布尔
	if newCol["Null"] != "" && newCol["Null"] != "true" && newCol["Null"] != "false" {
		return Error("Column 'Null' must be Boolean")
	}
	// 验证无符号为布尔
	if newCol["Unsigned"] != "" && newCol["Unsigned"] != "true" && newCol["Unsigned"] != "false" {
		return Error("Column 'Unsigned' must be Boolean")
	}
	// 验证自动递增为布尔
	if newCol["AutoIncrement"] != "" && newCol["AutoIncrement"] != "true" && newCol["AutoIncrement"] != "false" {
		return Error("Column 'AutoIncrement' must be Boolean")
	}
	return nil
}

// 验证Column的逻辑
func (v valid) validLogicCol() error {
	newCol := *v.col
	// 必须存在某种类型
	if newCol["AlterFieldType"] == "" && newCol["AlterKeyType"] == "" && newCol["KeyType"] == "" &&
		newCol["Type"] == "" {
		return Error("Column 'AlterFieldType', 'AlterKeyType', 'KeyType' and 'Type' must have one")
	}
	// 字段类型和键类型不能同时存在
	if newCol["Type"] != "" && newCol["KeyType"] != "" {
		return Error("Column 'Type' and 'KeyType' cannot exist at the same time")
	}
	// 判断字段必须存在
	if (newCol["Type"] != "" || newCol["AlterFieldType"] != "" ||
		(newCol["AlterKeyType"] == "Add" && newCol["KeyType"] != "FOREIGN")) && newCol["Field"] == "" {
		return Error("Column 'Field' must exist")
	}
	// 判断字段类型必须存在
	if newCol["AlterFieldType"] == "Add" && newCol["Type"] == "" {
		return Error("Column 'Type' must exist")
	}
	// 添加外键的时候
	if (newCol["AlterKeyType"] == "Add" || newCol["AlterKeyType"] == "") && newCol["KeyType"] == "FOREIGN" &&
		(newCol["Key"] == "" || newCol["KeyRelationTable"] == "" || newCol["KeyRelationField"] == "") {
		return Error("Column 'Key', 'KeyRelationTable' and 'KeyRelationField' must exist")
	}
	// 删除外键的时候
	if newCol["AlterKeyType"] == "Drop" && newCol["KeyType"] == "FOREIGN" && newCol["KeyConstraint"] == "" {
		return Error("Column 'KeyConstraint' must exist")
	}
	return nil
}
