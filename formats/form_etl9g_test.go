package formats

import "testing"

func TestReadETL9GFile(t *testing.T) {
	_, err := ReadETL9GFile("../etlcdb/ETL9G/ETL9G_01")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMakeETL9GDatasets(t *testing.T) {
	err := MakeETL9GDatasets("../etlcdb/ETL9G", "../datasets/ETL9G")
	if err != nil {
		t.Fatal(err)
	}
}
