package detect

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

type Resolver interface {
	Resolve(context.Context, engineimage.Upstream) (engineimage.ResolvedInput, error)
}

type ResolverFunc func(context.Context, engineimage.Upstream) (engineimage.ResolvedInput, error)

func (f ResolverFunc) Resolve(ctx context.Context, u engineimage.Upstream) (engineimage.ResolvedInput, error) {
	return f(ctx, u)
}

type ResolverSet map[string]Resolver

func DefaultResolvers(root string) ResolverSet {
	static := ResolverFunc(func(_ context.Context, u engineimage.Upstream) (engineimage.ResolvedInput, error) {
		value := u.Version
		if value == "" {
			value = u.Ref
		}
		if value == "" {
			return engineimage.ResolvedInput{}, fmt.Errorf("%s upstream %q requires version or ref", u.Type, u.Name)
		}
		return engineimage.ResolvedInput{Name: u.Name, Type: u.Type, Value: value, Source: sourceOf(u)}, nil
	})
	return ResolverSet{
		"static-version":   static,
		"github-release":   static,
		"github-tag":       static,
		"github-commit":    static,
		"container-digest": ResolverFunc(resolveContainerDigest),
		"pypi-version":     static,
		"local-file":       ResolverFunc(resolveLocalFile(root)),
	}
}

func (s ResolverSet) Resolve(ctx context.Context, u engineimage.Upstream) (engineimage.ResolvedInput, error) {
	r, ok := s[u.Type]
	if !ok {
		return engineimage.ResolvedInput{}, fmt.Errorf("unsupported upstream type %q", u.Type)
	}
	return r.Resolve(ctx, u)
}

func resolveContainerDigest(_ context.Context, u engineimage.Upstream) (engineimage.ResolvedInput, error) {
	if u.Image == "" || u.Digest == "" {
		return engineimage.ResolvedInput{}, fmt.Errorf("container-digest upstream %q requires image and digest", u.Name)
	}
	return engineimage.ResolvedInput{Name: u.Name, Type: u.Type, Value: u.Image, Digest: u.Digest, Source: u.Image}, nil
}

func resolveLocalFile(root string) ResolverFunc {
	return func(_ context.Context, u engineimage.Upstream) (engineimage.ResolvedInput, error) {
		if u.Path == "" {
			return engineimage.ResolvedInput{}, fmt.Errorf("local-file upstream %q requires path", u.Name)
		}
		path := u.Path
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return engineimage.ResolvedInput{}, err
		}
		sum := sha256.Sum256(b)
		return engineimage.ResolvedInput{Name: u.Name, Type: u.Type, Value: u.Path, Digest: fmt.Sprintf("sha256:%x", sum), Source: u.Path}, nil
	}
}

func sourceOf(u engineimage.Upstream) string {
	if u.Repository != "" {
		return u.Repository
	}
	if u.Package != "" {
		return u.Package
	}
	return u.Name
}
