package tracing

import (
	"fmt"

	"go.k6.io/k6/js/modules"
)

type Tracing struct {
	vu modules.VU

	*HTTP
}

type HTTP struct{}

func (h *HTTP) Get(url string, params map[string]string) {
	fmt.Printf("Getting %s", url)
}

func (t *Tracing) InstrumentHTTP() {
	// t.vu.Runtime().RunString(`var hello = 'world'`)
	// t.vu.Runtime().RunString(`
	// 	function intrumentedHTTPGet(url, params) {
	// 		console.log(Object.keys(global))
	// 	}
	// `)

	// t.vu.Runtime().Set("http", &HTTP{})
	// t.vu.Runtime().GlobalObject().Set("http", &HTTP{})
	fmt.Println("HERE")
	fmt.Println(t.vu.Runtime().GlobalObject().Get("http"))
	fmt.Println(t.vu.Runtime().GlobalObject().Get("get"))
	fmt.Println(t.vu.Runtime().Get("http"))
	fmt.Println(t.vu.Runtime().Get("get"))
	t.vu.Runtime().GlobalObject().Set("get", func(arg string) { fmt.Println("bonjourno: " + arg) })

}
