package blogger

import (
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestNewBlogger(t *testing.T) {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			blogger := NewBlogger("/tmp/test.log", L_DEBUG)
			blogger.AddBase("ip", "192.168.1.1")
			blogger.RequestLogid()
			blogger.Info("info")
			blogger.Debug("debug")
			blogger.Warning("warning")
			blogger.Flush()
		}()
	}
	wg.Wait()
}