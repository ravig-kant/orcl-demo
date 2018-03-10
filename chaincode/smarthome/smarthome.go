/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the SmartHome structure, with 4 properties.  Structure tags are used by encoding/json library
type SmartHome struct {
		Name   string `json:"name"`
    Tower  int `json:"tower"`
    Floor int `json:"floor"`
    CompletedFloor  int `json:"completedFloor"`
		BuildStatus	int	`json:"buildStatus"`
		BuilderPerc	int	`json:"builderPerc"`
		CustomerPerc	int	`json:"customerPerc"`
		Customer	string	`json:"customer"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartHome) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartHome) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryHome" {
		return s.queryHome(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createHome" {
		return s.createHome(APIstub, args)
	} else if function == "queryAllHomes" {
		return s.queryAllHomes(APIstub)
	} else if function == "changeHomeOwnership" {
		return s.changeHomeOwnership(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartHome) queryHome(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	homeAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(homeAsBytes)
}

func (s *SmartHome) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	homes := []SmartHome{
		SmartHome{Name: "101", Tower: 1, Floor: 0, CompletedFloor: 0, BuildStatus: 0, BuilderPerc: 100, CustomerPerc: 0, Customer: "customer.101@example.com"},
		SmartHome{Name: "102", Tower: 1, Floor: 0, CompletedFloor: 0, BuildStatus: 0, BuilderPerc: 100, CustomerPerc: 0, Customer: "customer.102@example.com"},
		SmartHome{Name: "103", Tower: 1, Floor: 0, CompletedFloor: 0, BuildStatus: 0, BuilderPerc: 100, CustomerPerc: 0, Customer: "customer.103@example.com"},
		SmartHome{Name: "104", Tower: 1, Floor: 0, CompletedFloor: 0, BuildStatus: 0, BuilderPerc: 100, CustomerPerc: 0, Customer: "customer.104@example.com"},
		SmartHome{Name: "201", Tower: 2, Floor: 0, CompletedFloor: 1, BuildStatus: 1, BuilderPerc: 80, CustomerPerc: 20, Customer: "customer.201@example.com"},
		SmartHome{Name: "202", Tower: 2, Floor: 0, CompletedFloor: 1, BuildStatus: 1, BuilderPerc: 80, CustomerPerc: 20, Customer: "customer.202@example.com"},
		SmartHome{Name: "203", Tower: 2, Floor: 0, CompletedFloor: 1, BuildStatus: 1, BuilderPerc: 80, CustomerPerc: 20, Customer: "customer.203@example.com"},
		SmartHome{Name: "204", Tower: 2, Floor: 0, CompletedFloor: 1, BuildStatus: 1, BuilderPerc: 80, CustomerPerc: 20, Customer: "customer.204@example.com"},
	}

	i := 0
	for i < len(homes) {
		fmt.Println("i is ", i)
		homeAsBytes, err := json.Marshal(homes[i])
		if err != nil {
			fmt.Println("error while converting ", err.Error())
			return shim.Error(err.Error())
		}
		APIstub.PutState("HOME"+strconv.Itoa(i), homeAsBytes)
		fmt.Println("Added", homes[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartHome) createHome(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

  iTower, _ := strconv.Atoi(args[2])
	iFloor, _ := strconv.Atoi(args[3])
	var home = SmartHome{Name: args[1], Tower: iTower, Floor: iFloor, CompletedFloor: 0, BuildStatus: 0, BuilderPerc: 100, CustomerPerc: 0, Customer: args[4]}

	homeAsBytes, _ := json.Marshal(home)
	APIstub.PutState(args[0], homeAsBytes)

	return shim.Success(nil)
}

func (s *SmartHome) queryAllHomes(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "HOME0"
	endKey := "HOME999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllHomes:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartHome) changeHomeOwnership(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	homeAsBytes, _ := APIstub.GetState(args[0])
	home := SmartHome{}

	json.Unmarshal(homeAsBytes, &home)
	home.BuilderPerc, _ = strconv.Atoi(args[1])
	home.CustomerPerc, _ = strconv.Atoi(args[2])

	homeAsBytes, _ = json.Marshal(home)
	APIstub.PutState(args[0], homeAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartHome))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
