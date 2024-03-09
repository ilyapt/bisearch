package bisearch

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrIncorrectLength         = errors.New("incorrect length")
	ErrOperationIsNotCompleted = errors.New("operation is not completed")
	ErrNotExist                = errors.New("key does not exist")
)

func Search(f io.ReadSeeker, from, count uint64, length int, key []byte) ([]byte, error) {
	if length < len(key) {
		return nil, ErrIncorrectLength
	}

	l, r := uint64(0), count-1
	record := make([]byte, length)

	for r <= count && l >= 0 && r >= l {
		m := l + (r-l)/2
		pos := int64(from + m*uint64(length))

		if n, err := f.Seek(pos, io.SeekStart); err != nil {
			return nil, err
		} else if n != pos {
			return nil, ErrOperationIsNotCompleted
		}

		if n, err := f.Read(record); err != nil {
			return nil, err
		} else if n != length {
			return nil, ErrOperationIsNotCompleted
		}

		if cmp := bytes.Compare(key, record[:len(key)]); cmp == 0 {
			return record, nil
		} else if cmp > 0 {
			l = m + 1
		} else {
			r = m - 1
		}

	}

	return nil, ErrNotExist
}
