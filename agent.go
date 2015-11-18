package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var helptext = `
	使用:

        agent command [argument] [argument] ......

    arguments 包括:

        run [command] [uri]	运行代理服务器

            run - 运行反向代理服务器
            Usage:
                run [ip/server]    设置需要代理的服务器，输出log信息，方便调试
`

var total = 1


type Anget struct {
	Server string
	Uri    string
	Param  string
	Method string
	resp   *http.Response
}

func (anget *Anget) result(w http.ResponseWriter, r *http.Request) {
	// 解析参数, 默认是不会解析的
	r.ParseForm()

	anget.Uri = r.URL.Path
	anget.Method = r.Method
	anget.Param = r.Form.Encode()

	u, _ := url.Parse("http://" + anget.Server + anget.Uri)

	fmt.Printf("\n=================================%d===================================\n", total)
	log.Println("地址：", u)
	log.Println("参数：", anget.Param)
	log.Println("方法：", anget.Method)
	total++

	switch anget.Method {
	case "GET":
		anget.resp, _ = http.Get(u.String())
	case "POST":
		anget.resp, _ = http.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(anget.Param))
	default:
		http.Error(w, http.StatusText(500), 500)
	}
	defer anget.resp.Body.Close()
	body, _ := ioutil.ReadAll(anget.resp.Body)

	// 这个写入到w的信息是输出到客户端的
	fmt.Fprintf(w, string(body))
	log.Println("返回：", string(body))
}

func (anget *Anget) Run(port string) {
	conn, err := net.Dial("tcp", anget.Server)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		os.Exit(0)
	}
	fmt.Println("已连接测试服务器～～～")
	conn.Close()

	http.HandleFunc("/"+anget.Uri, anget.result)
	fmt.Println("代理服务器开启，端口为：" + port)
	// 设置监听的端口
	err = http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	var ang Anget

	commands := os.Args

	if len(commands) < 3 {
		println(helptext)
		os.Exit(0)
	}

	switch commands[1] {
	case "help":
		println(helptext)
	case "run":
		switch commands[2] {
		case "--help", "-h":
			println(helptext)
			os.Exit(0)
		default:
			ang.Server = commands[2]
		}
		ang.Run("9900")
	default:
		println(helptext)
	}
}
