package blogger

import (
	"github.com/pkg/profile"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestNewBlogger(t *testing.T) {
	bFile := NewBFile("/tmp/test.log", L_INFO)
	defer profile.Start().Stop()
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			blogger := NewBlogger(bFile)
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