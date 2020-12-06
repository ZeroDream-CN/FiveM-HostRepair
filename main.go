package main

import (
	"fmt"
	"github.com/eddieivan01/nic"
	"log"
	"io/ioutil"
	"strings"
	"runtime"
	"os"
)

var HttpClient nic.Session

func main() {
	
	// 设置 Logger
	log.SetFlags(log.Ltime)
	log.Println("FiveM Host 修复工具 by Akkariin")
	
	// 获取服务器 IP
	log.Println("正在获取最新的 Host...")
	response, err := HttpClient.Get("https://api.zerodream.net/v1/fivemclient/host", nil)
	if err != nil {
		// 获取失败
		log.Println("无法获取服务器 IP，请稍后重试，如问题持续存在请联系 QQ 204034。")
		log.Fatalln(fmt.Sprintf("错误内容：%s", err.Error()))
	}
	log.Println(fmt.Sprintf("获得服务器 IP：%s", response.Text))
	
	// 设定 Hosts 路径
	hostPath := ""
	if runtime.GOOS == "windows" {
		hostPath = "C:/Windows/System32/drivers/etc/hosts"
	} else {
		hostPath = "/etc/hosts"
	}
	
	// 解析 Hosts 文件
	hostFile := file_get_contents(hostPath)
	allLines := strings.Split(hostFile, "\n")
	newHosts := ""
	serverIp := response.Text
	lambdaEx := false
	mirrorEx := false
	policyEx := false
	serverEx := false
	
	// 循环读取每一行内容
	for _, line := range allLines {
		if strings.Contains(line, "lambda.fivem.net") {
			log.Println("找到 lambda.fivem.net 的记录，已更新")
			newHosts += fmt.Sprintf("%s %s\n", serverIp, "lambda.fivem.net")
			lambdaEx = true
		} else if strings.Contains(line, "mirrors.fivem.net") {
			log.Println("找到 mirrors.fivem.net 的记录，已更新")
			newHosts += fmt.Sprintf("%s %s\n", serverIp, "mirrors.fivem.net")
			mirrorEx = true
		} else if strings.Contains(line, "policy-is.fivem.net") {
			log.Println("找到 policy-is.fivem.net 的记录，已更新")
			newHosts += fmt.Sprintf("%s %s\n", serverIp, "policy-is.fivem.net")
			policyEx = true
		} else if strings.Contains(line, "servers-frontend.fivem.net") {
			log.Println("找到 servers-frontend.fivem.net 的记录，已更新")
			newHosts += fmt.Sprintf("%s %s\n", serverIp, "servers-frontend.fivem.net")
			serverEx = true
		} else {
			newHosts += line + "\n"
		}
	}
	
	// 判断是否有不存在的记录
	if !lambdaEx {
		log.Println("未找到 lambda.fivem.net 的记录，已添加")
		newHosts += fmt.Sprintf("%s %s\n", serverIp, "lambda.fivem.net")
	}
	if !mirrorEx {
		log.Println("未找到 mirrors.fivem.net 的记录，已添加")
		newHosts += fmt.Sprintf("%s %s\n", serverIp, "mirrors.fivem.net")
	}
	if !policyEx {
		log.Println("未找到 policy-is.fivem.net 的记录，已添加")
		newHosts += fmt.Sprintf("%s %s\n", serverIp, "policy-is.fivem.net")
	}
	if !serverEx {
		log.Println("未找到 servers-frontend.fivem.net 的记录，已添加")
		newHosts += fmt.Sprintf("%s %s\n", serverIp, "servers-frontend.fivem.net")
	}
	
	// 写入到文件中
	if file_put_contents(hostPath, newHosts) {
		log.Println("Hosts 文件写入成功！")
	} else {
		log.Println("Hosts 文件写入失败，可能是因为没有权限，请尝试右键使用管理员身份运行！（Linux 系统下使用 sudo 命令）")
	}
	
	fmt.Printf("\n请按任意键退出...")
    os.Stdin.Read(make([]byte, 1))
}

// 读文件
func file_get_contents(fileName string) (string) {
    f, err := ioutil.ReadFile(fileName)
    if err != nil {
        log.Fatalln("Failed to read file", err)
    }
    return string(f)
}

// 写文件
func file_put_contents(fileName string, data string) (bool) {
	var bytes = []byte(data)
	err := ioutil.WriteFile(fileName, bytes, 0666)
	if err != nil {
		return false
	}
	return true
}