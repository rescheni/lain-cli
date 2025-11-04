package exec

import (
	"fmt"

	mui "github.com/rescheni/lain-cli/internal/ui"
	"github.com/rescheni/lain-cli/logs"

	"log"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"

	"github.com/gizak/termui/v3/widgets"
	"github.com/showwin/speedtest-go/speedtest"
)

type net_time struct {
	mbps float64
	time int64
}

// 回调函数
var callback func()

func runSpeedTest(downChan, upChan chan net_time) {
	var speedtestClient = speedtest.New()
	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	start := time.Now()

	for _, s := range targets {
		s.PingTest(nil)
		dm := s.Context.Manager
		// 设置SpeedTrst 下载回调
		dm.SetCallbackDownload(func(rate speedtest.ByteRate) {
			select {
			// 发送到下载缓冲区
			case downChan <- net_time{
				mbps: rate.Mbps(),
				time: time.Since(start).Milliseconds(),
			}:
			default:
			}
		})
		// 设置SpeedTrst 上传回调
		dm.SetCallbackUpload(func(rate speedtest.ByteRate) {
			select {
			// 发送到上传缓冲区
			case upChan <- net_time{
				mbps: rate.Mbps(),
				time: time.Since(start).Milliseconds(),
			}:
			default:
			}
		})
		// 并行执行上下行测试
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.DownloadTest()
		}()
		go func() {
			defer wg.Done()
			s.UploadTest()
		}()
		wg.Wait()

		// 优先接收 SpeetTest 测试完成的结果给 Tui 的输出表格的库
		callback = func() {
			mui.TuiPrintTable([]string{"Latency", "Download", "Upload"}, [][]string{
				{fmt.Sprint(s.Latency), fmt.Sprintf("%.2fMbps", s.DLSpeed.Mbps()), fmt.Sprintf("%.2fMbps", s.ULSpeed.Mbps())},
			})
		}
		s.Context.Reset()
	}
	close(downChan)
	close(upChan)
}

func tuiOpen(downChan, upChan chan net_time) {
	if err := ui.Init(); err != nil {
		log.Fatalf("termui init failed: %v", err)
	}
	defer ui.Close()

	lc := widgets.NewPlot()
	lc.Title = "Speed Test Monitor (press q)"
	lc.Data = [][]float64{{0}, {0}}
	lc.LineColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	lc.AxesColor = ui.ColorYellow
	lc.SetRect(0, 0, 80, 20)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	// 临时缓冲，用于接收测速数据
	var downBuf, upBuf []float64

	// 初始化 UI 事件
	uiEvents := ui.PollEvents()
	for {
		// 开启 UI 接收缓冲区数据同步到ui
		select {
		case d, ok := <-downChan:
			if ok {
				downBuf = append(downBuf, d.mbps)
				if len(downBuf) > 100 {
					downBuf = downBuf[1:]
				}
			} else {
				downChan = nil
			}
		case u, ok := <-upChan:
			if ok {
				upBuf = append(upBuf, u.mbps)
				if len(upBuf) > 100 {
					upBuf = upBuf[1:]
				}
			} else {
				upChan = nil
			}
		case <-ticker.C:
			lc.Data[0] = append([]float64(nil), downBuf...)
			lc.Data[1] = append([]float64(nil), upBuf...)
			// 如果某条线长度不足，就补 0 或保持前值
			if len(lc.Data[0]) < 2 {
				if len(lc.Data[0]) == 0 {
					lc.Data[0] = []float64{0, 0}
				} else {
					lc.Data[0] = append(lc.Data[0], lc.Data[0][0])
				}
			}
			if len(lc.Data[1]) < 2 {
				if len(lc.Data[1]) == 0 {
					lc.Data[1] = []float64{0, 0}
				} else {
					lc.Data[1] = append(lc.Data[1], lc.Data[1][0])
				}
			}
			// 渲染
			ui.Render(lc)
			// 设置 退出指令
		case e := <-uiEvents:
			if e.ID == "q" || e.ID == "<C-c>" {
				return
			}
		}
		// 接收值空
		if downChan == nil && upChan == nil {
			break
		}
	}
}

func RunSpeedTestUI(nui bool) {
	upChan := make(chan net_time, 1000)
	downChan := make(chan net_time, 1000)

	go func() {
		runSpeedTest(downChan, upChan)
	}()
	if nui {
		logs.Info("使用终端测试工具")
		for down := range downChan {
			up := <-upChan
			fmt.Println("[  Type  ]   Speed    Times")
			fmt.Printf("[Download]   %.2fMbps  %dms\n", down.mbps, down.time)
			fmt.Printf("[ Upload ]   %.2fMbps  %dms\n", up.mbps, up.time)
		}
	} else {
		logs.Info("使用ui测速工具")
		tuiOpen(downChan, upChan)
	}
	callback()
}
