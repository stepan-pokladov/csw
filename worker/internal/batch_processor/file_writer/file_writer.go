package file_writer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileWriter struct{}

func init() {
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
	}
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

// ProcessBatch writes the document to a file
func (fw *FileWriter) ProcessBatch(topic string, s []string) error {
	fo, err := os.Create("logs/" + topic + "-" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".json")
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	w := bufio.NewWriter(fo)
	for _, str := range s {
		_, err = w.WriteString(str + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return err
		}
	}

	if err = w.Flush(); err != nil {
		fmt.Printf("Error flushing to file: %v\n", err)
		return err
	}

	return nil
}
