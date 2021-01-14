package main

import (
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

type Config struct {
	filename string
}

func usage() {
	msg := `
Description: sort utf-8 text file by pinyin
Usage: psort <path/to/file>
For example: psort 1.txt
`
	fmt.Print(msg)
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	config := Config{
		filename: os.Args[1],
	}

	psort(config)
}

func psort(config Config) {
	// https://stackoverflow.com/questions/35080109/golang-how-to-read-input-filename-in-go/35080193#35080193
	data, err := ioutil.ReadFile(config.filename)
	if err != nil {
		fmt.Println("Can't read file:", os.Args[1])
		panic(err)
	}

	content := string(data)
	if !utf8.ValidString(content) {
		fmt.Println("Unsupported encoding, currently only support utf-8")
		os.Exit(0)
	}
	lines := strings.Split(content, "\n")
	pinyinArg := pinyin.NewArgs()
	// https://github.com/konglong87/golang_sort_by_chinese_pinyin/blob/master/sort_pinyin.go
	sort.Slice(lines, func(i, j int) bool {
		a := lines[i]
		b := lines[j]
		bLen := len(b)
		// https://stackoverflow.com/questions/15018545/how-to-index-characters-in-a-golang-string
		bRunes := []rune(b)
		// https://stackoverflow.com/questions/18130859/how-can-i-iterate-over-a-string-by-runes-in-go/18130921#18130921
		for idx, ai := range a {
			if idx > bLen-1 {
				return false
			}
			bi := bRunes[idx]
			if ai < 128 {
				if bi >= 128 {
					return true
				}

				if ai != bi {
					return ai < bi
				}
			}

			if ai >= 128 {
				if bi < 128 {
					return false
				}

				if ai != bi {
					aPinyin := pinyin.Pinyin(string(ai), pinyinArg)
					bPinyin := pinyin.Pinyin(string(bi), pinyinArg)
					return aPinyin[0][0] < bPinyin[0][0]
				}
			}
		}

		return true
	})
	fmt.Printf(strings.Join(lines, "\n"))
}
