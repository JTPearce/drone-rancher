package main

import (
	"errors"
	"time"

	"drone-rancher/rancher"

	log "github.com/Sirupsen/logrus"
)

// Plugin holds details passed to us by drone
type Plugin struct {
	URL            string
	Token          string
	Project        string
	Deployment     string
	Namespace      string
	DockerImage    string
	StartFirst     bool
	Confirm        bool
	Timeout        int
	IntervalMillis int64
	BatchSize      int64
	YamlVerified   bool
}

// Exec is the entrypoint of the plugin
func (p *Plugin) Exec() error {
	log.Info("Drone Rancher Initialised")

	if p.URL == "" || p.Token == "" || p.Project == "" || p.Deployment == "" || p.Namespace == "" {
		return errors.New("Error: Must have url, token, project, deployment and namespace defined")
	}

	id := "deployment:" + p.Namespace + ":" + p.Deployment

	r := rancher.New(p.Token, p.Project, p.URL)

	if p.DockerImage != "" {
		err := r.UpdateImage(id, p.DockerImage, p.Deployment)
		if err != nil {
			log.Errorf("Error while updating image: %v", err)
			return err
		}
		log.Info("Image Updated")
		return nil
	}

	err := r.Redeploy(id)
	if err != nil {
		log.Errorf("Error while redeploying container: %v", err)
		return err
	}
	log.Info("Deployment Updated")
	return nil
}

type retryFunc func() (interface{}, error)

func retry(f retryFunc, timeout time.Duration, interval time.Duration) (interface{}, error) {
	finish := time.After(timeout)
	for {
		result, err := f()
		if err == nil {
			return result, nil
		}
		select {
		case <-finish:
			return nil, err
		case <-time.After(interval):
		}
	}
}
