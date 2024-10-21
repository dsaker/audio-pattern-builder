package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
)

type Audio struct {
	Id     int
	Native bool
}

const (
	numberPhrases = 1000
	sliceSize     = 25 * numberPhrases
	limit         = 800
)

func main() {

	m := make([]Audio, sliceSize)
	multi := 10
	for i := 1; i <= numberPhrases; i++ {
		first := findNextFree(m, i)
		if first == numberPhrases {
			fmt.Println(i)
			break
		}
		m = fillNextIfFree(m, true, first, i)
		m = fillNextIfFree(m, false, first+1, i)
		m = fillNextIfFree(m, false, first+2, i)

		count := 1
		j := (first + 2) + multi*count
		for j < sliceSize {
			j = findNextFree(m, j)
			m = fillNextIfFree(m, true, j, i)
			m = fillNextIfFree(m, false, j+1, i)
			count = count + 1
			j = j + multi*count*(count)
		}
	}

	// remove all zeros from the array
	var newArray []Audio
	for i := 0; i < len(m); i++ {
		if m[i].Id != 0 {
			newArray = append(newArray, m[i])
		}
	}

	//fmt.Println(len(newArray))

	//for i := 0; i < 5; i++ {
	//	j := 0
	//	count := 0
	//	for j < len(newArray) {
	//		if newArray[j].Id == i && newArray[j].Native == true {
	//			count++
	//			fmt.Printf("%d: line: %d\n", i, j)
	//		}
	//		j++
	//	}
	//	fmt.Printf("%d: %d\n", i, count)
	//}

	// reduce all values by 1 so the pattern can map to other slices
	// (starts at zero)
	for i := 0; i < len(newArray); i++ {
		newArray[i].Id = newArray[i].Id - 1
	}

	file, err := os.OpenFile("/Users/dustysaker/GolandProjects/audioForLoop/audioPattern.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//fmt.Println(newArray)
	// Print the array
	for i := 0; i < len(newArray); i++ {
		printStructWithCommas(newArray[i], file)
		if i == limit {
			break
		}
	}
}

func fillNextIfFree(m []Audio, native bool, next, i int) []Audio {
	if next >= sliceSize {
		return m
	}
	if m[next].Id == 0 {
		m[next] = Audio{
			Id:     i,
			Native: native,
		}
		return m
	}
	m = slices.Insert(m, next, Audio{
		Id:     i,
		Native: native,
	})
	return m
}

func findNextFree(m []Audio, i int) int {
	for i < sliceSize && m[i].Id != 0 {
		i++
	}
	return i
}

func printStructWithCommas(s interface{}, file *os.File) {
	v := reflect.ValueOf(s)

	if v.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		return
	}

	// Write the string without a newline
	_, err := io.WriteString(file, fmt.Sprintf("{%v, %v},", v.Field(0).Interface(), v.Field(1).Interface()))
	if err != nil {
		panic(err)
	}

	//fmt.Printf("{%v, %v},", v.Field(0).Interface(), v.Field(1).Interface())
}
