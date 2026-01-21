package broadcast_test

import (
	"happyladysauce/broadcast"
	"testing"
)

func TestBroadcast(t *testing.T) {
	broadcast.Broadcast()
	broadcast.CutDownLatch()
	broadcast.CondSignal()
}