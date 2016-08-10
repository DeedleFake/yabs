package main

import (
	"context"
	"os"
	"os/signal"
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
