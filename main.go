package main

import (
    "fmt"
    "net"
    "os"
    "io"
    "os/exec"
    "syscall"
    "time"
)

func main() {
    for {
        // 检查3100端口是否存在
        if isPortOpen("127.0.0.1", 9100) {
            fmt.Println("Port 9100 is open.")
        } else {
            fmt.Println("Port 9100 is closed. Running a command...")
            // 在端口关闭时执行某个指令
            cmd := exec.Command("/opt/node_exporter/node_exporter")
            // 如果你需要将命令的输出打印到终端，可以使用以下代码
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
	    // 将命令设置为在后台运行
            cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	    // 创建一个管道，用于向命令的标准输入流写入数据
            stdin, err := cmd.StdinPipe()
            if err != nil {
                fmt.Printf("Error creating stdin pipe: %s\n", err)
                return
            }
            
            // 启动命令
            if err := cmd.Start(); err != nil {
                fmt.Printf("Command failed to start: %s\n", err)
            } else {
                // 向命令的标准输入流写入"Enter"键的换行符
                _, err := io.WriteString(stdin, "\n")
                if err != nil {
                    fmt.Printf("Error writing to stdin: %s\n", err)
                }
            }	
        }
        // 休眠一段时间后再次检查端口状态
        time.Sleep(5 * time.Second)
    }
}

// 检查指定主机和端口是否打开
func isPortOpen(host string, port int) bool {
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, 5*time.Second)
    if err != nil {
        return false
    }
    defer conn.Close()
    return true
}

