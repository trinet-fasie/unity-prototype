package repository

import (
	"database/sql"
)

type World struct {
	Id             uint
	Name           string `validate:"required"`
	Configurations uint
	CreatedAt      string
	UpdatedAt      string
}

func (w *World) Validate() error {
	return validate.Struct(w)
}

func (w *World) Insert(db *sql.DB) error {
	var err error
	const query = `
		INSERT INTO worlds(
			created_at,
			updated_at,
			name
		)
		VALUES (
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP,
			$1
		)
		RETURNING 
			id, 
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	err = db.QueryRow(
		query,
		w.Name,
	).Scan(
		&w.Id,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (w *World) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE worlds
		SET
			updated_at = CURRENT_TIMESTAMP,
			name = $1
		WHERE
			id = $2
		RETURNING
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	err = db.QueryRow(
		query,
		w.Name,
		w.Id,
	).Scan(
		&w.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetWorlds(db *sql.DB) (result []*World, err error) {
	result = make([]*World, 0)

	const query = `
		SELECT
			w.id,
			w.name,
 			(SELECT COUNT(id) FROM world_configurations WHERE world_id = w.id) AS configurations,
			to_char(w.created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(w.updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM worlds AS w
		ORDER BY w.id
`

	rows, err := db.Query(query)
	if err != nil && err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		w := new(World)
		err = rows.Scan(
			&w.Id,
			&w.Name,
			&w.Configurations,
			&w.CreatedAt,
			&w.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, w)
	}

	return result, nil
}

func DeleteWorld(db *sql.DB, o *World) error {
	const query = `
		DELETE 
		FROM worlds
		WHERE id = $1
	`

	var err error
	_, err = db.Exec(query, o.Id)

	return err
}

func GetWorldById(db *sql.DB, id uint) (w *World, err error) {
	const query = `
		SELECT
			w.id,
			to_char(w.created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(w.updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			w.name
		FROM worlds AS w
		WHERE
			id = $1
	`

	w = new(World)
	err = db.QueryRow(query, id).Scan(
		&w.Id,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.Name,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}
