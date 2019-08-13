/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */


package utils




func IsValidByteSlice(data []byte) bool {
	return (data != nil) && (len(data) != 0)
}

func IsPasswordValid(passphrase string) bool {
	return passphrase != ""
}