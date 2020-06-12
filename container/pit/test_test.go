package pit_test

import (
	"os"
	"testing"

	"ndn-dpdk/container/fib"
	"ndn-dpdk/container/fib/fibtest"
	"ndn-dpdk/container/pcct"
	"ndn-dpdk/container/pit"
	"ndn-dpdk/container/strategycode"
	"ndn-dpdk/core/testenv"
	"ndn-dpdk/dpdk/eal/ealtestenv"
	"ndn-dpdk/dpdk/pktmbuf"
	"ndn-dpdk/dpdk/pktmbuf/mbuftestenv"
	"ndn-dpdk/iface"
	"ndn-dpdk/ndn"
	"ndn-dpdk/ndn/ndntestenv"
)

func TestMain(m *testing.M) {
	mbuftestenv.Direct.Template.Update(pktmbuf.PoolConfig{Dataroom: 8000}) // needed for TestEntryLongName
	ealtestenv.InitEal()
	os.Exit(m.Run())
}

var (
	makeAR       = testenv.MakeAR
	makeInterest = ndntestenv.MakeInterest
	makeData     = ndntestenv.MakeData
)

type Fixture struct {
	Pit *pit.Pit

	fibFixture    *fibtest.Fixture
	emptyStrategy strategycode.StrategyCode
	EmptyFibEntry *fib.Entry
}

func NewFixture(pcctMaxEntries int) (fixture *Fixture) {
	fixture = new(Fixture)

	pcctCfg := pcct.Config{
		Id:         "TestPcct",
		MaxEntries: pcctMaxEntries,
	}
	pcct, e := pcct.New(pcctCfg)
	if e != nil {
		panic(e)
	}

	fixture.Pit = pit.FromPcct(pcct)

	fixture.fibFixture = fibtest.NewFixture(2, 4, 1)
	fixture.emptyStrategy = strategycode.MakeEmpty("empty")
	fixture.EmptyFibEntry = new(fib.Entry)
	return fixture
}

func (fixture *Fixture) Close() error {
	fixture.fibFixture.Close()
	return fixture.Pit.Pcct.Close()
}

// Return number of in-use entries in PCCT's underlying mempool.
func (fixture *Fixture) CountMpInUse() int {
	return fixture.Pit.GetMempool().CountInUse()
}

// Insert a PIT entry.
// Returns the PIT entry.
// If CS entry is found, returns nil and frees interest.
func (fixture *Fixture) Insert(interest *ndn.Interest) *pit.Entry {
	pitEntry, csEntry := fixture.Pit.Insert(interest, fixture.EmptyFibEntry)
	if csEntry != nil {
		ndntestenv.ClosePacket(interest)
		return nil
	}
	if pitEntry == nil {
		panic("Pit.Insert failed")
	}
	return pitEntry
}

func (fixture *Fixture) InsertFibEntry(name string, nexthop iface.FaceId) *fib.Entry {
	if _, e := fixture.fibFixture.Fib.Insert(fixture.fibFixture.MakeEntry(name,
		fixture.emptyStrategy, nexthop)); e != nil {
		panic(e)
	}
	return fixture.fibFixture.Fib.Find(ndn.MustParseName(name))
}
