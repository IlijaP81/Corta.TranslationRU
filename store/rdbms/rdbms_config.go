package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/store"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// persistance layer
//
// all functions go under one struct
//   why? because it will be easier to initialize and pass around
//
// each domain will be in it's own file
//
// connection logic will be built in the persistence layer (making pkg/db obsolete)
//

type (
	txRetryOnErrHandler func(int, error) bool
	columnPreprocFn     func(string, string) string
	valuePreprocFn      func(interface{}, string) interface{}
	errorHandler        func(error) error
	triggerKey          string

	rowScanner interface {
		Scan(...interface{}) error
	}

	TriggerHandlers map[triggerKey]interface{}

	Config struct {
		DriverName     string
		DataSourceName string
		DBName         string

		// MaxOpenConns sets maximum number of open connections to the database
		// defaults to same value as set in the db/sql
		MaxOpenConns int

		// ConnMaxLifetime sets the maximum amount of time a connection may be reused
		// defaults to same value as set in the db/sql
		ConnMaxLifetime time.Duration

		// MaxIdleConns sets the maximum number of connections in the idle connection pool
		// defaults to same value as set in the db/sql
		MaxIdleConns int

		// ConnTryPatience sets time window in which we do not complaining about failed connection tries
		ConnTryPatience time.Duration

		// ConnTryBackoffDelay sets backoff delay after failed try
		ConnTryBackoffDelay time.Duration

		// ConnTryTimeout sets timeout per try
		ConnTryTimeout time.Duration

		// ConnTryMax maximum number of retrys for getting the connection
		ConnTryMax int

		// PlaceholderFormat used by squirrel query generator
		PlaceholderFormat squirrel.PlaceholderFormat

		// Disable transactions
		TxDisabled bool

		// How many times should we retry failed transaction?
		TxMaxRetries int

		// TxRetryErrHandler should return true if transaction should be retried
		//
		// Because retry algorithm varies between concrete rdbms implementations
		//
		// Handler must return true if failed transaction should be replied
		// and false if we're safe to terminate it
		TxRetryErrHandler txRetryOnErrHandler

		ColumnPreprocessors map[string]columnPreprocFn
		ValuePreprocessors  map[string]valuePreprocFn

		ErrorHandler errorHandler

		// Implementations can override internal RDBMS row scanners
		RowScanners map[string]interface{}

		// Different store backend implementation might handle upsert differently...
		UpsertBuilder func(*Config, string, store.Payload, ...string) (squirrel.InsertBuilder, error)

		// TriggerHandlers handle various exceptions that can not be handled generally within RDBMS package.
		// see triggerKey type and defined constants to see where the hooks are and how can they be called
		TriggerHandlers TriggerHandlers

		// UniqueConstraintCheck flag controls if unique constraints should be explicitly checked within
		// store or is this handled inside the storage
		//
		//
		UniqueConstraintCheck bool

		// FunctionHandler takes care of translation & transformation of (sql) functions
		// and their parameters
		//
		// Functions are used in filters and aggregations
		SqlFunctionHandler func(f ql.Function) (ql.ASTNode, error)

		CastModuleFieldToColumnType func(field ModuleFieldTypeDetector, ident ql.Ident) (ql.Ident, error)
	}
)

var (
	dsnMasker = regexp.MustCompile("(.)(?:.*)(.):(.)(?:.*)(.)@")
)

// MaskedDSN replaces username & password from DSN string so it's usable for logging
func (c *Config) MaskedDSN() string {
	return dsnMasker.ReplaceAllString(c.DataSourceName, "$1****$2:$3****$4@")
}

func (c *Config) SetDefaults() {
	if c.PlaceholderFormat == nil {
		c.PlaceholderFormat = squirrel.Question
	}

	if c.TxMaxRetries == 0 {
		c.TxMaxRetries = TxRetryHardLimit
	}

	if c.TxRetryErrHandler == nil {
		// Default transaction retry handler
		c.TxRetryErrHandler = TxNoRetry
	}

	if c.ErrorHandler == nil {
		c.ErrorHandler = ErrHandlerFallthrough
	}

	if c.UpsertBuilder == nil {
		c.UpsertBuilder = UpsertBuilder
	}

	// ** ** ** ** ** ** ** ** ** ** ** ** ** **

	if c.MaxIdleConns == 0 {
		// Same as default in the db/sql
		c.MaxIdleConns = 32
	}

	if c.MaxOpenConns == 0 {
		// Same as default in the db/sql
		c.MaxOpenConns = 256
	}

	if c.ConnMaxLifetime == 0 {
		// Same as default in the db/sql
		c.ConnMaxLifetime = 10 * time.Minute
	}

	// ** ** ** ** ** ** ** ** ** ** ** ** ** **

	if c.ConnTryPatience == 0 {
		c.ConnTryPatience = 1 * time.Minute
	}

	if c.ConnTryBackoffDelay == 0 {
		c.ConnTryBackoffDelay = 10 * time.Second
	}

	if c.ConnTryTimeout == 0 {
		c.ConnTryTimeout = 30 * time.Second
	}

	if c.ConnTryMax == 0 {
		c.ConnTryMax = 99
	}

	if c.TriggerHandlers == nil {
		c.TriggerHandlers = TriggerHandlers{}
	}
}

// ParseExtra parses extra params (params starting with *)
// from DSN's querystring (after ?)
func (c *Config) ParseExtra() (err error) {
	// Make sure we only got qs
	const q = "?"
	var (
		dsn = c.DataSourceName
		qs  string
	)

	if pos := strings.LastIndex(dsn, q); pos == -1 {
		return nil
	} else {
		// Trim qs from DSN, we'll re-attach the remaining params
		c.DataSourceName, qs = dsn[:pos], dsn[pos+1:]
	}

	var vv url.Values
	if vv, err = url.ParseQuery(qs); err != nil {
		return err
	}

	var (
		val string

		parseInt = func(s string) (int, error) {
			if tmp, err := strconv.ParseInt(s, 10, 32); err != nil {
				return 0, err
			} else {
				return int(tmp), nil
			}

		}
	)

	for key := range vv {
		val = vv.Get(key)
		switch key {
		case "*connTryPatience":
			delete(vv, key)
			if c.ConnTryPatience, err = time.ParseDuration(val); err != nil {
				return
			}

		case "*connTryBackoffDelay":
			delete(vv, key)
			if c.ConnTryBackoffDelay, err = time.ParseDuration(val); err != nil {
				return
			}

		case "*connTryTimeout":
			delete(vv, key)
			if c.ConnTryTimeout, err = time.ParseDuration(val); err != nil {
				return
			}

		case "*connMaxTries":
			delete(vv, key)
			if c.ConnTryMax, err = parseInt(val); err != nil {
				return
			}

		case "*connMaxOpen":
			delete(vv, key)
			if c.MaxOpenConns, err = parseInt(val); err != nil {
				return
			}

		case "*connMaxLifetime":
			delete(vv, key)
			if c.ConnMaxLifetime, err = time.ParseDuration(val); err != nil {
				return
			}

		case "*connMaxIdle":
			delete(vv, key)
			if c.MaxIdleConns, err = parseInt(val); err != nil {
				return
			}
		}
	}

	// Encode QS back to DSN
	c.DataSourceName += q + vv.Encode()

	return nil
}
