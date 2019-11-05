package verbsbucket

import (
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type VerbsBucket struct {
	svc        *s3.S3
	bucketName string
}

// New creates a new VerbsBucket
// *VerbsBucket as a type definition means "this is a pointer to the value"
// *something used as a value means "dereference" i.e. get the underlying value of a pointer
// &something means take the address in memory of something (& is never used in type definintions)
func New(name string) *VerbsBucket {
	return &VerbsBucket{
		bucketName: name,
		svc:        s3.New(session.New()),
	}
}

// GetFile gets the content of the specified file
func (vb VerbsBucket) GetFile(key string) ([]byte, error) {
	out, err := vb.svc.GetObject(&s3.GetObjectInput{
		Bucket: &vb.bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(out.Body)
	if err != nil {
		log.Println("Error reading S3 result body: ", err)
		return nil, err
	}

	return result, nil
}
