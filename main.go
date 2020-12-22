package main

import (
	"fmt"
	"github.com/eddieivan01/nic"
	"github.com/json-iterator/go"
	"log"
	"io/ioutil"
	"strings"
	"runtime"
	"os"
)

var HttpClient nic.Session
var hostResult RequestResult
var hostPath   string

func main() {
	
	// 设置 Logger
	log.SetFlags(log.Ltime)
	log.Println("FiveM Host 修复工具 by Akkariin")
	
	// 设定 Hosts 路径
	if runtime.GOOS == "windows" {
		hostPath = "C:/Windows/System32/drivers/etc/hosts"
	} else {
		hostPath = "/etc/hosts"
	}
	
	log.Println("正在检测网络状态...")
	_, err := HttpClient.Get("http://ocsp.int-x3.letsencrypt.org/", nic.H {
		Timeout: 5,
	})

	if err != nil {
		if strings.Contains(err.Error(), "Timeout") {
			log.Println("检测到 Let's Encrypt 服务器可能被 DNS 污染，修复中...")
			if FixLetsEncrypt() {
				log.Println("Let's Encrypt 服务器 Hosts 修复成功")
			} else {
				log.Println("Let's Encrypt 服务器 Hosts 修复失败！")
				if runtime.GOOS == "windows" {
					log.Fatalln("请检查是否用 [右键 -> 以管理员身份运行] 本程序，如果使用管理员权限运行后还是出错，请手动打开 C:/Windows/System32/drivers/etc/ 并右键查看 hosts 文件是否被设置为了只读，如果是请将只读取消然后重试。如果问题持续存在，请联系 QQ 204034。")
				} else {
					log.Fatalln("请检查 /etc/hosts 文件是否是只读状态，以及您当前的用户组是否有读写权限，请尝试使用 root 身份执行本操作！")
				}
			}
		} else {
			log.Println("检测时发生了未知错误")
			log.Fatalln(err.Error())
		}
	}
	
	// 获取服务器 IP
	log.Println("正在获取最新的 Host...")
	response, err := HttpClient.Get("https://api.zerodream.net/v2/fivemclient/host", nil)
	if err != nil {
		// 获取失败
		log.Println("无法获取服务器 IP，请稍后重试，如问题持续存在请联系 QQ 204034。")
		log.Fatalln(fmt.Sprintf("错误内容：%s", err.Error()))
	}
	
	err = jsoniter.Unmarshal(response.Bytes, &hostResult)
	
	if err != nil {
		log.Fatalln(fmt.Sprintf("无法解析服务器返回的数据，请稍后重试。服务器返回内容：%s", response.Text))
	}
	
	log.Println(fmt.Sprintf("获得服务器 IP：%s", hostResult.Ip))
	
	// 解析 Hosts 文件
	hostFile := GetFileData(hostPath)
	
	for _, hostName := range hostResult.Host {
		log.Println(fmt.Sprintf("已更新 %s 的 IP 地址为 %s", hostName, hostResult.Ip))
		hostFile = SetHosts(hostFile, hostName, hostResult.Ip)
	}
	
	// 补齐结尾换行
	hostFile += "\n"
	
	// 写入到文件中
	if SetFileData(hostPath, hostFile) {
		log.Println("Hosts 文件写入成功！")
	} else {
		log.Println("Hosts 文件写入失败！")
		if runtime.GOOS == "windows" {
			log.Fatalln("请检查是否用 [右键 -> 以管理员身份运行] 本程序，如果使用管理员权限运行后还是出错，请手动打开 C:/Windows/System32/drivers/etc/ 并右键查看 hosts 文件是否被设置为了只读，如果是请将只读取消然后重试。如果问题持续存在，请联系 QQ 204034。")
		} else {
			log.Fatalln("请检查 /etc/hosts 文件是否是只读状态，以及您当前的用户组是否有读写权限，请尝试使用 root 身份执行本操作！")
		}
	}
	
	fmt.Printf("\n请按回车键退出...")
    os.Stdin.Read(make([]byte, 1))
}

// 修复 Let's Encrypt 服务器被劫持问题
func FixLetsEncrypt() (bool) {
	
	// 读取 Hosts 文件
	hostFile := GetFileData(hostPath)
	hostFile  = SetHosts(hostFile, "ocsp.int-x3.letsencrypt.org", "23.46.210.176")
	
	return SetFileData(hostPath, hostFile)
}

// 设置 Hosts
func SetHosts(hostData string, hostName string, ipAddr string) (string) {
	
	// 是否找到 Hosts
	foundHost := false
	
	allLines := strings.Split(hostData, "\n")
	newHosts := ""
	
	for _, line := range allLines {
		if strings.Contains(line, fmt.Sprintf(" %s", hostName)) {
			newHosts += fmt.Sprintf("%s %s\n", ipAddr, hostName)
			foundHost = true
		} else if line != "" {
			newHosts += line + "\n"
		}
	}
	
	if !foundHost {
		newHosts += fmt.Sprintf("%s %s\n", ipAddr, hostName)
	}
	
	return newHosts
}

// 读文件
func GetFileData(fileName string) (string) {
    f, err := ioutil.ReadFile(fileName)
    if err != nil {
        log.Fatalln("Failed to read file", err)
    }
    return string(f)
}

// 写文件
func SetFileData(fileName string, data string) (bool) {
	var bytes = []byte(data)
	err := ioutil.WriteFile(fileName, bytes, 0666)
	if err != nil {
		return false
	}
	return true
}

// 定义服务器返回文件类型
type RequestResult struct {
	Ip   string
	Host []string
}