# YOLOv8 inference using Go
本项目是 ：`github.com:AndreyGermanov/yolov8_onnx_go.git`的分支，在此基础上增加了视频检测功能,只提供了api没有实现前端页面。
## 环境
* 在运行之前请确认环境变量是否开启了cgo,windows 命令为`go env -w CGO_ENABLED=0`。如果之前没有开启过cgo，你需要安装MinGW。 安装链接为： <https://sourceforge.net/projects/mingw-w64/files/>
选择posix版本。解压后把bin目录加入到环境变量中。图片的存储路径请在global中修改
* 测试代码： 
```go'
package main

/*
#include <stdio.h>

static void SayHello(const char* s) {
puts(s);
}
*/
import "C"

func main() {
C.SayHello(C.CString("Hello, World\n"))
}
```
* 安装gocv库，对视频的解码和编码需要用到gocv库。安装方法为：`go get -u -d gocv.io/x/gocv`。安装过程中可能会出现一些问题，可以参考官方文档：<https://gocv.io/getting-started/windows/>，如果ip是国内，请把 `raw.githubusercontent.com`添加到自己host中。
  在C盘搜索hosts，然后添加下面的内容：185.199.108.133 raw.githubusercontent.com   #comments. put the address here

  
web 页面 [YOLOv8 目标检测神经网络](https://ultralytics.com/yolov8)
implemented on [Go](https://go.dev).

源代码 ["How to create YOLOv8-based object detection web service using Python, Julia, Node.js, JavaScript, Go and Rust"](https://dev.to/andreygermanov/how-to-create-yolov8-based-object-detection-web-service-using-python-julia-nodejs-javascript-go-and-rust-4o8e) tutorial.


## Run

Execute:

```
go mod tidy
go run .
```

它将在 http://localhost:8080 上启动一个网络服务器