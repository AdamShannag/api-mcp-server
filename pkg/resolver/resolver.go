package resolver

import (
	"context"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type CallToolRequest interface {
	RequireString(name string) (string, error)
	GetString(name, def string) string
	RequireInt(name string) (int, error)
	GetInt(name string, def int) int
	RequireFloat(name string) (float64, error)
	GetFloat(name string, def float64) float64
	RequireBool(name string) (bool, error)
	GetBool(name string, def bool) bool
}

type ArgResolver interface {
	Resolve(ctx context.Context, req CallToolRequest, arg types.Arg) (string, error)
	ToToolOption(arg types.Arg) mcp.ToolOption
}

type TypeResolverRegistry struct {
	resolvers map[string]ArgResolver
}

func NewDefaultTypeResolverRegistry() *TypeResolverRegistry {
	r := &TypeResolverRegistry{
		resolvers: make(map[string]ArgResolver),
	}

	r.Register("string", &StringResolver{})
	r.Register("int", &IntResolver{})
	r.Register("float", &FloatResolver{})
	r.Register("bool", &BoolResolver{})

	return r
}

func (r *TypeResolverRegistry) Register(argType string, resolver ArgResolver) {
	r.resolvers[argType] = resolver
}

func (r *TypeResolverRegistry) Resolve(ctx context.Context, req CallToolRequest, arg types.Arg) (string, error) {
	resolver, ok := r.resolvers[arg.Type]
	if !ok {
		resolver = r.resolvers["string"]
	}
	return resolver.Resolve(ctx, req, arg)
}

func (r *TypeResolverRegistry) ToToolOption(arg types.Arg) mcp.ToolOption {
	if res, ok := r.resolvers[arg.Type]; ok {
		return res.ToToolOption(arg)
	}
	return r.resolvers["string"].ToToolOption(arg)
}
