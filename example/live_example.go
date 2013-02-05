// Copyright 2013 Canonical Ltd.  This software is licensed under the
// GNU Lesser General Public License version 3 (see the file COPYING).

/*
This is an example on how the Go library gomaasapi can be used to interact with
a real MAAS server.
Note that this is a provided only as an example and that real code should probably do something more sensible with errors than ignoring them or panicking.
*/
package main

import (
	"fmt"
	"launchpad.net/gomaasapi"
	"net/url"
)

var apiKey string
var apiURL string

func init() {
	fmt.Println("Warning: this will create a node on the MAAS server; it should be deleted at the end of the run but if something goes wrong, that test node might be left over.  You've been warned.")
	fmt.Print("Enter API key: ")
	fmt.Scanf("%s", &apiKey)
	fmt.Print("Enter API URL: ")
	fmt.Scanf("%s", &apiURL)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	authClient, err := gomaasapi.NewAuthenticatedClient(apiURL, apiKey)
	checkError(err)

	maas := gomaasapi.NewMAAS(*authClient)

	nodeListing := maas.GetSubObject("nodes")

	// List nodes.
	fmt.Println("Fetching list of nodes...")
	listNodeObjects, err := nodeListing.CallGet("list", url.Values{})
	checkError(err)
	listNodes, err := listNodeObjects.GetArray()
	checkError(err)
	fmt.Printf("Got list of %v nodes\n", len(listNodes))
	for index, nodeObj := range listNodes {
		node, err := nodeObj.GetMAASObject()
		checkError(err)
		hostname, err := node.GetField("hostname")
		checkError(err)
		fmt.Printf("Node #%d is named '%v' (%v)\n", index, hostname, node.URL())
	}

	// Create a node.
	fmt.Println("Creating a new node...")
	params := url.Values{"architecture": {"i386/generic"}, "mac_addresses": {"AA:BB:CC:DD:EE:FF"}}
	newNodeObj, err := nodeListing.CallPost("new", params)
	checkError(err)
	newNode, err := newNodeObj.GetMAASObject()
	checkError(err)
	newNodeName, err := newNode.GetField("hostname")
	checkError(err)
	fmt.Printf("New node created: %s (%s)\n", newNodeName, newNode.URL())

	// Update the new node.
	fmt.Println("Updating the new node...")
	updateParams := url.Values{"hostname": {"mynewname"}}
	newNodeObj2, err := newNode.Update(updateParams)
	checkError(err)
	newNode2, err := newNodeObj2.GetMAASObject()
	checkError(err)
	newNodeName2, err := newNode2.GetField("hostname")
	checkError(err)
	fmt.Printf("New node updated, now named: %s\n", newNodeName2)

	// Count the nodes.
	listNodeObjects2, err := nodeListing.CallGet("list", url.Values{})
	checkError(err)
	listNodes2, err := listNodeObjects2.GetArray()
	checkError(err)
	fmt.Printf("We've got %v nodes\n", len(listNodes2))

	// Delete the new node.
	fmt.Println("Deleting the new node...")
	errDelete := newNode.Delete()
	checkError(errDelete)

	// Count the nodes.
	listNodeObjects3, err := nodeListing.CallGet("list", url.Values{})
	checkError(err)
	listNodes3, err := listNodeObjects3.GetArray()
	checkError(err)
	fmt.Printf("We've got %v nodes\n", len(listNodes3))

	fmt.Println("All done.")
}
