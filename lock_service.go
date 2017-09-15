package ddbsync

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type LockServicer interface {
	NewLock(string, int64, time.Duration) sync.Locker
}

type LockService struct {
	db DBer
}

var _ LockServicer = (*LockService)(nil) // Forces compile time checking of the interface

func NewLockService(tableName string, region string, endpoint string, disableSSL bool) *LockService {
	return &LockService{
		db: NewDatabase(tableName, region, endpoint, disableSSL),
	}
}

func NewLockServiceFromDDBAPI(ddb dynamodbiface.DynamoDBAPI, tableName string) *LockService {
	return &LockService{
		db: &database{
			client:    ddb,
			tableName: tableName,
		},
	}
}

// Create a new Lock/Mutex with a particular key and timeout
func (l *LockService) NewLock(name string, ttl int64, lockReattemptWait time.Duration) sync.Locker {
	return NewMutex(name, ttl, l.db, lockReattemptWait)
}
