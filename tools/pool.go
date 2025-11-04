package tools

import (
	"log"
	"sync"
)

type Pools struct {
	works chan func()
	nums  int
	wg    sync.WaitGroup
}

// 使用线程池进行端口测试
func NewPools(n int) *Pools {
	p := &Pools{
		works: make(chan func(), 1000), // 任务队列
		nums:  n,                       // worker 数量
	}
	return p
}
func NewDefaultPools() *Pools {
	p := &Pools{
		works: make(chan func(), 1000), // 任务队列
		nums:  100,                     // worker 数量
	}
	return p
}

// 像协程池添加任务
func (p *Pools) Add(f func()) {
	p.works <- f
}

// worker 的执行逻辑
func (p *Pools) exec() {
	defer p.wg.Done()

	for f := range p.works {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("worker panic: %v", r)
				}
			}()
			f()
		}()
	}
}

// 启动协程池
func (p *Pools) Run() {
	for i := 0; i < p.nums; i++ {
		p.wg.Add(1)
		go p.exec()
	}
}

// 安全关闭协程池
func (p *Pools) Stop() {
	close(p.works)
	p.wg.Wait()
}
