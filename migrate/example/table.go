package main

import "github.com/goodluckxu/go-lib/migrate"

type Table struct {
}

func (t Table) Up() {
	migrate.CreateTable("test", []migrate.Column{
		{Field: "id", Type: migrate.Type.Int, Length: 10, Null: false, Comment: "主键", Unsigned: true, AutoIncrement: true},
		{Field: "name", Type: migrate.Type.Varchar, Null: true, Default: "", Length: 20, Comment: "名称"},
		{Field: "id", KeyType: migrate.KeyType.Primary},
		{Field: "name", KeyType: migrate.KeyType.Unique, KeyFunc: migrate.KeyFunc.Btree},
		{Key: "name", KeyType: migrate.KeyType.Foreign, KeyRelationTable: "user", KeyRelationField: "id,name"},
	}, migrate.Args{
		Engine:  "InnoDB",
		Charset: "utf8mb4",
		Collate: "utf8mb4_unicode_ci",
		Comment: "用户表",
	})
	migrate.DropTable("test")
	migrate.ModifyTable("test", []migrate.Column{
		{
			AlterFieldType: migrate.AlterFieldType.Add,
			Field:          "icon",
			Type:           migrate.Type.Varchar,
			Length:         255,
			Comment:        "图标",
		},
		{AlterFieldType: migrate.AlterFieldType.Modify, Field: "icon", Type: migrate.Type.Varchar, Length: 255, Comment: "图标"},
		{AlterFieldType: migrate.AlterFieldType.Change, Field: "icon", ChangeField: "icon_change", Type: migrate.Type.Varchar, Length: 255, Comment: "图标"},
		{AlterFieldType: migrate.AlterFieldType.Drop, Field: "icon"},
		{AlterFieldType: migrate.AlterFieldType.Modify, Field: "id", Type: migrate.Type.Int},
		{AlterKeyType: migrate.AlterKeyType.Drop, KeyType: migrate.KeyType.Primary},
		{AlterKeyType: migrate.AlterKeyType.Drop, Key: "id"},
		{AlterKeyType: migrate.AlterKeyType.Add, Field: "id", KeyType: migrate.KeyType.Normal},
		{AlterKeyType: migrate.AlterKeyType.Add, KeyType: migrate.KeyType.Foreign, Key: "name", KeyRelationTable: "admin",
			KeyRelationField: "id", KeyConstraint: "aaa"},
		{AlterKeyType: migrate.AlterKeyType.Drop, KeyType: migrate.KeyType.Foreign, KeyConstraint: "aaa"},
	})
}

func (t Table) Down() {
	migrate.DropTable("aaa")
}
