# REST API Starter

## Installation
1. copy .env
2. create db_name
3. run `go get` 
4. run migration
5. run seeder
6. run `go run main.go`
7. copy Token
8. enjoy

## Example Request Body Login & Register
Register payload
```JSON
{
 	"name":"admin",
	"email":"admin@gmail.com",
	"password":"password123"
}
```
Login payload
```JSON
{
	"email":"user@gmail.com",
	"password":"password123"
}
```

## Example Request Body Product 
Create/Update Purchase
```JSON
{
	"name": "Tunai",
	"price":120000
}
```

## Example Request Body Cart
Add Cart
```JSON
{
	"product_id": 10,
	"qty": 10
}
```
Pay Cart
```JSON
{
	"payment_method": "QRIS",
	"cart_details":[
        {
            "product_id": 1,
            "qty": 2
        }
	]
}
```

## Example Request Body Purchase 
Create/Update Purchase
```JSON
{
	"payment_method": "Tunai",
	"purchase_details":[
        {
            "product_id": 1,
            "qty": 2
        }
	]
}
```
## Available Endpoint
### Auth
POST    | /api/auth/register
POST    | /api/auth/login
POST    | /api/auth/logout

### User 
GET     | /api/users/
POST    | /api/users/
GET     | /api/users/:id
PUT     | /api/users/:id
DELETE  | /api/users/:id

### Product 
GET     | /api/products/
POST    | /api/products/
GET     | /api/products/:id
PUT     | /api/products/:id
DELETE  | /api/products/:id

### Cart 
GET     | /api/carts/
POST    | /api/carts/add
POST    | /api/carts/remove/:productId
POST    | /api/carts/pay

### Purchase 
GET     | /api/purchase/
POST    | /api/purchase/
GET     | /api/purchase/:id
PUT     | /api/purchase/:id
DELETE  | /api/purchase/:id

## Migration
### Using command line
To create migration: 
`migrate create -ext=sql -dir=database/migrations -seq user`

To execute migration: 
`migrate -path=database/migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/tb_name?charset=utf8mb4&parseTime=True&loc=Local" -verbose up`

To rollback migration:
`migrate -path=database/migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/tb_name?charset=utf8mb4&parseTime=True&loc=Local" -verbose down (VERSION)`

### Using makefile
To create migration:
`make create_migration name=migrationName`

To execute migration:
`make migrate_up`

To rollback migration:
`make migrate_down`

To execute migration:
`make migrate_up_seed`

## Seeding
### Using command line
To seed all:
`go run database/seeder/main.go`

To seed specific, append the types at the back: 
`go run database/seeder/main.go account_type`
