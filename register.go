package tracing

import (
	"github.com/grafana/xk6-tracing/tracing"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/tracing", new(tracing.RootModule))
}
