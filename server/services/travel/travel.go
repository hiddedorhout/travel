package travel

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Service is the Journey base type
type Service struct {
	db              *sql.DB
	storeStartStmnt *sql.Stmt
	updateEndStmnt  *sql.Stmt
}

// New initate the Journey service and returns the service
func New(db *sql.DB) (*Service, error) {

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS journeys(
		ID TEXT NOT NULL PRIMARY KEY,
		Type TEXT NOT NULL,
		StartTime TIMESTAMP,
		EndTime TIMESTAMP,
		Price FLOAT
	)`); err != nil {
		return nil, err
	}

	storeStartStmnt, err := db.Prepare(`INSERT INTO journeys (ID, Type, StartTime)
	VALUES (?,?,?)`)
	if err != nil {
		return nil, err
	}

	updateEndStmnt, err := db.Prepare(`UPDATE journeys SET EndTime = ? WHERE ID= ?`)
	if err != nil {
		return nil, err
	}

	return &Service{
		db:              db,
		storeStartStmnt: storeStartStmnt,
		updateEndStmnt:  updateEndStmnt,
	}, nil
}

// StartJourney stores the start of a journey and returns the ID
func (s *Service) StartJourney(journeyType string, startTime time.Time) (ID *string, err error) {
	id := uuid.New().String()
	if _, err := s.storeStartStmnt.Exec(id, journeyType, startTime); err != nil {
		return nil, err
	}
	return &id, err
}

// EndJourney stores the ends the journey
func (s *Service) EndJourney(ID string, endTime time.Time) error {

	if _, err := s.updateEndStmnt.Exec(endTime, ID); err != nil {
		return err
	}
	return nil
}
