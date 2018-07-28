 package main

 import (
 "fmt"
 "github.com/hyperledger/fabric/core/chaincode/shim"
 pb "github.com/hyperledger/fabric/protos/peer"
 "encoding/json"
 )

 type PcXchg struct {
  }

 type PC struct {
     Snumber string
     Serie   string
     Others  string
     Status  string
 }

 func (c *PcXchg) Init(stub shim.ChaincodeStubInterface) pb.Response {
     return shim.Success(nil)
 }

 func (c *PcXchg) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
     
     function, args := stub.GetFunctionAndParameters()

     switch function {
     case "createPC":
         return c.createPC(stub, args)
     case "buyPC":
         return c.updateStatus(stub,args, "bought")
     case "handBackPC":
         return c.updateStatus(stub,args,"returned")
     case "queryStock":
         return c.queryStock(stub, args)
     case "queryDetail":
         return c.queryDetail(stub,args)
     default:
         return shim.Error("Available functions: createPC, buyPC, handBackPC, queryStock, queryDetail")
    }

  
 }

 func (c *PcXchg) createPC(stub shim.ChaincodeStubInterface, args []string) pb.Response {
     if len(args) != 3 {
        return shim.Error("createPC arguments usage: Serialnumber, Serie, Others")
     }

     pc := PC{args[0],args[1],args[2],"available"}
     
     pcAsBytes, err := json.Marshal(pc)

     if err != nil {
         return shim.Error(err.Error())
     } 
     
     err = stub.PutState(pc.Snumber, pcAsBytes)
     
     if err != nil {
         return shim.Error(err.Error())
     }

     return shim.Success(nil)
 }

 func (c *PcXchg) updateStatus(stub shim.ChaincodeStubInterface, args []string, status string) pb.Response {
     if len(args) != 1 {
        return shim.Error("This function needs the serial number as argument")
     }
     v,err := stub.GetState(args[0])
     if err !=nil {
      return shim.Error("Serial number " + args[0] + "not found")
     }

     var pc PC
     json.Unmarshal(v, &pc)
     pc.Status = status
     pcAsBytes, err := json.Marshal(pc)
     err = stub.PutState(pc.Snumber, pcAsBytes)
     if err != nil{
       return shim.Error(err.Error())
     }
 
     return shim.Success(nil)
 }

 // queryDetail gives all fields of stored data and wants to have the serial number
 func (c *PcXchg) queryDetail(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // Look for the serial number
    value, err := stub.GetState(args[0])
    if err != nil {
        return shim.Error("Serial number " + args[0] + " not found")
    }

    var pc PC
    // Decode value
    json.Unmarshal(value, &pc)

    fmt.Print(pc)
    // Response info
    return shim.Success([]byte(" SNMBR: " + pc.Snumber + " Serie: " + pc.Serie + " Others: " + pc.Others + " Status: " + pc.Status)) 
 } 



 // queryStock give all stored keys in the database
 func (c *PcXchg) queryStock(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // See stub.GetStateByRange in interfaces.go
    start, end := "",""

    if len(args) == 2 {
        start, end = args[0], args[1]
    } 

    // resultIterator is a StateQueryIteratorInterface
    resultsIterator, err := stub.GetStateByRange(start, end)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    keys := " \n"
    // This interface includes HasNext,Close and Next
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        keys+=queryResponse.Key + " \n"
    }

    fmt.Println(keys)

    return shim.Success([]byte(keys))
 }

 func main() {
    err := shim.Start(new(PcXchg))
    if err != nil {
        fmt.Printf("Error starting chaincode sample: %s", err)
    }
 }
