package main

import ort "github.com/yalue/onnxruntime_go"

var (
	UseCoreML    = false
	Blank        []float32
	ModelPath    = "./yolov8m.onnx"
	Yolo8Model   ModelSession
	ImagePath    = ""
	VideoPath    = "D:\\dialer\\yolov8-go\\IMG_3546.MP4"
	DstImagePath = ""
	DstVideoPath = "D:\\dialer\\yolov8-go\\3.MP4"
)

type ModelSession struct {
	Session *ort.AdvancedSession
	Input   *ort.Tensor[float32]
	Output  *ort.Tensor[float32]
}
