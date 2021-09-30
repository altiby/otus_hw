package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	inputFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed create input file %w", err)
	}
	defer inputFile.Close()

	if _, err := os.Stat(toPath); err == nil {
		err = os.Remove(toPath)
		if err != nil {
			return fmt.Errorf("failed remove old output file %w", err)
		}
	}

	outputFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("failed create output file %w", err)
	}
	defer outputFile.Close()

	stat, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed get input file size %w", err)
	}

	fileSize := stat.Size()
	if fileSize == 0 {
		return ErrUnsupportedFile
	}
	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	_, err = inputFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed seek input file %w", err)
	}

	buffer := make([]byte, 256)
	var writtenByte int64
	for {
		readByte, err := inputFile.Read(buffer)
		if err == nil || err == io.EOF {
			if limit > 0 && writtenByte+int64(readByte) > limit {
				readByte = int(limit - writtenByte)
				err = io.EOF
			}

			writeByte, writeErr := outputFile.Write(buffer[:readByte])
			writtenByte += int64(writeByte)
			if writeErr != nil {
				return fmt.Errorf("failed write output file %w", writeErr)
			}
			if err == io.EOF {
				break
			}
		} else {
			return fmt.Errorf("failed reade input file %w", err)
		}
	}
	return nil
}
