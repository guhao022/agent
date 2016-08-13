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
                run [ip/server:port]    设置需要代理的服务器，输出log信息，方便调试
`

var total = 1


type Agent struct {
	Server string
	Uri    string
	Param  string
	Method string
	Ip     string
	resp   *http.Response
}

func (agent *Agent) result(w http.ResponseWriter, r *http.Request) {
	
	r.Header.Set("Access-Control-Allow-Origin", "*")
	
	// 解析参数, 默认是不会解析的
	r.ParseForm()
	
	b, _ := ioutil.ReadAll(r.Body)
	
	

	agent.Uri = r.URL.Path
	agent.Method = r.Method
	agent.Param = r.Form.Encode()
	agent.Ip = agent.RemoteIp(r)

	u, _ := url.Parse("http://" + agent.Server + agent.Uri)

	fmt.Printf("\n=================================%d===================================\n", total)
	log.Println("地址：", u)
	log.Println("参数：", string(b))
	log.Println("方法：", agent.Method)
	log.Println("访问者：", agent.Ip)
	total++

	switch agent.Method {
	case "GET":
		agent.resp, _ = http.Get(u.String())
	case "POST":
		agent.resp, _ = http.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(agent.Param))
	default:
		http.Error(w, http.StatusText(500), 500)
	}
	defer agent.resp.Body.Close()
	body, _ := ioutil.ReadAll(agent.resp.Body)

	err := ioutil.WriteFile(agent.Ip + ".html", body, os.ModePerm)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	// 这个写入到w的信息是输出到客户端的
	fmt.Fprintf(w, string(body))
	log.Println("返回：", string(body))
}

func (agent *Agent) RemoteIp(r *http.Request) string {
	ip := strings.Split(r.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (agent *Agent) Run(port string) {
	conn, err := net.Dial("tcp", agent.Server)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		os.Exit(0)
	}
	fmt.Println("已连接测试服务器～～～")
	conn.Close()

	http.HandleFunc("/"+agent.Uri, agent.result)
	fmt.Println("代理服务器开启，端口为：" + port)
	// 设置监听的端口
	err = http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	var ang Agent

	commands := os.Args

	if len(commands) < 3 {
		fmt.Println(helptext)
		os.Exit(0)
	}

	switch commands[1] {
	case "help":
		fmt.Println(helptext)
	case "run":
		switch commands[2] {
		case "--help", "-h":
			fmt.Println(helptext)
			os.Exit(0)
		default:
			ang.Server = commands[2]
		}
		ang.Run("9900")
	default:
		fmt.Println(helptext)
	}
}
