package exec

import (
	"fmt"
	mui "lain-cli/ui"

	"log"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"

	"github.com/gizak/termui/v3/widgets"
	"github.com/showwin/speedtest-go/speedtest"
)

// /
type net_time struct {
	mbps float64
	time int64
}

var callback func()

func runSpeedTest(downChan, upChan chan net_time) {
	var speedtestClient = speedtest.New()
	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	start := time.Now()

	for _, s := range targets {
		s.PingTest(nil)
		dm := s.Context.Manager
		// 设置回调
		dm.SetCallbackDownload(func(rate speedtest.ByteRate) {
			select {
			case downChan <- net_time{
				mbps: rate.Mbps(),
				time: time.Since(start).Milliseconds(),
			}:
			default:
			}
		})
		dm.SetCallbackUpload(func(rate speedtest.ByteRate) {
			select {
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

		callback = func() {
			// fmt.Printf(
			// 	"Latency: %s, Download: %.2f Mbps, Upload: %.2f Mbps\n",
			// 	s.Latency,
			// 	s.DLSpeed.Mbps(),
			// 	s.ULSpeed.Mbps(),
			// )
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

	uiEvents := ui.PollEvents()

	for {
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
			ui.Render(lc)
		case e := <-uiEvents:
			if e.ID == "q" || e.ID == "<C-c>" {
				return
			}
		}
		if downChan == nil && upChan == nil {
			break
		}
	}
}

func RunSpeedTestUI() {
	upChan := make(chan net_time, 1000)
	downChan := make(chan net_time, 1000)

	go func() {
		runSpeedTest(downChan, upChan)

	}()
	tuiOpen(downChan, upChan)
	callback()
}

func RunSpeedTestNUI() {
	upChan := make(chan net_time, 1000)
	downChan := make(chan net_time, 1000)

	go func() {
		runSpeedTest(downChan, upChan)
	}()

	for down := range downChan {
		up := <-upChan
		fmt.Printf("[Download] %.2fMbps  %dms\n", down.mbps, down.time)
		fmt.Printf("[Upload]   %.2fMbps  %dms\n", up.mbps, up.time)
	}
	callback()
}
