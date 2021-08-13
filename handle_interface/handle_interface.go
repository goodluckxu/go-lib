package handle_interface

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Rule struct {
	// 需要修改的字段  .表示层级，*表示数组，如：a.*.b.c表示，a的map下面有数组，数组里的map有b，b的map有c，*可以用数字代表数组下标
	FindField string
	// 修改的内容，当Type为*时，UpdateValue如果填写内部字段的层级，则表示吧该内部字段赋予需要修改的字段，
	// 如果是 匿名函数 func(value interface{}) interface{}，则表示处理需要修改的字段本身内容
	UpdateValue interface{}
	// 类型 *表示修改内部字段，_表示删除字段
	Type string
}

// UpdateInterface 修改通用interface的内容
//    data    需要修改的interface
//    updates []Rule{{FindField:"",UpdateValue:interface{}{},Type:""}}
//    result   返回修改的内容
func UpdateInterface(data interface{}, updates []Rule) interface{} {
	dataByte, _ := json.Marshal(data)
	var newData interface{}
	_ = json.Unmarshal(dataByte, &newData)
	for _, update := range updates {
		switch update.Type {
		case "*":
			switch update.UpdateValue.(type) {
			case string:
				newData = updateInsideInterface(newData, update.FindField, update.UpdateValue.(string))
			case func(value interface{}) interface{}:
				newData = updateInsideFunInterface(newData, update.FindField, update.UpdateValue.(func(value interface{}) interface{}))
			}
		case "_":
			newData = delInterfaceField(newData, update.FindField)
		default:
			newData = updateUniversalInterface(newData, update.FindField, update.UpdateValue)
		}
	}
	return newData
}

// GetInterface  获取interface内容
//    data       原来的interface
//    findField  需要获取的interface
//       1. 里面的.表示层级，*表示数组，如：a.*.b.c表示，a的map下面有数组，数组里的map有b，b的map有c
//       2. 如果获取使用*获取内容，则将*下面的内容合并成一个数组，如果*下面的内容不为数组，则返回一位数组
func GetInterface(data interface{}, findField string) interface{} {
	dataByte, _ := json.Marshal(data)
	var newData interface{}
	_ = json.Unmarshal(dataByte, &newData)
	return getUniversalInterface(newData, findField)
}

func getUniversalInterface(data interface{}, findField string) interface{} {
	if data == nil {
		return data
	}
	findFieldList := strings.Split(findField, ".")
	isArrayNum := isInt(findFieldList[0]) && isSlice(data)
	if len(findFieldList) > 1 {
		if findFieldList[0] == "*" {
			newData := []interface{}{}
			newDataList, _ := data.([]interface{})
			if len(newDataList) == 0 {
				return nil
			}
			for _, v := range newDataList {
				newV := getUniversalInterface(
					v, strings.Join(findFieldList[1:], "."))
				if newV == nil {
					continue
				}
				if len(findFieldList) > 2 && isSlice(newV) {
					nNewV, _ := newV.([]interface{})
					for _, val := range nNewV {
						newData = append(newData, val)
					}
				} else {
					newData = append(newData, newV)
				}
			}
			return newData
		} else {
			if isArrayNum {
				num, _ := strconv.ParseInt(findFieldList[0], 10, 64)
				newDataList, _ := data.([]interface{})
				if len(newDataList) > int(num) {
					return getUniversalInterface(
						newDataList[num], strings.Join(findFieldList[1:], "."))
				}
				return nil
			} else {
				newDataMap, _ := data.(map[string]interface{})
				return getUniversalInterface(
					newDataMap[findFieldList[0]], strings.Join(findFieldList[1:], "."))
			}
		}
	} else {
		if findFieldList[0] != "*" {
			if isArrayNum {
				num, _ := strconv.ParseInt(findFieldList[0], 10, 64)
				newDataList, _ := data.([]interface{})
				if len(newDataList) > int(num) {
					return newDataList[num]
				}
				return nil
			} else {
				newDataMap, _ := data.(map[string]interface{})
				return newDataMap[findFieldList[0]]
			}
		}
	}
	return data
}

// 修改内部数据
func updateInsideInterface(data interface{}, findField string, updateValue string) interface{} {
	dataByte, _ := json.Marshal(data)
	var newData interface{}
	_ = json.Unmarshal(dataByte, &newData)
	findFieldList := strings.Split(findField, ".")
	updateValueList := strings.Split(updateValue, ".")
	commonFieldList := []string{}
	if len(findFieldList) < len(updateValueList) {
		for i := 0; i < len(findFieldList); i++ {
			if findFieldList[i] != updateValueList[i] {
				break
			}
			commonFieldList = append(commonFieldList, findFieldList[i])
		}
	} else {
		for i := 0; i < len(updateValueList); i++ {
			if findFieldList[i] != updateValueList[i] {
				break
			}
			commonFieldList = append(commonFieldList, findFieldList[i])
		}
	}
	commonField := strings.Join(commonFieldList, ".")
	findField = strings.Join(findFieldList[len(commonFieldList):], ".")
	updateValue = strings.Join(updateValueList[len(commonFieldList):], ".")
	newCommonFieldList := updateInsideCommonInterface(data, commonField, []string{})
	if len(newCommonFieldList) > 0 {
		for _, newCommonField := range newCommonFieldList {
			newFindField := newCommonField
			if findField != "" {
				newFindField += "." + findField
			}
			newUpdateValue := newCommonField
			if updateValue != "" {
				newUpdateValue += "." + updateValue
			}
			newData = updateUniversalInterface(newData, newFindField, GetInterface(data, newUpdateValue))
		}
	} else {
		newData = updateUniversalInterface(newData, commonField+"."+findField, GetInterface(data, updateValue))
	}
	return newData
}

func updateInsideCommonInterface(data interface{}, commonField string, list []string) []string {
	pIndex := strings.Index(commonField, ".")
	sIndex := strings.Index(commonField, "*")
	if sIndex == -1 && commonField != "" {
		list = append(list, commonField)
	}
	if sIndex != -1 {
		var childList interface{}
		if pIndex != -1 {
			childList = GetInterface(data, commonField[0:sIndex+1])
		} else {
			childList = GetInterface(data, commonField)
		}
		newChildList, _ := childList.([]interface{})
		num := len(newChildList)
		for i := 0; i < num; i++ {
			strI := strconv.FormatInt(int64(i), 10)
			newCommonField := strings.Replace(commonField, "*", strI, 1)
			list = updateInsideCommonInterface(data, newCommonField, list)
		}
	}
	return list
}

// 修改内部数据，自身使用函数赋值
func updateInsideFunInterface(data interface{}, findField string, updateValue func(value interface{}) interface{}) interface{} {
	findFieldList := strings.Split(findField, ".")
	isArrayNum := isInt(findFieldList[0]) && isSlice(data)
	if len(findFieldList) > 1 {
		if findFieldList[0] == "*" {
			newData := []interface{}{}
			newDataList, _ := data.([]interface{})
			for _, v := range newDataList {
				dataTmp := updateInsideFunInterface(v, strings.Join(findFieldList[1:], "."), updateValue)
				newData = append(newData, dataTmp)
			}
			data = newData
		} else {
			if isArrayNum {
				num, _ := strconv.ParseInt(findFieldList[0], 10, 64)
				newDataList, _ := data.([]interface{})
				if len(newDataList) > int(num) {
					newDataList[num] = updateInsideFunInterface(
						newDataList[num], strings.Join(findFieldList[1:], "."), updateValue)
				}
			} else {
				newDataMap, _ := data.(map[string]interface{})
				if newDataMap == nil {
					newDataMap = map[string]interface{}{}
				}
				newDataMap[findFieldList[0]] = updateInsideFunInterface(
					newDataMap[findFieldList[0]], strings.Join(findFieldList[1:], "."), updateValue)
				data = newDataMap
			}
		}
	} else {
		if findFieldList[0] == "*" {
			newData := []interface{}{}
			newDataList, _ := data.([]interface{})
			for _, v := range newDataList {
				newData = append(newData, updateValue(v))
			}
			data = newData
		} else {
			if isArrayNum {
				num, _ := strconv.ParseInt(findFieldList[0], 10, 64)
				newDataList, _ := data.([]interface{})
				if len(newDataList) > int(num) {
					newDataList[num] = updateValue(newDataList[num])
				}
			} else {
				newDataMap, _ := data.(map[string]interface{})
				if newDataMap == nil {
					newDataMap = map[string]interface{}{}
				}
				newDataMap[findFieldList[0]] = updateValue(newDataMap[findFieldList[0]])
				data = newDataMap
			}
		}
	}
	return data
}

// 删除数据
func delInterfaceField(data interface{}, findField string) interface{} {
	findFieldList := strings.Split(findField, ".")
	isArrayNum := isInt(findFieldList[0]) && isSlice(data)
	if len(findFieldList) > 1 {
		if findFieldList[0] == "*" {
			newData := []interface{}{}
			newDataList, _ := data.([]interface{})
			for _, v := range newDataList {
				newData = append(newData, delInterfaceField(
					v, strings.Join(findFieldList[1:], ".")))
			}
			data = newData
		} else {
			if isArrayNum {
				num, _ := strconv.ParseInt(findFieldList[0], 10, 10)
				newDataList, _ := data.([]interface{})
				if len(newDataList) > int(num) {
					newDataList[num] = delInterfaceField(
						newDataList[num], strings.Join(findFieldList[1:], "."))
				}
			} else {
				newDataMap, _ := data.(map[string]interface{})
				if newDataMap == nil {
					newDataMap = map[string]interface{}{}
				}
				newDataMap[findFieldList[0]] = delInterfaceField(
					newDataMap[findFieldList[0]], strings.Join(findFieldList[1:], "."))
				data = newDataMap
			}
		}
	} else {
		if findFieldList[0] == "*" {
			data = []interface{}{}
		} else {
			if isArrayNum {
				num, _ := strconv.ParseInt(findFieldList[0], 10, 64)
				newData := []interface{}{}
				newDataList, _ := data.([]interface{})
				for k, v := range newDataList {
					if int(num) == k {
						continue
					}
					newData = append(newData, v)
				}
				data = newData
			} else {
				newDataMap, _ := data.(map[string]interface{})
				delete(newDataMap, findFieldList[0])
				data = newDataMap
			}
		}
	}
	return data
}

// 修改通用数据
func updateUniversalInterface(data interface{}, findField string, updateValue interface{}) interface{} {
	findFieldList := strings.Split(findField, ".")
	isArrayNum := isInt(findFieldList[0]) && isSlice(data)
	if len(findFieldList) > 1 {
		if findFieldList[0] == "*" {
			newData := []interface{}{}
			newDataList, _ := data.([]interface{})
			for _, v := range newDataList {
				dataTmp := updateUniversalInterface(v, strings.Join(findFieldList[1:], "."), updateValue)
				newData = append(newData, dataTmp)
			}
			data = newData
		} else {
			if isArrayNum {
				num, _ := strconv.ParseInt(findFieldList[0], 10, 64)
				newDataList, _ := data.([]interface{})
				if len(newDataList) > int(num) {
					newDataList[num] = updateUniversalInterface(
						newDataList[num], strings.Join(findFieldList[1:], "."), updateValue)
				}
				data = newDataList
			} else {
				newDataMap, _ := data.(map[string]interface{})
				if newDataMap == nil {
					newDataMap = map[string]interface{}{}
				}
				newDataMap[findFieldList[0]] = updateUniversalInterface(
					newDataMap[findFieldList[0]], strings.Join(findFieldList[1:], "."), updateValue)
				data = newDataMap
			}
		}
	} else {
		if findFieldList[0] == "*" {
			newData := []interface{}{}
			newDataList, _ := data.([]interface{})
			for i := 0; i < len(newDataList); i++ {
				newData = append(newData, updateValue)
			}
			data = newData
		} else {
			if isArrayNum {
				num, _ := strconv.ParseInt(findFieldList[0], 10, 64)
				newDataList, _ := data.([]interface{})
				if len(newDataList) > int(num) {
					newDataList[num] = updateValue
				}
				data = newDataList
			} else {
				newDataMap, _ := data.(map[string]interface{})
				if newDataMap == nil {
					newDataMap = map[string]interface{}{}
				}
				newDataMap[findFieldList[0]] = updateValue
				data = newDataMap
			}
		}
	}
	return data
}

func isInt(value string) bool {
	_, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return true
}

func isSlice(value interface{}) bool {
	_, ok := value.([]interface{})
	return ok
}
