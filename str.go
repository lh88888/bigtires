package bigtires

import (
	"encoding/json"
	"strings"
)

/**
文本取左边（被查找的文本 string， 欲寻找的文本 string） 返回取到的文本 string
eg：StrGetLeft("123456", "4") 取4的左边，返回"123"，没取到则返回空文本。
*/
func StrGetLeft(orig string, findStr string) string {
	n := strings.Index(orig, findStr)
	if n == -1 {
		return ""
	}
	return string([]byte(orig)[:n])
}

/**
文本取右边（被查找的文本 string， 欲寻找的文本 string） 返回取到的文本 string
eg：StrGetRight("123456", "4") 取4的右边，返回"56"，没取到则返回空文本。
*/
func StrGetRight(orig string, findStr string) string {
	n := strings.Index(orig, findStr) + 1
	if n == -1 {
		return ""
	}
	return string([]byte(orig)[n:])
}

/**
文本取中间（被查找的文本 string， 前面文本 string， 后面文本 string） 返回取到的文本 string
eg：StrGetSub("123456", "12", "56") 取12和56的中间文本，返回"34"，没取到则返回空文本。
*/
func StrGetSub(orig string, first string, last string) string {
	firstIndex := strings.Index(orig, first)
	if firstIndex == -1 {
		return ""
	}
	firstIndex = firstIndex + len(first)
	lastIndex := strings.Index(orig[firstIndex:], last) + firstIndex
	for {
		if lastIndex <= firstIndex {
			orig = strings.Replace(orig, last, "", 1)
			firstIndex = strings.Index(orig, first) + len(first)
			lastIndex = strings.Index(orig, last)
		} else {
			break
		}
	}
	return string([]byte(orig)[firstIndex:lastIndex])
}

/**
文本逐字分割（需要分割的文本 string, 返回结果数组指针 *[]string）返回结果数组成员数 int
将指定文本,逐字分割成数组,保存为指定的变量数组中,返回成员个数,可识别换行符及全半角字符和汉字
*/
func StrSplitChinese(orig string, resArray *[]string) int {
	*resArray = (*resArray)[0:0]
	origByte, _ := EnCodeUtf8ToGbk([]byte(orig))
	c := len(origByte)
	n := 0
	z := 0
	for {
		if n < c {
			if n+1 > c {
				z = 1
			} else if origByte[n] > 128 {
				z = 2
			} else if origByte[n] != 13 {
				z = 1
			} else if origByte[n+1] == 10 {
				z = 2
			} else {
				z = 1
			}
			a, _ := EnCodeGbkToUtf8(origByte[n : n+z])
			*resArray = append(*resArray, string(a))
			n = n + z
		} else {
			break
		}
	}
	return len(origByte)
}

/**
将任意类型对象转为string
传参：
	obj：任意类型变量
返回：
	转换后的字符串，转换失败返回空字符串
*/
func StrToString(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
