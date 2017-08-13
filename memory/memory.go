// Package memory implements a darrel store in memory
package memory

import (
	"sync"

	"github.com/bakins/darrell/api"
)

type Memory struct {
	sync.Mutex
	artifacts      map[string]*api.Artifact
	artifactBuilds map[string]*api.ArtifactBuild
}

func New() *Memory {
	return &Memory{
		artifacts:      make(map[string]*api.Artifact),
		artifactBuilds: make(map[string]*api.ArtifactBuild),
	}
}
