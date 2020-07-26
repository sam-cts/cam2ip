// +build cv4,!cv2

// Package camera.
package camera

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pbnjay/pixfont"
	"gocv.io/x/gocv"
)

// Camera represents camera.
type Camera struct {
	opts    Options
	camera  *gocv.VideoCapture
	frame   *gocv.Mat
	imgchan chan image.Image
}

// New returns new Camera for given camera index.
func New(opts Options) (camera *Camera, err error) {
	camera = &Camera{}
	camera.opts = opts

	mat := gocv.NewMat()
	camera.frame = &mat

	camera.camera, err = gocv.VideoCaptureDevice(opts.Index)
	if err != nil {
		err = fmt.Errorf("camera: can not open camera %d: %s", opts.Index, err.Error())
	}

	camera.SetProperty(PropFrameWidth, opts.Width)
	camera.SetProperty(PropFrameHeight, opts.Height)

	go camera.StartStream()

	return
}

func (c *Camera) GetImage() image.Image {
	return <-c.imgchan
}

func (c *Camera) StartStream() {
	// var img image.Image
	var err error
	var img image.Image
	for {
		img, err = c.Read()
		if err != nil {
			return
		}
		c.imgchan <- img
	}
}

// Read reads next frame from camera and returns image.
func (c *Camera) Read() (img image.Image, err error) {
	ok := c.camera.Read(c.frame)
	if !ok {
		err = fmt.Errorf("camera: can not grab frame")
		return
	}

	img, e := c.frame.ToImage()
	if e != nil {
		err = fmt.Errorf("camera: %v", e)
		return
	}

	if c.frame == nil {
		err = fmt.Errorf("camera: can not retrieve frame")
		return
	}

	switch c.opts.Rotate {
	case 90:
		img = imaging.Rotate90(img)
	case 180:
		img = imaging.Rotate180(img)
	case 270:
		img = imaging.Rotate270(img)
	}

	if c.opts.Timestamp {
		dimg, ok := img.(draw.Image)
		if !ok {
			err = fmt.Errorf("camera: %T is not a drawable image type", img)
			return
		}

		pixfont.DrawString(dimg, 10, 10, time.Now().Format("2006-01-02 15:04:05"), color.White)
		img = dimg
	}

	return
}

// GetProperty returns the specified camera property.
func (c *Camera) GetProperty(id int) float64 {
	return c.camera.Get(gocv.VideoCaptureProperties(id))
}

// SetProperty sets a camera property.
func (c *Camera) SetProperty(id int, value float64) {
	c.camera.Set(gocv.VideoCaptureProperties(id), value)
}

// Close closes camera.
func (c *Camera) Close() (err error) {
	if c.camera == nil {
		err = fmt.Errorf("camera: camera is not opened")
		return
	}

	c.frame.Close()
	err = c.camera.Close()
	c.camera = nil
	return
}
