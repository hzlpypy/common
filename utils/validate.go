package utils

import (
	"github.com/hzlpypy/common"
	"regexp"
)

// 身份证校验
func ValidateCardNumber(cardNumber string) bool {
	card15ok, _ := regexp.Match(common.CardNumber15, []byte(cardNumber))
	card18ok, _ := regexp.Match(common.CardNumber18, []byte(cardNumber))
	if !card15ok && !card18ok {
		return false
	}
	return true
}

// 手机号码格式校验
func RegexpMobile(m string) bool {
	reg := regexp.MustCompile(common.Mobile)
	return reg.MatchString(m)
}

// 邮箱格式校验
func RegexpEmail(m string) bool {
	reg := regexp.MustCompile(common.Email)
	return reg.MatchString(m)
}

// 邮箱格式校验
func RegexpIPV4(ip string) bool {
	reg := regexp.MustCompile(common.Ipv4)
	return reg.MatchString(ip)
}
