package formats

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"path"

	"github.com/PyYoshi/etlcdb-tools/utils"
)

type ETLFormat string

const (
	ETLFormat1  ETLFormat = "1"
	ETLFormat2  ETLFormat = "2"
	ETLFormat3  ETLFormat = "3"
	ETLFormat4  ETLFormat = "4"
	ETLFormat5  ETLFormat = "5"
	ETLFormat6  ETLFormat = "6"
	ETLFormat7  ETLFormat = "7"
	ETLFormat8b ETLFormat = "8b"
	ETLFormat8g ETLFormat = "8g"
	ETLFormat9b ETLFormat = "9b"
	ETLFormat9g ETLFormat = "9g"
)

type Record interface {
	OutputImage(outputDir string) error
	DeallocImage()
}

// outputPng レコードに格納された画像をPNG形式で任意のディレクトリへ出力する
// - outputDir: 出力するディレクトリパス
// - imageName: 画像ファイル名
// - img: 画像データ
// - compLevel: 圧縮レベル
func outputPng(outputDir, imageName string, img image.Image, compLevel png.CompressionLevel) error {
	buf := bytes.Buffer{}

	enc := &png.Encoder{CompressionLevel: compLevel}
	err := enc.Encode(&buf, img)
	if err != nil {
		return err
	}

	err = utils.CreateIfNotExists(outputDir, true)
	if err != nil {
		return err
	}

	fpath := path.Join(outputDir, imageName)
	return ioutil.WriteFile(fpath, buf.Bytes(), 0644)
}
