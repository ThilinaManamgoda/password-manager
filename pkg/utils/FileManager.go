/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */

package utils

import (
	"io/ioutil"
	"os"
)

type PasswordFile struct {
	File string
}

func (p *PasswordFile) GetPasswords() ([]byte, error) {
	f, err := os.OpenFile(p.File, os.O_CREATE|os.O_RDONLY, 0640)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *PasswordFile) StorePasswords(data []byte) error {
	err := ioutil.WriteFile(p.File, data,0640 )
	if err != nil {
		return err
	}
	return nil
}
