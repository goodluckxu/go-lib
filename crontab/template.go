package crontab

import "time"

var (
	crontabType = struct {
		Minute uint8
		Hour   uint8
		Day    uint8
		Month  uint8
		Week   uint8
	}{1, 2, 3, 4, 5}
)

type BeforeTime struct {
	Time         time.Time // 上一次执行的时间
	CompareTypes []uint8   // 为空则判断全部类型(分钟,小时,天,月,周)
}

type Times struct {
	Interval int64 // 间隔时间
	Start    int64
	End      int64
}

type Between struct {
	Start int64
	End   int64
}
