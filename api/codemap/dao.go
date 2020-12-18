package codemap

import (
	"database/sql"
	"otter-cloud-ws/api/common"
	"otter-cloud-ws/db/mysql"
	"otter-cloud-ws/jobqueue"
	"otter-cloud-ws/po/codemappo"
	"strconv"
	"strings"
)

// Dao codemap dao
type Dao struct {
	gooq mysql.Gooq
}

// Add add codemap dao
func (dao *Dao) Add(addReqVo AddReqVo) error {
	run := func() interface{} {
		var g mysql.Gooq
		g.SQL.Insert(codemappo.Table, codemappo.Type, codemappo.Code, codemappo.Name, codemappo.SortNo, codemappo.Enable).
			Values("?", "?", "?", "?", "?")
		g.AddValues(addReqVo.Type, addReqVo.Code, addReqVo.Name, addReqVo.SortNo, addReqVo.Enable)

		if _, err := dao.gooq.Exec(g.SQL.GetSQL(), g.Args...); err != nil {
			return err
		}

		return nil
	}

	return jobqueue.Codemap.NewAddJob(run)
}

// Update update codemap dao
func (dao *Dao) Update(updateReqVo UpdateReqVo) error {
	var g mysql.Gooq
	g.SQL.Update(codemappo.Table).
		Set(c(codemappo.Code).Eq("?"), c(codemappo.Name).Eq("?"), c(codemappo.Type).Eq("?"), c(codemappo.SortNo).Eq("?"), c(codemappo.Enable).Eq("?")).
		Where(c(codemappo.ID).Eq("?"))
	g.AddValues(updateReqVo.Code, updateReqVo.Name, updateReqVo.Type, updateReqVo.SortNo, updateReqVo.Enable)
	g.AddValues(updateReqVo.ID)

	_, err := dao.gooq.Exec(g.SQL.GetSQL(), g.Args...)
	if err != nil {
		return err
	}

	return nil
}

// Delete update codemap dao
func (dao *Dao) Delete(deleteReqVo DeleteReqVo) error {
	var g mysql.Gooq
	g.SQL.Delete(codemappo.Table).Where(c(codemappo.ID).Eq("?"))
	g.AddValues(deleteReqVo.ID)

	_, err := dao.gooq.Exec(g.SQL.GetSQL(), g.Args...)
	if err != nil {
		return err
	}

	return nil
}

// List get codemap list
func (dao *Dao) List(listReqVo ListReqVo) (common.PageRespVo, error) {
	index := (listReqVo.Page - 1) * listReqVo.Limit

	var g mysql.Gooq
	g.SQL.Select(codemappo.ID, codemappo.Type, codemappo.Code, codemappo.Name, codemappo.SortNo, codemappo.Enable).From(codemappo.Table).
		Join("").Lp().
		/**/ Select(codemappo.PK).
		/**/ From(codemappo.Table).
		/**/ OrderBy(codemappo.ID).
		/**/ Limit(strconv.Itoa(index), strconv.Itoa(listReqVo.Limit)).
		Rp().As("t").
		Using(codemappo.PK)

	var where mysql.Gooq
	if len(listReqVo.Type) > 0 {
		where.SQL.And(c(codemappo.Type).Eq("?"))
		g.AddValues(listReqVo.Type)
	}
	if listReqVo.Enable == "true" {
		where.SQL.And(c(codemappo.Enable).Eq(true))
	}
	g.SQL.Add("WHERE" + strings.TrimPrefix(where.SQL.GetSQL(), "AND"))

	list := common.PageRespVo{
		Records: []interface{}{},
		Page:    listReqVo.Page,
		Limit:   listReqVo.Limit,
		Total:   0,
	}
	rowMapper := func(rows *sql.Rows) error {
		for rows.Next() {
			var record ListResVo
			err := rows.Scan(&record.ID, &record.Type, &record.Code, &record.Name, &record.SortNo, &record.Enable)
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

	var count mysql.Gooq
	count.SQL.Select(f.Count("*")).From(codemappo.Table)
	if listReqVo.Enable == "true" {
		count.SQL.Where(c(codemappo.Enable).Eq(true))
	}
	countRowMapper := func(row *sql.Row) error {
		return row.Scan(&(list.Total))
	}

	if err := dao.gooq.QueryRow(count.SQL.GetSQL(), countRowMapper); err != nil {
		return list, err
	}

	return list, nil
}
