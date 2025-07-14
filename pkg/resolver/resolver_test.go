package resolver_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AdamShannag/api-mcp-server/pkg/resolver"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/stretchr/testify/assert"
)

// mockCallToolRequest implements resolver.CallToolRequest interface for testing.
type mockCallToolRequest struct {
	stringVals map[string]string
	intVals    map[string]int
	floatVals  map[string]float64
	boolVals   map[string]bool

	requireStringErr map[string]error
	requireIntErr    map[string]error
	requireFloatErr  map[string]error
	requireBoolErr   map[string]error
}

func newMockCallToolRequest() *mockCallToolRequest {
	return &mockCallToolRequest{
		stringVals:       make(map[string]string),
		intVals:          make(map[string]int),
		floatVals:        make(map[string]float64),
		boolVals:         make(map[string]bool),
		requireStringErr: make(map[string]error),
		requireIntErr:    make(map[string]error),
		requireFloatErr:  make(map[string]error),
		requireBoolErr:   make(map[string]error),
	}
}

func (m *mockCallToolRequest) RequireString(name string) (string, error) {
	if err, ok := m.requireStringErr[name]; ok {
		return "", err
	}
	val, ok := m.stringVals[name]
	if !ok {
		return "", errors.New("missing required string")
	}
	return val, nil
}

func (m *mockCallToolRequest) GetString(name, def string) string {
	val, ok := m.stringVals[name]
	if !ok {
		return def
	}
	return val
}

func (m *mockCallToolRequest) RequireInt(name string) (int, error) {
	if err, ok := m.requireIntErr[name]; ok {
		return 0, err
	}
	val, ok := m.intVals[name]
	if !ok {
		return 0, errors.New("missing required int")
	}
	return val, nil
}

func (m *mockCallToolRequest) GetInt(name string, def int) int {
	val, ok := m.intVals[name]
	if !ok {
		return def
	}
	return val
}

func (m *mockCallToolRequest) RequireFloat(name string) (float64, error) {
	if err, ok := m.requireFloatErr[name]; ok {
		return 0, err
	}
	val, ok := m.floatVals[name]
	if !ok {
		return 0, errors.New("missing required float")
	}
	return val, nil
}

func (m *mockCallToolRequest) GetFloat(name string, def float64) float64 {
	val, ok := m.floatVals[name]
	if !ok {
		return def
	}
	return val
}

func (m *mockCallToolRequest) RequireBool(name string) (bool, error) {
	if err, ok := m.requireBoolErr[name]; ok {
		return false, err
	}
	val, ok := m.boolVals[name]
	if !ok {
		return false, errors.New("missing required bool")
	}
	return val, nil
}

func (m *mockCallToolRequest) GetBool(name string, def bool) bool {
	val, ok := m.boolVals[name]
	if !ok {
		return def
	}
	return val
}

func TestStringResolver(t *testing.T) {
	r := &resolver.StringResolver{}

	ctx := context.Background()
	mreq := newMockCallToolRequest()
	mreq.stringVals["foo"] = "bar"

	arg := types.Arg{Name: "foo", Type: "string", Required: true}
	val, err := r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "bar", val)

	mreq.requireStringErr["foo"] = errors.New("required error")
	_, err = r.Resolve(ctx, mreq, arg)
	assert.Error(t, err)
	delete(mreq.requireStringErr, "foo")

	arg.Required = false
	arg.DefaultValue = "default"
	mreq.stringVals = map[string]string{}
	val, err = r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "default", val)

	arg.DefaultValue = nil
	val, err = r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}

func TestIntResolver(t *testing.T) {
	r := &resolver.IntResolver{}

	ctx := context.Background()
	mreq := newMockCallToolRequest()
	mreq.intVals["foo"] = 42

	arg := types.Arg{Name: "foo", Type: "int", Required: true}
	val, err := r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "42", val)

	mreq.requireIntErr["foo"] = errors.New("required error")
	_, err = r.Resolve(ctx, mreq, arg)
	assert.Error(t, err)
	delete(mreq.requireIntErr, "foo")

	arg.Required = false
	arg.DefaultValue = 123
	mreq.intVals = map[string]int{}
	val, err = r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "123", val)

	arg.DefaultValue = nil
	val, err = r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "0", val)
}

func TestFloatResolver(t *testing.T) {
	r := &resolver.FloatResolver{}

	ctx := context.Background()
	mreq := newMockCallToolRequest()
	mreq.floatVals["foo"] = 3.14

	arg := types.Arg{Name: "foo", Type: "float", Required: true}
	val, err := r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "3.140000", val)

	mreq.requireFloatErr["foo"] = errors.New("required error")
	_, err = r.Resolve(ctx, mreq, arg)
	assert.Error(t, err)
	delete(mreq.requireFloatErr, "foo")

	arg.Required = false
	arg.DefaultValue = 2.718
	mreq.floatVals = map[string]float64{}
	val, err = r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "2.718000", val)

	arg.DefaultValue = nil
	val, err = r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "0.000000", val)
}

func TestBoolResolver(t *testing.T) {
	r := &resolver.BoolResolver{}

	ctx := context.Background()
	mreq := newMockCallToolRequest()
	mreq.boolVals["foo"] = true

	arg := types.Arg{Name: "foo", Type: "bool", Required: true}
	val, err := r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "true", val)

	mreq.requireBoolErr["foo"] = errors.New("required error")
	_, err = r.Resolve(ctx, mreq, arg)
	assert.Error(t, err)
	delete(mreq.requireBoolErr, "foo")

	arg.Required = false
	arg.DefaultValue = true
	mreq.boolVals = map[string]bool{}
	val, err = r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "true", val)

	arg.DefaultValue = nil
	val, err = r.Resolve(ctx, mreq, arg)
	assert.NoError(t, err)
	assert.Equal(t, "false", val)
}

func TestTypeResolverRegistry(t *testing.T) {
	r := resolver.NewDefaultTypeResolverRegistry()
	ctx := context.Background()

	for _, arg := range []types.Arg{
		{Name: "str", Type: "string"},
		{Name: "int", Type: "int"},
		{Name: "flt", Type: "float"},
		{Name: "bl", Type: "bool"},
	} {
		_, err := r.Resolve(ctx, newMockCallToolRequest(), arg)
		assert.NoError(t, err)
	}

	arg := types.Arg{Name: "unknown", Type: "unknown"}
	val, err := r.Resolve(ctx, newMockCallToolRequest(), arg)
	assert.NoError(t, err)
	assert.Equal(t, "", val)

	opt := r.ToToolOption(types.Arg{Name: "foo", Type: "string", Description: "desc", Required: true})
	assert.NotNil(t, opt)

	opt = r.ToToolOption(types.Arg{Name: "foo", Type: "unknown"})
	assert.NotNil(t, opt)
}
