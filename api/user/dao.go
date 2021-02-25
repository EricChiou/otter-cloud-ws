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
type Dao struct{}

// CheckAccExisting dao
func (dao *Dao) CheckAccExisting(acc string) (bool, error) {
	var g mysql.Gooq

	g.SQL.Select(userpo.Acc).From(userpo.Table).Where(c(userpo.Acc).Eq("?"))
	g.AddValues(acc)

	if err := g.QueryRow(func(row *sql.Row) error {
		var acc string
		return row.Scan(&acc)
	}); err != nil {
		return false, err
	}

	return true, nil
}

// SignUp dao
func (dao *Dao) SignUp(signUp SignUpReqVo, activeCode string) error {
	run := func() interface{} {
		var g mysql.Gooq

		// encrypt password
		encryptPwd := sha3.Encrypt(signUp.Pwd)

		// minio bucket id
		bucketName := minio.GetUserBucketName(signUp.Acc)

		g.SQL.
			Insert(
				userpo.Table,
				userpo.Acc,
				userpo.Pwd,
				userpo.Name,
				userpo.Status,
				userpo.BucketName,
				userpo.ActiveCode,
			).
			Values("?", "?", "?", "?", "?", "?")
		g.AddValues(signUp.Acc, encryptPwd, signUp.Name, userstatus.Inactive, bucketName, activeCode)

		if _, err := g.Exec(); err != nil {
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
	if err := g.QueryRow(rowMapper); err != nil {
		return signInBo, err
	}

	return signInBo, nil
}

// Update dao
func (dao *Dao) Update(updateData UpdateReqVo, acc string) error {
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

	g.SQL.Update(userpo.Table).Set(conditions...).Where(c(userpo.Acc).Eq("?"))
	g.AddValues(acc)

	if _, err := g.Exec(); err != nil {
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

	if err := g.Query(rowMapper); err != nil {
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

	if err := countG.QueryRow(countRowMapper); err != nil {
		return list, err
	}

	return list, nil
}

// GetBucketName by user account
func (dao *Dao) GetBucketName(acc string) (string, error) {
	var g mysql.Gooq

	g.SQL.Select(userpo.BucketName).From(userpo.Table).Where(c(userpo.Acc).Eq("?"))
	g.AddValues(acc)

	var bucketName string
	err := g.QueryRow(func(row *sql.Row) error {
		return row.Scan(&bucketName)
	})
	if err != nil {
		return "", err
	}

	return bucketName, nil
}

// GetUserFuzzyList by keyword
func (dao *Dao) GetUserFuzzyList(keyword string) ([]string, error) {
	var g mysql.Gooq

	g.SQL.
		Select(userpo.Acc).
		From(userpo.Table).
		Where(c(userpo.Acc).Like("?"))

	g.AddValues(keyword + "%")

	var accountList []string
	err := g.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var account string
			if err := rows.Scan(&account); err == nil {
				accountList = append(accountList, account)
			}
		}
		return nil
	})

	return accountList, err
}

// ActivateAcc by active code
func (dao *Dao) ActivateAcc(activeCode string) error {
	var g mysql.Gooq

	g.SQL.
		Select(userpo.Acc).
		From(userpo.Table).
		Where(c(userpo.ActiveCode).Eq("?"))
	g.AddValues(activeCode)

	var account string
	if err := g.QueryRow(func(row *sql.Row) error {
		return row.Scan(&account)
	}); err != nil {
		return err
	}

	g = mysql.Gooq{}
	conditions := []gooq.Condition{c(userpo.Status).Eq("?")}
	g.SQL.Update(userpo.Table).Set(conditions...).Where(c(userpo.Acc).Eq("?"))
	g.AddValues(userstatus.Active, account)

	if _, err := g.Exec(); err != nil {
		return err
	}

	return nil
}

// SendActivationCode by account
func (dao *Dao) SendActivationCode(account, activeCode string) (userName string, err error) {
	var g mysql.Gooq

	conditions := []gooq.Condition{c(userpo.ActiveCode).Eq("?")}
	g.SQL.Update(userpo.Table).Set(conditions...).Where(c(userpo.Acc).Eq("?"))
	g.AddValues(activeCode, account)

	if _, err = g.Exec(); err != nil {
		return "", err
	}

	g = mysql.Gooq{}
	g.SQL.Select(userpo.Name).From(userpo.Table).Where(c(userpo.Acc).Eq("?"))
	g.AddValues(account)

	err = g.QueryRow(func(row *sql.Row) error {
		return row.Scan(&userName)
	})

	return userName, err
}

// ResetPwd by account
func (dao *Dao) ResetPwd(account, newPwd string) (userName string, err error) {
	var g mysql.Gooq

	// encrypt password
	encryptNewPwd := sha3.Encrypt(newPwd)

	conditions := []gooq.Condition{c(userpo.Pwd).Eq("?")}
	g.SQL.Update(userpo.Table).Set(conditions...).Where(c(userpo.Acc).Eq("?"))
	g.AddValues(encryptNewPwd, account)

	if _, err := g.Exec(); err != nil {
		return "", err
	}

	g = mysql.Gooq{}
	g.SQL.Select(userpo.Name).From(userpo.Table).Where(c(userpo.Acc).Eq("?"))
	g.AddValues(account)

	err = g.QueryRow(func(row *sql.Row) error {
		return row.Scan(&userName)
	})

	return userName, err
}
