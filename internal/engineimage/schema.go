package engineimage

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ManifestSchema      = "kuvryn.engine-image/v1"
	LockSchema          = "kuvryn.engine-lock/v1"
	RepositoryPrefix    = "ghcr.io/azrtydxb/kuvryn-ai-engine/"
	DefaultLockFilename = "upstream.lock"
)

type ImageManifest struct {
	Schema       string            `yaml:"schema" json:"schema"`
	Name         string            `yaml:"name" json:"name"`
	Engine       string            `yaml:"engine" json:"engine"`
	Vendor       string            `yaml:"vendor" json:"vendor"`
	Variant      string            `yaml:"variant,omitempty" json:"variant,omitempty"`
	Image        Image             `yaml:"image" json:"image"`
	Build        Build             `yaml:"build" json:"build"`
	Capabilities Capabilities      `yaml:"capabilities" json:"capabilities"`
	Upstreams    []Upstream        `yaml:"upstreams" json:"upstreams"`
	Publish      Publish           `yaml:"publish" json:"publish"`
	Tests        []TestDeclaration `yaml:"tests" json:"tests"`
	Labels       map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

type Image struct {
	Repository string   `yaml:"repository" json:"repository"`
	Platforms  []string `yaml:"platforms" json:"platforms"`
	Base       Base     `yaml:"base" json:"base"`
}

type Base struct {
	Image string `yaml:"image" json:"image"`
	Tag   string `yaml:"tag" json:"tag"`
}

type Build struct {
	Context    string            `yaml:"context" json:"context"`
	Dockerfile string            `yaml:"dockerfile" json:"dockerfile"`
	Args       map[string]string `yaml:"args,omitempty" json:"args,omitempty"`
}

type Capabilities struct {
	Accelerators []string `yaml:"accelerators" json:"accelerators"`
	Collective   string   `yaml:"collective" json:"collective"`
	Distributed  bool     `yaml:"distributed" json:"distributed"`
	RDMA         string   `yaml:"rdma,omitempty" json:"rdma,omitempty"`
	RayRequired  bool     `yaml:"rayRequired,omitempty" json:"rayRequired,omitempty"`
}

type Upstream struct {
	Name       string `yaml:"name" json:"name"`
	Type       string `yaml:"type" json:"type"`
	Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
	Image      string `yaml:"image,omitempty" json:"image,omitempty"`
	Package    string `yaml:"package,omitempty" json:"package,omitempty"`
	Version    string `yaml:"version,omitempty" json:"version,omitempty"`
	Ref        string `yaml:"ref,omitempty" json:"ref,omitempty"`
	Path       string `yaml:"path,omitempty" json:"path,omitempty"`
	Digest     string `yaml:"digest,omitempty" json:"digest,omitempty"`
}

type Publish struct {
	Certification string   `yaml:"certification" json:"certification"`
	Tags          []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

type TestDeclaration struct {
	Name     string   `yaml:"name" json:"name"`
	Type     string   `yaml:"type" json:"type"`
	Required bool     `yaml:"required" json:"required"`
	Runner   []string `yaml:"runner,omitempty" json:"runner,omitempty"`
}

type UpstreamLock struct {
	Schema         string          `yaml:"schema" json:"schema"`
	Flavor         string          `yaml:"flavor" json:"flavor"`
	ResolvedInputs []ResolvedInput `yaml:"resolvedInputs" json:"resolvedInputs"`
	ImageDigest    string          `yaml:"imageDigest,omitempty" json:"imageDigest,omitempty"`
	PublishedTags  []string        `yaml:"publishedTags,omitempty" json:"publishedTags,omitempty"`
	TestEvidence   []TestEvidence  `yaml:"testEvidence,omitempty" json:"testEvidence,omitempty"`
}

type ResolvedInput struct {
	Name     string `yaml:"name" json:"name"`
	Type     string `yaml:"type" json:"type"`
	Value    string `yaml:"value" json:"value"`
	Digest   string `yaml:"digest,omitempty" json:"digest,omitempty"`
	Source   string `yaml:"source,omitempty" json:"source,omitempty"`
	Resolved string `yaml:"resolved,omitempty" json:"resolved,omitempty"`
}

type TestEvidence struct {
	Name   string `yaml:"name" json:"name"`
	Status string `yaml:"status" json:"status"`
	Report string `yaml:"report,omitempty" json:"report,omitempty"`
}

func LoadManifest(path string) (ImageManifest, error) {
	var m ImageManifest
	if err := loadStrictYAML(path, &m); err != nil {
		return m, err
	}
	return m, ValidateManifest(m)
}

func LoadLock(path string) (UpstreamLock, error) {
	var l UpstreamLock
	if err := loadStrictYAML(path, &l); err != nil {
		return l, err
	}
	if l.Schema != LockSchema {
		return l, fmt.Errorf("lock schema must be %q", LockSchema)
	}
	if l.Flavor == "" {
		return l, errors.New("lock flavor is required")
	}
	return l, nil
}

func loadStrictYAML(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}
