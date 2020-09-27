package sensitive

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	fileName = "sensitive.txt"
	TrieFilter *Trie
)

func init() {
	TrieFilter = NewTrie()
	if err := readSensitiveWords(); err != nil {
		panic(err)
	}
	fmt.Println("Load sensitiveWord success")
}

func readSensitiveWords()  error{
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return err
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	list := make([]string, 0)
	for {
		line, _, err := buf.ReadLine()
		list = append(list, string(line))
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Read file error!", err)
				return err
			}
		}
	}
	TrieFilter.Add(list...)
	return nil
}

func Replace(input string, character rune)  string{
	return TrieFilter.Replace(input, character)
}

