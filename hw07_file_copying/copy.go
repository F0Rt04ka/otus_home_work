package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	fileCloseFn := func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}

	defer fileCloseFn(fileFrom)

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileCloseFn(fileTo)

	fileStat, err := fileFrom.Stat()
	if err != nil {
		return err
	}

	fileSize := fileStat.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if fileStat.Mode()&os.ModeDevice != 0 {
		return ErrUnsupportedFile
	}

	if limit > fileSize || limit == 0 {
		limit = fileSize
	}

	copyFn := func(from, to *os.File) error {
		bufferSize := int64(100)
		if bufferSize > limit {
			bufferSize = limit
		}

		buffer := make([]byte, bufferSize)
		alreadyReadBytes := int64(0)

		bar := pb.StartNew(0)
		bar.SetTotal(limit)
		if limit+offset > fileSize {
			bar.SetTotal(fileSize - offset)
		}

		defer func() {
			bar.Finish()
		}()

		for alreadyReadBytes < limit {
			readBytes, err := from.ReadAt(buffer, offset+alreadyReadBytes)
			bar.Add(readBytes)

			if err != nil {
				if errors.Is(err, io.EOF) {
					alreadyReadBytes = limit
				} else {
					return err
				}
			}

			alreadyReadBytes += int64(readBytes)

			_, err = to.Write(buffer[:readBytes])
			if err != nil {
				return err
			}
		}

		return nil
	}

	copyFn(fileFrom, fileTo)

	return nil
}
