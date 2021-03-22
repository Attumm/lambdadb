package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"runtime"
)

var STORAGEFUNCS storageFuncs
var RETRIEVEFUNCS retrieveFuncs

func init() {
	STORAGEFUNCS = make(storageFuncs)
	STORAGEFUNCS["bytes"] = saveAsBytes // currently default
	STORAGEFUNCS["bytesz"] = saveAsBytesCompressed
	STORAGEFUNCS["json"] = saveAsJsonZipped
	STORAGEFUNCS["jsonz"] = saveAsJsonZipped

	RETRIEVEFUNCS = make(retrieveFuncs)
	RETRIEVEFUNCS["bytes"] = loadAsBytes // currently default
	RETRIEVEFUNCS["bytesz"] = loadAsBytesCompressed
	RETRIEVEFUNCS["json"] = loadAsJsonZipped
	RETRIEVEFUNCS["jsonz"] = loadAsJsonZipped
}

func saveAsJsonZipped(items Items, filename string) (int64, error) {
	var b bytes.Buffer
	writer := gzip.NewWriter(&b)
	itemJSON, _ := json.Marshal(ITEMS)
	writer.Write(itemJSON)
	writer.Flush()
	writer.Close()
	err := ioutil.WriteFile(filename, b.Bytes(), 0666)
	if err != nil {
		return 0, err
	}
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	size := fi.Size()
	return size, nil
}

func saveAsBytes(items Items, filename string) (int64, error) {
	data := EncodeItems(items)
	WriteToFile(data, filename)
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	size := fi.Size()
	return size, nil
}

func saveAsBytesCompressed(items Items, filename string) (int64, error) {
	data := EncodeItems(items)
	data = Compress(data)
	WriteToFile(data, filename)
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	size := fi.Size()
	return size, nil
}

func EncodeItems(items Items) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(items)
	if err != nil {
		fmt.Println("error encoding", err)
	}
	return buf.Bytes()
}

func Compress(s []byte) []byte {
	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(s)
	zipped.Close()
	return zipbuf.Bytes()
}

func Decompress(s []byte) []byte {
	//TODO check empty
	reader, _ := gzip.NewReader(bytes.NewReader(s))
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("Unable to Decompress", err)
	}
	reader.Close()
	return data
}

func DecodeToItems(s []byte) Items {
	items := make(Items, 0, 100*1000)
	decoder := gob.NewDecoder(bytes.NewReader(s))
	err := decoder.Decode(&items)
	if err != nil {
		fmt.Println("Unable to DecodeToItems", err)
	}
	return items
}

func WriteToFile(s []byte, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Unable to WriteToFile", err)
	}
	f.Write(s)
}

func ReadFromFile(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Unable to ReadFromFile", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Unable to ReadFromFile1", err)
	}
	return data
}

func loadAsBytes(items Items, filename string) (int, error) {
	d := ReadFromFile(filename)
	items = DecodeToItems(d)
	ITEMS = items
	runtime.GC()
	items = nil
	return len(ITEMS), nil
}

func loadAsBytesCompressed(items Items, filename string) (int, error) {
	d := ReadFromFile(filename)
	d = Decompress(d)
	items = DecodeToItems(d)
	ITEMS = make(Items, 0, 100*1000)
	ITEMS = items
	return len(items), nil
}

func loadAsJsonZipped(items Items, filename string) (int, error) {
	fi, err := os.Open(filename)
	if err != nil {
		_, err2 := os.Getwd()
		if err2 != nil {
			return 0, err2
		}
		return 0, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return 0, err
	}
	defer fz.Close()

	// TODO buffered instead of one big chunk
	s, err := ioutil.ReadAll(fz)

	if err != nil {
		return 0, err
	}

	ITEMS = make(Items, 0, 100*1000)
	json.Unmarshal(s, &ITEMS)

	// GC friendly
	s = nil
	return len(ITEMS), nil
}

func loadAtStart(storagename string, filename string, indexed bool) {

	retrievefunc, found := RETRIEVEFUNCS[storagename]
	if !found {
		fmt.Println("storage mehod not found, trying with bytes")
		storagename := "bytes"
		retrievefunc = RETRIEVEFUNCS[storagename]
	}

	filename = fmt.Sprintf("%s.%s", FILENAME, storagename)
	msg := fmt.Sprintf("retrieving with: %s, with filename: %s", storagename, filename)
	fmt.Printf(WarningColorN, msg)

	start := time.Now()
	itemsAdded, err := retrievefunc(ITEMS, filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("could not open %s reason %s", filename, err))
	}

	diff := time.Since(start)
	msg = fmt.Sprint("Loaded in memory amount: ", itemsAdded, " time: ", diff)
	fmt.Printf(WarningColorN, msg)
	runtime.GC()

	if indexed {
		start = time.Now()
		msg := fmt.Sprint("Creating index")
		fmt.Printf(WarningColorN, msg)
		makeIndex()
		diff = time.Since(start)
		msg = fmt.Sprint("Index set time: ", diff)
		fmt.Printf(WarningColorN, msg)
		runtime.GC()
	}
}
