# lib
模块 [handle_interface](#handle_interface),[crontab](#crontab),[migrate](#migrate)

## 模块 <a id="handle_interface">handle_interface</a>
~~~go
jsonStr = `[
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/storage/upload/20210602/fe065f1b894984259724806fd82a5c2b.jpg",
        "description": "1",
        "id": 1,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 2,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 3,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 4,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 5,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 6,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 7,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 8,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/storage/upload/20210602/fe065f1b894984259724806fd82a5c2b.jpg",
        "description": "1",
        "id": 9,
        "subject": "测试",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 10,
        "subject": "1",
        "twelve_hours_activity": [
            {
                "info": {
                    "cover": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/111",
                    "id": 1,
                    "title": "喝水"
                },
                "twelve_hours_activity_id": 1
            }
        ]
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 11,
        "subject": "1",
        "twelve_hours_activity": []
    },
    {
        "background_img": "https://350sz-oss.oss-cn-shanghai.aliyuncs.com/1",
        "description": "1",
        "id": 12,
        "subject": "1",
        "twelve_hours_activity": []
    }
]`
var jsonObj interface{}
err := json.Unmarshal([]byte(jsonStr), &jsonObj)
if err != nil {
    fmt.Println(err)
    os.Exit(0)
}
b := handle_interface.GetInterface("*.twelve_hours_activity")

a := handle_interface.UpdateInterface(jsonObj, []handle_interface.Rule{
    {
        FindField:   "*.twelve_hours_activity",
        UpdateValue: "*.twelve_hours_activity.*.info",
        Type:        "*",
    },
})

fmt.Println(b)
fmt.Println(jsonObj)
fmt.Println(a)
~~~

## 模块 <a id="crontab">crontab</a>
~~~go
fmt.Println(crontab.New().IsRun("0,18 */1 * * */2", crontab.BeforeTime{
    Time: time.Now().Add(30 * time.Minute),
    CompareTypes: []uint8{
        crontab.CrontabType.Hour,
        crontab.CrontabType.Minute,
    },
}))
~~~

## 模块 <a id="migrate">migrate</a>
用法: key相关的为键的，目前只支持mysql生成

运行 crontab/example的实例获取数据
~~~sql
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `nickname` varchar(255) NULL COMMENT '昵称',
  `account` varchar(20) NOT NULL COMMENT '账号',
  `status` tinyint(1) NOT NULL COMMENT '状态',
  `price` decimal(10,2) NOT NULL DEFAULT '0' COMMENT '价格',
  PRIMARY KEY (`id`),
  KEY `my_foreign` (`username`),
  FULLTEXT KEY `nickname` (`nickname`),
  UNIQUE KEY `my_unique` (`account`,`status`),
  CONSTRAINT `user_center_user_my_foreign` FOREIGN KEY (`my_foreign`) REFERENCES `center_user` (`nick_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT '用户表'
DROP TABLE `test`
ALTER TABLE `user_info` ADD COLUMN `icon` varchar(255) COMMENT '图标'
ALTER TABLE `user_info` MODIFY COLUMN `icon` varchar(255) COMMENT '图标'
ALTER TABLE `user_info` CHANGE COLUMN `icon` `icon_change` varchar(255) COMMENT '图标'
ALTER TABLE `user_info` DROP COLUMN `icon`
ALTER TABLE `user_info` ADD PRIMARY KEY (`id`)
ALTER TABLE `user_info` DROP PRIMARY KEY
ALTER TABLE `user_info` ADD UNIQUE KEY `test_1` (`test`)
ALTER TABLE `user_info` DROP INDEX `test_1`
ALTER TABLE `user_info` ADD CONSTRAINT `aaa` FOREIGN KEY (`name`) REFERENCES `admin` (`id`)
ALTER TABLE `user_info` DROP CONSTRAINT `aaa`
~~~