package dbs

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type mysql struct {
    instance *sql.DB
}

//type dbConfigs struct {
//    MysqlConfig mysqlConfig `json:"mysql"`
//    RedisConfig redisConfig `json:"redis"`
//}
//
//type mysqlConfig struct {
//    User     string `json:"user"`
//    Password string `json:"password"`
//    Db       string `json:"db"`
//}
//
//type redisConfig struct {
//}

var mysqlInstance *mysql

func GetMysql() *sql.DB {
    if mysqlInstance == nil {
        conf:=getConfig("mysql")
        db, err := sql.Open(
            "mysql",
            fmt.Sprintf("%s:%s@/%s",
                conf.User,
                conf.Password,
                conf.Db,
            ),
        )
        if err != nil {
            log.Fatal(err)
        }
        mysqlInstance = &mysql{instance: db}
        //defer db.Close()
    }
    return mysqlInstance.instance
}

//func InsertUser(user moonshard.User) bool {
//    db := GetMysql()
//    stmtIns, err := db.Prepare("INSERT INTO user VALUES( ?, ? )")
//    if err!=nil {
//        log.Println("Got error inserting user: ", err)
//        return false
//    }
//    defer stmtIns.Close()
//    stmtIns.Exec(user.Id, user.Data)
//    return true
//}

func GetUser(id string) string {
    db := GetMysql()
    stmtOut, _ := db.Prepare("select data from user where user.id = ?")
    defer stmtOut.Close()
    var data string
    stmtOut.QueryRow(id).Scan(&data)
    return data
}
