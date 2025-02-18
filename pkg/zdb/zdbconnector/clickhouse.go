package zdbconnector

import (
	"fmt"
	clickhouse2 "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/zondax/golem/pkg/zdb/zdbconfig"
	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"strings"
)

const httpsProtocol = "https"

type ClickHouseConnector struct{}

func (c *ClickHouseConnector) Connect(config *zdbconfig.Config) (*gorm.DB, error) {
	dsn := buildClickhouseDSN(config.ConnectionParams)
	gormConfig := zdbconfig.BuildGormConfig(config.LogConfig)

	var dbConn *gorm.DB
	var err error
	dbConn, err = gorm.Open(clickhouse.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func buildClickhouseDSN(params zdbconfig.ConnectionParams) string {
	protocol := ""

	switch {
	case strings.EqualFold(params.Protocol, clickhouse2.HTTP.String()):
		protocol = clickhouse2.HTTP.String()
	case strings.EqualFold(params.Protocol, httpsProtocol):
		protocol = httpsProtocol
	case strings.EqualFold(params.Protocol, clickhouse2.Native.String()), params.Protocol == "":
		protocol = clickhouse2.Native.String()
	default:
		zap.S().Errorf("Failed to identify connection protocol [%s]", params.Protocol)
	}

	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%v/%s",
		protocol,
		params.User,
		params.Password,
		params.Host,
		params.Port,
		params.Name,
	)

	if params.Params != "" {
		dsn = fmt.Sprintf("%s?%s", dsn, params.Params)
	}

	return dsn
}

func (c *ClickHouseConnector) VerifyConnection(db *gorm.DB) error {
	return db.Exec("SELECT 1").Error
}
