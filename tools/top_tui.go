package tools

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

const (
	PID = iota
	USERNAME
	STATUS
	PROGARM_NAME
	MEM
	CPU
	TIMES
	PROGARM_LINE
)

// 排序规则
var sortFlag = MEM

func style() {
	tview.Styles.PrimitiveBackgroundColor = tcell.Color16 // 背景黑
	tview.Styles.ContrastBackgroundColor = tcell.Color16  // 较暗背景也黑
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorBlack
	tview.Styles.BorderColor = tcell.ColorWhite       // 边框灰
	tview.Styles.TitleColor = tcell.ColorWhite        // 标题白
	tview.Styles.PrimaryTextColor = tcell.ColorWhite  // 主文本白
	tview.Styles.SecondaryTextColor = tcell.ColorGray // 次文本灰
}

func OpenPerformance() {
	app := tview.NewApplication()
	style()
	// 左列
	netBox := tview.NewTextView()
	diskBox := tview.NewTextView()
	left := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(netBox, 0, 1, false).
		AddItem(diskBox, 0, 1, false)
	// 中间
	processTree := tview.NewTable()

	// 右列
	memBox := tview.NewTextView()
	cpuBox := tview.NewTextView()

	right := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(memBox, 0, 1, false).
		AddItem(cpuBox, 0, 1, false)

	// 整体横向布局
	mainFlex := tview.NewFlex().
		AddItem(left, 0, 2, false).       // 左列占比 2
		AddItem(processTree, 0, 5, true). // 中间主区占比 5
		AddItem(right, 0, 2, false)       // 右列占比 2

	cpuView(cpuBox, app)
	memView(memBox, app)
	processView(processTree)
	diskView(diskBox)
	netView(netBox)
	// 键盘退出
	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyCtrlC || e.Rune() == 'q' || e.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		}
		switch e.Rune() {
		case 'c', 'C':
			sortFlag = CPU
		case 'm', 'M':
			sortFlag = MEM
		case 't', 'T':
			sortFlag = TIMES
		case 'p', 'P':
			sortFlag = PROGARM_NAME
		}
		if e.Key() == tcell.KeyDelete {
			row, _ := processTree.GetSelection()
			if row <= 0 { // 跳过表头
				return e
			}
			pids := processTree.GetCell(row, 0)
			pidText := pids.Text
			pid, _ := strconv.Atoi(pidText)
			proc, err := os.FindProcess(pid)
			if err != nil {
				fmt.Println("no find Process")
				return e
			}
			err = syscall.Kill(pid, syscall.SIGKILL)
			if err != nil {
				proc.Kill()
			}
		}

		return e
	})

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func getCpuStatus(flag bool) <-chan []float64 {

	info := make(chan []float64, 100)

	go func() {
		cpu.Percent(0, flag)
		for {
			cpus, err := cpu.Percent(time.Second, flag)
			if err != nil {
				// continue
			}
			info <- cpus
			if len(info) > 0 {
				// 将数据发送到 channel
				info <- cpus

			}
		}
	}()
	return info
}

type memInfo struct {
	total     float64
	available float64
	used      float64
	free      float64
	cached    float64
}

var G = 1024 * 1024 * 1024.0

func getMemStatus() <-chan memInfo {

	memch := make(chan memInfo, 100)
	go func() {
		for {
			vm, _ := mem.VirtualMemory()

			tmp := memInfo{
				total:     float64(vm.Total) / G,
				available: float64(vm.Available) / G,
				used:      float64(vm.Used) / G,
				free:      float64(vm.Free) / G,
				cached:    float64(vm.Cached) / G,
			}
			memch <- tmp
			time.Sleep(1 * time.Second / 2)
		}

	}()
	return memch
}

type procinfo struct {
	info []*tview.TableCell
}
type SortByprocinfo []procinfo

func (a SortByprocinfo) Len() int      { return len(a) }
func (a SortByprocinfo) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByprocinfo) Less(i, j int) bool {
	return a[i].info[sortFlag].Text > a[j].info[sortFlag].Text
}

type procInfos struct {
	proinfo SortByprocinfo
}

func getProcessinfo() <-chan procInfos {
	ch := make(chan procInfos, 100)
	go func() {

		for {
			procs, err := process.Processes()
			if err != nil {
				fmt.Println("Get Proc err")
				return
			}
			procinfov := make(SortByprocinfo, 0)
			for _, v := range procs {
				chv := []*tview.TableCell{}
				U, _ := v.Username() // 判断用户
				osusername := os.Getenv("USER")
				if U != osusername && osusername != "root" {
					continue
				}

				// 1 PID
				pid := tview.NewTableCell(fmt.Sprint(v.Pid))
				chv = append(chv, pid)
				u, _ := v.Username()
				//2 USERNAME
				user := tview.NewTableCell(u)
				chv = append(chv, user)
				//3 STATUS
				S, _ := v.Status()
				ss := ""
				for _, i := range S {
					ss += i
				}
				s := tview.NewTableCell(ss)
				chv = append(chv, s)

				//4 PROGARMNAME
				progarm, _ := v.Name()
				l := len(progarm)
				if l >= 8 {
					l = 8
				}
				pro := tview.NewTableCell(progarm[:l])
				chv = append(chv, pro)
				//5 MEM
				mem, _ := v.MemoryInfo()
				smem := fmt.Sprintf("%5dM", mem.RSS/1024/1024)
				memery := tview.NewTableCell(smem)
				chv = append(chv, memery)

				//6 CPU
				cpup, _ := v.CPUPercent()
				scpup := fmt.Sprintf("%7.2f%%", cpup)
				scpuss := tview.NewTableCell(scpup)
				chv = append(chv, scpuss)

				//7 TIMES
				time, _ := v.Times()
				totalCPUTime := time.User + time.System // 这是一个 float64，例如 123.456
				minutes := int64(totalCPUTime / 60)
				seconds := totalCPUTime - (float64(minutes) * 60)
				timeStr := fmt.Sprintf("%5d:%05.2f", minutes, seconds)
				stime := tview.NewTableCell(timeStr).SetAlign(tview.AlignCenter) // 时间/数字通常右对齐
				chv = append(chv, stime)

				//8 PROGARM LINE
				cmd, _ := v.Cmdline()
				cmdline := tview.NewTableCell(cmd)
				chv = append(chv, cmdline)

				procinfov = append(procinfov, procinfo{
					info: chv,
				})
			}
			cha := procInfos{
				proinfo: procinfov,
			}
			sort.Sort(cha.proinfo)
			ch <- cha
			time.Sleep(1 * time.Second / 2)
		}
	}()
	return ch
}

func formatBytePerSec(bytes float64) string {

	const (
		KB = 1024
		MB = KB * 1024
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f G/s", bytes/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f M/s", bytes/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f K/s", bytes/KB)
	default:
		return fmt.Sprintf("%.2f B/s", bytes)
	}

}

type netinfo struct {
	name     string
	upload   float64 // 上行速度 (Bytes/sec)
	download float64 // 下行速度 (Bytes/sec)
}

func getNetInfo() <-chan netinfo {
	ch := make(chan netinfo, 10)

	go func() {
		var lastTime time.Time
		var lastSent, lastRecv uint64

		statsList, err := net.IOCounters(false)
		if err == nil && len(statsList) > 0 {
			lastSent = statsList[0].BytesSent
			lastRecv = statsList[0].BytesRecv

		}
		lastTime = time.Now()
		for {
			time.Sleep(1 * time.Second)
			newStatsList, err := net.IOCounters(false)
			if err != nil {
				continue
			}
			newTime := time.Now()
			duration := newTime.Sub(lastTime).Seconds()
			if duration == 0 {
				continue // 避免除以零
			}
			newStat := newStatsList[0]
			ch <- netinfo{
				name:     newStat.Name,
				upload:   float64(newStat.BytesSent-lastSent) / duration,
				download: float64(newStat.BytesRecv-lastRecv) / duration,
			}
			lastTime = newTime
			lastSent = newStat.BytesSent
			lastRecv = newStat.BytesRecv
		}

	}()
	return ch
}

type diskinfo struct {
	name  string
	write float64
	read  float64
}

func getdiskInfo() <-chan diskinfo {

	ch := make(chan diskinfo, 100)

	go func() {

		for {
			ctx := context.Background()

			counters1, err := disk.IOCountersWithContext(ctx)
			if err != nil {
				panic(err)
			}
			t1 := time.Now()
			time.Sleep(1 * time.Second)
			counters2, err := disk.IOCountersWithContext(ctx)
			if err != nil {
				panic(err)
			}
			t2 := time.Now()
			elapsed := t2.Sub(t1).Seconds()
			for name, cnt2 := range counters2 {
				cnt1, ok := counters1[name]
				if !ok {
					continue
				}
				deltaReadBytes := float64(cnt2.ReadBytes - cnt1.ReadBytes)
				deltaWriteBytes := float64(cnt2.WriteBytes - cnt1.WriteBytes)
				readSpeed := deltaReadBytes / elapsed   // 字节/秒
				writeSpeed := deltaWriteBytes / elapsed // 字节/秒
				ch <- diskinfo{
					name:  name,
					read:  readSpeed,
					write: writeSpeed,
				}
			}
		}

	}()

	return ch
}

func diskView(box *tview.TextView) {
	box.SetBorder(true).SetTitle("磁盘IO")
	ch := getdiskInfo()
	diskstat, _ := disk.Usage("/")
	diskinfos := fmt.Sprintf("Used: %s\nTotal:%s\n", formatBytePerSec(float64(diskstat.Used)), formatBytePerSec(float64(diskstat.Total)))

	go func() {
		for val := range ch {
			view := diskinfos
			view += "disk_name:" + val.name + "\n"
			view += "Read: " + formatBytePerSec(val.read) + "\n"
			view += "Write: " + formatBytePerSec(val.write) + "\n"
			box.SetText(view)
		}
	}()

}
func netView(box *tview.TextView) {
	box.SetBorder(true).SetTitle("网络IO")

	ch := getNetInfo()
	go func() {
		for val := range ch {
			format := ""
			format += "net-interface:" + val.name + "\n"
			format += "download:" + formatBytePerSec(val.download) + "\n"
			format += "upload:" + formatBytePerSec(val.upload) + "\n"
			box.SetText(format)
		}
	}()

}
func processView(box *tview.Table) {
	box.SetBorder(true).SetTitle("Process Info")
	box.SetFixed(1, 0)             // 固定第一行（表头）
	box.SetSelectable(true, false) // 允许行选
	headers := []string{"PID", "USER", " S ", "Program", "Mem", "CPU", "Time+", "Cmdline"}

	for col, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorBeige).
			SetAlign(tview.AlignCenter).
			SetBackgroundColor(tcell.ColorWhite).
			SetExpansion(1)
		// SetSelectable(false)
		box.SetCell(0, col, cell)
	}
	ch := getProcessinfo()

	go func() {
		for info := range ch {
			for i, val := range info.proinfo {
				for j, vval := range val.info {
					box.SetCell(i+1, j, vval)
				}
			}
		}
	}()
}
func memView(box *tview.TextView, app *tview.Application) {
	box.SetBorder(true).SetTitle("内存使用率")
	ch := getMemStatus()
	box.SetChangedFunc(func() {
		app.Draw()
	})
	go func() {
		for v := range ch {
			textV := ""
			textV += fmt.Sprintf("Total: %.1f G\n", v.total)
			textV += fmt.Sprintf("available: %.2f G\n", v.available)
			textV += fmt.Sprintf("used: %.2f G\n", v.used)
			textV += fmt.Sprintf("cached: %.2f G\n", v.cached)
			textV += fmt.Sprintf("free: %.2f G\n", v.free)
			box.SetText(textV)
		}
	}()
}

func cpuView(box *tview.TextView, app *tview.Application) {
	box.SetBorder(true).SetTitle("CPU使用率")
	ch := getCpuStatus(true)
	ch0 := getCpuStatus(false)
	box.SetChangedFunc(func() {
		app.Draw()
	})
	go func() {
		for vs := range ch {
			usageText := fmt.Sprintf("CPU %6.2f%%\n", (<-ch0)[0])
			usageText += "\n"
			for i, v := range vs {
				usageText += fmt.Sprintf("C %d %6.2f%%\t\t", i, v)
				if i%2 == 1 {
					usageText += "\n"
				}
			}
			app.QueueUpdate(func() {
				box.SetText(usageText)
			})
		}
	}()
}
