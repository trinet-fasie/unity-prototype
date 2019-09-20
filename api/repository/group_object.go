package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type GroupObject struct {
	Id           uint
	ParentId     sql.NullInt64
	Position     uint
	GroupId      uint `validate:"required"`
	ObjectId     uint `validate:"required"`
	InstanceId   uint `validate:"required"`
	Name         string
	Data         map[string]interface{}
	Locked       bool           `json:"-"`
	GroupObjects []*GroupObject `validate:"dive"`
	CreatedAt    string
	UpdatedAt    string
}

func (o *GroupObject) Validate() error {
	return validate.Struct(o)
}

func (o *GroupObject) InsertOrUpdate(db *sql.DB) error {
	if o.Id == 0 {
		return o.Insert(db)
	} else {
		return o.Update(db)
	}
}

func (o *GroupObject) Insert(db *sql.DB) error {
	var err error
	const query = `
		INSERT INTO group_objects(
			created_at,
			updated_at,
			parent_id,
			position,
			group_id,
			object_id,
			instance_id,
			name,
			data
		)
		VALUES (
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP,
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
		RETURNING 
			id, 
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
	`

	if o.Name == "" {
		o.Name = strconv.Itoa(int(o.InstanceId))
	}

	var data, _ = json.Marshal(o.Data)
	err = db.QueryRow(
		query,
		o.ParentId,
		o.Position,
		o.GroupId,
		o.ObjectId,
		o.InstanceId,
		o.Name,
		data,
	).Scan(
		&o.Id,
		&o.CreatedAt,
		&o.UpdatedAt,
	)

	if err != nil {
		return errors.New(fmt.Sprintf("cannot insert object: %s", err))
	}

	return nil
}

func (o *GroupObject) Update(db *sql.DB) error {
	var err error
	const query = `
		UPDATE group_objects
		SET
			parent_id = $1,
			position = $2,
			instance_id = $3,
			name = $4,
			data = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $6
		RETURNING updated_at
	`

	var data, _ = json.Marshal(o.Data)

	err = db.QueryRow(
		query,
		o.ParentId,
		o.Position,
		o.InstanceId,
		o.Name,
		string(data),
		o.Id,
	).Scan(
		&o.UpdatedAt,
	)

	if err != nil {
		return errors.New(fmt.Sprintf("cannot update object: %s", err))
	}

	return nil
}

func DeleteGroupObject(db *sql.DB, o *GroupObject) error {
	const query = `
		DELETE 
		FROM group_objects
		WHERE id = $1
	`

	var err error
	_, err = db.Exec(query, o.Id)

	return err
}

func GetGroupObjectById(db *sql.DB, id uint) (o *GroupObject, err error) {
	const query = `
		SELECT
			id,
			parent_id,
			position,
			group_id,
			object_id,
			instance_id,
			name,
			data,
			locked,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM group_objects
		WHERE
			id = $1
	`

	var data []byte
	o = new(GroupObject)
	err = db.QueryRow(query, id).Scan(
		&o.Id,
		&o.ParentId,
		&o.Position,
		&o.GroupId,
		&o.ObjectId,
		&o.InstanceId,
		&o.Name,
		&data,
		&o.Locked,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &o.Data)

	return o, nil
}

func GetGroupObjectsMapByGroupId(db *sql.DB, groupId uint) (result map[uint]*GroupObject, err error) {
	result = make(map[uint]*GroupObject)

	const query = `
		SELECT
			id,
			parent_id,
			position,
			group_id,
			object_id,
			instance_id,
			name,
			data,
			locked,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM group_objects
		WHERE
			group_id = $1
	`

	rows, err := db.Query(query, groupId)
	if err != nil && err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []byte

	for rows.Next() {
		o := new(GroupObject)
		err = rows.Scan(
			&o.Id,
			&o.ParentId,
			&o.Position,
			&o.GroupId,
			&o.ObjectId,
			&o.InstanceId,
			&o.Name,
			&data,
			&o.Locked,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(data, &o.Data)

		result[o.Id] = o
	}

	return result, nil
}

func GetGroupObjectsTreeByGroupId(db *sql.DB, groupId uint) (result []*GroupObject, err error) {
	groupObjectsByIdMap, err := GetGroupObjectsMapByGroupId(db, groupId)
	if err != nil {
		return result, err
	}

	return BuildGroupObjectsTree(groupObjectsByIdMap)
}

func BuildGroupObjectsTree(groupObjectsByIdMap map[uint]*GroupObject) (result []*GroupObject, err error) {
	result = make([]*GroupObject, 0)

	for _, o := range groupObjectsByIdMap {
		if o.ParentId.Valid {
			parentObject, ok := groupObjectsByIdMap[uint(o.ParentId.Int64)]
			if !ok {
				return nil, errors.New(fmt.Sprintf("Cannot get objects tree. Unknown parent with id: %d", o.ParentId.Int64))
			}
			parentObject.GroupObjects = append(parentObject.GroupObjects, o)
		} else {
			result = append(result, o)
		}
	}

	for _, o := range groupObjectsByIdMap {
		sort.Slice(o.GroupObjects, func(i, j int) bool {
			return o.GroupObjects[i].Position < o.GroupObjects[j].Position
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Position < result[j].Position
	})

	return result, nil
}

func UpdateGroupObjectsLocks(db *sql.DB, groupId uint, lockedInstances []uint) error {
	var err error
	_, err = db.Exec(`
		UPDATE group_objects
		SET locked = false
		WHERE group_id = $1
	`, groupId)

	if err != nil {
		return err
	}

	if len(lockedInstances) > 0 {
		var ids []string
		for _, id := range lockedInstances {
			ids = append(ids, strconv.Itoa(int(id)))
		}

		_, err = db.Exec(fmt.Sprintf(`
			UPDATE group_objects
			SET locked = true
			WHERE 
				group_id = $1
				AND instance_id IN(%s)
		`, strings.Join(ids, ",")), groupId)

		if err != nil {
			return err
		}
	}

	return nil
}
