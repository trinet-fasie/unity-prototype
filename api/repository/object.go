package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Object struct {
	Id        uint
	Guid      string
	CreatedAt string
	UpdatedAt string
	Usages    uint
	Config    map[string]interface{}
	Tags      []*ObjectTag
	Resources struct {
		Bundle   string
		Manifest string
		Config   string
		Icon     string
		DllPath  string
	}
}

func (o *Object) fillResources() {
	resPrefix := fmt.Sprintf("/data/objects/resources/%s/%s", o.Guid[:2], o.Guid)
	o.Resources.Bundle = resPrefix + "/bundle"
	o.Resources.Manifest = resPrefix + "/bundle.manifest"
	o.Resources.Config = resPrefix + "/bundle.json"
	o.Resources.Icon = resPrefix + "/bundle.png"
	o.Resources.DllPath = resPrefix
}

func (o *Object) ResourcesDirectory() string {
	return fmt.Sprintf("data/objects/resources/%s/%s", o.Guid[:2], o.Guid)
}

func (o *Object) Insert(db *sql.DB) error {
	var err error

	var config, _ = json.Marshal(o.Config)
	const query = `
		INSERT INTO objects(
			created_at,
			updated_at,
			guid,
			config
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
		o.Guid,
		config,
	).Scan(
		&o.Id,
		&o.CreatedAt,
		&o.UpdatedAt,
	)

	if err != nil {
		return err
	}

	o.fillResources()

	return nil
}

func (o *Object) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE objects
		SET
			config = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $2
		RETURNING
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	var config, _ = json.Marshal(o.Config)

	err = db.QueryRow(
		query,
		config,
		o.Id,
	).Scan(
		&o.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (o *Object) LoadTags(db *sql.DB) error {
	o.Tags = make([]*ObjectTag, 0)

	const query = `
		SELECT
			ot.id,
			ot.created_at,
			ot.updated_at,
			ot.text
		FROM object_tags AS ot
		INNER JOIN object_tag_objects AS oto ON ot.id = oto.object_tag_id
		WHERE oto.object_id = $1
`

	rows, err := db.Query(query, o.Id)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
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
			return err
		}
		o.Tags = append(o.Tags, t)
	}

	return nil
}

func UpdateObjectTags(db *sql.DB, o *Object, oldTags, newTags []*ObjectTag) error {
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
		_, err := GetObjectTagById(db, newTag.Id)
		if err != nil {
			return err
		}
		tagsToInsert = append(tagsToInsert, fmt.Sprintf("(%d, %d, CURRENT_TIMESTAMP)", newTag.Id, o.Id))
	}

	if len(tagsToInsert) > 0 {
		_, err := db.Exec(fmt.Sprintf(`
			INSERT INTO object_tag_objects(
				object_tag_id,
				object_id,
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
			DELETE FROM object_tag_objects
			WHERE 
				object_tag_id IN(%s) AND
				object_id = $1
		`, strings.Join(tagsToDelete, ",")), o.Id)
		if err != nil {
			return err
		}
	}

	o.Tags = newTags

	return nil
}

func GetObjects(db *sql.DB) (result []*Object, err error) {
	result = make([]*Object, 0)

	const query = `
		SELECT
			id,
			guid,
			config,
 			(SELECT COUNT(id) FROM group_objects WHERE object_id = id),
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM objects
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

	var config []byte
	for rows.Next() {
		o := new(Object)
		err = rows.Scan(
			&o.Id,
			&o.Guid,
			&config,
			&o.Usages,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(config, &o.Config)

		o.fillResources()

		result = append(result, o)
	}

	return result, nil
}

func GetObjectByGuid(db *sql.DB, guid string) (o *Object, err error) {
	const query = `
		SELECT
			id,
			guid,
			config,
 			(SELECT COUNT(id) FROM group_objects WHERE object_id = id),
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM objects
		WHERE guid = $1
`

	row := db.QueryRow(query, guid)
	if err != nil {
		return nil, err
	}

	var config []byte
	o = new(Object)
	err = row.Scan(
		&o.Id,
		&o.Guid,
		&config,
		&o.Usages,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(config, &o.Config)

	o.fillResources()

	return o, nil
}

func GetObjectById(db *sql.DB, id uint) (o *Object, err error) {
	const query = `
		SELECT
			id,
			guid,
			config,
 			(SELECT COUNT(id) FROM group_objects WHERE object_id = id),
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM objects
		WHERE id = $1
`

	row := db.QueryRow(query, id)
	if err != nil {
		return nil, err
	}

	var config []byte
	o = new(Object)
	err = row.Scan(
		&o.Id,
		&o.Guid,
		&config,
		&o.Usages,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(config, &o.Config)

	o.fillResources()

	return o, nil
}

func DeleteObject(db *sql.DB, o *Object) error {
	const query = `
		DELETE 
		FROM objects
		WHERE id = $1
	`

	var err error
	_, err = db.Exec(query, o.Id)

	return err
}
