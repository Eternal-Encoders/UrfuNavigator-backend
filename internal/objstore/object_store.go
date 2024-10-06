package objstore

import "mime/multipart"

type ObjectStore interface {
	GetFile(fileName string) ([]byte, error)
	PostFile(file multipart.FileHeader) (string, string, error)
	DeleteFile(fileName string)
}
