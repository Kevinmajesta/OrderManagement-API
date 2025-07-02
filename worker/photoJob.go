package worker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type PhotoJob struct {
	Src      io.Reader
	Filename string
	Ext      string
}

var PhotoQueue = make(chan PhotoJob, 100) // buffer 100 foto

func StartPhotoWorker() {
	for i := 0; i < 3; i++ {
		go func(workerID int) {
			for job := range PhotoQueue {
				savePath := filepath.Join("assets", "photos", job.Filename+job.Ext)

				dst, err := os.Create(savePath)
				if err != nil {
					fmt.Printf("[PhotoWorker %d] Failed to create photo file: %v\n", workerID, err)
					continue
				}

				_, err = io.Copy(dst, job.Src)
				if err != nil {
					fmt.Printf("[PhotoWorker %d] Failed to copy photo: %v\n", workerID, err)
				} else {
					fmt.Printf("[PhotoWorker %d] Photo saved: %s\n", workerID, savePath)
				}
				dst.Close()
			}
		}(i)
	}
}
