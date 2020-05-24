// Copyright © 2019 Thilina Manamgoda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package passwords

import (
	"encoding/csv"
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

const (
	//CSVExporterID is the CSV exporter ID.
	CSVExporterID = "CSV_EXPORTER"
	//CSVImporterID is the CSV importer ID.
	CSVImporterID = "CSV_IMPORTER"
	// CSVSeparator represents the separator for CSV file.
	CSVSeparator = ","
	//ConfKeyCSVFilePath represents the CSV path key in the conf map.
	ConfKeyCSVFilePath = "CONF_KEY_CSV_PATH"
)

// CSVExporter struct is used to export password entries.
type CSVExporter struct {
	path string
}

// CSVImporter struct is used to import password entries.
type CSVImporter struct {
	path string
}

// Init initialize the CSV Importer.
func (i *CSVImporter) Init(conf map[string]string) {
	i.path = conf[ConfKeyCSVFilePath]
}

// Init initialize the CSV Exporter.
func (e *CSVExporter) Init(conf map[string]string) {
	e.path = conf[ConfKeyCSVFilePath]
}

// Import imports password entries from a CSV file.
func (i *CSVImporter) Import() ([]ImportExportEntry, error) {
	csvFile, err := os.Open(i.path)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't open the CSV file")
	}
	r := csv.NewReader(csvFile)

	var entries []ImportExportEntry
	first := true
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "unable to read CSV file")
		}
		if first {
			first = false
			continue
		}
		id := record[0]
		uN := record[1]
		password := record[2]
		desc := record[3]
		labels := strings.Split(record[4], CSVSeparator)
		entries = append(entries, ImportExportEntry{
			Entry: Entry{ID: id,
				Username:    uN,
				Password:    password,
				Description: desc,
			},
			Labels: labels,
		})
	}
	return entries, nil
}

// Export exports password entries to a CSV file.
func (e *CSVExporter) Export(entries []ImportExportEntry) error {
	exists, err := fileio.IsFileExists(e.path)
	if err != nil {
		return errors.Wrapf(err, "cannot inspect the given CSV file path")
	}

	if exists {
		return errors.New("CSV file already exists")
	}

	csvFile, err := os.Create(e.path)
	if err != nil {
		return errors.Wrapf(err, "unable to create the CSV file")
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	firstLine := true
	for _, entry := range entries {
		var csvEntry []string
		if firstLine {
			csvEntry = []string{"id", "username", "password", "description", "labels"}
			firstLine = false
		} else {
			csvEntry = []string{entry.ID, entry.Username, entry.Password, entry.Description, strings.Join(entry.Labels, CSVSeparator)}
		}
		err = writer.Write(csvEntry)
		if err != nil {
			return errors.Wrap(err, "couldn't write to CSV file")
		}
	}
	return nil
}
