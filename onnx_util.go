package main

import (
	ort "github.com/yalue/onnxruntime_go"
	"runtime"
)

func InitYolo8Session(input []float32) (ModelSession, error) {
	ort.SetSharedLibraryPath(getSharedLibPath())
	err := ort.InitializeEnvironment()
	if err != nil {
		return ModelSession{}, err
	}

	inputShape := ort.NewShape(1, 3, 640, 640)
	inputTensor, err := ort.NewTensor(inputShape, input)
	if err != nil {
		return ModelSession{}, err
	}

	outputShape := ort.NewShape(1, 84, 8400)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return ModelSession{}, err
	}

	options, e := ort.NewSessionOptions()
	if e != nil {
		return ModelSession{}, err
	}

	if UseCoreML { // If CoreML is enabled, append the CoreML execution provider
		e = options.AppendExecutionProviderCoreML(0)
		if e != nil {
			options.Destroy()
			return ModelSession{}, err
		}
		defer options.Destroy()
	}

	session, err := ort.NewAdvancedSession(ModelPath,
		[]string{"images"}, []string{"output0"},
		[]ort.ArbitraryTensor{inputTensor}, []ort.ArbitraryTensor{outputTensor}, options)

	if err != nil {
		return ModelSession{}, err
	}

	modelSes := ModelSession{
		Session: session,
		Input:   inputTensor,
		Output:  outputTensor,
	}

	return modelSes, err
}

func getSharedLibPath() string {
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			return "./third_party/onnxruntime.dll"
		}
	}
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return "./third_party/onnxruntime_arm64.dylib"
		}
	}
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "arm64" {
			return "../third_party/onnxruntime_arm64.so"
		}
		return "./third_party/onnxruntime.so"
	}
	panic("Unable to find a version of the onnxruntime library supporting this system.")
}

// Array of YOLOv8 class labels
var yolo_classes = []string{
	"人", "自行车", "汽车", "摩托车", "飞机", "公共汽车", "火车", "卡车", "船",
	"交通灯", "消防栓", "停车标志", "停车收费表", "长凳", "鸟", "猫", "狗", "马",
	"羊", "牛", "大象", "熊", "斑马", "长颈鹿", "背包", "雨伞", "手提包", "领带",
	"手提箱", "飞盘", "滑雪板", "滑雪板", "运动球", "风筝", "棒球棒", "棒球手套",
	"滑板", "冲浪板", "网球拍", "瓶子", "酒杯", "杯子", "叉子", "刀", "勺子",
	"碗", "香蕉", "苹果", "三明治", "橙子", "西兰花", "胡萝卜", "热狗", "披萨", "甜甜圈",
	"蛋糕", "椅子", "沙发", "盆栽植物", "床", "餐桌", "厕所", "电视", "笔记本电脑", "鼠标",
	"遥控器", "键盘", "手机", "微波炉", "烤箱", "烤面包机", "水槽", "冰箱", "书",
	"时钟", "花瓶", "剪刀", "泰迪熊", "吹风机", "牙刷",
}

func runInference(modelSes ModelSession, input []float32) ([]float32, error) {
	inTensor := modelSes.Input.GetData()
	copy(inTensor, input)
	err := modelSes.Session.Run()
	if err != nil {
		return nil, err
	}
	return modelSes.Output.GetData(), nil
}
