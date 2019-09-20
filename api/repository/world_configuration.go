package repository

import (
	"database/sql"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
)

type WorldConfiguration struct {
	Id                   uint
	Sid                  string
	Name                 string `validate:"required"`
	WorldId              uint   `validate:"required"`
	StartWorldLocationId uint   `validate:"required"`
	GroupIds             []uint `validate:"required"`
	CreatedAt            string
	UpdatedAt            string
}

func (wc *WorldConfiguration) Validate() error {
	return validate.Struct(wc)
}

func (wc *WorldConfiguration) Insert(db *sql.DB) error {
	var err error
	const query = `
		INSERT INTO world_configurations(
			created_at,
			updated_at,
			sid,
			world_id,
			name,
			start_world_location_id
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

	if wc.Sid == "" {
		wc.Sid = uuid.Must(uuid.NewV4()).String()
	}

	err = db.QueryRow(
		query,
		wc.Sid,
		wc.WorldId,
		wc.Name,
		wc.StartWorldLocationId,
	).Scan(
		&wc.Id,
		&wc.CreatedAt,
		&wc.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return wc.InsertGroups(db)
}

func (wc *WorldConfiguration) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE world_configurations
		SET
			updated_at = CURRENT_TIMESTAMP,
			name = $1,
			start_world_location_id = $2
		WHERE
			id = $3
		RETURNING
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	err = db.QueryRow(
		query,
		wc.Name,
		wc.StartWorldLocationId,
		wc.Id,
	).Scan(
		&wc.UpdatedAt,
	)

	if err != nil {
		return err
	}

	err = DeleteWorldConfigurationGroups(db, wc)
	if err != nil {
		return err
	}

	return wc.InsertGroups(db)
}

func (wc *WorldConfiguration) InsertGroups(db *sql.DB) error {
	var err error
	var placeholders []string
	var values []interface{}
	var placeholderNum = 1
	for _, groupId := range wc.GroupIds {
		placeholders = append(placeholders, fmt.Sprintf("(CURRENT_TIMESTAMP, $%d, $%d)", placeholderNum, placeholderNum+1))
		values = append(values, wc.Id, groupId)
		placeholderNum += 2
	}

	var query = fmt.Sprintf(`
		INSERT INTO world_configuration_groups(
			created_at,
			world_configuration_id,
			group_id
		)
		VALUES %s
	`, strings.Join(placeholders, ","))

	_, err = db.Exec(query, values...)

	return err
}

func DeleteWorldConfigurationGroups(db *sql.DB, wc *WorldConfiguration) error {
	const query = `
		DELETE 
		FROM world_configuration_groups
		WHERE world_configuration_id = $1
	`

	var err error
	_, err = db.Exec(query, wc.Id)

	return err
}

func DeleteWorldConfiguration(db *sql.DB, wc *WorldConfiguration) error {
	const query = `
		DELETE 
		FROM world_configurations
		WHERE id = $1
	`

	err := DeleteWorldConfigurationGroups(db, wc)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, wc.Id)

	return err
}

func GetWorldConfigurationById(db *sql.DB, id uint) (wc *WorldConfiguration, err error) {
	const query = `
		SELECT
			id,
			sid,
			world_id,
			name,
			start_world_location_id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM world_configurations
		WHERE
			id = $1
	`

	wc = new(WorldConfiguration)
	err = db.QueryRow(query, id).Scan(
		&wc.Id,
		&wc.Sid,
		&wc.WorldId,
		&wc.Name,
		&wc.StartWorldLocationId,
		&wc.CreatedAt,
		&wc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	wc.GroupIds, err = GetWorldConfigurationGroups(db, wc.Id)
	if err != nil {
		return nil, err
	}

	return wc, nil
}

func GetWorldConfigurationsByWorldId(db *sql.DB, worldId uint) (result []*WorldConfiguration, err error) {
	result = make([]*WorldConfiguration, 0)

	const query = `
		SELECT
			id,
			sid,
			world_id,
			name,
			start_world_location_id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM world_configurations
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
		wc := new(WorldConfiguration)
		err = rows.Scan(
			&wc.Id,
			&wc.Sid,
			&wc.WorldId,
			&wc.Name,
			&wc.StartWorldLocationId,
			&wc.CreatedAt,
			&wc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		wc.GroupIds, err = GetWorldConfigurationGroups(db, wc.Id)
		if err != nil {
			return nil, err
		}

		result = append(result, wc)
	}

	return result, nil
}

func GetWorldConfigurationGroups(db *sql.DB, worldConfigurationId uint) ([]uint, error) {
	result := make([]uint, 0)

	const query = `
		SELECT group_id
		FROM world_configuration_groups
		WHERE world_configuration_id = $1
`
	rows, err := db.Query(query, worldConfigurationId)
	if err != nil && err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupId uint
	for rows.Next() {
		err = rows.Scan(&groupId)
		if err != nil {
			return nil, err
		}
		result = append(result, groupId)
	}

	return result, nil
}
