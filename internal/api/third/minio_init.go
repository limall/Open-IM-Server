package apiThird

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/policy"
	url2 "net/url"
)

var (
	minioClient *minio.Client
)

func MinioInit() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "minio config: ", config.Config.Credential.Minio)
	var initUrl string
	if config.Config.Credential.Minio.EndpointInnerEnable {
		initUrl = config.Config.Credential.Minio.EndpointInner
	} else {
		initUrl = config.Config.Credential.Minio.Endpoint
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "use initUrl: ", initUrl)
	minioUrl, err := url2.Parse(initUrl)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "parse failed, please check config/config.yaml", err.Error())
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "Parse ok ", config.Config.Credential.Minio)
	minioClient, err = minio.New(minioUrl.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Config.Credential.Minio.AccessKeyID, config.Config.Credential.Minio.SecretAccessKey, ""),
		Secure: false,
	})
	log.NewInfo(operationID, utils.GetSelfFuncName(), "new ok ", config.Config.Credential.Minio)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "init minio client failed", err.Error())
		return
	}
	opt := minio.MakeBucketOptions{
		Region:        config.Config.Credential.Minio.Location,
		ObjectLocking: false,
	}
	err = minioClient.MakeBucket(context.Background(), config.Config.Credential.Minio.Bucket, opt)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "MakeBucket failed ", err.Error())
		exists, err := minioClient.BucketExists(context.Background(), config.Config.Credential.Minio.Bucket)
		if err == nil && exists {
			log.NewWarn(operationID, utils.GetSelfFuncName(), "We already own ", config.Config.Credential.Minio.Bucket)
		} else {
			if err != nil {
				log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
			}
			log.NewError(operationID, utils.GetSelfFuncName(), "create bucket failed and bucket not exists")
			return
		}
	}
	// 自动化桶public的代码
	err = minioClient.SetBucketPolicy(context.Background(), config.Config.Credential.Minio.Bucket, policy.BucketPolicyReadWrite)
	if err != nil {
		log.NewDebug("", utils.GetSelfFuncName(), "SetBucketPolicy failed please set in web", err.Error())
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "minio create and set policy success")
}
