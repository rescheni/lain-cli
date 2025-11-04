package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/rescheni/lain-cli/config"
	"github.com/rescheni/lain-cli/internal/utils"
	"github.com/rescheni/lain-cli/logs"
)

type llmctx struct {
	context  string
	ppid     int
	tempfile string
}

var LLMCTX llmctx

func (c *llmctx) Init() {
	CleanStaleContextFiles()
	c.ppid = os.Getppid()
	// 查找context
	c.tempfile = os.TempDir() + "/" + fmt.Sprint(c.ppid) + "-" + config.Conf.Context.Local
	_, err := os.Stat(c.tempfile)
	if err == nil {
		// 文件不存在
		// fmt.Println("文件存在")
		c.read()
	} else {
		// 文件存在
		// fmt.Println("文件不存在")
		c.context = utils.Getprompt()
		c.Add(c.context)

	}
	// fmt.Println(c.tempfile)
}

func (c *llmctx) Add(ctx string) {

	file, err := os.OpenFile(c.tempfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logs.Err("open file err", err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte("\n" + ctx))
	if err != nil {
		logs.Err("Write file err", err)
		return
	}
}

func (c *llmctx) read() {
	bctx, err := os.ReadFile(c.tempfile)
	if err != nil {
		logs.Err("Read file err", err)
		return
	}
	c.context = string(bctx)
}
func (c *llmctx) Getcontext() string {
	return c.context + "\n"
}

// / 上下文控制
func processExists(pid int) bool {
	if pid <= 0 {
		return false
	}
	// kill 0 不会终止进程，仅用于检测是否存在（Unix）
	err := syscall.Kill(pid, 0)
	if err == nil {
		return true
	}
	if err == syscall.ESRCH {
		return false
	}
	// EPERM 等其他错误说明进程存在但无权限操作，认为存在
	return true
}

func CleanStaleContextFiles() {
	// 查找上下文缓存
	pattern := filepath.Join(os.TempDir(), "*-"+config.Conf.Context.Local)
	files, err := filepath.Glob(pattern)
	if err != nil {
		logs.Err("glob err:", err)
		return
	}
	for _, f := range files {
		base := filepath.Base(f)
		parts := strings.SplitN(base, "-", 2)
		if len(parts) < 2 {
			continue
		}
		pidStr := parts[0]
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}
		if !processExists(pid) {
			if err := os.Remove(f); err != nil {
				// fmt.Println("remove stale context failed:", f, err)
			} else {
				// fmt.Println("removed stale context:", f)
			}
		}
	}
}
