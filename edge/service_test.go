package edge

import (
	"reflect"
	"testing"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/service/mucp"
)

type MyStruct struct {
	Name string
	Age  int32
}

func TestInterfaceEqual(t *testing.T) {

	var ms = micro.NewService() //this is a micro.Service
	var ss = mucp.NewService()  //this is a service.Service

	if reflect.TypeOf(ms).Kind() != reflect.TypeOf(ss).Kind() {
		t.Errorf("different type %v : %v", ms, ss)
	}
	//so they are equl type

}

var b = map[string]MyStruct{}

func TestRessignVars(t *testing.T) {

	var a = map[string]interface{}{}

	var c = map[string]MyStruct{
		"c": MyStruct{
			Name: "im c",
		},
	}
	//so they are equl type
	if reflect.TypeOf(b).Kind() == reflect.TypeOf(a).Kind() {
		t.Logf("equal type %v : %v", b, a)
	}

	b = c

	if reflect.TypeOf(b).Kind() != reflect.TypeOf(c).Kind() {
		t.Errorf("different type %v : %v", b, c)
	}

	t.Log("b.name:", b["c"].Name)
}
