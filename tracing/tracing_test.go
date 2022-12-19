package tracing

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/modulestest"
)

func TestTracingNormalizeArgs(t *testing.T) {
	t.Parallel()

	t.Run("methods with a params arg left nil should initialize it", func(t *testing.T) {
		t.Parallel()

		rt := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: rt.VU}
		// wantParams := rt.VU.Runtime().NewObject()

		gotArgs, gotParams := tracing.normalizeArgs(k6HTTPGetMethodName, nil)

		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 1)
		// assert.Equal(t, gotArgs[0], wantParams)
		// assert.True(t, gotParams == wantParams)
	})

	t.Run("methods with a params argument should leave it untouched", func(t *testing.T) {
		t.Parallel()

		rt := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: rt.VU}
		headers := rt.VU.Runtime().NewObject()
		err := headers.Set("test-header", "test-value")
		require.NoError(t, err)
		params := rt.VU.Runtime().NewObject()
		err = params.Set("headers", headers)
		require.NoError(t, err)

		gotArgs, gotParams := tracing.normalizeArgs(k6HTTPGetMethodName, params)

		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 1)
		assert.Equal(t, gotArgs[0], params)
		assert.True(t, gotParams == params)
	})

	t.Run("no arg provided to a method with an optional argument should add an empty params", func(t *testing.T) {
		t.Parallel()

		rt := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: rt.VU}

		// the http.get method has a single argument, besides url,
		// which is the optional params object. We use it here to
		// verify that we create a default params object under the
		// hood.
		args, params := tracing.normalizeArgs(k6HTTPGetMethodName)

		assert.NotNil(t, params)
		assert.NotNil(t, args)
		assert.Len(t, args, 2)
		assert.Equal(t, args[0], goja.Null())
		assert.Equal(t, args[1], params)
	})
}

func TestTracingNormalizeParamsHeaders(t *testing.T) {
	t.Parallel()

	t.Run("params object with headers should return them", func(t *testing.T) {
		t.Parallel()

		rt := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: rt.VU}
		params := rt.VU.Runtime().NewObject()
		headers := rt.VU.Runtime().NewObject()
		err := params.Set("headers", headers)
		require.NoError(t, err)

		gotHeaders, gotErr := tracing.normalizeParamsHeaders(params)

		assert.NoError(t, gotErr)
		assert.Equal(t, gotHeaders, params.Get("headers"))
	})

	t.Run("params object without headers should create one", func(t *testing.T) {
		t.Parallel()

		rt := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: rt.VU}
		params := rt.VU.Runtime().NewObject()

		gotHeaders, gotErr := tracing.normalizeParamsHeaders(params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, params.Get("headers"))
		assert.NotNil(t, gotHeaders)
	})

	t.Run("params object with undefined headers should set one", func(t *testing.T) {
		t.Parallel()

		rt := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: rt.VU}
		params := rt.VU.Runtime().NewObject()
		err := params.Set("headers", goja.Undefined())
		require.NoError(t, err)

		gotHeaders, gotErr := tracing.normalizeParamsHeaders(params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, params.Get("headers"))
		assert.NotNil(t, gotHeaders)
	})

	t.Run("params object with null headers should set one", func(t *testing.T) {
		t.Parallel()

		rt := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: rt.VU}
		params := rt.VU.Runtime().NewObject()
		err := params.Set("headers", goja.Null())
		require.NoError(t, err)

		gotHeaders, gotErr := tracing.normalizeParamsHeaders(params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, params.Get("headers"))
		assert.NotNil(t, gotHeaders)
	})
}
