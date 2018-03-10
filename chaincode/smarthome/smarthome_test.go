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


import (
	"fmt"
	"testing"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	 sc "github.com/hyperledger/fabric/protos/peer"
)

type QueryResults struct {
		Key string
		Record SmartHome
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) sc.Response{

	res := stub.MockInvoke("1",args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}

	return res
}

func checkHome(t *testing.T, stub *shim.MockStub, id string, name string ) {

	res := stub.MockInvoke("1",[][]byte{[]byte("queryHome"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("QueryHome", id, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("QueryHome", id, "failed to get value")
		t.FailNow()
	}

  homeAsBytes := res.Payload
	home := SmartHome{}

	json.Unmarshal(homeAsBytes, &home)

	if home.Name != name {
		fmt.Println("QueryHome", id, "expecting name: ", name, " but found: ", home.Name)
		t.FailNow()
	}

}

func TestQueryHome(t *testing.T) {

	scc := new(SmartHome)
	stub := shim.NewMockStub("ex01",scc)

	checkInvoke(t, stub, [][]byte{[]byte("initLedger")})
	checkHome(t, stub, "HOME0", "101")

}

func TestCreateHome(t *testing.T) {

	scc := new(SmartHome)
	stub := shim.NewMockStub("ex01",scc)

	checkInvoke(t, stub, [][]byte{[]byte("createHome"), []byte("TESTHOME"), []byte("301"), []byte("3"), []byte("0"), []byte("customer.301@example.com")})
	checkHome(t, stub, "TESTHOME", "301")

}

func TestQueryAllHomes(t *testing.T) {

	scc := new(SmartHome)
	stub := shim.NewMockStub("ex01",scc)
  allhomesbytes :=make([]QueryResults,0)
  checkInvoke(t, stub, [][]byte{[]byte("initLedger")})
	all_homes:=checkInvoke(t, stub, [][]byte{[]byte("queryAllHomes")})
	fmt.Sprintf("Return payload is %s ", all_homes.Payload)
	json.Unmarshal(all_homes.Payload, &allhomesbytes)

  num_of_homes := len(allhomesbytes)
	fmt.Println("Number of records Found ", num_of_homes)
	if num_of_homes < 8 {
		fmt.Println("Expected more number of rows")
		t.FailNow()
	}

	i := 0
	//allhomes := make([]SmartHome,num_of_homes)
	for i < num_of_homes {
		fmt.Println("Key val ", allhomesbytes[i].Key)
		//json.Unmarshal(allhomesbytes[i].Record,&allhomes[i])
		fmt.Println("Found ", allhomesbytes[i].Record)
		i = i + 1
	}
}

func TestChangeHomeOwnership(t *testing.T) {

	scc := new(SmartHome)
	stub := shim.NewMockStub("ex01",scc)

	checkInvoke(t, stub, [][]byte{[]byte("changeHomeOwnership"), []byte("HOME0"), []byte("80"), []byte("20")})
	res := checkInvoke(t, stub, [][]byte{[]byte("queryHome"), []byte("HOME0")})

	homeAsBytes := res.Payload
	home := SmartHome{}

	json.Unmarshal(homeAsBytes, &home)

	if home.BuilderPerc != 80 {
		fmt.Println("Incorrect percentage ")
		t.FailNow()
	}
}
