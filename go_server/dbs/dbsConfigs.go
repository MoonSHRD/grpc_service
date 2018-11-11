package dbs

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
)

type dbConfigs struct {
    MysqlConfig configData `json:"mysql"`
    RedisConfig configData `json:"redis"`
    MongoConfig configData `json:"Mongo"`
}

type configData struct {
    User     string `json:"user"`
    Password string `json:"password"`
    Db       string `json:"db"`
    Host     string `json:"host"`
    Port     string `json:"port"`
}

func getConfig(db string) configData {
    jsonFile, err := os.Open("go_server/database.json")
    if err != nil {
        log.Fatal("Create database.json")
    }
    
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var conf dbConfigs
    json.Unmarshal(byteValue, &conf)
    var ret configData
    switch db {
    case "Mongo":
        ret = conf.MongoConfig
    case "mysql":
        ret = conf.MysqlConfig
    case "redis":
        ret = conf.RedisConfig
    }
    return ret
}