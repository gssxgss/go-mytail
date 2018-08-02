//https://www.linkedin.com/pulse/golang-simple-tail-tagir-ismagilov/
//https://github.com/hpcloud/tail/blob/master/tail.go
//https://gobyexample.com/command-line-flags
//https://gobyexample.com/command-line-arguments
//https://gobyexample.com/reading-files
//https://medium.freecodecamp.org/writing-command-line-applications-in-go-2bc8c0ace79d
//https://github.com/coreymgilmore/gotail/blob/master/gotail.go
//https://gobyexample.com/command-line-flags
package main

import (
	"os"
	"io"
	"fmt"
	"errors"
	"flag"
	)

func tail(file *os.File, line int) {

	// NOTE io.SeekEnd(=> 2) means relative to the end
	offset, err := file.Seek(0, io.SeekEnd, )

	if err != nil {
		panic(err)
	}

	lineBreakCount := 0
	taskFinished := false
	var tailContent string

	for offset > 0 {

		// offset should larger than 0
		offset -= 1024
		if offset < 0 {
			offset = 0
		}

		readBuffer := make([]byte, 1024)
		readBytes, err := file.ReadAt(readBuffer, offset)
		currentIndex := readBytes - 1

		if err != nil && err != io.EOF {
			panic(err)
		}

		//fmt.Println("# read bytes: " + strconv.Itoa(readBytes))
		//fmt.Println("# current index: " + strconv.Itoa(currentIndex))
		//fmt.Println("# first byte equal \\n: " + strconv.FormatBool(readBuffer[currentIndex] == '\n'))
		//fmt.Println("# lines to print: " + strconv.Itoa(line))

		for currentIndex > -1 {
			currentByte := readBuffer[currentIndex]

			if currentByte == '\n' {
				lineBreakCount += 1
			}

			if lineBreakCount >= line {
				//fmt.Println("# line break count: " + strconv.Itoa(lineBreakCount))
				tailContent = string(readBuffer[currentIndex+1:readBytes]) + tailContent
				taskFinished = true
				break
			} else if currentIndex == 0 && err == io.EOF {
				tailContent = string(readBuffer[currentIndex:readBytes]) + tailContent
				taskFinished = true
				break
			}
			currentIndex--
		}

		if taskFinished {
			break
		}

	}
	fmt.Print(tailContent)

}

func main() {
	line := 10

	flag.IntVar(&line, "n", line, "lines to tail")
	flag.Parse() // NOTE should put directly under flags

	files := flag.Args()

	//fmt.Println(line)
	//fmt.Println(files)

	if line < 1 {
		panic(errors.New("no lines to tail"))
	}

	filesCount := len(files)
	if len(files) < 1 {
		panic(errors.New("no files to tail"))
	}

	for i := 0; i < filesCount; i++ {
		file, err := os.Open(files[i])

		if err != nil {
			panic(err)
		}

		fmt.Println("==> " + files[i] + "<==")
		tail(file, line)
		file.Close()
		if i != filesCount-1 {
			fmt.Println()
		}
	}

}
