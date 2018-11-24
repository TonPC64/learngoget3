package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func main() {
	e := echo.New()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	mongoHost := viper.GetString("mongo.host")
	mongoUser := viper.GetString("mongo.user")
	mongoPass := viper.GetString("mongo.pass")
	port := ":" + viper.GetString("port")

	e.Use(middleware.Logger())

	connString := fmt.Sprintf("%s:%s@%s", mongoUser, mongoPass, mongoHost)
	session, err := mgo.Dial(connString)
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
	h := handlers{session}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"status": "ok",
		})
	})
	e.GET("/todos", h.list)
	e.GET("/todos/:id", h.view)
	e.PUT("/todos/:id", h.done)
	e.POST("/todos", h.create)
	e.DELETE("/todos/:id", h.delete)
	e.Logger.Fatal(e.Start(port))
}

type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

type handlers struct {
	m *mgo.Session
}

func (h *handlers) create(c echo.Context) error {
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	session := h.m.Copy()
	defer session.Close()
	col := session.DB("workshop").C("todos")

	t.ID = bson.NewObjectId()
	if err := col.Insert(t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handlers) list(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	var ts []todo
	col := session.DB("workshop").C("todos")
	if err := col.Find(nil).All(&ts); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ts)
}

func (h *handlers) view(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	var t todo
	col := session.DB("workshop").C("todos")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handlers) done(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	var t todo
	col := session.DB("workshop").C("todos")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}
	t.Done = true
	if err := col.UpdateId(id, t); err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handlers) delete(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	col := session.DB("workshop").C("todos")
	if err := col.RemoveId(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"result": "success",
	})
}
