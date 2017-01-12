package formats

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"

	//"image/png"
	// _ "image/gif"

	"path"

	"github.com/k0kubun/pp"
)

const (
	// etl9gRecordSize レコードサイズ
	etl9gRecordSize = 8199

	// etl9gSampleWidth サンプリング画像の幅
	etl9gSampleWidth = 128

	// etl9gSampleHeight サンプリング画像の高さ
	etl9gSampleHeight = 127

	// etl9gSamplePixelNum サンプル画像ピクセル数
	etl9gSamplePixelNum = etl9gSampleWidth * etl9gSampleHeight

	// etl9gSampleSize サンプル画像サイズ
	etl9gSampleSize = etl9gSamplePixelNum / 2

	// etl9gFileNum ファイル数
	etl9gFileNum = 50

	// etl9gRecordNum レコード数
	etl9gRecordNum = 12144
)

// FormatETL9G etl9g用フォーマット
type FormatETL9G struct {
	recordSize   int
	sampleWidth  int
	sampleHeight int
	fileNum      int
	recordNum    int
}

// NewFormatETL9G FormatETL9Gを生成する
func NewFormatETL9G() FormatETL9G {
	return FormatETL9G{
		recordSize:   etl9gRecordSize,
		sampleWidth:  etl9gSampleWidth,
		sampleHeight: etl9gSampleHeight,
		fileNum:      etl9gFileNum,
		recordNum:    etl9gRecordNum,
	}
}

// ReadFiles 指定ディレクトリに存在するすべてのETL9Gファイルを読み込む
// - dir: ETL9Gファイルがあるディレクトリパス
func (f *FormatETL9G) ReadFiles(dir string) error {
	for i := 1; i <= f.fileNum; i++ {
		fpath := path.Join(dir, fmt.Sprintf("ETL9G_%02d", i))
		f.ReadFile(fpath)
	}
	return nil
}

func (f *FormatETL9G) parseRecord(r *BinReader) error {
	sheetIndex, err := r.ReadShort(false)
	if err != nil {
		return err
	}
	pp.Println("sheetIndex:", sheetIndex)

	jisKanjiCode, err := r.ReadShort(false)
	if err != nil {
		return err
	}
	pp.Println("jisKanjiCode:", jisKanjiCode)

	jisTypicalReading, err := r.ReadBytes(8, false)
	if err != nil {
		return err
	}
	pp.Println("jisTypicalReading:", string(jisTypicalReading))

	serialDataNumber, err := r.ReadInt(false)
	if err != nil {
		return err
	}
	pp.Println("serialDataNumber:", serialDataNumber)

	evaluationOfIndividualCharacterImage, err := r.ReadChar(false)
	if err != nil {
		return err
	}
	pp.Println("evaluationOfIndividualCharacterImage:", evaluationOfIndividualCharacterImage)

	evaluationOfCharacterGroup, err := r.ReadChar(false)
	if err != nil {
		return err
	}
	pp.Println("evaluationOfCharacterGroup:", evaluationOfCharacterGroup)

	maleFemaleCode, err := r.ReadChar(false)
	if err != nil {
		return err
	}
	pp.Println("maleFemaleCode:", maleFemaleCode)

	ageOfWriter, err := r.ReadChar(false)
	if err != nil {
		return err
	}
	pp.Println("ageOfWriter:", ageOfWriter)

	industryClassificawtionCode, err := r.ReadShort(false)
	if err != nil {
		return err
	}
	pp.Println("industryClassificawtionCode:", industryClassificawtionCode)

	occupationClassificationCode, err := r.ReadShort(false)
	if err != nil {
		return err
	}
	pp.Println("occupationClassificationCode:", occupationClassificationCode)

	sheetGatherringDate, err := r.ReadShort(false)
	if err != nil {
		return err
	}
	pp.Println("sheetGatherringDate(YYMM):", sheetGatherringDate)

	scanningDate, err := r.ReadShort(false)
	if err != nil {
		return err
	}
	pp.Println("scanningDate(YYMM):", scanningDate)

	samplePositionXOnSheet, err := r.ReadChar(false)
	if err != nil {
		return err
	}
	pp.Println("samplePositionXOnSheet:", samplePositionXOnSheet)

	samplePositionYOnSheet, err := r.ReadChar(false)
	if err != nil {
		return err
	}
	pp.Println("samplePositionYOnSheet:", samplePositionYOnSheet)

	_, err = r.ReadBytes(34, false)
	if err != nil {
		return err
	}

	sampleImageRawReader, err := r.ReNew(etl9gSampleSize, false)
	if err != nil {
		return err
	}

	pixels := []uint8{}

	imB, err := r.ReadBytes(etl9gSampleSize, false)
	if err != nil {
		return err
	}

	for _, twPxT := range imB {
		twPx := int(twPxT)

		// 8bitから4bit取得
		// http://stackoverflow.com/questions/29583024/reading-8-bits-from-a-reader-in-golang
		px1 := twPx >> 4
		px2 := twPx & 0x0F

		// 4bitグレイスケールを8bitへ
		px1 = px1 * 256 / 16
		px2 = px2 * 256 / 16

		pixels = append(pixels, uint8(px1))
		pixels = append(pixels, uint8(px2))
	}

	sampleImage := image.NewGray(image.Rect(0, 0, etl9gSampleWidth, etl9gSampleHeight))
	sampleImage.Pix = pixels
	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, sampleImage, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("./hoge.jpg", buf.Bytes(), 0644)
}

// ReadFile 指定ファイルパスのETL9Gファイルを読み込む
// - fpath: ETL9Gファイルパス
func (f *FormatETL9G) ReadFile(fpath string) error {
	pp.Println(fpath)
	r, err := NewBinReaderFromFilePath(fpath)
	if err != nil {
		return err
	}

	for i := 0; i < etl9gRecordNum; i++ {
		r2, err := r.ReNew(etl9gRecordSize, false)
		if err != nil {
			return err
		}

		err = f.parseRecord(r2)
		if err != nil {
			return err
		}

		break
	}

	return nil
}
