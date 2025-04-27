# Sales Analytics Project

## Project Overview

This project handles the processing of sales data from CSV files into a Maria database, using Golang 1.24 for backend services. It automates refreshing the database with customer, product, order, and order item details, along with maintaining a refresh log.

---

## Database Schema

### 1. Customers
| Field         | Type           | Details       |
|---------------|----------------|---------------|
| customer_id   | varchar(100)    | Primary Key   |
| name          | varchar(100)    |               |
| email         | varchar(100)    |               |
| address       | varchar(255)    |               |

### 2. Products
| Field         | Type           | Details       |
|---------------|----------------|---------------|
| product_id    | varchar(100)    | Primary Key   |
| product_name  | varchar(255)    |               |
| category      | varchar(100)    |               |

### 3. Orders
| Field         | Type           | Details       |
|---------------|----------------|---------------|
| order_id      | bigint          | Primary Key   |
| customer_id   | varchar(100)    | Foreign Key (Customers) |
| payment_method| varchar(100)    |               |
| date_of_sale  | datetime        |               |
| region        | varchar(100)    |               |

### 4. Order Items
| Field           | Type          | Details       |
|-----------------|---------------|---------------|
| order_item_id   | bigint         | Primary Key   |
| order_id        | bigint         | Foreign Key (Orders) |
| product_id      | varchar(100)   | Foreign Key (Products) |
| quantity_sold   | int            |               |
| unit_price      | decimal(10,2)  |               |
| shipping_cost   | decimal(10,2)  |               |
| discount        | decimal(10,2)  |               |

### 5. Sales Refresh Logs
| Field         | Type           | Details       |
|---------------|----------------|---------------|
| id            | int             | Primary Key   |
| status        | varchar(20)     |               |
| type          | varchar(20)     | (manaul/auto) |
| created_date  | datetime        |               |
| updated_date  | datetime        |               |

---

## Key Features

- **CSV Data Upload**: Parses and loads large CSV files into the database.
- **Auto Refresh Scheduler**: Uses Go routines and ticker to refresh the database daily.
- **API Endpoint**:  
  - `POST /refreshdata` → Triggers manual refresh from CSV.
  - `GET /revenue/total/2023-01-01/2024-01-01` → Total revenue based on given date range.
  - `GET /revenue//product/P123/2023-01-01/2024-01-01` → Revenue based on product and given date range.

- **Logging**: Status of refresh operations is tracked in `sales_refresh_logs`.

---

## Technologies Used

- **Golang** (Gin/Gorilla Mux for APIs)
- **MySQL** (Database)
- **GORM** (ORM for Go)
- **gocsv** (CSV Parsing)
- **toml** (TOML Configuration)


---
