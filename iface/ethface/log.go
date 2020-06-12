package ethface

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"ndn-dpdk/core/logger"
	"ndn-dpdk/dpdk/ethdev"
)

var (
	log           = logger.New("ethface")
	makeLogFields = logger.MakeFields
	addressOf     = logger.AddressOf
)

func newPortLogger(ethDev ethdev.EthDev) logrus.FieldLogger {
	return logger.NewWithPrefix("ethface", fmt.Sprintf("port %s", ethDev))
}
