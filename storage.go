package flipper

import (
	pb "github.com/flipperdevices/go-flipper/internal/proto"
	pbstorage "github.com/flipperdevices/go-flipper/internal/proto/storage"
)

const maxPayloadLength = 512

type storage struct {
	f *Flipper
}

type FileType pbstorage.File_FileType

const (
	FileTypeFile = FileType(pbstorage.File_FILE)
	FileTypeDir  = FileType(pbstorage.File_DIR)
)

type File struct {
	Type FileType
	Name string
	Size uint32
	Data []byte
}

func (s *storage) Info(path string) (total_space, free_space uint64, err error) {
	req := &pb.Main{
		Content: &pb.Main_StorageInfoRequest{
			StorageInfoRequest: &pbstorage.InfoRequest{
				Path: path,
			},
		},
	}
	res, err := s.f.call(nil, req)
	if err != nil {
		return 0, 0, err
	}

	f := res[0].(*pb.Main_StorageInfoResponse).StorageInfoResponse
	total_space = f.TotalSpace
	free_space = f.FreeSpace

	return
}

func (s *storage) Stat(path string) (*File, error) {
	req := &pb.Main{
		Content: &pb.Main_StorageStatRequest{
			StorageStatRequest: &pbstorage.StatRequest{
				Path: path,
			},
		},
	}
	res, err := s.f.call(nil, req)
	if err != nil {
		return nil, err
	}

	f := res[0].(*pb.Main_StorageStatResponse).StorageStatResponse.File

	return &File{
		Type: FileType(f.Type),
		Name: f.Name,
		Size: f.Size,
		Data: f.Data,
	}, nil
}

func (s *storage) List(path string) ([]*File, error) {
	req := &pb.Main{
		Content: &pb.Main_StorageListRequest{
			StorageListRequest: &pbstorage.ListRequest{
				Path: path,
			},
		},
	}
	res, err := s.f.call(nil, req)
	if err != nil {
		return nil, err
	}

	var items []*File
	for _, r := range res {
		for _, i := range r.(*pb.Main_StorageListResponse).StorageListResponse.File {
			items = append(items, &File{
				Type: FileType(i.Type),
				Name: i.Name,
				Size: i.Size,
				Data: i.Data,
			})
		}
	}

	return items, nil
}

func (s *storage) Read(path string, progressCallback func(bytesRead uint32)) ([]byte, error) {
	req := &pb.Main{
		Content: &pb.Main_StorageReadRequest{
			StorageReadRequest: &pbstorage.ReadRequest{
				Path: path,
			},
		},
	}
	res, err := s.f.call(func(read, written uint32) {
		if progressCallback != nil {
			progressCallback(read * maxPayloadLength)
		}
	}, req)
	if err != nil {
		return nil, err
	}

	var data []byte

	for _, part := range res {
		cf := part.(*pb.Main_StorageReadResponse).StorageReadResponse.File
		data = append(data, cf.Data...)
	}

	return data, nil
}

func (s *storage) Write(path string, data []byte, progressCallback func(bytesWritten uint32)) error {
	var requests []*pb.Main

	for i := 0; i < len(data); i += maxPayloadLength {
		end := i + maxPayloadLength
		if end > len(data) {
			end = len(data)
		}
		req := &pb.Main{
			Content: &pb.Main_StorageWriteRequest{
				StorageWriteRequest: &pbstorage.WriteRequest{
					Path: path,
					File: &pbstorage.File{
						Data: data[i:end],
					},
				},
			},
		}
		requests = append(requests, req)
	}

	_, err := s.f.call(func(read, written uint32) {
		if progressCallback == nil {
			return
		}
		written *= maxPayloadLength
		if written > uint32(len(data)) {
			written = uint32(len(data))
		}
		progressCallback(written)
	}, requests...)
	return err
}

func (s *storage) Delete(path string, recursive bool) error {
	req := &pb.Main{
		Content: &pb.Main_StorageDeleteRequest{
			StorageDeleteRequest: &pbstorage.DeleteRequest{
				Path:      path,
				Recursive: recursive,
			},
		},
	}
	_, err := s.f.call(nil, req)
	return err
}

func (s *storage) Mkdir(path string) error {
	req := &pb.Main{
		Content: &pb.Main_StorageMkdirRequest{
			StorageMkdirRequest: &pbstorage.MkdirRequest{
				Path: path,
			},
		},
	}
	_, err := s.f.call(nil, req)
	return err
}

func (s *storage) GetMd5Sum(path string) (string, error) {
	req := &pb.Main{
		Content: &pb.Main_StorageMd5SumRequest{
			StorageMd5SumRequest: &pbstorage.Md5SumRequest{
				Path: path,
			},
		},
	}
	res, err := s.f.call(nil, req)
	if err != nil {
		return "", err
	}
	return res[0].(*pb.Main_StorageMd5SumResponse).StorageMd5SumResponse.Md5Sum, nil
}

func (s *storage) Rename(oldPath, newPath string) error {
	req := &pb.Main{
		Content: &pb.Main_StorageRenameRequest{
			StorageRenameRequest: &pbstorage.RenameRequest{
				OldPath: oldPath,
				NewPath: newPath,
			},
		},
	}
	_, err := s.f.call(nil, req)
	return err
}
