package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetMachineNo 获取机器号，每个POD的ip都不一致，不同的pod生成唯一的机器号，可以提供给雪花算法！
// Node IP：Node节点的IP地址，即物理网卡的IP地址。
// Pod IP：Pod的IP地址，即docker容器的IP地址，此为虚拟IP地址。
// Cluster IP：Service的IP地址，此为虚拟IP地址。
func GetMachineNo() int {
	ip := os.Getenv("HOST_IP")
	if ip == "" {
		ip = "127.0.0.1"
	}
	fmt.Println("HOST_IP:", ip)
	ipSegs := strings.Split(ip, ".")
	number := ipSegs[len(ipSegs)-1]
	no, _ := strconv.Atoi(number)
	return no
}
