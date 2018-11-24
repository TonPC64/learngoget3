package main

import (
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	session, err := mgo.Dial("localhost")
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
	h := handlers{session}

	e.GET("/todos", h.list)
	e.GET("/todos/:id", h.view)
	e.POST("/todos", h.create)
	e.Logger.Fatal(e.Start(":1323"))
}

type handlers struct {
	m *mgo.Session
}

type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
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
