package main

import "github.com/goodluckxu/go-lib/migrate"

type MyTable struct {
}

func (m MyTable) Up() {
	migrate.CreateTable("user", []migrate.Column{
		// 字段
		{Field: "id", Type: migrate.Type.Int, Length: 10, Unsigned: true, AutoIncrement: true, Null: false,
			Comment: "主键"},
		{Field: "username", Type: migrate.Type.Varchar, Length: 255, Null: false, Default: "", Comment: "用户名"},
		{Field: "nickname", Type: migrate.Type.Varchar, Length: 255, Null: true, Comment: "昵称"},
		{Field: "account", Type: migrate.Type.Varchar, Length: 20, Null: false, Comment: "账号"},
		{Field: "status", Type: migrate.Type.Tinyint, Length: 1, Null: false, Comment: "状态"},
		{Field: "price", Type: migrate.Type.Decimal, Length: 10, DecimalPoint: 2, Null: false, Default: "0",
			Comment: "价格"},
		// 主键
		{Field: "id", KeyType: migrate.KeyType.Primary},
		// 索引
		{Field: "username", Key: "my_foreign", KeyType: migrate.KeyType.Normal},
		{Field: "nickname", KeyType: migrate.KeyType.Fulltext},
		{Field: "account,status", Key: "my_unique", KeyType: migrate.KeyType.Unique},
		// 外键，添加外键必须先添加索引
		{Key: "my_foreign", KeyType: migrate.KeyType.Foreign, KeyRelationTable: "center_user",
			KeyRelationField: "nick_name"},
	}, migrate.Args{
		Engine:  "InnoDB",
		Charset: "utf8mb4",
		Collate: "utf8mb4_unicode_ci",
		Comment: "用户表",
	})
	migrate.DropTable("test")
	migrate.ModifyTable("user_info", []migrate.Column{
		// 添加字段
		{AlterFieldType: migrate.AlterFieldType.Add, Field: "icon", Type: migrate.Type.Varchar, Length: 255,
			Comment: "图标", AlterFieldFirst: true},
		// 修改字段，不能修改字段名
		{AlterFieldType: migrate.AlterFieldType.Modify, Field: "icon", Type: migrate.Type.Varchar, Length: 255,
			Comment: "图标", AlterFieldAfter: "id"},
		// 修改字段，可以修改字段名
		{AlterFieldType: migrate.AlterFieldType.Change, Field: "icon", ChangeField: "icon_change",
			Type: migrate.Type.Varchar, Length: 255, Comment: "图标"},
		// 删除字段
		{AlterFieldType: migrate.AlterFieldType.Drop, Field: "icon"},
		// 添加主键
		{AlterKeyType: migrate.AlterKeyType.Add, Field: "id", KeyType: migrate.KeyType.Primary},
		// 删除主键
		{AlterKeyType: migrate.AlterKeyType.Drop, KeyType: migrate.KeyType.Primary},
		// 添加索引
		{AlterKeyType: migrate.AlterKeyType.Add, Field: "test", Key: "test_1", KeyType: migrate.KeyType.Unique},
		// 删除索引
		{AlterKeyType: migrate.AlterKeyType.Drop, Key: "test_1"},
		// 添加外键
		{AlterKeyType: migrate.AlterKeyType.Add, KeyType: migrate.KeyType.Foreign, Key: "name", KeyRelationTable: "admin",
			KeyRelationField: "id", KeyConstraint: "aaa"},
		// 删除外键
		{AlterKeyType: migrate.AlterKeyType.Drop, KeyType: migrate.KeyType.Foreign, KeyConstraint: "aaa"},
	})
}

func (m MyTable) Down() {
	migrate.DropTable("user")
}