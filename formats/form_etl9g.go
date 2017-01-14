package formats

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/PyYoshi/etlcdb-tools/tables"
	"github.com/PyYoshi/etlcdb-tools/utils"
	"github.com/k0kubun/pp"
	"github.com/syndtr/goleveldb/leveldb"

	"encoding/json"

	"github.com/disintegration/imaging"
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

// RecordETL9G ETL9G用レコード
// http://etlcdb.db.aist.go.jp/?page_id=1711
type RecordETL9G struct {
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

// DeallocImage RecordETL9G.Imageにnilを代入する
func (r *RecordETL9G) DeallocImage() {
	r.Image = nil
}

// NewRecordETL9G RecordETL9Gを生成する
func NewRecordETL9G(
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
) RecordETL9G {
	return RecordETL9G{
		Format:      ETLFormat9g,
		Character:   string(tables.JIS0208[jisCharacterCode]),
		Image:       img,
		ImageHash:   imgHash,
		ImageName:   fmt.Sprintf("ETL9G_0x%x_%s.png", jisCharacterCode, imgHash),
		ImageWidth:  etl9gSampleWidth,
		ImageHeight: etl9gSampleHeight,

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
func (r *RecordETL9G) OutputImage(outputDir string, width, height int) error {
	if r.Image == nil {
		return errors.New("RecordETL9G.Image is nil")
	}

	// リサイズ
	var dstImage image.Image
	if !(width == etl9gSampleWidth && height == etl9gSampleHeight) {
		dstImage = imaging.Resize(r.Image, width, height, imaging.Lanczos)
	} else {
		dstImage = r.Image
	}

	return outputPng(outputDir, r.ImageName, dstImage, png.BestCompression)
}

// GetKey ETL9Gレコード全体でユニークなキー
func (r *RecordETL9G) GetKey() string {
	return r.ImageName
}

func parseETL9GRecord(fp io.Reader) (Record, error) {
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
	undefined1 := make([]byte, 34)
	err = binary.Read(fp, binary.BigEndian, &undefined1)
	if err != nil {
		return nil, err
	}

	sampleImage := image.NewGray(image.Rect(0, 0, etl9gSampleWidth, etl9gSampleHeight))
	pxIndex := 0
	for i := 0; i < etl9gSampleSize; i++ {
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
	uncertain1 := make([]byte, 7)
	err = binary.Read(fp, binary.BigEndian, &uncertain1)
	if err != nil {
		return nil, err
	}

	imgHash := sha256.New()
	imgHash.Write(sampleImage.Pix)

	record := NewRecordETL9G(
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

// ReadETL9GFile 指定ファイルパスのETL9Gファイルを読み込む
// - fpath: ETL9Gファイルパス
func ReadETL9GFile(fpath string) ([]Record, error) {
	// 一度ファイル全体をメモリへ載せることによってファイルへ逐一アクセスせずに済むので処理が早くなるがメモリは余計に使う
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	records := make([]Record, etl9gRecordNum)
	for i := 0; i < etl9gRecordNum; i++ {
		record, err := parseETL9GRecord(r)
		if err != nil {
			return nil, err
		}
		records[i] = record
	}

	return records, nil
}

type jobWorkerMakeETL9GDatasets struct {
	outputDir                           string
	outputImageWidth, outputImageHeight int
	ldb                                 *leveldb.DB
	mu                                  *sync.Mutex
}

func (w *jobWorkerMakeETL9GDatasets) start(wg *sync.WaitGroup, q chan string) {
	defer wg.Done()
	for {
		fpath, ok := <-q // closeされると ok が false になる
		if !ok {
			return
		}

		log.Printf("ETL9G: reading %s\n", fpath)

		// 一度ファイル全体をメモリへ載せることによってファイルへ逐一アクセスせずに済むので処理が早くなるがメモリは余計に使う
		b, err := ioutil.ReadFile(fpath)
		if err != nil {
			log.Fatal(err)
		}
		r := bytes.NewReader(b)

		ldbBatch := new(leveldb.Batch)
		for i := 0; i < etl9gRecordNum; i++ {
			record, err := parseETL9GRecord(r)
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

		w.mu.Lock()
		err = w.ldb.Write(ldbBatch, nil)
		if err != nil {
			log.Fatal(err)
		}
		w.mu.Unlock()
	}
}

// MakeETL9GDatasets 指定ディレクトリに存在するすべてのETL9Gファイルからデータセットを作成する
// - inputDir: ETL9Gファイルがあるディレクトリパス
// - outputDir: ETL9Gのデータセットを出力するディレクトリパス
// - outputImageWidth: 出力する画像の幅
// - outputImageHeight: 出力する画像の高さ
// - workerNum: 並行して実行する数
func MakeETL9GDatasets(inputDir, outputDir string, outputImageWidth, outputImageHeight, workerNum int) error {
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
	defer ldb.Close()

	jobWorker := jobWorkerMakeETL9GDatasets{
		outputDir:         outputDir,
		outputImageWidth:  outputImageWidth,
		outputImageHeight: outputImageHeight,
		ldb:               ldb,
		mu:                &sync.Mutex{},
	}

	q := make(chan string, etl9gFileNum)

	wg := &sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go jobWorker.start(wg, q)
	}

	for i := 1; i <= etl9gFileNum; i++ {
		q <- path.Join(inputDir, fmt.Sprintf("ETL9G_%02d", i))
	}
	close(q)

	// 処理待ち
	wg.Wait()

	ldbIter := ldb.NewIterator(nil, nil)
	var records []Record
	for ldbIter.Next() {
		rjb := ldbIter.Value()
		record := &RecordETL9G{}
		err = json.Unmarshal(rjb, record)
		if err != nil {
			return err
		}
		records = append(records, record)
	}

	pp.Println(len(records))
	bj, err := json.MarshalIndent(records, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(outputDir, "etl9g.json"), bj, 0644)
	if err != nil {
		return err
	}

	return os.RemoveAll(ldbPath)
}
