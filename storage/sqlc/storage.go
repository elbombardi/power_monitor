package storage

import (
	"context"
	"database/sql"
	"log"

	"github.com/elbombardi/power_monitor/models"
	"github.com/elbombardi/power_monitor/util"
	_ "github.com/lib/pq"
)

type DataStore interface {
	Store(ctx context.Context, powerStateChange *models.PowerStateChange) error
}

type PostgresStore struct {
	*Queries
	db *sql.DB
}

func (s *PostgresStore) Store(ctx context.Context, powerStateChange *models.PowerStateChange) error {
	_, err := s.InsertPowerLog(ctx, InsertPowerLogParams{
		Timestamp:  powerStateChange.Time,
		PowerState: string(powerStateChange.PowerState),
	})
	return err
}

var dbInstance *sql.DB

func GetStoreInstance() (DataStore, error) {
	if dbInstance == nil {
		log.Println("Initializing database conntection..")
		var err error
		dbInstance, err = sql.Open(util.ConfigDBDriver(), *util.ConfigDBSource())
		if err != nil {
			return nil, err
		}
		err = dbInstance.Ping()
		if err != nil {
			return nil, err
		}
		value, _ := util.ConfigDBMaxIdleConns()
		dbInstance.SetMaxIdleConns(value)
		value, _ = util.ConfigDBMaxOpenConns()
		dbInstance.SetMaxOpenConns(value)
		duration, _ := util.ConfigDBConnMaxIdleTime()
		dbInstance.SetConnMaxIdleTime(duration)
		duration, _ = util.ConfigDBConnMaxLifeTime()
		dbInstance.SetConnMaxLifetime(duration)
	}

	return &PostgresStore{
		db:      dbInstance,
		Queries: &Queries{db: dbInstance},
	}, nil
}

func Finalize() error {
	if dbInstance == nil {
		return nil
	}
	return dbInstance.Close()
}
