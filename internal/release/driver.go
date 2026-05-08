package release

import (
	"database/sql"
	"os"
	"strings"

	hullerr "github.com/ebogdum/hull/internal/errors"
	"k8s.io/client-go/kubernetes"
)

// SelectStorage chooses a storage driver based on HULL_DRIVER (with fallback
// to HELM_DRIVER for compatibility). Recognized values:
//
//   - "secret" / ""           Kubernetes Secrets (default)
//   - "configmap"             Kubernetes ConfigMaps
//   - "memory"                in-process (tests, dry runs)
//   - "sql"                   SQL backend; requires HULL_DRIVER_SQL_DRIVER
//                             (e.g. "postgres", "sqlite", "mysql") and
//                             HULL_DRIVER_SQL_DSN (driver-specific DSN).
//
// For the SQL driver the caller's binary must already import a database/sql
// driver matching HULL_DRIVER_SQL_DRIVER (e.g. github.com/lib/pq for postgres).
func SelectStorage(clientset kubernetes.Interface, namespace string) (Storage, error) {
	driver := os.Getenv("HULL_DRIVER")
	if "" == driver {
		driver = os.Getenv("HELM_DRIVER")
	}
	switch strings.ToLower(driver) {
	case "", "secret", "secrets":
		return NewSecretStorage(clientset, namespace), nil
	case "configmap", "configmaps":
		return NewConfigMapStorage(clientset, namespace), nil
	case "memory":
		return NewMemoryStorage(), nil
	case "sql":
		return openSQLFromEnv(namespace)
	}
	return nil, hullerr.NewErrorf(hullerr.ErrCLIValidation, "unknown HULL_DRIVER %q", driver)
}

func openSQLFromEnv(namespace string) (Storage, error) {
	sqlDriver := os.Getenv("HULL_DRIVER_SQL_DRIVER")
	if "" == sqlDriver {
		sqlDriver = os.Getenv("HELM_DRIVER_SQL_USER_DATABASE_DRIVER")
	}
	dsn := os.Getenv("HULL_DRIVER_SQL_DSN")
	if "" == dsn {
		dsn = os.Getenv("HELM_DRIVER_SQL_CONNECTION_STRING")
	}
	if "" == sqlDriver || "" == dsn {
		return nil, hullerr.NewError(hullerr.ErrCLIValidation,
			"SQL driver requires HULL_DRIVER_SQL_DRIVER and HULL_DRIVER_SQL_DSN env vars")
	}
	db, err := sql.Open(sqlDriver, dsn)
	if nil != err {
		return nil, hullerr.WrapErrorf(hullerr.ErrCLIValidation, err, "failed to open SQL driver %q", sqlDriver)
	}
	dialect := DialectStandard
	if "postgres" == strings.ToLower(sqlDriver) || "pgx" == strings.ToLower(sqlDriver) {
		dialect = DialectPostgres
	}
	return NewSQLStorageWithDialect(db, namespace, dialect)
}
