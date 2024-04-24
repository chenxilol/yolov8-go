package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// DrawTextInfo 图片绘字信息
type DrawTextInfo struct {
	Text string
	X    int
	Y    int
}

// DrawRectInfo 图片画框信息
type DrawRectInfo struct {
	X1          int
	Y1          int
	X2          int
	Y2          int
	ObjectType  string
	Probability float32
}

// TextBrush 字体相关
type TextBrush struct {
	FontType  *truetype.Font
	FontSize  float64
	FontColor *image.Uniform
	TextWidth int
}

// NewTextBrush 新生成笔刷
func NewTextBrush(FontFilePath string, FontSize float64, FontColor *image.Uniform, textWidth int) (*TextBrush, error) {
	fontFile, err := os.ReadFile(FontFilePath)
	if err != nil {
		return nil, err
	}
	fontType, err := truetype.Parse(fontFile)
	if err != nil {
		return nil, err
	}
	if textWidth <= 0 {
		textWidth = 20
	}
	return &TextBrush{FontType: fontType, FontSize: FontSize, FontColor: FontColor, TextWidth: textWidth}, nil
}

// DrawFontOnRGBA 图片插入文字
func (fb *TextBrush) DrawFontOnRGBA(rgba *image.RGBA, pt image.Point, content string) {

	c := freetype.NewContext()
	c.SetDPI(102)
	c.SetFont(fb.FontType)
	c.SetHinting(font.HintingFull)
	c.SetFontSize(fb.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fb.FontColor)

	_, _ = c.DrawString(content, freetype.Pt(pt.X, pt.Y))

}

// Image2RGBA Image2RGBA
func Image2RGBA(img image.Image) *image.RGBA {

	baseSrcBounds := img.Bounds().Max

	newWidth := baseSrcBounds.X
	newHeight := baseSrcBounds.Y

	des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板
	//首先将一个图片信息存入jpg
	draw.Draw(des, des.Bounds(), img, img.Bounds().Min, draw.Over)

	return des
}

// DrawStringOnImageAndSave 图片上写文字
func DrawStringOnImageAndSave(imagePath string, imageData []byte, infos []*DrawTextInfo) (err error) {
	//判断图片类型
	var backgroud image.Image
	filetype := http.DetectContentType(imageData)
	switch filetype {
	case "image/jpeg", "image/jpg":
		backgroud, err = jpeg.Decode(bytes.NewReader(imageData))
		if err != nil {
			fmt.Println("jpeg error")
			return err
		}

	case "image/gif":
		backgroud, err = gif.Decode(bytes.NewReader(imageData))
		if err != nil {
			return err
		}

	case "image/png":
		backgroud, err = png.Decode(bytes.NewReader(imageData))
		if err != nil {
			return err
		}
	default:
		return err
	}
	des := Image2RGBA(backgroud)

	//新建笔刷
	textBrush, _ := NewTextBrush("ttf/arial.ttf", 50, image.Black, 50)

	//Px Py 绘图开始坐标 text要绘制的文字
	//调整颜色
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(textBrush.FontType)
	c.SetHinting(font.HintingFull)
	c.SetFontSize(textBrush.FontSize)
	c.SetClip(des.Bounds())
	c.SetDst(des)
	textBrush.FontColor = image.NewUniform(color.RGBA{
		R: 0xFF,
		G: 0,
		B: 0,
		A: 255,
	})
	c.SetSrc(textBrush.FontColor)

	for _, info := range infos {
		_, _ = c.DrawString(info.Text, freetype.Pt(info.X, info.Y))
	}

	//保存图片
	fSave, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer fSave.Close()

	err = jpeg.Encode(fSave, des, nil)

	if err != nil {
		return err
	}
	return nil
}

// DrawRectOnImageAndSave 图片上画多个框
func DrawRectOnImageAndSave(imagePath string, imageData []byte, infos []*DrawRectInfo) (err error) {
	//判断图片类型
	var backgroud image.Image
	filetype := http.DetectContentType(imageData)
	switch filetype {
	case "image/jpeg", "image/jpg":
		backgroud, err = jpeg.Decode(bytes.NewReader(imageData))
		if err != nil {
			fmt.Println("jpeg error")
			return err
		}

	case "image/gif":
		backgroud, err = gif.Decode(bytes.NewReader(imageData))
		if err != nil {
			return err
		}

	case "image/png":
		backgroud, err = png.Decode(bytes.NewReader(imageData))
		if err != nil {
			return err
		}
	default:
		return err
	}
	des := Image2RGBA(backgroud)
	//新建笔刷
	textBrush, _ := NewTextBrush("D:\\workspace\\wrest-chat\\wclient\\yolov8\\HONORSansDesign-PC.ttf", 20, image.Black, 30)
	for _, info := range infos {
		var c *freetype.Context
		c = freetype.NewContext()
		c.SetDPI(122)
		c.SetFont(textBrush.FontType)
		c.SetHinting(font.HintingFull)
		c.SetFontSize(textBrush.FontSize)
		c.SetClip(des.Bounds())
		c.SetDst(des)
		cGreen := image.NewUniform(color.RGBA{
			R: 0,
			G: 0xFF,
			B: 0,
			A: 255,
		})

		c.SetSrc(cGreen)

		for i := info.X1; i < info.X2; i++ {
			_, _ = c.DrawString("·", freetype.Pt(i, info.Y1))
			_, _ = c.DrawString("·", freetype.Pt(i, info.Y2))
		}
		for j := info.Y1; j < info.Y2; j++ {
			_, _ = c.DrawString("·", freetype.Pt(info.X1, j))
			_, _ = c.DrawString("·", freetype.Pt(info.X2, j))
		}
		textBrush.FontColor = image.NewUniform(color.RGBA{
			R: 0xFF,
			G: 0,
			B: 0,
			A: 255,
		})
		c.SetSrc(textBrush.FontColor)

		for _, info := range infos {
			_, _ = c.DrawString(fmt.Sprintf("类型 %s", info.ObjectType), freetype.Pt(info.X1, info.Y1-70))
			_, _ = c.DrawString(fmt.Sprintf("精度 %.2f", info.Probability), freetype.Pt(info.X1, info.Y1-30))
		}
	}

	//保存图片
	fSave, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer fSave.Close()

	err = jpeg.Encode(fSave, des, nil)

	if err != nil {
		return err
	}
	return nil
}

// DrawRectOnImage 图片上画多个框并返回帧
func DrawRectOnImage(imageData []byte, infos []*DrawRectInfo) (imageBytes []byte, err error) {
	//判断图片类型
	var backgroud image.Image
	filetype := http.DetectContentType(imageData)
	switch filetype {
	case "image/jpeg", "image/jpg":
		backgroud, err = jpeg.Decode(bytes.NewReader(imageData))
		if err != nil {
			fmt.Println("jpeg error")
			return nil, err
		}

	case "image/gif":
		backgroud, err = gif.Decode(bytes.NewReader(imageData))
		if err != nil {
			return nil, err
		}

	case "image/png":
		backgroud, err = png.Decode(bytes.NewReader(imageData))
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	des := Image2RGBA(backgroud)
	//新建笔刷
	textBrush, _ := NewTextBrush("D:\\workspace\\wrest-chat\\wclient\\yolov8\\HONORSansDesign-PC.ttf", 20, image.Black, 30)
	for _, info := range infos {
		var c *freetype.Context
		c = freetype.NewContext()
		c.SetDPI(122)
		c.SetFont(textBrush.FontType)
		c.SetHinting(font.HintingFull)
		c.SetFontSize(textBrush.FontSize)
		c.SetClip(des.Bounds())
		c.SetDst(des)
		cGreen := image.NewUniform(color.RGBA{
			R: 0,
			G: 0xFF,
			B: 0,
			A: 255,
		})

		c.SetSrc(cGreen)

		for i := info.X1; i < info.X2; i++ {
			_, _ = c.DrawString("·", freetype.Pt(i, info.Y1))
			_, _ = c.DrawString("·", freetype.Pt(i, info.Y2))
		}
		for j := info.Y1; j < info.Y2; j++ {
			_, _ = c.DrawString("·", freetype.Pt(info.X1, j))
			_, _ = c.DrawString("·", freetype.Pt(info.X2, j))
		}
		textBrush.FontColor = image.NewUniform(color.RGBA{
			R: 0xFF,
			G: 0,
			B: 0,
			A: 255,
		})
		c.SetSrc(textBrush.FontColor)

		for _, info := range infos {
			_, _ = c.DrawString(fmt.Sprintf("类型 %s", info.ObjectType), freetype.Pt(info.X1, info.Y1-70))
			_, _ = c.DrawString(fmt.Sprintf("精度 %.2f", info.Probability), freetype.Pt(info.X1, info.Y1-30))
		}
	}

	var buf []byte
	// 创建一个缓冲区，并将图像写入其中
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, des)
	if err != nil {
		// 处理错误
	}
	// 将缓冲区中的内容复制到字节切片中
	buf = buffer.Bytes()
	return buf, nil
}
