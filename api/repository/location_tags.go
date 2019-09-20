package repository

import (
	"database/sql"
)

type LocationTag struct {
	Id        uint
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	Text      string
}

func (t *LocationTag) Insert(db *sql.DB) error {
	var err error
	const query = `
		INSERT INTO location_tags(
			created_at,
			updated_at,
			text
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
		t.Text,
	).Scan(
		&t.Id,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *LocationTag) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE location_tags
		SET
			updated_at = CURRENT_TIMESTAMP,
			text = $1
		WHERE
			id = $2
		RETURNING
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	err = db.QueryRow(
		query,
		t.Text,
		t.Id,
	).Scan(
		&t.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetLocationTagByText(db *sql.DB, text string) (t *LocationTag, err error) {
	const query = `
		SELECT
			id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			text
		FROM location_tags
		WHERE
			text = $1
	`

	t = new(LocationTag)
	err = db.QueryRow(query, text).Scan(
		&t.Id,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.Text,
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func GetLocationTags(db *sql.DB, text string, limit uint) (result []*LocationTag, err error) {
	result = make([]*LocationTag, 0)

	const query = `
		SELECT
			id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			text
		FROM location_tags
		WHERE
			text LIKE '%' || $1 || '%'
		ORDER BY text
		LIMIT $2
`

	rows, err := db.Query(query, text, limit)
	if err != nil && err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := new(LocationTag)
		err = rows.Scan(
			&t.Id,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.Text,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func GetLocationTagById(db *sql.DB, id uint) (t *LocationTag, err error) {
	const query = `
		SELECT
			id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			text
		FROM location_tags
		WHERE
			id = $1
	`

	t = new(LocationTag)
	err = db.QueryRow(query, id).Scan(
		&t.Id,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.Text,
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}
