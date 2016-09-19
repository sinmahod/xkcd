package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//用于包装xkcd返回的json ，因为暂时制需要img图片数据所以先注释掉其他无用的属性
type xkcd struct {
	//Month      string `json:"month"`
	//Num        int    `json:"num"`
	//Link       string `json:"link"`
	//Year       string `json:"year"`
	//News       string `json:"news"`
	//Safe_title string `json:"safe_title"`
	//Transcript string `json:"transcript"`
	//Alt        string `json:"alt"`
	Img string `json:"img"`
	//Title      string `json:"title"`
	//Day        string `json:"day"`
}

//建立两个常量存放URL中的固定字符串
const (
	XKCDUrl  string = "https://xkcd.com/"
	XKCDFile string = "info.0.json"
)

func getXKCD(num string) (x *xkcd, err error) {
	var v xkcd
	resp, err := http.Get(XKCDUrl + num + "/" + XKCDFile)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return &v, nil
}

func form(w http.ResponseWriter, r *http.Request) {
	if data, err := ioutil.ReadFile("form.html"); err == nil {
		html := string(data)
		fmt.Fprintln(w, html)
	}
}

func handlers(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		//读取参数
		if k == `img` {
			if len(v) > 0 {
				x, err := getXKCD(v[0])
				if err != nil {
					fmt.Printf("%s", err)
					return
				}
				js, _ := json.Marshal(x)
				fmt.Fprintf(w, "%s", js)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", form)
	http.HandleFunc("/img", handlers)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
