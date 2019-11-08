INCLUDEFLAGS= "-I/usr/local/include/dpdk -I/usr/include/dpdk"
BPFCC = "clang-6.0"
BPFFLAGS = "-O2 -target bpf #{INCLUDEFLAGS} -Wno-int-to-void-pointer-cast"

desc "Generate **/cgostruct.go"
task "cgostruct"
Rake::FileList["**/cgostruct.in.go"].each do |f|
  name = f.pathmap("%d/cgostruct.go")
  file name => f do |t|
    sh "cd #{f.pathmap("%d")}; go tool cgo -godefs -- #{INCLUDEFLAGS} cgostruct.in.go | gofmt -s > cgostruct.go; rm -rf _obj"
  end
  task "cgostruct" => name
end

CDeps = {}
CDeps["app/fwdp"] = ["container/fib", "container/ndt", "container/pcct", "iface"]
CDeps["app/fetch"] = ["container/mintmr"]
CDeps["app/ping"] = ["app/pingclient", "app/pingserver"]
CDeps["app/pingclient"] = ["iface"]
CDeps["app/pingserver"] = ["iface"]
CDeps["appinit"] = ["dpdk"]
CDeps["container/cs"] = ["container/pcct"]
CDeps["container/diskstore"] = ["spdk", "ndn"]
CDeps["container/fib"] = ["container/strategycode", "core/urcu", "dpdk", "ndn"]
CDeps["container/mintmr"] = ["dpdk"]
CDeps["container/mintmr/mintmrtest"] = ["container/mintmr"]
CDeps["container/ndt"] = ["ndn"]
CDeps["container/pcct"] = ["container/fib", "container/mintmr"]
CDeps["container/pit"] = ["container/pcct"]
CDeps["container/strategycode"] = ["core"]
CDeps["core"] = []
CDeps["core/coretest"] = ["core"]
CDeps["core/running_stat"] = []
CDeps["core/urcu"] = []
CDeps["dpdk"] = ["core"]
CDeps["dpdk/dpdktest"] = ["dpdk"]
CDeps["iface"] = ["ndn"]
CDeps["iface/ethface"] = ["iface"]
CDeps["iface/ifacetest"] = ["iface"]
CDeps["iface/mockface"] = ["iface"]
CDeps["iface/socketface"] = ["iface"]
CDeps["ndn"] = ["dpdk"]
CDeps["spdk"] = ["dpdk"]
CDeps["strategy"] = ["container/fib", "container/pcct", "ndn"]

desc "Generate **/cgoflags.go"
task "cgoflags"
CgoflagsPathmap = "%p/cgoflags.go"
CDeps.each do |key,value|
  name = key.pathmap(CgoflagsPathmap)
  file name => value.map{|v| v.pathmap(CgoflagsPathmap)} do |t|
    sh "mk/cgoflags.sh #{key} #{value.join(" ")}"
  end
  task "cgoflags" => name
end
Rake::Task["strategy".pathmap(CgoflagsPathmap)].clear

desc "Compile build/libndn-dpdk-*.a"
task "cbuilds"
ClibPathmap = "build/libndn-dpdk-%n.a"
CDeps.each do |key,value|
  name = key.pathmap(ClibPathmap)
  cSrc = Rake::FileList["#{key}/*.h", "#{key}/*.c"]
  cSrc = Rake::FileList["#{key}/api-*"] if key == "strategy"
  deps = value.map{|v| v.pathmap(ClibPathmap)} + cSrc
  if cSrc.length > 0 && !key.end_with?("test")
    file name => deps do |t|
      sh "mk/cbuild.sh #{key}"
    end
  else
    task name => deps
  end
  task "cbuilds" => name
end
Rake::Task["container/mintmr/mintmrtest".pathmap(ClibPathmap)].clear

file "ndn/error.h" => "ndn/error.tsv" do
  sh "ndn/make-error.sh"
end
file "ndn/tlv-type.h" => "ndn/tlv-type.tsv" do
  sh "ndn/make-tlv-type.sh"
end
task "ndn".pathmap(ClibPathmap) => ["ndn/error.h", "ndn/tlv-type.h"]

desc "Build forwarding strategies"
task "strategies" => "strategy/strategy_elf/bindata.go"
SgBpfPath = "build/strategy-bpf"
directory SgBpfPath
file "strategy/strategy_elf/bindata.go" do |t|
  sh "go-bindata -nomemcopy -pkg strategy_elf -prefix #{SgBpfPath} -o /dev/stdout #{SgBpfPath} | gofmt -s > #{t.name}"
end
SgDeps = [SgBpfPath, "strategy".pathmap(ClibPathmap)] + ["ndn", "container/fib", "container/pcct"].map{|v| v.pathmap(ClibPathmap)}
SgSrc = Rake::FileList["strategy/*.c"]
SgSrc.exclude("strategy/api*")
SgSrc.each do |f|
  name = f.pathmap("build/strategy-bpf/%n.o")
  file name => [f] + SgDeps do |t|
    sh "#{BPFCC} #{BPFFLAGS} -c #{t.source} -o #{t.name}"
  end
  task "strategy/strategy_elf/bindata.go" => name
end