package cartservice

import "errors"

var (
	ErrItemNotFound    = errors.New("cart item not found")
	ErrInvalidQuantity = errors.New("invalid quantity")
)

type Cart struct {
	cartId    string
	userId    string
	cartItems []*CartItem
}

func NewCart(cartId, userId string) *Cart {
	return &Cart{
		cartId:    cartId,
		userId:    userId,
		cartItems: make([]*CartItem, 0),
	}
}

func (c *Cart) GetCartID() string {
	return c.cartId
}

func (c *Cart) GetUserID() string {
	return c.userId
}

func (c *Cart) GetCartItems() []*CartItem {
	items := make([]*CartItem, len(c.cartItems))
	copy(items, c.cartItems)
	return items
}

func (c *Cart) GetTotal() int64 {
	var total int64

	for _, item := range c.cartItems {
		total += item.GetTotalPrice()
	}

	return total
}

func (c *Cart) FindItem(itemId string) (int, *CartItem) {
	for i, item := range c.cartItems {
		if item.itemId == itemId {
			return i, item
		}
	}

	return -1, nil
}

func (c *Cart) AddItem(item *CartItem) error {
	if item == nil || item.itemQuantity <= 0 {
		return ErrInvalidQuantity
	}

	c.cartItems = append(c.cartItems, item)
	return nil
}

func (c *Cart) RemoveItem(itemId string) error {
	index, _ := c.FindItem(itemId)

	if index == -1 {
		return ErrItemNotFound
	}

	c.cartItems = append(
		c.cartItems[:index],
		c.cartItems[index+1:]...,
	)

	return nil
}

func (c *Cart) IncreaseQuantity(itemId string) error {
	_, item := c.FindItem(itemId)

	if item == nil {
		return ErrItemNotFound
	}

	item.increaseQuantityByOne()
	return nil
}

func (c *Cart) DecreaseQuantity(itemId string) error {
	index, item := c.FindItem(itemId)

	if item == nil {
		return ErrItemNotFound
	}

	if item.itemQuantity <= 1 {
		return c.RemoveItem(itemId)
	}

	item.DecreaseQuantityByOne()
	_ = index

	return nil
}

func (c *Cart) Clear() {
	c.cartItems = make([]*CartItem, 0)
}
