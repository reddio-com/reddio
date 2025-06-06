package s3

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3Config(folder, bucket string) (*S3ConfigClient, error) {
	s := &S3ConfigClient{}
	if err := s.Init(folder, bucket); err != nil {
		return nil, fmt.Errorf("init s3 client err: %v", err)
	}
	if err := s.LoadAllConfig(); err != nil {
		return nil, fmt.Errorf("load config from s3 err: %v", err)
	}
	return s, nil
}

type S3ConfigClient struct {
	Folder string
	Bucket string
	client *s3.Client
	cd     *ConfigData
}

func (s *S3ConfigClient) Init(Folder, Bucket string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	s.client = client
	s.Folder = Folder
	s.Bucket = Bucket
	s.cd = &ConfigData{}
	return nil
}

func (s *S3ConfigClient) GetConfig() *ConfigData {
	return s.cd
}

func (s *S3ConfigClient) LoadAllConfig() error {
	var err error
	s.cd.EvmCfg, err = s.LoadConfig(filepath.Join(s.Folder, "evm.toml"))
	if err != nil {
		return err
	}
	s.cd.YuCfg, err = s.LoadConfig(filepath.Join(s.Folder, "yu.toml"))
	if err != nil {
		return err
	}
	s.cd.PoaCfg, err = s.LoadConfig(filepath.Join(s.Folder, "poa.toml"))
	if err != nil {
		return err
	}
	s.cd.ConfigCfg, err = s.LoadConfig(filepath.Join(s.Folder, "config.toml"))
	if err != nil {
		return err
	}
	return nil
}

func (s *S3ConfigClient) LoadConfig(key string) ([]byte, error) {
	resp, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("load %v error %v", key, err)
	}
	return data, nil
}

type ConfigData struct {
	EvmCfg    []byte
	YuCfg     []byte
	PoaCfg    []byte
	ConfigCfg []byte
}
