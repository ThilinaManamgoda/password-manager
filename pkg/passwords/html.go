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
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	"github.com/pkg/errors"
	"html/template"
	"os"
	"strings"
)

const (
	//ConfKeyHTMLFilePath represents the HTML file path key in the conf map.
	ConfKeyHTMLFilePath = "CONF_KEY_HTML_PATH"
	//HTMLExporterID is the HTML exporter ID.
	HTMLExporterID = "HTML_EXPORTER"
	// HTMLTemplate represents the HTML template for password database.
    HTMLTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Password Manager</title>
	<style>
	table, th, td {
		border: 1px solid black;
		text-align: center;
	}
	</style>
</head>
 <body>
	 <h1>Password Manager database</h1>
	<table>
	 <tr>
		<th>ID</th>
		<th>Username</th>
		<th>Password</th>
		<th>Description</th>
		<th>Labels</th>
	  </tr>
	{{range $val := .}}
		<tr>
			<td>{{$val.ID}}</td>
			<td>{{$val.Username}}</td>
			<td>{{$val.Password}}</td>
			<td>{{$val.Description}}</td>
			<td>{{labels $val.Labels ","}}</td>
	 	 </tr>
	{{end}}
	</table>
 </body>
</html>
`
)
// HTMLExporter struct is used to export password entries.
type HTMLExporter struct {
	path string
}

// Init initialize the HTML Exporter.
func (e *HTMLExporter) Init(conf map[string]string) {
	e.path = conf[ConfKeyHTMLFilePath]
}


// Export exports password entries to a HTML file.
func (e *HTMLExporter) Export(entries []ImportExportEntry) error {
	exists, err := fileio.IsFileExists(e.path)
	if err != nil {
		return errors.Wrapf(err, "cannot inspect the given HTML file path")
	}

	if exists {
		return errors.New(fmt.Sprintf("Given HTML file: %s already exists", e.path))
	}

	htmlFile, err := os.Create(e.path)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("unable to create the HTML file: %s", e.path))
	}
	funcMap := template.FuncMap{
		"labels": strings.Join,
	}
	t := template.Must(template.New("html-tmpl").Funcs(funcMap).Parse(HTMLTemplate))
	err= t.Execute(htmlFile, entries)
	if err != nil {
		return errors.Wrap(err, "couldn't write to the HTML file")
	}
	return nil
}
