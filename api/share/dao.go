package share

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

// GetShareFolder by ownerAcc
func (dao *Dao) GetShareFolder(ownerAcc string) []GetShareFolderResVo {
	var g mysql.Gooq
	g.SQL.
		Select(
			sharedpo.Table+"."+sharedpo.ID,
			sharedpo.SharedAcc,
			userpo.Name,
			sharedpo.Prefix,
			sharedpo.Permission,
		).
		From(sharedpo.Table).
		Join(userpo.Table).On(c(sharedpo.SharedAcc).Eq(userpo.Acc)).
		Where(c(sharedpo.OwnerAcc).Eq("?"))
	g.AddValues(ownerAcc)

	var shareFolderList []GetShareFolderResVo
	g.Query(func(rows *sql.Rows) error {
		if rows.Next() {
			var shareFolder GetShareFolderResVo
			rows.Scan(
				&shareFolder.ID,
				&shareFolder.SharedAcc,
				&shareFolder.SharedName,
				&shareFolder.Prefix,
				&shareFolder.Permission,
			)

			shareFolderList = append(shareFolderList, shareFolder)
		}
		return nil
	})

	return shareFolderList
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

// Get shared folder by user account
func (dao *Dao) Get(sharedAcc string) []GetResVo {
	var g mysql.Gooq
	g.SQL.
		Select(sharedpo.ID, sharedpo.Prefix, sharedpo.Permission, userpo.Acc, userpo.Name).
		From(sharedpo.Table).
		Join(userpo.Table).On(c(sharedpo.OwnerAcc).Eq(userpo.Acc)).
		Where(c(sharedpo.SharedAcc).Eq("?"))
	g.AddValues(sharedAcc)

	var getResVos []GetResVo
	g.Query(func(rows *sql.Rows) error {
		if rows.Next() {
			var getResVo GetResVo
			if err := rows.Scan(
				&getResVo.ID,
				&getResVo.Prefix,
				&getResVo.Permission,
				&getResVo.OwnerAcc,
				&getResVo.OwnerName,
			); err != nil {
				return err
			}
			getResVos = append(getResVos, getResVo)
		}
		return nil
	})

	return getResVos
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
			entity.ID,
			entity.SharedAcc,
			entity.BucketName,
			entity.Prefix,
			entity.Permission,
		)
	})

	if entity.SharedAcc != sharedAcc ||
		strings.Index(prefix, entity.Prefix) != 0 {
		return entity, errors.New("permission denied")
	}

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
