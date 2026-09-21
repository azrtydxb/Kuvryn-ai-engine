package clusterprobe

import (
	"bufio"
	"fmt"
	"strings"
)

type ProbeSections map[string][]string

func ParseSections(output string) (ProbeSections, error) {
	sections := ProbeSections{}
	current := ""
	sawEnd := false
	s := bufio.NewScanner(strings.NewReader(output))
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "===KUVRYN:") && strings.HasSuffix(line, "===") {
			name := strings.TrimSuffix(strings.TrimPrefix(line, "===KUVRYN:"), "===")
			if name == "end" {
				sawEnd = true
				current = ""
				continue
			}
			current = name
			sections[current] = nil
			continue
		}
		if current != "" {
			sections[current] = append(sections[current], line)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	for _, required := range []string{"accelerators", "network", "fabric", "runtime"} {
		if _, ok := sections[required]; !ok {
			return nil, fmt.Errorf("missing probe section %q", required)
		}
	}
	if !sawEnd {
		return nil, fmt.Errorf("missing probe end sentinel")
	}
	return sections, nil
}

func RecordFromSections(host string, sections ProbeSections) HostRecord {
	r := HostRecord{Host: host}
	if len(sections["accelerators"]) == 0 {
		r.Limitations = append(r.Limitations, "no accelerator command output")
	}
	for _, f := range sections["fabric"] {
		f = strings.TrimSpace(f)
		if f != "" {
			r.Fabric = append(r.Fabric, FabricInterface{Name: f, Kind: "infiniband", RDMADevice: f})
		}
	}
	if len(r.Fabric) == 0 {
		r.Limitations = append(r.Limitations, "no RDMA fabric discovered")
	}
	return r
}
