/*
 * Copyright (c) 2020 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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
	"github.com/pkg/errors"
)

// ImportExportEntry struct represents an import/export entry.
type ImportExportEntry struct {
	Entry
	Labels []string
}

// Exporter interface is used to export data.
type Exporter interface {
	// Init initialises the exporter.
	Init(conf map[string]string)
	// Export exports the importExport entries.
	Export(entries [] ImportExportEntry) error
}

// Importer interface is used to import data.
type Importer interface {
	// Init initialises the importer.
	Init(conf map[string]string)
	// Import imports the importExport entries.
	Import() ([] ImportExportEntry, error)
}

// ExporterFactory used to create an Exporter.
type ExporterFactory struct {
	ID string
}

// Exporter returns the requested an Exporter.
func (e *ExporterFactory) Exporter() Exporter {
	switch e.ID {
	case CSVExporterID:
		return &CSVExporter{}
	default:
		return &CSVExporter{}
	}
}

// ImporterFactory used to create an Importer.
type ImporterFactory struct {
	ID string
}

// Importer returns the requested an Importer.
func (e *ImporterFactory) Importer() Importer {
	switch e.ID {
	case CSVImporterID:
		return &CSVImporter{}
	default:
		return &CSVImporter{}
	}
}

// Import imports the password database.
func (p *Repository) Import(importerID string, conf map[string]string) error {
	var importer Importer
	if importerID == CSVImporterID {
		imFac := ImporterFactory{
			ID: CSVImporterID,
		}
		importer = imFac.Importer()
		importer.Init(map[string]string{ConfKeyCSVPath: conf[ConfKeyCSVPath]})
	} else {
		return errors.New("No supported import medium provided")
	}

	entries, err := importer.Import()
	if err != nil {
		return errors.Wrap(err, "unable to import from the CSV file")
	}

	for _, entry := range entries {
		err = p.addPasswordEntryToRepo(entry.ID, entry.Username, entry.Password, entry.Labels)
		if err != nil {
			return errors.Wrap(err, "cannot import passwords")
		}
	}

	err = p.saveDB()
	if err != nil {
		return err
	}
	return nil
}

// Export exports the password database.
func (p *Repository) Export(exporterID string, conf map[string]string) error {
	var exporter Exporter
	if exporterID == CSVExporterID {
		exFac := ExporterFactory{
			ID: CSVImporterID,
		}
		exporter = exFac.Exporter()
		exporter.Init(map[string]string{ConfKeyCSVPath: conf[ConfKeyCSVPath]})
	} else {
		return errors.New("No export medium provided")
	}
	var importExportEntries []ImportExportEntry
	for _, entry := range p.db.Entries {
		labels, err := p.searchLabelsForID(entry.ID)
		if err != nil {
			return errors.Wrapf(err, "couldn't find the labels for id %s", entry.ID)
		}
		importExportEntries = append(importExportEntries, ImportExportEntry{
			Entry:  Entry{ID: entry.ID, Username: entry.Username, Password: entry.Password},
			Labels: labels,
		})
	}
	err := exporter.Export(importExportEntries)
	if err != nil {
		return errors.Wrap(err, "cannot export password DB")
	}
	return nil
}
