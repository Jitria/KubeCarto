// SPDX-License-Identifier: Apache-2.0

package collector

import (
	"Operator/protobuf"
	"fmt"
	"io"
)

// 먼저 gRPC 서버 생성 후 서비스 등록 시에
// protobuf.RegisterOperatorServer(gRPCServer, ColH)
// 를 호출하여 ColHandler가 해당 RPC를 처리하도록 등록합니다.

// GiveAPILog implements the server-side streaming RPC for APILog.
func (ch *ColHandler) GiveAPILog(stream protobuf.Operator_GiveAPILogServer) error {
	for {
		// Receive APILog from stream.
		apiLog, err := stream.Recv()
		if err == io.EOF {
			// 클라이언트에서 전송 종료시 응답 전송
			return stream.SendAndClose(&protobuf.Response{Msg: 0})
		}
		if err != nil {
			return fmt.Errorf("GiveAPILog recv error: %v", err)
		}
		// 채널에 넣어 비동기로 처리 (ProcessAPILogs에서 처리)
		ch.apiLogChan <- apiLog
	}
}

// GiveEnvoyMetrics implements the server-side streaming RPC for EnvoyMetrics.
func (ch *ColHandler) GiveEnvoyMetrics(stream protobuf.Operator_GiveEnvoyMetricsServer) error {
	for {
		// Receive EnvoyMetrics from stream.
		envoyMetrics, err := stream.Recv()
		if err == io.EOF {
			// 전송 종료 시 응답 전송
			return stream.SendAndClose(&protobuf.Response{Msg: 0})
		}
		if err != nil {
			return fmt.Errorf("GiveEnvoyMetrics recv error: %v", err)
		}
		// metricsChan으로 보내 ProcessEnvoyMetrics에서 처리하도록 함
		ch.metricsChan <- envoyMetrics
	}
}
