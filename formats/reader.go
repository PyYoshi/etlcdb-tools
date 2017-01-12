package formats

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
)

type BinReader struct {
	raw    *[]byte
	buffer *bytes.Reader
}

func NewBinReader(b []byte) (*BinReader, error) {
	buf := bytes.NewReader(b)
	return &BinReader{
		&b,
		buf,
	}, nil
}

func NewBinReaderFromFilePath(path string) (*BinReader, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewBinReader(b)
}

func (r *BinReader) ReadByte() (byte, error) {
	return r.buffer.ReadByte()
}

func selectByteOrder(isLittleEndian bool) binary.ByteOrder {
	var endian binary.ByteOrder
	if isLittleEndian {
		endian = binary.LittleEndian
	} else {
		endian = binary.BigEndian
	}
	return endian
}

func (r *BinReader) ReadChar(isLittleEndian bool) (int8, error) {
	var i int8
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *BinReader) ReadUchar(isLittleEndian bool) (uint8, error) {
	var i uint8
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *BinReader) ReadShort(isLittleEndian bool) (int16, error) {
	var i int16
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *BinReader) ReadUshort(isLittleEndian bool) (uint16, error) {
	var i uint16
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *BinReader) ReadInt(isLittleEndian bool) (int32, error) {
	var i int32
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *BinReader) ReadUint(isLittleEndian bool) (uint32, error) {
	var i uint32
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *BinReader) ReadLong(isLittleEndian bool) (int64, error) {
	var i int64
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *BinReader) ReadUlong(isLittleEndian bool) (uint64, error) {
	var i uint64
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *BinReader) ReadBytes(size int, isLittleEndian bool) ([]byte, error) {
	b := make([]byte, size)
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func (r *BinReader) ReNew(size int, isLittleEndian bool) (*BinReader, error) {
	b := make([]byte, size)
	err := binary.Read(r.buffer, selectByteOrder(isLittleEndian), &b)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(b)
	return &BinReader{
		&b,
		buf,
	}, nil
}
