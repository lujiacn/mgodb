package mgodb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/revel/revel"
)

// learn from leanote, put all db actions in one file
var (
	MgoSession *mgo.Session
	MgoDBName  string
	Dial       string
)

//Init setup mgo connection
func MgoDBInit() {
	//init mgoDB
	MgoDBConnect()

	//New bind type
	objID := bson.NewObjectId()
	revel.TypeBinders[reflect.TypeOf(objID)] = ObjectIDBinder
}

//MgoDBConnect do mgo connection
func MgoDBConnect() {
	var err error
	var found bool

	Dial = revel.Config.StringDefault("mongodb.dial", "localhost")
	if MgoDBName, found = revel.Config.String("mongodb.name"); !found {
		urls := strings.Split(Dial, "/")
		MgoDBName = urls[len(urls)-1]
	}

	MgoSession, err = mgo.Dial(Dial)

	if err != nil {
		panic("Cannot connect to database")
	}

	if MgoSession == nil {
		MgoSession, err = mgo.Dial(Dial)
		if err != nil {
			panic("Cannot connect to database")
		}
	}
}

func NewMgoSession() *mgo.Session {
	s := MgoSession.Clone()
	return s
}

//MgoControllerInit should be put in controller init function
func MgoControllerInit() {
	revel.InterceptMethod((*MgoController).Begin, revel.BEFORE)
	revel.InterceptMethod((*MgoController).End, revel.FINALLY)
}

//MgoController including the mgo session
type MgoController struct {
	*revel.Controller
	MgoSession *mgo.Session
}

//Begin do mgo connection
func (c *MgoController) Begin() revel.Result {
	if MgoSession == nil {
		MgoDBConnect()
	}

	c.MgoSession = MgoSession.Clone()
	return nil
}

//End close mgo session
func (c *MgoController) End() revel.Result {
	if c.MgoSession != nil {
		c.MgoSession.Close()
	}
	return nil
}

// ObjectIDBinder do binding
var ObjectIDBinder = revel.Binder{
	// Make a ObjectId from a request containing it in string format.
	Bind: revel.ValueBinder(func(val string, typ reflect.Type) reflect.Value {
		if len(val) == 0 {
			return reflect.Zero(typ)

		}
		if bson.IsObjectIdHex(val) {
			objID := bson.ObjectIdHex(val)
			return reflect.ValueOf(objID)

		}

		revel.AppLog.Error("ObjectIDBinder.Bind - invalid ObjectId!")
		return reflect.Zero(typ)

	}),
	// Turns ObjectId back to hexString for reverse routing
	Unbind: func(output map[string]string, name string, val interface{}) {
		var hexStr string
		hexStr = fmt.Sprintf("%s", val.(bson.ObjectId).Hex())
		// not sure if this is too carefull but i wouldn't want invalid ObjectIds in my App
		if bson.IsObjectIdHex(hexStr) {
			output[name] = hexStr

		} else {
			revel.AppLog.Error("ObjectIDBinder.Unbind - invalid ObjectId!")
			output[name] = ""

		}

	},
}
