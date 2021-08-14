package main

import (
	"errors"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileDoesNotExists     = errors.New("file does not exists")
	ErrUnexpectedError       = errors.New("unexpected error")
	ErrWrongArguments        = errors.New("wrong or negative arguments")
)

// Очевидно задачу можно было бы намного проще реализовать используя функцию
// io.CopyN, но необходимо было реализовать кастомный прогресс-бар, для чего
// пришлось создавать цикл записи, через который и можно отслеживать этот
// самый прогресс копирования.

const (
	stepSize = 150
)

var displayProgressBar = true

func openFile(path string, flag int) (*os.File, error) {
	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		if os.IsNotExist(err) && flag&os.O_CREATE == 0 {
			return nil, ErrFileDoesNotExists
		} else if os.IsNotExist(err) {
			return nil, ErrUnexpectedError
		}
	}
	if flag&os.O_CREATE != 0 {
		file.Truncate(0)
	}

	return file, err
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := openFile(fromPath, os.O_RDONLY)
	if err != nil {
		return err
	}
	to, err := openFile(toPath, os.O_WRONLY|os.O_CREATE)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	size := fileInfo.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}
	limit += offset
	if limit > size || limit == 0 {
		limit = size
	}
	if offset < 0 || limit < 0 {
		return ErrWrongArguments
	}

	fileBytes := make([]byte, limit)
	file.ReadAt(fileBytes, offset)

	pb := NewProgressBar(uint(limit-offset), displayProgressBar)
	step := int64(stepSize)
	for offset < limit {
		if offset+step > limit {
			step = limit - offset
		}
		written, _ := to.Write(fileBytes[:step])
		fileBytes = fileBytes[step:]
		offset += step

		pb.Process(uint(written))
	}
	file.Close()
	to.Close()

	return nil
}
