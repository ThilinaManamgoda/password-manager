/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package passwords

import (
	"encoding/csv"
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

// CSVSeparator represents the separator for CSV file.
const CSVSeparator = ","

// ImportFromCSV imports the passwords entries from the given CSV file.
func (p *Repository) ImportFromCSV(csvFilePath string) error {
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return errors.Wrap(err, "Couldn't open the csv file")
	}
	r := csv.NewReader(csvFile)

	first := true
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if first {
			first = false
			continue
		}
		id := record[0]
		uN := record[1]
		password := record[2]
		labels := strings.Split(record[3], CSVSeparator)

		err = p.addPasswordEntryToRepo(id, uN, password, labels)
		if err != nil {
			return err
		}
	}
	err = p.saveDB()
	if err != nil {
		return err
	}
	return nil
}

// ExportToCSV exports the password database to the given CSV file.
func (p *Repository) ExportToCSV(csvFilePath string) error {
	exists, err := fileio.IsFileExists(csvFilePath)
	if err != nil {
		return errors.Wrapf(err, "cannot inspect the given CSV file path")
	}

	if exists {
		return errors.New("CSV file already exists")
	}

	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		return errors.Wrapf(err, "unable to create the CSV file")
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	firstLine := true
	for id, entry := range p.db.Entries {
		var csvEntry []string
		if firstLine {
			csvEntry = []string{"id", "username", "password", "labels"}
			firstLine = false
		} else {
			labels, err := p.searchLabelsForID(id)
			if err != nil {
				return errors.Wrapf(err, "couldn't find the labels for id %s", id)
			}
			csvEntry = []string{id, entry.Username, entry.Password, strings.Join(labels, CSVSeparator)}
		}
		err = writer.Write(csvEntry)
		if err != nil {
			return errors.Wrap(err, "couldn't write to CSV file")
		}
	}
	return nil
}
