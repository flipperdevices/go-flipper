package flipper

import (
	pb "github.com/flipperdevices/go-flipper/internal/proto"
)

type status struct {
	f *Flipper
}

func (s *status) Ping() error {
	req := &pb.Main{
		Content: &pb.Main_PingRequest{},
	}
	_, err := s.f.call(nil, req)
	return err
}
