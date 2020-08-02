package boltstore

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Human struct {
	Name   string
	Height float64
}

func testFile() *os.File {
	f, err := ioutil.TempFile(".", "BoltStore")
	if err != nil {
		panic(err)
	}
	return f
}

func TestGeneral(t *testing.T) {
	f := testFile()
	defer os.Remove(f.Name())
	ks, err := Open(f.Name())
	assert.Nil(t, err)
	err = ks.Set("hello", "world")
	assert.Nil(t, err)

	var a string
	err = ks.Get("hello", &a)
	assert.Nil(t, err)
	assert.Equal(t, "world", a)

	// Set a object, using a Gzipped JSON
	type Human struct {
		Name   string
		Height float64
	}
	err = ks.Set("human:1", Human{"Dante", 5.4})
	assert.Nil(t, err)

	var human Human
	err = ks.Get("human:1", &human)
	assert.Nil(t, err)
	assert.Equal(t, 5.4, human.Height)

	err = ks.Get("DOES NOT EXIST", &human)
	assert.NotNil(t, err)

	keys := ks.Keys()
	assert.Equal(t, []string{"hello", "human:1"}, keys)

	err = ks.Set("human:2", Human{"Dante2", 5.5})
	assert.Nil(t, err)
	err = ks.Set("human:3", Human{"Dante3", 5.6})
	assert.Nil(t, err)
	err = ks.GetAll(&human, func(key string) error {
		fmt.Println(key, human)
		return nil
	})
	assert.Nil(t, err)

	err = ks.Delete("human:1")
	assert.Nil(t, err)
	err = ks.Get("human:1", &human)
	assert.NotNil(t, err)

}

func BenchmarkGet(b *testing.B) {
	f := testFile()
	defer os.Remove(f.Name())
	ks, err := Open(f.Name())
	if err != nil {
		panic(err)
	}
	err = ks.Set("human:1", Human{"Dante", 5.4})
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var human Human
		ks.Get("human:1", &human)
	}
}

func BenchmarkSet(b *testing.B) {
	f := testFile()
	defer os.Remove(f.Name())
	ks, err := Open(f.Name())
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	// set a key to any object you want
	for i := 0; i < b.N; i++ {
		err := ks.Set("human:"+strconv.Itoa(i), Human{"Dante", 5.4})
		if err != nil {
			panic(err)
		}
	}
}
