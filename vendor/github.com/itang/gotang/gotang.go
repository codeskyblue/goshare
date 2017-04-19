package gotang

import (
	"fmt"
	"time"
)

func Time(block func()) {
	start := time.Now().UnixNano()
	block()
	end := time.Now().UnixNano()

	ms := (end - start) / 1000000.0
	fmt.Printf("Elapsed time: %v msecs", ms)
}

// 执行IO操作并设定超时时间，超时返回超时错误
func DoIOWithTimeout(iof func() error, t time.Duration) error {
	timeout := time.NewTicker(t)
	defer timeout.Stop()

	done := make(chan error)
	go func() {
		done <- iof()
	}()

	select {
	case <-timeout.C:
		return fmt.Errorf("Do IO timeout: %v", t)
	case err := <-done:
		return err
	}
}

func TruncStr(s string, le int, a string) string {
	if ulen(s) < le {
		return substr(s, 0, le)
	}
	return substr(s, 0, le) + a
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func ulen(s string) int {
	runes := []rune(s)
	return len(runes)
}
