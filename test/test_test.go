package test_test

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/showwin/speedtest-go/speedtest"
)

type Net_time struct {
	mbps float64
	time int64
}

func RunSpeedTest(downChan, upChan chan Net_time) {
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
			case downChan <- Net_time{
				mbps: rate.Mbps(),
				time: time.Since(start).Milliseconds(),
			}:
			default:
			}
		})
		dm.SetCallbackUpload(func(rate speedtest.ByteRate) {
			select {
			case upChan <- Net_time{
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

		fmt.Printf(
			"Latency: %s, Download: %.2f Mbps, Upload: %.2f Mbps\n",
			s.Latency,
			s.DLSpeed.Mbps(),
			s.ULSpeed.Mbps(),
		)

		s.Context.Reset()
	}

	close(downChan)
	close(upChan)
}

func TuiOpen(downChan, upChan chan Net_time) {
	if err := ui.Init(); err != nil {
		log.Fatalf("termui init failed: %v", err)
	}
	defer ui.Close()

	lc := widgets.NewPlot()
	lc.Title = "Speed Test Monitor (press q)"
	lc.Data = [][]float64{{0}, {0}}
	lc.LineColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	lc.AxesColor = ui.ColorYellow
	lc.SetRect(0, 0, 80, 10)

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
			if len(lc.Data) < 2 || len(lc.Data[0]) == 0 || len(lc.Data[1]) == 0 {
				continue
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

func TestSpeedTest(t *testing.T) {
	upChan := make(chan Net_time, 1000)
	downChan := make(chan Net_time, 1000)

	go func() {
		RunSpeedTest(downChan, upChan)

	}()
	TuiOpen(downChan, upChan)

}
