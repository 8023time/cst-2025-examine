// package SQL

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// var dbPool *pgxpool.Pool

// func ConnectSQL() (*pgxpool.Pool, error) {
// 	connUrl := "postgres://testuser:123456@localhost:6900/testdb"

// 	if dbPool == nil {
// 		var err error
// 		dbPool, err = pgxpool.New(context.Background(), connUrl)
// 		if err != nil {
// 			return nil, fmt.Errorf("无法连接到数据库: %w", err)
// 		}
// 		log.Println("成功连接到数据库")
// 	}

// 	return dbPool, nil
// }

// db/db.go
package SQL

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// ConnectDB 函数返回一个数据库连接
func ConnectSQL() (*pgx.Conn, error) {
	connUrl := "postgres://testuser:123456@localhost:6900/testdb"
	conn, err := pgx.Connect(context.Background(), connUrl) // context.Background()：通常用于没有特定上下文的情况下，作为根上下文传递
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return conn, nil
}
