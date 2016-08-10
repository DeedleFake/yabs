package main

import (
	"context"
	"os"
	"os/signal"
	"time"
)

func SignalContext(ctx context.Context, sig ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		c := make(chan os.Signal, 1)
		signal.Notify(c, sig...)
		defer signal.Stop(c)

		select {
		case <-c:
		case <-ctx.Done():
		}
	}()

	return ctx
}

type FileInfoByName []os.FileInfo

func (fi FileInfoByName) Len() int {
	return len(fi)
}

func (fi FileInfoByName) Swap(i1, i2 int) {
	fi[i1], fi[i2] = fi[i2], fi[i1]
}

func (fi FileInfoByName) Less(i1, i2 int) bool {
	return fi[i1].Name() < fi[i2].Name()
}

type FileInfoByTimestamp struct {
	fi []os.FileInfo
	f  string
}

func (fi FileInfoByTimestamp) Len() int {
	return len(fi.fi)
}

func (fi FileInfoByTimestamp) Swap(i1, i2 int) {
	fi.fi[i1], fi.fi[i2] = fi.fi[i2], fi.fi[i1]
}

func (fi FileInfoByTimestamp) Less(i1, i2 int) bool {
	t1, err := time.Parse(fi.f, fi.fi[i1].Name())
	if err != nil {
		return true
	}

	t2, err := time.Parse(fi.f, fi.fi[i2].Name())
	if err != nil {
		return false
	}

	return t1.After(t2)
}
