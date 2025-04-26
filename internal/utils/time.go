package utils

import "time"

const TimeFormatStr = "2006-01-02 15:04:05"

var TimeLocation *time.Location

func init() {
	var err error
	TimeLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
}

func GetCurrentTime() time.Time {
	return time.Now().In(TimeLocation)
}

func GetCurrentTimeString() string {
	return time.Now().In(TimeLocation).Format(TimeFormatStr)
}

// 当前毫秒级别时间戳
func GetCurrentTimeStampMilli() int64 {
	return time.Now().In(TimeLocation).UnixMilli()
}

// 当前秒级别时间戳
func GetCurrentTimeStampSec() int64 {
	return time.Now().In(TimeLocation).Unix()
}

// 间戳转时间系列
// 第一个参数表示秒级别的时间戳，第二个参数表示纳秒级别的时间戳，如果是毫秒级别的时间戳需要先乘以1e6然后放入第二个参数即可
func timestamp2Time(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).In(TimeLocation)
}

// 秒级别时间戳转时间
func TimestampSec2Time(sec int64) time.Time {
	return timestamp2Time(sec, 0).In(TimeLocation)
}

// 毫秒级别时间戳转时间
func TimestampMilli2Time(stamp int64) time.Time {
	return timestamp2Time(0, stamp*1e6).In(TimeLocation)
}

// 纳秒级别时间戳转时间
func TimestampNano2Time(stamp int64) time.Time {
	return timestamp2Time(0, stamp).In(TimeLocation)
}

// 返回参数1: t1时间减去t2时间的时间差
// 返回参数2: t1时间是否在t2时间的后面
func GetTimeSubSecs(t1, t2 time.Time) time.Duration {
	// Notice 后面不要加Seconds()方法，有时候即使t1在t2后面返回的也有可能是一个负数!
	return t1.Sub(t2)
}
