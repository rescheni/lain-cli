package exec

import (
	"fmt"

	"lain-cli/logs"
	"lain-cli/tools"
	mui "lain-cli/ui"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var CompletedTasks atomic.Uint64
var TotalTasks uint64
var mux sync.Mutex

var Rfun func()

func RunNmap(ip string, begin, end int) {

	if begin > end {
		logs.Err(fmt.Sprintln(begin, ">", end))
		return
	}
	pool := tools.NewDefaultPools()
	pool.Run()

	testNums := end - begin
	TotalTasks = uint64(testNums)

	for i := begin; i <= end; i++ {
		pool.Add(func() {
			defer CompletedTasks.Add(1)
			tools.Scan(ip, i)
		})
	}
	pool.Stop()

}
func RunDefaultNmap(ip string) {

	pool := tools.NewDefaultPools()
	pool.Run()

	ports := []int{
		21, //	端口：FTP 文件传输服务
		22, //	端口：SSH协议、SCP（文件传输）、端口号重定向
		23, //	/tcp端口：TELNET 终端仿真服务
		25, //	端口：SMTP 简单邮件传输服务
		53, //	端口：DNS 域名解析服务
		69, //	/udp：TFTP
		80, //	/8080/3128/8081/9098端口：HTTP协议代理服务器
		11, //	/tcp端口：POP3（E-mail）
		11, //	端口：Network
		12, //	端口：NTP（网络时间协议）
		// 135、137、138、139端口： 局域网相关默认端口，应关闭
		161,   //	端口：SNMP（简单网络管理协议）
		389,   //	端口：LDAP（轻量级目录访问协议）、ILS（定位服务）
		443,   //	/tcp 443/udp：HTTPS服务器
		465,   //	端口：SMTP（简单邮件传输协议）
		873,   //	端口：rsync
		1080,  //端口：SOCKS代理协议服务器常用端口号、QQ
		1158,  //端口：ORACLE EMCTL
		1433,  ///tcp/udp端口：MS SQL*SERVER数据库server、MS SQL*SERVER数据库monitor
		1521,  //端口：Oracle 数据库
		2100,  //端口：Oracle XDB FTP服务
		3389,  //端口：WIN2003远程登录
		3306,  //端口：MYSQL数据库端口
		5432,  //	端口：postgresql数据库端口
		5601,  //	端口：kibana
		6379,  //	端口：Redis数据库端口
		8080,  //	端口：TCP服务端默认端口、JBOSS、TOMCAT、Oracle XDB（XML 数据库）
		8081,  //	端口：Symantec AV/Filter for MSE
		8888,  //	端口：Nginx服务器的端口
		9000,  //	端口：php-fpm
		9080,  //	端口：Webshpere应用程序
		9090,  //	端口：webshpere管理工具
		9200,  //	端口：Elasticsearch服务器端口
		10050, //	端口：zabbix_server 10050
		10051, //	端口：zabbix_agent
		11211, //	端口：memcache（高速缓存系统）
		27017, //	端口：mongoDB数据库默认端口
		22122, //	端口：fastdfs服务器默认端口
	}
	testNums := len(ports)
	TotalTasks = uint64(testNums)

	for _, i := range ports {
		pool.Add(func() {
			defer CompletedTasks.Add(1)
			tools.Scan(ip, i)
		})
	}
	pool.Stop()

}
func RunNmapPorts(ip string, ports ...int) {
	TotalTasks = uint64(len(ports))
	pool := tools.NewDefaultPools()
	pool.Run()

	for _, v := range ports {
		pool.Add(func() {
			defer CompletedTasks.Add(1)
			tools.Scan(ip, v)
		})
	}
	pool.Stop()
}
func OutScanner() {
	mux.Lock()
	info := make([][]string, 0)
	i := 0
	for k, v := range tools.Set {
		i++
		t := []string{fmt.Sprint(i), k.IP, fmt.Sprint(k.PORT), v}
		info = append(info, t)
	}
	mui.TuiPrintTable([]string{"----", "IP/Domain Name", "PORT", "PROC"}, info)
	mux.Unlock()

}

const (
	padding  = 0
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func Run() {
	m := models{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		logs.Err("Oh no!", err)
		os.Exit(1)
	}
	fmt.Println()
	OutScanner()

}

func runScanCmd() tea.Cmd {
	return func() tea.Msg {
		Rfun()
		return "scan finished"
	}
}

type tickMsg time.Time

type models struct {
	progress progress.Model
}

func (m models) Init() tea.Cmd {

	return tea.Batch(
		runScanCmd(),
		tickCmd(),
	)
}

func (m models) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:

		currentDone := CompletedTasks.Load()
		percent := float64(currentDone) / float64(TotalTasks)
		if percent >= 1.0 {

			return m, tea.Quit
		}
		cmd := m.progress.SetPercent(percent)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressmodels, cmd := m.progress.Update(msg)
		m.progress = progressmodels.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m models) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press ctrl + c to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
