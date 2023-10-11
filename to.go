package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	//"github.com/hyperledger/fabric-chaincode-go/shim"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//const (
//	minUnicodeRuneValue   = 0            //U+0000
//)

// SmartContract Chaincode implementation
type SmartContract struct {
	contractapi.Contract
}

// Account describes basic details of what makes up a car
type Account struct {
	//Id   string `json:"id"`
	Name          string  `json:"name"`
	SavingBalance float64 `json:"savingbalance"`

	Info   string `json:"info"`
	Status int    `json:"status"`
}
type Transaction struct {
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
	Transfer string  `json:"transfer"`
}

// QueryTransactionResult structure used for handling result of query
type QueryTransactionResult struct {
	Key    string `json:"Key"`
	Record *Transaction
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Account
}

func (t *SmartContract) Init(ctx contractapi.TransactionContextInterface, number int) error {
	// Initialize the chaincode
	fmt.Printf("Number%d Account Init\n", number)
	// Write the state to the ledger
	for j := 0; j <= 2; j++ {
		for i := 0; i <= number; {
			Name := "m" + strconv.Itoa(j) + "_" + "a" + strconv.Itoa(i)
			accountAsBytes, err := ctx.GetStub().GetState(Name)
			if err != nil {
				return fmt.Errorf("Failed to get account: " + err.Error())
			} else if accountAsBytes != nil {
				fmt.Println("This account already exists: " + Name)
				return fmt.Errorf("This account already exists: " + Name)
			}
			SavingBalance := 0.0
			Info := ""
			Status := -1
			account := &Account{Name, SavingBalance, Info, Status}
			accountJSONasBytes, err := json.Marshal(account)
			if err != nil {
				return err
			}
			/*accountJSONasString := `{"name": "` + Name + `", "savingbalance": "` + strconv.Itoa(SavingBalance)+ `", "checkingbalance": ` + strconv.Itoa(CheckingBalance) + `, "status": "` + strconv.Itoa(Status) + `"}`
			accountJSONasBytes := []byte(accountJSONasString)*/

			err = ctx.GetStub().PutState("M"+strconv.Itoa(j)+"_"+"A"+strconv.Itoa(i), accountJSONasBytes)

			if err != nil {
				return fmt.Errorf("Failed to put to world state. %s", err.Error())
			} else {
				i += 1
			}
		}
	}

	return nil
}
func (t *SmartContract) PutFinal(ctx contractapi.TransactionContextInterface, key string, value string, op string) error {
	fmt.Printf("in:%s", op)
	account, err := t.QueryAccount(ctx, key)
	if err != nil {
		return err
	}
	old := account.SavingBalance
	v, _ := strconv.ParseFloat(value, 64)

	fmt.Println(account)
	switch op {
	case "+":
		account.SavingBalance += v
	case "-":
		account.SavingBalance -= v
	default:
		return fmt.Errorf("Unrecognized operation %s", op)
	}

	accountAsBytes, _ := json.Marshal(account)
	err = ctx.GetStub().PutState(key, accountAsBytes)
	if err != nil {
		return err
	}

	fmt.Sprintf("Successfully  %s :%v %s %d =%v \n", key, old, op, value, account.SavingBalance)
	return nil
}

//	func (t *SmartContract)  InitMiddle(ctx contractapi.TransactionContextInterface, middle_key string ) error {
//		params := []string{"QueryMiddleByCompositeKey", middle_key}
//		queryArgs := make([][]byte, len(params))
//		for i, arg := range params {
//			queryArgs[i] = []byte(arg)
//		}
//
//		response := ctx.GetStub().InvokeChaincode("from", queryArgs, "channel1")
//		if response.Status != 200 {
//			return fmt.Errorf("Failed to query chaincode. Got error: %s", response.Payload)
//		}
//		fmt.Printf("channel1  --> response.Status=%d,response.Payload = %s\n", response.Status,string(response.Payload))
//
//		// Initialize the chaincode
//		// Write the state to the ledger
//		middleAsBytes:=response.Payload
//
//		/*middle := Account{}
//		err := json.Unmarshal(middleAsBytes, &middle)
//		if err != nil {
//
//			return err
//		}
//		middleAsBytes, _ := json.Marshal(middle)*/
//		err :=ctx.GetStub().PutState(middle_key, middleAsBytes)
//		if err != nil {
//			return err
//		}
//
//
//
//		fmt.Printf("init %s status finished\n",middle_key)
//
//		return nil
//	}
func (t *SmartContract) InitAccount(ctx contractapi.TransactionContextInterface, number int) error {

	// Initialize the chaincode
	fmt.Printf("Number%d Account Init\n", number)
	// Write the state to the ledger
	for i := 0; i <= number; {
		Name := "account" + strconv.Itoa(i)
		accountAsBytes, err := ctx.GetStub().GetState(Name)
		if err != nil {
			return fmt.Errorf("Failed to get account: " + err.Error())
		} else if accountAsBytes != nil {
			fmt.Println("This account already exists: " + Name)
			return fmt.Errorf("This account already exists: " + Name)
		}
		SavingBalance := 0.0

		Info := ""
		Status := -1
		account := &Account{Name, SavingBalance, Info, Status}
		accountJSONasBytes, err := json.Marshal(account)
		if err != nil {
			return err
		}
		//accountJSONasString := `{"name": "` + Name + `", "savingbalance": "` + strconv.Itoa(SavingBalance)+ `", "checkingbalance": ` + strconv.Itoa(CheckingBalance) + `, "status": "` + strconv.Itoa(Status) + `"}`
		//accountJSONasBytes := []byte(accountJSONasString)

		err = ctx.GetStub().PutState("ACCOUNT"+strconv.Itoa(i), accountJSONasBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		} else {
			i += 1
		}
	}

	return nil
}
func (t *SmartContract) InitMiddle(ctx contractapi.TransactionContextInterface, number int) error {

	// Initialize the chaincode
	fmt.Printf("Number%d Middle Init\n", number)
	// Write the state to the ledger
	for i := 0; i < number; {
		Name := "middle" + strconv.Itoa(i)
		accountAsBytes, err := ctx.GetStub().GetState(Name)
		if err != nil {
			return fmt.Errorf("Failed to get account: " + err.Error())
		} else if accountAsBytes != nil {
			fmt.Println("This account already exists: " + Name)
			return fmt.Errorf("This account already exists: " + Name)
		}
		SavingBalance := 1000.0

		Info := ""
		Status := -1
		middle := &Account{Name, SavingBalance, Info, Status}
		middleJSONasBytes, err := json.Marshal(middle)
		if err != nil {
			return err
		}
		//accountJSONasString := `{"name": "` + Name + `", "savingbalance": "` + strconv.Itoa(SavingBalance)+ `", "checkingbalance": ` + strconv.Itoa(CheckingBalance) + `, "status": "` + strconv.Itoa(Status) + `"}`
		//accountJSONasBytes := []byte(accountJSONasString)

		err = ctx.GetStub().PutState("MIDDLE"+strconv.Itoa(i), middleJSONasBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		} else {
			i += 1
		}
	}

	return nil
}

//func (t *SmartContract) TransferOrdinary(ctx contractapi.TransactionContextInterface,	operation string,  key string, value string) error {
//
//	amount,_:=strconv.ParseFloat(value,64)
//
//	account, err := t.QueryAccount(ctx, key)
//	if err != nil {
//		return err
//	}
//
//	/* middle, err := t.QueryAccount(ctx, middle_key)
//	if err != nil {
//		return err
//	}
//	*/
//
//	switch operation {
//	case "+":
//		account.SavingBalance += amount
//	case "-":
//		account.SavingBalance -= amount
//	default:
//		return fmt.Errorf("Unrecognized operation %s", operation)
//	}
//
//
//	fmt.Printf("TransferOrdinary accout%v  amount=%d\n ",account,account.SavingBalance)
//	accountAsBytes, _ := json.Marshal(account)
//	err=ctx.GetStub().PutState(key, accountAsBytes)
//	if err != nil {
//		return err
//	}
//	/*middleAsBytes, _ := json.Marshal(middle)
//	err=ctx.GetStub().PutState(middle_key, middleAsBytes)
//	if err != nil {
//		return err
//	}*/
//	fmt.Sprintf("Successfully %s  %s  %d \n", key,operation, amount)
//
//	return nil
//}
//

//	func (t *SmartContract) TransferM(ctx contractapi.TransactionContextInterface, middle_key string) ([]string,error) {
//		//fmt.Printf(middle_key+"haha")
//		params := []string{"QueryCompositeKey", middle_key}
//		queryArgs := make([][]byte, len(params))
//		for i, arg := range params {
//			queryArgs[i] = []byte(arg)
//		}
//
//		response := ctx.GetStub().InvokeChaincode("from", queryArgs, "channel1")
//		if response.Status != 200 {
//			return nil,fmt.Errorf("Failed to query chaincode. Got error: %s", response.Payload)
//		}
//		fmt.Printf("channel1  --> response.Status=%d,response.Payload = %s\n", response.Status,string(response.Payload))
//		/*
//			channel1  --> response.Status=200,response.Payload = ["\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000013\u0000ACCOUNT23\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000016\u0000ACCOUNT6\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000016\u0000ACCOUNT9\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000020\u0000ACCOUNT6\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000022\u0000ACCOUNT15\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u00009\u0000ACCOUNT43\u0000"]
//		*/
//		transactionSetAsBytes:=response.Payload
//		//transactionSet := formatJSON(transactionSetAsBytes)
//
//
//		results := []string{}
//
//
//		err := json.Unmarshal(transactionSetAsBytes, &results)
//		//fmt.Printf("fuhejian ==%v   size= %d \n",results, len(results))
//
//
//		if err != nil {
//			return nil,err
//		}
//		var m_finalValue float64
//		//var a_finalValue int
//		for i:=0;i<len(results);i++{
//
//
//			compositeKey:=results[i]
//			//fmt.Printf("中继转出交易：%s\n,len=%d\n",compositeKey,len(compositeKey))
//			_, keyParts, splitKeyErr := ctx.GetStub().SplitCompositeKey(compositeKey)
//			if splitKeyErr != nil {
//				return nil,fmt.Errorf(splitKeyErr.Error())
//			}
//
//			// Retrieve the delta value and operation
//			op:=keyParts[1]
//			valueStr := keyParts[2]
//			to_account:=keyParts[3]
//
//			// Convert the value string and perform the operation
//			//value, _:= strconv.Atoi(valueStr)
//			value,_:=strconv.ParseFloat(valueStr,64)
//			m_finalValue-=value
//			//fmt.Printf("keyParts[0]=%s,keyParts[1]=%s,keyParts[2]=%s,finalval=%d\n",keyParts[0],keyParts[1],keyParts[2],value)
//			//keyParts[0]=MIDDLE2,keyParts[1]=+,keyParts[2]=17,finalval=17
//			/*account2, err := t.QueryAccount(ctx, to_account)//ACCOUNT49
//			if err != nil {
//				return nil,err
//			}
//			m_finalValue-=value
//
//			account2.SavingBalance += value
//			account2AsBytes, _ := json.Marshal(account2)
//			err=ctx.GetStub().PutState(to_account, account2AsBytes)*/
//
//			//fmt.Printf("账户 %s余额%d 向账户 %s 余额%d 转账%d\n", middle,account1.SavingBalance,account, account2.SavingBalance,amount)
//
//			typeName := "account~op~value"
//
//			// Create the composite key
//			compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{to_account, op,valueStr})
//
//			if compositeErr != nil {
//				return nil,fmt.Errorf("Could not create a composite key for %s: %s\n", typeName, compositeErr.Error())
//			}
//
//			compositePutErr :=ctx.GetStub().PutState(compositeKey, []byte{0x00})
//			if compositePutErr != nil {
//				return nil,fmt.Errorf("Could not put operation for %s in the ledger: %s", middle_key, compositePutErr.Error())
//			}
//			//fmt.Printf("compositeKey:%v\n",compositeKey)
//
//		}
//		middle, err := t.QueryAccount(ctx, middle_key)//ACCOUNT49
//		middle.SavingBalance += m_finalValue //-
//		middleAsBytes, _ := json.Marshal(middle)
//		err=ctx.GetStub().PutState(middle_key, middleAsBytes)
//		//fmt.Printf("middle.SavingBalance=%d,finalval=%d\n",middle.SavingBalance,m_finalValue)
//
//
//		//fmt.Printf("查询到的另一个通道的账户状态 %s \n", transactionSet)
//
//	   return results,nil
//
// }
func (t *SmartContract) TransferTransaction(ctx contractapi.TransactionContextInterface, middle_key string) ([]string, error) {
	//fmt.Printf(middle_key+"haha")
	params := []string{"QueryCompositeKey", middle_key}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode("fromckks", queryArgs, "channel1")
	if response.Status != 200 {
		return nil, fmt.Errorf("Failed to query chaincode. Got error: %s", response.Payload)
	}
	fmt.Printf("channel1  --> response.Status=%d,response.Payload = %s\n", response.Status, string(response.Payload)) //channel1  --> response.Status=200,response.Payload = ["ACCOUNT2:10.82","ACCOUNT8:15.77","ACCOUNT4:48.90"]

	/*
		channel1  --> response.Status=200,response.Payload = ["\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000013\u0000ACCOUNT23\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000016\u0000ACCOUNT6\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000016\u0000ACCOUNT9\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000020\u0000ACCOUNT6\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000022\u0000ACCOUNT15\u0000","\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u00009\u0000ACCOUNT43\u0000"]
	*/
	transactionSetAsBytes := response.Payload
	//transactionSet := formatJSON(transactionSetAsBytes)

	results := []string{}

	err := json.Unmarshal(transactionSetAsBytes, &results)
	//fmt.Printf("fuhejian ==%v   size= %d \n",results, len(results))//fuhejian ==[ACCOUNT2:10.82 ACCOUNT8:15.77 ACCOUNT4:48.90]   size= 3

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(results); i++ {

		temp := strings.Split(results[i], ":")
		destId := temp[0]
		amountstr := temp[1]
		fmt.Printf("中继转出交易：%v==%s|%s \n", results[i], destId, amountstr)

		// Convert the value string and perform the operation
		//value, _:= strconv.Atoi(valueStr)
		value, _ := strconv.ParseFloat(amountstr, 64)

		dest, err := t.QueryAccount(ctx, destId) //ACCOUNT49
		if err != nil {
			return nil, err
		}
		dest.SavingBalance += value

		destAsBytes, _ := json.Marshal(dest)
		err = ctx.GetStub().PutState(destId, destAsBytes)

		//fmt.Printf("账户 %s余额%d 向账户 %s 余额%d 转账%d\n", middle,account1.SavingBalance,account, account2.SavingBalance,amount)

		typeName := "M~V"

		// Create the composite key
		compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{middle_key, amountstr})

		if compositeErr != nil {
			return nil, fmt.Errorf("Could not create a composite key for %s: %s\n", typeName, compositeErr.Error())
		}

		compositePutErr := ctx.GetStub().PutState(compositeKey, []byte{0x00})
		if compositePutErr != nil {
			return nil, fmt.Errorf("Could not put operation for %s in the ledger: %s", middle_key, compositePutErr.Error())
		}
		//fmt.Printf("compositeKey:%v\n",compositeKey)//compositeKey:M~VMIDDLE013.83

	}

	//fmt.Printf("middle.SavingBalance=%d,finalval=%d\n",middle.SavingBalance,m_finalValue)

	return results, nil

}
func (t *SmartContract) Prune(ctx contractapi.TransactionContextInterface, middle_key string) (string, error) {

	// Retrieve the name of the variable to prune
	name := middle_key

	// Get all delta rows for the variable
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("M~V", []string{name})
	if deltaErr != nil {
		return "", fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return "", fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set computing final value while iterating and deleting each key
	//var finalVal float64
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return "", fmt.Errorf(nextErr.Error())
		}

		// Split the key into its composite parts
		_, keyParts, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)

		if splitKeyErr != nil {
			return "", fmt.Errorf(splitKeyErr.Error())
		}

		valueStr := keyParts[1]
		fmt.Print(valueStr)

		// Convert the value to a int
		//value, convErr := strconv.ParseFloat(valueStr,64)
		//if convErr != nil {
		//	return "",fmt.Errorf(convErr.Error())
		//}
		//
		//// Delete the row from the ledger
		////deltaRowDelErr := ctx.GetStub().DelState(responseRange.Key)
		////if deltaRowDelErr != nil {
		////	return "",fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
		////}
		//
		//// Add the value of the deleted row to the final aggregate
		//
		//
		//finalVal -= value

	}
	// Update the ledger with the final value
	//finalvalueStr:=strconv.FormatFloat(finalVal,'f',20,64)
	//updateRespErr := ctx.GetStub().PutState(name,[]byte(finalvalueStr))
	//if updateRespErr !=nil {
	//	return "",fmt.Errorf("Could not update the final value of the variable after pruning")
	//}
	//
	//fmt.Printf("Successfully pruned variable %s, final value is %f, %d rows pruned", name, finalVal, i)

	return "finalvalueStr", nil
}

func (t *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, middle_key int) error {
	//middleName~op~value~accountMIDDLE1+10ACCOUNT2
	fmt.Printf("kualian \n")
	params := []string{"QueryTransactionSet", "MIDDLE" + strconv.Itoa(middle_key)}

	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode("from", queryArgs, "channel1")
	if response.Status != 200 {
		return fmt.Errorf("Failed to query chaincode. Got error: %s", response.Payload)
	}
	//fmt.Printf("channel1  --> response.Status=%d,response.Payload = %s\n", response.Status,formatJSON(response.Payload))

	transactionSetAsBytes := response.Payload
	//fmt.Printf("获得事务大小 %d \n", len(formatJSON(transactionSetAsBytes)))
	results := []QueryTransactionResult{}

	err := json.Unmarshal(transactionSetAsBytes, &results)
	if err != nil {
		return err
	}
	//fmt.Printf("len=%d   result==>:%v\n", len(results),results)
	typeName := "middleName~op~value~account"
	for i := 0; i < len(results); i++ {

		//fmt.Printf("中继转出交易：%v\n",results[i].Record)
		transaction := results[i].Record
		account := transaction.Transfer
		middle := transaction.Receiver
		amount := transaction.Amount

		// Create the composite key
		compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{middle, "-", strconv.FormatFloat(amount, 'f', 2, 64), account})
		//fmt.Printf("compositeKey:%v\n",compositeKey)
		if compositeErr != nil {
			return fmt.Errorf("Could not create a composite key for %s: %s\n", typeName, compositeErr.Error())
		}
		/*account1, err := t.QueryAccount(ctx, middle)//"receiver": "MIDDLE0"
		if err != nil {
			return err
		}*/

		account2, err := t.QueryAccount(ctx, account) //ACCOUNT49
		if err != nil {
			return err
		}
		//account1.SavingBalance -= amount
		account2.SavingBalance += amount
		account2.Info += "{M" + strconv.Itoa(middle_key) + "转入" + strconv.FormatFloat(amount, 'f', 2, 64) + "};"
		//fmt.Printf("账户 %s余额%d 向账户 %s 余额%d 转账%d\n", middle,account1.SavingBalance,account, account2.SavingBalance,amount)

		compositePutErr := ctx.GetStub().PutState(compositeKey, []byte{0x00})
		if compositePutErr != nil {
			return fmt.Errorf("Could not put operation for %s in the ledger: %s", middle_key, compositePutErr.Error())
		}

		//fmt.Printf("Successfully M%d added %s%d to %s\n", middle_key,"+", amount, account2.Name)
		/*account1AsBytes, _ := json.Marshal(account1)
		err=ctx.GetStub().PutState(middle, account1AsBytes)
		if err != nil {
			return err
		}*/

		account2AsBytes, _ := json.Marshal(account2)
		err = ctx.GetStub().PutState(account, account2AsBytes)
		if err != nil {
			return err
		}

	}

	return nil
}
func (t *SmartContract) Get(ctx contractapi.TransactionContextInterface, account_key string) (*Account, error) {
	account, err := t.QueryAccount(ctx, account_key) //ACCOUNT49
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	fmt.Printf(account.Name)
	// Check we have a valid number of args

	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("account~op~value", []string{account_key})
	if deltaErr != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", account_key, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return nil, fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", account_key))
	}

	// Iterate through result set and compute final value
	var finalVal float64
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return nil, fmt.Errorf(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return nil, fmt.Errorf(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		//operation := keyParts[1]
		valueStr := keyParts[2]

		//value,_ := strconv.Atoi(valueStr)
		value, _ := strconv.ParseFloat(valueStr, 64)

		finalVal += value

	}

	account.SavingBalance += finalVal
	accountAsBytes, _ := json.Marshal(account)
	err = ctx.GetStub().PutState(account_key, accountAsBytes)
	fmt.Printf("middle.SavingBalance=%d,finalval=%d\n", account.SavingBalance, finalVal)

	//account := new(Account)
	//_ = json.Unmarshal(accountvalbytes, account)
	//fmt.Printf("查询到的另一个通道的账户状态 %s \n", transactionSet)

	return account, nil

}

func (t *SmartContract) QueryMiddleByCompositeKey(ctx contractapi.TransactionContextInterface, middle_key string) (float64, error) {
	//var err error
	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("middleName~op~value~account", []string{middle_key})
	if deltaErr != nil {
		return -1, fmt.Errorf("Could not retrieve value for %s: %s", middle_key, deltaErr.Error())
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return -1, fmt.Errorf("No variable by the name %s exists", middle_key)
	}
	// Iterate through result set and compute final value
	var finalVal float64
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return -1, fmt.Errorf(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return -1, fmt.Errorf(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		valueStr := keyParts[2]

		// Convert the value string and perform the operation
		//value, _:= strconv.Atoi(valueStr)
		value, _ := strconv.ParseFloat(valueStr, 64)
		finalVal -= value

	}
	fmt.Printf("finalval=%d\n", finalVal)

	return finalVal, nil

}

// Query callback representing the query of a chaincode
func (t *SmartContract) QueryAccount(ctx contractapi.TransactionContextInterface, accountId string) (*Account, error) {
	var err error
	// Get the state from the ledger
	accountvalbytes, err := ctx.GetStub().GetState(accountId)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if accountvalbytes == nil {
		return nil, fmt.Errorf("%s does not exist", accountId)
	}

	account := new(Account)
	_ = json.Unmarshal(accountvalbytes, account)

	return account, nil
}

// Delete  an entity from state
func (t *SmartContract) Delete(ctx contractapi.TransactionContextInterface, account string) error {

	// Delete the key from the state in ledger
	err := ctx.GetStub().DelState(account)
	if err != nil {
		return fmt.Errorf("Failed to delete state")
	}

	return nil
}

// QueryAllAccounts returns all cars found in world state
func (t *SmartContract) QueryAllAccounts(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		account := new(Account)
		_ = json.Unmarshal(queryResponse.Value, account)

		queryResult := QueryResult{Key: queryResponse.Key, Record: account}
		results = append(results, queryResult)
	}

	return results, nil
}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

func main() {
	cc, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting SmartContract chaincode: %s", err)
	}
}
