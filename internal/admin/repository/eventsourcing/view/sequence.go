package view

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	sequencesTable = "adminapi.current_sequences"
)

func (v *View) saveCurrentSequence(viewName string, event *models.Event) error {
	return repository.SaveCurrentSequence(v.Db, sequencesTable, viewName, event.InstanceID, event.Sequence, event.CreationDate)
}

func (v *View) latestSequence(viewName, instanceID string) (*repository.CurrentSequence, error) {
	return repository.LatestSequence(v.Db, sequencesTable, viewName, instanceID)
}

func (v *View) latestSequences(viewName string, instanceIDs []string) ([]*repository.CurrentSequence, error) {
	return repository.LatestSequences(v.Db, sequencesTable, viewName, instanceIDs)
}

func (v *View) AllCurrentSequences(db, instanceID string) ([]*repository.CurrentSequence, error) {
	return repository.AllCurrentSequences(v.Db, db+".current_sequences", instanceID)
}

func (v *View) updateSpoolerRunSequence(viewName string, instanceIDs []string) error {
	currentSequences, err := repository.LatestSequences(v.Db, sequencesTable, viewName, instanceIDs)
	if err != nil {
		return err
	}
	for _, currentSequence := range currentSequences {
		if currentSequence.ViewName == "" {
			currentSequence.ViewName = viewName
		}
		currentSequence.LastSuccessfulSpoolerRun = time.Now()
	}
	return repository.UpdateCurrentSequences(v.Db, sequencesTable, currentSequences)
}

func (v *View) ClearView(db, viewName string) error {
	truncateView := db + "." + viewName
	sequenceTable := db + ".current_sequences"
	return repository.ClearView(v.Db, truncateView, sequenceTable)
}
