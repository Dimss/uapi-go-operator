package controller

import (
	"github.com/uapi-go-operator/pkg/controller/uapi"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, uapi.Add)
}
