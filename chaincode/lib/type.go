package lib

const (
	ToBeResponded = "ToBeResponded"
	Pass          = "Pass"
	Veto          = "Veto"
)

var (
	ResidualValue   = 0 // 残值
	KeyOfSupplier   = "supplier"
	KeyOfSchemesMap = "replenishmentSchemes"
)

// Retailer 零售商
type Retailer struct {
	RetailerName       string  `json:"retailer_name"`        // 零售商名称
	UnitPrice          float64 `json:"unit_price"`           // 订货单价
	LeadTime           int     `json:"lead_time"`            // 提前期
	Inventory          int     `json:"inventory"`            // 库存量
	AverageDemand      int     `json:"average_demand"`       // 需求量均值
	UpdateCycle        int     `json:"update_cycle"`         // 上传数据的周期
	InventoryValue     float64 `json:"inventory_value"`      // 库存商品价值
	AnnualInterestRate float64 `json:"annual_interest_rate"` // 年利率
	FixedOrderCost     float64 `json:"fixed_order_cost"`     // 固定订货成本
	ReviewCycle        int     `json:"review_cycle"`         // 审查周期
	State              string  `json:"state"`                // 帐号状态（待审核、通过、否决）
}

// ReplenishmentScheme 补货方案
type ReplenishmentScheme struct {
	RetailerName    string  `json:"retailer_name"`    // 零售商名称
	ReorderQuantity int     `json:"reorder_quantity"` // 补货数量
	UnitPrice       float64 `json:"unit_price"`       // 单价
	ResponseResults string  `json:"response_results"` // 回应结果
}
