package shutdown

import (
	"os"
	"syscall"
	"testing"
)

func TestGracefulShutdown(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
		sig  []os.Signal
	}{
		{
			name: "default signal",
			fn:   func() {},
			sig:  []os.Signal{},
		},
		{
			name: "custom signal",
			fn:   func() {},
			sig:  []os.Signal{os.Interrupt},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stop := GracefulShutdown(tt.fn, tt.sig...)
			if stop == nil {
				t.Error("stop channel is nil")
			}

			err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			if err != nil {
				t.Error(err)
			}

			_, ok := <-stop
			if ok {
				t.Error("stop channel is not closed")
			}
		})
	}
}
