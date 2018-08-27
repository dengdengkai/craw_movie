// transer
package main

import (
	"fmt"
	"regexp"
	"strconv"
	//	"strings"
)

func Transer(newstr string) string {

	return strconv.QuoteToASCII(newstr)

}
func main() {
	fmt.Println("将要转换的编码中文字是： 今天有烤鸭吃，哈哈！")

	str := "今天有烤鸭吃，哈哈！ ! abc"
	str1 := Transer(str)
	fmt.Printf("中文变为unicode码， 略略略：“ 今天有烤鸭吃，哈哈！”变： %s", str1)
	fmt.Println()
	str2 := "小是正则表达式要匹配的字符串包含中英文vsdfvsva猪cas飞天demicgwegr//lwe;rgc了"
	fmt.Println(str2)
	fmt.Println("请匹配：字符串是否以“小”开头中间按顺序有“d”，有“猪”，以“了“结尾”")

	fmt.Println(strconv.QuoteToASCII("小")) //转换为unicode
	fmt.Println(strconv.QuoteToASCII("猪"))
	fmt.Println(strconv.QuoteToASCII("了"))
	//regex := `^[\u5c0f].*?d.*[\\u732a].*[\u4e86]$`
	regex := `[\\u5c0f]` + `.*[\\u732a].*d.*[\\u4e86]`
	//regex := strconv.QuoteToASCII("小")
	fmt.Printf("匹配规则如下:\n%s", regex)
	s1 := regexp.MustCompile(regex)
	fmt.Println()

	fmt.Println(s1.MatchString(str2))

}
