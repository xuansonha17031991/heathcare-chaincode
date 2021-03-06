package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type HeathCare_Chaincode struct {
}

type PatientInformation struct {
	ObjectType                   string `json:"doctype"`
	ID                           string `json:"photo_id"`
	InsuranceCard                string `json:"insurance_card"`
	CurrentMedicationInformation string `json:"current_medication_information"`
	RelatedMedicalRecords        string `json:"related_medical_records"`
	MakeNoteOfAppointmentDate    string `json:"make_note_of_appointment_date"`
}

type MedicalRecord struct {
	ObjectType                        string `json:"doctype"`
	ID                                string `json:"id"`
	PersonalIdentificationInformation string `json:"personal_identification"`
	MedicalHistory                    string `json:"medical_history"`
	FamilyMedicalHistory              string `json:"family_medical_history"`
	MedicationHistory                 string `json:"medication_history"`
	TreatmentHistory                  string `json:"treatment_history"`
	MedicalDirectives                 string `json:"medical_directives"`
}

type DrugInformation struct {
	ObjectType     string `json:"doctype"`
	ID             string `json:"id"`
	PatientName    string `json:"patient_name"`
	DrugName       string `json:"drug_name"`
	ExpirationDate string `json:"expiration_date"`
	Quantity       string `json:"quantity"`
	PrescribedBy   string `json:"prescribed_by"`
}

type HospitalFees struct {
	ObjectType               string `json:"docType"`
	ID                       string `json:"id"`
	PatientName              string `json:"patient_name"`
	Account                  string `json:"account"`
	DateOfService            string `json:"date_of_service"`
	PatientService           string `json:"patient_service"`
	PrimaryInsuranceBilled   string `json:"primary_insurance_billed"`
	SecondaryInsuranceBilled string `json:"secondary_insurance_billed"`
	Pharmacy                 string `json:"pharmacy"`
	Room                     string `json:"room"`
	AmountDue                string `json:"amount_due"`
}

type Query struct {
	ObjectType string `json:"docType"`
	UserID     string `json:"userid"`
	PatientID  string `json:"patientid"`
	Location   string `json:"location"`
	Time       string `json:"time"`
	Purpose    string `json:"purpose"`
}

/*main*/
func main() {
	err := shim.Start(new(HeathCare_Chaincode))
	if err != nil {
		fmt.Printf("cannot initiate heathcare chaincode: %s", err)
	}
}

// Init chaincode
func (t *HeathCare_Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke
func (t *HeathCare_Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke is running" + function)

	switch function {
	case "createDrugInformation":
		return t.createDrugInformation(stub, args)
	case "createHospitalFees":
		return t.createHospitalFees(stub, args)
	case "createMedicalRecord":
		return t.createMedicalRecord(stub, args)
	case "createPatientInformation":
		return t.createPatientInformation(stub, args)
	case "historyModify":
		return t.historyModify(stub, args)
	case "historyQuery":
		return t.historyQuery(stub, args)
	case "modifyDrugData":
		return t.modifyDrugData(stub, args)
	case "modifyMedicalData":
		return t.modifyMedicalData(stub, args)
	case "modifyPatientInformation":
		return t.modifyPatientInformation(stub, args)
	case "query":
		return t.query(stub, args)

	default:
		fmt.Println("Invoke did not find function: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

//###########################################################################################
/**
 * modify data of medical record and store with id of user execute query
 * @param: patientId
 * @param: personalIdentificationInformation
 * @param: medicalHistory
 * @param: familyMedicalHistory
 * @param: medicationHistory
 * @param: treatmentHistory
 * @param: medicalDirectives
 */
//###########################################################################################
func (t *HeathCare_Chaincode) createMedicalRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createMedicalRecord function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	if len(args) != 7 {
		return shim.Error("there must be 7 argument")
	}

	for i := 0; i < len(args); i++ {
		if len(args[i]) == 0 {
			return shim.Error("argument " + strconv.Itoa(i+1) + " must be declare")
		}
	}
	patientId := args[0]
	personalIdentificationInformation := args[1]
	medicalHistory := args[2]
	familyMedicalHistory := args[3]
	medicationHistory := args[4]
	treatmentHistory := args[5]
	medicalDirectives := args[6]

	//convert variable to json
	objectType := "MedicalRecord"
	medialRecord := &MedicalRecord{objectType, patientId, personalIdentificationInformation,
		medicalHistory, familyMedicalHistory, medicationHistory,
		treatmentHistory, medicalDirectives}

	//convert data to byte
	MedicalRecordAsByte, errMedicalRecordAsByte := json.Marshal(medialRecord)
	if errMedicalRecordAsByte != nil {
		return shim.Error(errMedicalRecordAsByte.Error())
	}

	//save to database
	errMedicalRecordAsByte = stub.PutPrivateData("MedicalRecordCollection", patientId, MedicalRecordAsByte)
	if errMedicalRecordAsByte != nil {
		return shim.Error(errMedicalRecordAsByte.Error())
	}

	//create index key
	indexName := "id"
	medicalRecordIndexKey, errMedicalRecordIndexKey := stub.CreateCompositeKey(indexName, []string{medialRecord.ID, medialRecord.PersonalIdentificationInformation, medialRecord.MedicalHistory, medialRecord.FamilyMedicalHistory, medialRecord.MedicationHistory, medialRecord.TreatmentHistory, medialRecord.MedicalDirectives})
	if errMedicalRecordIndexKey != nil {
		return shim.Error(errMedicalRecordIndexKey.Error())
	}

	//save index
	value := []byte{0x00}
	stub.PutPrivateData("MedicalRecordCollection", medicalRecordIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction createMedicalRecord")
	fmt.Printf("time start: %s", start.String())
	fmt.Printf("time end: %s", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end createMedicalRecord function ===============")
	return shim.Success(nil)
}

//###########################################################################################
/**
 * modify drug's data of patient
 * @param: userid
 * @param: patientid
 * @param: location
 * @param: collection
 * @param: newPatientName
 * @param: newDrugName
 * @param: newExpirationDate
 * @param: newQuantity
 * @param: newPrescribedBy
 */
//###########################################################################################
func (t *HeathCare_Chaincode) modifyDrugData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start modifyDrugData function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	var jsonResp string

	if len(args) != 9 {
		return shim.Error("expecting 9 argument")
	}

	userid := args[0]
	patientid := args[1]
	location := args[2]
	collection := args[3]

	newPatientName := args[4]
	newDrugName := args[5]
	newExpirationDate := args[6]
	newQuantity := args[7]
	newPrescribedBy := args[8]

	timeQuery := time.Now().String()

	//get user identity before query
	userIdentityAsBytes, errUserIdentityAsByte := stub.GetPrivateData(collection, userid)
	if errUserIdentityAsByte != nil {
		return shim.Error("cannot get user identity")
	} else if userIdentityAsBytes == nil {
		return shim.Error("user does not exist")
	}

	objectType := "Query"
	query := &Query{objectType, userid, patientid, location, timeQuery, "modify"}
	queryAsByte, errQueryAsByte := json.Marshal(query)
	if errQueryAsByte != nil {
		return shim.Error(errQueryAsByte.Error())
	}

	//save to database
	errQueryAsByte = stub.PutPrivateData("modifyCollection", userid, queryAsByte)
	if errQueryAsByte != nil {
		return shim.Error(errQueryAsByte.Error())
	}

	//create index key
	indexName := "userid~patientid"
	queryIndexKey, errQueryIndexKey := stub.CreateCompositeKey(indexName, []string{query.UserID, query.PatientID, query.Location, query.Purpose})
	if errQueryIndexKey != nil {
		return shim.Error(errQueryIndexKey.Error())
	}

	//save index
	value := []byte{0x00}
	stub.PutPrivateData("modifyCollection", queryIndexKey, value)

	//get data
	drugAsBytes, errDrugAsByte := stub.GetPrivateData("drugInformationCollection", patientid)
	if errDrugAsByte != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + patientid + ": " + errDrugAsByte.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if errDrugAsByte == nil {
		return shim.Error("patient's data does not exist")
	}

	//convert data of patient to json
	drug := &DrugInformation{}
	errDrugAsByte = json.Unmarshal(drugAsBytes, drug)

	//change data
	drug.PatientName = newPatientName
	drug.DrugName = newDrugName
	drug.ExpirationDate = newExpirationDate
	drug.Quantity = newQuantity
	drug.PrescribedBy = newPrescribedBy

	//store new data of drug information
	newDrugInformationAsByte, errNewDruvInformationAsByte := json.Marshal(drug)
	errNewDruvInformationAsByte = stub.PutPrivateData("drugInformationCollection", patientid, newDrugInformationAsByte)
	if errNewDruvInformationAsByte != nil {
		return shim.Error("cannot store new drug's data")
	}

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("function modifyDrugData")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end modifyDrugData function ===============")

	return shim.Success(nil)
}

//###########################################################################################
/**
 * modify drug's data of patient
 * @param: patientid
 * @param: patientName
 * @param: drugName
 * @param: expirationDate
 * @param: quantity
 * @param: prescribedBy
 */
//###########################################################################################
func (t *HeathCare_Chaincode) createDrugInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createDrugInformation function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	if len(args) != 6 {
		return shim.Error("expecting 6 argument")
	}

	for i := 0; i < len(args); i++ {
		if len(args[i]) == 0 {
			return shim.Error("argument " + strconv.Itoa(i+1) + " must be declare")
		}
	}

	//define argument
	patientId := args[0]
	patientName := args[1]
	drugName := args[2]
	expirationDate := args[3]
	quantity := args[4]
	prescribedBy := args[5]

	//convert to json
	objectType := "DrugInformation"
	drugInformation := &DrugInformation{objectType, patientId, patientName, drugName,
		expirationDate, quantity, prescribedBy}
	drugInformationAsByte, errDrugInformationAsByte := json.Marshal(drugInformation)
	if errDrugInformationAsByte != nil {
		return shim.Error(errDrugInformationAsByte.Error())
	}

	//save to ledger
	errDrugInformationAsByte = stub.PutPrivateData("DrugInformationCollection", patientId, drugInformationAsByte)
	if errDrugInformationAsByte != nil {
		return shim.Error(errDrugInformationAsByte.Error())
	}

	//create and save key
	indexName := "id~patient_name"
	DrugInformationIndexKey, errDrugInformationIndexKey := stub.CreateCompositeKey(indexName, []string{drugInformation.ID, drugInformation.PatientName, drugInformation.DrugName, drugInformation.ExpirationDate, drugInformation.Quantity, drugInformation.ExpirationDate})
	if errDrugInformationIndexKey != nil {
		return shim.Error(errDrugInformationIndexKey.Error())
	}
	value := []byte{0x00}
	stub.PutPrivateData("DrugInformationCollection", DrugInformationIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction createDrugInformation")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end createDrugInformation function ===============")

	return shim.Success(nil)
}

//###########################################################################################
/**
 * create patient information
 * @param: patientid
 * @param: insuranceCard
 * @param: currentMedicationInformation
 * @param: relatedMedicalRecords
 * @param: makeNoteOfAppointmentDate
 */
//###########################################################################################
func (t *HeathCare_Chaincode) createPatientInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createPatientInformation function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	if len(args) != 5 {
		return shim.Error("there must be 5 argument")
	}

	for i := 0; i < len(args); i++ {
		if len(args[i]) == 0 {
			return shim.Error("argument " + strconv.Itoa(i+1) + " must be declare")
		}
	}

	patientId := args[0]
	insuranceCard := args[1]
	currentMedicationInformation := args[2]
	relatedMedicalRecords := args[3]
	makeNoteOfAppointmentDate := args[4]

	//convert variable to json
	objectType := "PatientInformation"
	// user := &User{}
	patient := &PatientInformation{objectType, patientId, insuranceCard,
		currentMedicationInformation, relatedMedicalRecords, makeNoteOfAppointmentDate}
	// patientData := user.Name + " " + strconv.Itoa(user.Age) + " " + user.Number + " " + user.Address + " " + patient.Data
	PatientInformationAsByte, errPatientInformationAsByte := json.Marshal(patient)
	if errPatientInformationAsByte != nil {
		return shim.Error(errPatientInformationAsByte.Error())
	}

	//save to database
	errPatientInformationAsByte = stub.PutPrivateData("PatientInformationCollection", patientId, PatientInformationAsByte)
	if errPatientInformationAsByte != nil {
		return shim.Error(errPatientInformationAsByte.Error())
	}

	//create index key
	indexName := "id~insurance_card"
	patientIndexKey, errPatientIndexKey := stub.CreateCompositeKey(indexName, []string{patient.ID, patient.InsuranceCard, patient.CurrentMedicationInformation, patient.RelatedMedicalRecords, patient.MakeNoteOfAppointmentDate})
	if errPatientIndexKey != nil {
		return shim.Error(errPatientIndexKey.Error())
	}

	//save index
	value := []byte{0x00}
	stub.PutPrivateData("PatientInformationCollection", patientIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction createPatientInformation")
	fmt.Printf("time start: %s", start.String())
	fmt.Printf("time end: %s", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end createPatientInformation function ===============")
	return shim.Success(nil)
}

//###########################################################################################
/**
 * create hospital fees of patients
 * @param: userid
 * @param: patientName
 * @param: account
 * @param: dateOfService
 * @param: patientService
 * @param: primaryInsuranceBilled
 * @param: secondaryInsuranceBilled
 * @param: pharmacy
 * @param: room
 * @param: amountDue
 * ouput: nil
 */
//###########################################################################################
func (t *HeathCare_Chaincode) createHospitalFees(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start createHospitalFees function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	//check length of data
	if len(args) != 3 {
		return shim.Error("expecting 3 argument")
	}

	//define data variable
	id := args[0]
	patientName := args[1]
	account := args[2]
	dateOfService := args[3]
	patientService := args[4]
	primaryInsuranceBilled := args[5]
	secondaryInsuranceBilled := args[6]
	pharmacy := args[7]
	room := args[8]
	amountDue := args[9]

	ObjectType := "HospitalFees"
	hospitalFees := &HospitalFees{ObjectType, id, patientName, account, dateOfService,
		patientService, primaryInsuranceBilled, secondaryInsuranceBilled, pharmacy,
		room, amountDue}

	//marshal delivery to byte
	hospitalFeesAsByte, errHospitalFeesAsByte := json.Marshal(hospitalFees)
	if errHospitalFeesAsByte != nil {
		return shim.Error("cannot marshal pharmacy's data")
	}

	//put data to ledger
	errHospitalFeesAsByte = stub.PutPrivateData("HospitalFeesCollection", id, hospitalFeesAsByte)
	if errHospitalFeesAsByte != nil {
		return shim.Error("cannot put private data of pharmacy")
	}

	//create index key
	indexKey := "id~patient_name"
	hospitalFeesIndexKey, errHospitalFeesIndexKey := stub.CreateCompositeKey(indexKey, []string{hospitalFees.ID, hospitalFees.PatientName, hospitalFees.Account, hospitalFees.DateOfService, hospitalFees.PatientService, hospitalFees.PrimaryInsuranceBilled, hospitalFees.SecondaryInsuranceBilled, hospitalFees.Pharmacy, hospitalFees.Room, hospitalFees.AmountDue})
	if errHospitalFeesIndexKey != nil {
		return shim.Error("cannot create index key of delivery")
	}

	//save key
	value := []byte{0x00}
	stub.PutPrivateData("HospitalFeesCollection", hospitalFeesIndexKey, value)

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction createHospitalFees")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end createHospitalFees function ===============")
	time.Sleep(time.Second)

	return shim.Success(nil)
}

//###########################################################################################
/**
 * modify data of medical record and store with id of user execute query
 * @param: userid
 * @param: patientid
 * @param: location
 * @param: collection of user execute query
 * @param: newPersonalIdentificationInformation
 * @param: newMedicalHistory
 * @param: newFamilyMedicalHistory
 */
//###########################################################################################
func (t *HeathCare_Chaincode) modifyMedicalData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start modifyMedicalData function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	var jsonResp string

	if len(args) != 8 {
		return shim.Error("expecting 4 argument")
	}

	//define identity of query-er and new value of medical record
	userid := args[0]
	patientid := args[1]
	location := args[2]
	collection := args[3]

	newPersonalIdentificationInformation := args[4]
	newMedicalHistory := args[5]
	newFamilyMedicalHistory := args[6]
	newMedicationHistory := args[7]
	newTreatmentHistory := args[8]
	newMedicalDirectives := args[9]
	timeQuery := time.Now().String()

	//get user identity before query
	userIdentityAsBytes, errUserIdentityAsByte := stub.GetPrivateData(collection, userid)
	if errUserIdentityAsByte != nil {
		return shim.Error("cannot get user identity")
	} else if userIdentityAsBytes == nil {
		return shim.Error("user does not exist")
	}

	//create query object with purpose: modify
	objectType := "Query"
	query := &Query{objectType, userid, patientid, location, timeQuery, "modify"}
	queryAsByte, errQueryAsByte := json.Marshal(query)
	if errQueryAsByte != nil {
		return shim.Error(errQueryAsByte.Error())
	}

	//save to database
	errQueryAsByte = stub.PutPrivateData("modifyCollection", userid, queryAsByte)
	if errQueryAsByte != nil {
		return shim.Error(errQueryAsByte.Error())
	}

	//create index key
	indexName := "userid~patientid"
	queryIndexKey, errQueryIndexKey := stub.CreateCompositeKey(indexName, []string{query.UserID, query.PatientID, query.Location, query.Purpose})
	if errQueryIndexKey != nil {
		return shim.Error(errQueryIndexKey.Error())
	}

	//save index
	value := []byte{0x00}
	stub.PutPrivateData("modifyCollection", queryIndexKey, value)

	//get medical record data
	medicalRecordAsBytes, errMedicalRecordAsByte := stub.GetPrivateData("MedicalRecordCollection", patientid)
	if errMedicalRecordAsByte != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + patientid + ": " + errMedicalRecordAsByte.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if errMedicalRecordAsByte == nil {
		return shim.Error("patient's data does not exist")
	}

	//convert data of patient to json
	medicalRecord := &MedicalRecord{}
	errMedicalRecordAsByte = json.Unmarshal(medicalRecordAsBytes, medicalRecord)

	//change data
	medicalRecord.PersonalIdentificationInformation = newPersonalIdentificationInformation
	medicalRecord.MedicalHistory = newMedicalHistory
	medicalRecord.FamilyMedicalHistory = newFamilyMedicalHistory
	medicalRecord.MedicationHistory = newMedicationHistory
	medicalRecord.TreatmentHistory = newTreatmentHistory
	medicalRecord.MedicalDirectives = newMedicalDirectives

	//convert new medical record data to byte
	newMedicalRecordAsByte, errNewMedicalRecordAsByte := json.Marshal(medicalRecord)

	//store new data
	errNewMedicalRecordAsByte = stub.PutPrivateData("MedicalRecordCollection", patientid, newMedicalRecordAsByte)
	if errNewMedicalRecordAsByte != nil {
		return shim.Error("cannot save new medical record's data")
	}

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("function modifyMedicalData")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end modifyMedicalData function ===============")

	return shim.Success(nil)
}

//###########################################################################################
/**
 * modify data of patient and save id of user execute query
 * @param: userid
 * @param: patientid
 * @param: location
 * @param: collection
 * @param: newInsuranceCard
 * @param: newCurrentMedicationInformation
 * @param: newRelatedMedicalRecords
 * @param: newmakeNoteOfAppointmentDate
 */
//###########################################################################################
func (t *HeathCare_Chaincode) modifyPatientInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start modifyPatientInformation function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	var jsonResp string

	if len(args) != 8 {
		return shim.Error("expecting 8 argument")
	}

	userid := args[0]
	patientid := args[1]
	location := args[2]
	collection := args[3]

	newInsuranceCard := args[4]
	newCurrentMedicationInformation := args[5]
	newRelatedMedicalRecords := args[6]
	newmakeNoteOfAppointmentDate := args[7]
	timeQuery := time.Now().String()

	//get user identity before query
	userIdentityAsBytes, errUserIdentityAsByte := stub.GetPrivateData(collection, userid)
	if errUserIdentityAsByte != nil {
		return shim.Error("cannot get user identity")
	} else if userIdentityAsBytes == nil {
		return shim.Error("user does not exist")
	}

	objectType := "Query"
	query := &Query{objectType, userid, patientid, location, timeQuery, "modify"}
	queryAsByte, errQueryAsByte := json.Marshal(query)
	if errQueryAsByte != nil {
		return shim.Error(errQueryAsByte.Error())
	}

	//save to database
	errQueryAsByte = stub.PutPrivateData("modifyCollection", userid, queryAsByte)
	if errQueryAsByte != nil {
		return shim.Error(errQueryAsByte.Error())
	}

	//create index key
	indexName := "userid~patientid"
	queryIndexKey, errQueryIndexKey := stub.CreateCompositeKey(indexName, []string{query.UserID, query.PatientID, query.Location, query.Purpose})
	if errQueryIndexKey != nil {
		return shim.Error(errQueryIndexKey.Error())
	}

	//save index
	value := []byte{0x00}
	stub.PutPrivateData("modifyCollection", queryIndexKey, value)

	//get data
	patientAsBytes, errPatientAsByte := stub.GetPrivateData("PatientInformationCollection", patientid)
	if errPatientAsByte != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + patientid + ": " + errPatientAsByte.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if errPatientAsByte == nil {
		return shim.Error("patient's data does not exist")
	}

	//convert data of patient to json
	patient := &PatientInformation{}
	errPatientAsByte = json.Unmarshal(patientAsBytes, patient)

	//change data
	patient.InsuranceCard = newInsuranceCard
	patient.CurrentMedicationInformation = newCurrentMedicationInformation
	patient.RelatedMedicalRecords = newRelatedMedicalRecords
	patient.MakeNoteOfAppointmentDate = newmakeNoteOfAppointmentDate

	patientAsByte, errPatientAsByte := json.Marshal(patient)

	errPatientAsByte = stub.PutPrivateData("PatientInformationCollection", patientid, patientAsByte)
	if errPatientAsByte != nil {
		return shim.Error("cannot patient's data")
	}

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("function modifyPatientInformation")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end modifyPatientInformation function ===============")

	return shim.Success(nil)
}

//###########################################################################################
/**
 * view history query of user
 * params: userid
 * params: patientid
 * params: queryCollection
 * params: modifyCollection
 */
//###########################################################################################
func (t *HeathCare_Chaincode) historyQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start historyQuery function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	// check require argument
	if len(args) != 1 {
		return shim.Error("expecting 1 argument")
	}

	for i := 0; i < len(args); i++ {
		if len(args[i]) == 0 {
			return shim.Error("argument " + strconv.Itoa(i+1) + " must be delare")
		}
	}

	//define argument
	userid := args[0]

	//query
	queryDataAsBytes, errQueryDataAsByte := stub.GetPrivateData("queryCollection", userid)
	if errQueryDataAsByte != nil {
		return shim.Error("cannot get data of query")
	} else if queryDataAsBytes == nil {
		return shim.Error("query data history does not exist")
	}

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction historyQuery")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end historyQuery function ===============")
	time.Sleep(time.Second)

	return shim.Success(queryDataAsBytes)
}

//###########################################################################################
/**
 * view history modify data of user
 * params: userid
 * params: patientid
 * params: queryCollection
 * params: modifyCollection
 */
//###########################################################################################
func (t *HeathCare_Chaincode) historyModify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start historyModify function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	// check require argument
	if len(args) != 1 {
		return shim.Error("expecting 1 argument")
	}

	for i := 0; i < len(args); i++ {
		if len(args[i]) == 0 {
			return shim.Error("argument " + strconv.Itoa(i+1) + " must be delare")
		}
	}

	//define argument
	userid := args[0]

	//modify
	modifyDataAsBytes, errModifyDataAsByte := stub.GetPrivateData("modifyCollection", userid)
	if errModifyDataAsByte != nil {
		return shim.Error("cannot get modify data")
	} else if modifyDataAsBytes == nil {
		return shim.Error("err2")
	}

	end := time.Now()
	elapsed := time.Since(start)

	fmt.Println("\nfunction historyModify")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end historyModify function ===============")
	time.Sleep(time.Second)

	return shim.Success(modifyDataAsBytes)
}

//###########################################################################################
/**
 * modify data of medical record and store with id of user execute query
 * @param: userid
 * @param: patientid
 * @param: location
 * @param: collection
 */
//###########################################################################################
func (t *HeathCare_Chaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("\n=============== start query function ===============")
	start := time.Now()
	time.Sleep(time.Second)

	var jsonResp string

	if len(args) != 4 {
		return shim.Error("expecting 4 argument")
	}

	userid := args[0]
	patientid := args[1]
	location := args[2]
	collection := args[3]
	timeQuery := time.Now().String()

	//get user identity before query
	userIdentityAsBytes, errUserIdentityAsByte := stub.GetPrivateData(collection, userid)
	if errUserIdentityAsByte != nil {
		return shim.Error("cannot get user identity")
	} else if userIdentityAsBytes == nil {
		return shim.Error("user does not exist")
	}

	objectType := "Query"
	query := &Query{objectType, userid, patientid, location, timeQuery, "query"}
	queryAsByte, errQueryAsByte := json.Marshal(query)
	if errQueryAsByte != nil {
		return shim.Error(errQueryAsByte.Error())
	}

	//save to database
	errQueryAsByte = stub.PutPrivateData("queryCollection", userid, queryAsByte)
	if errQueryAsByte != nil {
		return shim.Error(errQueryAsByte.Error())
	}

	//create index key
	indexName := "userid~patientid"
	queryIndexKey, errQueryIndexKey := stub.CreateCompositeKey(indexName, []string{query.UserID, query.PatientID, query.Location, query.Purpose})
	if errQueryIndexKey != nil {
		return shim.Error(errQueryIndexKey.Error())
	}

	//save index
	value := []byte{0x00}
	stub.PutPrivateData("queryCollection", queryIndexKey, value)

	//get data
	valueAsBytes, errValueAsByte := stub.GetPrivateData("PatientInformationCollection", patientid)
	if errValueAsByte != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + patientid + ": " + errValueAsByte.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valueAsBytes == nil {
		jsonResp = "user id does not exist"
		return shim.Error(jsonResp)
	}

	end := time.Now()
	elapsed := time.Since(start)
	fmt.Println("function query")
	fmt.Println("time start: ", start.String())
	fmt.Println("time end: ", end.String())
	fmt.Println("time execute: ", elapsed.String())
	fmt.Println("=============== end query function ===============")

	return shim.Success(valueAsBytes)
}
