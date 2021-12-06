package flipper

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/flipperdevices/go-flipper/internal/delimited"
	pb "github.com/flipperdevices/go-flipper/internal/proto"
)

type Flipper struct {
	rd      *delimited.Reader
	wr      *delimited.Writer
	mutex   sync.Mutex
	pending map[uint32]*rpcCall
	counter uint32
	timeout time.Duration

	System  system
	Storage storage
	App     app
	Gui     gui
}

type rpcCall struct {
	written          uint32
	progressCallback updateProgress
	res              []*pb.Main
	done             chan bool
	timer            *time.Timer
	err              error
}

type updateProgress func(read, written uint32)

func Connect(rw io.ReadWriter) (*Flipper, error) {
	return ConnectWithTimeout(rw, 10*time.Second)
}

func ConnectWithTimeout(rw io.ReadWriter, timeout time.Duration) (*Flipper, error) {
	f := &Flipper{
		rd:      delimited.NewReader(rw),
		wr:      delimited.NewWriter(rw),
		pending: make(map[uint32]*rpcCall, 1),
		counter: 1,
		timeout: timeout,
	}
	f.System = system{f: f}
	f.Storage = storage{f: f}
	f.App = app{f: f}
	f.Gui = gui{f: f}

	go f.read()

	err := f.System.Ping()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Flipper) StopSession() error {
	req := &pb.Main{
		Content: &pb.Main_StopSession{},
	}
	_, err := f.call(nil, req)
	return err
}

func (f *Flipper) sendUnsolicited(req *pb.Main) error {
	f.mutex.Lock()
	err := f.wr.PutProto(req)
	f.mutex.Unlock()
	return err
}

func (f *Flipper) call(progressCallback updateProgress, req ...*pb.Main) ([]interface{}, error) {
	f.mutex.Lock()

	id := f.counter
	f.counter++

	call := &rpcCall{
		done:             make(chan bool),
		progressCallback: progressCallback,
	}
	f.pending[id] = call

	for i, r := range req {
		r.CommandId = id
		r.HasNext = len(req) > i+1
		err := f.wr.PutProto(r)
		call.written = uint32(i + 1)
		if call.progressCallback != nil {
			call.progressCallback(uint32(len(call.res)), call.written)
		}
		if err != nil {
			delete(f.pending, id)
			f.mutex.Unlock()
			return nil, err
		}
	}
	f.mutex.Unlock()
	call.timer = time.NewTimer(f.timeout)
	select {
	case <-call.done:
	case <-call.timer.C:
		call.err = errors.New("request timeout")
	}
	if call.err != nil {
		return nil, call.err
	}
	if call.res[len(call.res)-1].CommandStatus != pb.CommandStatus_OK {
		return nil, errors.New(call.res[len(call.res)-1].CommandStatus.String())
	}
	var content []interface{}
	for _, r := range call.res {
		if !isValidResponse(req[0].Content, r.Content) {
			return nil, errors.New("wrong response type")
		}
		content = append(content, r.Content)
	}
	return content, nil
}

func (f *Flipper) handleUnsolicited(res *pb.Main) {
	switch c := res.Content.(type) {
	case *pb.Main_GuiScreenFrame:
		if f.Gui.frameCallback != nil {
			f.Gui.frameCallback(ScreenFrame{c.GuiScreenFrame.Data})
		}
		break
	default:
		//TODO handle unknown unsolicited packets?
	}
}

func (f *Flipper) read() {
	var err error
	for err == nil {
		var res pb.Main
		err = f.rd.NextProto(&res)
		if err != nil {
			err = errors.New("error reading message: " + err.Error())
			continue
		}
		if res.CommandId == 0 {
			f.handleUnsolicited(&res)
			continue
		}
		f.mutex.Lock()
		call := f.pending[res.CommandId]
		if !res.HasNext {
			delete(f.pending, res.CommandId)
		}
		f.mutex.Unlock()
		if call == nil {
			err = errors.New("no corresponding request found")
			continue
		}
		call.res = append(call.res, &res)
		if call.progressCallback != nil {
			call.progressCallback(uint32(len(call.res)), call.written)
		}
		if res.HasNext {
			call.timer.Reset(f.timeout)
		} else {
			call.done <- true
		}
	}
	f.mutex.Lock()
	for _, call := range f.pending {
		call.err = err
		call.done <- true
	}
	f.mutex.Unlock()
}
