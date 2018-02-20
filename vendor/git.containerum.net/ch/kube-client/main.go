package main

import (
	"fmt"

	"git.containerum.net/ch/kube-client/pkg/cmd"
)

var (
	emptyQuery = make(map[string]string)
)

//ONLY for FIRST TESTS
func main() {
	client, err := cmd.CreateCmdClient(cmd.User{
		Role: "admin",
	})
	if err != nil {
		panic(err)
	}

	nsList, err := client.GetNamespaceList(emptyQuery)
	if err != nil {
		panic(err)
	}
	fmt.Println(nsList)

	ns, err := client.GetNamespace(nsList[0].Name)
	if err != nil {
		panic(err)
	}
	fmt.Println(ns)
}
