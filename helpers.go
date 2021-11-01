package flipper

import (
	pb "github.com/flipperdevices/go-flipper/internal/proto"
	"reflect"
)

var reqResMap = map[reflect.Type]reflect.Type{
	reflect.TypeOf(&pb.Main_StopSession{}):                 reflect.TypeOf(&pb.Main_Empty{}),
	reflect.TypeOf(&pb.Main_PingRequest{}):                 reflect.TypeOf(&pb.Main_PingResponse{}),
	reflect.TypeOf(&pb.Main_StorageStatRequest{}):          reflect.TypeOf(&pb.Main_StorageStatResponse{}),
	reflect.TypeOf(&pb.Main_StorageListRequest{}):          reflect.TypeOf(&pb.Main_StorageListResponse{}),
	reflect.TypeOf(&pb.Main_StorageReadRequest{}):          reflect.TypeOf(&pb.Main_StorageReadResponse{}),
	reflect.TypeOf(&pb.Main_StorageWriteRequest{}):         reflect.TypeOf(&pb.Main_Empty{}),
	reflect.TypeOf(&pb.Main_StorageDeleteRequest{}):        reflect.TypeOf(&pb.Main_Empty{}),
	reflect.TypeOf(&pb.Main_StorageMkdirRequest{}):         reflect.TypeOf(&pb.Main_Empty{}),
	reflect.TypeOf(&pb.Main_StorageMd5SumRequest{}):        reflect.TypeOf(&pb.Main_StorageMd5SumResponse{}),
	reflect.TypeOf(&pb.Main_AppStartRequest{}):             reflect.TypeOf(&pb.Main_Empty{}),
	reflect.TypeOf(&pb.Main_AppLockStatusRequest{}):        reflect.TypeOf(&pb.Main_AppLockStatusResponse{}),
	reflect.TypeOf(&pb.Main_GuiStartScreenStreamRequest{}): reflect.TypeOf(&pb.Main_Empty{}),
	reflect.TypeOf(&pb.Main_GuiStopScreenStreamRequest{}):  reflect.TypeOf(&pb.Main_Empty{}),
	reflect.TypeOf(&pb.Main_GuiSendInputEventRequest{}):    reflect.TypeOf(&pb.Main_Empty{}),
}

func isValidResponse(req interface{}, res interface{}) bool {
	v, ok := reqResMap[reflect.TypeOf(req)]
	if !ok {
		return false
	}
	return v == reflect.TypeOf(res)
}
