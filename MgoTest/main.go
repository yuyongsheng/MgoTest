// MgoTest project main.go
package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}
type Men struct {
	Persons []Person
}

const (
	URL = "127.0.0.1:27017"
)

func catchPanic() {
	if v := recover(); v != nil {
		fmt.Println("recover(), catch error :", v)
	}
}

func Test() {
	defer catchPanic()

	dial := mgo.DialInfo{
		Addrs:    []string{URL},
		Username: "auth01",
		Password: "auth01",
		Database: "authTestDb01",
		Timeout:  10 * time.Second}

	session, err := mgo.DialWithInfo(&dial)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	db := session.DB("authTestDb01") //数据库名称
	collection := db.C("person")     //如果该集合已经存在的话，则直接返回

	// create user
	usr := mgo.User{
		Username:   "auth02",
		Password:   "auth02",
		Roles:      []mgo.Role{mgo.RoleReadWrite},
		CustomData: map[string]string{"addr": "tianjin", "gender": "male", "mail": "1234@123.com"},
	}
	err = db.UpsertUser(&usr)
	if err != nil {
		panic(err)
	}

	//*****集合中元素数目********
	countNum, err := collection.Count()
	if err != nil {
		panic(err)
	}
	fmt.Println("Things objects count: ", countNum)

	//*******插入元素*******
	temp := &Person{
		Phone: "18811577546",
		Name:  "zhangzheHero"}

	//一次可以插入多个对象 插入两个Person对象
	err = collection.Insert(&Person{"Ale", "+55 53 8116 9639"}, temp)
	if err != nil {
		panic(err)
	}

	//*****查询单条数据*******
	result := Person{}
	err = collection.Find(bson.M{"phone": "456"}).One(&result)
	fmt.Println("Phone:", result.Name, result.Phone)

	//*****查询多条数据*******
	var personAll Men //存放结果
	iter := collection.Find(nil).Iter()
	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.Name)
		personAll.Persons = append(personAll.Persons, result)
	}

	//*******更新数据**********
	err = collection.Update(bson.M{"name": "ccc"}, bson.M{"$set": bson.M{"name": "ddd"}})
	err = collection.Update(bson.M{"name": "ddd"}, bson.M{"$set": bson.M{"phone": "12345678"}})
	err = collection.Update(bson.M{"name": "aaa"}, bson.M{"phone": "1245", "name": "bbb"})

	//******删除数据************
	_, err = collection.RemoveAll(bson.M{"name": "Ale"})
}

func main() {
	fmt.Println("Hello World!")

	Test()

	fmt.Println("END")
}
