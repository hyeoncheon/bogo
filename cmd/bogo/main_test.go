package main

import (
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/hyeoncheon/bogo/internal/common"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	r := require.New(t)
	r.Nil(nil)

	opts := common.DefaultOptions()
	r.NotNil(opts)
	c, _ := common.NewDefaultContext(&opts)
	r.NotNil(c)

	indicator := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		run(c, &opts)
		indicator = true
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	wg.Wait()
	r.True(indicator)
}
