// Package memory implements a darrel store in memory
package memory

import (
	"sync"

	"github.com/bakins/alecton"
	"github.com/bakins/alecton/api"
)

type Memory struct {
	sync.Mutex
	artifacts    map[string]*api.Artifact
	applications map[string]*api.Application
	deployments  map[string]*api.Deployment
}

func New() *Memory {
	return &Memory{
		artifacts:    make(map[string]*api.Artifact),
		applications: make(map[string]*api.Application),
		deployments:  make(map[string]*api.Deployment),
	}
}

func provider(map[string]interface{}) (alecton.StorageProvider, error) {
	return New(), nil
}

func init() {
	alecton.RegisterStorageProvider("memory", provider)
}
