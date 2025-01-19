package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"context"
	"sync"
	"golang.org/x/sync/semaphore"
)

var first = []int{
	0, 1, 2, 3,
	8, 9, 10, 11,
	16, 17, 18, 19,
	24, 25, 26, 27}

var second = []int{
	7, 6, 5, 4,
	15, 14, 13, 12,
	23, 22, 21, 20,
	31, 30, 29, 28}

var third = []int{
	56, 57, 58, 59,
	48, 49, 50, 51,
	40, 41, 42, 43,
	32, 33, 34, 35}

var fourth = []int{
	63, 62, 61, 60,
	55, 54, 53, 52,
	47, 46, 45, 44,
	39, 38, 37, 36}

type record struct {
	BPName string `json:"blue player name"`
	ORName string `json:"orange player name"`
	Record []int  `json:"record"`
}
func main() {
	exec()
}
const concurrency = 30
func exec(){
	ctx := context.TODO()
	var s *semaphore.Weighted = semaphore.NewWeighted(concurrency)
	var wg sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := s.Acquire(ctx, 1); err != nil {
				return
			}
			defer s.Release(1)
			doRotate(i)
		}(i)
	}
	wg.Wait()
}

func doRotate(filenum int) {
	fmt.Println("File number:", filenum)
	filename := strconv.Itoa(filenum) + ".json"
	//filenameのファイルを開く
	file, err := os.Open("./record/" + filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	//ファイルの中身を読み込む
	var buf = make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	//bufをrecordに変換
	var rec record
	err = json.Unmarshal(buf[:n], &rec)
	if err != nil {
		fmt.Println("Error unmarshaling:", err)
		return
	}
	//rec.Record[0]がどの配列にあるかを判定
	if slices.Contains(first, rec.Record[0]) {
		replace(rec.Record, 1)
	} else if slices.Contains(second, rec.Record[0]) {
		replace(rec.Record, 2)
	} else if slices.Contains(third, rec.Record[0]) {
		replace(rec.Record, 3)
	} else if slices.Contains(fourth, rec.Record[0]) {
		replace(rec.Record, 4)
	} else {
		panic("Error: invalid record")
	}
	//recをjsonに変換
	buf, err = json.Marshal(rec)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}
	//filenameのファイルを開く
	file, err = os.Create("./result/" + filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	//bufをファイルに書き込む
	_, err = file.Write(buf)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}

func replace(rec []int, place int) {
	switch place {
	case 1:
		return
	case 2:
		for i, item := range rec {
			rec[i] = item - (item%8)*2 + 7
		}
	case 3:
		for i, item := range rec {
			rec[i] = item + 8*(7-(item/8)*2)
		}
	case 4:
		for i, item := range rec {
			rec[i] = 63 - item
		}
	default:
		panic("Error: invalid place")
	}
}

//test