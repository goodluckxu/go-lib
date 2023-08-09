package crontab

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func New() *crontab {
	cron := new(crontab)
	return cron
}

type crontab struct {
}

// IsRun 判断规则当前是否可执行
// 规则为 * * * * * 形式linux模式
func (c crontab) IsRun(rules string, beforeTime *time.Time) bool {
	minuteList, hourList, dayList, monthList, weekList, err := c.parse(rules)
	if err != nil {
		return false
	}
	// 判断当前时间是否满足
	minute, hour, day, month, week := c.timeSeparate(time.Now())
	if !(c.inArrayInt64(minute, minuteList) &&
		c.inArrayInt64(hour, hourList) &&
		c.inArrayInt64(day, dayList) &&
		c.inArrayInt64(month, monthList) &&
		c.inArrayInt64(week, weekList)) {
		return false
	}
	// 判断前时间是否已经运行
	if beforeTime == nil {
		return true
	}
	beforeMinute, beforeHour, beforeDay, beforeMonth, beforeWeek := c.timeSeparate(*beforeTime)
	if minute == beforeMinute && hour == beforeHour &&
		day == beforeDay && month == beforeMonth && week == beforeWeek {
		return false
	}
	return true
}

// 解析* * * * *类型定时任务
func (c crontab) parse(rules string) (minute, hour, day, month, week []int64, err error) {
	ruleList := strings.Split(rules, " ")
	if len(ruleList) != 5 {
		err = errors.New("参数错误")
		return
	}
	minute, err = c.parseMinute(ruleList[0])
	if err != nil {
		minute = []int64{}
		hour = []int64{}
		day = []int64{}
		month = []int64{}
		week = []int64{}
		return
	}
	hour, err = c.parseHour(ruleList[1])
	if err != nil {
		minute = []int64{}
		hour = []int64{}
		day = []int64{}
		month = []int64{}
		week = []int64{}
		return
	}
	day, err = c.parseDay(ruleList[2])
	if err != nil {
		minute = []int64{}
		hour = []int64{}
		day = []int64{}
		month = []int64{}
		week = []int64{}
		return
	}
	month, err = c.parseMonth(ruleList[3])
	if err != nil {
		minute = []int64{}
		hour = []int64{}
		day = []int64{}
		month = []int64{}
		week = []int64{}
		return
	}
	week, err = c.parseWeek(ruleList[4])
	if err != nil {
		minute = []int64{}
		hour = []int64{}
		day = []int64{}
		month = []int64{}
		week = []int64{}
		return
	}
	return
}

func (c crontab) parseMinute(rule string) ([]int64, error) {
	return c.parseSingle(rule, crontabType.Minute)
}

func (c crontab) parseHour(rule string) ([]int64, error) {
	return c.parseSingle(rule, crontabType.Hour)
}

func (c crontab) parseDay(rule string) ([]int64, error) {
	return c.parseSingle(rule, crontabType.Day)
}
func (c crontab) parseMonth(rule string) ([]int64, error) {
	return c.parseSingle(rule, crontabType.Month)
}

func (c crontab) parseWeek(rule string) ([]int64, error) {
	return c.parseSingle(rule, crontabType.Week)
}

func (c crontab) parseSingle(rule string, crontabType uint8) (rs []int64, rsErr error) {
	rsErr = errors.New("参数错误: " + fmt.Sprintf("%d", crontabType))
	for _, r := range strings.Split(rule, ",") {
		if strings.Index(r, "/") != -1 {
			// 解析每多少时间times
			rList := strings.Split(r, "/")
			if len(rList) != 2 {
				return
			}
			twoR, err := strconv.ParseInt(rList[1], 10, 64)
			if err != nil {
				return
			}
			if twoR < 1 {
				return
			}
			// 每1时间次执行1次
			if rList[0] == "*" {
				maxBetween := c.getMaxBetween(crontabType)
				if err := c.times(Times{
					Interval: twoR,
					Start:    maxBetween.Start,
					End:      maxBetween.End,
				}, &rs, crontabType); err != nil {
					return
				}
				continue
			}
			_, err = strconv.ParseInt(rList[0], 10, 64)
			if err == nil {
				return
			}
			charList := strings.Split(rList[0], "-")
			if len(charList) != 2 {
				return
			}
			oneNum, err := strconv.ParseInt(charList[0], 10, 64)
			if err != nil {
				return
			}
			twoNum, err := strconv.ParseInt(charList[1], 10, 64)
			if err != nil {
				return
			}
			// 时间段内，每1时间次执行1次
			if err := c.times(Times{
				Interval: twoR,
				Start:    oneNum,
				End:      twoNum,
			}, &rs, crontabType); err != nil {
				return
			}
		} else if strings.Index(r, "-") != -1 {
			// 解释时间段between
			rList := strings.Split(r, "-")
			if len(rList) != 2 {
				return
			}
			oneR, err := strconv.ParseInt(rList[0], 10, 64)
			if err != nil {
				return
			}
			twoR, err := strconv.ParseInt(rList[1], 10, 64)
			if err != nil {
				return
			}
			if err := c.between(Between{
				Start: oneR,
				End:   twoR,
			}, &rs, crontabType); err != nil {
				return
			}
		} else {
			if r == "*" {
				// 每1时间次执行1次
				maxBetween := c.getMaxBetween(crontabType)
				if err := c.times(Times{
					Interval: 1,
					Start:    maxBetween.Start,
					End:      maxBetween.End,
				}, &rs, crontabType); err != nil {
					return
				}
				continue
			}
			charNum, err := strconv.ParseInt(r, 10, 64)
			if err != nil {
				return
			}
			if err := c.between(Between{
				Start: charNum,
				End:   charNum,
			}, &rs, crontabType); err != nil {
				return
			}
		}
	}
	rsErr = nil
	return
}

func (c crontab) times(times Times, list *[]int64, crontabType uint8) (rsErr error) {
	rsErr = errors.New("参数错误")
	maxBetween := c.getMaxBetween(crontabType)
	if times.Interval < 1 || times.Start < maxBetween.Start || times.End > maxBetween.End {
		return
	}
	listData := *list
	for i := times.Start; i <= times.End; i += times.Interval {
		if !c.inArrayInt64(i, listData) {
			listData = append(listData, i)
		}
	}
	*list = listData
	rsErr = nil
	return
}

func (c crontab) between(between Between, list *[]int64, crontabType uint8) (rsErr error) {
	rsErr = errors.New("参数错误")
	maxBetween := c.getMaxBetween(crontabType)
	if between.Start < maxBetween.Start || between.End > maxBetween.End {
		return
	}
	listData := *list
	for i := between.Start; i <= between.End; i++ {
		if !c.inArrayInt64(i, listData) {
			listData = append(listData, i)
		}
	}
	*list = listData
	rsErr = nil
	return
}

func (c crontab) getMaxBetween(ct uint8) Between {
	var start int64
	var end int64
	switch ct {
	case crontabType.Minute, crontabType.Hour:
		start = 0
		end = 60
	case crontabType.Day:
		start = 1
		end = 31
	case crontabType.Month:
		start = 1
		end = 12
	case crontabType.Week:
		start = 0
		end = 6
	}
	return Between{
		Start: start,
		End:   end,
	}
}

func (c crontab) inArrayInt64(val int64, arr []int64) bool {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}

func (c crontab) timeSeparate(time time.Time) (minute, hour, day, month, week int64) {
	minute = int64(time.Minute())
	hour = int64(time.Hour())
	day = int64(time.Day())
	month = int64(time.Month())
	week = int64(time.Weekday())
	return
}
