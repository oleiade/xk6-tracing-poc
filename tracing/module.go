package tracing

import (
	"fmt"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/modules"
)

type (
	// RootModule is the global module instance that will create Client
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		vu modules.VU

		*Tracing
	}
)

// Ensure the interfaces are implemented correctly
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface and returns
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	vu.Runtime().SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	return &ModuleInstance{
		vu: vu,
		Tracing: &Tracing{
			vu:   vu,
			HTTP: &HTTP{},
		},
	}
}

// Exports implements the modules.Instance interface and returns
// the exports of the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{Named: map[string]interface{}{
		"tracing":        mi.Tracing,
		"instrumentHTTP": mi.Tracing.InstrumentHTTP,
		"get":            func(arg string) { fmt.Println("hello: " + arg) },
	}}
}
