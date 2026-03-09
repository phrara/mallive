package logging

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/phrara/mallive/common/tracing"
	"github.com/sirupsen/logrus"
)



func WhenGRPC(ctx context.Context, grpcServiceName string, args ...any) (
	func (any, error)) {
	fields := logrus.Fields{
		"grpcSvcName": grpcServiceName,
		Args: parseArgs(args...),
	}
	start := time.Now()
	traceId := tracing.TraceID(ctx)
	return func (resp any, err error)  {
		level, msg := logrus.InfoLevel, "grpc_success"
		fields[Cost] = time.Since(start)
		fields[Response] = resp
		fields["trace_id"] = traceId
		if err != nil {
			level, msg = logrus.ErrorLevel, "grpc_error"
			fields[Error] = err.Error()
		}
		logrus.WithContext(ctx).WithFields(fields).Logf(level, "%s", msg)
	}
}


func parseArgs(args ...any) string {
	var item []string
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			item = append(item, v)
		default:
			if v_str, ok := v.(interface{ String() string }); ok {
				item = append(item, v_str.String())
			} else {
				v_, err := json.Marshal(v)
				if err == nil {
					item = append(item, string(v_))
				}
			}
		}
		item = append(item, formatMySQLArg(arg))
	}
	return strings.Join(item, "||")
}