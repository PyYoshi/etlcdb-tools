package formats

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
}
