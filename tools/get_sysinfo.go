package tools

import (
	"bufio"
	"fmt"
	"lain-cli/config"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

var (
	logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E90FF")).Bold(true)
	keyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#bcbc05ff")).Bold(true)
	valStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffffff"))
)
var wg sync.WaitGroup

// TUI Color
func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

func BasePrint() {
	txt, err := os.OpenFile(config.Conf.Logo.Logo_txt, os.O_RDONLY, 0660)
	if err != nil {
		fmt.Println("open logo txt err")
		return
	}
	defer txt.Close()

	keys := []string{
		"",
		"host",
		"OS",
		"Kernel",
		"Uptime",
		"Packages",
		"Shell",
		"Disk",
		"CPU",
		"GPU",
		"RAM",
		"IP",
	}
	vals := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	}
	tasks := []func([]string){
		getHostinfo,
		getDiskinfo,
		getIPinfo,
		getGpuinfo,
		getShellinfo,
		getPackageinfo,
		getCpuinfo,
		getMeminfo,
	}
	for _, task := range tasks {
		wg.Add(1)
		go func() {
			task(vals)
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println()
	txtScanner := bufio.NewScanner(txt)
	maxline := ""
	for i := range keys {

		if !txtScanner.Scan() {
			if err := txtScanner.Err(); err != nil {
				fmt.Println("Read logo err", err)
				return
			}
			break
		}

		line := txtScanner.Text()
		if i == 0 {
			maxline = fmt.Sprint(len(string(line)))
			fmt.Printf("%"+maxline+"s", logoStyle.Render(string(line)))
			fmt.Println()
			continue
		}

		// fmt.Printf("%"+maxline+"s"+"\t%s:%s", string(line), keys[i], vals[i])
		fmt.Printf("%"+maxline+"s"+"\t%s: %s", logoStyle.Render(string(line)), keyStyle.Render(keys[i]), valStyle.Render(vals[i]))
		fmt.Println()

	}

	for txtScanner.Scan() {
		line := txtScanner.Text()
		fmt.Printf("%"+maxline+"s", logoStyle.Render(string(line)))
		fmt.Println()
	}
	if err := txtScanner.Err(); err != nil {
		fmt.Println("Read logo err:", err)
	}

}

func getCpuinfo(ss []string) {
	cinfo := ""
	num := 8
	info, err := cpu.Info()
	if err != nil {
		fmt.Println("Get cpu info err", err)
		return
	}
	if len(info) == 0 {
		ss[num] = "nil"
	}
	cinfo += fmt.Sprintf("%s Cores %d", info[0].ModelName+" "+fmt.Sprintf("%.1f MHZ", info[0].Mhz), info[0].Cores)
	ss[num] = cinfo

}

func getMeminfo(ss []string) {

	num := 10
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Get mem info err", err)
		ss[num] = "nil"
		return
	}

	all := memInfo.Total / 1024 / 1024
	used := memInfo.Used / 1024 / 1024

	ss[num] = fmt.Sprintf("[%.2f%%] %d MB/%d MB", memInfo.UsedPercent, used, all)

}

func getHostinfo(ss []string) {
	num_host := 1
	num_os := 2
	num_kernel := 3
	num_uptime := 4
	ht, err := host.Info()
	if err != nil {
		fmt.Println("Get sys Info err")
		return
	}
	kerversion, _ := host.KernelVersion()
	ss[num_host] = os.Getenv("USER") + "@" + ht.Hostname
	ss[num_os] = ht.OS + " " + ht.PlatformVersion + " " + ht.PlatformFamily
	ss[num_kernel] = fmt.Sprintf("%s %s", ht.KernelArch, kerversion)
	uptime := ht.Uptime
	days := uptime / 86400
	hours := (uptime % 86400) / 3600
	mins := (uptime % 3600) / 60
	ss[num_uptime] = fmt.Sprintf("%d days %d hours %d mins", days, hours, mins)

}

func getDiskinfo(ss []string) {

	num_disk := 7
	path := "/"
	diskUsed, _ := disk.Usage(path)
	g := uint64(1024 * 1024 * 1024)
	all := diskUsed.Total / g
	used := diskUsed.Used / g
	ss[num_disk] = fmt.Sprintf("[%.2f%%] %d GB / %d GB (%s)", diskUsed.UsedPercent, used, all, path)

}

func getIPinfo(ss []string) {
	num_ip := 11

	ifaces, err := net.Interfaces()
	if err != nil {
		ss[num_ip] = "unknown"

		return
	}

	parts := make([]string, 0, 4)
	for _, iface := range ifaces {
		// 跳过 loopback 或未启动的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		name := iface.Name
		// 只关注 en0, en1, eth0 以及常见 wifi 前缀（wlan/wlp/wl）或包含 wifi 的名
		lname := strings.ToLower(name)
		ok := name == "en0" || name == "en1" || name == "eth0" ||
			strings.HasPrefix(lname, "wl") || strings.HasPrefix(lname, "wlan") || strings.Contains(lname, "wifi")
		if !ok {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			// 只取 IPv4，若需要 IPv6 可移除此判断
			if ip.To4() == nil {
				continue
			}
			parts = append(parts, fmt.Sprintf("%s:%s", name, ip.String()))
		}
	}

	if len(parts) == 0 {
		ss[num_ip] = "unknown"
	} else {
		ss[num_ip] = strings.Join(parts, " ")
	}

}
func getShellinfo(ss []string) {

	ss[6] = os.Getenv("SHELL")

}
func getPackageinfo(ss []string) {
	cmd := exec.Command("sh", "-c", "brew list 2>/dev/null | wc -l")
	out, _ := cmd.Output()
	packages := strings.TrimSpace(string(out))
	ss[5] = packages

}

func getGpuinfo(ss []string) {
	cmd := exec.Command("sh", "-c", "system_profiler SPDisplaysDataType | grep 'Chipset Model' | head -1 | cut -d: -f2 | xargs")
	out, _ := cmd.Output()
	gpu := strings.TrimSpace(string(out))
	ss[9] = gpu

}
