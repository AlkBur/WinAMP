package ui

import (
	"os"
	"fmt"
	"sync"
)

type Bufferer interface {
	Buffer() Buffer
}

var renderJobs chan []Bufferer
var renderLock sync.Mutex

func render(bs ...Bufferer) {
	defer func() {
		if e := recover(); e != nil {
			Close()
			fmt.Fprintf(os.Stderr, "Captured a panic(value=%v) when rendering Bufferer. Exit termui and clean terminal...\nPrint stack trace:\n\n", e)
			//debug.PrintStack()
			//gs, err := stack.ParseDump(bytes.NewReader(debug.Stack()), os.Stderr)
			//if err != nil {
			//	debug.PrintStack()
			//	os.Exit(1)
			//}
			//p := &stack.Palette{}
			//buckets := stack.SortBuckets(stack.Bucketize(gs, stack.AnyValue))
			//srcLen, pkgLen := stack.CalcLengths(buckets, false)
			//for _, bucket := range buckets {
			//	io.WriteString(os.Stdout, p.BucketHeader(&bucket, false, len(buckets) > 1))
			//	io.WriteString(os.Stdout, p.StackLines(&bucket.Signature, srcLen, pkgLen, false))
			//}
			os.Exit(1)
		}
	}()
	for _, b := range bs {

		buf := b.Buffer()
		// set cels in buf
		for p, c := range buf.CellMap {
			if p.In(buf.Area) {

				SetCell(p.X, p.Y, c.Ch, Attribute(c.Fg), Attribute(c.Bg))

			}
		}

	}

	renderLock.Lock()
	// render
	Flush()
	renderLock.Unlock()
}

func Render(bs ...Bufferer) {
	//go func() { renderJobs <- bs }()
	renderJobs <- bs
}



