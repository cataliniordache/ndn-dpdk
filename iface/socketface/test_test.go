package socketface_test

import (
	"net"
	"os"
	"testing"

	"ndn-dpdk/core/testenv"
	"ndn-dpdk/dpdk/eal/ealtestenv"
	"ndn-dpdk/iface/socketface"
)

var socketfaceCfg socketface.Config

func TestMain(m *testing.M) {
	ealtestenv.InitEal()

	socketfaceCfg = socketface.Config{
		TxqPkts:   64,
		TxqFrames: 64,
	}

	os.Exit(m.Run())
}

var makeAR = testenv.MakeAR

// Create net.Conn from file descriptor.
func makeConnFromFd(fd int) net.Conn {
	file := os.NewFile(uintptr(fd), "")
	if file == nil {
		panic(fd)
	}
	defer file.Close()

	conn, e := net.FileConn(file)
	if e != nil {
		panic(e)
	}
	return conn
}
