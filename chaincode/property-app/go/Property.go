// 부동산(Property) - Seller - Buyer
// Add Property
// Query a Property By ID or All Properties
// Transfer Property Ownership

// Import Section
package main

// 함수가 호출할 때마다 GO로 작성된 내용의 인풋과 아웃풋이 JSON 형태로 되도록 만들기 위해서 사용된다.

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Property Transfer smart contract to show the property transfer transactions
type PropertyTransferSmartContract struct {
	contractapi.Contract
}

// Property describes basic details
type Property struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Area int `json:"area"`
	OwnerName string `json:"ownerName"`
	Value int `json:"value"`
}

func (pc *PropertyTransferSmartContract) AddProperty(
	ctx contractapi.TransactionContextInterface, 
	id string, 
	name string,
	area int,
	ownerName string,
	value int,
) error {
	propertyJSON, err := ctx.GetStub().GetState(id)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state", err)
	}

	if propertyJSON != nil {
		return fmt.Errorf("the property %s already exists", id)
	}
	
	prop := Property {
		ID: id,
		Name: name, 
		Area: area,
		OwnerName: ownerName, 
		Value: value,
	}

	// JSON 인코딩 
	propertyBytes, err := json.Marshal(prop)

	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, propertyBytes)
}

// Query All Properties 
// This function returns all the existing properties
func (pc *PropertyTransferSmartContract) QueryAllProperties (
		ctx contractapi.TransactionContextInterface,
	) ([]*Property, error) {
	propertyIterator, err := ctx.GetStub().GetStateByRange("","")

	if err != nil {
		return nil, err
	}

	var properties []*Property

	for propertyIterator.HasNext() {
		propertyResponse, err := propertyIterator.Next()

	
		if err != nil {
			return nil, err
		}

		var property *Property
		err = json.Unmarshal(propertyResponse.Value, &property)

		if err != nil {
			return nil, err
		}
		properties = append(properties, property)
	}
	return properties, nil
}

// Query Property By ID
func (pc *PropertyTransferSmartContract) QueryPropertyById(
	ctx contractapi.TransactionContextInterface, id string,
	) (*Property, error) {
	// GetState로 Id전달 해서 PropertyJSON 가져오기
	propertyJSON, err := ctx.GetStub().GetState(id)

	// GetState를 호출할 때 에러가 발생이 되면, 에러 처리
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read the data from the world state", err)
	}
	
	// ID가 존재하지 않을 때, 에러 처리
	if propertyJSON == nil {
		return nil, fmt.Errorf("The property %s does not exist", id)
	}
	// property 변수 설정
	var property *Property // see through Property

	// Unmarshal을 통해서 PropertyJSON 전달
	err = json.Unmarshal(propertyJSON, &property)

	// Unmarshal을 함수 실행할 때 에러가 발생이되면 에러 처리
	if err != nil {
		return nil, err
	}
	// 최종적으로 property 반환
	return property, nil
}

// This function helps to transfer the ownership of the property
func (pc *PropertyTransferSmartContract) TransferProperty(
	ctx contractapi.TransactionContextInterface, 
	id string,
	newOwner string,
	) error {
		property, err := pc.QueryPropertyById(ctx, id)

		if err != nil {
			return err
		}

		property.OwnerName = newOwner
		propertyJSON, err := json.Marshal(property)
		
		if err != nil {
			return err
		}
		return ctx.GetStub().PutState(id, propertyJSON)
}

// Change Property Value Function
func (pc *PropertyTransferSmartContract) ChangePropertyValue(
	ctx contractapi.TransactionContextInterface, 
	id string,
	newValue int,
	) error {
		property, err := pc.QueryPropertyById(ctx, id)

		if err != nil {
			return err
		}

		property.Value = newValue
		propertyJSON, err := json.Marshal(property)
		
		if err != nil {
			return err
		}
		return ctx.GetStub().PutState(id, propertyJSON)
}

func main() {
	propertyTransferSmartContract := new(PropertyTransferSmartContract)

	cc, err := contractapi.NewChaincode(propertyTransferSmartContract)

	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}





