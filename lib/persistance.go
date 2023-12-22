package lib

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

type Persistance struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

// create or open the file based on its existence
// read from the file into the Persistance type
// start a goroutine to sync the file to disk every 1 second while the server is running
// this still has problem of durability, talking about milliseconds here. for a full durable system we can push data to a file after every commands but that can be costly as io operations are expensive
func NewPersistance(path string) (*Persistance, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	persistant := &Persistance{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// start go routine to sync persistantFile to disk every 1 second
	go func() {
		for {
			persistant.mu.Lock()
			persistant.file.Sync()
			persistant.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return persistant, nil
}

// close the file properly when the server shuts down
func (persistantFile *Persistance) Close() error {
	persistantFile.mu.Lock()
	defer persistantFile.mu.Unlock()

	return persistantFile.file.Close()
}

// write commands to the file
func (persistantFile *Persistance) Write(value Value) error {
	persistantFile.mu.Lock()
	defer persistantFile.mu.Unlock()

	_, err := persistantFile.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

// read commands/logs from the file
func (persistantFile *Persistance) Read(fn func(value Value)) error {
	persistantFile.mu.Lock()
	defer persistantFile.mu.Unlock()

	persistantFile.file.Seek(0, io.SeekStart)

	reader := NewResp(persistantFile.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
		fn(value)
	}

	return nil
}
