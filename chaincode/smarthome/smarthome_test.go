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
	checkHome(t, stub, "101", "101")

}

func TestCreateHome(t *testing.T) {

	scc := new(SmartHome)
	stub := shim.NewMockStub("ex01",scc)

	checkInvoke(t, stub, [][]byte{[]byte("createHome"), []byte("301"), []byte("C"), []byte("1")})
	checkHome(t, stub, "301", "301")

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
  checkInvoke(t, stub, [][]byte{[]byte("initLedger")})
	checkInvoke(t, stub, [][]byte{[]byte("changeHomeOwnership"), []byte("104"), []byte("Test.Customer@example.com")})
	res := checkInvoke(t, stub, [][]byte{[]byte("queryHome"), []byte("104")})

	if res.Status != shim.OK {
		fmt.Println("Query Home failed", string(res.Message))
		t.FailNow()
	}

	homeAsBytes := res.Payload
	home := SmartHome{}

	json.Unmarshal(homeAsBytes, &home)

	if home.Customer != "Test.Customer@example.com" {
		fmt.Println("Incorrect customer ")
		t.FailNow()
	}
}

func TestNotifyFloorCompletion(t *testing.T) {

	scc := new(SmartHome)
	stub := shim.NewMockStub("ex01",scc)

	checkInvoke(t, stub, [][]byte{[]byte("initLedger")})
	checkInvoke(t, stub, [][]byte{[]byte("notifyFloorCompletion"), []byte("C"),[]byte("5")})
	towerAsBytes,_ := stub.GetState("C")
	tower := Tower{}

	json.Unmarshal(towerAsBytes, &tower)

	if tower.BuildStatus != "COM" {
		fmt.Println("Incorrect status ")
		t.FailNow()
	}
}

func TestVerifyFloorCompletion(t *testing.T) {

	scc := new(SmartHome)
	stub := shim.NewMockStub("ex01",scc)

	checkInvoke(t, stub, [][]byte{[]byte("initLedger")})
	checkInvoke(t, stub, [][]byte{[]byte("verifyFloorCompletion"), []byte("C"),[]byte("5"),[]byte("OK")})

	keyname := "tower~floor~bank"
	key,err := stub.CreateCompositeKey(keyname, []string{"C","5","bank1"})
	if err != nil {
		fmt.Println("Error forming composite key")
		t.FailNow()
	}

	endorsementAsBytes,_ := stub.GetState(key)

	fmt.Println(string(endorsementAsBytes))
	if string(endorsementAsBytes) != "OK" {
		fmt.Println("Couldn't verify status ")
		t.FailNow()
	}
}

func TestObtainCompletionVerification(t *testing.T) {
	scc := new(SmartHome)
	stub := shim.NewMockStub("ex01",scc)

	checkInvoke(t, stub, [][]byte{[]byte("initLedger")})
  checkInvoke(t, stub, [][]byte{[]byte("notifyFloorCompletion"), []byte("C"),[]byte("5")})
	checkInvoke(t, stub, [][]byte{[]byte("verifyFloorCompletion"), []byte("C"),[]byte("5"),[]byte("OK")})
	checkInvoke(t, stub, [][]byte{[]byte("obtainCompletionVerification"), []byte("C"),[]byte("5")})
	towerAsBytes,_ := stub.GetState("C")
	tower := Tower{}

	json.Unmarshal(towerAsBytes, &tower)

	if tower.BuildStatus != "VER" {
		fmt.Println("Invalid status expecting verified")
		t.FailNow()
	}
  checkInvoke(t, stub, [][]byte{[]byte("notifyFloorCompletion"), []byte("C"),[]byte("6")})
	checkInvoke(t, stub, [][]byte{[]byte("verifyFloorCompletion"), []byte("C"),[]byte("6"),[]byte("NOK")})
	res := checkInvoke(t, stub, [][]byte{[]byte("obtainCompletionVerification"), []byte("C"),[]byte("6")})

	if res.Status == shim.OK {
		fmt.Println("Invalid verify status")
		t.FailNow()
	}
}
