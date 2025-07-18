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
	// Hanya meluncurkan SATU goroutine worker
	go func() { // Tidak perlu workerID jika hanya ada satu
		for job := range PhotoQueue {
			savePath := filepath.Join("assets", "photos", job.Filename+job.Ext)

			dst, err := os.Create(savePath)
			if err != nil {
				// Gunakan fmt.Println atau fmt.Printf tanpa workerID
				fmt.Printf("Failed to create photo file: %v\n", err)
				continue
			}

			_, err = io.Copy(dst, job.Src)
			if err != nil {
				fmt.Printf("Failed to copy photo: %v\n", err)
			} else {
				fmt.Printf("Photo saved: %s\n", savePath)
			}
			dst.Close()
		}
	}() // Langsung panggil fungsi anonim
}
