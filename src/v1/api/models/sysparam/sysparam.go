package sysparam

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (s *SysParam) Update(v string) (*sql.Result, error) {

	sql := `
		UPDATE sysparam
			SET value = ?
		WHERE ` + "`key`" + ` = ?
	`
	sqlResult, err := database.DBInstancePublic.Exec(sql, v, s.Key)

	return &sqlResult, err
}

func (s *SysParam) GetAll() (*ResponseSysParam, error) {
	var responseSysParam ResponseSysParam

	sql := "SELECT `key`, `value` FROM sysparam"
	err := database.DBInstancePublic.Select(&responseSysParam, sql)

	fmt.Println("err", err)

	return &responseSysParam, err
}

func (s *SysParam) GetByKey(k string) (*SysParam, error) {
	var sysParam []SysParam

	sql := "SELECT `key`, `value` FROM sysparam WHERE `key` = ?"
	err := database.DBInstancePublic.Select(&sysParam, sql, k)

	if err != nil {
		return nil, err
	}

	if len(sysParam) == 0 {
		return nil, errors.New("no value found for key: " + k)
	}

	if sysParam[0].Value == "" {
		return &sysParam[0], errors.New("no value found for key: " + k)
	}

	return &sysParam[0], err
}