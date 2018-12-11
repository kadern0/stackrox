package enrichanddetect

import (
	"github.com/stackrox/rox/central/detection/lifecycle"
	"github.com/stackrox/rox/central/enrichment"
	"github.com/stackrox/rox/generated/storage"
)

type enricherAndDetectorImpl struct {
	enricher enrichment.Enricher
	manager  lifecycle.Manager
}

// EnrichAndDetect runs enrichment and detection on a deployment.
func (e *enricherAndDetectorImpl) EnrichAndDetect(deployment *storage.Deployment) error {
	updated, err := e.enricher.Enrich(deployment)
	if err != nil {
		return err
	}
	if updated {
		_, _, err := e.manager.DeploymentUpdated(deployment)
		return err
	}
	return nil
}
