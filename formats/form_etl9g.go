package formats

import (
	"fmt"
	"image"

	"strings"

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

// RecordETL9G ETL9G用レコード
// http://etlcdb.db.aist.go.jp/?page_id=1711
type RecordETL9G struct {
	Format                                      ETLFormat   `json:"format"`
	SerialSheetNumber                           uint16      `json:"serial_sheet_number"`
	JisCharacterCode                            uint16      `json:"jis_character_code"`
	JisTypicalReading                           string      `json:"jis_typical_reading"`
	SerialDataNumber                            uint32      `json:"serial_data_number"`
	QualityEvaluationOfIndividualCharacterImage uint8       `json:"quality_evaluation_of_individual_character_image"`
	QualityEvaluationOfCharacterGroup           uint8       `json:"quality_evaluation_of_character_group"`
	GenderOfWriter                              uint8       `json:"gender_of_writer"`
	AgeOfWriter                                 uint8       `json:"age_of_writer"`
	IndustryClassificationCode                  uint16      `json:"industry_classification_code"`
	OccupationClassificationCode                uint16      `json:"occupation_classification_code"`
	DateOfCollection                            uint16      `json:"date_of_collection"`
	DateOfScan                                  uint16      `json:"date_of_scan"`
	XCoordinateOfSampleOnSheet                  uint8       `json:"x_coordinate_of_sample_on_sheet"`
	YCoordinateOfSampleOnSheet                  uint8       `json:"y_coordinate_of_sample_on_sheet"`
	Image                                       image.Image `json:"-"`
	ImageName                                   string      `json:"image_name"`
	ImageWidth                                  int         `json:"image_width"`
	ImageHeight                                 int         `json:"image_height"`
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
) RecordETL9G {
	return RecordETL9G{
		Format:                                      ETLFormat9g,
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
		Image:       img,
		ImageName:   fmt.Sprintf("ETL9G_%d_%x.png", serialSheetNumber, jisCharacterCode),
		ImageWidth:  etl9gSampleWidth,
		ImageHeight: etl9gSampleHeight,
	}
}

func parseETL9GRecord(r *BinReader) (Record, error) {
	serialSheetNumber, err := r.ReadUshort(false)
	if err != nil {
		return nil, err
	}

	jisCharacterCode, err := r.ReadUshort(false)
	if err != nil {
		return nil, err
	}

	jisTypicalReading, err := r.ReadBytes(8, false)
	if err != nil {
		return nil, err
	}

	serialDataNumber, err := r.ReadUint(false)
	if err != nil {
		return nil, err
	}

	qualityEvaluationOfIndividualCharacterImage, err := r.ReadUchar(false)
	if err != nil {
		return nil, err
	}

	qualityEvaluationOfCharacterGroup, err := r.ReadUchar(false)
	if err != nil {
		return nil, err
	}

	genderOfWriter, err := r.ReadUchar(false)
	if err != nil {
		return nil, err
	}

	ageOfWriter, err := r.ReadUchar(false)
	if err != nil {
		return nil, err
	}

	industryClassificationCode, err := r.ReadUshort(false)
	if err != nil {
		return nil, err
	}

	occupationClassificationCode, err := r.ReadUshort(false)
	if err != nil {
		return nil, err
	}

	dateOfCollection, err := r.ReadUshort(false)
	if err != nil {
		return nil, err
	}

	dateOfScan, err := r.ReadUshort(false)
	if err != nil {
		return nil, err
	}

	xCoordinateOfSampleOnSheet, err := r.ReadUchar(false)
	if err != nil {
		return nil, err
	}

	yCoordinateOfSampleOnSheet, err := r.ReadUchar(false)
	if err != nil {
		return nil, err
	}

	// undefined
	_, err = r.ReadBytes(34, false)
	if err != nil {
		return nil, err
	}

	sampleImageRawReader, err := r.ReNew(etl9gSampleSize, false)
	if err != nil {
		return nil, err
	}

	pixels := []uint8{}

	sampleImageRawBitReader := NewBitReader(sampleImageRawReader)
	for i := 0; i < etl9gSampleSize; i++ {
		px1, err := sampleImageRawBitReader.ReadUint(4)
		if err != nil {
			return nil, err
		}

		px2, err := sampleImageRawBitReader.ReadUint(4)
		if err != nil {
			return nil, err
		}

		// 4bitグレイスケールを8bitへ
		px1 = px1 * 256 / 16
		px2 = px2 * 256 / 16

		pixels = append(pixels, uint8(px1))
		pixels = append(pixels, uint8(px2))
	}

	sampleImage := image.NewGray(image.Rect(0, 0, etl9gSampleWidth, etl9gSampleHeight))
	sampleImage.Pix = pixels

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
	)

	return &record, nil
}

// ReadETL9GFile 指定ファイルパスのETL9Gファイルを読み込む
// - fpath: ETL9Gファイルパス
func ReadETL9GFile(fpath string) error {
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

		record, err := parseETL9GRecord(r2)
		if err != nil {
			return err
		}
		pp.Println(record)

		break
	}

	return nil
}
