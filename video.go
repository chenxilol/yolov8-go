package main

import (
	"bytes"
	"fmt"
	"gocv.io/x/gocv"
	"io"
	"log"
)

// VideoDecoder is a struct to decode video frames.
type VideoDecoder struct {
	videoCapture *gocv.VideoCapture
	width        int
	height       int
	framerate    float64
}

// NewVideoDecoder creates a new VideoDecoder instance.
func NewVideoDecoder(path string) (*VideoDecoder, error) {
	videoCapture, err := gocv.VideoCaptureFile(path) // 传入视频文件路径或者视频流
	if err != nil {
		return nil, fmt.Errorf("failed to open video: %w", err)
	}

	width := int(videoCapture.Get(gocv.VideoCaptureFrameWidth))
	height := int(videoCapture.Get(gocv.VideoCaptureFrameHeight))
	framerate := videoCapture.Get(gocv.VideoCaptureFPS)

	return &VideoDecoder{
		videoCapture: videoCapture,
		width:        width,
		height:       height,
		framerate:    framerate,
	}, nil
}

// NextFrame reads the next frame from the video.
func (vd *VideoDecoder) NextFrame() ([]byte, error) {
	frame := gocv.NewMat()
	defer frame.Close()

	if ok := vd.videoCapture.Read(&frame); !ok {
		if vd.videoCapture.Get(gocv.VideoCapturePosFrames) == vd.videoCapture.Get(gocv.VideoCaptureFrameCount) {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("failed to read frame")
	}

	frameBytes, err := gocv.IMEncode(".jpg", frame)
	if err != nil {
		return nil, fmt.Errorf("failed to encode frame: %w", err)
	}

	return frameBytes.GetBytes(), nil
}

// Width returns the width of the video.
func (vd *VideoDecoder) Width() int {
	return vd.width
}

// Height returns the height of the video.
func (vd *VideoDecoder) Height() int {
	return vd.height
}

// Framerate returns the framerate of the video.
func (vd *VideoDecoder) Framerate() float64 {
	return vd.framerate
}

// Close closes the VideoDecoder.
func (vd *VideoDecoder) Close() {
	vd.videoCapture.Close()
}
func DetectVideo(path, dstPath string) {
	fmt.Println(path)
	// 创建视频解码器
	videoDecoder, err := NewVideoDecoder(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer videoDecoder.Close()

	// 创建视频编码器
	videoEncoder, err := NewVideoEncoder(dstPath, videoDecoder.Width(), videoDecoder.Height(), videoDecoder.Framerate())
	if err != nil {
		log.Println(err)
		return
	}
	defer videoEncoder.Close()

	// 逐帧进行目标检测
	for {
		frame, err := videoDecoder.NextFrame()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return
		}

		// 对当前帧进行目标检测
		boxes, err := Detect_objects_on_image(bytes.NewReader(frame))
		if err != nil {
			log.Println(err)
			return
		}
		drawinfo := make([]*DrawRectInfo, 0)
		for i := 0; i < len(boxes); i++ {
			d := &DrawRectInfo{}
			d.X1 = int(boxes[i][0].(float64))
			d.Y1 = int(boxes[i][1].(float64))
			d.X2 = int(boxes[i][2].(float64))
			d.Y2 = int(boxes[i][3].(float64))
			d.ObjectType = boxes[i][4].(string)
			d.Probability = boxes[i][5].(float32)
			drawinfo = append(drawinfo, d)
		}
		// 绘制边界框,
		imageBytes, err := DrawRectOnImage(frame, drawinfo)
		if err != nil {
			log.Println(err)
			return
		}

		// 将处理后的帧写入视频
		err = videoEncoder.WriteFrame(imageBytes)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

type VideoEncoder struct {
	videoWriter *gocv.VideoWriter
	width       int
	height      int
	framerate   float64
}

// NewVideoEncoder creates a new VideoEncoder instance.
func NewVideoEncoder(path string, width, height int, framerate float64) (*VideoEncoder, error) {
	videoWriter, err := gocv.VideoWriterFile(path, "MP4V", framerate, width, height, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create video writer: %w", err)
	}

	return &VideoEncoder{
		videoWriter: videoWriter,
		width:       width,
		height:      height,
		framerate:   framerate,
	}, nil
}

// WriteFrame writes a frame to the video.
func (ve *VideoEncoder) WriteFrame(frame []byte) error {
	img, err := gocv.IMDecode(frame, gocv.IMReadColor)
	if err != nil {
		return fmt.Errorf("failed to decode frame: %w", err)
	}
	defer img.Close()

	if !ve.videoWriter.IsOpened() {
		return fmt.Errorf("video writer is not opened")
	}

	if err := ve.videoWriter.Write(img); err != nil {
		return fmt.Errorf("failed to write frame")
	}

	return nil
}

// Close closes the VideoEncoder.
func (ve *VideoEncoder) Close() {
	ve.videoWriter.Close()
}
