package taobaoke

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	var err error

	Login()
	GetUnionPubContextInfo()

	go func() {
		//启动线程,定时访问alimam 保持cookies
		log.Println("cookies save")
		for {
			KeepLogin()
			time.Sleep(5 * 60 * time.Second)
		}
	}()

	http.HandleFunc("/", http112233)
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func http112233(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数, 默认是不会解析的

	p := Search(strings.Join(r.Form["k"], ""))
	if p.ID == "" {
		log.Println("找不到")
		return
	}

	log.Println(p)
	a, b := NewSelfAdzone2(p.ID)
	SelfAdzoneCreate(a, b)
	l := GetAuctionCode(p.ID, a, b)
	log.Println(l)
	json1, _ := json.Marshal(p)
	json2, _ := json.Marshal(l)
	fmt.Fprintf(w, string(json1)+string(json2)) //输出到客户端的信息
}
