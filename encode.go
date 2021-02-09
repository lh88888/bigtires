package bigtires

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// 编码GBK到UTF8(GBK 字节集) 返回UTF8 字节集, 错误信息 error
func EnCodeGbkToUtf8(orig []byte) ([]byte, error) {
	I := bytes.NewReader(orig)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 编码UTF8到GBK(UTF8 字节集) 返回GBK 字节集, 错误信息 error
func EnCodeUtf8ToGbk(orig []byte) ([]byte, error) {
	I := bytes.NewReader(orig)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 编码BIG5到UTF8(BIG5 字节集) 返回UTF8 字节集, 错误信息 error
func EnCodeBig5ToUtf8(orig []byte) ([]byte, error) {
	I := bytes.NewReader(orig)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 编码UTF-8到BIG5(UTF-8 字节集) 返回BIG5 字节集, 错误信息 error
func EnCodeUtf8ToBig5(orig []byte) ([]byte, error) {
	I := bytes.NewReader(orig)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
