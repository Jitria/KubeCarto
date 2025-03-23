// SPDX-License-Identifier: Apache-2.0

package processor

import (
	"log"
	"sync"

	"github.com/Jitria/SentryFlow/protobuf"

	"Agent/uploader"
)

// InsertAPILog Function
func InsertAPILog(data interface{}) {
	LogH.apiLogChan <- data
}

// processAPILogs Function
func processAPILogs(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case logType, ok := <-LogH.apiLogChan:
			if !ok {
				log.Print("[LogProcessor] Failed to process an API log")
			}

			go uploader.UploadAPILog(logType.(*protobuf.APILog))

		case <-LogH.stopChan:
			wg.Done()
			return
		}
	}
}

// == //
