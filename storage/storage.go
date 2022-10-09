package storage

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName string
	maxKeys    int
	region     string
)

func init() {
	flag.StringVar(&bucketName, "bucket", "nftwswap", "")
	flag.IntVar(&maxKeys, "max-keys", 3, "")
	flag.StringVar(&region, "region", "ap-southeast-1", "")
}

func ListFiles() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Printf("failed to load SDK config, %v", err)
	}
	client := s3.NewFromConfig(cfg)

	params := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}
	pages := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		if v := int32(maxKeys); v != 0 {
			o.Limit = v
		}
	})
	var i int
	for pages.HasMorePages() {
		i++

		page, err := pages.NextPage(context.TODO())
		if err != nil {
			log.Printf("failed to get page, %v", err)
		}

		for _, obj := range page.Contents {
			fmt.Println("Object:", *obj.Key)
		}
	}
}

func UploadFile(path string, dest string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Printf("failed to load SDK config, %v", err)
	}
	client := s3.NewFromConfig(cfg)

	// read file
	file, err := os.Open(path)
	if err != nil {
		log.Println("Unable to open file" + path)
	}

	// create put object input
	filename := filepath.Join(dest, filepath.Base(path))
	input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &filename,
		Body:   file,
	}

	result, err := client.PutObject(context.TODO(), input)
	if err != nil {
		log.Printf("Got error uploading file %v", err)
	}
	_ = result
}
