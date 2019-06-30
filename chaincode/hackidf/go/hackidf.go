package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"crypto/sha256"
)
var (
	fileName = "hackidf"
)

//structure of chaincode
type HackidfChaincode struct{
}

// User Structure
type User struct{
	Username string `json:"Username"`
	PasswordHash [32]byte `json:"PasswordHash"`
	Email string `json:"Email"`
	Ph string `json:"Ph"`
	IsVerified string `json:"IsVerified"`
}

// Organisation Structure
type Organisation struct{
	OrgName string `json:"OrgName"`
	PasswordHash [32]byte `json:"PasswordHash"`
	IsVerified string `json:"IsVerified"`
}

// Claim 
type Claim struct{
	UserID [32]byte `json:"UserID"`
	OrgID [32]byte `json:"OrgID"`
	Skill [32]byte `json:"Skill"`
	Comments string `json:"Comments"`
	Timestamp [32]byte `json:"Timestamp"`
	IsVerified string `json:"IsVerified"`
}

// Issuing a certificate
type Cert struct{
	OrgID string `json:"OrgID"`
	CertHash [32]byte `json:"CertHash"`
}

//initialization function
func (t *HackidfChaincode) Init(stub shim.ChaincodeStubInterface)pb.Response{
	// Whatever variable initialisation you want can be done here //
	return shim.Success(nil)
}

// invoking functions
func  (t *HackidfChaincode) Invoke(stub shim.ChaincodeStubInterface)pb.Response{
	// IF-ELSE-IF all the functions 
	function, args := stub.GetFunctionAndParameters()
	if function == "CreateUser" {
		return t.CreateUser(stub, args)
	}else if function == "VerifyUser" {
		return t.VerifyUser(stub, args)
	}else if function == "CreateOrg" {
		return t.CreateOrg(stub, args)
	}else if function == "VerifyOrg" {
		return t.VerifyOrg(stub, args)
	}else if function == "MakeClaim" {
		return t.MakeClaim(stub, args)
	}else if function == "VerifyClaim"{
		return t.VerifyClaim(stub, args)
	}else if function == "IssueCert"{
		return t.IssueCert(stub, args)
	}else if function == "Query"{
		return t.Query(stub, args)
	}else if function == "QueryWithCert"{
		return t.QueryWithCert(stub, args)
	}
	fmt.Println("invoke did not find func : " + function) //error
	return shim.Error("Received unknown function invocation")
	// end of all functions
}

// Check KYC
func CheckUser(stub shim.ChaincodeStubInterface, UserID string) int {
	UserAsBytes, err := stub.GetState(UserID)
	var User User
	err = json.Unmarshal(UserAsBytes, &User)
	if err != nil {
		return 0
	}
	if User.IsVerified == "True" {
        return 1
	}
	return 0
}

// Check If Org exists and is verified
func CheckOrg(stub shim.ChaincodeStubInterface, OrgID string) int {
	OrgAsBytes, err := stub.GetState(OrgID)
	if err != nil {
		return 0
	}else if OrgAsBytes == nil{
		return 0
	}
	var Organisation Organisation
	err = json.Unmarshal(OrgAsBytes, &Organisation)
	if err != nil {
		return 0
	}
	if Organisation.IsVerified == "True" {
        return 1
    }
	return 0
}

// Org Verify Password
func OrgVerifyPassword(stub shim.ChaincodeStubInterface, OrgID string, Password string) int {
	OrgAsBytes, err := stub.GetState(OrgID)
	if err != nil {
		return 0
	}else if OrgAsBytes == nil{
		return 0
	}
	var Organisation Organisation
	err = json.Unmarshal(OrgAsBytes, &Organisation)
	if err != nil {
		return 0
	}
	if Organisation.PasswordHash == sha256.Sum256([]byte(Password)) {
        return 1
    }
	return 0
}

// User Verify Password
func UserVerifyPassword(stub shim.ChaincodeStubInterface, UserID string, Password string) int {
	UserAsBytes, err := stub.GetState(UserID)
	if err != nil {
		return 0
	}else if UserAsBytes == nil{
		return 0
	}
	var User User
	err = json.Unmarshal(UserAsBytes, &User)
	if err != nil {
		return 0
	}
	if User.PasswordHash == sha256.Sum256([]byte(Password)) {
        return 1
    }
	return 0
}

// Adding info about a user
// args as : UserID, Username, Password, Email, Phonenumber
func  (t *HackidfChaincode) CreateUser(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var UserID = args[0]
	var Username = args[1]
	var Password = args[2]
	var Email = args[3]
	var Ph = args[4]
	var IsVerified = "False"
	PasswordHash := sha256.Sum256([]byte(Password))
	// checking for an error or if the user already exists
	UserAsBytes, err := stub.GetState(Username)
	if err != nil {
		return shim.Error("Failed to get Username:" + err.Error())
	}else if UserAsBytes != nil{
		return shim.Error("User with current username already exists")
	}

	var User = &User{Username:Username, PasswordHash:PasswordHash, Email:Email, Ph:Ph, IsVerified:IsVerified}
	UserJsonAsBytes, err :=json.Marshal(User)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(UserID, UserJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Creating User")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Do KYC for peer
// args as: UserID, Password(Hardcoded to "Password")
func  (t *HackidfChaincode) VerifyUser(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var UserID = args[0]
	var Password = args[1]
	PasswordHash := sha256.Sum256([]byte(Password))
	if PasswordHash != sha256.Sum256([]byte("Password")) {
		return shim.Error("WRONG PASSWORD ALERT!")
	}
	UserAsBytes, err := stub.GetState(UserID)
	if err != nil {
		return shim.Error("Failed to get User:" + err.Error())
	}else if UserAsBytes == nil{
		return shim.Error("Please give a valid User-ID")
	}
	if CheckUser(stub, UserID) == 1 {
		return shim.Error("The user is already verified")
	}
	var User User
	err = json.Unmarshal(UserAsBytes, &User)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	User.IsVerified = "True"
	UserJsonAsBytes, err :=json.Marshal(User)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(UserID, UserJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("VERIFIED!!")
	return shim.Success(nil)
}

// Adding info about an Organisations
// args as: OrgID, OrgName, Password(of Org)
func  (t *HackidfChaincode) CreateOrg(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var OrgID = args[0]
	var OrgName = args[1]
	var Password = args[2]
	var IsVerified = "False"
	PasswordHash := sha256.Sum256([]byte(Password))
	// checking for an error or if the user already exists
	OrgAsBytes, err := stub.GetState(OrgID)
	if err != nil {
		return shim.Error("Failed to get Organisation:" + err.Error())
	}else if OrgAsBytes != nil{
		return shim.Error("Organisation is already registered")
	}
	var Organisation = &Organisation{OrgName:OrgName, PasswordHash:PasswordHash, IsVerified:IsVerified}
	OrgJsonAsBytes, err :=json.Marshal(Organisation)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(OrgID, OrgJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Creating Organisation")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Verify Organisation
// args as: OrgID, Password(Hardcoded to "Password")
func  (t *HackidfChaincode) VerifyOrg(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var OrgID = args[0]
	var Password = args[1]
	PasswordHash := sha256.Sum256([]byte(Password))
	if PasswordHash != sha256.Sum256([]byte("Password")) {
		return shim.Error("WRONG PASSWORD ALERT!")
	}
	OrgAsBytes, err := stub.GetState(OrgID)
	if err != nil {
		return shim.Error("Failed to get Organisation:" + err.Error())
	}else if OrgAsBytes == nil{
		return shim.Error("Organisation not registered")
	}
	var Organisation Organisation
	err = json.Unmarshal(OrgAsBytes, &Organisation)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	Organisation.IsVerified = "True"
	OrgJsonAsBytes, err :=json.Marshal(Organisation)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(OrgID, OrgJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("VERIFIED!!")
	return shim.Success(nil)
}

// Make Claim
// args as: Hash(string only :P), UserID, Password, OrgID, Skill, Timeline(Eg. 2017-18)
func  (t *HackidfChaincode) MakeClaim(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var Hash = args[0]
	var UserID = args[1]
	var Password = args[2]
	var OrgID = args[3]
	var Skill = args[4]
	var Timestamp = args[5]
	var IsVerified = "False"
	if CheckUser(stub, UserID) == 0 {
		return shim.Error("Please finish your KYC procedure with InfoEaze.")
	}
	if CheckOrg(stub, OrgID) == 0 {
		return shim.Error("Org isn't verified by InfoEaze.")
	}
	if UserVerifyPassword(stub, UserID, Password) == 0 {
		return shim.Error("Password doesn't match user")
	}
	var Claim = &Claim{UserID:sha256.Sum256([]byte(UserID)), OrgID:sha256.Sum256([]byte(OrgID)), Skill:sha256.Sum256([]byte(Skill)), Comments:"NIL", Timestamp:sha256.Sum256([]byte(Timestamp)) ,IsVerified:IsVerified}
	ClaimJsonAsBytes, err :=json.Marshal(Claim)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(Hash, ClaimJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Making Claim")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Verify Claim
// args as: Hash, OrgID, Password(of Organisation)
func  (t *HackidfChaincode) VerifyClaim(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var Hash = args[0]
	var OrgID = args[1]
	var Password = args[2]
	if OrgVerifyPassword(stub, OrgID, Password) == 0 {
		return shim.Error("Password doesn't match Organisation")
	}
	ClaimAsBytes, err := stub.GetState(Hash)
	if err != nil {
		return shim.Error("Failed to get Claim:" + err.Error())
	}else if ClaimAsBytes == nil{
		return shim.Error("Claim not made")
	}
	var Claim Claim
	err = json.Unmarshal(ClaimAsBytes, &Claim)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	if Claim.OrgID != sha256.Sum256([]byte(OrgID)) {
		return shim.Error("Not authorised to verify this claim")
	}
	Claim.IsVerified = "True"
	ClaimJsonAsBytes, err :=json.Marshal(Claim)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(Hash, ClaimJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("VERIFIED!!")
	return shim.Success(nil)
}

// Issue Certificate
// args as: Hash(Token), OrgID, Password, Certificate(a string itself)
func  (t *HackidfChaincode) IssueCert(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var Hash = args[0]
	var OrgID = args[1]
	var Password = args[2]
	var Certificate = args[3]
	if CheckOrg(stub, OrgID) == 0 {
		return shim.Error("Org isn't verified by InfoEaze.")
	}
	if OrgVerifyPassword(stub, OrgID, Password) == 0 {
		return shim.Error("Password doesn't match Organisation")
	}
	var Cert = &Cert{OrgID:OrgID, CertHash:sha256.Sum256([]byte(Certificate))}
	CertJsonAsBytes, err :=json.Marshal(Cert)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(Hash, CertJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Making Claim")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Query Function
// args as: any_string
func  (t *HackidfChaincode) Query(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	DataAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Error encountered")
	}else if DataAsBytes == nil {
		return shim.Error("No Data")
	}
	return shim.Success(DataAsBytes)
}

// Query Function when cert is given
// args as: Hash(basically token), Certificate(basically string)
func  (t *HackidfChaincode) QueryWithCert(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	DataAsBytes, err := stub.GetState(args[0])
	var Certificate = args[1]
	if err != nil {
		return shim.Error("Error encountered")
	}else if DataAsBytes == nil {
		return shim.Error("Hash Token provided is not valid.")
	}
	var Cert Cert
	err = json.Unmarshal(DataAsBytes, &Cert)
	if Cert.CertHash != sha256.Sum256([]byte(Certificate)){
		return shim.Error("Certificate doesn't match with details of hash")
	}
	return shim.Success(DataAsBytes)
}

// MAIN FUNCTION
func  main() {
	err := shim.Start(new(HackidfChaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}