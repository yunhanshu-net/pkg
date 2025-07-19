package typex

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Time time.Time

func NewTime(t time.Time) Time {
	return Time(t)
}
func NewTimeP(t time.Time) *Time {
	tp := Time(t)
	return &tp
}

const ctLayout = "2006-01-02 15:04:05"

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	// 处理空字符串情况 - 设置为零值时间
	if s == "" || s == "null" {
		*t = Time(time.Time{})
		return nil
	}

	// 支持多种时间格式
	formats := []string{
		ctLayout,                   // "2006-01-02 15:04:05"
		time.RFC3339,               // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano,           // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02T15:04:05Z",     // ISO 8601 UTC格式
		"2006-01-02T15:04:05.000Z", // 带3位毫秒的UTC格式
		"2006-01-02T15:04:05.999Z", // 带3位毫秒的UTC格式（另一种写法）
		"2006-01-02",               // 仅日期格式
		"15:04:05",                 // 仅时间格式
	}

	var nt time.Time
	for _, format := range formats {
		nt, err = time.Parse(format, s)
		if err == nil {
			*t = Time(nt)
			return nil
		}
	}

	// 如果所有格式都解析失败，返回最后一个错误
	return fmt.Errorf("无法解析时间格式: %s, 支持的格式: %v", s, formats)
}

func (t Time) GetUnix() int64 {
	return time.Time(t).Unix()
}

func (t Time) MarshalJSON() ([]byte, error) {
	// 处理零值时间 - 返回空字符串
	if time.Time(t).IsZero() {
		return []byte(`""`), nil
	}
	return []byte(t.String()), nil
}

func (t Time) String() string {
	return fmt.Sprintf("%q", time.Time(t).Format(ctLayout))
}

func (date *Time) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Time(nullTime.Time)
	return
}

func (date Time) Value() (driver.Value, error) {
	ti := time.Time(date)
	y, m, d := ti.Date()
	h := ti.Hour()
	minute := ti.Minute()
	s := ti.Second()
	return time.Date(y, m, d, h, minute, s, 0, time.Time(date).Location()), nil
}
