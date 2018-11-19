package ndnping

/*
#include "input.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"ndn-dpdk/appinit"
	"ndn-dpdk/dpdk"
	"ndn-dpdk/iface"
	"ndn-dpdk/iface/createface"
)

type App struct {
	Tasks []Task
	rxls  []iface.IRxLooper
}

func NewApp(cfg []TaskConfig) (app *App, e error) {
	app = new(App)

	appinit.StartRxl = app.addRxl
	if e = appinit.EnableCreateFace(createface.GetDefaultConfig()); e != nil {
		return nil, e
	}

	var faceCreateArgs []createface.CreateArg
	for _, taskCfg := range cfg {
		faceCreateArgs = append(faceCreateArgs, taskCfg.Face)
	}
	faces, e := createface.Create(faceCreateArgs...)
	if e != nil {
		return nil, e
	}

	for i, taskCfg := range cfg {
		task, e := newTask(taskCfg, faces[i])
		if e != nil {
			return nil, fmt.Errorf("[%d] init error: %v", i, e)
		}
		if faceKind := task.Face.GetFaceId().GetKind(); faceKind != iface.FaceKind_Eth {
			return nil, fmt.Errorf("[%d] FaceKind %v is not supported", i, faceKind)
		}
		app.Tasks = append(app.Tasks, task)
	}

	return app, nil
}

func (app *App) addRxl(rxl iface.IRxLooper) (usr interface{}, e error) {
	app.rxls = append(app.rxls, rxl)
	return nil, nil
}

func (app *App) Launch() {
	for _, rxl := range app.rxls {
		app.launchRxl(rxl)
	}
	for _, task := range app.Tasks {
		task.Launch()
	}
}

func (app *App) launchRxl(rxl iface.IRxLooper) {
	minFaceId := iface.FACEID_MAX
	maxFaceId := iface.FACEID_MIN
	for _, faceId := range rxl.ListFacesInRxLoop() {
		if faceId < minFaceId {
			minFaceId = faceId
		}
		if faceId > maxFaceId {
			maxFaceId = faceId
		}
	}

	inputC := C.NdnpingInput_New(C.uint16_t(minFaceId), C.uint16_t(maxFaceId), C.unsigned(rxl.GetNumaSocket()))
	for i, task := range app.Tasks {
		entryC := C.__NdnpingInput_GetEntry(inputC, C.uint16_t(task.Face.GetFaceId()))
		if entryC == nil {
			continue
		}
		if task.Client != nil {
			queue, e := dpdk.NewRing(fmt.Sprintf("client-rx-%d", i), 256,
				task.Face.GetNumaSocket(), true, true)
			if e != nil {
				panic(e)
			}
			entryC.clientQueue = (*C.struct_rte_ring)(queue.GetPtr())
			task.Client.c.rxQueue = entryC.clientQueue
		}
		if task.Server != nil {
			queue, e := dpdk.NewRing(fmt.Sprintf("server-rx-%d", i), 256,
				task.Face.GetNumaSocket(), true, true)
			if e != nil {
				panic(e)
			}
			entryC.serverQueue = (*C.struct_rte_ring)(queue.GetPtr())
			task.Server.c.rxQueue = entryC.serverQueue
		}
	}

	appinit.MustLaunch(func() int {
		rxl.RxLoop(64, unsafe.Pointer(C.NdnpingInput_FaceRx), unsafe.Pointer(inputC))
		return 0
	}, rxl.GetNumaSocket())
}

type Task struct {
	Face   iface.IFace
	Client *Client
	Server *Server
}

func newTask(cfg TaskConfig, face iface.IFace) (task Task, e error) {
	task.Face = face

	if cfg.Client != nil {
		task.Client, e = newClient2(task.Face, *cfg.Client)
		if e != nil {
			task.Close()
			return Task{}, e
		}
	}

	if cfg.Server != nil {
		task.Server = newServer(task.Face, *cfg.Server)
	}

	return task, nil
}

func (task *Task) Launch() {
	numaSocket := task.Face.GetNumaSocket()
	if task.Server != nil {
		appinit.MustLaunch(task.Server.Run, numaSocket)
	}
	if task.Client != nil {
		appinit.MustLaunch(task.Client.RunRx, numaSocket)
		appinit.MustLaunch(task.Client.RunTx, numaSocket)
	}
}

func (task *Task) Close() error {
	if task.Server != nil {
		task.Server.Close()
	}
	if task.Client != nil {
		task.Server.Close()
	}
	task.Face.Close()
	return nil
}