package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"os"

	tfExample "github.com/tensorflow/tensorflow/tensorflow/go/core/example/example_protos_go_proto"
	"google.golang.org/protobuf/proto"
)

// https://github.com/tensorflow/tensorflow/blob/051a96f3ec4fc38b248e8ae8ad2f8ad124eda59b/tensorflow/core/lib/hash/crc32c.h
const maskDelta uint32 = 0xa282ead8

// https://github.com/tensorflow/tensorflow/blob/051a96f3ec4fc38b248e8ae8ad2f8ad124eda59b/tensorflow/core/lib/hash/crc32c.h#L53-L56
func mask(crc uint32) uint32 {
	return ((crc >> 15) | (crc << 17)) + maskDelta
}

var crc32Table = crc32.MakeTable(crc32.Castagnoli)

func crc32Hash(data []byte) uint32 {
	return crc32.Checksum(data, crc32Table)
}

func uint64ToBytes(x uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, x)
	return b
}

func Write(w io.Writer, data []byte) (int, error) {
	// Write based on format specified in https://github.com/tensorflow/tensorflow/blob/051a96f3ec4fc38b248e8ae8ad2f8ad124eda59b/tensorflow/core/lib/io/record_writer.cc#L124-L128
	//  uint64    length
	//  uint32    masked crc of length
	//  byte      data[length]
	//  uint32    masked crc of data

	length := uint64(len(data))
	lengthCRC := mask(crc32Hash(uint64ToBytes(uint64(len(data)))))
	dataCRC := mask(crc32Hash(data))

	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return 0, err
	}

	if err := binary.Write(w, binary.LittleEndian, lengthCRC); err != nil {
		return 0, err
	}

	if _, err := w.Write(data); err != nil {
		return 0, err
	}

	if err := binary.Write(w, binary.LittleEndian, dataCRC); err != nil {
		return 0, err
	}

	// return size written
	return binary.Size(dataCRC) + len(data) + binary.Size(length) + binary.Size(lengthCRC), nil
}

func Read(r io.Reader) (data []byte, err error) {
	var (
		length         uint64
		lengthChecksum uint32
		dataChecksum   uint32
	)

	// get data length
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	// get length checksum
	if err := binary.Read(r, binary.LittleEndian, &lengthChecksum); err != nil {
		return nil, err
	}

	// get data
	data = make([]byte, length)
	if _, err := r.Read(data); err != nil {
		return nil, err
	}

	// get data checksum
	if err := binary.Read(r, binary.LittleEndian, &dataChecksum); err != nil {
		return nil, err
	}

	// check checksum length
	if actual := mask(crc32Hash(uint64ToBytes(length))); actual != lengthChecksum {
		return nil, errors.New("corrupted record, length checksum doesn't match")
	}

	// check data checksum
	if actual := mask(crc32Hash(data)); actual != dataChecksum {
		return nil, errors.New("corrupted record, data checksum doesn't match")
	}

	return data, nil
}

func main() {
	tfrecordFile := "example.tfrecord"
	exampleWriter, err := os.Create(tfrecordFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	example := tfExample.Example{
		Features: &tfExample.Features{
			Feature: map[string]*tfExample.Feature{
				"x": &tfExample.Feature{
					Kind: &tfExample.Feature_Int64List{
						Int64List: &tfExample.Int64List{
							Value: []int64{1, 2, 3, 4, 5},
						},
					},
				},
				"y": &tfExample.Feature{
					Kind: &tfExample.Feature_BytesList{
						BytesList: &tfExample.BytesList{
							Value: [][]byte{[]byte("hello")},
						},
					},
				},
				"z": &tfExample.Feature{
					Kind: &tfExample.Feature_FloatList{
						FloatList: &tfExample.FloatList{
							Value: []float32{0.1, 0.2, 0.3, 0.4, 0.5},
						},
					},
				},
			},
		},
	}

	exampleBytes, err := proto.Marshal(&example)
	if err != nil {
		log.Fatal(err.Error())
	}

	Write(exampleWriter, exampleBytes)

	exampleReader, err := os.Open(tfrecordFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	exampleReaderBytes, err := ioutil.ReadAll(exampleReader)
	if err != nil {
		log.Fatal(err.Error())
	}

	exampleBytesReader := bytes.NewReader(exampleReaderBytes)

	for {
		var content []byte

		content, err = Read(exampleBytesReader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err.Error())
		}

		m := tfExample.Example{}
		err = proto.Unmarshal(content, &m)
		if err != nil {
			log.Fatal(err.Error())
		}

		x := m.Features.GetFeature()["x"].GetInt64List().Value
		y := m.Features.GetFeature()["y"].GetBytesList().Value
		z := m.Features.GetFeature()["z"].GetFloatList().Value

		fmt.Println(x, string(y[0]), z)
	}
}
