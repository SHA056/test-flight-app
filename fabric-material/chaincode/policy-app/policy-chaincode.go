
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

)


type SmartContract struct{

}

type Userdetails struct {
	Userid string `json:"userid"`
	Name string `json:"name"`
	Email string `json:"email"`
	Age string `json:"age"`
}

type Policydetails struct {
	Polid string `json:"policyid"`
	Polname string `json:"polname"`
	Polvalue string `json:"polvalue"`
	Polvalidity string `json:"polvalidity"`
}

type Flightdetails struct {
	Flightnum string `json:"flightnum"`
	Airlinename string `json:"airlinename"`
	Arrdep string `json:"arrdep"`
	Iatacode string `json:"iatacode"` 
	Time string `json:"time"`
}

// Create schema for block entrY
type Record struct {
	Userdet Userdetails `json:"userdet"`
	Poldet Policydetails `json:"policydet"`
	Flightdet Flightdetails `json:"flightdet"`
	Claimstatus string `json:"claimstatus"`
}



func (c *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (c *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, arg := stub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "newUser" {
		return c.newUser(stub, arg)
	} else if function == "selectPolicy" {
		return c.selectPolicy(stub, arg) 
	} else if function == "flightDet" {
		return c.flightDet(stub, arg) 
	} else if function == "rollClaim" {
		return c.rollClaim(stub)
	} else if function == "viewLedger" {
		return c.viewLedger(stub) 
	} else if function == "initLedger" {
		return c.initLedger(stub)
	}

//  Add functions to query by features (Policy name/number, flight, customer name/ID, claim status)
	

	return shim.Error("Invalid function name. Call \"newUser\", \"selectPolicy\", \"flightDet\", \"rollClaim\", \"viewLedger\", \"initLedger\".")
}



func (c *SmartContract) newUser(stub shim.ChaincodeStubInterface, arg []string) pb.Response {
	if len(arg) != 5 {
		return shim.Error("Expecting 5 values: Key, UserID, Name, Email, Age")
	}

	var record = Record{Userdet: Userdetails{Userid:arg[1], Name: arg[2], Email: arg[3], Age: arg[4]}, Claimstatus: "UNCLAIMED"}

	recordAsBytes, _ := json.Marshal(record)
	err := stub.PutState(arg[0], recordAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create a new record: %s", arg[0]))
	}

	return shim.Success(nil)
}



func (c *SmartContract) selectPolicy(stub shim.ChaincodeStubInterface, arg []string) pb.Response {
	if len(arg) != 5 {
		return shim.Error("Expecting 5 values: Key, ID, Name, Value and Validity of the policy selected.")
	}

	recordAsBytes,_ := stub.GetState(arg[0])
	if recordAsBytes == nil {
		return shim.Error("Could not locate user.")
	}

	var record = Record{}

	json.Unmarshal(recordAsBytes,&record)

	record.Poldet.Polid = arg[1]
	record.Poldet.Polname = arg[2]
	record.Poldet.Polvalue = arg[3]
	record.Poldet.Polvalidity = arg[4]

	recordAsBytes, _ = json.Marshal(record)
	err := stub.PutState(arg[0],recordAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error during Policy Selection: %s", arg[0]))
	}

	return shim.Success(nil)
}



func (c *SmartContract) flightDet(stub shim.ChaincodeStubInterface, arg []string) pb.Response {
	if len(arg) != 6 {
		return shim.Error("Expecting 6 values: Key, Flight Number, Airline Name, arrival/departure, Airport iata_code, Time of")
	}

	recordAsBytes,_ := stub.GetState(arg[0])
	if recordAsBytes == nil {
		return shim.Error("Could not locate user.")
	}

	var record = Record{}

	json.Unmarshal(recordAsBytes,&record)

	record.Flightdet.Flightnum = arg[1]
	record.Flightdet.Airlinename = arg[2]
	record.Flightdet.Arrdep = arg[3]
	record.Flightdet.Iatacode = arg[4]
	record.Flightdet.Time = arg[5]

	recordAsBytes, _ = json.Marshal(record)
	err := stub.PutState(arg[0],recordAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error during Flight Entry: %s", arg[0]))
	}

	return shim.Success(nil)
}



func (c *SmartContract) rollClaim(stub shim.ChaincodeStubInterface) pb.Response {
	startKey := "0"
	endKey := "999"

//	The block is traversed and filtered to get a single block variable called record.
//	Record gets overwritten on every iteration.	
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var record = Record{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(queryResponse.Value,&record)

/*		The variables used to compare against the API response:

		record.Claimstatus 
		record.Flightdet.Arrdep
		record.Flightdet.Flightnum
		record.Flightdet.Iatacode
		record.Flightdet.Time
*/
		if record.Claimstatus == "UNCLAIMED" || record.Claimstatus == "ACTIVE" {
			
			// call API here
			APIstring := "http://aviation-edge.com/api/public/timetable?key=73a96e-e3e7e8-3b6580-f37bdf-94c179&iataCode="+record.Flightdet.Iatacode+"&type="+record.Flightdet.Arrdep

			Resp,_ := http.Get(APIstring)

			respBody,_ := ioutil.ReadAll(Resp.Body)

	    	var data []map[string]interface{}
    		err := json.Unmarshal([]byte(respBody), &data)
    		if err != nil {
        		panic(err)
    		}

 		   	for i := 0; i<len(data); i++ { // The maximum number iterations have to be adjusted to true values
    			status := data[i]["status"]
    			flightnumber := data[i]["flight"].(map[string]interface{})["iataNumber"]
    			timestamp := data[i]["departure"].(map[string]interface{})["scheduledTime"]

    			if flightnumber == record.Flightdet.Flightnum && timestamp == record.Flightdet.Time {
    				if status == "cancelled" {
    					record.Claimstatus = "Claim ACCEPTED: flight cancelled"
    				} else if status == "landed" {
    					record.Claimstatus = "NOT VALID FOR CLAIM: flight on-time"
    				} else if status == "scheduled" || status == "active" {
    					record.Claimstatus = "ACTIVE"
    				}

    			}
    		}

    	}

		recordAsBytes, _ := json.Marshal(record)
		err = stub.PutState(queryResponse.Key,recordAsBytes)
		if err != nil {
			return shim.Error("Error while saving Claimstatus")
		}
	}

	return shim.Success(nil)
}



func (c *SmartContract) viewLedger(stub shim.ChaincodeStubInterface) pb.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(",\"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- viewLedger:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}



func (c *SmartContract) initLedger(stub shim.ChaincodeStubInterface) pb.Response {
	record := []Record{
		Record{Userdet: Userdetails{Userid:"1", Name: "Bruce Banner", Email: "hulk@q.com", Age: "43"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"111", Polname: "ShortFlightCancellation", Polvalue: "3000", Polvalidity: "23085434218"}, Flightdet: Flightdetails{Flightnum: "EY8868", Airlinename: "Etihad Airways", Arrdep: "departure", Iatacode: "BLR", Time: "2018-10-10T08:00:00.000"}},
		Record{Userdet: Userdetails{Userid:"2", Name: "Stephen Strange", Email: "strange@q.com", Age: "49"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"121", Polname: "LongFlightCancellation", Polvalue: "15000", Polvalidity: "43862643284"}, Flightdet: Flightdetails{Flightnum: "DL8708", Airlinename: "Delta Air Lines", Arrdep: "departure", Iatacode: "CDG", Time: "2018-10-10T09:30:00.000"}},
		Record{Userdet: Userdetails{Userid:"3", Name: "Natalia Romanova", Email: "widow@q.com", Age: "35"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"111", Polname: "ShortFlightCancellation", Polvalue: "4000", Polvalidity: "23085434218"}, Flightdet: Flightdetails{Flightnum: "V2K2000", Airlinename: "Anisec", Arrdep: "arrival", Iatacode: "CDG", Time: "2018-10-10T07:00:00.000"}},
		Record{Userdet: Userdetails{Userid:"4", Name: "Wade Wilson", Email: "deadpool@q.com", Age: "28"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"121", Polname: "LongFlightCancellation", Polvalue: "20000", Polvalidity: "783262347857."}, Flightdet: Flightdetails{Flightnum: "9W838", Airlinename: "Jet Airways (India)", Arrdep: "arrival", Iatacode: "DEL", Time: "2018-10-10T13:55:00.000"}},
		Record{Userdet: Userdetails{Userid:"2", Name: "Stephen Strange", Email: "strange@q.com", Age: "49"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"121", Polname: "ShortFlightCancellation", Polvalue: "5000", Polvalidity: "693750379453"}, Flightdet: Flightdetails{Flightnum: "BA5867", Airlinename: "British Airways", Arrdep: "arrival", Iatacode: "COK", Time: "2018-10-10T06:00:00.000"}},
		Record{Userdet: Userdetails{Userid:"3", Name: "Natalia Romanova", Email: "widow@q.com", Age: "35"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"121", Polname: "LongFlightCancellation", Polvalue: "18000", Polvalidity: "432982582870"}, Flightdet: Flightdetails{Flightnum: "AT9526", Airlinename: "Royal Air Maroc", Arrdep: "departure", Iatacode: "JFK", Time: "2018-10-10T22:00:00.000"}},
		Record{Userdet: Userdetails{Userid:"1", Name: "Bruce Banner", Email: "hulk@q.com", Age: "43"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"111", Polname: "ShortFlightCancellation", Polvalue: "4000", Polvalidity: "985943755432"}, Flightdet: Flightdetails{Flightnum: "7C1104", Airlinename: "Jeju Air", Arrdep: "departure", Iatacode: "ICN", Time: "2018-10-10T15:05:00.000"}},
		Record{Userdet: Userdetails{Userid:"5", Name: "Carol Danvers", Email: "marvel@q.com", Age: "37"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"121", Polname: "LongFlightCancellation", Polvalue: "50000", Polvalidity: "67327537803"}, Flightdet: Flightdetails{Flightnum: "R5501", Airlinename: "Jordan Aviation", Arrdep: "arrival", Iatacode: "IEV", Time: "2018-10-10T10:00:00.000"}},
		Record{Userdet: Userdetails{Userid:"3", Name: "Natalia Romanova", Email: "widow@q.com", Age: "35"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"111", Polname: "ShortFlightCancellation", Polvalue: "3000", Polvalidity: "5436246932847"}, Flightdet: Flightdetails{Flightnum: "BA4558", Airlinename: "British Airways", Arrdep: "departure", Iatacode: "SGN", Time: "2018-10-10T11:15:00.000"}},
		Record{Userdet: Userdetails{Userid:"4", Name: "Wade Wilson", Email: "deadpool@q.com", Age: "28"}, Claimstatus:"UNCLAIMED", Poldet: Policydetails{Polid:"111", Polname: "ShortFlightCancellation", Polvalue: "6000", Polvalidity: "783538493343"}, Flightdet: Flightdetails{Flightnum: "HA397", Airlinename: "Hawaiian Airlines", Arrdep: "departure", Iatacode: "KOA", Time: "2018-10-10T21:45:00.000"}},
	}

	i := 000
	for i < len(record) {
		fmt.Println("i is ", i)
		recordAsBytes, _ := json.Marshal(record[i])
		stub.PutState(strconv.Itoa(i+1), recordAsBytes)
		fmt.Println("Added", record[i])
		i = i + 1
	}

	return shim.Success(nil)
}



func main() {
	// Creates a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating a new Smart Contract: %s", err)
	}
}
