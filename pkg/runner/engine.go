package runner

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/monkey92t/afrog/v2/pkg/config"
	"github.com/monkey92t/afrog/v2/pkg/poc"
	"github.com/monkey92t/afrog/v2/pkg/protocols/http/retryhttpclient"
	"github.com/monkey92t/afrog/v2/pkg/result"
	"github.com/panjf2000/ants/v2"
)

var CheckerPool = sync.Pool{
	New: func() any {
		return &Checker{
			Options: &config.Options{},
			// OriginalRequest: &http.Request{},
			VariableMap: make(map[string]any),
			Result:      &result.Result{},
			CustomLib:   NewCustomLib(),
		}
	},
}

func (e *Engine) AcquireChecker() *Checker {
	c := CheckerPool.Get().(*Checker)
	c.Options = e.options
	c.Result.Output = e.options.Output
	return c
}

func (e *Engine) ReleaseChecker(c *Checker) {
	// *c.OriginalRequest = http.Request{}
	c.VariableMap = make(map[string]any)
	c.Result = &result.Result{}
	c.CustomLib = NewCustomLib()
	CheckerPool.Put(c)
}

type Engine struct {
	options *config.Options
	ticker  *time.Ticker
}

func NewEngine(options *config.Options) *Engine {
	engine := &Engine{
		options: options,
	}
	return engine
}

func (runner *Runner) Execute() {

	options := runner.options

	pocSlice := options.CreatePocList()

	options.Count += options.Targets.Len() * len(pocSlice)

	if options.Smart {
		options.SmartControl()
	}

	runner.engine.ticker = time.NewTicker(time.Second / time.Duration(options.RateLimit))
	var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(options.Concurrency, func(p any) {
		defer wg.Done()
		<-runner.engine.ticker.C

		tap := p.(*TransData)

		if len(tap.Target) > 0 && len(tap.Poc.Id) > 0 {
			runner.executeExpression(tap.Target, &tap.Poc)
		}

	})
	defer p.Release()

	for _, poc := range pocSlice {
		for _, t := range runner.options.Targets.List() {
			wg.Add(1)
			p.Invoke(&TransData{Target: t.(string), Poc: poc})
		}
	}

	wg.Wait()
}

func (runner *Runner) executeExpression(target string, poc *poc.Poc) {
	c := runner.engine.AcquireChecker()
	defer runner.engine.ReleaseChecker(c)

	defer func() {
		// https://github.com/monkey92t/afrog/v2/issues/7
		if r := recover(); r != nil {
			c.Result.IsVul = false
			runner.OnResult(c.Result)
		}
	}()

	c.Check(target, poc)
	runner.OnResult(c.Result)
}

type TransData struct {
	Target string
	Poc    poc.Poc
}

func JndiTest() bool {
	url := "http://" + config.ReverseJndi + ":" + config.ReverseApiPort + "/?api=test"
	resp, _, err := retryhttpclient.Get(url)
	if err != nil {
		return false
	}
	if strings.Contains(string(resp), "no") || strings.Contains(string(resp), "yes") {
		return true
	}
	return false
}

func CeyeTest() bool {
	url := fmt.Sprintf("http://%s.%s", "test", config.ReverseCeyeDomain)
	resp, _, err := retryhttpclient.Get(url)
	if err != nil {
		return false
	}
	if strings.Contains(string(resp), "\"meta\":") || strings.Contains(string(resp), "201") {
		return true
	}
	return false
}
