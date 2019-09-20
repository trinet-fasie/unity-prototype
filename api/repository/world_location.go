package repository

import (
	"database/sql"
	"github.com/satori/go.uuid"
)

type WorldLocation struct {
	Id         uint
	Sid        string
	Name       string `validate:"required"`
	WorldId    uint   `validate:"required"`
	LocationId uint   `validate:"required"`
	CreatedAt  string
	UpdatedAt  string
}

func (wl *WorldLocation) Validate() error {
	return validate.Struct(wl)
}

func (wl *WorldLocation) Insert(db *sql.DB) error {
	var err error
	const query = `
		INSERT INTO world_locations(
			created_at,
			updated_at,
			sid,
			world_id,
			location_id,
			name
		)
		VALUES (
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP,
			$1,
			$2,
			$3,
			$4

		)
		RETURNING 
			id, 
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	if wl.Sid == "" {
		wl.Sid = uuid.Must(uuid.NewV4()).String()
	}

	err = db.QueryRow(
		query,
		wl.Sid,
		wl.WorldId,
		wl.LocationId,
		wl.Name,
	).Scan(
		&wl.Id,
		&wl.CreatedAt,
		&wl.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (wl *WorldLocation) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE world_locations
		SET
			updated_at = CURRENT_TIMESTAMP,
			location_id = $1,
			name = $2
		WHERE
			id = $3
		RETURNING
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	err = db.QueryRow(
		query,
		wl.LocationId,
		wl.Name,
		wl.Id,
	).Scan(
		&wl.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func DeleteWorldLocation(db *sql.DB, wl *WorldLocation) error {
	const query = `
		DELETE 
		FROM world_locations
		WHERE id = $1
	`

	var err error
	_, err = db.Exec(query, wl.Id)

	return err
}

func GetWorldLocationById(db *sql.DB, id uint) (wl *WorldLocation, err error) {
	const query = `
		SELECT
			id,
			sid,
			world_id,
			location_id,
			name,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM world_locations
		WHERE
			id = $1
	`

	wl = new(WorldLocation)
	err = db.QueryRow(query, id).Scan(
		&wl.Id,
		&wl.Sid,
		&wl.WorldId,
		&wl.LocationId,
		&wl.Name,
		&wl.CreatedAt,
		&wl.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

func GetWorldLocationsByWorldId(db *sql.DB, worldId uint) (result []*WorldLocation, err error) {
	result = make([]*WorldLocation, 0)

	const query = `
		SELECT
			id,
			sid,
			world_id,
			location_id,
			name,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM world_locations
		WHERE
			world_id = $1
		ORDER BY id
	`

	rows, err := db.Query(query, worldId)
	if err != nil && err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		wl := new(WorldLocation)
		err = rows.Scan(
			&wl.Id,
			&wl.Sid,
			&wl.WorldId,
			&wl.LocationId,
			&wl.Name,
			&wl.CreatedAt,
			&wl.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, wl)
	}

	return result, nil
}
