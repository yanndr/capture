package capture

import (
	"bufio"
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/opennota/screengen"
	"golang.org/x/net/context"
)

type ExtractRequest struct {
	Video  []byte
	Name   string
	Time   int64
	Width  int32
	Height int32
}

type ExtractResponse struct {
	Data []byte
}

type Service interface {
	Extract(context.Context, ExtractRequest) (ExtractResponse, error)
}

func NewService(logger log.Logger, extracts metrics.Counter) Service {
	var svc Service
	svc = VideoCaptureService{}
	svc = LoggingMiddleware(logger)(svc)
	svc = InstrumentingMiddleware(extracts)(svc)

	return svc
}

//VideoCaptureService expose funcions for video.
type VideoCaptureService struct {
}

func New() VideoCaptureService {
	return VideoCaptureService{}
}

//Extract extract an image from a video.
func (s VideoCaptureService) Extract(ctx context.Context, request ExtractRequest) (ExtractResponse, error) {

	name, err := saveFile(request.Name, request.Video)
	if err != nil {
		return ExtractResponse{}, err
	}

	g, err := screengen.NewGenerator(name)
	if err != nil {
		return ExtractResponse{}, err
	}

	img, err := g.ImageWxH(request.Time, int(request.Width), int(request.Height))
	if err != nil {
		return ExtractResponse{}, err
	}

	// var imgResult image.Image
	// if request.OverlayImage != nil {
	// 	imgResult, err = addImageOverlay(img, in.OverlayImage)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	imgResult = img
	// }

	result, err := saveToPng(img)
	if err != nil {
		return ExtractResponse{}, err
	}

	os.Remove(name)

	return ExtractResponse{Data: result}, nil
}

func saveFile(name string, data []byte) (string, error) {
	fo, err := os.Create(name)
	if err != nil {
		return "", err
	}

	w := bufio.NewWriter(fo)

	if _, err := w.Write(data); err != nil {
		return "", err
	}

	if err := w.Flush(); err != nil {
		return "", err
	}

	return name, nil
}

type overlayImage interface {
	GetPath() string
	GetX() int32
	GetY() int32
}

func addImageOverlay(img image.Image, overlayImage overlayImage) (*image.RGBA, error) {
	flogo, err := os.Open(overlayImage.GetPath())
	if err != nil {
		return nil, err
	}

	defer flogo.Close()
	logo, _, err := image.Decode(flogo)
	if err != nil {
		return nil, err
	}

	m := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X, img.Bounds().Max.Y))
	draw.Draw(m, m.Bounds(), img, image.Point{0, 0}, draw.Src)
	draw.Draw(m, m.Bounds(), logo, image.Point{int(overlayImage.GetX()), int(overlayImage.GetY())}, draw.Src)

	return m, nil
}

func saveToPng(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	// log.Println("Wrote ", len(buf.Bytes()))
	return buf.Bytes(), nil
}
