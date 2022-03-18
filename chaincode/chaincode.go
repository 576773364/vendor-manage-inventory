package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/vendor-manage-inventory/chaincode/lib"
	"github.com/vendor-manage-inventory/chaincode/utils"
	"strconv"
)

type MedicalSystem struct {
}

func (t *MedicalSystem) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// 初始化供应商
	supplierBytes := []byte("supplierAdmin")
	err := stub.PutState(lib.KeyOfSupplier, supplierBytes)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
	}
	// 初始化补货方案map
	schemesMap := make(map[string]lib.ReplenishmentScheme)
	schemesMapJSON,err := json.Marshal(schemesMap)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
	}
	err = stub.PutState(lib.KeyOfSchemesMap,schemesMapJSON)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
	}

	return pb.Response{Status: 200, Message: "Initialize successful", Payload: nil}
}

func (t *MedicalSystem) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "retailerRegistration" {
		// 零售商注册账号
		return t.retailerRegistration(stub, args)
	} else if function == "retailerViewScheme" {
		// 零售商查看供应商补货方案
		return t.retailerViewScheme(stub, args)
	} else if function == "retailerResponseScheme" {
		// 零售商回应补货方案
		return t.retailerResponseScheme(stub, args)
	} else if function == "retailerUpdateInventory" {
		// 零售商更新库存
		return t.retailerUpdateInventory(stub, args)
	} else if function == "supplierAuditRegistration" {
		// 供货商通过与拒绝零售商注册
		return t.supplierAuditRegistration(stub, args)
	} else if function == "supplierViewSchemes" {
		// 供货商查看零售商们补货方案
		return t.supplierViewSchemes(stub, args)
	}

	return shim.Error("Invalid invoke function name.")
}

// 零售商注册账号
// 参数： 零售商名称 订货单价 提前期 初始库存 需求量均值 上传数据的周期 库存商品价值 年利率 固定订货成本 审查周期
// 返回： 空
func (t *MedicalSystem) retailerRegistration(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数个数
	if len(args) != 10 {
		return pb.Response{Status: 400, Message: "Incorrect number of arguments. Expecting 10", Payload: nil}
	}
	// 判断参数合法性（每个参数都不能为空）
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" || args[5] == "" || args[6] == "" || args[7] == "" || args[8] == "" || args[9] == "" {
		return pb.Response{Status: 400, Message: "The parameter cannot be empty", Payload: nil}
	}
	// 赋值变量，与函数前参数说明顺序一致。部分参数需转换数据类型（与 Retailer 结构体各属性的数据类型一致）
	retailerName := args[0] // 零售商名称
	// 将 string 转换为 float64，下同
	unitPrice, err := strconv.ParseFloat(args[1], 64) // 订货单价
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	// 将 string 转换为 int，下同
	leadTime, err := strconv.Atoi(args[2]) // 提前期
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	inventory, err := strconv.Atoi(args[3]) // 库存量
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	averageDemand, err := strconv.Atoi(args[4]) // 需求量均值
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	updateCycle, err := strconv.Atoi(args[5]) // 上传数据的周期
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	inventoryValue, err := strconv.ParseFloat(args[6], 64) // 库存商品价值
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	annualInterestRate, err := strconv.ParseFloat(args[7], 64) // 年利率
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	fixedOrderCost, err := strconv.ParseFloat(args[8], 64) // 固定订货成本
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	reviewCycle, err := strconv.Atoi(args[9]) // 审查周期
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}
	// 创建零售商对象
	retailer := lib.Retailer{
		RetailerName:       retailerName,
		UnitPrice:          unitPrice,
		LeadTime:           leadTime,
		Inventory:          inventory,
		AverageDemand:      averageDemand,
		UpdateCycle:        updateCycle,
		InventoryValue:     inventoryValue,
		AnnualInterestRate: annualInterestRate,
		FixedOrderCost:     fixedOrderCost,
		ReviewCycle:        reviewCycle,
		State:              lib.ToBeResponded,
	}
	// 序列化对象
	retailerJSON, err := json.Marshal(retailer)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
	}
	// 写入账本
	err = stub.PutState(retailer.RetailerName, retailerJSON)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
	}



	return pb.Response{Status: 200, Message: "Register successful", Payload: nil}
}

// 零售商查看供应商补货方案
// 参数： 零售商名称
// 返回： 补货方案对象
func (t *MedicalSystem) retailerViewScheme(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return pb.Response{Status: 400, Message: "Incorrect number of arguments. Expecting 1", Payload: nil}
	}
	if args[0] == "" {
		return pb.Response{Status: 400, Message: "The parameter cannot be empty", Payload: nil}
	}
	retailerName := args[0]

	// 读取账本，获取该零售商对象
	retailerJSON, err := stub.GetState(retailerName)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	} else if retailerJSON == nil {
		return pb.Response{Status: 400, Message: "The retailer does not exist", Payload: nil}
	}
	retailer := new(lib.Retailer)
	// 反序列化对象
	err = json.Unmarshal(retailerJSON, retailer)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
	}

	// 验证零售商信息已由供应商审核通过
	if retailer.State != lib.Pass {
		return pb.Response{Status: 400, Message: "The retailer failed the audit", Payload: nil}
	}

	replenishmentSchemeJSON, err := stub.GetState(utils.ConstructSchemeKey(retailerName))
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	}

	return pb.Response{Status: 200, Message: "View successful", Payload: replenishmentSchemeJSON}
}

// 零售商回应补货方案
// 参数： 零售商名称 回应（0或1）
// 返回： 空 或 补货前与补货后库存量
func (t *MedicalSystem) retailerResponseScheme(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数个数
	if len(args) != 2 {
		return pb.Response{Status: 400, Message: "Incorrect number of arguments. Expecting 2", Payload: nil}
	}
	// 判断参数合法性（每个参数都不能为空， 第二个参数必须为 0 或 1）
	if args[0] == "" || args[1] == "" {
		return pb.Response{Status: 400, Message: "The parameter cannot be empty", Payload: nil}
	}
	if args[1] != "0" && args[1] != "1" {
		return pb.Response{Status: 400, Message: "The response result must be 0 or 1", Payload: nil}
	}
	// 赋值变量，与函数前参数说明顺序一致。
	retailerName := args[0]
	result := args[1]

	// 读取账本，获取该零售商对象
	retailerJSON, err := stub.GetState(retailerName)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	} else if retailerJSON == nil {
		return pb.Response{Status: 400, Message: "The retailer does not exist", Payload: nil}
	}
	retailer := new(lib.Retailer)
	// 反序列化对象
	err = json.Unmarshal(retailerJSON, retailer)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
	}

	// 验证零售商信息已由供应商审核通过
	if retailer.State != lib.Pass {
		return pb.Response{Status: 400, Message: "The retailer failed the audit", Payload: nil}
	}

	// 读取账本，获取该零售商的补货方案对象
	replenishmentSchemeJSON, err := stub.GetState(utils.ConstructSchemeKey(retailerName))
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	}
	var replenishmentScheme lib.ReplenishmentScheme
	// 反序列化对象
	err = json.Unmarshal(replenishmentSchemeJSON, &replenishmentScheme)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
	}

	// 更新供应商的补货方案Map
	// 获取补货方案map
	schemesMapJSON,err := stub.GetState(lib.KeyOfSchemesMap)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	}
	// 反序列化补货方案map
	schemesMap := make(map[string]lib.ReplenishmentScheme)
	err = json.Unmarshal(schemesMapJSON, &schemesMap)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
	}

	// 判断回应
	// 如果为0，即为不同意
	if result == "0" {
		// 修改补货方案的回应结果为 不同意
		replenishmentScheme.ResponseResults = lib.Veto
		// 序列化对象
		replenishmentSchemeJSON, err = json.Marshal(replenishmentScheme)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
		}
		// 写入账本
		err = stub.PutState(utils.ConstructSchemeKey(retailerName), replenishmentSchemeJSON)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
		}

		// 更新当前零售商补货方案到map中
		schemesMap[retailerName] = replenishmentScheme
		// 序列化
		schemesMapJSON,err = json.Marshal(schemesMap)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
		}
		// 写入账本
		err = stub.PutState(lib.KeyOfSchemesMap,schemesMapJSON)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
		}

		return pb.Response{Status: 200, Message: "Veto successful", Payload: nil}
	} else { // 如果为1，即为同意

		// 获取旧的库存量
		oldInventory := retailer.Inventory
		// 修改库存量
		retailer.Inventory += replenishmentScheme.ReorderQuantity

		// 修改补货方案的回应结果为 同意
		replenishmentScheme.ResponseResults = lib.Pass
		// 序列化对象
		replenishmentSchemeJSON, err = json.Marshal(replenishmentScheme)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
		}
		// 写入账本
		err = stub.PutState(utils.ConstructSchemeKey(retailerName), replenishmentSchemeJSON)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
		}

		// 序列化对象
		retailerJSON, err = json.Marshal(retailer)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
		}
		// 写入账本
		err = stub.PutState(retailerName, retailerJSON)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
		}

		// 更新当前零售商补货方案到map中
		schemesMap[retailerName] = replenishmentScheme
		// 序列化
		schemesMapJSON,err = json.Marshal(schemesMap)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
		}
		// 写入账本
		err = stub.PutState(lib.KeyOfSchemesMap,schemesMapJSON)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
		}
		
		// 使用匿名结构体存储要返回的内容
		res := struct {
			OldInventory int // 补货前库存量
			NewInventory int // 补货后库存量
		}{oldInventory, retailer.Inventory}
		// 序列化返回值
		resJSON, err := json.Marshal(res)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
		}
		
		return pb.Response{Status: 200, Message: "Pass successful", Payload: resJSON}
	}
}

// 零售商更新库存
// 参数： 零售商名称 新的库存量
// 返回： 空
func (t *MedicalSystem) retailerUpdateInventory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数个数
	if len(args) != 2 {
		return pb.Response{Status: 400, Message: "Incorrect number of arguments. Expecting 2", Payload: nil}
	}
	// 判断参数合法性（每个参数都不能为空）
	if args[0] == "" || args[1] == "" {
		return pb.Response{Status: 400, Message: "The parameter cannot be empty", Payload: nil}
	}
	// 将参数赋值给变量，部分参数需要转换数据类型
	retailerName := args[0] // 零售商名称
	// 将 string 转换为 int
	newInventory, err := strconv.Atoi(args[1]) // 新的库存量
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Conversion of data type failed: %s", err), Payload: nil}
	}

	// 读取账本，获取该零售商对象
	retailerJSON, err := stub.GetState(retailerName)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	} else if retailerJSON == nil {
		return pb.Response{Status: 400, Message: "The retailer does not exist", Payload: nil}
	}
	retailer := new(lib.Retailer)
	// 反序列化对象
	err = json.Unmarshal(retailerJSON, retailer)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
	}

	// 验证零售商信息已由供应商审核通过
	if retailer.State != lib.Pass {
		return pb.Response{Status: 400, Message: "The retailer failed the audit", Payload: nil}
	}

	// 更新库存量
	retailer.Inventory = newInventory

	// 序列化对象
	retailerJSON, err = json.Marshal(retailer)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
	}
	// 写入账本
	err = stub.PutState(retailerName, retailerJSON)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
	}

	// 计算订购点s
	// 订购点s = 提前期 x 需求量均值 + 残值
	s := retailer.LeadTime*retailer.AverageDemand + lib.ResidualValue
	// 如果库存量低于s时，计算补货数量
	var reorderQuantity int
	if retailer.Inventory < s {
		// 补货数量 = (T + 提前期) * 需求量 + 残值 - 当前报告的库存量
		reorderQuantity = (retailer.ReviewCycle+retailer.LeadTime)*retailer.AverageDemand + lib.ResidualValue - retailer.Inventory
	} else {
		reorderQuantity = 0
	}
	// 创建补货方案对象
	replenishmentScheme := lib.ReplenishmentScheme{
		RetailerName:    retailerName,
		ReorderQuantity: reorderQuantity,
		UnitPrice:       retailer.UnitPrice,
		ResponseResults: lib.ToBeResponded,
	}
	// 序列化对象
	replenishmentSchemeJSON, err := json.Marshal(replenishmentScheme)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
	}
	// 写入账本
	err = stub.PutState(utils.ConstructSchemeKey(retailerName), replenishmentSchemeJSON)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
	}

	// 更新供应商的补货方案Map
	// 获取补货方案map
	schemesMapJSON,err := stub.GetState(lib.KeyOfSchemesMap)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	}
	// 反序列化补货方案map
	schemesMap := make(map[string]lib.ReplenishmentScheme)
	err = json.Unmarshal(schemesMapJSON, &schemesMap)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
	}
	// 添加当前零售商补货方案到map中
	schemesMap[retailerName] = replenishmentScheme
	// 序列化
	schemesMapJSON,err = json.Marshal(schemesMap)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
	}
	// 写入账本
	err = stub.PutState(lib.KeyOfSchemesMap,schemesMapJSON)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
	}

	return pb.Response{Status: 200, Message: "Update successful", Payload: nil}
}

// 供货商通过与拒绝零售商注册
// 参数： 供应商名称 零售商名称 回应（0或1）
// 返回： 空
func (t *MedicalSystem) supplierAuditRegistration(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		errMes := "Incorrect number of arguments. Expecting 3"
		return pb.Response{Status: 400, Message: errMes, Payload: []byte(errMes)}
	}
	if args[0] == "" || args[1] == "" || args[2] == "" {
		return pb.Response{Status: 400, Message: "The parameter cannot be empty", Payload: nil}
	}
	if args[2] != "0" && args[2] != "1" {
		return pb.Response{Status: 400, Message: "The response result must be 0 or 1", Payload: nil}
	}
	supplierName := args[0]
	retailerName := args[1]
	result := args[2]

	// 读取账本，获取供应商
	supplierBytes, err := stub.GetState("supplier")
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	}
	// 验证供应商名称是否正确
	if string(supplierBytes) != supplierName {
		return pb.Response{Status: 400, Message: "Incorrect supplier name", Payload: nil}
	}

	// 读取账本，获取该零售商对象
	retailerJSON, err := stub.GetState(retailerName)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	} else if retailerJSON == nil {
		return pb.Response{Status: 400, Message: "The retailer does not exist", Payload: nil}
	}
	retailer := new(lib.Retailer)
	// 反序列化对象
	err = json.Unmarshal(retailerJSON, retailer)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
	}

	// 根据回应更改零售商状态
	if result == "0" {
		retailer.State = lib.Veto
	} else if result == "1" {
		retailer.State = lib.Pass

		// 生成该零售商的补货方案
		// 计算订购点 s ：  订购点s = 提前期 x 需求量均值 + 残值
		s := retailer.LeadTime*retailer.AverageDemand + lib.ResidualValue
		var reorderQuantity int
		// 如果库存量低于s时，计算补货数量
		if retailer.Inventory < s {
			// 补货数量 = (T + 提前期) * 需求量 + 残值 - 当前报告的库存量
			reorderQuantity = (retailer.ReviewCycle +retailer.LeadTime)*retailer.AverageDemand + lib.ResidualValue - retailer.Inventory
		} else {
			reorderQuantity = 0
		}
		// 创建补货方案对象
		replenishmentScheme := lib.ReplenishmentScheme{
			RetailerName:    retailerName,
			ReorderQuantity: reorderQuantity,
			UnitPrice:       retailer.UnitPrice,
			ResponseResults: lib.ToBeResponded,
		}
		// 序列化对象
		replenishmentSchemeJSON, err := json.Marshal(replenishmentScheme)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
		}
		// 写入账本
		err = stub.PutState(utils.ConstructSchemeKey(retailerName), replenishmentSchemeJSON)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
		}

		// 更新供应商的补货方案Map
		// 获取补货方案map
		schemesMapJSON,err := stub.GetState(lib.KeyOfSchemesMap)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
		}
		// 反序列化补货方案map
		schemesMap := make(map[string]lib.ReplenishmentScheme)
		err = json.Unmarshal(schemesMapJSON, &schemesMap)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
		}
		// 添加当前零售商补货方案到map中
		schemesMap[retailerName] = replenishmentScheme
		// 序列化
		schemesMapJSON,err = json.Marshal(schemesMap)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
		}
		// 写入账本
		err = stub.PutState(lib.KeyOfSchemesMap,schemesMapJSON)
		if err != nil {
			return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
		}
	}

	// 序列化对象
	retailerJSON, err = json.Marshal(retailer)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
	}
	// 写入账本
	err = stub.PutState(retailerName, retailerJSON)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("PutState error: %s", err), Payload: nil}
	}

	return pb.Response{Status: 200, Message: "Audit successful", Payload: nil}
}

// 供货商查看零售商们补货方案
// 参数： 供应商名称
// 返回： 补货方案列表
func (t *MedicalSystem) supplierViewSchemes(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数个数
	if len(args) != 1 {
		return pb.Response{Status: 400, Message: "Incorrect number of arguments. Expecting 1", Payload: nil}
	}
	// 判断参数合法性（每个参数都不能为空）
	if args[0] == "" {
		return pb.Response{Status: 400, Message: "The parameter cannot be empty", Payload: nil}
	}
	
	// 获取补货方案map
	schemesMapJSON,err := stub.GetState(lib.KeyOfSchemesMap)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("GetState error: %s", err), Payload: nil}
	}
	// 反序列化补货方案map
	schemesMap := make(map[string]lib.ReplenishmentScheme)
	err = json.Unmarshal(schemesMapJSON, &schemesMap)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Unmarshal error: %s", err), Payload: nil}
	}
	// 声明补货方案列表
	var schemesList []lib.ReplenishmentScheme
	// 遍历补货方案map，将所有补货方案添加到补货方案列表中
	for _, scheme := range schemesMap {
		schemesList = append(schemesList, scheme)
	}
	// 序列化补货方案列表
	schemesListJSON,err := json.Marshal(schemesList)
	if err != nil {
		return pb.Response{Status: 500, Message: fmt.Sprintf("Marshal error: %s", err), Payload: nil}
	}

	return pb.Response{Status: 200, Message: "View successful", Payload: schemesListJSON}
}

func main() {
	err := shim.Start(new(MedicalSystem))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
