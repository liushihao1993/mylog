package mylog

import (
	"context"
	"testing"
)

type UserInfo struct {
	Name string
}

func TestCtx(t *testing.T) {
	Println("Hello World")
	Ctx(context.Background()).WithField("heelo", "world").
		WithFields("nihao", map[string]string{"k1": "v1"},
			"k2", 222, "UserInfo", UserInfo{Name: "dsf"}, 333, "addf", "as").Info("Hello World")
	Ctx(context.Background()).WithField("heelo", "world").
		WithFields("nihao", map[string]string{"k1": "v1"},
			"k2", 222, "UserInfo", UserInfo{Name: "dsf"}, "addf", "as").Infof("Hello World")
}

// TestPrintf
func TestPrintf(t *testing.T) {
	Printf("Hello World")
	Printf("Hello %s", "World")
	Printf("Hello %s %d", "World", 123)
	Printf("Hello %s %d %v", "World", 123, map[string]string{"k1": "v1"})
	Printf("Hello %s %d %v %v", "World", 123, map[string]string{"k1": "v1"}, UserInfo{Name: "dsf"})
}

// TestPrintln
func TestPrintln(t *testing.T) {
	Println("Hello World")
	Println("Hello", "World")
	Println("Hello", "World", 123)
	Println("Hello", "World", 123, map[string]string{"k1": "v1"})
	Println("Hello", "World", 123, map[string]string{"k1": "v1"}, UserInfo{Name: "dsf"})
}

// TestError
func TestError(t *testing.T) {
	Ctx(context.Background()).Error("Hello World")
	Ctx(context.Background()).Error("Hello", "World")
	Ctx(context.Background()).Error("Hello", "World", 123)
	Ctx(context.Background()).Error("Hello", "World", 123, map[string]string{"k1": "v1"})
	Ctx(context.Background()).Error("Hello", "World", 123, map[string]string{"k1": "v1"}, UserInfo{Name: "dsf"})

}

// TestErrorf
func TestErrorf(t *testing.T) {
	Ctx(context.Background()).Errorf("Hello World")
	Ctx(context.Background()).Errorf("Hello %s", "World")
	Ctx(context.Background()).Errorf("Hello %s %d", "World", 123)
	Ctx(context.Background()).Errorf("Hello %s %d %v", "World", 123, map[string]string{"k1": "v1"})
	Ctx(context.Background()).Errorf("Hello %s %d %v %v", "World", 123, map[string]string{"k1": "v1"}, UserInfo{Name: "dsf"})
}
