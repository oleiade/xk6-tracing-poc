package tracing

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/modulestest"
)

func TestTracingGetOrCreateParams(t *testing.T) {
	t.Parallel()

	t.Run("a provided body and params leaves params untouched", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		body := testSetup.VU.Runtime().NewObject()
		params := testSetup.VU.Runtime().NewObject()
		headers := testSetup.VU.Runtime().NewObject()
		err := headers.Set("Content-Type", "application/json")
		require.NoError(t, err)
		err = params.Set("headers", headers)
		require.NoError(t, err)

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPPostMethodName, body, params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 2)
		assert.Equal(t, body, gotArgs[0])
		assert.Equal(t, params, gotArgs[1])
	})

	t.Run("a provided body and null params arg intializes params", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		body := testSetup.VU.Runtime().NewObject()

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPPostMethodName, body, goja.Null())

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 2)
		assert.Equal(t, body, gotArgs[0])
		assert.Equal(t, gotParams, gotArgs[1])
	})

	t.Run("a provided body and undefined params arg intializes params", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		body := testSetup.VU.Runtime().NewObject()

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPPostMethodName, body, goja.Undefined())

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 2)
		assert.Equal(t, body, gotArgs[0])
		assert.Equal(t, gotParams, gotArgs[1])
	})

	t.Run("a provided body and nil params arg intializes params", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		body := testSetup.VU.Runtime().NewObject()

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPPostMethodName, body, nil)

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 2)
		assert.Equal(t, body, gotArgs[0])
		assert.Equal(t, gotParams, gotArgs[1])
	})

	t.Run("a provided body and no params arg intializes params", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		body := testSetup.VU.Runtime().NewObject()

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPPostMethodName, body)

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 2)
		assert.Equal(t, body, gotArgs[0])
		assert.Equal(t, gotParams, gotArgs[1])
	})

	t.Run("a provided nil body and no params arg intializes params", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPPostMethodName, nil)

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 2)
		assert.Nil(t, gotArgs[0])
		assert.Equal(t, gotParams, gotArgs[1])
	})

	t.Run("a provided null params argument should initialize it", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		params := goja.Null()

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPGetMethodName, params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 1)
		assert.Equal(t, gotParams, gotArgs[0])
	})

	t.Run("a provided undefined params argument should initialize it", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		params := goja.Undefined()

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPGetMethodName, params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 1)
		assert.Equal(t, gotParams, gotArgs[0])
	})

	t.Run("a provided nil params argument should initialize it", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPGetMethodName, nil)

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 1)
		assert.Equal(t, gotParams, gotArgs[0])
	})

	t.Run("a provided params argument should leave it untouched", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		headers := testSetup.VU.Runtime().NewObject()
		err := headers.Set("test-header", "test-value")
		require.NoError(t, err)
		params := testSetup.VU.Runtime().NewObject()
		err = params.Set("headers", headers)
		require.NoError(t, err)

		gotArgs, gotParams, gotErr := tracing.getOrCreateParams(k6HTTPGetMethodName, params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, gotParams)
		assert.NotNil(t, gotArgs)
		assert.Len(t, gotArgs, 1)
		assert.Equal(t, gotArgs[0], params)
		assert.True(t, gotParams == params)
	})

	t.Run("no arguments should initialize a params argument", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}

		// the http.get method has a single argument, besides url,
		// which is the optional params object. We use it here to
		// verify that we create a default params object under the
		// hood.
		args, params, gotErr := tracing.getOrCreateParams(k6HTTPGetMethodName)

		assert.NoError(t, gotErr)
		assert.NotNil(t, params)
		assert.NotNil(t, args)
		assert.Len(t, args, 2)
		assert.Equal(t, args[0], goja.Null())
		assert.Equal(t, args[1], params)
	})
}

func TestTracingGetOrCreateHeaders(t *testing.T) {
	t.Parallel()

	t.Run("params object with headers should return them", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		params := testSetup.VU.Runtime().NewObject()
		headers := testSetup.VU.Runtime().NewObject()
		err := params.Set("headers", headers)
		require.NoError(t, err)

		gotHeaders, gotErr := tracing.getOrCreateHeaders(params)

		assert.NoError(t, gotErr)
		assert.Equal(t, gotHeaders, params.Get("headers"))
	})

	t.Run("params object without headers should create one", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		params := testSetup.VU.Runtime().NewObject()

		gotHeaders, gotErr := tracing.getOrCreateHeaders(params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, params.Get("headers"))
		assert.NotNil(t, gotHeaders)
	})

	t.Run("params object with undefined headers should set one", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		params := testSetup.VU.Runtime().NewObject()
		err := params.Set("headers", goja.Undefined())
		require.NoError(t, err)

		gotHeaders, gotErr := tracing.getOrCreateHeaders(params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, params.Get("headers"))
		assert.NotNil(t, gotHeaders)
	})

	t.Run("params object with null headers should set one", func(t *testing.T) {
		t.Parallel()

		testSetup := modulestest.NewRuntime(t)
		tracing := &Tracing{vu: testSetup.VU}
		params := testSetup.VU.Runtime().NewObject()
		err := params.Set("headers", goja.Null())
		require.NoError(t, err)

		gotHeaders, gotErr := tracing.getOrCreateHeaders(params)

		assert.NoError(t, gotErr)
		assert.NotNil(t, params.Get("headers"))
		assert.NotNil(t, gotHeaders)
	})
}
