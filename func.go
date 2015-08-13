package foolgo

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

var date_replace_pattern = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

func StrToTime(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(date_replace_pattern...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

func Date(format string, t time.Time) string {
	replacer := strings.NewReplacer(date_replace_pattern...)
	format = replacer.Replace(format)
	return t.Format(format)
}

func Time() int64 {
	return time.Now().Unix()
}

func Ip2long(ip_str string) int64 {
	bits := strings.Split(ip_str, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

func Long2ip(ip_int int64) string {
	var bytes [4]byte
	bytes[0] = byte(ip_int & 0xFF)
	bytes[1] = byte((ip_int >> 8) & 0xFF)
	bytes[2] = byte((ip_int >> 16) & 0xFF)
	bytes[3] = byte((ip_int >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

func Md5(str string) string {
	t := md5.New()
	io.WriteString(t, str)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func Sha1(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}
