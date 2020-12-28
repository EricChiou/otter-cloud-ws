package user

import (
	"database/sql"
	"otter-cloud-ws/api/common"
	"otter-cloud-ws/bo/userbo"
	"otter-cloud-ws/constants/userstatus"
	"otter-cloud-ws/db/mysql"
	"otter-cloud-ws/jobqueue"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/po/rolepo"
	"otter-cloud-ws/po/userpo"
	"otter-cloud-ws/service/sha3"
	"strconv"

	"github.com/EricChiou/gooq"
)

// Dao user dao
type Dao struct {
	gooq mysql.Gooq
}

// SignUp dao
func (dao *Dao) SignUp(signUp SignUpReqVo) error {
	run := func() interface{} {
		var g mysql.Gooq

		// encrypt password
		encryptPwd := sha3.Encrypt(signUp.Pwd)

		// minio bucket id
		bucketName := minio.GetUserBucketName(signUp.Acc)

		g.SQL.Insert(userpo.Table, userpo.Acc, userpo.Pwd, userpo.Name, userpo.Status, userpo.BucketName).
			Values("?", "?", "?", "?", "?")
		g.AddValues(signUp.Acc, encryptPwd, signUp.Name, userstatus.Active, bucketName)

		if _, err := dao.gooq.Exec(g.SQL.GetSQL(), g.Args...); err != nil {
			return err
		}

		return nil
	}

	return jobqueue.User.NewSignUpJob(run)
}

// SignIn dao
func (dao *Dao) SignIn(signInReqVo SignInReqVo) (userbo.SignInBo, error) {
	var g mysql.Gooq
	var signInBo userbo.SignInBo

	g.SQL.Select(
		userpo.Table+"."+userpo.ID,
		userpo.Acc,
		userpo.Pwd,
		userpo.Table+"."+userpo.Name,
		userpo.RoleCode,
		userpo.Status,
		rolepo.Table+"."+rolepo.Name,
		userpo.BucketName,
	).
		From(userpo.Table).
		Join(rolepo.Table).On(c(userpo.RoleCode).Eq(rolepo.Code)).
		Where(c(userpo.Acc).Eq("?"))
	g.AddValues(signInReqVo.Acc)

	rowMapper := func(row *sql.Row) error {
		if err := row.Scan(
			&signInBo.ID,
			&signInBo.Acc,
			&signInBo.Pwd,
			&signInBo.Name,
			&signInBo.RoleCode,
			&signInBo.Status,
			&signInBo.RoleName,
			&signInBo.BucketName,
		); err != nil {
			return err
		}
		return nil
	}

	// check account existing
	if err := dao.gooq.QueryRow(g.SQL.GetSQL(), rowMapper, g.Args...); err != nil {
		return signInBo, err
	}

	return signInBo, nil
}

// Update dao
func (dao *Dao) Update(updateData UpdateReqVo) error {
	var g mysql.Gooq

	var conditions []gooq.Condition
	if len(updateData.Name) != 0 {
		conditions = append(conditions, c(userpo.Name).Eq("?"))
		g.AddValues(updateData.Name)
	}
	if len(updateData.Pwd) != 0 {
		conditions = append(conditions, c(userpo.Pwd).Eq("?"))
		g.AddValues(sha3.Encrypt(updateData.Pwd))
	}

	g.SQL.Update(userpo.Table).Set(conditions...).Where(c(userpo.ID).Eq("?"))
	g.AddValues(updateData.ID)

	if _, err := dao.gooq.Exec(g.SQL.GetSQL(), g.Args...); err != nil {
		return err
	}

	return nil
}

// List dao
func (dao *Dao) List(listReqVo ListReqVo) (common.PageRespVo, error) {
	index := (listReqVo.Page - 1) * listReqVo.Limit

	var g mysql.Gooq
	g.SQL.
		Select(userpo.ID, userpo.Acc, userpo.Name, userpo.RoleCode, userpo.Status).From(userpo.Table).
		Join("").Lp().
		/**/ Select(userpo.PK).From(userpo.Table).
		/**/ OrderBy(userpo.ID).
		/**/ Limit(strconv.Itoa(index), strconv.Itoa(listReqVo.Limit)).
		Rp().As("t").
		Using(userpo.PK)

	if listReqVo.Active == "true" {
		g.SQL.Where(c(userpo.Status).Eq("?"))
		g.AddValues(string(userstatus.Active))
	}

	list := common.PageRespVo{
		Records: []interface{}{},
		Page:    listReqVo.Page,
		Limit:   listReqVo.Limit,
		Total:   0,
	}
	rowMapper := func(rows *sql.Rows) error {
		for rows.Next() {
			var record ListResVo
			err := rows.Scan(&record.ID, &record.Acc, &record.Name, &record.RoleCode, &record.Status)
			if err != nil {
				return err
			}
			list.Records = append(list.Records, record)
		}
		return nil
	}

	if err := dao.gooq.Query(g.SQL.GetSQL(), rowMapper, g.Args...); err != nil {
		return list, err
	}

	var countG mysql.Gooq
	countG.SQL.Select(f.Count("*")).From(userpo.Table)
	if listReqVo.Active == "true" {
		countG.SQL.Where(c(userpo.Status).Eq("?"))
		countG.AddValues(string(userstatus.Active))
	}

	countRowMapper := func(row *sql.Row) error {
		return row.Scan(&(list.Total))
	}

	if err := dao.gooq.QueryRow(countG.SQL.GetSQL(), countRowMapper, countG.Args...); err != nil {
		return list, err
	}

	return list, nil
}
