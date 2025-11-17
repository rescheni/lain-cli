package tools

import (
	"bufio"
	"fmt"

	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	config "github.com/rescheni/lain-cli/config"
	"github.com/rescheni/lain-cli/internal/utils"
	logs "github.com/rescheni/lain-cli/logs"
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
func InfoInit() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

const (
	i_info = iota
	i_Host
	i_OS
	i_Kernel
	i_Uptime
	i_Packages
	i_Shell
	i_Disk
	i_CPU
	i_GPU
	i_RAM
	i_IP
)

func BasePrint() {
	logoConfurl := utils.ExpandPath(config.Conf.Logo.Logo_txt)
	txt, err := os.OpenFile(logoConfurl, os.O_RDONLY, 0660)
	if err != nil {
		logs.Err("open logo txt err", err)
		return
	}
	defer txt.Close()

	keys := []string{
		"",
		"Host",
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
	txtScanner := bufio.NewScanner(txt)
	maxline := ""
	for i := range keys {
		if !txtScanner.Scan() {
			if err := txtScanner.Err(); err != nil {
				logs.Err("Read logo err", err)
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
		fmt.Printf("%"+maxline+"s"+"\t%s: %s", logoStyle.Render(string(line)), keyStyle.Render(keys[i]), valStyle.Render(vals[i]))
		fmt.Println()
	}
	for txtScanner.Scan() {
		line := txtScanner.Text()
		fmt.Printf("%"+maxline+"s", logoStyle.Render(string(line)))
		fmt.Println()
	}
	if err := txtScanner.Err(); err != nil {
		logs.Err("Read logo:", err)
	}

}

func getCpuinfo(ss []string) {
	cinfo := ""
	info, err := cpu.Info()
	if err != nil {
		logs.Err("Get cpu info", err)
		return
	}
	if len(info) == 0 {
		ss[i_CPU] = "nil"
	}
	cinfo += fmt.Sprintf("%s Cores %d", info[0].ModelName+" "+fmt.Sprintf("%.1f MHZ", info[0].Mhz), info[0].Cores)
	ss[i_CPU] = cinfo

}

func getMeminfo(ss []string) {

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		logs.Err("Get mem info err", err)
		ss[i_RAM] = "nil"
		return
	}

	all := memInfo.Total / 1024 / 1024
	used := memInfo.Used / 1024 / 1024

	ss[i_RAM] = fmt.Sprintf("[%.2f%%] %d MB/%d MB", memInfo.UsedPercent, used, all)

}

func getHostinfo(ss []string) {
	ht, err := host.Info()
	if err != nil {
		logs.Err("Get sys Info")
		return
	}
	kerversion, _ := host.KernelVersion()
	ss[i_Host] = os.Getenv("USER") + "@" + ht.Hostname
	ss[i_OS] = ht.OS + " " + ht.PlatformVersion + " " + ht.PlatformFamily
	ss[i_Kernel] = fmt.Sprintf("%s %s", ht.KernelArch, kerversion)
	uptime := ht.Uptime
	days := uptime / 86400
	hours := (uptime % 86400) / 3600
	mins := (uptime % 3600) / 60
	ss[i_Uptime] = fmt.Sprintf("%d days %d hours %d mins", days, hours, mins)
}
func getDiskinfo(ss []string) {
	path := "/"
	diskUsed, _ := disk.Usage(path)
	g := uint64(1024 * 1024 * 1024)
	all := diskUsed.Total / g
	used := diskUsed.Used / g
	ss[i_Disk] = fmt.Sprintf("[%.2f%%] %d GB / %d GB (%s)", diskUsed.UsedPercent, used, all, path)
}
func getIPinfo(ss []string) {

	ifaces, err := net.Interfaces()
	if err != nil {
		ss[i_IP] = "unknown"
		return
	}

	parts := make([]string, 0, 4)
	for _, iface := range ifaces {
		// 跳过 loopback 或未启动的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		name := iface.Name
		// 只关注  en, eth 以及常见 wifi 前缀（wlan/wlp/wl）或包含 wifi 的名
		lname := strings.ToLower(name)
		ok := strings.HasPrefix(lname, "en") || strings.HasPrefix(lname, "eth") ||
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
		ss[i_IP] = "unknown"
	} else {
		ss[i_IP] = strings.Join(parts, " ")
	}
}
func getShellinfo(ss []string) {
	ss[i_Shell] = os.Getenv("SHELL")
}
func getPackageinfo(ss []string) {
	packages := "0"
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", "brew list 2>/dev/null | wc -l")
		out, _ := cmd.Output()
		packages = strings.TrimSpace(string(out))
	} else if runtime.GOOS == "linux" {
		cmd_dpkg := exec.Command("sh", "-c", "dpkg -l | grep -c '^ii'")
		cmd_out, _ := cmd_dpkg.Output()
		// cmd_snap := exec.Command("sh", "-c", "snap list | wc -l")
		// snap_out, _ := cmd_snap.Output()
		packages = "dpkg:" + strings.TrimSpace(string(cmd_out)) // + " snap:" + strings.TrimSpace(string(snap_out))
	} else {

	}
	ss[i_Packages] = packages
}
func getGpuinfo(ss []string) {

	gpu := "unknown"
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", "system_profiler SPDisplaysDataType | grep 'Chipset Model' | head -1 | cut -d: -f2 | xargs")
		out, _ := cmd.Output()
		gpu = strings.TrimSpace(string(out))
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "nvidia-smi --query-gpu=name --format=csv,noheader,nounits | head -1 | xargs")
		out, err := cmd.Output()
		if err != nil {
			cmd := exec.Command("sh", "-c", "lspci | grep -i vga | head -1 | cut -d: -f3- | xargs")
			out, err := cmd.Output()
			if err != nil {
				gpu = "unknown"
			} else {
				gpu = strings.TrimSpace(string(out))
			}
		} else {
			gpu = strings.TrimSpace(string(out))
		}
	} else {

	}
	ss[i_GPU] = gpu
}
