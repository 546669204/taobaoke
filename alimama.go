package taobaoke

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp/runner"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/tuotoo/qrcode"

	httpdo "github.com/546669204/golang-http-do"

	"github.com/tidwall/gjson"
)

type Products struct {
	ID      string
	Title   string
	Bili    string
	Yongjin string
	Jiage   string
	Youhui  string
	Pic     string
}
type Link struct {
	ShortLinkUrl       string
	TaoToken           string
	CouponShortLinkUrl string
	CouponLinkTaoToken string
}
type UserInfoModel struct {
	Name      string
	Pic       string
	LastLogin string
}
type SelfAdzone struct {
	Siteid   string
	Adzoneid []string
}
type Entry struct {
	//use form github.com/546669204/golang-http-do/cookies.go type entry
	Name          string
	Value         string
	Domain        string
	Path          string
	Secure        bool
	HttpOnly      bool
	Persistent    bool
	HostOnly      bool
	Expires       time.Time
	Creation      time.Time
	LastAccess    time.Time
	Updated       time.Time
	CanonicalHost string
}

var ChromeUserDataDIR = ""

var tb_token = "e5b7657bb757esdc"
var pvid = "10_"
var UserInfo UserInfoModel
var conTextMap map[string]int64
var c *chromedp.CDP
var ctxt context.Context
var cancel context.CancelFunc
var browserRun int

func init() {
	httpdo.Autocookieflag = true
	browserRun = 0
}
func Login(QrcodeStr *string, lg *string) bool {
	op := httpdo.Default()
	var timestamp int64 = time.Now().UnixNano() / 1000000
	/*
		https://qrlogin.taobao.com/qrcodelogin/generateQRCode4Login.do?from=alimama&appkey=00000000&_ksTS=1518060289319_30&callback=jsonp31

		(function(){jsonp31({"success":true,"message":"null","url":"//img.alicdn.com/tfscom/TB1m2dpXwKTBuNkSne1wu1JoXXa.png","lgToken":"de99458ca8b8ea36121b060b0366ba45","adToken":"fbfa92d66fc980a2a1a4a20dac55a4fa"});})();
	*/
	/*
	   {"xv":"3.3.7","xt":"C1456801506861066291436551527262916933767","etf":"u","xa":"taobao_login","siteId":"","uid":"","eml":"AA","etid":"","esid":"","type":"pc","nce":true,"plat":"Win32","nacn":"Mozilla","nan":"Netscape","nlg":"zh-CN","sw":1366,"sh":768,"saw":1366,"sah":728,"bsw":1349,"bsh":919,"eloc":"https%3A%2F%2Flogin.taobao.com%2Fmember%2Flogin.jhtml","etz":480,"ett":1527262917061,"ecn":"00d28072d163458ac4e5f76257988ce89e70c44e","eca":"B5RVD3WPjnkCAWq6GBzzSjQY","est":2,"xs":"2BA477700510A7DFDD385CBE84379362010373F6C7D0D0ECA1391732D5273DC933A7BE54B94CA04CCD43AD3E795C914CE03B127657950E82D778D535531B3C4D","ms":"1147","erd":"default,communications,4535fcec252bd8f8d01cec8f5d24e8066556152d2ee086f85f9aaf47e35c0def,c8fb49d17fa46987570d448e09d428f03c0b3d44680d89eb4445c75d07011680,default,communications,285726e33527b95a362ad5e680c4142cb553d4cc51bed714f108b27790d9d0fa","cacheid":"7f0abc57f372979a","xh":"","ips":"192.168.0.105","epl":4,"ep":"28d732090cbbaf0e0f09c3c1ad8f6b8c5ae59712","epls":"C370c307f4aca7858493dfe322254e5cb438be944,N0fcd6e18ff6df74f98a698b7f6b6d838a6c11e69,W31cc12636d0f0ccb723b965c719df6faa8883435","esl":false}
	*/
	op = httpdo.Default()
	op.Url = fmt.Sprintf(`https://login.taobao.com/member/login.jhtml?redirectURL=https%%3A%%2F%%2Fwww.taobao.com%%2F`)
	httpbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return false
	}
	op = httpdo.Default()
	op.Url = fmt.Sprintf(`https://ynuf.alipay.com/uid`)
	httpbyte, err = httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return false
	}
	cacheid := strings.Split(string(httpbyte), "\"")[1]
	op = httpdo.Default()
	op.Url = fmt.Sprintf(`https://ynuf.alipay.com/service/um.json`)
	op.Method = "POST"
	op.Data = `data=ENCODE~~V01~~` + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"xv":"3.3.7","xt":"C1456801506861066291436551527262916933767","etf":"u","xa":"taobao_login","siteId":"","uid":"","eml":"AA","etid":"","esid":"","type":"pc","nce":true,"plat":"Win32","nacn":"Mozilla","nan":"Netscape","nlg":"zh-CN","sw":1366,"sh":768,"saw":1366,"sah":728,"bsw":1349,"bsh":919,"eloc":"https%%3A%%2F%%2Flogin.taobao.com%%2Fmember%%2Flogin.jhtml","etz":480,"ett":%d,"ecn":"00d28072d163458ac4e5f76257988ce89e70c44e","eca":"B5RVD3WPjnkCAWq6GBzzSjQY","est":2,"xs":"2BA477700510A7DFDD385CBE84379362010373F6C7D0D0ECA1391732D5273DC933A7BE54B94CA04CCD43AD3E795C914CE03B127657950E82D778D535531B3C4D","ms":"1147","erd":"default,communications,4535fcec252bd8f8d01cec8f5d24e8066556152d2ee086f85f9aaf47e35c0def,c8fb49d17fa46987570d448e09d428f03c0b3d44680d89eb4445c75d07011680,default,communications,285726e33527b95a362ad5e680c4142cb553d4cc51bed714f108b27790d9d0fa","cacheid":"%s","xh":"","ips":"192.168.0.105","epl":4,"ep":"28d732090cbbaf0e0f09c3c1ad8f6b8c5ae59712","epls":"C370c307f4aca7858493dfe322254e5cb438be944,N0fcd6e18ff6df74f98a698b7f6b6d838a6c11e69,W31cc12636d0f0ccb723b965c719df6faa8883435","esl":false}`, timestamp, cacheid)))
	op.Header = "Content-Type:application/x-www-form-urlencoded; charset=UTF-8"
	httpbyte, err = httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return false
	}
	op = httpdo.Default()
	op.Url = fmt.Sprintf(`https://qrlogin.taobao.com/qrcodelogin/generateQRCode4Login.do?from=alimama&appkey=00000000&_ksTS=%d_30&callback=jsonp31&umid_token=%s`, timestamp, gjson.ParseBytes(httpbyte).Get("tn").String())
	op.Header = fmt.Sprintf("authority:qrlogin.taobao.com\nmethod:GET\npath:/qrcodelogin/generateQRCode4Login.do?from=alimama&appkey=00000000&_ksTS=%d_30&callback=jsonp31\nscheme:https", timestamp)
	httpbyte, err = httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return false
	}
	reg := regexp.MustCompile(`^\(function\(\)\{jsonp31\((\S+)\);\}\)\(\)`)
	matches := reg.FindStringSubmatch(string(httpbyte))
	if len(matches) != 2 {
		log.Println("返回格式不正确")
		return false
	}

	qrurl := gjson.Get(matches[1], "url").String()
	lgToken := gjson.Get(matches[1], "lgToken").String()
	//adToken := gjson.Get(matches[1], "adToken").String()

	op = httpdo.Default()
	op.Url = "https:" + qrurl
	op.Header = fmt.Sprintf("authority:img.alicdn.com\nmethod:GET\npath:%s\nscheme:https", strings.Replace(qrurl, "//img.alicdn.com", "", -1))
	httpbyte, err = httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return false
	}

	M, err := qrcode.Decode(bytes.NewReader(httpbyte))

	*QrcodeStr = M.Content
	*lg = lgToken
	//qrterminal.GenerateHalfBlock(M.Content, qrterminal.L, os.Stdout)

	return true
}
func BrowserLogin() string {
	var err error

	conTextMap = make(map[string]int64)
	// 创建内容
	ctxt, cancel = context.WithCancel(context.Background())
	//defer cancel()
	if len(ChromeUserDataDIR) > 0 {
		c, err = chromedp.New(ctxt, chromedp.WithRunnerOptions(
			runner.Flag("headless", true),
			runner.Flag("disable-gpu", true),
			runner.Flag("no-first-run", true),
			runner.Flag("no-sandbox", true),
			runner.Flag("no-default-browser-check", true),
			// runner.Flag("disable-popup-blocking", false),                                 //关闭弹窗拦截
			// runner.Flag("disable-web-security", true),                                    //安全策略 跨域之类
			runner.Flag("user-data-dir", ChromeUserDataDIR),
		), chromedp.WithLog(BrowserHandler))
	} else {
		c, err = chromedp.New(ctxt, chromedp.WithRunnerOptions(
			runner.Flag("headless", true),
			runner.Flag("disable-gpu", true),
			runner.Flag("no-first-run", true),
			runner.Flag("no-sandbox", true),
			runner.Flag("no-default-browser-check", true),
			// runner.Flag("disable-popup-blocking", false),                                 //关闭弹窗拦截
			// runner.Flag("disable-web-security", true),                                    //安全策略 跨域之类
		), chromedp.WithLog(BrowserHandler))
	}
	// 创建chrome实例

	if err != nil {
		log.Fatal("chromedp.New", err)
	}

	// 运行任务
	var img string
	err = c.Run(ctxt, getQrcode(&img))
	if err != nil {
		log.Fatal("getQrcode", err)
	}
	browserRun = 1
	log.Println("输出测试 ", img)
	return img
}
func BrowserHandler(a string, b ...interface{}) {
	log.Printf(a, b...)
	if len(b) >= 1 && reflect.TypeOf(b[0]).String() == "string" {
		data := gjson.Parse(b[0].(string))
		if "Runtime.executionContextCreated" == data.Get("method").String() {
			conTextMap[data.Get("params").Get("context").Get("auxData").Get("frameId").String()] = data.Get("params").Get("context").Get("id").Int()
		}
	}

}
func getQrcode(img *string) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(`https://www.alimama.com/member/login.htm?forward=http%3A%2F%2Fpub.alimama.com%2Fmyunion.htm%3Fspm%3Da219t.7900221%2F1.a214tr8.2.446dfb5b8vg0Sx`),
		chromedp.WaitVisible(`.mm-logo`, chromedp.ByQuery),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(z context.Context, h cdp.Executor) error {
			for {
				frameTree, err := page.GetFrameTree().Do(z, h)

				if err != nil {
					log.Println(err.Error())
				}
				for _, v := range frameTree.ChildFrames {
					if v.Frame.Name != "taobaoLoginIfr" {
						continue
					}
					cid, ok := conTextMap[string(v.Frame.ID)]
					if !ok {
						continue
					}
					res, _, err := runtime.Evaluate(`
							function test (){
								document.querySelector("#J_LoginBox") && document.querySelector("#J_LoginBox").className && document.querySelector("#J_LoginBox").className.indexOf("module-quick") == -1 && (document.querySelector(".login-switch .quick").click());
								if (!document.querySelector("#J_QRCodeImg img")){return 0}
								if (document.querySelector("#J_QRCodeImg img").naturalWidth==0){return 0}
								if (document.querySelector("#J_QRCodeImg img").naturalHeight==0){return 0}
								if (document.querySelector("#J_QRCodeImg img").width==0){return 0}
								if (document.querySelector("#J_QRCodeImg img").height==0){return 0}
								var img = document.querySelector("#J_QRCodeImg img");
								return '{"width":'+img.width+',"height":'+img.height+',"src":"'+img.src+'"}'
							}
							test()
						`).WithContextID(runtime.ExecutionContextID(cid)).Do(z, h)

					if err != nil {
						log.Println(err.Error())
					}
					if string(res.Value) == "0" {
						continue
					}
					//
					*img = string(res.Value)
					log.Println("二维码输出成功", string(res.Value), v.Frame.Name)
					goto ForEnd
				}
				time.Sleep(3 * time.Second)
			}
		ForEnd:
			return nil
		}),
	}
}
func getcookies() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
			cookies, err := network.GetAllCookies().Do(ctxt, h)
			if err != nil {
				return err
			}

			cookiestr := ""
			var b map[string]map[string]Entry
			b = make(map[string]map[string]Entry)
			for _, c := range cookies {
				//if c.Domain == ".alimama.com" { //筛选作用域
				cookiestr += c.Name + "=" + c.Value + ";"
				var d Entry
				domain := c.Domain
				CanonicalHost := c.Domain
				if domain[:1] == "." {
					domain = domain[1:]
					CanonicalHost = "www" + CanonicalHost
				}

				rootdomain := domain
				domaina := strings.Split(domain, ".")
				if len(domaina) == 3 {
					rootdomain = domaina[1] + "." + domaina[2]
				}
				d.Domain = domain
				d.Path = c.Path
				d.Name = c.Name
				d.Value = c.Value
				d.CanonicalHost = CanonicalHost
				expires, _ := time.Parse("2006-01-02 15:04:05", "9999-01-02 15:04:05")
				d.Expires = expires
				if _, ok := b[rootdomain]; !ok {
					b[rootdomain] = make(map[string]Entry)
				}
				b[rootdomain][domain+":/:"+c.Name] = d
			}

			asd, _ := json.Marshal(b)
			file, err := os.OpenFile("cookies.data", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0)
			if err != nil {
				log.Println(err)
			}
			file.Write(asd)
			file.Close()
			httpdo.LoadCookies()
			return nil
		}),
	}
}
func BrowserCheckLogin() (status bool, msg string) {
	var site string
	var err error
	status = false
	msg = ""
	if browserRun == 0 {
		msg = "还未启动浏览器"
		return
	}
	err = c.Run(ctxt, chromedp.Location(&site))
	if err != nil {
		log.Fatal(err)
	}
	// 循环判断网址是否是登陆成功后的网址
	site = strings.Replace(site, "https://", "http://", -1)
	if site[:34] == "http://pub.alimama.com/myunion.htm" {
		err = c.Run(ctxt, getcookies())
		if err != nil {
			log.Fatal(err)
		}
		browserRun = 2
		// 关闭浏览器
		err = c.Shutdown(ctxt)
		if err != nil {
			log.Fatal(err)
		}

		// 等待浏览器完全关闭
		err = c.Wait()
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
	}

	if browserRun == 1 {
		msg = "请扫描二维码"
	}
	if browserRun == 2 {
		status = true
		msg = "登录成功"
	}
	return

}
func CheckLogin(lgToken string) (status bool, msg string) {
	status = false
	msg = ""
	op := httpdo.Default()
	var timestamp int64 = time.Now().UnixNano() / 1000000
	/*
		https://qrlogin.taobao.com/qrcodelogin/qrcodeLoginCheck.do?lgToken=de99458ca8b8ea36121b060b0366ba45&defaulturl=http%3A%2F%2Flogin.taobao.com%2Fmember%2Ftaobaoke%2Flogin.htm%3Fis_login%3D1&_ksTS=1518060293480_85&callback=jsonp86

		(function(){jsonp86({"code":"10000","message":"login start state","success":true});})();
	*/
	op.Url = fmt.Sprintf(`https://qrlogin.taobao.com/qrcodelogin/qrcodeLoginCheck.do?lgToken=%s&defaulturl=http%%3A%%2F%%2Flogin.taobao.com%%2Fmember%%2Ftaobaoke%%2Flogin.htm%%3Fis_login%%3D1&_ksTS=%d_30&callback=jsonp31`, lgToken, timestamp)
	op.Header = fmt.Sprintf("authority:qrlogin.taobao.com\nmethod:GET\npath:/qrcodelogin/qrcodeLoginCheck.do?lgToken=%s&defaulturl=http%%3A%%2F%%2Flogin.taobao.com%%2Fmember%%2Ftaobaoke%%2Flogin.htm%%3Fis_login%%3D1&_ksTS=%d_30&callback=jsonp31\nscheme:https", lgToken, timestamp)
	httpbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return
	}
	reg := regexp.MustCompile(`^\(function\(\)\{jsonp31\((.+)\);\}\)\(\)`)
	matches := reg.FindStringSubmatch(string(httpbyte))
	if len(matches) != 2 {
		log.Println("返回格式不正确", string(httpbyte))
		return
	}
	//log.Println(matches[1])
	code := gjson.Get(matches[1], "code").Int()
	switch code {
	case 10000:
		msg = "请扫描二维码"
		break
	case 10001:
		msg = "请在手机上确认"
		break
	case 10004:
		msg = "二维码已失效"
		break
	case 10006:
		op = httpdo.Default()
		op.Url = gjson.Get(matches[1], "url").String()
		op.Header = fmt.Sprintf("authority:login.taobao.com\nmethod:GET\npath:%s\nscheme:https", strings.Replace(op.Url, "https://login.taobao.com", "", -1))
		_, err := httpdo.HttpDo(op)
		if err != nil {
			log.Println(err)
			return
		}
		status = true
		msg = "登录成功"
		return
		break
	}
	return
}

func KeepLogin() {
	op := httpdo.Default()
	op.Url = `https://pub.alimama.com/common/getUnionPubContextInfo.json`
	_, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println(string(htmlbyte))
	/*op.Url = fmt.Sprintf(`https://pub.alimama.com/report/getTbkPaymentDetails.json?startTime=%s&endTime=%s&payStatus=&queryType=1&toPage=1&perPageSize=20`, time.Now().AddDate(0, 0, -90).Format("2006-01-02"), time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	_, err = httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return
	}*/
	//log.Println(string(htmlbyte))
}
func GetUnionPubContextInfo() {
	op := httpdo.Default()
	op.Url = "http://pub.alimama.com/common/getUnionPubContextInfo.json"
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return
	}
	ok := gjson.Get(string(htmlbyte), "ok")
	if !ok.Exists() || !ok.Bool() {
		log.Println(gjson.Get(string(htmlbyte), "info").Get("message").String())
		return
	}
	data := gjson.Get(string(htmlbyte), "data")
	log.Println(data.String())
	UserInfo.Name = data.Get("mmNick").String()
	UserInfo.Pic = data.Get("avatar").String()
	UserInfo.LastLogin = time.Now().Format("2006-01-02 15:04:05")
	pvid = "10_" + gjson.Get(string(htmlbyte), "ip").String() + "_7878_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	return
}
func Search(keyword string) (product Products) {
	op := httpdo.Default()
	op.Url = "http://pub.alimama.com/items/search.json?q=" + url.QueryEscape(keyword) + "&auctionTag=&perPageSize=50&shopTag=" + token()
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return
	}
	ok := gjson.Get(string(htmlbyte), "ok")
	if !ok.Exists() || !ok.Bool() {
		log.Println(gjson.Get(string(htmlbyte), "info").Get("message").String())
		return
	}
	data := gjson.Get(string(htmlbyte), "data")
	array := data.Get("pageList").Array()
	if len(array) <= 0 {
		return
	}
	product.ID = array[0].Get("auctionId").String()
	product.Title = array[0].Get("title").String()
	product.Bili = array[0].Get("tkRate").String()
	product.Yongjin = array[0].Get("tkCommFee").String()
	product.Jiage = array[0].Get("zkPrice").String()
	product.Youhui = array[0].Get("couponAmount").String()
	product.Pic = array[0].Get("pictUrl").String()

	/*log.Println("产品id", array[0].Get("auctionId").String())
	log.Println("标题", array[0].Get("title").String())
	log.Println("佣金比例", array[0].Get("tkRate").String())
	log.Println("佣金", array[0].Get("tkCommFee").String())
	log.Println("购买价格", array[0].Get("zkPrice").String())
	log.Println("优惠卷", array[0].Get("couponAmount").String())
	log.Println("图片", array[0].Get("pictUrl").String())*/
	return

}

func NewSelfAdzone2(itemId string) []SelfAdzone {
	op := httpdo.Default()
	op.Url = "http://pub.alimama.com/common/adzone/newSelfAdzone2.json?tag=29&itemId=" + itemId + "&blockId=" + token()
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return nil
	}
	ok := gjson.Get(string(htmlbyte), "ok")
	if !ok.Exists() || !ok.Bool() {
		log.Println(gjson.Get(string(htmlbyte), "info").Get("message").String())
		return nil
	}

	data := gjson.Get(string(htmlbyte), "data")
	array := data.Get("otherAdzones").Array()
	var ret []SelfAdzone
	for _, value := range array {
		var adzone []string
		for _, value2 := range value.Get("sub").Array() {
			ad := value2.Get("id").String()
			adzone = append(adzone, ad)
		}
		var tmp = SelfAdzone{Siteid: value.Get("id").String(), Adzoneid: adzone}
		ret = append(ret, tmp)
	}

	return ret
}

func SelfAdzoneCreate(siteid, adzoneid string) {
	op := httpdo.Default()
	op.Method = "POST"
	op.Url = "http://pub.alimama.com/common/adzone/selfAdzoneCreate.json"
	op.Data = "tag=29&gcid=8&siteid=" + siteid + "&selectact=sel&adzoneid=" + adzoneid + token()
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return
	}
	ok := gjson.Get(string(htmlbyte), "ok")
	if !ok.Exists() || !ok.Bool() {
		log.Println(gjson.Get(string(htmlbyte), "info").Get("message").String())
		return
	}
	return
}
func GetAuctionCode(itemId, siteid, adzoneid string) (l Link) {
	op := httpdo.Default()
	op.Url = "http://pub.alimama.com/common/code/getAuctionCode.json?auctionid=" + itemId + "&adzoneid=" + adzoneid + "&siteid=" + siteid + "&scenes=1" + token()
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return
	}
	ok := gjson.Get(string(htmlbyte), "ok")
	if !ok.Exists() || !ok.Bool() {
		log.Println(gjson.Get(string(htmlbyte), "info").Get("message").String())
		return
	}
	data := gjson.Get(string(htmlbyte), "data")
	l.ShortLinkUrl = data.Get("shortLinkUrl").String()
	l.TaoToken = data.Get("taoToken").String()
	l.CouponShortLinkUrl = data.Get("couponShortLinkUrl").String()
	l.CouponLinkTaoToken = data.Get("couponLinkTaoToken").String()
	return
	//log.Println(shortLinkUrl, taoToken, couponShortLinkUrl, couponLinkTaoToken)
}
func token() string {
	return "&t=" + strconv.FormatInt(time.Now().UnixNano(), 10) + "&_tb_token_=" + tb_token + "&pvid=" + pvid
}

func SaveLogin() {
	httpdo.SaveCookies()
	return
}
func LoadLogin() bool {
	httpdo.LoadCookies()
	CookiesTotb_token()
	if tb_token == "e5b7657bb757esdc" {
		return false
	}
	return true
}
func IsLogin() bool {
	op := httpdo.Default()
	op.Url = `https://pub.alimama.com/common/getUnionPubContextInfo.json`
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return false
	}
	data := gjson.ParseBytes(htmlbyte)
	if data.Get("data").Get("noLogin").Bool() == true {
		return false
	}
	return true
}
func CookiesTotb_token() {
	u, _ := url.Parse("http://*.alimama.com/")
	for _, value := range httpdo.Autocookie.Cookies(u) {
		if value.Name == "_tb_token_" {
			tb_token = value.Value
		}
	}
}

func Download(startTime, endTime string) []byte {
	op := httpdo.Default()
	op.Url = fmt.Sprintf("http://pub.alimama.com/report/getTbkPaymentDetails.json?queryType=1&payStatus=&DownloadID=DOWNLOAD_REPORT_INCOME_NEW&startTime=%s&endTime=%s", startTime, endTime)
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return nil
	}
	return htmlbyte
}

func MediaRpt(startTime, endTime string) []byte {
	/*
		http://pub.alimama.com/report/mediaRpt.json?gcId=&siteType=&siteId=&startTime=2018-03-09&endTime=2018-03-15&t=1521185338929&pvid=&_tb_token_=7b703b3b7ee13&_input_charset=utf-8

		{"data":{"datas":[{"siteId":null,"thedate":"2018-03-09","siteName":null,"alipayNum":0,"rec":0.00,"mixClick":0,"alipayRec":0.00,"mixPv":null,"mixCtr":null,"mixEcpm":null,"alipayAmt":null,"mixRphc":"--"},{"siteId":null,"thedate":"2018-03-10","siteName":null,"alipayNum":0,"rec":0.00,"mixClick":0,"alipayRec":0.00,"mixPv":null,"mixCtr":null,"mixEcpm":null,"alipayAmt":null,"mixRphc":"--"},{"siteId":null,"thedate":"2018-03-11","siteName":null,"alipayNum":0,"rec":0.09,"mixClick":0,"alipayRec":0.00,"mixPv":null,"mixCtr":null,"mixEcpm":null,"alipayAmt":null,"mixRphc":"--"},{"siteId":null,"thedate":"2018-03-12","siteName":null,"alipayNum":0,"rec":0.00,"mixClick":0,"alipayRec":0.00,"mixPv":null,"mixCtr":null,"mixEcpm":null,"alipayAmt":null,"mixRphc":"--"},{"siteId":null,"thedate":"2018-03-13","siteName":null,"alipayNum":0,"rec":0.00,"mixClick":0,"alipayRec":0.00,"mixPv":null,"mixCtr":null,"mixEcpm":null,"alipayAmt":null,"mixRphc":"--"},{"siteId":null,"thedate":"2018-03-14","siteName":null,"alipayNum":1,"rec":0.00,"mixClick":2,"alipayRec":0.72,"mixPv":null,"mixCtr":null,"mixEcpm":null,"alipayAmt":null,"mixRphc":"0.00"},{"siteId":null,"thedate":"2018-03-15","siteName":null,"alipayNum":0,"rec":0.00,"mixClick":0,"alipayRec":0.00,"mixPv":null,"mixCtr":null,"mixEcpm":null,"alipayAmt":null,"mixRphc":"--"}],"webSites":[],"query":{"siteId":null,"siteIds":null,"startTime":"2018-03-09","endTime":"2018-03-15","pubId":112672261,"length":10,"offset":0,"pageNo":0,"pageSize":0,"toPage":0,"perPageSize":0,"startRow":0},"apps":[],"countMap":{"totalRec":0.09,"totalAlipayNum":1,"totalAlipayRec":0.72,"totalMixClick":2},"softs":[],"guides":[{"name":"亲朋好友推","id":42520773},{"name":"爱分享(手机客户端专享)_112672261","id":42488830}],"allMedias":[{"name":"亲朋好友推","id":42520773},{"name":"爱分享(手机客户端专享)_112672261","id":42488830}]},"info":{"message":null,"ok":true},"ok":true,"invalidKey":null}
	*/
	op := httpdo.Default()
	var timestamp int64 = time.Now().UnixNano() / 1000000
	op.Url = fmt.Sprintf("http://pub.alimama.com/report/mediaRpt.json?gcId=&siteType=&siteId=&startTime=%s&endTime=%s&t=%d&pvid=&_tb_token_=%s&_input_charset=utf-8", startTime, endTime, timestamp, tb_token)
	op.Header = "Content-Type:application/x-www-form-urlencoded; charset=UTF-8\nX-Requested-With:XMLHttpRequest"
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return nil
	}
	return htmlbyte
}

func GetTbkPaymentDetails(startTime, endTime string) []byte {
	/*
		http://pub.alimama.com/report/getTbkPaymentDetails.json?startTime=2017-12-23&endTime=2018-03-22&payStatus=&queryType=1&toPage=1&perPageSize=20&total=&t=1521784873057&pvid=&_tb_token_=e60bd0555e318&_input_charset=utf-8

		{"data":{"paginator":{"length":8,"offset":0,"page":1,"beginIndex":1,"endIndex":8,"items":8,"pages":1,"itemsPerPage":20,"firstPage":1,"lastPage":1,"previousPage":1,"nextPage":1,"slider":[1]},"paymentList":[{"createTime":"2018-03-21 09:31:03","bizType":200,"shareRate":"100.00 %","earningTime":null,"tkBizTag":3,"tk3rdPubShareFee":0.00,"tk3rdTypeStr":null,"auctionUrl":"http://item.taobao.com:80/item.htm?id=37256767583","auctionId":37256767583,"tkShareRate":0.00,"payStatus":12,"auctionTitle":"【狂欢价】蚊帐蒙古包1.8m床1.5双人家用加密加厚三开门1.2米床单人学生宿舍","exShopTitle":"金喜路家纺旗舰店","realPayFeeString":"0","auctionNum":1,"payPrice":363.00,"taobaoTradeParentId":"139715127486773509","exNickName":"金喜路家纺旗舰店","tkShareRateToString":"0.00","tkPubShareFeeString":"0","feeString":"0.79","terminalType":"无线","discountAndSubsidyToString":"2.00 %","finalDiscountToString":"1.50","exMemberId":807576160,"realPayFee":null,"totalAlipayFeeString":"39.00"},{"createTime":"2018-03-20 11:08:56","bizType":200,"shareRate":"100.00 %","earningTime":null,"tkBizTag":2,"tk3rdPubShareFee":0.00,"tk3rdTypeStr":null,"auctionUrl":"http://item.taobao.com:80/item.htm?id=545892162708","auctionId":545892162708,"tkShareRate":0.00,"payStatus":12,"auctionTitle":"【狂欢价】太力真空压缩袋11件套收纳袋衣物棉被子大号整理打包袋抽气真空袋","exShopTitle":"太力家居旗舰店","realPayFeeString":"0","auctionNum":2,"payPrice":139.90,"taobaoTradeParentId":"139127293725773509","exNickName":"太力家居旗舰店","tkShareRateToString":"0.00","tkPubShareFeeString":"0","feeString":"4.84","terminalType":"无线","discountAndSubsidyToString":"5.50 %","finalDiscountToString":"5.00","exMemberId":138960927,"realPayFee":null,"totalAlipayFeeString":"88.00"},{"createTime":"2018-03-20 10:29:08","bizType":200,"shareRate":"100.00 %","earningTime":null,"tkBizTag":2,"tk3rdPubShareFee":0.00,"tk3rdTypeStr":null,"auctionUrl":"http://item.taobao.com:80/item.htm?id=527852105881","auctionId":527852105881,"tkShareRate":0.00,"payStatus":12,"auctionTitle":"加固搬家袋牛津布编织袋加厚行李收纳袋子大容量防水旅行包托运袋","exShopTitle":"爱贝办公旗舰店","realPayFeeString":"0","auctionNum":2,"payPrice":19.80,"taobaoTradeParentId":"139111109680773509","exNickName":"爱贝办公旗舰店","tkShareRateToString":"0.00","tkPubShareFeeString":"0","feeString":"0.40","terminalType":"无线","discountAndSubsidyToString":"2.01 %","finalDiscountToString":"2.00","exMemberId":903487543,"realPayFee":null,"totalAlipayFeeString":"19.80"},{"createTime":"2018-03-18 22:56:51","bizType":200,"shareRate":"100.00 %","earningTime":"2018-03-20 22:08:31","tkBizTag":2,"tk3rdPubShareFee":0.00,"tk3rdTypeStr":null,"auctionUrl":"http://item.taobao.com:80/item.htm?id=522587894416","auctionId":522587894416,"tkShareRate":0.00,"payStatus":3,"auctionTitle":"CHANDO/自然堂滋润滋养护唇膏 保湿滋润修护干裂润唇膏 补水保湿","exShopTitle":"自然堂旗舰店","realPayFeeString":"34.00","auctionNum":1,"payPrice":35.00,"taobaoTradeParentId":"138710289409773509","exNickName":"自然堂旗舰店","tkShareRateToString":"0.00","tkPubShareFeeString":"1.36","feeString":"1.36","terminalType":"无线","discountAndSubsidyToString":"4.00 %","finalDiscountToString":"3.00","exMemberId":1652554937,"realPayFee":34.00,"totalAlipayFeeString":"34.00"},{"createTime":"2018-03-14 09:50:32","bizType":200,"shareRate":"100.00 %","earningTime":null,"tkBizTag":2,"tk3rdPubShareFee":0.00,"tk3rdTypeStr":null,"auctionUrl":"http://item.taobao.com:80/item.htm?id=39657356892","auctionId":39657356892,"tkShareRate":0.00,"payStatus":12,"auctionTitle":"上福 999纯银耳钉 女养耳棒耳针韩国简约学生隐形耳棍养耳洞耳钉","exShopTitle":"上福旗舰店","realPayFeeString":"0","auctionNum":1,"payPrice":40.00,"taobaoTradeParentId":"137273301578773509","exNickName":"上福旗舰店","tkShareRateToString":"0.00","tkPubShareFeeString":"0","feeString":"0.72","terminalType":"无线","discountAndSubsidyToString":"3.61 %","finalDiscountToString":"3.11","exMemberId":2048244473,"realPayFee":null,"totalAlipayFeeString":"20.00"},{"createTime":"2018-03-08 19:25:44","bizType":200,"shareRate":"100.00 %","earningTime":null,"tkBizTag":2,"tk3rdPubShareFee":0.00,"tk3rdTypeStr":null,"auctionUrl":"http://item.taobao.com:80/item.htm?id=38315823364","auctionId":38315823364,"tkShareRate":0.00,"payStatus":13,"auctionTitle":"品奥可充电电子称体重秤家用成人精准人体秤减肥称测体重计称重器","exShopTitle":"品奥旗舰店","realPayFeeString":"0","auctionNum":1,"payPrice":143.60,"taobaoTradeParentId":"135561997865773509","exNickName":"品奥旗舰店","tkShareRateToString":"0.00","tkPubShareFeeString":"0","feeString":"0","terminalType":"无线","discountAndSubsidyToString":"1.80 %","finalDiscountToString":"1.50","exMemberId":664889198,"realPayFee":null,"totalAlipayFeeString":"0"},{"createTime":"2018-03-02 13:45:19","bizType":200,"shareRate":"100.00 %","earningTime":null,"tkBizTag":1,"tk3rdPubShareFee":0.00,"tk3rdTypeStr":null,"auctionUrl":"http://item.taobao.com:80/item.htm?id=563003718928","auctionId":563003718928,"tkShareRate":0.00,"payStatus":13,"auctionTitle":"一分钱自动发货运营0.01元秒发秒评壁纸0.1商品宝贝1分钱套餐图片","exShopTitle":"雨中阳光饰品店","realPayFeeString":"0","auctionNum":1,"payPrice":0.10,"taobaoTradeParentId":"120905540384628052","exNickName":"zgc13149955","tkShareRateToString":"0.00","tkPubShareFeeString":"0","feeString":"0","terminalType":"无线","discountAndSubsidyToString":"1.50 %","finalDiscountToString":"1.50","exMemberId":759654030,"realPayFee":null,"totalAlipayFeeString":"0"},{"createTime":"2018-03-01 14:10:45","bizType":200,"shareRate":"100.00 %","earningTime":"2018-03-11 12:16:14","tkBizTag":2,"tk3rdPubShareFee":0.00,"tk3rdTypeStr":null,"auctionUrl":"http://item.taobao.com:80/item.htm?id=17689518523","auctionId":17689518523,"tkShareRate":0.00,"payStatus":3,"auctionTitle":"洗照片 1寸2寸证件照 照片冲印冲洗网上洗相片手机大头贴 幼儿园","exShopTitle":"益好旗舰店","realPayFeeString":"4.80","auctionNum":4,"payPrice":2.00,"taobaoTradeParentId":"120604933632628052","exNickName":"益好旗舰店","tkShareRateToString":"0.00","tkPubShareFeeString":"0.09","feeString":"0.09","terminalType":"无线","discountAndSubsidyToString":"2.00 %","finalDiscountToString":"1.50","exMemberId":854149508,"realPayFee":4.80,"totalAlipayFeeString":"4.80"}],"hascount":true},"info":{"message":null,"ok":true},"ok":true,"invalidKey":null}
	*/
	op := httpdo.Default()
	var timestamp int64 = time.Now().UnixNano() / 1000000
	op.Url = fmt.Sprintf("http://pub.alimama.com/report/getTbkPaymentDetails.json?startTime=%s&endTime=%s&payStatus=&queryType=1&toPage=1&perPageSize=20&total=&t=%d&pvid=&_tb_token_=%s&_input_charset=utf-8", startTime, endTime, timestamp, tb_token)
	//op.Header = "Content-Type:application/x-www-form-urlencoded; charset=UTF-8\nX-Requested-With:XMLHttpRequest"
	htmlbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return nil
	}
	return htmlbyte
}
