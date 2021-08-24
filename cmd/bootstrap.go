package cmd

import (
	"github.com/pkg/errors"

	"github.com/evalphobia/logrus_sentry"
	"github.com/kopjenmbeng/evermos_online_store/internal/infrastructure"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/jwe_auth"
	"github.com/kopjenmbeng/goconf"
	newrelic "github.com/newrelic/go-agent"
	log "github.com/sirupsen/logrus"
)

var (
	logger    *log.Logger
	Redis     infrastructure.IRedisClient
	telemetry newrelic.Application
	Db        infrastructure.IPostgresSql
	jw        *jwe_auth.JWE
)

const (
	CfgDatabaseRead  = "database.read"
	CfgDatabaseWrite = "database.write"

	CfgRedisHost = "database.redis.master.address"
	CfgRedisPass = "database.redis.master.password"
	CfgRedisDB   = "database.redis.master.db"

	CfgJweKeySrc     = "jwe.key.value"
	CfgJweKeySrcType = "jwe.key.type"
	CfgJweExpire     = "jwe.expire"

	CfgSentryKey   = "sentry_dsn"
	CfgNewRelicKey = "newrelic.key"

	CfgNewRelicDebug = "newrelic.debug"

	TelemetryId = "newrelic.telemetry_id"
)

func init() {
	logger = NewLogger()
	Db = infrastructure.NewPostgresSql(goconf.GetString(CfgDatabaseRead), goconf.GetString(CfgDatabaseWrite), logger)
	Redis = infrastructure.NewRedisClient(goconf.GetStringSlice(CfgRedisHost), goconf.GetString(CfgRedisPass), goconf.GetInt(CfgRedisDB), logger)
	telemetry = NewTelemetry(logger)
	jw = NewJWE(goconf.GetString(CfgJweKeySrc), goconf.GetString(CfgJweKeySrcType), goconf.GetInt(CfgJweExpire))

}

func NewJWE(value string, typ string, exp int) *jwe_auth.JWE {
	key := jwe_auth.Key{Value: value, Type: jwe_auth.KeyStorageType(typ)}
	pk, err := key.GetPrivate()
	if pk == nil || err != nil {
		panic(errors.New("could not get private key"))
	}
	return jwe_auth.NewJWE(pk, exp)
}

func NewLogger() *log.Logger {
	log.SetFormatter(&log.JSONFormatter{})
	l := log.StandardLogger()
	if dsn := goconf.GetString(CfgSentryKey); len(dsn) > 0 {
		hook, err := logrus_sentry.NewSentryHook(dsn, []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		})
		if err == nil {
			hook.StacktraceConfiguration.Enable = true
			hook.SetRelease("partnership@v0.0.1")
			l.Hooks.Add(hook)
		}
	}
	return l
}

func NewTelemetry(l *log.Logger) newrelic.Application {
	k := goconf.GetString(CfgNewRelicKey)
	e := l.WithField("component", "newrelic")
	if len(k) == 0 {
		e.Warnf("configuration %s is not defined", CfgNewRelicKey)
		return nil
	}

	appName := goconf.GetString(TelemetryId)
	config := newrelic.NewConfig(appName, k)
	if isDebug := goconf.Config().GetBool(CfgNewRelicDebug); isDebug {
		l.SetLevel(log.DebugLevel)
	}

	app, err := newrelic.NewApplication(config)
	if err != nil {
		e.Info(errors.Cause(err))
		return nil
	}
	return app
}
