package log

import (
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offWidth   uint64 = 4
	posWidth   uint64 = 8
	totalWidth        = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())
	// truncate needed since files being extended at runtime
	// last record not guaranteed to be at end of file enless extra space removed
	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}
	return idx, nil
}

func (i *index) Close() error {
	// sync data to persisted file
	if err := i.mmap.Sync(gommap.MS_ASYNC); err != nil {
		return err
	}
	// flush contents to stable storage
	if err := i.file.Sync(); err != nil {
		return err
	}
	// remove extra memory at the end
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	return i.file.Close()
}
