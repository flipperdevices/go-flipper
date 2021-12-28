package flipper

import (
	"time"

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

func (s *system) GetDateTime() (*time.Time, error) {
	req := &pb.Main{
		Content: &pb.Main_SystemGetDatetimeRequest{},
	}
	res, err := s.f.call(nil, req)
	if err != nil {
		return nil, err
	}
	content := res[0].(*pb.Main_SystemGetDatetimeResponse).SystemGetDatetimeResponse.Datetime
	t := time.Date(int(content.Year), time.Month(content.Month), int(content.Day),
		int(content.Hour), int(content.Minute), int(content.Second), 0,
		time.Local)
	return &t, nil
}

func (s *system) SetDateTime(t time.Time) error {
	req := &pb.Main{
		Content: &pb.Main_SystemSetDatetimeRequest{
			SystemSetDatetimeRequest: &pbsystem.SetDateTimeRequest{
				Datetime: &pbsystem.DateTime{
					Hour:    uint32(t.Hour()),
					Minute:  uint32(t.Minute()),
					Second:  uint32(t.Second()),
					Day:     uint32(t.Day()),
					Month:   uint32(t.Month()),
					Year:    uint32(t.Year()),
					Weekday: uint32(t.Weekday()),
				},
			},
		},
	}
	_, err := s.f.call(nil, req)
	return err
}

func (s *system) PlayAudiovisualAlert() error {
	req := &pb.Main{
		Content: &pb.Main_SystemPlayAudiovisualAlertRequest{
			SystemPlayAudiovisualAlertRequest: &pbsystem.PlayAudiovisualAlertRequest{},
		},
	}
	_, err := s.f.call(nil, req)
	return err
}
