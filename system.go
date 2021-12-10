package flipper

import (
	pb "github.com/flipperdevices/go-flipper/internal/proto"
	pbsystem "github.com/flipperdevices/go-flipper/internal/proto/system"
)

type system struct {
	f *Flipper
}

type RebootMode pbsystem.RebootRequest_RebootMode

const (
	RebootModeOs  = RebootMode(pbsystem.RebootRequest_OS)
	RebootModeDfu = RebootMode(pbsystem.RebootRequest_DFU)
)

func (s *system) Ping() error {
	req := &pb.Main{
		Content: &pb.Main_SystemPingRequest{},
	}
	_, err := s.f.call(nil, req)
	return err
}

func (s *system) Reboot(mode RebootMode) error {
	req := &pb.Main{
		Content: &pb.Main_SystemRebootRequest{
			SystemRebootRequest: &pbsystem.RebootRequest{
				Mode: pbsystem.RebootRequest_RebootMode(mode),
			},
		},
	}
	_, err := s.f.call(nil, req)
	return err
}

func (s *system) DeviceInfo() (map[string]string, error) {
	req := &pb.Main{
		Content: &pb.Main_SystemDeviceInfoRequest{},
	}
	res, err := s.f.call(nil, req)
	if err != nil {
		return nil, err
	}

	deviceInfo := make(map[string]string)
	for _, r := range res {
		pair := r.(*pb.Main_SystemDeviceInfoResponse).SystemDeviceInfoResponse
		deviceInfo[pair.Key] = pair.Value
	}

	return deviceInfo, nil
}

func (s *system) FactoryReset() error {
	req := &pb.Main{
		Content: &pb.Main_SystemFactoryResetRequest{},
	}
	_, err := s.f.call(nil, req)
	return err
}
