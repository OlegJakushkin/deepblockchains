package smt

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb"
)

const DefaultChunkstorePath = "/tmp/cloudstore"

type Cloudstore struct {
	ldb *leveldb.DB
	// providers [4]ICloudstore
	filepath string
}

func NewCloudstore(path string) (self Cloudstore, err error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return self, fmt.Errorf("NewCloudstore:OpenFile %s: %v\n", path, err)
	}
	log.Info("NewCloudstore", "path", path)
	self = Cloudstore{
		ldb:      ldb,
		filepath: path,
	}
	/*
		Layer 2 interfaces with multiple cloud storage providers:
		self.providers[0], _ = cloud.NewCloudGoogleDatastore(cloud.DefaultGoogleProject)
		self.providers[1], _ = cloud.NewCloudAlibabaTablestore(cloud.DefaultAlibabaAccessKeyId, cloud.DefaultAlibabaAccessKeySecret)
		self.providers[2], _ = cloud.NewCloudMicrosoftAzure(cloud.DefaultMicrosoftAzureResourceString)
		self.providers[3], _ = cloud.NewCloudAmazonDynamo(cloud.DefaultAmazonRegion)
	*/
	return self, nil
}

func (self Cloudstore) RetrieveChunk(k []byte) (v []byte, err error) {
	// Layer 2 interfaces with providers like ...
	//	v, err = self.providers[0].RetrieveChunk(k)
	val, err := self.ldb.Get(k, nil)
	if err == leveldb.ErrNotFound {
		return val, leveldb.ErrNotFound
	} else if err != nil {
		return val, err
	}
	return val, nil
}

func (self Cloudstore) StoreChunk(k []byte, v []byte) (err error) {
	// Layer 2 interfaces with providers like ...
	//  err = self.providers[0].StoreChunk(k, v)
	err = self.ldb.Put(k, v, nil)
	if err != nil {
		return err
	}
	return nil
}

func (self Cloudstore) Close() (err error) {
	err = self.ldb.Close()
	if err != nil {
		return err
	}
	return nil
}
