package formats

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/PyYoshi/etlcdb-tools/tables"
	"github.com/PyYoshi/etlcdb-tools/utils"
	"github.com/disintegration/imaging"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	// etl8gRecordSize レコードサイズ
	etl8gRecordSize = 8199

	// etl8gSampleWidth サンプリング画像の幅
	etl8gSampleWidth = 128

	// etl8gSampleHeight サンプリング画像の高さ
	etl8gSampleHeight = 127

	// etl8gSamplePixelNum サンプル画像ピクセル数
	etl8gSamplePixelNum = etl8gSampleWidth * etl8gSampleHeight

	// etl8gSampleSize サンプル画像サイズ
	etl8gSampleSize = etl8gSamplePixelNum / 2

	// etl8gFileNum ファイル数
	etl8gFileNum = 33

	// etl8gRecordNum レコード数
	etl8gRecordNum = 4780

	// etl8gRecordTotalNum レコード総数
	etl8gRecordTotalNum = (etl8gRecordNum * (etl8gFileNum - 1)) + 956
)

// RecordETL8G ETL8G用レコード
// http://etlcdb.db.aist.go.jp/?page_id=2461
type RecordETL8G struct {
	Format      ETLFormat   `json:"format"`
	Character   string      `json:"character"`
	Image       image.Image `json:"-"`
	ImageHash   string      `json:"-"`
	ImageName   string      `json:"image_name"`
	ImageWidth  int         `json:"image_width"`
	ImageHeight int         `json:"image_height"`

	SerialSheetNumber                           uint16 `json:"serial_sheet_number"`
	JisCharacterCode                            uint16 `json:"jis_character_code"`
	JisTypicalReading                           string `json:"jis_typical_reading"`
	SerialDataNumber                            uint32 `json:"serial_data_number"`
	QualityEvaluationOfIndividualCharacterImage uint8  `json:"quality_evaluation_of_individual_character_image"`
	QualityEvaluationOfCharacterGroup           uint8  `json:"quality_evaluation_of_character_group"`
	GenderOfWriter                              uint8  `json:"gender_of_writer"`
	AgeOfWriter                                 uint8  `json:"age_of_writer"`
	IndustryClassificationCode                  uint16 `json:"industry_classification_code"`
	OccupationClassificationCode                uint16 `json:"occupation_classification_code"`
	DateOfCollection                            uint16 `json:"date_of_collection"`
	DateOfScan                                  uint16 `json:"date_of_scan"`
	XCoordinateOfSampleOnSheet                  uint8  `json:"x_coordinate_of_sample_on_sheet"`
	YCoordinateOfSampleOnSheet                  uint8  `json:"y_coordinate_of_sample_on_sheet"`
}

// DeallocImage RecordETL8G.Imageにnilを代入する
func (r *RecordETL8G) DeallocImage() {
	r.Image = nil
}

// NewRecordETL8G RecordETL8Gを生成する
func NewRecordETL8G(
	serialSheetNumber uint16,
	jisCharacterCode uint16,
	jisTypicalReading string,
	serialDataNumber uint32,
	qualityEvaluationOfIndividualCharacterImage uint8,
	qualityEvaluationOfCharacterGroup uint8,
	genderOfWriter uint8,
	ageOfWriter uint8,
	industryClassificationCode uint16,
	occupationClassificationCode uint16,
	dateOfCollection uint16,
	dateOfScan uint16,
	xCoordinateOfSampleOnSheet uint8,
	yCoordinateOfSampleOnSheet uint8,
	img image.Image,
	imgHash string,
) RecordETL8G {
	return RecordETL8G{
		Format:      ETLFormat8g,
		Character:   string(tables.JIS0208[jisCharacterCode]),
		Image:       img,
		ImageHash:   imgHash,
		ImageName:   fmt.Sprintf("ETL8G_0x%x_%s.png", jisCharacterCode, imgHash),
		ImageWidth:  etl8gSampleWidth,
		ImageHeight: etl8gSampleHeight,

		SerialSheetNumber:                           serialSheetNumber,
		JisCharacterCode:                            jisCharacterCode,
		JisTypicalReading:                           jisTypicalReading,
		SerialDataNumber:                            serialDataNumber,
		QualityEvaluationOfIndividualCharacterImage: qualityEvaluationOfIndividualCharacterImage,
		QualityEvaluationOfCharacterGroup:           qualityEvaluationOfCharacterGroup,
		GenderOfWriter:                              genderOfWriter,
		AgeOfWriter:                                 ageOfWriter,
		IndustryClassificationCode:                  industryClassificationCode,
		OccupationClassificationCode:                occupationClassificationCode,
		DateOfCollection:                            dateOfCollection,
		DateOfScan:                                  dateOfScan,
		XCoordinateOfSampleOnSheet:                  xCoordinateOfSampleOnSheet,
		YCoordinateOfSampleOnSheet:                  yCoordinateOfSampleOnSheet,
	}
}

// OutputImage レコードに格納された画像任意のディレクトリへ出力する
// - outputDir: 出力するディレクトリパス
// - width: 出力する画像の幅
// - height: 出力する画像の高さ
func (r *RecordETL8G) OutputImage(outputDir string, width, height int) error {
	if r.Image == nil {
		return errors.New("RecordETL8G.Image is nil")
	}

	// リサイズ
	var dstImage image.Image
	if !(width == etl8gSampleWidth && height == etl8gSampleHeight) {
		dstImage = imaging.Resize(r.Image, width, height, imaging.Lanczos)
	} else {
		dstImage = r.Image
	}

	return outputPng(outputDir, r.ImageName, dstImage, png.BestCompression)
}

// GetKey ETL8Gレコード全体でユニークなキー
func (r *RecordETL8G) GetKey() string {
	return r.ImageName
}

func ParseETL8GRecord(fp io.Reader) (Record, error) {
	var err error

	var serialSheetNumber uint16
	err = binary.Read(fp, binary.BigEndian, &serialSheetNumber)
	if err != nil {
		return nil, err
	}

	var jisCharacterCode uint16
	err = binary.Read(fp, binary.BigEndian, &jisCharacterCode)
	if err != nil {
		return nil, err
	}

	jisTypicalReading := make([]byte, 8)
	err = binary.Read(fp, binary.BigEndian, &jisTypicalReading)
	if err != nil {
		return nil, err
	}

	var serialDataNumber uint32
	err = binary.Read(fp, binary.BigEndian, &serialDataNumber)
	if err != nil {
		return nil, err
	}

	var qualityEvaluationOfIndividualCharacterImage uint8
	err = binary.Read(fp, binary.BigEndian, &qualityEvaluationOfIndividualCharacterImage)
	if err != nil {
		return nil, err
	}

	var qualityEvaluationOfCharacterGroup uint8
	err = binary.Read(fp, binary.BigEndian, &qualityEvaluationOfCharacterGroup)
	if err != nil {
		return nil, err
	}

	var genderOfWriter uint8
	err = binary.Read(fp, binary.BigEndian, &genderOfWriter)
	if err != nil {
		return nil, err
	}

	var ageOfWriter uint8
	err = binary.Read(fp, binary.BigEndian, &ageOfWriter)
	if err != nil {
		return nil, err
	}

	var industryClassificationCode uint16
	err = binary.Read(fp, binary.BigEndian, &industryClassificationCode)
	if err != nil {
		return nil, err
	}

	var occupationClassificationCode uint16
	err = binary.Read(fp, binary.BigEndian, &occupationClassificationCode)
	if err != nil {
		return nil, err
	}

	var dateOfCollection uint16
	err = binary.Read(fp, binary.BigEndian, &dateOfCollection)
	if err != nil {
		return nil, err
	}

	var dateOfScan uint16
	err = binary.Read(fp, binary.BigEndian, &dateOfScan)
	if err != nil {
		return nil, err
	}

	var xCoordinateOfSampleOnSheet uint8
	err = binary.Read(fp, binary.BigEndian, &xCoordinateOfSampleOnSheet)
	if err != nil {
		return nil, err
	}

	var yCoordinateOfSampleOnSheet uint8
	err = binary.Read(fp, binary.BigEndian, &yCoordinateOfSampleOnSheet)
	if err != nil {
		return nil, err
	}

	// undefined
	undefined1 := make([]byte, 30)
	err = binary.Read(fp, binary.BigEndian, &undefined1)
	if err != nil {
		return nil, err
	}

	sampleImage := image.NewGray(image.Rect(0, 0, etl8gSampleWidth, etl8gSampleHeight))
	pxIndex := 0
	for i := 0; i < etl8gSampleSize; i++ {
		var pxT uint8
		err = binary.Read(fp, binary.BigEndian, &pxT)
		if err != nil {
			return nil, err
		}

		// 8bitから4bit取得
		// http://stackoverflow.com/questions/29583024/reading-8-bits-from-a-reader-in-golang
		px1 := pxT >> 4
		px2 := pxT & 0x0F

		// 4bitグレイスケールを8bitへ
		px1 = px1 * (256 / 16)
		px2 = px2 * (256 / 16)

		sampleImage.Pix[pxIndex] = px1
		pxIndex++
		sampleImage.Pix[pxIndex] = px2
		pxIndex++
	}

	// uncertain
	uncertain1 := make([]byte, 11)
	err = binary.Read(fp, binary.BigEndian, &uncertain1)
	if err != nil {
		return nil, err
	}

	imgHash := sha256.New()
	imgHash.Write(sampleImage.Pix)

	record := NewRecordETL8G(
		serialSheetNumber,
		jisCharacterCode,
		strings.TrimSpace(string(jisTypicalReading)),
		serialDataNumber,
		qualityEvaluationOfIndividualCharacterImage,
		qualityEvaluationOfCharacterGroup,
		genderOfWriter,
		ageOfWriter,
		industryClassificationCode,
		occupationClassificationCode,
		dateOfCollection,
		dateOfScan,
		xCoordinateOfSampleOnSheet,
		yCoordinateOfSampleOnSheet,
		sampleImage,
		hex.EncodeToString(imgHash.Sum(nil)),
	)

	return &record, nil
}

type jobWorkerMakeETL8GDatasets struct {
	outputDir                           string
	outputImageWidth, outputImageHeight int
	ldb                                 *leveldb.DB
	mu                                  *sync.Mutex
}

func (w *jobWorkerMakeETL8GDatasets) start(wg *sync.WaitGroup, q chan string) {
	defer wg.Done()
	for {
		fpath, ok := <-q // closeされると ok が false になる
		if !ok {
			return
		}

		log.Printf("ETL8G: reading %s\n", fpath)

		// ファイルを開く
		r, err := os.Open(fpath)
		if err != nil {
			log.Fatal(err)
		}

		ldbBatch := new(leveldb.Batch)
		for i := 0; i < etl8gRecordNum; i++ {
			// レコードサイズ分メモリへ読み込み, そこから処理を行う
			rb := make([]byte, etl8gRecordSize)
			_, err = r.Read(rb)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			}
			r2 := bytes.NewReader(rb)

			record, err := ParseETL8GRecord(r2)
			if err != nil {
				log.Fatal(err)
			}

			// 画像を生成
			err = record.OutputImage(w.outputDir, w.outputImageWidth, w.outputImageHeight)
			if err != nil {
				log.Fatal(err)
			}

			// DeallocImageを逐一呼び出ししないとメモリ不足で落ちる
			record.DeallocImage()

			rjb, err := json.Marshal(record)
			if err != nil {
				log.Fatal(err)
			}
			ldbBatch.Put([]byte(record.GetKey()), rjb)
		}
		r.Close()

		w.mu.Lock()
		err = w.ldb.Write(ldbBatch, nil)
		if err != nil {
			log.Fatal(err)
		}
		w.mu.Unlock()
	}
}

// MakeETL8GDatasets 指定ディレクトリに存在するすべてのETL8Gファイルからデータセットを作成する
// - inputDir: ETL8Gファイルがあるディレクトリパス
// - outputDir: ETL8Gのデータセットを出力するディレクトリパス
// - outputImageWidth: 出力する画像の幅
// - outputImageHeight: 出力する画像の高さ
// - workerNum: 並行して実行する数
func MakeETL8GDatasets(inputDir, outputDir string, outputImageWidth, outputImageHeight, workerNum int) error {
	err := utils.CreateIfNotExists(outputDir, true)
	if err != nil {
		return err
	}

	ldbPath := path.Join(outputDir, ".ldb")
	err = utils.CreateIfNotExists(ldbPath, true)
	if err != nil {
		return err
	}

	ldb, err := leveldb.OpenFile(ldbPath, nil)
	if err != nil {
		return err
	}

	jobWorker := jobWorkerMakeETL8GDatasets{
		outputDir:         outputDir,
		outputImageWidth:  outputImageWidth,
		outputImageHeight: outputImageHeight,
		ldb:               ldb,
		mu:                &sync.Mutex{},
	}

	q := make(chan string, etl8gFileNum)

	wg := &sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go jobWorker.start(wg, q)
	}

	for i := 1; i <= etl8gFileNum; i++ {
		q <- path.Join(inputDir, fmt.Sprintf("ETL8G_%02d", i))
	}
	close(q)

	// 処理待ち
	wg.Wait()

	// etl8g.json用io.Writer
	wj, err := os.Create(path.Join(outputDir, "etl8g.json"))
	if err != nil {
		return err
	}

	// etl8g.json先頭に`[`を付加
	_, err = wj.Write([]byte("[\n"))
	if err != nil {
		return err
	}

	ldbIter := ldb.NewIterator(nil, nil)
	ldbIterIndex := 0
	for ldbIter.Next() {
		rjb := ldbIter.Value()
		_, err = wj.Write(rjb)
		if err != nil {
			return err
		}

		// 最後のレコードの場合はカンマを付けないように
		if ldbIterIndex < etl8gRecordTotalNum-1 {
			_, err = wj.Write([]byte(",\n"))
			if err != nil {
				return err
			}
		} else {
			_, err = wj.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}
		ldbIterIndex++
	}
	ldbIter.Release()
	err = ldbIter.Error()
	if err != nil {
		return err
	}

	// etl8g.json終端に`[`を付加
	_, err = wj.Write([]byte("]"))
	if err != nil {
		return err
	}

	// etl8g.json用io.Writerを閉じる
	wj.Close()

	// leveldbで利用したファイルを削除
	ldb.Close()
	return os.RemoveAll(ldbPath)
}
