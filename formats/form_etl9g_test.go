package formats

import "testing"

func TestReadETL9GFile(t *testing.T) {
	_, err := ReadETL9GFile("../etlcdb/ETL9G/ETL9G_01")
	if err != nil {
		t.Fatal(err)
	}
}
