package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var mg MongoInstance

const DBName = "fire-hrms"
const mongoURI = "mongodb://localhost:27017/" + DBName

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id, omitempty"` //bson because we're working woth mongodb
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

func Connect() error {
	client, err1 := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err1 != nil {
		return err1
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err2 := client.Connect(ctx)
	db := client.Database(DBName)
	if err2 != nil {
		return err2
	}

	mg = MongoInstance{
		Client: client,
		DB:     db,
	}
	fmt.Println("DB Connection Success")
	return nil
}

func GetEmployees(c *fiber.Ctx) {
	query := bson.D{{}}
	cursor, err := mg.DB.Collection("employees").Find(c.Context(), query)
	if err != nil {
		c.Status(500).SendString(err.Error())
	}

	var employees []Employee = make([]Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		c.Status(500).SendString(err.Error())
	}

	c.JSON(employees)
}

func NewEmployee(c *fiber.Ctx) {
	collection := mg.DB.Collection("employees")

	employee := new(Employee)

	if err := c.BodyParser(employee); err != nil {
		c.Status(400).SendString(err.Error())
	}

	insertResult, err := collection.InsertOne(c.Context(), employee)

	if err != nil {
		c.Status(500).SendString(err.Error())
	}

	//get employee data with employee ID just created
	emp_query := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	NewEmployeeRecord := collection.FindOne(c.Context(), emp_query)

	createdEmployee := &Employee{}
	NewEmployeeRecord.Decode(createdEmployee)

	c.Status(201).JSON(createdEmployee)

}

func UpdateEmployee(c *fiber.Ctx) {
	idParam := c.Params("id")

	empID, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		c.Status(400)
	}

	employee := new(Employee)

	if err := c.BodyParser(employee); err != nil {
		c.Status(400).Send(err)
	}

	query := bson.D{{Key: "_id", Value: empID}}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			},
		},
	}

	err = mg.DB.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.SendStatus(400)
		}
		c.SendStatus(500)
	}

	employee.ID = idParam

	c.Status(200).JSON(employee)
}

func deleteEmployee(c *fiber.Ctx) {
	idParam := c.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		c.SendStatus(500)
	}

	query := bson.D{{Key: "_id", Value: employeeID}}

	result, err := mg.DB.Collection("employees").DeleteOne(c.Context(), query)

	if err != nil {
		c.SendStatus(500)
	}

	if result.DeletedCount < 1 {
		c.SendStatus(404)
	}

	c.Status(200).JSON("Record Deleted")
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employee", GetEmployees)
	app.Post("/employee", NewEmployee)
	app.Put("/employee/:id", UpdateEmployee)
	app.Delete("/employee/:id", deleteEmployee)

	log.Fatal(app.Listen(":3000"))
}
