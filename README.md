# API CRUD gRPC over HTTP

This repository contains an application of open source to provide HTTP CRUD operations over a gRPC server. All contributions are shared under the MIT license.

Please contact vadgun@hotmail.com to learn more.

- Go
- gRPC
- Protocol buffers

## What it does

The API CRUD gRPC over HTTP application, recieve incoming HTTP request to CREATE, READ, DELETE, and EDIT a product, and also you can create Orders with that products.

## Installation
Build the image
```bash
docker build -t my-go-app .
```

Run the image on detach exposing 8080 port
```bash
docker run -d -p 8080:8080 --name my-go-running-app my-go-app
```
## Functional Details

- First you will need to get a JWT to get access to CRUD operations, the application queries a POST request to /login to get access and responds a JWT, based on username and password parameters like this

```json
{
   "endpoint":"/login",
   "method":"POST",
   "body":{
      "username":"admin",
      "password":"password"
   },
   "200":{
      "token":"eyJh.....frS23wvrOBQ"
   },
   "401":{
      "error":"Invalid credentials"
   }
}
```

## Products CRUD

You will get a similar JWT which you need to copy for next requests like CREATE, READ, EDIT, DELETE

Copy token and include it in your request header named *Authorization* for the following endpoints.

- TO CREATE a product: 
```json
{
   "endpoint":"/products",
   "description":"This endpoint will create a product and retrieve its Id",
   "method":"POST",
   "headers": {
        "Authorization":"eyJh.....frS23wvrOBQ",
   },
   "body":{
      "name":"Air Pods",
      "description":"Apple headphones",
      "price":350,
      "quantity":4
   },
   "200":{
      "message":"Product created successfully!",
      "product":"d64d8a13-0a03-4cbb-a658-d2c86437b2ab"
   },
   "401":{
      "error":"Invalid Token"
   }
}
```

- TO READ products:
```json
{
   "endpoint":"/products",
   "description":"This endpoint will retrieve all products",
   "method":"GET",
   "headers":{
      "Authorization":"eyJh.....frS23wvrOBQ"
   },
   "body":{
   },
   "200":[
      {
         "id":"d64d8a13-0a03-4cbb-a658-d2c86437b2ab",
         "name":"Air Pods",
         "description":"Apple headphones",
         "price":350,
         "quantity":4
      }
   ],
   "401":{
      "error":"Invalid Token"
   }
}
```

- TO READ a single product by ID:
```json
{
   "endpoint":"/products/:id",
   "description":"This endpoint will retrieve a single product by ID",
   "method":"GET",
   "headers":{
      "Authorization":"eyJh.....frS23wvrOBQ"
   },
   "body":{
   },
   "200":{
      "id":"d64d8a13-0a03-4cbb-a658-d2c86437b2ab",
      "name":"Air Pods",
      "description":"Apple headphones",
      "price":350,
      "quantity":4
   },
   "404":{
      "message":"Product not found"
   }
}
```

- TO EDIT a product by ID:
```json
{
   "endpoint":"/products/:id",
   "description":"This endpoint will update a single product by ID",
   "method":"PATCH",
   "headers":{
      "Authorization":"eyJh.....frS23wvrOBQ"
   },
   "body":{
    "name":"Air Podzzzz",
    "description": "Apple headphones",
    "price":350,
    "quantity": 4
},
   "200":{
      "id":"d64d8a13-0a03-4cbb-a658-d2c86437b2ab",
      "name":"Air Podzzzz",
      "description":"Apple headphones",
      "price":350,
      "quantity":4
   },
   "404":{
      "message":"Product not found"
   }
}
```

- TO DELETE a product by ID:
```json
{
   "endpoint":"/products/:id",
   "description":"This endpoint will delete a single product by ID",
   "method":"DELETE",
   "headers":{
      "Authorization":"eyJh.....frS23wvrOBQ"
   },
   "body":{
   },
   "200":{
      "message":"Product deleted"
   },
   "404":{
      "message":"Product not found"
   }
}
```

## ORDERS CRUD
Orders check if its quantity can be provided and if the product exists to create the order, if any item in the order cannot be supplied, it will not be included in the created order instead of those that do have quantity and exist.


- TO CREATE an order:
```json
{
   "endpoint":"/orders/",
   "description":"This endpoint will create an order",
   "method":"POST",
   "headers":{
      "Authorization":"eyJh.....frS23wvrOBQ"
   },
   "body":{
      "items":[
         {
            "product_id":"ae1332e3-14fa-46b7-ab95-9122d362b867",
            "quantity":1
         },
         {
            "product_id":"d64d8a13-0a03-4cbb-a658-d2c86437b2ab",
            "quantity":1
         }
      ]
   },
   "200":{
      "message":"Order created successfully!",
      "order":"88166084-fb27-4d5d-ab3f-ed5468dcb2a3"
   },
   "401":{
      "error":"Invalid Token"
   }
}
```

- TO READ orders:
```json
{
   "endpoint":"/orders/",
   "description":"This endpoint will retrieve all orders",
   "method":"GET",
   "headers":{
      "Authorization":"eyJh.....frS23wvrOBQ"
   },
   "body":{
      
   },
   "200":[
      {
         "id":"88166084-fb27-4d5d-ab3f-ed5468dcb2a3",
         "items":[
            {
               "product_id":"d64d8a13-0a03-4cbb-a658-d2c86437b2ab",
               "quantity":1
            }
         ]
      }
   ],
   "401":{
      "error":"Invalid Token"
   }
}
```


