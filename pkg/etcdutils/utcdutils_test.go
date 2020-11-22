package etcdutils

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	e := NewETC("../../certs")

	r := e.EtcdRun()
	fmt.Println(r)
}
