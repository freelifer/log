package main

import (
	"testing"
)

func Test_App_Get(t *testing.T) {

	var app App
	apps, err := app.Get()

	if err != nil {
		t.Error(err)
	} else {
		t.Log(apps)
	}
}
