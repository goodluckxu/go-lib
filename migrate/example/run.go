package main

import (
	"fmt"
	"github.com/goodluckxu/go-lib/migrate"
)

func main() {
	sqlList, err := migrate.GetSql("./table.go", migrate.RunType.Up)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, sql := range sqlList {
		fmt.Println(sql)
	}
}
