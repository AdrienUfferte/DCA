// ****************************************************************************
//
// MIT License
//
// Copyright (c) 2021 Adrien Ufferte
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// ****************************************************************************

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func dca(date time.Time, keys []time.Time, datas map[time.Time]int) (float64, int, float64) {
	// log.Println("dca(", date, ")")
	totalInvested := 0
	totalResult := float64(0)
	dateIndex := 0
	for i := range keys {
		if keys[i] == date {
			dateIndex = i
			// log.Println("dateIndex = ", dateIndex, "/", len(keys))
			break
		}
	}
	ratio := float64(100)
	for dateIndex < len(keys)-1 {
		totalInvested += 100
		totalResult += 100
		ratio = float64(datas[keys[dateIndex+1]]) / float64(datas[keys[dateIndex]])
		//log.Println("ratio = ", ratio)
		totalResult *= ratio
		dateIndex++
	}
	// log.Println("totalInvested = ", totalInvested)
	// log.Println("totalResult = ", totalResult)
	if totalResult <= 0 {
		return 1, 100, 100
	}
	return (totalResult / float64(totalInvested)), totalInvested, totalResult
}

func main() {

	absolutePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("error: 1\n")
	}

	dataPath := filepath.Join(absolutePath, "Weekly.txt")
	osFile, err := os.Open(dataPath)
	if err != nil {
		log.Fatal("error: 2\n")
	}

	data, err := ioutil.ReadAll(osFile)
	if err != nil {
		log.Fatal(err)
	}

	stringData := strings.ReplaceAll(string(data), "\r\n", "\n")

	tempDatas := strings.Split(stringData, "\n")

	finalDatas := make(map[time.Time]int)

	layout := "02/01/2006"
	for _, data := range tempDatas {
		data1 := strings.Split(data, ", ")
		time, err := time.Parse(layout, data1[0])
		if err != nil {
			log.Fatal(err)
		}
		value, err := strconv.Atoi(data1[1])
		if err != nil {
			log.Fatal(err)
		}
		finalDatas[time] = value
	}

	// log.Println(finalDatas)
	keys := make([]time.Time, 0)

	for date := range finalDatas {
		keys = append(keys, date)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	for _, date := range keys {
		y, m, d := date.Date()
		fmt.Printf("* %-4v", y)
		fmt.Printf(" %-10v", m)
		fmt.Printf("%02d", d)
		ratio, invested, result := dca(date, keys, finalDatas)
		fmt.Printf(fmt.Sprintf(" = ratio : %05.2f", ratio))
		fmt.Printf(" ( invested : %d, result : %.2f )\n", invested, result)
		//fmt.Println(" = ", dca(date, keys, finalDatas))
	}
}
