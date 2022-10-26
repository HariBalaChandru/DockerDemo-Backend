package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ddosify/go-faker/faker"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// type profile struct {
// 	Name     string `json:"name"`
// 	Gender   string `json:"gender"`
// 	Mail     string `json:"mail"`
// 	Avatar   string `json:"avator"`
// 	Username string `json:"username"`
// 	Mobile   string `json:"mobile"`
// 	Color    string `json:"color"`
// 	DOB      string `json:"dob"`
// }

func createNewProfile() *bson.D {
	faker := faker.NewFaker()
	MaleOrFemale := func(val bool) string {
		if val {
			return "Male"
		} else {
			return "Female"
		}
	}
	profile := bson.D{
		{Key: "name", Value: faker.RandomPersonFullName()},
		{Key: "mail", Value: faker.RandomEmail()},
		{Key: "avatar", Value: faker.RandomAbstractImage()},
		{Key: "gender", Value: MaleOrFemale(faker.RandomBoolean())},
		{Key: "username", Value: faker.RandomUsername()},
		{Key: "mobile", Value: faker.RandomPhoneNumber()},
		{Key: "color", Value: faker.RandomSafeColorHex()},
		{Key: "dob", Value: faker.RandomDatePast()},
	}
	return &profile
}

func createNProfiles(N int) []interface{} {
	profiles := make([]interface{}, N)
	for i := 0; i < N; i++ {
		profiles[i] = createNewProfile()
	}
	return profiles
}

func clearDB(collection string) {
	customer := MCLI.Database("customer")
	profiles := customer.Collection(collection)
	err := profiles.Drop(MCTX)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB Cleared")
}

func populateDB(data []interface{}) {
	customer := MCLI.Database("customer")
	profiles := customer.Collection("profiles")
	profileResult, err := profiles.InsertMany(MCTX, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Inserted in DB : %v\n", profileResult)
}

func addProfiles(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var profiles []interface{}
	if err := json.Unmarshal(x, &profiles); err != nil {
		log.Fatal(err)
	}
	populateDB(profiles)
}

func readProfile() []bson.M {
	customer := MCLI.Database("customer")
	profiles := customer.Collection("profiles")
	cursor, err := profiles.Find(MCTX, bson.M{})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	var userprofiles []bson.M

	err = cursor.All(MCTX, &userprofiles)

	if err != nil {
		log.Fatal(err)
		return nil
	}
	return userprofiles
}

func closeDB() {
	defer MCLI.Disconnect(MCTX)
	defer MCAN()
}

func watcherReadProfile(c *gin.Context) {
	initDB()
	data := readProfile()
	closeDB()
	c.IndentedJSON(http.StatusOK, data)
}

func watcherCloseDB(c *gin.Context) {
	closeDB()
}

func watcherPopulateDB(c *gin.Context) {
	initDB()
	clearDB("profiles")
	populateDB(createNProfiles(10))
	closeDB()
}

func watcherClearDB(c *gin.Context) {
	initDB()
	clearDB("profiles")
	closeDB()
}

func watcherAddProfiles(c *gin.Context) {
	initDB()
	addProfiles(c)
	closeDB()
}

func watcherHealthCheck(c *gin.Context) {
	data := map[string]interface{}{
		"message": "success",
	}
	log.Println("Health Check...")
	c.IndentedJSON(http.StatusOK, data)
}
