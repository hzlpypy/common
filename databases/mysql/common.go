package mysql

import (
	"regexp"
	"strings"
)

// -- gorm v1/v2 通用方法 --

// GetColumnName:根据 gorm中的column 获取字段名
func GetColumnName(jsonName string) string {
	for _, name := range strings.Split(jsonName, ";") {
		if strings.Index(name, "column") == -1 {
			continue
		}
		return strings.Replace(name, "column:", "", 1)
	}
	return ""
}

// GetLimitOffset:page limit
func GetLimitOffset(page, pageSize int32) (uint32, uint32) {
	limit := 10 //每页记录条数
	if pageSize != 0 {
		limit = int(pageSize)
	}
	offset := 0 // 查询页数
	if page > 0 {
		offset = (int(page) - 1) * limit
	} else {
		offset = 0
	}
	return uint32(limit), uint32(offset)
}

// FilteredSQLInject:正则过滤sql注入的方法
// 参数 : 要匹配的语句
func FilteredSQLInject(toMatchStrs []string) bool {
	//过滤 ‘
	//ORACLE 注解 --  /**/
	//关键字过滤 update ,delete
	// 正则的字符串, 不能用 " " 因为" "里面的内容会转义
	str := `\b(and|exec|insert|select|drop|grant|alter|delete|update|count|chr|mid|master|truncate|char|declare|or)\b|(\*|;|\+|'|%)`
	re, _ := regexp.Compile(str)
	for _, toMatchStr := range toMatchStrs {
		if !re.MatchString(toMatchStr) {
			return false
		}
	}
	return true
}

