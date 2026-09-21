package detect

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

type Options struct {
	Root      string
	Resolvers ResolverSet
}

func Detect(ctx context.Context, opts Options) (UpdatePlan, error) {
	root := opts.Root
	if root == "" {
		root = "."
	}
	resolvers := opts.Resolvers
	if resolvers == nil {
		resolvers = DefaultResolvers(root)
	}
	paths, err := discoverManifests(root)
	if err != nil {
		return UpdatePlan{}, err
	}
	plan := UpdatePlan{Schema: UpdatePlanSchema}
	for _, path := range paths {
		m, err := engineimage.LoadManifest(path)
		if err != nil {
			plan.Errors = append(plan.Errors, PlanError{Manifest: rel(root, path), Error: err.Error()})
			continue
		}
		resolved, err := resolveAll(ctx, resolvers, m.Upstreams)
		if err != nil {
			plan.Errors = append(plan.Errors, PlanError{Manifest: rel(root, path), Flavor: m.Name, Error: err.Error()})
			continue
		}
		lockPath := filepath.Join(filepath.Dir(path), engineimage.DefaultLockFilename)
		changed, reason, err := compareLock(lockPath, m.Name, resolved)
		if err != nil {
			plan.Errors = append(plan.Errors, PlanError{Manifest: rel(root, path), Flavor: m.Name, Error: err.Error()})
			continue
		}
		manifestRel := rel(root, path)
		plan.Flavors = append(plan.Flavors, FlavorPlan{
			Name:           m.Name,
			Engine:         m.Engine,
			Vendor:         m.Vendor,
			Manifest:       manifestRel,
			Dockerfile:     dockerfileRel(manifestRel, m.Build.Dockerfile),
			Image:          m.Image.Repository,
			Certification:  m.Publish.Certification,
			Changed:        changed,
			Reason:         reason,
			ResolvedInputs: resolved,
			Tests:          m.Tests,
		})
	}
	sort.Slice(plan.Flavors, func(i, j int) bool { return plan.Flavors[i].Name < plan.Flavors[j].Name })
	if len(plan.Errors) > 0 {
		return plan, fmt.Errorf("detector failed with %d error(s)", len(plan.Errors))
	}
	return plan, nil
}

func discoverManifests(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(filepath.Join(root, "engines"), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(path) == "image.yaml" {
			paths = append(paths, path)
		}
		return nil
	})
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	sort.Strings(paths)
	return paths, err
}

func resolveAll(ctx context.Context, resolvers ResolverSet, upstreams []engineimage.Upstream) ([]engineimage.ResolvedInput, error) {
	resolved := make([]engineimage.ResolvedInput, 0, len(upstreams))
	for _, u := range upstreams {
		r, err := resolvers.Resolve(ctx, u)
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, r)
	}
	sort.Slice(resolved, func(i, j int) bool { return resolved[i].Name < resolved[j].Name })
	return resolved, nil
}

func compareLock(path, flavor string, resolved []engineimage.ResolvedInput) (bool, string, error) {
	lock, err := engineimage.LoadLock(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, "missing-lock", nil
		}
		return false, "", err
	}
	if lock.Flavor != flavor {
		return false, "", fmt.Errorf("lock flavor %q does not match manifest flavor %q", lock.Flavor, flavor)
	}
	locked := append([]engineimage.ResolvedInput(nil), lock.ResolvedInputs...)
	sort.Slice(locked, func(i, j int) bool { return locked[i].Name < locked[j].Name })
	if !reflect.DeepEqual(locked, resolved) {
		return true, "upstream-changed", nil
	}
	return false, "unchanged", nil
}

func dockerfileRel(manifestRel, dockerfile string) string {
	if filepath.IsAbs(dockerfile) || filepath.Dir(dockerfile) != "." {
		return filepath.ToSlash(dockerfile)
	}
	return filepath.ToSlash(filepath.Join(filepath.Dir(manifestRel), dockerfile))
}

func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(r)
}
