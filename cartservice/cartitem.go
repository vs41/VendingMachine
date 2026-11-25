package cartservice

type CartItem struct {
	itemId       string
	itemQuantity int
	unitPrice    int64
	slotId       string
	floorId      string
	itemName     string
}

func NewCartItem(
	itemId string,
	itemQuantity int,
	unitPrice int64,
	slotId string,
	floorId string,
	itemName string,
) *CartItem {
	return &CartItem{
		itemId:       itemId,
		itemQuantity: itemQuantity,
		unitPrice:    unitPrice,
		slotId:       slotId,
		floorId:      floorId,
		itemName:     itemName,
	}
}

func (ci *CartItem) GetItemID() string {
	return ci.itemId
}

func (ci *CartItem) GetItemQuantity() int {
	return ci.itemQuantity
}

func (ci *CartItem) GetTotalPrice() int64 {
	return ci.unitPrice * int64(ci.itemQuantity)
}

func (ci *CartItem) GetSlotID() string {
	return ci.slotId
}

func (ci *CartItem) GetFloorID() string {
	return ci.floorId
}

func (ci *CartItem) GetItemName() string {
	return ci.itemName
}

func (ci *CartItem) GetUnitPrice() int64 {
	return ci.unitPrice
}

func (ci *CartItem) DecreaseQuantityByOne() {
	ci.itemQuantity--
}

func (ci *CartItem) increaseQuantityByOne() {
	ci.itemQuantity++
}
