package tools

import (
	"crypto/tls"
	"fmt"
	"lain-cli/base"
	"net"
	"strconv"
	"sync"
	"time"
)

type info struct {
	IP   string
	PORT int
}

var mux sync.Mutex
var c sync.Once
var Set = make(map[info]string, 0)
var openFlag = false
var LocalIP string

// 可能是并发安全的set
func writeSet(ip string, port int, proc string) {
	mux.Lock()
	if proc == "" {
		proc = "NO PROC"
	}
	data := info{
		IP:   ip,
		PORT: port,
	}
	Set[data] = proc
	mux.Unlock()
}

func SetScannerOpen() {
	openFlag = true
}

func scanPortHttp(ip string, port int, timeout time.Duration) {

	address := net.JoinHostPort(ip, strconv.Itoa(port))
	// fmt.Println(address)

	tlsflag := true
	var buf [1024]byte
	for k, v := range base.Proc {

		// localAddr := &net.TCPAddr{IP: net.ParseIP(LocalIP)}
		// dialer := &net.Dialer{
		// 	LocalAddr: localAddr,
		// 	Timeout:   2 * time.Second, // 建议 2s-3s 更稳定
		// }
		// conn, err := dialer.Dial("tcp", address)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err != nil {
			// conn.Close()
			// fmt.Println("conn err", err)
			continue
		}
		// 清空
		buf = [1024]byte{}
		conn.Write(v)
		conn.SetReadDeadline(time.Now().Add(timeout))
		n, err := conn.Read(buf[:])
		if err != nil {
			if openFlag {
				writeSet(ip, port, "OPEN,NOPROCESS")
			}
			conn.Close()

			continue
		}
		// fmt.Println(string(buf[:n]))
		proc := base.GetInfo(k, buf[:n])
		if proc == "" {
			conn.Close()
			continue
		}
		writeSet(ip, port, proc)
		tlsflag = false
		conn.Close()
		break
	}
	if tlsflag {
		localAddr := &net.TCPAddr{IP: net.ParseIP(LocalIP)}
		dialer := &net.Dialer{LocalAddr: localAddr, Timeout: timeout}
		for k, v := range base.Proc {
			tlsConn, tlsErr := tls.DialWithDialer(dialer, "tcp", address, &tls.Config{
				InsecureSkipVerify: true,
			})
			if tlsErr != nil {
				// tlsConn.Close()
				// fmt.Println("tls conn error", tlsErr)
				continue
			}
			// 清空
			buf = [1024]byte{}
			tlsConn.Write(v)
			tlsConn.SetReadDeadline(time.Now().Add(timeout))
			n, err := tlsConn.Read(buf[:])
			if err != nil {
				// fmt.Println(k, "[read err]", err)
				if openFlag {
					writeSet(ip, port, "OPEN,NOPROCESS")
				}
				tlsConn.Close()
				continue
			}
			proc := base.GetInfo(k, buf[:n])
			if proc == "" {
				tlsConn.Close()
				continue
			}
			writeSet(ip, port, proc)
			tlsConn.Close()
			break
		}

	}
}
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("没有找到合适的本地 IP 地址")
}

func Scan(ip string, port int) {
	c.Do(func() {
		base.InitProc(ip)
		i, err := getLocalIP()
		if err != nil {
			panic("Local IP Got error")
		}
		LocalIP = i
	})
	scanPortHttp(ip, port, 1*time.Second)
}
