package pingclient

import (
	"time"

	"ndn-dpdk/ndn"
)

// Client config.
type Config struct {
	Patterns []Pattern     // traffic patterns
	Interval time.Duration // sending interval
}

// Client pattern definition.
type Pattern struct {
	Weight int // weight of random choice, minimum is 1

	Prefix           *ndn.Name     // name prefix
	CanBePrefix      bool          // whether to set CanBePrefix
	MustBeFresh      bool          // whether to set MustBeFresh
	InterestLifetime time.Duration // InterestLifetime value, zero means default
	HopLimit         int           // HopLimit value, zero means default

	// If non-zero, request cached Data. This must appear after a pattern without SeqNumOffset.
	// The client derives sequece number by subtracting SeqNumOffset from the previous pattern's
	// sequence number. Sufficient CS capacity is necessary for Data to actually come from CS.
	SeqNumOffset int
}
