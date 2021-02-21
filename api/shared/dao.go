package shared

import (
	"database/sql"
	"errors"
	"otter-cloud-ws/db/mysql"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/po/sharedpo"
	"otter-cloud-ws/po/userpo"
	"strings"
)

// Dao shared dao
type Dao struct{}

// CheckFolder existing
func (dao *Dao) CheckFolder(bucketName, prefix string) error {
	objList := minio.ListObjects(bucketName, prefix, false)
	if len(objList) > 0 {
		return nil
	}

	return errors.New("the folder not existing")
}

// CheckShare not duplicated
func (dao *Dao) CheckShare(ownerAcc, sharedAcc, prefix string) error {
	var g mysql.Gooq
	g.SQL.
		Select(sharedpo.ID).
		From(sharedpo.Table).
		Where(c(sharedpo.OwnerAcc).Eq("?")).
		And(c(sharedpo.SharedAcc).Eq("?")).
		And(c(sharedpo.Prefix).Eq("?"))
	g.AddValues(ownerAcc, sharedAcc, prefix)

	return g.QueryRow(func(row *sql.Row) error {
		var id int
		return row.Scan(&id)
	})
}

// Add share folder
func (dao *Dao) Add(ownerAcc, sharedAcc, bucketName, prefix, permission string) error {
	var g mysql.Gooq
	g.SQL.
		Insert(
			sharedpo.Table,
			sharedpo.OwnerAcc,
			sharedpo.SharedAcc,
			sharedpo.BucketName,
			sharedpo.Prefix,
			sharedpo.Permission,
		).
		Values("?", "?", "?", "?", "?")
	g.AddValues(ownerAcc, sharedAcc, bucketName, prefix, permission)

	if _, err := g.Exec(); err != nil {
		return err
	}

	return nil
}

// GetSharedFolder by user acc
func (dao *Dao) GetSharedFolder(userAcc string) []GetSharedFolderResVo {
	var g mysql.Gooq
	g.SQL.
		Select(
			sharedpo.Table+"."+sharedpo.ID,
			sharedpo.OwnerAcc,
			"owner."+userpo.Name+" as owner_name",
			sharedpo.SharedAcc,
			"shared."+userpo.Name+" as shared_name",
			sharedpo.Prefix,
			sharedpo.Permission,
		).
		From(sharedpo.Table).
		Join(userpo.Table).As("owner").On(c(sharedpo.OwnerAcc).Eq("owner." + userpo.Acc)).
		Join(userpo.Table).As("shared").On(c(sharedpo.SharedAcc).Eq("shared." + userpo.Acc)).
		Where(c(sharedpo.OwnerAcc).Eq("?")).Or(c(sharedpo.SharedAcc).Eq("?"))
	g.AddValues(userAcc, userAcc)

	var sharedFolderList []GetSharedFolderResVo
	g.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var sharedFolder GetSharedFolderResVo
			rows.Scan(
				&sharedFolder.ID,
				&sharedFolder.OwnerAcc,
				&sharedFolder.OwnerName,
				&sharedFolder.SharedAcc,
				&sharedFolder.SharedName,
				&sharedFolder.Prefix,
				&sharedFolder.Permission,
			)

			if userAcc == sharedFolder.SharedAcc {
				perfixSep := strings.SplitAfter(sharedFolder.Prefix, "/")
				sharedFolder.Prefix = perfixSep[len(perfixSep)-2]
			}

			sharedFolderList = append(sharedFolderList, sharedFolder)
		}
		return nil
	})

	return sharedFolderList
}

// Remove shared folder
func (dao *Dao) Remove(sharedID int, ownerAcc string) error {
	var g mysql.Gooq
	g.SQL.
		Delete(sharedpo.Table).
		Where(c(sharedpo.ID).Eq("?")).And(c(sharedpo.OwnerAcc).Eq("?"))
	g.AddValues(sharedID, ownerAcc)

	if _, err := g.Exec(); err != nil {
		return err
	}

	return nil
}

// CheckPermission by shared id , shared acc and prefix
func (dao *Dao) CheckPermission(sharedID int, sharedAcc, prefix string) (sharedpo.Entity, error) {
	var g mysql.Gooq
	g.SQL.
		Select(
			sharedpo.ID,
			sharedpo.SharedAcc,
			sharedpo.BucketName,
			sharedpo.Prefix,
			sharedpo.Permission,
		).
		From(sharedpo.Table).
		Where(c(sharedpo.ID).Eq("?"))
	g.AddValues(sharedID)

	var entity sharedpo.Entity
	err := g.QueryRow(func(row *sql.Row) error {
		return row.Scan(
			&entity.ID,
			&entity.SharedAcc,
			&entity.BucketName,
			&entity.Prefix,
			&entity.Permission,
		)
	})

	entityPrefixSep := strings.SplitAfter(entity.Prefix, "/")
	prefixSep := strings.SplitAfterN(prefix, "/", 2)
	if entity.SharedAcc != sharedAcc || entityPrefixSep[len(entityPrefixSep)-2] != prefixSep[0] {
		return entity, errors.New("permission denied")
	}

	entity.Prefix = entity.Prefix + prefixSep[1]

	return entity, err
}

// CheckWritePermission by shared id , shared acc and prefix
func (dao *Dao) CheckWritePermission(sharedID int, sharedAcc, prefix string) (sharedpo.Entity, error) {
	sharedEntity, err := dao.CheckPermission(sharedID, sharedAcc, prefix)
	if err != nil {
		return sharedEntity, err
	}

	if sharedEntity.Permission != "write" {
		return sharedEntity, errors.New("permission denied")
	}

	return sharedEntity, nil
}
