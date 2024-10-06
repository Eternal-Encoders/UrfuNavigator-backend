package objstore

import (
	"context"
	"io"
	"log"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOS3 struct {
	Client     *minio.Client
	BucketName string
}

func Connect(endpoint string, access string, secret string, bucketName string) *MinIOS3 {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(access, secret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &MinIOS3{
		Client:     minioClient,
		BucketName: bucketName,
	}
}

func (s *MinIOS3) GetFile(fileName string) ([]byte, error) {
	option := minio.GetObjectOptions{}

	file, err := s.Client.GetObject(context.TODO(), s.BucketName, fileName, option)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := make([]byte, 8)
	res := []byte{}
	for {
		i, err := file.Read(buf)
		res = append(res, buf[:i]...)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return res, nil
		}
	}
}

func (s *MinIOS3) PostFile(file multipart.FileHeader) (string, string, error) {
	option := minio.PutObjectOptions{ContentType: file.Header["Content-Type"][0]}
	fileName := file.Filename

	buff, err := file.Open()
	if err != nil {
		log.Fatalln(err)
	}

	defer buff.Close()

	info, err := s.Client.PutObject(context.Background(), s.BucketName, fileName, buff, file.Size, option)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Uploaded", fileName, " of size: ", info.Size, "Successfully.")

	return fileName, fileName, err
}

func (s *MinIOS3) DeleteFile(fileName string) {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	err := s.Client.RemoveObject(context.Background(), s.BucketName, fileName, opts)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Success")
}
