package routes

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	db                *sql.DB
	createRouteStmnt  *sql.Stmt
	getRouteListStmnt *sql.Stmt
	getRouteStmnt     *sql.Stmt
}

type route struct {
	ID          string `json:"routeID"`
	Description string `json:"description"`
}

func New(db *sql.DB) (*Service, error) {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS routes(
		ID TEXT NOT NULL PRIMARY KEY,
		Description TEXT NOT NULL,
		Route BLOB NOT NULL
	)`); err != nil {
		return nil, err
	}

	createRouteStmnt, err := db.Prepare(`INSERT INTO routes (ID, Description, Route) VALUES (?,?,?)`)
	if err != nil {
		return nil, err
	}

	getRouteListStmnt, err := db.Prepare(`SELECT ID, Description FROM routes`)
	if err != nil {
		return nil, err
	}

	getRouteStmnt, err := db.Prepare(`SELECT Route FROM routes WHERE ID=?`)
	if err != nil {
		return nil, err
	}

	return &Service{
		createRouteStmnt:  createRouteStmnt,
		getRouteListStmnt: getRouteListStmnt,
		getRouteStmnt:     getRouteStmnt,
	}, nil

}

func (s *Service) CreateRoute(description string, route []byte) (ID *string, err error) {

	if !json.Valid(route) {
		return nil, errors.New("Invalid json")
	}

	id := uuid.New().String()
	if _, err := s.createRouteStmnt.Exec(id, description, route); err != nil {
		return nil, err
	}
	return &id, nil
}

func (s *Service) GetRouteList() (jsonList *[]byte, err error) {

	routelist := make([]route, 0)

	rows, err := s.getRouteListStmnt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r route
		if err := rows.Scan(&r.ID, r.Description); err != nil {
			return nil, err
		}
		routelist = append(routelist, r)
	}
	defer rows.Close()

	result, err := json.Marshal(routelist)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) GetRoute(ID string) (route *[]byte, err error) {

	var r []byte
	if err := s.getRouteStmnt.QueryRow(ID).Scan(&r); err != nil {
		return nil, err
	}
	if !json.Valid(r) {
		return nil, errors.New("Invalid json")
	}
	return &r, nil
}
