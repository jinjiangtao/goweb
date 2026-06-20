package main

import (
	"fmt"
	"logistics-server/database"
	"logistics-server/models"
	"math/rand"
	"time"
)

func main() {
	database.Init()

	rand.Seed(time.Now().UnixNano())

	orders := []models.LogisticsOrder{
		{
			OrderNo:     "WL202401001",
			Sender:      "张三",
			Receiver:    "李四",
			Origin:      "北京市朝阳区",
			Destination: "上海市浦东新区",
			GoodsName:   "电子产品",
			Weight:      2.5,
			Status:      models.OrderStatusInTransit,
		},
		{
			OrderNo:     "WL202401002",
			Sender:      "王五",
			Receiver:    "赵六",
			Origin:      "广州市天河区",
			Destination: "深圳市南山区",
			GoodsName:   "服装",
			Weight:      1.2,
			Status:      models.OrderStatusDelivered,
		},
		{
			OrderNo:     "WL202401003",
			Sender:      "孙七",
			Receiver:    "周八",
			Origin:      "成都市武侯区",
			Destination: "重庆市渝中区",
			GoodsName:   "食品",
			Weight:      5.0,
			Status:      models.OrderStatusException,
		},
	}

	for i := range orders {
		database.DB.Create(&orders[i])

		nodes := generateNodes(orders[i].ID, orders[i].Status)
		for _, node := range nodes {
			database.DB.Create(&node)
		}
	}

	fmt.Println("测试数据初始化完成！")
}

func generateNodes(orderID uint, orderStatus string) []models.LogisticsNode {
	baseTime := time.Now().AddDate(0, 0, -3)

	nodes := []models.LogisticsNode{
		{
			OrderID:     orderID,
			NodeName:    "订单创建",
			Status:      models.NodeStatusPending,
			Location:    "系统",
			Description: "订单已创建，等待揽收",
			Operator:    "系统",
			IsAbnormal:  false,
			OccurredAt:  baseTime,
		},
		{
			OrderID:     orderID,
			NodeName:    "已揽收",
			Status:      models.NodeStatusPicked,
			Location:    "发货网点",
			Description: "快递员已揽收包裹",
			Operator:    "张快递员",
			IsAbnormal:  false,
			OccurredAt:  baseTime.Add(2 * time.Hour),
		},
	}

	if orderStatus != models.OrderStatusPending {
		nodes = append(nodes, models.LogisticsNode{
			OrderID:     orderID,
			NodeName:    "运输中",
			Status:      models.NodeStatusInTransit,
			Location:    "中转中心",
			Description: "包裹已到达中转中心，正在分拣",
			Operator:    "分拣中心",
			IsAbnormal:  false,
			OccurredAt:  baseTime.Add(8 * time.Hour),
		})
	}

	if orderStatus == models.OrderStatusException {
		nodes = append(nodes, models.LogisticsNode{
			OrderID:        orderID,
			NodeName:       "异常件",
			Status:         models.NodeStatusException,
			Location:       "运输途中",
			Description:    "包裹在运输过程中出现异常",
			Operator:       "异常处理组",
			IsAbnormal:     true,
			AbnormalReason: "地址信息不完整，需要联系收件人确认",
			OccurredAt:     baseTime.Add(24 * time.Hour),
		})
	}

	if orderStatus == models.OrderStatusArrived || orderStatus == models.OrderStatusDelivered {
		nodes = append(nodes, models.LogisticsNode{
			OrderID:     orderID,
			NodeName:    "已到达",
			Status:      models.NodeStatusArrived,
			Location:    "目的城市",
			Description: "包裹已到达目的城市，正在派送",
			Operator:    "派送网点",
			IsAbnormal:  false,
			OccurredAt:  baseTime.Add(48 * time.Hour),
		})
	}

	if orderStatus == models.OrderStatusDelivered {
		nodes = append(nodes, models.LogisticsNode{
			OrderID:     orderID,
			NodeName:    "已签收",
			Status:      models.NodeStatusDelivered,
			Location:    "收件地址",
			Description: "包裹已成功签收",
			Operator:    "李快递员",
			IsAbnormal:  false,
			OccurredAt:  baseTime.Add(52 * time.Hour),
		})
	}

	return nodes
}
