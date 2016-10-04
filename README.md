*A*PIM *R*EST *C*lient

Currently you can add and get API Products

1. Install [Go](https://golang.org/) if you havent already and clone the repository

2. Run the command `go build arc.go` to build

3. Usage
    
    To generate template *data/product.json* for creating a product - `./arc -create-data=product`
    To create a new product using the data in *data/product.json* -  `./arc -resource=publisher:create:product`
    To get list of existing products - `./arc -resource=publisher:view:products`