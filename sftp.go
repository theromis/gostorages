package gostorages

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"time"
)

type SFTPStorage struct {
	*BaseStorage
	*sftp.Client
}

type SFTPFile struct {
	*sftp.File
	Storage  Storage
	FileInfo os.FileInfo
}

func NewSFTPStorage(username string, password string, host string, port int, location string, baseURL string) (Storage, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)

	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(sshClient)

	if err != nil {
		return nil, err
	}

	return &SFTPStorage{NewBaseStorage(location, baseURL), client}, nil
}

func NewSFTPFile(storage Storage, file *sftp.File) (File, error) {
	fileInfo, err := file.Stat()

	if err != nil {
		return nil, err
	}

	return &SFTPFile{
		file,
		storage,
		fileInfo,
	}, nil
}

func (s *SFTPStorage) Open(filepath string) (File, error) {
	file, err := s.Open(s.Path(filepath))

	if err != nil {
		return nil, err
	}

	return NewSFTPFile(s, file)
}

func (f *SFTPFile) Size() int64 {
	return f.FileInfo.Size()
}

func (f *SFTPFile) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(f)
}

func (s *SFTPStorage) Delete(filepath string) error {
	return s.Remove(s.Path(filepath))
}

func (s *SFTPStorage) Exists(filepath string) bool {
	_, err := s.Lstat(s.Path(filepath))
	return err == nil
}

func (s *SFTPStorage) ModifiedTime(filepath string) (time.Time, error) {
	fi, err := s.Lstat(s.Path(filepath))
	if err != nil {
		return time.Time{}, err
	}

	return fi.ModTime(), nil
}
