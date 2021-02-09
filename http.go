package bigtires

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//Http请求结构体参数
type ReqParms struct {
	Url           string      //请求地址
	Mode          string      //提交方式：GET POST HEAD PUT OPTIONS DELETE TRACE CONNECT，为空默认为GET
	DataStr       string      //提交字符串数据，POST方式本参数有效，Data与DataByte参数二选一传入即可。
	DataByte      []byte      //提交字节集数据，POST方式本参数有效，Data与DataByte参数二选一传入即可。
	Cookies       string      //附加Cookies，把浏览器中开发者工具中Cookies复制传入即可
	Headers       string      //附加协议头，直接将浏览器抓包的协议头复制下来传入即可，无需调整格式，User-Agent也是在此处传入，如果为空默认为Chrome的UA。
	RetHeaders    http.Header //返回协议头，http.Header类型，需导入"net/http"包，返回协议头的参数通过本变量.Get(参数名 string)获取
	RetStatusCode int         //返回状态码
	Redirect      bool        //是否禁止重定向，true为禁止重定向
	ProxyIP       string      //代理IP，格式IP:端口，如：127.0.0.1:8888
	ProxyUser     string      //代理IP账户
	ProxyPwd      string      //代理IP密码
	TimeOut       int         //超时时间，单位：秒，默认30秒，如果提供大于0的数值，则修改操作超时时间
}

//发送HTTP请求(ReqParms) 返回resStr（请求返回字符串结果 string），resByte（请求返回字节集结果 []byte）和err（错误信息 error）
func HttpSend(suReqs *ReqParms) (resStr string, resByte []byte, err error) {
	//设置超时时间
	if suReqs.TimeOut == 0 {
		suReqs.TimeOut = 30
	}
	client := &http.Client{Timeout: time.Duration(suReqs.TimeOut) * time.Second}
	//判断是否重定向
	if suReqs.Redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	//判断是否有代理IP
	if suReqs.ProxyIP != "" {
		proxyAddr := ""
		if suReqs.ProxyUser == "" {
			proxyAddr = "http://" + suReqs.ProxyIP + "/"
		} else {
			proxyAddr = "http://" + suReqs.ProxyUser + ":" + suReqs.ProxyPwd + "@" + suReqs.ProxyIP + "/"
		}
		proxy, err := url.Parse(proxyAddr)
		if err != nil {
			log.Fatal(err)
		}
		netTransport := &http.Transport{
			Proxy:                 http.ProxyURL(proxy),
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Duration(suReqs.TimeOut) * time.Second,
		}
		client.Timeout = time.Duration(suReqs.TimeOut) * time.Second
		client.Transport = netTransport
	}
	if suReqs.Mode == "" {
		suReqs.Mode = "GET"
	}
	var req *http.Request
	if suReqs.Mode == "POST" || suReqs.Mode == "PUT" || suReqs.Mode == "OPTIONS" || suReqs.Mode == "DELETE" {
		if suReqs.DataStr == "" {
			req, err = http.NewRequest(suReqs.Mode, suReqs.Url, bytes.NewReader(suReqs.DataByte))
			req.Header.Set("Content-Length", strconv.Itoa(len(suReqs.DataByte)))
		} else {
			req, err = http.NewRequest(suReqs.Mode, suReqs.Url, strings.NewReader(suReqs.DataStr))
			req.Header.Set("Content-Length", strconv.Itoa(len(suReqs.DataStr)))
		}
	} else {
		req, err = http.NewRequest(suReqs.Mode, suReqs.Url, nil)
	}
	if err != nil {
		log.Println(err)
		return
	}
	//添加headers
	if strings.Index(suReqs.Headers, "User-Agent") == -1 && strings.Index(suReqs.Headers, "user-agent") == -1 {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.3")
	}
	strSplit := strings.Split(suReqs.Headers, "\n")
	for _, val := range strSplit {
		val = strings.Replace(strings.Replace(val, ": ", ":", 1), "\t", "", -1)
		if val != "" {
			req.Header.Set(StrGetLeft(val, ":"), StrGetRight(val, ":"))
		}
	}
	//添加Cookies
	req.Header.Set("Cookie", suReqs.Cookies)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	//合并Cookies
	suReqs.Cookies = HttpMergeCookies(suReqs.Cookies, HttpCookiesToStr(resp.Cookies()))
	resByte, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	//判断是否需要Gzip解压
	if strings.Index(resp.Header.Get("Content-Encoding"), "gzip") != -1 {
		resByte = HttpGzipUn(resByte)
	}
	suReqs.RetHeaders = resp.Header
	suReqs.RetStatusCode = resp.StatusCode
	//判断是否需要转码，Golang默认UTF8编码，如果网站采用GBK则需要转换为UTF8后Golang才能识别
	resStr = string(resByte)
	if strings.Index(resp.Header.Get("Content-Type:"), "charset=gb") != -1 || strings.Index(resStr, "charset=\"gb") != -1 || strings.Index(resStr, "charset=gb") != -1 {
		resByte, _ = EnCodeGbkToUtf8(resByte)
		resStr = string(resByte)
	}
	return
}

//将http的[]Cookie类型转为Cookies字符串
func HttpCookiesToStr(cookies []*http.Cookie) string {
	res := ""
	for _, v := range cookies {
		res = res + v.Name + "=" + v.Value + "; "
	}
	if res != "" {
		res = res[0 : len(res)-2]
	}
	return res
}

//合并文本Cookies，返回合并后的文本Cookies
func HttpMergeCookies(oldCookies string, newCookies string) string {
	//初步格式化
	oldCookies = strings.TrimSpace(oldCookies)
	if oldCookies != "" && oldCookies[len(oldCookies)-1:len(oldCookies)] == ";" {
		oldCookies = oldCookies + " "
	}

	newCookies = strings.TrimSpace(newCookies)
	if newCookies != "" && newCookies[len(newCookies)-1:len(newCookies)] == ";" {
		newCookies = newCookies + " "
	}
	//开始合并Cookies
	oldArray := strings.Split(oldCookies, "; ")
	for _, val := range oldArray {
		if strings.Index(newCookies, StrGetLeft(val, "=")+"=") == -1 {
			newCookies = newCookies + "; " + val
		}
	}
	return strings.ReplaceAll(newCookies, "; ; ", "; ")
}

//Gzip压缩：传入准备压缩的数据，返回压缩后的数据
func HttpGzipPack(data []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	w.Write(data)
	w.Flush()
	return b.Bytes()
}

//Gzip解压，传入准备解压的数据，返回解压后的数据
func HttpGzipUn(data []byte) []byte {
	var b bytes.Buffer
	b.Write(data)
	r, _ := gzip.NewReader(&b)
	defer r.Close()
	unRes, _ := ioutil.ReadAll(r)
	return unRes
}
