package utils

import (
	"archive/zip"
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

type File struct {
	Mu        *sync.Mutex
	FileInfos []*FileInfo
}

type FileInfo struct {
	Filename string
	FileByte []byte
}

func ZipFilesWithBytes(needZipPkgName string, file *File) ([]byte, error) {
	file.Mu.Lock()
	defer file.Mu.Unlock()
	var err error
	var zipBuffer *bytes.Buffer = new(bytes.Buffer)
	var zipWriter *zip.Writer = zip.NewWriter(zipBuffer)
	var zipEntry io.Writer
	// Create entry in zip file
	for _, fileInfo := range file.FileInfos {
		zipEntry, err = zipWriter.Create(fileInfo.Filename)
		if err != nil {
			log.Errorf("ZipFilesWithBytes,Create error,err=%v", err)
			return nil, err
		}
		// Write content into zip entry
		_, err = zipEntry.Write(fileInfo.FileByte)
		if err != nil {
			log.Errorf("ZipFilesWithBytes,Write error,err=%v", err)
			return nil, err
		}
	}
	// Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		log.Errorf("ZipFilesWithBytes,Close error,err=%v", err)
		return nil, err
	}
	return zipBuffer.Bytes(), nil
	// Write the zip file to the disk
	//err = ioutil.WriteFile(needZipPkgName, zipBuffer.Bytes(), 0644)
	//if err != nil {
	//	return nil, err
	//}
	//defer os.Remove(needZipPkgName)
	//return ioutil.ReadFile(needZipPkgName)
}
