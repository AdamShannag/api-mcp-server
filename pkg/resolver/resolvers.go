package resolver

import (
	"context"
	"fmt"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type StringResolver struct{}

func (r *StringResolver) Resolve(_ context.Context, req CallToolRequest, arg types.Arg) (string, error) {
	if arg.Required {
		return req.RequireString(arg.Name)
	}
	def, ok := arg.DefaultValue.(string)
	if !ok {
		def = ""
	}
	return req.GetString(arg.Name, def), nil
}

func (r *StringResolver) ToToolOption(arg types.Arg) mcp.ToolOption {
	opts := propertyOptions(arg)
	return mcp.WithString(arg.Name, opts...)
}

type IntResolver struct{}

func (r *IntResolver) Resolve(_ context.Context, req CallToolRequest, arg types.Arg) (string, error) {
	var val int
	var err error

	if arg.Required {
		val, err = req.RequireInt(arg.Name)
	} else {
		def, ok := arg.DefaultValue.(int)
		if !ok {
			def = 0
		}
		val = req.GetInt(arg.Name, def)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", val), nil
}

func (r *IntResolver) ToToolOption(arg types.Arg) mcp.ToolOption {
	opts := propertyOptions(arg)
	return mcp.WithNumber(arg.Name, opts...)
}

type FloatResolver struct{}

func (r *FloatResolver) Resolve(_ context.Context, req CallToolRequest, arg types.Arg) (string, error) {
	var val float64
	var err error

	if arg.Required {
		val, err = req.RequireFloat(arg.Name)
	} else {
		def, ok := arg.DefaultValue.(float64)
		if !ok {
			def = 0.0
		}
		val = req.GetFloat(arg.Name, def)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%f", val), nil
}

func (r *FloatResolver) ToToolOption(arg types.Arg) mcp.ToolOption {
	opts := propertyOptions(arg)
	return mcp.WithNumber(arg.Name, opts...)
}

type BoolResolver struct{}

func (r *BoolResolver) Resolve(_ context.Context, req CallToolRequest, arg types.Arg) (string, error) {
	var val bool
	var err error

	if arg.Required {
		val, err = req.RequireBool(arg.Name)
	} else {
		def, ok := arg.DefaultValue.(bool)
		if !ok {
			def = false
		}
		val = req.GetBool(arg.Name, def)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%t", val), nil
}

func (r *BoolResolver) ToToolOption(arg types.Arg) mcp.ToolOption {
	opts := propertyOptions(arg)
	return mcp.WithBoolean(arg.Name, opts...)
}

func propertyOptions(arg types.Arg) []mcp.PropertyOption {
	var opts []mcp.PropertyOption
	if arg.Description != "" {
		opts = append(opts, mcp.Description(arg.Description))
	}
	if arg.Required {
		opts = append(opts, mcp.Required())
	}
	return opts
}
