package apriltag

// #cgo CFLAGS: -std=gnu99 -fPIC -I. -O4 -fno-strict-overflow
// #cgo LDFLAGS: -lm
// #include "apriltag.h"
// #include "tag36h11.h"
// #include "common/image_u8.h"
import "C"

import (
	"image"
	"image/color"
	"unsafe"
)

// Detector encapsulates an apriltag detector initialized with a tag family.
type Detector struct {
	family   *C.apriltag_family_t
	detector *C.apriltag_detector_t
}

// New creates a new 36h11 detector.
func New() *Detector {
	f := C.tag36h11_create()
	d := C.apriltag_detector_create()
	C.apriltag_detector_add_family(d, f)
	return &Detector{
		family:   f,
		detector: d,
	}
}

// Close frees the underlying resources of the detector.
func (d *Detector) Close() error {

	if d.detector != nil {
		C.apriltag_detector_destroy(d.detector)
		d.detector = nil
	}
	if d.family != nil { //TODO: switch when more than tag36h11 supported
		C.tag36h11_destroy(d.family)
		d.family = nil
	}
	return nil
}

func grayToC(img *image.Gray) *C.image_u8_t {
	b := img.Bounds().Size()
	cImg := new(C.image_u8_t)
	cImg.width = C.int32_t(b.X)
	cImg.height = C.int32_t(b.Y)
	cImg.stride = C.int32_t(img.Stride)
	cImg.buf = (*C.uint8_t)(C.CBytes(img.Pix))
	return cImg
}

// Finding represents a single apriltag detected in an image.
type Finding struct {
	// The decoded ID of the tag
	ID int

	// The number of error bits that were corrected.
	Hamming  int
	Goodness float32
	// A measure of the quality of the binary decoding process. Higher == better.
	DecisionMargin float32

	// Centeroid of the detected tag
	CenterX float64
	CenterY float64
	// Bounding edges of the detected tag. [0] = X, [1] = Y.
	Corners [4][2]float64
}

// ImgToGrayscale converts a image.Image into an image.Gray
func ImgToGrayscale(img image.Image) *image.Gray {
	if g, alreadyGray := img.(*image.Gray); alreadyGray {
		return g
	}

	out := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			out.Set(x, y, color.GrayModel.Convert(img.At(x, y)).(color.Gray))
		}
	}
	return out
}

// Find returns a slice of all the apriltags found in the provided image.
func (d *Detector) Find(goImage *image.Gray) []Finding {
	img := grayToC(goImage)
	detections := C.apriltag_detector_detect(d.detector, img)
	numDetections := int(C.zarray_size(detections))

	var findings []Finding
	for i := 0; i < numDetections; i++ {
		var detection *C.apriltag_detection_t
		C.zarray_get(detections, C.int(i), unsafe.Pointer(&detection))

		finding := Finding{
			ID:             int(detection.id),
			CenterX:        float64(detection.c[0]),
			CenterY:        float64(detection.c[1]),
			Hamming:        int(detection.hamming),
			Goodness:       float32(detection.goodness),
			DecisionMargin: float32(detection.decision_margin),
		}
		for i := range detection.p {
			finding.Corners[i][0] = float64(detection.p[i][0])
			finding.Corners[i][1] = float64(detection.p[i][1])
		}

		findings = append(findings, finding)
	}

	C.apriltag_detections_destroy(detections)
	C.free(unsafe.Pointer(img.buf))
	return findings
}

// WritableImage represents a drawing surface.
type WritableImage interface {
	image.Image
	Set(x, y int, c color.Color)
}

func drawLineHorizontal(img WritableImage, startX, y, length int, color color.Color) {
	for i := 0; i < length; i++ {
		img.Set(startX+i, y, color)
	}
}
func drawLineVertical(img WritableImage, x, startY, length int, color color.Color) {
	for i := 0; i < length; i++ {
		img.Set(x, startY+i, color)
	}
}

// DrawFindings draws all found apriltags on the provided image.
func DrawFindings(img WritableImage, findings []Finding, centerColor, cornerColor color.Color) {
	for _, finding := range findings {
		drawLineHorizontal(img, int(finding.CenterX)-7, int(finding.CenterY), 14, centerColor)
		drawLineVertical(img, int(finding.CenterX), int(finding.CenterY)-7, 14, centerColor)

		for _, corner := range finding.Corners {
			drawLineHorizontal(img, int(corner[0])-7, int(corner[1]), 14, cornerColor)
			drawLineVertical(img, int(corner[0]), int(corner[1])-7, 14, cornerColor)
		}
	}
}
