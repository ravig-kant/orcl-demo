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

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInvoke(t *testing, stub *shim.MockStub, args [][]byte) {

	res := stub.MockInvoke("1",args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkHome(t *testing, stub *shim.MockStub, string id, string name) {

	res := stub.MockInvoke("1",[][]byte{[]byte("queryHome"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("QueryHome", id, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("QueryHome", id, "failed to get value")
		t.FailNow()
	}
  homeAsBytes, _ := res.Payload
	home := SmartHome{}

	json.Unmarshal(homeAsBytes, &home)
	if home.name != name {
		fmt.Println("QueryHome", id, "expecting name: ", name, " but found: ", home.name)
		t.FailNow()
	}

}

func testQueryHome(t *testing.T) {

	scc := new(SmartHome)
	stub = shim.NewMockStub("ex01",scc)

	checkInvoke(t, stub, [][]byte{[]byte("initLedger")})
	checkHome(t, stub, "HOME1", "101")

}
