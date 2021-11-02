package flipper

import (
	pb "github.com/flipperdevices/go-flipper/internal/proto"
	pbgui "github.com/flipperdevices/go-flipper/internal/proto/gui"
	"image"
	"image/color"
	"image/draw"
)

type gui struct {
	f             *Flipper
	frameCallback updateFrame
}

type updateFrame func(frame ScreenStreamFrame)

type InputKey pbgui.InputKey

const (
	InputKeyUp    = InputKey(pbgui.InputKey_UP)
	InputKeyDown  = InputKey(pbgui.InputKey_DOWN)
	InputKeyRight = InputKey(pbgui.InputKey_RIGHT)
	InputKeyLeft  = InputKey(pbgui.InputKey_LEFT)
	InputKeyOk    = InputKey(pbgui.InputKey_OK)
	InputKeyBack  = InputKey(pbgui.InputKey_BACK)
)

type InputType pbgui.InputType

const (
	InputTypePress   = InputType(pbgui.InputType_PRESS)
	InputTypeRelease = InputType(pbgui.InputType_RELEASE)
	InputTypeShort   = InputType(pbgui.InputType_SHORT)
	InputTypeLong    = InputType(pbgui.InputType_LONG)
	InputTypeRepeat  = InputType(pbgui.InputType_REPEAT)
)

func (g *gui) StartScreenStream(callback updateFrame) error {
	req := &pb.Main{
		Content: &pb.Main_GuiStartScreenStreamRequest{},
	}
	g.frameCallback = callback
	_, err := g.f.call(nil, req)
	return err
}

func (g *gui) StopScreenStream() error {
	req := &pb.Main{
		Content: &pb.Main_GuiStopScreenStreamRequest{},
	}
	_, err := g.f.call(nil, req)
	return err
}

func (g *gui) SendInputEvent(key InputKey, eventType InputType) error {
	req := &pb.Main{
		Content: &pb.Main_GuiSendInputEventRequest{
			GuiSendInputEventRequest: &pbgui.SendInputEventRequest{
				Key:  pbgui.InputKey(key),
				Type: pbgui.InputType(eventType),
			},
		},
	}
	_, err := g.f.call(nil, req)
	return err
}

type ScreenStreamFrame struct {
	buffer []byte
}

func (sf ScreenStreamFrame) Bytes() []byte {
	return sf.buffer
}

func (sf ScreenStreamFrame) IsPixelSet(x, y int) bool {
	i := (y / 8) * 128
	y &= 7
	i += x
	return (sf.buffer[i] & (1 << y)) != 0
}

func (sf ScreenStreamFrame) ToImage(foreground, background color.Color) image.Image {
	width := 128
	height := 64

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})

	if background != nil {
		draw.Draw(img, img.Bounds(), &image.Uniform{C: background}, image.Point{}, draw.Src)
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if sf.IsPixelSet(x, y) {
				img.Set(x, y, foreground)
			}
		}
	}

	return img
}
