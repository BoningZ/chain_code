package main

import (
	//	"crypto/elliptic"
	//	"crypto/rand"
	//	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//	"github.com/tjfoc/gmsm/sm2"
	//	"math/big"
	"strconv"
)

// Order 交易数据：
//订单ID
//商品名称
//商品数量
//订单金额
//下单时间
//买家ID
//卖家ID
//
//属性信息：
//订单状态（待支付、已支付、已发货、已完成等）
//物流信息（运输公司、运单号等）
//买家评价
//卖家评价
type Order struct {
	OrderID string `json:"orderId"`
	//	Name         []byte  `json:"name"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Amount       float64 `json:"amount"`
	OrderTime    string  `json:"orderTime"`
	BuyerID      string  `json:"buyerId"`
	SellerID     string  `json:"sellerId"`
	Status       string  `json:"status"`
	Logistics    string  `json:"logistics"`
	BuyerReview  string  `json:"buyerReview"`
	SellerReview string  `json:"sellerReview"`
}

type OrderChaincode struct {
}

func (cc *OrderChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *OrderChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createOrder" {
		return cc.createOrder(stub, args)
	} else if function == "getOrder" {
		return cc.getOrder(stub, args)
	} else if function == "updateOrder" {
		return cc.updateOrderStatus(stub, args)
	} else if function == "addBuyerReview" {
		return cc.addBuyerReview(stub, args)
	} else if function == "addSellerReview" {
		return cc.addSellerReview(stub, args)
	} else {
		return shim.Error("Invalid function name.")
	}
}

func (cc *OrderChaincode) createOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//if len(args) != 8 {
	//	return shim.Error("Expecting 8 arguments: OrderID, Name, Quantity, Amount, OrderTime, BuyerID, SellerID, Pubkey")
	//}
	if len(args) != 7 {
		return shim.Error("Expecting 7 arguments: OrderID, Name, Quantity, Amount, OrderTime, BuyerID, SellerID, now have " + strconv.Itoa(len(args)))
	}

	orderID := args[0]
	name := args[1]
	quantity, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid quantity. Expecting an integer value.")
	}
	amount, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return shim.Error("Invalid amount. Expecting a float value.")
	}
	orderTime := args[4]
	buyerID := args[5]
	sellerID := args[6]
	//pubkey := args[7]

	// 对敏感字段进行加密
	//encryptedName, err := encryptData([]byte(name), pubkey)
	encryptedName := name
	if err != nil {
		return shim.Error("Failed to encrypt name.")
	}

	order := Order{
		OrderID:      orderID,
		Name:         encryptedName,
		Quantity:     quantity,
		Amount:       amount,
		OrderTime:    orderTime,
		BuyerID:      buyerID,
		SellerID:     sellerID,
		Status:       "Pending",
		Logistics:    "",
		BuyerReview:  "",
		SellerReview: "",
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return shim.Error("Failed to marshal order.")
	}

	err = stub.PutState(orderID, orderJSON)
	if err != nil {
		return shim.Error("Failed to save order.")
	}

	return shim.Success(nil)
}

func (cc *OrderChaincode) getOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//if len(args) != 2 {
	//	return shim.Error("Expecting 2 argument: OrderID, PrivateKey")
	//}
	if len(args) != 1 {
		return shim.Error("Expecting 1 argument: OrderID")
	}

	orderID := args[0]
	//privateKey := args[1]

	orderJSON, err := stub.GetState(orderID)
	if err != nil {
		return shim.Error("Failed to get order.")
	}

	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return shim.Error("Failed to unmarshal order.")
	}

	// 对敏感字段进行解密
	//decryptedName, err := decryptData([]byte(order.Name), privateKey)
	decryptedName := order.Name
	if err != nil {
		return shim.Error("Failed to decrypt name.")
	}

	order.Name = decryptedName

	orderJSON, err = json.Marshal(order)
	if err != nil {
		return shim.Error("Failed to marshal order.")
	}

	return shim.Success(orderJSON)
}
func (cc *OrderChaincode) updateOrderStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Expecting 2 arguments: OrderID, Status")
	}

	orderID := args[0]
	status := args[1]

	orderJSON, err := stub.GetState(orderID)
	if err != nil {
		return shim.Error("Failed to get order.")
	}

	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return shim.Error("Failed to unmarshal order.")
	}

	// 更新订单状态
	order.Status = status

	orderJSON, err = json.Marshal(order)
	if err != nil {
		return shim.Error("Failed to marshal order.")
	}

	err = stub.PutState(orderID, orderJSON)
	if err != nil {
		return shim.Error("Failed to update order.")
	}

	return shim.Success(nil)
}

func (cc *OrderChaincode) addBuyerReview(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Expecting 2 arguments: OrderID, BuyerReview")
	}

	orderID := args[0]
	buyerReview := args[1]

	orderJSON, err := stub.GetState(orderID)
	if err != nil {
		return shim.Error("Failed to get order.")
	}

	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return shim.Error("Failed to unmarshal order.")
	}

	// 添加买家评价
	order.BuyerReview = buyerReview

	orderJSON, err = json.Marshal(order)
	if err != nil {
		return shim.Error("Failed to marshal order.")
	}

	err = stub.PutState(orderID, orderJSON)
	if err != nil {
		return shim.Error("Failed to update order.")
	}

	return shim.Success(nil)
}

func (cc *OrderChaincode) addSellerReview(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Expecting 2 arguments: OrderID, SellerReview")
	}

	orderID := args[0]
	sellerReview := args[1]

	orderJSON, err := stub.GetState(orderID)
	if err != nil {
		return shim.Error("Failed to get order.")
	}

	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return shim.Error("Failed to unmarshal order.")
	}

	// 添加卖家评价
	order.SellerReview = sellerReview

	orderJSON, err = json.Marshal(order)
	if err != nil {
		return shim.Error("Failed to marshal order.")
	}

	err = stub.PutState(orderID, orderJSON)
	if err != nil {
		return shim.Error("Failed to update order.")
	}

	return shim.Success(nil)
}

//func encryptData(data []byte, publicKeyStr string) ([]byte, error) {
//	// 解码公钥字符串
//	publicKeyBytes, err := hex.DecodeString(publicKeyStr)
//	if err != nil {
//		return nil, err
//	}
//	// 解析公钥
//	x, y := elliptic.Unmarshal(sm2.P256Sm2(), publicKeyBytes)
//	publicKey := &sm2.PublicKey{
//		Curve: sm2.P256Sm2(),
//		X:     x,
//		Y:     y,
//	}
//	// 加密数据
//	ciphertext, err := sm2.EncryptAsn1(publicKey, data, rand.Reader)
//	if err != nil {
//		return nil, err
//	}
//	return ciphertext, nil
//}
//
//func decryptData(ciphertext []byte, privateKeyStr string) ([]byte, error) {
//	// 解码私钥字符串
//	privateKeyBytes, err := hex.DecodeString(privateKeyStr)
//	if err != nil {
//		return nil, err
//	}
//	// 解析私钥
//	privateKey := new(sm2.PrivateKey)
//	privateKey.Curve = sm2.P256Sm2()
//	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
//	privateKey.PublicKey.Curve = privateKey.Curve
//	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.Curve.ScalarBaseMult(privateKeyBytes)
//	// 解密数据
//	plaintext, err := privateKey.DecryptAsn1(ciphertext)
//	if err != nil {
//		return nil, err
//	}
//	return plaintext, nil
//}

func main() {
	err := shim.Start(new(OrderChaincode))
	if err != nil {
		fmt.Printf("Error starting OrderChaincode: %s", err)
	}
}
