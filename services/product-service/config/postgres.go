package config

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/daffaromero/retries/services/common/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	host               = utils.GetEnv("DB_HOST")
	port               = utils.GetEnv("DB_PORT")
	username           = utils.GetEnv("DB_USERNAME")
	password           = utils.GetEnv("DB_PASSWORD")
	dbName             = utils.GetEnv("DB_NAME")
	minConns           = utils.GetEnv("DB_MIN_CONNS")
	maxConns           = utils.GetEnv("DB_MAX_CONNS")
	TimeOutDuration, _ = strconv.Atoi(utils.GetEnv("DB_CONNECTION_TIMEOUT"))
)

func NewPGDatabase() *pgxpool.Pool {
	conn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, dbName)

	poolConf, err := pgxpool.ParseConfig(conn)
	if err != nil {
		log.Print("failed to parse conn string", conn)
	}

	minConnsInt, err := strconv.Atoi(minConns)
	if err != nil {
		log.Print("expected DB_MIN_CONNS to be int", minConns)
	}

	maxConnsInt, err := strconv.Atoi(maxConns)
	if err != nil {
		log.Print("expected DB_MAX_CONNS to be int", minConns)
	}

	poolConf.MinConns = int32(minConnsInt)
	poolConf.MaxConns = int32(maxConnsInt)
	poolConf.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeDescribeExec

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConf)
	if err != nil {
		log.Print("failed to apply pool configuration", conn)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		log.Print(err)
	}

	log.Print("database connected", conn)

	return pool
}
