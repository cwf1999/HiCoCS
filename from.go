package main

import (
	"bytes"
	"encoding/json"
	"math"

	//"golang.org/x/text/date"

	//"unicode/utf8"

	//"errors"
	"fmt"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/tuneinsight/lattigo/v4/ckks"
	"github.com/tuneinsight/lattigo/v4/rlwe"

	//"github.com/ldsec/lattigo/ckks"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

var params, _ = ckks.NewParametersFromLiteral(ckks.PN13QP218)

// we need logQP = 600 and greater for 128-bit security
// var params = ckks.DefaultParams[ckks.PN12QP109] // logQP = 109
// var params = ckks.DefaultParams[ckks.PN13QP218] // // logQP = 218
// var params = ckks.DefaultParams[ckks.PN14QP438] // logQP = 438
// var params = ckks.DefaultParams[ckks.PN15QP880] // logQP = 880
var encoder = ckks.NewEncoder(params)

// Keys
var kgen = ckks.NewKeyGenerator(params)
var sk, pk = kgen.GenKeyPair()

// Relinearization key
var rlk = kgen.GenRelinearizationKey(sk, 1)

// Encryptor
var encryptor = ckks.NewEncryptor(params, pk)

// Decryptor
var decryptor = ckks.NewDecryptor(params, sk)

// Evaluator
var evaluator = ckks.NewEvaluator(params, rlwe.EvaluationKey{Rlk: rlk})

const slots = 4096 // Number of homomorphic operations

func enc_value(value float64) *rlwe.Ciphertext {
	// encrypt with pk : ciphertext = [pk[0]*u + m + e_0, pk[1]*u + e_1]
	// encrypt with sk : ciphertext = [-a*sk + m + e, a]
	//var values []float64
	values := make([]float64, 1)
	values[0] = value

	plaintext := encoder.EncodeNew(values, params.MaxLevel(), params.DefaultScale(), params.LogSlots())
	ciphertext1 := encryptor.EncryptNew(plaintext)

	return ciphertext1

}
func enc(values []float64) *rlwe.Ciphertext {
	// encrypt with pk : ciphertext = [pk[0]*u + m + e_0, pk[1]*u + e_1]
	// encrypt with sk : ciphertext = [-a*sk + m + e, a]

	plaintext := encoder.EncodeNew(values, params.MaxLevel(), params.DefaultScale(), params.LogSlots())
	ciphertext1 := encryptor.EncryptNew(plaintext)

	return ciphertext1

}
func dec_val(ciphertext *rlwe.Ciphertext) []float64 {

	valuesTest := encoder.Decode(decryptor.DecryptNew(ciphertext), params.LogSlots())
	realval := make([]float64, slots)

	for i := 0; i < int(slots); i++ {
		realval[i] = real(valuesTest[i])
	}

	return (realval)
}
func CiphertextToString(ciphertext *rlwe.Ciphertext) string {

	valuesTest := encoder.Decode(decryptor.DecryptNew(ciphertext), params.LogSlots())
	realval := make([]float64, slots)

	for i := 0; i < int(slots); i++ {
		realval[i] = real(valuesTest[i])
	}

	return (strconv.FormatFloat(realval[0], 'f', 10, 64))
}

// func StringToCiphertext(str string) (string ,*rlwe.Ciphertext){
func StringToCiphertext(str string) string {
	// encrypt with pk : ciphertext = [pk[0]*u + m + e_0, pk[1]*u + e_1]
	// encrypt with sk : ciphertext = [-a*sk + m + e, a]
	// 解密
	cyptde := crypto.
		FromBase64String(str).
		SetKey("dfertf12dfertf12").
		Aes().
		ECB().
		PKCS7Padding().
		Decrypt().
		ToString()
	//fmt.Printf("orig:%v\n",cyptde)

	//values:=make([]float64,1)
	//values[0],_=strconv.ParseFloat(str,64)
	//plaintext := encoder.EncodeNew(values, params.MaxLevel(), params.DefaultScale(), params.LogSlots())
	//ciphertext1 := encryptor.EncryptNew(plaintext)

	//return cyptde,ciphertext1
	return cyptde

}
func StringToCiphertextSingle(str string) *rlwe.Ciphertext {
	// encrypt with pk : ciphertext = [pk[0]*u + m + e_0, pk[1]*u + e_1]
	// encrypt with sk : ciphertext = [-a*sk + m + e, a]
	// 解密
	cyptde := crypto.
		FromBase64String(str).
		SetKey("dfertf12dfertf12").
		Aes().
		ECB().
		PKCS7Padding().
		Decrypt().
		ToString()
	fmt.Printf("orig:%v\n", cyptde)

	values := make([]float64, 1)
	values[0], _ = strconv.ParseFloat(cyptde, 64)
	plaintext := encoder.EncodeNew(values, params.MaxLevel(), params.DefaultScale(), params.LogSlots())
	ciphertext1 := encryptor.EncryptNew(plaintext)

	return ciphertext1

}

// SmartContract Chaincode implementation
type SmartContract struct {
	contractapi.Contract
}

// Account describes basic details of what makes up a car
type Account struct {
	//Id   string `json:"id"`
	Name          string  `json:"name"`
	SavingBalance float64 `json:"savingbalance"`
	//
	Info   string `json:"info"`
	Status int    `json:"status"`
}

// type CompositeAccount struct {
//
//	Key  string `json:"key"`
//	Ledger *TransactionInfo
//
// }
type PendingTransation struct {
	Receiver  string           `json:"receiver"`
	Cipertext *rlwe.Ciphertext `json:"cihpertext"`
}
type PendingPool struct {
	Key                string              `json:"key"`
	pendingTransations []PendingTransation `json:"pendingTransations"`
}
type Transaction struct {
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
	Transfer string  `json:"transfer"`
}

// QueryResult structure used for handling result of query
/*type QueryTransactionResult struct {
	Key    string `json:"Key"`
	Record *Transaction
}
*/

// QueryResult structure used for handling result of query
/*type QueryResult struct {
	Key    string `json:"Key"`
	Record *Account
}*/
/*type Collection struct {
	 Composite  map[string][]int
	 Middle map[string][]int
	 //TransactionCount int
}

var(
	composite  =map[string][]int{}
	middle =map[string][]int{}
	transactionCount =0
)*/
func (t *SmartContract) RegenesisTest(ctx contractapi.TransactionContextInterface, from_key string) (map[string][]string, error) {

	// Retrieve the name of the variable to prune
	name := from_key

	// Get all delta rows for the variable
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("f~m~t~v", []string{name})
	if deltaErr != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return nil, fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set computing final value while iterating and deleting each key
	pending := make(map[string][]string)

	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return nil, fmt.Errorf(nextErr.Error())
		}

		// Split the key into its composite parts
		_, keyParts, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)

		if splitKeyErr != nil {
			return nil, fmt.Errorf(splitKeyErr.Error())
		}

		// Retrieve the operation and value
		oringin := keyParts[0]
		middler := keyParts[1]
		destination := keyParts[2]
		amount := keyParts[3]
		temp := destination + ":" + amount
		key := middler + oringin
		if v, ok := pending[key]; ok {
			fmt.Printf("key:%s  value:%v\n", key, v)
			pending[key] = append(pending[key], temp)
		} else {
			fmt.Printf("key not found\n")
			pending[key] = append(pending[key], temp)
		}

		//fmt.Printf("i=%d,orin=%v,cipherTextStr:%s\n",i,oringin,cipherTextStr)
		//val,_:=strconv.ParseFloat(oringin,64)
		//values=append(values,val)

		// Delete the row from the ledger
		deltaRowDelErr := ctx.GetStub().DelState(responseRange.Key)
		if deltaRowDelErr != nil {
			return nil, fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
		}

	}
	for k, v := range pending {
		fmt.Println(k, v)
	}
	return pending, nil
	//fmt.Println()
	//fmt.Printf("Values     : %2f %2f %2f %2f...\n", round(values[0]), round(values[1]), round(values[2]), round(values[3]))
	//fmt.Println()
	// Update the ledger with the final value
	//fmt.Printf("length=%d\n", len(values))
}
func (t *SmartContract) PruneSISDTest(ctx contractapi.TransactionContextInterface, middle_key string) (string, error) {

	// Retrieve the name of the variable to prune
	name := middle_key

	// Get all delta rows for the variable
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("m~f~t~v", []string{name})
	if deltaErr != nil {
		return "", fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return "", fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set computing final value while iterating and deleting each key
	ciphertext := enc_value(0.0)

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

		// Retrieve the operation and value
		cipherTextStr := keyParts[3]
		curCiphertext := StringToCiphertextSingle(cipherTextStr)
		fmt.Printf("i=%d,cipherTextStr:%s\n", i, cipherTextStr)
		evaluator.Add(ciphertext, curCiphertext, ciphertext)

		if err := evaluator.Rescale(ciphertext, params.DefaultScale(), ciphertext); err != nil {
			return "", err
		}
		r := dec_val(ciphertext)
		fmt.Printf("curValue: %.6f \n", r[0])

		// Convert the value to a int
		//value, convErr := strconv.Atoi(valueStr)
		//if convErr != nil {
		//	return fmt.Errorf(convErr.Error())
		//}

		// Delete the row from the ledger
		deltaRowDelErr := ctx.GetStub().DelState(responseRange.Key)
		if deltaRowDelErr != nil {
			return "", fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
		}

	}
	// Update the ledger with the final value
	res := dec_val(ciphertext)
	fmt.Printf("finalValue: %.6f \n", res[0])
	// Update the ledger with the final value
	finalvalueStr := strconv.FormatFloat(res[0], 'f', 20, 64)
	updateRespErr := ctx.GetStub().PutState(name, []byte(finalvalueStr))
	if updateRespErr != nil {
		return "", fmt.Errorf("Could not update the final value of the variable after pruning")
	}

	//fmt.Sprintf("Successfully pruned variable %s, final value is %f, %d rows pruned", name, finalVal, i)

	return finalvalueStr, nil
}
func (t *SmartContract) PruneSIMDTest(ctx contractapi.TransactionContextInterface, middle_key string) (string, error) {

	// Retrieve the name of the variable to prune
	name := middle_key

	// Get all delta rows for the variable
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("m~f~t~v", []string{name})
	if deltaErr != nil {
		return "", fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return "", fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set computing final value while iterating and deleting each key
	values := make([]float64, 0)

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

		// Retrieve the operation and value
		cipherTextStr := keyParts[3]
		oringin := StringToCiphertext(cipherTextStr)
		fmt.Printf("i=%d,orin=%v,cipherTextStr:%s\n", i, oringin, cipherTextStr)
		val, _ := strconv.ParseFloat(oringin, 64)
		values = append(values, val)

		// Delete the row from the ledger
		deltaRowDelErr := ctx.GetStub().DelState(responseRange.Key)
		if deltaRowDelErr != nil {
			return "", fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
		}

	}
	fmt.Println()
	fmt.Printf("Values     : %2f %2f %2f %2f...\n", round(values[0]), round(values[1]), round(values[2]), round(values[3]))
	fmt.Println()
	// Update the ledger with the final value
	fmt.Printf("length=%d\n", len(values))
	ciphertext := enc(values)
	//res:=dec_val(ciphertext)
	//fmt.Println()
	//fmt.Printf("sum Values     : %2f %2f %2f %2f...\n", round(res[0]), round(res[1]), round(res[2]), round(res[3]))
	//fmt.Println()

	//fmt.Sprintf("Successfully pruned variable %s, final value is %f, %d rows pruned", name, finalVal, i)
	batch := 1
	n := len(values)
	rotationsKey := kgen.GenRotationKeysForRotations(params.RotationsForInnerSum(batch, n), false, sk)
	eval := evaluator.WithKey(rlwe.EvaluationKey{Rlk: rlk, Rtks: rotationsKey})
	eval.InnerSum(ciphertext, batch, n, ciphertext)
	r1 := dec_val(ciphertext)
	fmt.Println()
	fmt.Printf("sum Values     : %2f %2f %2f %2f...\n", round(r1[0]), round(r1[1]), round(r1[2]), round(r1[3]))
	fmt.Println()
	// Update the ledger with the final value
	finalvalueStr := strconv.FormatFloat(r1[0], 'f', 20, 64)
	updateRespErr := ctx.GetStub().PutState(name, []byte(finalvalueStr))
	if updateRespErr != nil {
		return "", fmt.Errorf("Could not update the final value of the variable after pruning")
	}
	fmt.Printf("finalvalueStr Values     : %v ...\n", finalvalueStr)

	return finalvalueStr, nil
}
func (t *SmartContract) CollectorTest(ctx contractapi.TransactionContextInterface, accountFrom_key, middle_key, accountTo_key string, amountStr string) error {

	ciphertextStr := StringToCiphertext(amountStr)
	//tmp :=dec_val(ciphertext)
	amount, _ := strconv.ParseFloat(ciphertextStr, 64)
	fmt.Printf(" %.6f~ \n", amount)
	accountFrom, err := t.QueryAccount(ctx, accountFrom_key)
	if err != nil {
		return fmt.Errorf("query account error:%s\n", err)
	}

	/* middle, err := t.QueryAccount(ctx, middle_key)
	if err != nil {
		return err
	}
	*/

	if accountFrom.SavingBalance < amount {
		return fmt.Errorf("Balance is not enough!\n")
	}

	typeName := "M~F~V"
	// Create the composite key
	compositeMKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{middle_key, accountFrom_key, strconv.FormatFloat(amount, 'g', 20, 64)})
	//compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{middle_key,accountFrom_key})
	fmt.Printf("compositeMKey:%v\n", compositeMKey)
	if compositeErr != nil {
		return fmt.Errorf("Could not create a composite key for %s: %s\n", typeName, compositeErr.Error())
	}
	// Save the composite key index
	compositePutErr := ctx.GetStub().PutState(compositeMKey, []byte{0x00})
	if compositePutErr != nil {
		return fmt.Errorf(fmt.Sprintf("Could not put operation for %s in the ledger: %s", compositeMKey, compositePutErr.Error()))
	}

	compositeKey := middle_key + accountFrom_key
	//compositeKey:=middle_key+"~"+accountFrom_key

	accountFrom.SavingBalance -= amount

	accountFrom1AsBytes, _ := json.Marshal(accountFrom)
	err = ctx.GetStub().PutState(accountFrom_key, accountFrom1AsBytes)
	if err != nil {
		return err
	}
	//fmt.Printf("账户 %s余额%d 向中间账户 %s  转账%f\n", accountFrom_key,accountFrom.SavingBalance,middle_key, amount)

	info := accountTo_key + ":" + strconv.FormatFloat(amount, 'f', 2, 64)

	//check
	ledgerBytes, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	//fmt.Printf("compositeKey:%v -- %s --info: %v\n",compositeKey,compositeKey,info)//compositeKey:M~FromMIDDLE2ACCOUNT8 -- M~FromMIDDLE2ACCOUNT8 --info: ACCOUNT9:93
	if ledgerBytes == nil {
		//return fmt.Errorf("compositeKey:%s does not exist", compositeKey)
		fmt.Println("compositeKey:%s does not exist,creating...", compositeKey)
		err = ctx.GetStub().PutState(compositeKey, []byte(info))
		return nil

	}
	fmt.Printf("already ledger: %v \n", string(ledgerBytes))
	//infos=append(infos,string(ledgerBytes))
	//infos=append(infos,info)
	info = string(ledgerBytes) + "," + info
	//fmt.Printf("ledgerBytes: %v--%v \n",ledgerBytes,string(ledgerBytes))//ledgerBytes: [65 67 67 79 85 78 84 56 58 49 51 46 48 48]--ACCOUNT8:13.00
	//compositePutErr :=ctx.GetStub().PutState(compositeKey, []byte{0x00})
	compositePutErr = ctx.GetStub().PutState(compositeKey, []byte(info))
	if compositePutErr != nil {
		return fmt.Errorf("Could not put operation for %s in the ledger: %s", middle_key, compositePutErr.Error())
	}
	fmt.Printf("compositeKey:%v --info: %v\n", compositeKey, info)
	//fmt.Printf("info: %v \n",info)

	//fmt.Sprintf("Successfully %s added %s%v to %s", accountFrom_key,"+", amount, middle_key)

	return nil
}
func (t *SmartContract) Collector(ctx contractapi.TransactionContextInterface, accountFrom_key, middle_key, accountTo_key string, amountCiphertextStr string) (string, error) {

	ciphertextOrig := StringToCiphertext(amountCiphertextStr)
	//tmp :=dec_val(ciphertext)
	amount, _ := strconv.ParseFloat(ciphertextOrig, 64)
	fmt.Printf("orig amount=%.6f\n", amount)

	accountFrom, err := t.QueryAccount(ctx, accountFrom_key)
	if err != nil {
		return "", err
	}
	if accountFrom.SavingBalance < amount {
		return "", fmt.Errorf("Balance is not enough!\n")
	}

	senderTypeName := "f~m~t~v"
	middlerTypeName := "m~f~t~v"
	//value:=strconv.FormatFloat(amount,'f',2,64)

	// Create the composite key base mainkey sender
	senderCompositeKey, scompositeErr := ctx.GetStub().CreateCompositeKey(senderTypeName, []string{accountFrom_key, middle_key, accountTo_key, amountCiphertextStr})
	if scompositeErr != nil {
		return "", fmt.Errorf("Could not create a composite key for %s: %s\n", senderTypeName, scompositeErr.Error())
	}

	//fmt.Printf("senderCompositeKey:%v\n",senderCompositeKey)

	// Create the composite key base mainkey middler
	middlerCompositeKey, mcompositeErr := ctx.GetStub().CreateCompositeKey(middlerTypeName, []string{middle_key, accountFrom_key, accountTo_key, amountCiphertextStr})
	if mcompositeErr != nil {
		return "", fmt.Errorf("Could not create a composite key for %s: %s\n", middlerTypeName, mcompositeErr.Error())
	}

	//fmt.Printf("middlerCompositeKey:%v\n",middlerCompositeKey)
	//fmt.Printf("账户 %s余额%d 向中间账户 %s 余额%d 转账%d\n", accountFrom_key,accountFrom.SavingBalance,middle_key, middle.SavingBalance,amount)
	accountFrom.SavingBalance -= amount

	accountFrom1AsBytes, _ := json.Marshal(accountFrom)
	err = ctx.GetStub().PutState(accountFrom_key, accountFrom1AsBytes)
	if err != nil {
		return "", err
	}

	senderCompositePutErr := ctx.GetStub().PutState(senderCompositeKey, []byte{0x00})
	if senderCompositePutErr != nil {
		return "", senderCompositePutErr
	}
	middlerCompositePutErr := ctx.GetStub().PutState(middlerCompositeKey, []byte{0x00})
	if middlerCompositePutErr != nil {
		return "", middlerCompositePutErr
	}

	//fmt.Printf("Successfully %s added %s%s to %s", accountFrom_key,"+", value, middle_key)
	receipt := accountFrom_key + middle_key + accountTo_key + ciphertextOrig
	return receipt, nil
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
		SavingBalance := 1000.0

		Info := ""
		Status := -1
		account := &Account{Name, SavingBalance, Info, Status}
		accountJSONasBytes, err := json.Marshal(account)
		if err != nil {
			return err
		}
		/*accountJSONasString := `{"name": "` + Name + `", "savingbalance": "` + strconv.Itoa(SavingBalance)+ `", "checkingbalance": ` + strconv.Itoa(CheckingBalance) + `, "status": "` + strconv.Itoa(Status) + `"}`
		accountJSONasBytes := []byte(accountJSONasString)*/

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
	fmt.Printf("create %d Middle\n", number)
	// Initialize the chaincode
	// Write the state to the ledger
	for i := 0; i <= number; {
		Name := "middle" + strconv.Itoa(i)
		middleAsBytes, err := ctx.GetStub().GetState(Name)
		if err != nil {
			return fmt.Errorf("Failed to get middle: " + err.Error())
		} else if middleAsBytes != nil {
			fmt.Println("This middle already exists: " + Name)
			return fmt.Errorf("This middle already exists: " + Name)
		}
		SavingBalance := 0.0

		Info := ""
		Status := -1
		account := &Account{Name, SavingBalance, Info, Status}
		middleJSONasBytes, err := json.Marshal(account)
		if err != nil {
			return err
		}
		//middleJSONasString := `{"name": "` + Name + `", "savingbalance": "` + strconv.Itoa(SavingBalance)+ `", "checkingbalance": ` + strconv.Itoa(CheckingBalance) + `, "status": "` + strconv.Itoa(Status) + `"}`
		//middleJSONasBytes := []byte(middleJSONasString)

		err = ctx.GetStub().PutState("MIDDLE"+strconv.Itoa(i), middleJSONasBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		} else {
			i += 1
		}
	}

	return nil
}

func (t *SmartContract) Prune(ctx contractapi.TransactionContextInterface, middle_key string) (string, error) {

	// Retrieve the name of the variable to prune
	name := middle_key
	fmt.Printf("middler:%s\n", name)

	// Get all delta rows for the variable
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("M~F~V", []string{name})
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

	ciphertext := enc_value(0.0)
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return "", fmt.Errorf(nextErr.Error())
		}

		// Split the key into its composite parts
		_, keyParts, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		fmt.Printf("responseRange.Key:%v\n", responseRange.Key)

		if splitKeyErr != nil {
			return "", fmt.Errorf(splitKeyErr.Error())
		}

		valueStr := keyParts[2]
		fmt.Print(valueStr)
		value, _ := strconv.ParseFloat(valueStr, 64)
		curciphertext := enc_value(value) //init

		evaluator.Add(ciphertext, curciphertext, ciphertext)
		if err := evaluator.Rescale(ciphertext, params.DefaultScale(), ciphertext); err != nil {
			return "", err
		}

		//Convert the value to a int
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
	res := dec_val(ciphertext)
	fmt.Printf("finalValue: %.6f \n", res[0])
	// Update the ledger with the final value
	finalvalueStr := strconv.FormatFloat(res[0], 'f', 20, 64)
	updateRespErr := ctx.GetStub().PutState(name, []byte(finalvalueStr))
	if updateRespErr != nil {
		return "", fmt.Errorf("Could not update the final value of the variable after pruning")
	}

	//fmt.Printf("Successfully pruned variable %s, final value is %f, %d rows pruned", name, finalVal, i)

	return "finalvalueStr", nil
}
func (t *SmartContract) update(ctx contractapi.TransactionContextInterface, args []string) error {
	// Check we have a valid number of args
	if len(args) != 3 {
		return fmt.Errorf("Incorrect number of arguments, expecting 3")
	}

	// Extract the args
	name := args[0]
	op := args[2]
	_, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("Provided value was not a number")
	}

	// Make sure a valid operator is provided
	if op != "+" && op != "-" {
		return fmt.Errorf(fmt.Sprintf("Operator %s is unrecognized", op))
	}

	// Retrieve info needed for the update procedure
	txid := ctx.GetStub().GetTxID()
	compositeIndexName := "varName~op~value~txID"

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey(compositeIndexName, []string{name, op, args[1], txid})
	if compositeErr != nil {
		return fmt.Errorf(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := ctx.GetStub().PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return fmt.Errorf(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr.Error()))
	}
	fmt.Sprintf("Successfully added %s%s to %s", op, args[1], name)
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
func (t *SmartContract) DeleteCompositeKey(ctx contractapi.TransactionContextInterface, args []string) error {
	// Check we have a valid number of args
	if len(args) != 1 {
		return fmt.Errorf("Incorrect number of arguments, expecting 3")
	}

	// Retrieve the variable name
	name := args[0]

	// Delete all delta rows
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("middleName~op~value~account", []string{name})
	if deltaErr != nil {
		return fmt.Errorf(fmt.Sprintf("Could not retrieve delta rows for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Ensure the variable exists
	if !deltaResultsIterator.HasNext() {
		return fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set and delete all indices
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return fmt.Errorf(fmt.Sprintf("Could not retrieve next delta row: %s", nextErr.Error()))
		}

		deltaRowDelErr := ctx.GetStub().DelState(responseRange.Key)
		if deltaRowDelErr != nil {
			return fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
		}
	}

	fmt.Sprintf("Deleted %s, %d rows removed", name, i)

	return nil
}
func (t *SmartContract) Regenesis(ctx contractapi.TransactionContextInterface, from_key, middle_key string) ([]string, error) {

	//pendingtransaction:=[]rlwe.Ciphertext{}
	startKey := from_key + middle_key + "ACCOUNT0000"
	endKey := from_key + middle_key + "ACCOUNT9999"
	//var err error
	// Get all deltas for the variable
	//fmt.Printf("middle_key :%s\n",middle_key)//middle_key :MIDDLE2
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByRange(startKey, endKey)
	fmt.Printf("startKey :%s  endKey:%s\n", startKey, endKey) //startKey :MIDDLE0ACCOUNT000  endKey:MIDDLE0ACCOUNT999
	//fmt.Printf(" QueryCompositeKey deltaResultsIterator--->%v\n",deltaResultsIterator)
	if deltaErr != nil {
		return nil, fmt.Errorf("Could not retrieve value for %s: %s", middle_key, deltaErr.Error())
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return nil, fmt.Errorf("No variable by the name %s exists", middle_key)
	}

	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		fmt.Printf("responseRange.Key:%v\n", responseRange.Key) //MIDDLE0ACCOUNT6
		//ciphertext:=new(rlwe.Ciphertext)
		//json.Unmarshal(responseRange.Value,ciphertext)

		//fmt.Printf("ciphertext :%v\n",ciphertext)//ACCOUNT8:5.41,ACCOUNT7:54.97

		if nextErr != nil {
			return nil, fmt.Errorf(nextErr.Error())
		}
		/*compositevalbytes, err := ctx.GetStub().GetState(responseRange.Key)
		if err != nil {
			return nil,fmt.Errorf("Failed to get compositevalbytes state")
		}
		if compositevalbytes == nil {
			return nil,fmt.Errorf("compositekey Entity not found")
		}*/
		//fmt.Printf("QueryCompositeKey==ledgerBytes: %v--%v \n",compositevalbytes,string(compositevalbytes))//QueryCompositeKey==ledgerBytes: [65 67 67 79 85 78 84 51 58 53 57 46 51 53]--ACCOUNT3:59.35

		// Split the composite key into its component parts
		//_, _, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)

		//fmt.Printf("SplitCompositeKey str:%s\n",str)//middleName~op~value~account

		//fmt.Printf("responseRange.Key :%s\n",responseRange.Key)//responseRange.Key :M~FromMIDDLE2ACCOUNT0
		//fmt.Printf(" keyParts :%s\n",keyParts)//
		//if splitKeyErr != nil {
		//	return nil,fmt.Errorf(splitKeyErr.Error())
		//}

		//compositeKey[i]=responseRange.Key
		//pendingtransaction=append(pendingtransaction,ciphertext)

	}

	// Iterate through result set computing final value while iterating and deleting each key
	fmt.Printf("Regenesis")

	//fmt.Sprintf("Successfully pruned variable %s, final value is %f, %d rows pruned", name, finalVal, i)

	return nil, nil
}

func (t *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, accountFrom_key, middle_key, accountTo_key string, amount float64) (string, error) {

	accountFrom, err := t.QueryAccount(ctx, accountFrom_key)
	if err != nil {
		return "", err
	}
	if accountFrom.SavingBalance < amount {
		return "", fmt.Errorf("Balance is not enough!\n")
	}

	/* middle, err := t.QueryAccount(ctx, middle_key)
	if err != nil {
		return err
	}
	*/
	//typeName := "f~m~t"
	value := strconv.FormatFloat(amount, 'f', 2, 64)

	// Create the composite key
	//compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{accountFrom_key,middle_key,accountTo_key})
	//if compositeErr != nil {
	//	return "",fmt.Errorf("Could not create a composite key for %s: %s\n", typeName, compositeErr.Error())
	//}
	compositeKey := accountFrom_key + middle_key + accountTo_key
	fmt.Printf("compositeKey:%v\n", compositeKey)

	//fmt.Printf("账户 %s余额%d 向中间账户 %s 余额%d 转账%d\n", accountFrom_key,accountFrom.SavingBalance,middle_key, middle.SavingBalance,amount)
	accountFrom.SavingBalance -= amount
	//accountFrom.Info+="{转出"+strconv.Itoa(amount)+middle_key+"}"
	//middle.SavingBalance += amount
	//middle.Info+="{"+accountFrom_key+"转进"+strconv.Itoa(amount)+"}"
	//middle.Status=0 // -1 menas 初始状态， 0 means已锁定，1 menas 已转移

	accountFrom1AsBytes, _ := json.Marshal(accountFrom)
	err = ctx.GetStub().PutState(accountFrom_key, accountFrom1AsBytes)
	if err != nil {
		return "", err
	}
	/*middleAsBytes, _ := json.Marshal(middle)
	err=ctx.GetStub().PutState(middle_key, middleAsBytes)
	if err != nil {
		return err
	}*/

	//compositePutErr :=ctx.GetStub().PutState(compositeKey, []byte{0x00})
	//str:=account_key+"+"+strconv.Itoa(amount)
	//fmt.Printf("add compositeKey:%s-->%v\n",compositeKey,[]byte(str))
	ciphertext := enc_value(amount) //init
	cipherStr := CiphertextToString(ciphertext)
	fmt.Printf("cipherStr:%v\n", cipherStr) //cipherStr:90.0145407026

	cipherStrJSONasBytes, _ := json.Marshal(cipherStr)
	//fmt.Printf("transactionJSONasBytes-->%v\n",transactionJSONasBytes)
	//if err != nil {
	//	return err
	//}
	//pend:=new(PendingTransation)
	//_ = json.Unmarshal(transactionJSONasBytes, pend)
	//fmt.Printf("transactionUnJSON-->%v:%v\n",pend.Receiver,pend.Cipertext)
	compositePutErr := ctx.GetStub().PutState(compositeKey, cipherStrJSONasBytes)
	//compositePutErr :=ctx.GetStub().PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return "", fmt.Errorf("Could not put operation for %s in the ledger: %s", middle_key, compositePutErr.Error())
	}

	fmt.Printf("Successfully %s added %s%s to %s", accountFrom_key, "+", value, middle_key)
	receipt := compositeKey + ":" + value
	return receipt, nil
}
func (t *SmartContract) TransferUseDivCompositeKey(ctx contractapi.TransactionContextInterface, accountFrom_key, middle_key, accountTo_key string, amount float64) error {

	//tmp :=dec_val(ciphertext)
	fmt.Printf(" %.6f~ \n", amount)
	accountFrom, err := t.QueryAccount(ctx, accountFrom_key)
	if err != nil {
		return fmt.Errorf("query account error:%s\n", err)
	}

	/* middle, err := t.QueryAccount(ctx, middle_key)
	if err != nil {
		return err
	}
	*/

	if accountFrom.SavingBalance < amount {
		return fmt.Errorf("Balance is not enough!\n")
	}

	//fmt.Printf("账户 %s余额%d 向中间账户 %s 余额%d 转账%d\n", accountFrom_key,accountFrom.SavingBalance,middle_key, middle.SavingBalance,amount)

	//fmt.Printf("accountFrom.SavingBalance= %f  amount%v\n", accountFrom.SavingBalance,amount)//accountFrom.SavingBalance= 1000.000000  amount95.40112327802981

	typeName := "M~F~V"
	// Create the composite key
	compositeMKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{middle_key, accountFrom_key, strconv.FormatFloat(amount, 'g', 20, 64)})
	//compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{middle_key,accountFrom_key})
	fmt.Printf("compositeMKey:%v\n", compositeMKey)
	if compositeErr != nil {
		return fmt.Errorf("Could not create a composite key for %s: %s\n", typeName, compositeErr.Error())
	}
	// Save the composite key index
	compositePutErr := ctx.GetStub().PutState(compositeMKey, []byte{0x00})
	if compositePutErr != nil {
		return fmt.Errorf(fmt.Sprintf("Could not put operation for %s in the ledger: %s", compositeMKey, compositePutErr.Error()))
	}

	compositeKey := middle_key + accountFrom_key
	//compositeKey:=middle_key+"~"+accountFrom_key

	accountFrom.SavingBalance -= amount

	accountFrom1AsBytes, _ := json.Marshal(accountFrom)
	err = ctx.GetStub().PutState(accountFrom_key, accountFrom1AsBytes)
	if err != nil {
		return err
	}
	//fmt.Printf("账户 %s余额%d 向中间账户 %s  转账%f\n", accountFrom_key,accountFrom.SavingBalance,middle_key, amount)

	info := accountTo_key + ":" + strconv.FormatFloat(amount, 'f', 2, 64)

	//check
	ledgerBytes, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	//fmt.Printf("compositeKey:%v -- %s --info: %v\n",compositeKey,compositeKey,info)//compositeKey:M~FromMIDDLE2ACCOUNT8 -- M~FromMIDDLE2ACCOUNT8 --info: ACCOUNT9:93
	if ledgerBytes == nil {
		//return fmt.Errorf("compositeKey:%s does not exist", compositeKey)
		fmt.Println("compositeKey:%s does not exist,creating...", compositeKey)
		err = ctx.GetStub().PutState(compositeKey, []byte(info))
		return nil

	}
	fmt.Printf("already ledger: %v \n", string(ledgerBytes))
	//infos=append(infos,string(ledgerBytes))
	//infos=append(infos,info)
	info = string(ledgerBytes) + "," + info
	//fmt.Printf("ledgerBytes: %v--%v \n",ledgerBytes,string(ledgerBytes))//ledgerBytes: [65 67 67 79 85 78 84 56 58 49 51 46 48 48]--ACCOUNT8:13.00
	//compositePutErr :=ctx.GetStub().PutState(compositeKey, []byte{0x00})
	compositePutErr = ctx.GetStub().PutState(compositeKey, []byte(info))
	if compositePutErr != nil {
		return fmt.Errorf("Could not put operation for %s in the ledger: %s", middle_key, compositePutErr.Error())
	}
	fmt.Printf("compositeKey:%v --info: %v\n", compositeKey, info)
	//fmt.Printf("info: %v \n",info)

	//fmt.Sprintf("Successfully %s added %s%v to %s", accountFrom_key,"+", amount, middle_key)

	return nil
}

////func (t *SmartContract) GetDivCompositeKey(ctx contractapi.TransactionContextInterface) (map[string][]int,error) {
//func (t *SmartContract) GetDivCompositeKey(ctx contractapi.TransactionContextInterface,middle_key string) (Collection,error) {
//
//
//
//
//	middle, err := t.QueryAccount(ctx, middle_key)
//	//fmt.Printf("middle:%v\n",middle)
//	if err != nil {
//		return nil,err
//	}
//
//	//var err error
//	// Get all deltas for the variable
//	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("middleName~op~value~account", []string{middle_key})
//	fmt.Printf("deltaResultsIterator--->%v\n",deltaResultsIterator)
//	if deltaErr != nil {
//		return nil,fmt.Errorf("Could not retrieve value for %s: %s", middle_key, deltaErr.Error())
//	}
//	defer deltaResultsIterator.Close()
//
//	// Check the variable existed
//	if !deltaResultsIterator.HasNext() {
//		return nil,fmt.Errorf("No variable by the name %s exists", middle_key)
//	}
//	// Iterate through result set and compute final value
//	var finalVal float64
//	var i int
//	for i = 0; deltaResultsIterator.HasNext(); i++ {
//		// Get the next row
//		responseRange, nextErr := deltaResultsIterator.Next()
//		if nextErr != nil {
//			return nil,fmt.Errorf(nextErr.Error())
//		}
//
//		// Split the composite key into its component parts
//		_, keyParts, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)
//		startKey := responseRange.Key
//		const maxUnicodeRuneValue int32 = utf8.MaxRune
//		endKey := responseRange.Key + string(maxUnicodeRuneValue)
//		fmt.Printf("[startKey:%s  -- endKey:%s]\n",startKey,endKey)
//		//fmt.Printf("SplitCompositeKey str:%s\n",str)//middleName~op~value~account
//		fmt.Printf("responseRange.Key :%s\n",responseRange.Key)//
//		fmt.Printf("keyParts :%s\n",keyParts)//
//		if splitKeyErr != nil {
//			return nil,fmt.Errorf(splitKeyErr.Error())
//		}
//
//		// Retrieve the delta value and operation
//		operation := keyParts[1]
//		valueStr := keyParts[2]
//		from_account :=keyParts[3]
//		fmt.Printf("operation:%s  valueStr:%s  --> %s\n",operation,valueStr,from_account)//operation:+  valueStr:10
//		// Convert the value string and perform the operation
//		value, convErr := strconv.ParseFloat(valueStr,64)
//		if convErr != nil {
//			return nil,fmt.Errorf(convErr.Error())
//		}
//
//		switch operation {
//		case "+":
//			finalVal += value
//		case "-":
//			finalVal -= value
//		default:
//			return nil,fmt.Errorf("Unrecognized operation %s", operation)
//		}
//		fmt.Printf("i=%d,finalval=%d\n",i,finalVal)
//	}
//
//	middle.SavingBalance = finalVal
//	middleAsBytes, _ := json.Marshal(middle)
//	err=ctx.GetStub().PutState(middle_key, middleAsBytes)
//	return middle,nil
//
//
//
//	return collection,nil
//}

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
//
//	func (t *SmartContract) QueryAllAccounts(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
//		startKey := ""
//		endKey := ""
//
//		resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
//
//		if err != nil {
//			return nil, err
//		}
//		defer resultsIterator.Close()
//
//		results := []QueryResult{}
//
//		for resultsIterator.HasNext() {
//			queryResponse, err := resultsIterator.Next()
//
//			if err != nil {
//				return nil, err
//			}
//
//			account := new(Account)
//			_ = json.Unmarshal(queryResponse.Value, account)
//
//			queryResult := QueryResult{Key: queryResponse.Key, Record: account}
//			results = append(results, queryResult)
//		}
//
//		return results, nil
//	}
//
// Query callback representing the query of a chaincode
func (t *SmartContract) QueryAccount(ctx contractapi.TransactionContextInterface, accountId string) (*Account, error) {
	//func (t *SmartContract) QueryAccount(ctx contractapi.TransactionContextInterface, accountId string) (map[string][]int, error) {
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

// Query callback representing the query of a chaincode
func (t *SmartContract) QueryValue(ctx contractapi.TransactionContextInterface, Id string) (string, error) {
	var err error
	fmt.Println(" Id:%v\n", Id)
	// Get the state from the ledger
	valbytes, err := ctx.GetStub().GetState(Id)
	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if valbytes == nil {
		return "", fmt.Errorf("%s does not exist", Id)
	}
	fmt.Printf("%v\n", string(valbytes))

	return string(valbytes), nil
}

// Query callback representing the query of a chaincode
//func (t *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, middle_key,account_key string) (*Transaction, error) {
//	var err error
//	// Get the state from the ledger
//	typeName := "middleName~op~value~account"
//	value :=strconv.Itoa(10)
//	//compositeKey := typeName+middle_key+"+"+value+account_key  //不能这么用
//	compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey(typeName, []string{middle_key, "+",value,account_key})
//	fmt.Printf("Transaction compositeKey:%v\n",compositeKey)
//	if compositeErr != nil {
//		return nil,fmt.Errorf("Could not create a composite key for %s: %s\n", typeName, compositeErr.Error())
//	}
//	valbytes, err := ctx.GetStub().GetState(compositeKey)
//
//	if err != nil {
//		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
//	}
//
//	if valbytes == nil {
//		return nil, fmt.Errorf("%s does not exist", compositeKey)
//	}
//	jsonResp := "{\"Name\":\"" + compositeKey + "\",\"Transaction\":\"" + string(valbytes) + "\"}"
//	fmt.Printf("Query Response:%s\n", jsonResp)
//	transaction := new(Transaction)
//	_ = json.Unmarshal(valbytes, transaction)
//	return transaction, nil
//}
// Query callback representing the query of a chaincode
/**
Result:[{"Key":"\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000012\u0000ACCOUNT1\u0000","Record":{"receiver":"MIDDLE2","amount":12,"transfer":"ACCOUNT28"}},{"Key":"\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000013\u0000ACCOUNT1\u0000","Record":{"receiver":"MIDDLE2","amount":13,"transfer":"ACCOUNT28"}},{"Key":"\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000015\u0000ACCOUNT1\u0000","Record":{"receiver":"MIDDLE2","amount":15,"transfer":"ACCOUNT15"}},{"Key":"\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000019\u0000ACCOUNT1\u0000","Record":{"receiver":"MIDDLE2","amount":19,"transfer":"ACCOUNT45"}},{"Key":"\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u00002\u0000ACCOUNT1\u0000","Record":{"receiver":"MIDDLE2","amount":2,"transfer":"ACCOUNT16"}},{"Key":"\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000022\u0000ACCOUNT1\u0000","Record":{"receiver":"MIDDLE2","amount":22,"transfer":"ACCOUNT12"}},{"Key":"\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000024\u0000ACCOUNT1\u0000","Record":{"receiver":"MIDDLE2","amount":24,"transfer":"ACCOUNT33"}},{"Key":"\u0000middleName~op~value~account\u0000MIDDLE2\u0000+\u000025\u0000ACCOUNT1\u0000","Record":{"receiver":"MIDDLE2","amount":25,"transfer":"ACCOUNT0"}}]
*/
//func (t *SmartContract) QueryTransactionSet(ctx contractapi.TransactionContextInterface, middle_key string) ([]QueryTransactionResult, error) {
//	//var err error
//	// Get all deltas for the variable
//
//	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("middleName~op~value~account", []string{middle_key})
//	if deltaErr != nil {
//		return nil,fmt.Errorf("Could not retrieve value for %s: %s", middle_key, deltaErr.Error())
//	}
//	defer deltaResultsIterator.Close()
//
//
//	results := []QueryTransactionResult{}
//	// Check the variable existed
//	if !deltaResultsIterator.HasNext() {
//		return nil,fmt.Errorf("No variable by the name %s exists", middle_key)
//	}
//	// Iterate through result set and compute final value
//	//var finalVal int
//	var i int
//	for i = 0; deltaResultsIterator.HasNext(); i++ {
//		// Get the next row
//		queryResponse, nextErr := deltaResultsIterator.Next()
//		if nextErr != nil {
//			return nil,fmt.Errorf(nextErr.Error())
//		}
//
//		transaction := new(Transaction)
//		_ = json.Unmarshal(queryResponse.Value, transaction)
//		//fmt.Printf("queryResponse.Value===>%v,transaction:%s",queryResponse.Value,transaction)
//		queryResult := QueryTransactionResult{Key: queryResponse.Key, Record: transaction}
//		results = append(results, queryResult)
//
//	}
//	//fmt.Printf("finalval=%d\n",finalVal)
//	fmt.Printf("QueryTransactionResult=%d\n",results)
//
//
//	return results,nil
//
//}
// Query callback representing the query of a chaincode
//func (t *SmartContract) QueryMiddleByCompositeKey(ctx contractapi.TransactionContextInterface, middle_key string) (*Account, error) {
//	middle, err := t.QueryAccount(ctx, middle_key)
//	//fmt.Printf("middle:%v\n",middle)
//	if err != nil {
//		return nil,err
//	}
//
//	//var err error
//	// Get all deltas for the variable
//	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey("middleName~op~value~account", []string{middle_key})
//	fmt.Printf("deltaResultsIterator--->%v\n",deltaResultsIterator)
//	if deltaErr != nil {
//		return nil,fmt.Errorf("Could not retrieve value for %s: %s", middle_key, deltaErr.Error())
//	}
//	defer deltaResultsIterator.Close()
//
//	// Check the variable existed
//	if !deltaResultsIterator.HasNext() {
//		return nil,fmt.Errorf("No variable by the name %s exists", middle_key)
//	}
//	// Iterate through result set and compute final value
//	var finalVal float64
//	var i int
//	for i = 0; deltaResultsIterator.HasNext(); i++ {
//		// Get the next row
//		responseRange, nextErr := deltaResultsIterator.Next()
//		if nextErr != nil {
//			return nil,fmt.Errorf(nextErr.Error())
//		}
//
//		// Split the composite key into its component parts
//		_, keyParts, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)
//		startKey := responseRange.Key
//		const maxUnicodeRuneValue int32 = utf8.MaxRune
//		endKey := responseRange.Key + string(maxUnicodeRuneValue)
//		fmt.Printf("[startKey:%s  -- endKey:%s]\n",startKey,endKey)
//		//fmt.Printf("SplitCompositeKey str:%s\n",str)//middleName~op~value~account
//		fmt.Printf("responseRange.Key :%s\n",responseRange.Key)//
//		fmt.Printf("keyParts :%s\n",keyParts)//
//		if splitKeyErr != nil {
//			return nil,fmt.Errorf(splitKeyErr.Error())
//		}
//
//		// Retrieve the delta value and operation
//		operation := keyParts[1]
//		valueStr := keyParts[2]
//		from_account :=keyParts[3]
//		fmt.Printf("operation:%s  valueStr:%s  --> %s\n",operation,valueStr,from_account)//operation:+  valueStr:10
//		// Convert the value string and perform the operation
//		value, convErr := strconv.ParseFloat(valueStr,64)
//		if convErr != nil {
//			return nil,fmt.Errorf(convErr.Error())
//		}
//
//		switch operation {
//		case "+":
//			finalVal += value
//		case "-":
//			finalVal -= value
//		default:
//			return nil,fmt.Errorf("Unrecognized operation %s", operation)
//		}
//		fmt.Printf("i=%d,finalval=%d\n",i,finalVal)
//	}
//
//	middle.SavingBalance = finalVal
//	middleAsBytes, _ := json.Marshal(middle)
//	err=ctx.GetStub().PutState(middle_key, middleAsBytes)
//	return middle,nil
//
//}

func (t *SmartContract) QueryCompositeKey(ctx contractapi.TransactionContextInterface, middle_key string) ([]string, error) {

	compositeKey := []string{}
	startKey := middle_key + "ACCOUNT0000"
	endKey := middle_key + "ACCOUNT9999"
	//var err error
	// Get all deltas for the variable
	//fmt.Printf("middle_key :%s\n",middle_key)//middle_key :MIDDLE2
	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByRange(startKey, endKey)
	fmt.Printf("startKey :%s  endKey:%s\n", startKey, endKey) //startKey :MIDDLE0ACCOUNT000  endKey:MIDDLE0ACCOUNT999
	//fmt.Printf(" QueryCompositeKey deltaResultsIterator--->%v\n",deltaResultsIterator)
	if deltaErr != nil {
		return nil, fmt.Errorf("Could not retrieve value for %s: %s", middle_key, deltaErr.Error())
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return nil, fmt.Errorf("No variable by the name %s exists", middle_key)
	}

	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		//fmt.Println(responseRange.Key)//MIDDLE0ACCOUNT6
		fmt.Println(string(responseRange.Value)) //ACCOUNT8:5.41,ACCOUNT7:54.97

		if nextErr != nil {
			return nil, fmt.Errorf(nextErr.Error())
		}
		/*compositevalbytes, err := ctx.GetStub().GetState(responseRange.Key)
		if err != nil {
			return nil,fmt.Errorf("Failed to get compositevalbytes state")
		}
		if compositevalbytes == nil {
			return nil,fmt.Errorf("compositekey Entity not found")
		}*/
		//fmt.Printf("QueryCompositeKey==ledgerBytes: %v--%v \n",compositevalbytes,string(compositevalbytes))//QueryCompositeKey==ledgerBytes: [65 67 67 79 85 78 84 51 58 53 57 46 51 53]--ACCOUNT3:59.35

		// Split the composite key into its component parts
		//_, _, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)

		//fmt.Printf("SplitCompositeKey str:%s\n",str)//middleName~op~value~account

		//fmt.Printf("responseRange.Key :%s\n",responseRange.Key)//responseRange.Key :M~FromMIDDLE2ACCOUNT0
		//fmt.Printf(" keyParts :%s\n",keyParts)//
		//if splitKeyErr != nil {
		//	return nil,fmt.Errorf(splitKeyErr.Error())
		//}

		//compositeKey[i]=responseRange.Key
		compositeKey = append(compositeKey, string(responseRange.Value))

	}

	return compositeKey, nil

}

//func (t *SmartContract) QueryRange(ctx contractapi.TransactionContextInterface, middle_key string) ([]string, error) {
//
//    compositeKey:=[]string{}
//    startKey:=middle_key
//    endKey:="ACCOUNT222"
//	//var err error
//	// Get all deltas for the variable
//	//fmt.Printf("middle_key :%s\n",middle_key)//middle_key :MIDDLE2
//	deltaResultsIterator, deltaErr := ctx.GetStub().GetStateByRange(startKey,endKey)
//	fmt.Printf("startKey :%s  endKey:%s\n",startKey,endKey)//startKey :ACCOUNT1  endKey:ACCOUNT222
//	//fmt.Printf(" QueryCompositeKey deltaResultsIterator--->%v\n",deltaResultsIterator)
//	if deltaErr != nil {
//		return nil,fmt.Errorf("Could not retrieve value for %s: %s", middle_key, deltaErr.Error())
//	}
//	defer deltaResultsIterator.Close()
//
//	// Check the variable existed
//	if !deltaResultsIterator.HasNext() {
//		return nil,fmt.Errorf("No variable by the name %s exists", middle_key)
//	}
//
//	var i int
//	for i = 0; deltaResultsIterator.HasNext(); i++ {
//		// Get the next row
//		responseRange, nextErr := deltaResultsIterator.Next()
//		//fmt.Println(responseRange.Key)//ACCOUNT1
//		//fmt.Println(responseRange.Value)//[123 34 105 110 102 111 34 58 34 34 44 34 110 97 109 101 34 58 34 97 99 99 111 117 110 116 49 34 44 34 115 97 118 105 110 103 98 97 108 97 110 99 101 34 58 49 48 48 48 44 34 115 116 97 116 117 115 34 58 45 49 125]
//		//fmt.Println(string(responseRange.Value))//
//		if nextErr != nil {
//			return nil,fmt.Errorf(nextErr.Error())
//		}
//		/*compositevalbytes, err := ctx.GetStub().GetState(responseRange.Key)
//		if err != nil {
//			return nil,fmt.Errorf("Failed to get compositevalbytes state")
//		}
//		if compositevalbytes == nil {
//			return nil,fmt.Errorf("compositekey Entity not found")
//		}*/
//		//fmt.Printf("QueryCompositeKey==ledgerBytes: %v--%v \n",compositevalbytes,string(compositevalbytes))//QueryCompositeKey==ledgerBytes: [65 67 67 79 85 78 84 51 58 53 57 46 51 53]--ACCOUNT3:59.35
//
//
//		// Split the composite key into its component parts
//		//_, _, splitKeyErr := ctx.GetStub().SplitCompositeKey(responseRange.Key)
//
//		//fmt.Printf("SplitCompositeKey str:%s\n",str)//middleName~op~value~account
//
//		//fmt.Printf("responseRange.Key :%s\n",responseRange.Key)//responseRange.Key :M~FromMIDDLE2ACCOUNT0
//		//fmt.Printf(" keyParts :%s\n",keyParts)//
//		//if splitKeyErr != nil {
//		//	return nil,fmt.Errorf(splitKeyErr.Error())
//		//}
//
//
//		//compositeKey[i]=responseRange.Key
//        compositeKey=append(compositeKey,string(responseRange.Value))
//
//
//	}
//
//	return compositeKey,nil
//
//}

//func (s *SmartContract) putSK(ctx contractapi.TransactionContextInterface, id string) (string, error) {
//
//
//
//
//	var sk string
//	err = json.Unmarshal(assetJSON, &asset)
//	if err != nil {
//		return nil, err
//	}
//
//	return &asset, nil
//}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
func (pendingpool PendingPool) MarshalJSON(data []byte, err error) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
func round(x float64) float64 {
	return math.Round(x*100000000) / 100000000
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
