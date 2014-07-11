package controllers
//
//import (
//	"code.google.com/p/goprotobuf/proto"
////	"common"
//	"fmt"
//	"controllers/models"
//	"controllers/models/configs"
//	pd "protodata"
//	pb "protodata/protodata"
//)
//
//// 10291 充值列表
//func (role *Role) RechargeList(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	response := &pb.RechargeListResponse{
//		ECode: proto.Int32(1),
//	}
//
//	for _, pc := range configs.ConfigPayCenterList() {
//
//		response.RechargeId = append(response.RechargeId, int32(pc.Id))
//		response.Ingot = append(response.Ingot, int32(pc.Ingot))
//		response.Money = append(response.Money, int32(pc.Rmb))
//		response.Name = append(response.Name, pc.Name)
//	}
//
//	response.Count = proto.Int32(int32(len(response.RechargeId)))
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 10292 创建订单
//func (role *Role) RechargeOrder(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	request := &pb.RechargeRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	payCenterId := int(request.GetRechargeId())
//
//	payData := configs.GetPayCenterById(payCenterId)
//
//	if payData.Id == 0 {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("PayCenterId Error ! ")
//	}
//
//	orderId, err := models.CreateOrder(player.UniqueId, ServerId, payCenterId, payData.Rmb, payData.Ingot)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK,
//		&pb.RechargeResponse{
//			ECode:        proto.Int32(1),
//			OrderId:      proto.String(orderId),
//			ProductIndex: proto.Int32(request.GetProductIndex()),
//		}), nil
//
//}

// 10293 核对订单
//func (role *Role) CheckRecharge(player *Player, cq *pb.CommandRequest) (string, error) {
//	request := &pb.CheckRechargeRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	// read order from db
//	order, err := models.GetOrderById(request.GetOrderId())
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	if order.Unique != player.UniqueId {
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.CheckRechargeResponse{
//			ECode: proto.Int32(2),
//			EStr:  proto.String(common.ESTR_order_not_exit),
//		}), nil
//	}
//
//	var result bool
//	if order.Status == 5 {
//		result = true
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.CheckRechargeResponse{
//		ECode:  proto.Int32(1),
//		Result: proto.Bool(result),
//		Ingot:  proto.Int32(int32(order.PcIngot)),
//	}), nil
//}
