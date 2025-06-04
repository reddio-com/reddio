package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3Config(bucket string) (*S3ConfigClient, error) {
	s := &S3ConfigClient{}
	if err := s.Init(bucket); err != nil {
		return nil, err
	}
	return s, nil
}

type S3ConfigClient struct {
	Bucket string
	client *s3.Client
	cd     *ConfigData
}

func (s *S3ConfigClient) Init(Bucket string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	s.client = client
	s.Bucket = Bucket
	s.cd = &ConfigData{}
	return nil
}

func (s *S3ConfigClient) GetConfig() *ConfigData {
	return s.cd
}

func (s *S3ConfigClient) LoadAllConfig() error {
	var err error
	s.cd.EvmCfg, err = s.LoadConfig("evm.toml")
	if err != nil {
		return err
	}
	s.cd.YuCfg, err = s.LoadConfig("yu.toml")
	if err != nil {
		return err
	}
	s.cd.PoaCfg, err = s.LoadConfig("poa.toml")
	if err != nil {
		return err
	}
	s.cd.ConfigCfg, err = s.LoadConfig("config.toml")
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
	return io.ReadAll(resp.Body)
}

type ConfigData struct {
	EvmCfg    []byte
	YuCfg     []byte
	PoaCfg    []byte
	ConfigCfg []byte
}
