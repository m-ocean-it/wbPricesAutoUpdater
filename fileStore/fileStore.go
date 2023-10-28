package fileStore

import (
	"fmt"
	"os"
)

type FileStore struct {
	filePath string
	jobQueue jobQueue
}

func (fs *FileStore) popJob() job {
	var j job
	j, fs.jobQueue = fs.jobQueue[0], fs.jobQueue[1:]

	return j
}

func (fs *FileStore) pushJob(j job) {
	fs.jobQueue = append(fs.jobQueue, j)
}

func (fs *FileStore) processJob(j job) {
	response := jobResponse{}

	switch j.jobType {
	case "read":
		{
			fileContent, err := os.ReadFile(fs.filePath)
			if err != nil {
				response.err = err
				break
			}

			fileStat, err := os.Stat(fs.filePath)
			if err != nil {
				response.err = err
				break
			}

			response.fileContent = fileContent
			response.fileStat = fileStat
		}
	case "write":
		{
			jsonBytes := j.requestContent

			file, err := os.Create(fs.filePath) // Create or open the file for writing (truncating if it already exists)
			if err != nil {
				response.err = err
				break
			}
			defer file.Close()

			_, err = file.Write(jsonBytes)
			if err != nil {
				err = fmt.Errorf("error writing to file: %w", err)
				response.err = err
				break
			}
		}
	default:
		panic("")
	}

	j.responseChannel <- response
}

type jobQueue []job

type job struct {
	jobType         string
	requestContent  []byte
	responseChannel chan jobResponse // best when when capacity is 1
}

type jobResponse struct {
	fileContent []byte
	fileStat    os.FileInfo
	err         error
}
