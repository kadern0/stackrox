package reconciliation

import (
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/set"
)

var (
	log = logging.LoggerForModule()
)

// Perform factors out some of the common reconciliation logic to avoid duplication.
// It does the reconciliation logic, and closes the passed store at the end.
func Perform(store Store, existingIDs set.StringSet, resourceType string, removeFunc func(id string) error) error {
	defer store.Close()
	idsToDelete := existingIDs.Difference(store.GetSet()).AsSlice()
	if len(idsToDelete) == 0 {
		return nil
	}

	log.Infof("Deleting %d %s as a part of reconciliation", len(idsToDelete), resourceType)

	errList := errorhelpers.NewErrorList("Network Policy reconciliation")
	for _, id := range idsToDelete {
		errList.AddError(removeFunc(id))
	}
	return errList.ToError()
}
