package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Location struct {
	Id        uint
	Guid      string
	Name      string
	Usages    uint
	CreatedAt string
	UpdatedAt string
	Tags      []*LocationTag
	Resources struct {
		Bundle   string
		Manifest string
		Config   string
		Icon     string
		DllPath  string
	}
}

func (l *Location) ResourcesDirectory() string {
	return fmt.Sprintf("data/locations/resources/%s/%s", l.Guid[:2], l.Guid)
}

func (l *Location) fillResources() {
	resPrefix := fmt.Sprintf("/data/locations/resources/%s/%s", l.Guid[:2], l.Guid)
	l.Resources.Bundle = resPrefix + "/bundle"
	l.Resources.Manifest = resPrefix + "/bundle.manifest"
	l.Resources.Config = resPrefix + "/bundle.json"
	l.Resources.Icon = resPrefix + "/bundle.png"
	l.Resources.DllPath = resPrefix
}

func (l *Location) Insert(db *sql.DB) error {
	var err error

	const query = `
		INSERT INTO locations(
			created_at,
			updated_at,
			guid,
			name
		)
		VALUES (
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP,
			$1,
			$2
		)
		RETURNING 
			id, 
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	err = db.QueryRow(
		query,
		l.Guid,
		l.Name,
	).Scan(
		&l.Id,
		&l.CreatedAt,
		&l.UpdatedAt,
	)

	if err != nil {
		return err
	}

	l.fillResources()

	return nil
}

func (l *Location) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE locations
		SET
			name = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $2
		RETURNING
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	err = db.QueryRow(
		query,
		l.Name,
		l.Id,
	).Scan(
		&l.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (l *Location) LoadTags(db *sql.DB) error {
	l.Tags = make([]*LocationTag, 0)

	const query = `
		SELECT
			ot.id,
			ot.created_at,
			ot.updated_at,
			ot.text
		FROM location_tags AS ot
		INNER JOIN location_tag_locations AS oto ON ot.id = oto.location_tag_id
		WHERE oto.location_id = $1
`

	rows, err := db.Query(query, l.Id)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
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
			return err
		}
		l.Tags = append(l.Tags, t)
	}

	return nil
}

func UpdateLocationTags(db *sql.DB, l *Location, oldTags, newTags []*LocationTag) error {
	var tagsToDelete []string
oldTagLoop:
	for _, oldTag := range oldTags {
		for _, newTag := range newTags {
			if newTag.Id == oldTag.Id {
				continue oldTagLoop
			}
		}
		tagsToDelete = append(tagsToDelete, strconv.Itoa(int(oldTag.Id)))
	}

	var tagsToInsert []string
newTagLoop:
	for _, newTag := range newTags {
		for _, oldTag := range oldTags {
			if oldTag.Id == newTag.Id {
				continue newTagLoop
			}
		}
		_, err := GetLocationTagById(db, newTag.Id)
		if err != nil {
			return err
		}
		tagsToInsert = append(tagsToInsert, fmt.Sprintf("(%d, %d, CURRENT_TIMESTAMP)", newTag.Id, l.Id))
	}

	if len(tagsToInsert) > 0 {
		_, err := db.Exec(fmt.Sprintf(`
			INSERT INTO location_tag_locations(
				location_tag_id,
				location_id,
				created_at
			)
			VALUES %s
		`, strings.Join(tagsToInsert, ",")))
		if err != nil {
			return err
		}
	}

	if len(tagsToDelete) > 0 {
		_, err := db.Exec(fmt.Sprintf(`
			DELETE FROM location_tag_locations
			WHERE 
				location_tag_id IN(%s) AND
				location_id = $1
		`, strings.Join(tagsToDelete, ",")), l.Id)
		if err != nil {
			return err
		}
	}

	l.Tags = newTags

	return nil
}

func GetLocationById(db *sql.DB, id uint) (l *Location, err error) {
	const query = `
		SELECT
			id,
			guid,
			name,
			(SELECT COUNT(id) FROM world_locations WHERE location_id = id),
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM locations
		WHERE
			id = $1
	`

	l = new(Location)
	err = db.QueryRow(query, id).Scan(
		&l.Id,
		&l.Guid,
		&l.Name,
		&l.Usages,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	l.fillResources()

	return l, nil
}

func GetLocationByGuid(db *sql.DB, guid string) (l *Location, err error) {
	const query = `
		SELECT
			id,
			guid,
			name,
			(SELECT COUNT(id) FROM world_locations WHERE location_id = id),
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM locations
		WHERE
			guid = $1
	`

	l = new(Location)
	err = db.QueryRow(query, guid).Scan(
		&l.Id,
		&l.Guid,
		&l.Name,
		&l.Usages,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	l.fillResources()

	return l, nil
}

func GetLocations(db *sql.DB) (result []*Location, err error) {
	result = make([]*Location, 0)

	const query = `
		SELECT
			id,
			guid,
			name,
			(SELECT COUNT(id) FROM world_locations WHERE location_id = id),
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM locations
		ORDER BY id
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
		l := new(Location)
		err = rows.Scan(
			&l.Id,
			&l.Guid,
			&l.Name,
			&l.Usages,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		l.fillResources()

		result = append(result, l)
	}

	return result, nil
}

func DeleteLocation(db *sql.DB, l *Location) error {
	const query = `
		DELETE 
		FROM locations
		WHERE id = $1
	`

	var err error
	_, err = db.Exec(query, l.Id)

	return err
}
