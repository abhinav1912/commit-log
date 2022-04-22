package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

func (s *store) Append(data []byte) (bytesWritten uint64, position uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	position = s.size
	// writing the length of data to file first
	if err := binary.Write(s.buf, enc, uint64(len(data))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(data)
	if err != nil {
		return 0, 0, err
	}
	// account for the length of data written previously
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), position, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// flush any records that might be in the buffer
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	// read the record size from file
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	data := make([]byte, enc.Uint64(size))
	// read actual data from file
	if _, err := s.File.ReadAt(data, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *store) ReadAt(p []byte, offset int64) (int, error) {
	// helper func to read data at given offset into argument byte array
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, offset)
}

func (s *store) Close() error {
	// flush anu buffered data and close the file
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
