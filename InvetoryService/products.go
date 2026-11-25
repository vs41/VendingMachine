package invetoryservice

type Product struct {
	productID    string
	productName  string
	productCode  string
	productPrice int64 // Price in paise
}

func NewProduct(
	productID string,
	productName string,
	productCode string,
	productPrice int64,
) *Product {
	return &Product{
		productID:    productID,
		productName:  productName,
		productCode:  productCode,
		productPrice: productPrice,
	}
}

func (p *Product) GetProductID() string {
	return p.productID
}

func (p *Product) GetProductCode() string {
	return p.productCode
}

func (p *Product) GetProductName() string {
	return p.productName
}

func (p *Product) GetProductPrice() int64 {
	return p.productPrice
}
