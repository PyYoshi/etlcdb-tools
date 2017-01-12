package formats

import "testing"
import "log"

func TestNewFormatETL9G(t *testing.T) {
	// tests := []struct {
	// 	name string
	// 	want FormatETL9G
	// }{
	// // TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		if got := NewFormatETL9G(); !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("NewFormatETL9G() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}

func TestFormatETL9G_ReadFiles(t *testing.T) {
	format := NewFormatETL9G()
	format.ReadFiles("../etlcdb/ETL9G/")
	// type fields struct {
	// 	recordSize   int
	// 	sampleWidth  int
	// 	sampleHeight int
	// 	fileNum      int
	// 	recordNum    int
	// }
	// type args struct {
	// 	dir string
	// }
	// tests := []struct {
	// 	name   string
	// 	fields fields
	// 	args   args
	// }{
	// // TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		f := &FormatETL9G{
	// 			recordSize:   tt.fields.recordSize,
	// 			sampleWidth:  tt.fields.sampleWidth,
	// 			sampleHeight: tt.fields.sampleHeight,
	// 			fileNum:      tt.fields.fileNum,
	// 			recordNum:    tt.fields.recordNum,
	// 		}
	// 		f.ReadFiles(tt.args.dir)
	// 	})
	// }
}

func TestBit(t *testing.T) {
	v := 0x10110110
	log.Printf("値: 0x%x", v)
	log.Printf("下位: 0x%x", v&0x00001111)
	log.Printf("上位: 0x%x", v&0x11110000>>16)
}

func TestFormatETL9G_ReadFile(t *testing.T) {
	format := NewFormatETL9G()
	err := format.ReadFile("../etlcdb/ETL9G/ETL9G_01")
	if err != nil {
		t.Fatal(err)
	}
	// type fields struct {
	// 	recordSize   int
	// 	sampleWidth  int
	// 	sampleHeight int
	// 	fileNum      int
	// 	recordNum    int
	// }
	// type args struct {
	// 	fpath string
	// }
	// tests := []struct {
	// 	name   string
	// 	fields fields
	// 	args   args
	// }{
	// // TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		f := &FormatETL9G{
	// 			recordSize:   tt.fields.recordSize,
	// 			sampleWidth:  tt.fields.sampleWidth,
	// 			sampleHeight: tt.fields.sampleHeight,
	// 			fileNum:      tt.fields.fileNum,
	// 			recordNum:    tt.fields.recordNum,
	// 		}
	// 		f.ReadFile(tt.args.fpath)
	// 	})
	// }
}
