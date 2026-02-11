package ceng_env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

/*
Envs is a struct containing all the available environment variables available inside the application.
These variables are set by the .env file and overwritten by input (e.g. via Docker compose)
It is possible to set each variable as mandatory as input or optional, in which case it is needed
to define a default value.
*/
type Envs struct {
	DbHost                           string
	DbPort                           int
	DbUsername                       string
	DbPassword                       string
	DbName                           string
	DbSslMode                        string
	DbLogSlowQueryThreshold          int
	AppPort                          int
	AppMode                          string
	AppCorsOrigin                    string
	PrinterEnabled                   bool
	SearchRelevanceThreshold         float64
	PubSubPersistEventsOnDb          bool
	PubSubPersistEventsRetentionDays int
	PubSubSyncMode                   bool
	PickerCorrelationValidityHours   int
	AuthJwtSecret                    string
	AuthJwtAccessTokenDuration       int
	AuthJwtRefreshTokenDuration      int
}

/*
ReadEnvs function reads all the Env Variables.
*/
func ReadEnvs() *Envs {
	godotenv.Load()
	envs := Envs{
		DbHost:                           getMandatoryStringValue("DB_HOST"),
		DbPort:                           getMandatoryIntValue("DB_PORT"),
		DbUsername:                       getMandatoryStringValue("DB_USERNAME"),
		DbPassword:                       getMandatoryStringValue("DB_PASSWORD"),
		DbName:                           getMandatoryStringValue("DB_NAME"),
		DbSslMode:                        getMandatoryStringValue("DB_SSL_MODE"),
		DbLogSlowQueryThreshold:          getMandatoryIntValue("DB_LOG_SLOW_QUERY_THRESHOLD"),
		AppPort:                          getMandatoryIntValue("APP_PORT"),
		AppMode:                          getMandatoryStringValue("APP_MODE"),
		PrinterEnabled:                   getMandatoryBooleanValue("PRINTER_ENABLED"),
		AppCorsOrigin:                    getMandatoryStringValue("APP_CORS_ORIGIN"),
		SearchRelevanceThreshold:         getMandatoryFloatValue("SEARCH_RELEVANCE_THRESHOLD"),
		PubSubPersistEventsOnDb:          getMandatoryBooleanValue("PUBSUB_PERSIST_EVENTS_ON_DB"),
		PubSubPersistEventsRetentionDays: getMandatoryIntValue("PUBSUB_PERSIST_EVENTS_RETENTION_DAYS"),
		PubSubSyncMode:                   getMandatoryBooleanValue("PUBSUB_SYNC_MODE"),
		AuthJwtSecret:                    getMandatoryStringValue("AUTH_JWT_SECRET"),
		AuthJwtAccessTokenDuration:       getMandatoryIntValue("AUTH_JWT_ACCESS_TOKEN_DURATION"),
		AuthJwtRefreshTokenDuration:      getMandatoryIntValue("AUTH_JWT_REFRESH_TOKEN_DURATION"),
	}

	return &envs
}

/*
Read a mandatory integer field, otherwise raise a panic error.
*/
func getMandatoryIntValue(field string) int {
	val := os.Getenv(field)
	if val == "" {
		zap.L().Error(fmt.Sprintf("Missing Mandatory %s field value", field), zap.String("service", "envs-service"))
		panic(fmt.Sprintf("Missing Mandatory %s field value", field))
	}
	intValue, err := strconv.Atoi(val)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Invalid %s field. Value is not an integer", field), zap.String("service", "envs-service"), zap.Error(err))
		panic(fmt.Sprintf("Invalid %s field.  Value is not an integer", field))
	}
	return intValue
}

/*
Read a mandatory string field, otherwise raise a panic error.
*/
func getMandatoryStringValue(field string) string {
	val := os.Getenv(field)
	if val == "" {
		zap.L().Error(fmt.Sprintf("Missing Mandatory %s field value", field), zap.String("service", "envs-service"))
		panic(fmt.Sprintf("Missing Mandatory %s field value", field))
	}
	return val
}

/*
Read a mandatory float field, otherwise raise a panic error.
*/
func getMandatoryFloatValue(field string) float64 {
	val := os.Getenv(field)
	if val == "" {
		zap.L().Error(fmt.Sprintf("Missing Mandatory %s field value", field), zap.String("service", "envs-service"))
		panic(fmt.Sprintf("Missing Mandatory %s field value", field))
	}
	floatValue, err := strconv.ParseFloat(val, 64)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Invalid %s field. Value is not a float", field), zap.String("service", "envs-service"), zap.Error(err))
		panic(fmt.Sprintf("Invalid %s field.  Value is not a float", field))
	}
	return floatValue
}

/*
Read a mandatory boolean field, otherwise raise a panic error.
*/
func getMandatoryBooleanValue(field string) bool {
	val := os.Getenv(field)
	if val == "" {
		zap.L().Error(fmt.Sprintf("Missing Mandatory %s field value", field), zap.String("service", "envs-service"))
		panic(fmt.Sprintf("Missing Mandatory %s field value", field))
	}
	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Invalid %s field. Value is not a boolean", field), zap.String("service", "envs-service"), zap.Error(err))
		panic(fmt.Sprintf("Invalid %s field.  Value is not a boolean", field))
	}
	return boolValue
}
