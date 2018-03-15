package taobaoke

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

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

var tb_token = "e5b7657bb757esdc"
var pvid = "10_"
var UserInfo UserInfoModel

func init() {
	httpdo.Autocookieflag = true
}
func Login(QrcodeStr *string, lg *string) bool {
	op := httpdo.Default()
	var timestamp int64 = time.Now().UnixNano() / 1000000
	/*
		https://qrlogin.taobao.com/qrcodelogin/generateQRCode4Login.do?from=alimama&appkey=00000000&_ksTS=1518060289319_30&callback=jsonp31

		(function(){jsonp31({"success":true,"message":"null","url":"//img.alicdn.com/tfscom/TB1m2dpXwKTBuNkSne1wu1JoXXa.png","lgToken":"de99458ca8b8ea36121b060b0366ba45","adToken":"fbfa92d66fc980a2a1a4a20dac55a4fa"});})();
	*/
	op.Url = fmt.Sprintf(`https://qrlogin.taobao.com/qrcodelogin/generateQRCode4Login.do?from=alimama&appkey=00000000&_ksTS=%d_30&callback=jsonp31`, timestamp)
	op.Header = fmt.Sprintf("authority:qrlogin.taobao.com\nmethod:GET\npath:/qrcodelogin/generateQRCode4Login.do?from=alimama&appkey=00000000&_ksTS=%d_30&callback=jsonp31\nscheme:https", timestamp)
	httpbyte, err := httpdo.HttpDo(op)
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
func CheckLogin(lgToken string) (status bool, msg string) {
	status = false
	msg = ""
	op := httpdo.Default()
	var timestamp int64 = time.Now().UnixNano() / 1000000
	/*
		https://qrlogin.taobao.com/qrcodelogin/qrcodeLoginCheck.do?lgToken=de99458ca8b8ea36121b060b0366ba45&defaulturl=http%3A%2F%2Flogin.taobao.com%2Fmember%2Ftaobaoke%2Flogin.htm%3Fis_login%3D1&_ksTS=1518060293480_85&callback=jsonp86

		(function(){jsonp86({"code":"10000","message":"login start state","success":true});})();
	*/
	op.Url = fmt.Sprintf(`https://qrlogin.taobao.com/qrcodelogin/qrcodeLoginCheck.do?lgToken=%s&defaulturl=http%3A%2F%2Flogin.taobao.com%2Fmember%2Ftaobaoke%2Flogin.htm%3Fis_login%3D1&_ksTS=%d_30&callback=jsonp31`, lgToken, timestamp)
	op.Header = fmt.Sprintf("authority:qrlogin.taobao.com\nmethod:GET\npath:/qrcodelogin/qrcodeLoginCheck.do?lgToken=%s&defaulturl=http%3A%2F%2Flogin.taobao.com%2Fmember%2Ftaobaoke%2Flogin.htm%3Fis_login%3D1&_ksTS=%d_30&callback=jsonp31\nscheme:https", lgToken, timestamp)
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
	op.Url = fmt.Sprintf(`https://pub.alimama.com/report/getTbkPaymentDetails.json?startTime=%s&endTime=%s&payStatus=&queryType=1&toPage=1&perPageSize=20`, time.Now().AddDate(0, 0, -90).Format("2006-01-02"), time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	_, err = httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return
	}
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
	//log.Println(data.String())
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
