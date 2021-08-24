package migrate

type Column struct {
	Field            string // 字段名
	Type             uint8  // 字段类型
	Length           int    // 长度
	DecimalPoint     int    // 小数点
	Null             bool   // 是null
	Comment          string // 注释
	Default          string // 默认值
	Unsigned         bool   // 无符号
	AutoIncrement    bool   // 自动递增
	Key              string // 键
	KeyType          uint8  // 索引类型
	KeyFunc          uint8  // 索引方法
	KeyRelationTable string // 外键关联表
	KeyRelationField string // 外键关联字段
	KeyConstraint    string // 约束
	AlterFieldType   uint8  // 字段操作类型
	ChangeField      string // 修改字段
	AlterKeyType     uint8  // 索引操作类型
}

type Args struct {
	Charset string // 编码
	Collate string // 校对规则
	Comment string // 注释
	Engine  string // 引擎
}

type fieldType struct {
	Tinyint    uint8
	Smallint   uint8
	Int        uint8
	Bigint     uint8
	Decimal    uint8
	Float      uint8
	Char       uint8
	Varchar    uint8
	Tinytext   uint8
	Mediumtext uint8
	Text       uint8
	Longtext   uint8
	Date       uint8
	Time       uint8
	Year       uint8
	Datetime   uint8
	Timestamp  uint8
	Enum       uint8
	Json       uint8
}

type keyType struct {
	Primary  uint8 // 主键
	Normal   uint8 // 常规
	Fulltext uint8 // 全文
	Spatial  uint8 // 空间
	Unique   uint8 // 唯一
	Foreign  uint8 // 外键
}

type keyFunc struct {
	Hash  uint8
	Btree uint8
}

type alterFieldType struct {
	Add    uint8 // 添加字段
	Modify uint8 // 修改字段
	Change uint8 // 修改字段(可修改字段名)
	Drop   uint8 // 删除字段
}

type alterKeyType struct {
	Add  uint8 // 添加
	Drop uint8 // 删除
}

type runType struct {
	Up   uint8
	Down uint8
}

var RunType runType
var Type fieldType
var KeyType keyType
var KeyFunc keyFunc
var AlterFieldType alterFieldType
var AlterKeyType alterKeyType
var FilePath string
var LastLineNum int

func init() {
	RunType = runType{
		Up:   1,
		Down: 2,
	}
}
