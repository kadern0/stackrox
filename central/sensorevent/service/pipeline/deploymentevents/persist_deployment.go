package deploymentevents

import (
	"github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
)

func newPersistDeployment(deployments datastore.DataStore) *persistDeploymentImpl {
	return &persistDeploymentImpl{
		deployments: deployments,
	}
}

type persistDeploymentImpl struct {
	deployments datastore.DataStore
}

func (s *persistDeploymentImpl) do(action v1.ResourceAction, deployment *storage.Deployment) error {
	switch action {
	case v1.ResourceAction_CREATE_RESOURCE:
		if err := s.deployments.UpsertDeployment(deployment); err != nil {
			log.Errorf("unable to add deployment %s: %s", deployment.GetId(), err)
			return err
		}
	case v1.ResourceAction_UPDATE_RESOURCE:
		if err := s.deployments.UpsertDeployment(deployment); err != nil {
			log.Errorf("unable to update deployment %s: %s", deployment.GetId(), err)
			return err
		}
	case v1.ResourceAction_REMOVE_RESOURCE:
		if err := s.deployments.RemoveDeployment(deployment.GetId()); err != nil {
			log.Errorf("unable to remove deployment %s: %s", deployment.GetId(), err)
			return err
		}
	default:
		log.Warnf("unknown action: %s", action)
		return nil // Be interoperable: don't reject these requests.
	}
	return nil
}
