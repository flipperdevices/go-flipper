package flipper

import (
	"image"
	"image/color"
	"image/draw"

	pb "github.com/flipperdevices/go-flipper/internal/proto"
	pbgui "github.com/flipperdevices/go-flipper/internal/proto/gui"
)

type gui struct {
	f             *Flipper
	frameCallback updateFrame
}

type updateFrame func(frame ScreenFrame)

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

func (g *gui) StartVirtualDisplay(buf []byte) error {
	var frame *pbgui.ScreenFrame
	if buf != nil {
		frame = &pbgui.ScreenFrame{Data: buf}
	}

	req := &pb.Main{
		Content: &pb.Main_GuiStartVirtualDisplayRequest{
			GuiStartVirtualDisplayRequest: &pbgui.StartVirtualDisplayRequest{
				FirstFrame: frame,
			},
		},
	}
	_, err := g.f.call(nil, req)
	return err
}

func (g *gui) StopVirtualDisplay() error {
	req := &pb.Main{
		Content: &pb.Main_GuiStopVirtualDisplayRequest{},
	}
	_, err := g.f.call(nil, req)
	return err
}

func (g *gui) UpdateVirtualDisplay(buf []byte) error {
	req := &pb.Main{
		Content: &pb.Main_GuiScreenFrame{
			GuiScreenFrame: &pbgui.ScreenFrame{Data: buf},
		},
	}
	return g.f.sendUnsolicited(req)
}

type ScreenFrame struct {
	buffer []byte
}

func (sf ScreenFrame) Bytes() []byte {
	return sf.buffer
}

func (sf ScreenFrame) IsPixelSet(x, y int) bool {
	i := (y / 8) * 128
	y &= 7
	i += x
	return (sf.buffer[i] & (1 << y)) != 0
}

func (sf ScreenFrame) ToImage(foreground, background color.Color) image.Image {
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
