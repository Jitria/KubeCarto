// SPDX-License-Identifier: Apache-2.0

package uploader

// == //

// UplH global reference for Uploader Handler
var ColH *UplHandler

// init Function
func init() {
	ColH = NewUploaderHandler()
}

// UplHandler Structure
type UplHandler struct {
}

// NewUploaderHandler Function
func NewUploaderHandler() *UplHandler {
	ch := &UplHandler{}
	return ch
}

// == //

// StartUploader Function
func StartUploader() bool {

	return true
}

// StopUploader Function
func StopUploader() bool {

	return true
}

// == //

// trash
func InsertEnvoyMetrics(_ interface{}) {

}

func InsertAPILog(_ interface{}) {

}
