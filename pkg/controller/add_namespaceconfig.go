package controller

import (
	"github.com/raffaelespazzoli/namespace-configuration-controller/pkg/controller/namespaceconfig"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, namespaceconfig.Add)
}
