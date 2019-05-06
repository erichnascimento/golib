package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type CleanWatchDog func() error

func NewCleaner(db *sql.DB) *Cleaner {
	return &Cleaner{
		db:       db,
		WatchDog: UnauthorizedWatchDog,
	}
}

type Cleaner struct {
	db       *sql.DB
	WatchDog CleanWatchDog
}

func (c Cleaner) MustDeleteByID(table string, id interface{}) {
	if ok, err := c.canDeleteByID(); !ok {
		panic(fmt.Errorf("unable to delete by id: %v", err))
	}

	_, err := c.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = $1", table), id)
	if err != nil {
		panic(err)
	}
}

func (c Cleaner) canDeleteByID() (bool, error) {
	if c.WatchDog == nil {
		return true, nil
	}

	err := c.WatchDog()
	if err != nil {
		return false, err
	}

	return true, nil
}

func UnauthorizedWatchDog() error {
	return errors.New("UnauthorizedWatchDog can not authorize this operation")
}

func NewAllowedByEnvWatchDog(getEnvFunc func() string, allowedEnvs ...string) CleanWatchDog {
	if getEnvFunc == nil {
		panic("invalid getEnvFunc")
	}

	watchDog := func() error {
		currentEnv := getEnvFunc()
		for _, env := range allowedEnvs {
			if env == currentEnv {
				return nil
			}
		}
		return fmt.Errorf("operation not allowed for environment %s", currentEnv)
	}

	return watchDog
}
