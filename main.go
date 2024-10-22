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
	numberPhrases = 1000               // numberPhrases is set higher than limit to allow for more repetitions
	sliceSize     = 20 * numberPhrases // you can adjust sliceSize smaller or bigger to change final file size
	limit         = 800                // set this to where you want to set the limit of phrase
)

func main() {

	m := make([]Audio, sliceSize)
	// increasing the multiplier will increase the distance between repetitions of a phrase
	multi := 8
	for i := 1; i <= numberPhrases; i++ {
		// find the next free index
		first := findNextFree(m, i)

		// the firsts time you hear a phrase it will be said once in your native language
		// then repeated twice in the language you want to learn
		m = fillNextIfFree(m, true, first, i)
		m = fillNextIfFree(m, false, first+1, i)
		m = fillNextIfFree(m, false, first+2, i)

		// add one native phrase and one of the language you want to learn spread out
		// across the array with increasing space between
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

	// uncomment to print info describing the array
	//for i := 0; i < 5; i++ {
	//	j := 0
	//	count := 0
	//	for j < len(newArray) {
	//		if newArray[j].Id == i && newArray[j].Native == true {
	//			count++
	//			// this will give you an idea of the spacing between repetition
	//			fmt.Printf("%d: line: %d\n", i, j)
	//		}
	//		j++
	//	}
	//	// total number of times it appears in the array
	//	fmt.Printf("%d: %d\n", i, count)
	//}

	// reduce all values by 1 so the pattern can map to other slices
	// (starts at zero)
	for i := 0; i < len(newArray); i++ {
		newArray[i].Id = newArray[i].Id - 1
	}

	// get working directory to print out pattern
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// opens a file to write to in append mode, emptying it first
	file, err := os.OpenFile(pwd+"/audioPattern.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Print the array to the file
	for i := 0; i < len(newArray); i++ {
		printStructWithCommas(newArray[i], file)
		if i == limit {
			break
		}
	}
}

// fillNextIfFree adds the next Audio struct to the next index if free or inserts it if not
func fillNextIfFree(m []Audio, native bool, next, i int) []Audio {
	if next >= sliceSize {
		return m
	}
	if m[next].Id == 0 {
		m[next].Id = i
		m[next].Native = native
		return m
	}
	m = slices.Insert(m, next, Audio{
		Id:     i,
		Native: native,
	})
	return m
}

// findNextFree finds the next free index in the Audio struct slice
func findNextFree(m []Audio, i int) int {
	for i < sliceSize && m[i].Id != 0 {
		i++
	}
	return i
}

// printStructWithCommas writes the struct to a file in the same way it needs to appear in code
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
}
