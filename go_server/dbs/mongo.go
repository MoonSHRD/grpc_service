package dbs

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"grpc_service/go_server/models"
	"log"
)

type Mongo struct {
	Instance *mgo.Database
}

var mongoInstance *Mongo

var limit = 15

func GetMongo() *Mongo {
	if mongoInstance == nil {
		conf := getConfig("Mongo")
		mongo_db, err := mgo.Dial(
			fmt.Sprintf("%s:%s",
				conf.Host,
				conf.Port,
			),
		)
		if err != nil {
			log.Fatal(err)
		}
		db := mongo_db.DB(conf.Db)
		mongoInstance = &Mongo{Instance: db}
		//defer db.Close()
	}
	return mongoInstance
}

func (db Mongo) GetUsers(ptr int, str string) []models.User {
	skip := int(ptr) * limit
	var users []models.User
	mongoInstance.Instance.C("users").Find(bson.M{"$or": []bson.M{
		{"firstname": bson.M{"$regex": bson.RegEx{str, ""}}},
		{"lastname": bson.M{"$regex": bson.RegEx{str, ""}}},
	}}).Skip(skip).Limit(limit).All(&users)
	return users
}

func (db Mongo) GetCommunities(ptr int, str string) []models.Community {
	skip := int(ptr) * limit
	var communities []models.Community
	mongoInstance.Instance.C("communities").Find(bson.M{
		"name": bson.M{"$regex": bson.RegEx{str, ""}},
	}).Skip(skip).Limit(limit).All(&communities)
	return communities
}

func (db Mongo) GetAll(ptr int, str string) ([]models.User, []models.Community) {
	skip := int(ptr) * limit
	var users []models.User
	mongoInstance.Instance.C("users").Find(bson.M{"$or": []bson.M{
		{"firstname": bson.M{"$regex": bson.RegEx{str, ""}}},
		{"lastname": bson.M{"$regex": bson.RegEx{str, ""}}},
	}}).Skip(skip).Limit(limit).All(&users)
	var communities []models.Community
	mongoInstance.Instance.C("communities").Find(bson.M{
		"name": bson.M{"$regex": bson.RegEx{str, ""}},
	}).Skip(skip).Limit(limit).All(&communities)
	return users, communities
}
