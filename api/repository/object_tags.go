package repository

import (
	"database/sql"
)

type ObjectTag struct {
	Id        uint
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	Text      string
}

func (t *ObjectTag) Insert(db *sql.DB) error {
	var err error
	const query = `
		INSERT INTO object_tags(
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

func (t *ObjectTag) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE object_tags
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

func GetObjectTagByText(db *sql.DB, text string) (t *ObjectTag, err error) {
	const query = `
		SELECT
			id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			text
		FROM object_tags
		WHERE
			text = $1
	`

	t = new(ObjectTag)
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

func GetObjectTags(db *sql.DB, text string, limit uint) (result []*ObjectTag, err error) {
	result = make([]*ObjectTag, 0)

	const query = `
		SELECT
			id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			text
		FROM object_tags
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
		t := new(ObjectTag)
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

func GetObjectTagById(db *sql.DB, id uint) (t *ObjectTag, err error) {
	const query = `
		SELECT
			id,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS'),
			text
		FROM object_tags
		WHERE
			id = $1
	`

	t = new(ObjectTag)
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
