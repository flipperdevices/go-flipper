package flipper

import (
	pb "github.com/flipperdevices/go-flipper/internal/proto"
	pbapp "github.com/flipperdevices/go-flipper/internal/proto/app"
)

type app struct {
	f *Flipper
}

func (a *app) Start(name, args string) error {
	req := &pb.Main{
		Content: &pb.Main_AppStartRequest{
			AppStartRequest: &pbapp.StartRequest{
				Name: name,
				Args: args,
			},
		},
	}
	_, err := a.f.call(nil, req)
	return err
}

func (a *app) IsLocked() (bool, error) {
	req := &pb.Main{
		Content: &pb.Main_AppLockStatusRequest{},
	}
	res, err := a.f.call(nil, req)
	if err != nil {
		return false, err
	}
	return res[0].(*pb.Main_AppLockStatusResponse).AppLockStatusResponse.Locked, nil
}
