# GORM Oracle Driver

GORM Oracle driver using [godror](https://github.com/godror/godror)

## Quick Start

```go
import (
	"github.com/gabrielkim13/oracle"
	"gorm.io/gorm"
)

// https://github.com/godror/godror#connection
dsn := `user="USER"
        password="root"
        connectString="localhost:1521/ORCLPDB1"
        libDir="C:\\instantclient_19_16"`
db, err := gorm.Open(oracle.Open(dsn), &gorm.Config{})
```

## Configuration

> Currently, this driver does not accept any configurations other than the DSN string.

To customize the underlying driver's behaviour (i.e. `godror`), please refer to its documentation: 
https://pkg.go.dev/github.com/godror/godror#pkg-overview.

```go
import (
    "github.com/gabrielkim13/oracle"
    "gorm.io/gorm"
)

dsn := `user="USER"
        password="root"
        connectString="localhost:1521/ORCLPDB1"
        libDir="C:\\instantclient_19_16"`
db, err := gorm.Open(oracle.New(oracle.Config{
    DSN: dsn,
}), &gorm.Config{})
```

Checkout [https://gorm.io](https://gorm.io) for details.
