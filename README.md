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
~~~
CREATE TABLE `test` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(20) NULL DEFAULT '' COMMENT '名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE,
  CONSTRAINT `test_user_name` FOREIGN KEY (`name`) REFERENCES `user` (`id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT '用户表'
DROP TABLE `test`
ALTER TABLE `test` ADD COLUMN `icon` varchar(255) COMMENT '图标'
ALTER TABLE `test` MODIFY COLUMN `icon` varchar(255) COMMENT '图标'
ALTER TABLE `test` CHANGE COLUMN `icon` `icon_change` varchar(255) COMMENT '图标'
ALTER TABLE `test` DROP COLUMN `icon`
ALTER TABLE `test` MODIFY COLUMN `id` int
ALTER TABLE `test` DROP PRIMARY KEY
ALTER TABLE `test` DROP INDEX `id`
ALTER TABLE `test` ADD KEY `id` (`id`)
ALTER TABLE `test` ADD CONSTRAINT `aaa` FOREIGN KEY (`name`) REFERENCES `admin` (`id`)
ALTER TABLE `test` DROP CONSTRAINT `aaa`
~~~