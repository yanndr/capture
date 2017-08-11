package capture

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/opennota/screengen"
	"golang.org/x/net/context"
)

//VideoCaptureService expose funcions for video.
type VideoCaptureService struct {
}

func New() VideoCaptureService {
	return VideoCaptureService{}
}

type ExtractRequest struct {
	Path   string
	Time   int64
	Width  int32
	Height int32
}

type ExtractResponse struct {
	Data []byte
	Err  error
}

type Service interface {
	Extract(context.Context, ExtractRequest) ExtractResponse
}

func NewService(logger log.Logger) Service {
	var svc Service
	svc = VideoCaptureService{}
	svc = LoggingMiddleware(logger)(svc)

	return svc
}

//Extract extract an image from a video.
func (s VideoCaptureService) Extract(ctx context.Context, request ExtractRequest) ExtractResponse {

	g, err := screengen.NewGenerator(request.Path)
	if err != nil {
		return ExtractResponse{nil, err}
	}

	img, err := g.ImageWxH(request.Time, int(request.Width), int(request.Height))
	if err != nil {
		return ExtractResponse{nil, err}
	}

	// var imgResult image.Image
	// if in.OverlayImage != nil {
	// 	imgResult, err = addImageOverlay(img, in.OverlayImage)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	imgResult = img
	// }

	result, err := saveToPng(img)
	if err != nil {
		return ExtractResponse{nil, err}
	}

	return ExtractResponse{result, nil}
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
