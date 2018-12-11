// Copyright 2018 Wolk Inc.
// This file is part of the SMT library.
//
// The SMT library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The SMT library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the plasmacash library. If not, see <http://www.gnu.org/licenses/>.
package smt

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

const DefaultChunkstorePath = "/tmp/cloudstore"

type Cloudstore struct {
	ldb      *leveldb.DB
	filepath string
	// providers
}

func NewCloudstore(path string) (self Cloudstore, err error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return self, fmt.Errorf("NewCloudstore:OpenFile %s: %v\n", path, err)
	}

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

//func (self Cloudstore) GetChunk(k []byte) (v []byte, err error) {
func (self Cloudstore) GetChunk(k []byte) (chunk []byte, ok bool, err error) {

	// Layer 2 interfaces with providers like ...
	//	v, err = self.providers[0].GetChunk(k)
	if self.ldb == nil {
		return []byte(""), false, fmt.Errorf("cache unavailable")
	}
	v, err := self.ldb.Get(k, nil)
	if err == leveldb.ErrNotFound {
		return v, false, nil
	} else if err != nil {
		return v, false, err
	}
	return v, true, nil
}

func (self Cloudstore) SetChunk(k []byte, v []byte) (err error) {
	// Layer 2 interfaces with providers like ...
	//  err = self.providers[0].StoreChunk(k, v)
	if self.ldb == nil {
		return fmt.Errorf("cache unavailable")
	}
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
