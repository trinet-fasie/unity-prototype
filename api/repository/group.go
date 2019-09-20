package repository

import (
	"database/sql"
	"encoding/json"
)

type Group struct {
	Id              uint
	WorldLocationId uint   `validate:"required"`
	Name            string `validate:"required"`
	Code            string
	EditorData      map[string]interface{}
	CreatedAt       string
	UpdatedAt       string
}

func NewGroup() *Group {
	g := new(Group)
	g.EditorData = make(map[string]interface{})

	return g
}

func (g *Group) Validate() error {
	return validate.Struct(g)
}

func (g *Group) Insert(db *sql.DB) error {
	var err error
	const query = `
		INSERT INTO groups(
			created_at,
			updated_at,
			world_location_id,
			name,
			code,
			editor_data
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

	var editorData, _ = json.Marshal(g.EditorData)
	err = db.QueryRow(
		query,
		g.WorldLocationId,
		g.Name,
		g.Code,
		editorData,
	).Scan(
		&g.Id,
		&g.CreatedAt,
		&g.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (g *Group) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE groups
		SET
			name = $1,
			code = $2,
			editor_data = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $4
		RETURNING
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	var editorData, _ = json.Marshal(g.EditorData)
	err = db.QueryRow(
		query,
		g.Name,
		g.Code,
		editorData,
		g.Id,
	).Scan(
		&g.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func DeleteGroup(db *sql.DB, o *Group) error {
	const query = `
		DELETE 
		FROM groups
		WHERE id = $1
	`

	var err error
	_, err = db.Exec(query, o.Id)

	return err
}

func GetGroupById(db *sql.DB, id uint) (g *Group, err error) {
	const query = `
		SELECT
			id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			world_location_id,
			name,
			code,
			editor_data
		FROM groups AS g
		WHERE
			id = $1
	`

	var editorData []byte
	g = new(Group)
	err = db.QueryRow(query, id).Scan(
		&g.Id,
		&g.CreatedAt,
		&g.UpdatedAt,
		&g.WorldLocationId,
		&g.Name,
		&g.Code,
		&editorData,
	)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(editorData, &g.EditorData)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func GetGroupsByWorldLocationId(db *sql.DB, worldLocationId uint) (result []*Group, err error) {
	result = make([]*Group, 0)

	const query = `
		SELECT
			id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			world_location_id,
			name,
			code,
			editor_data
		FROM groups
		WHERE
			world_location_id = $1
		ORDER BY id
	`

	rows, err := db.Query(query, worldLocationId)
	if err != nil && err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var editorData []byte
	for rows.Next() {
		g := new(Group)
		err = rows.Scan(
			&g.Id,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.WorldLocationId,
			&g.Name,
			&g.Code,
			&editorData,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(editorData, &g.EditorData)
		if err != nil {
			return nil, err
		}

		result = append(result, g)
	}

	return result, nil
}
