package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Order struct {
	BuyerId  int
	GoodId   int
	SellerId int
	Quantity int
}

func checkError(err error, msg string) {
	if err != nil {
		log.Println(err)
		defer log.Println(msg)
		panic(err)
	}
}

func (o Order) AddOrder() {
	db, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/golang")
	checkError(err, "Can't connect to mysql")
	defer db.Close()

	_, err = db.Exec("INSERT INTO `order`(buyer_id,good_id,seller_id,quantity)VALUES (?,?,?,?)", o.BuyerId, o.GoodId, o.SellerId, o.Quantity)
	checkError(err, "Can't add order")
}

func (o Order) ChangeInventory() {
	db, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/golang")
	checkError(err, "Can't connect to mysql")
	defer db.Close()

	var inventory int
	rows, err := db.Query("SELECT quantity FROM inventory  WHERE shop_id=?&&good_id=?", o.SellerId, o.GoodId)
	for rows.Next() {
		err := rows.Scan(&inventory)
		checkError(err, "Can't find inventory")
	}
	rows.Close()

	_, err = db.Exec("UPDATE inventory SET quantity=? WHERE shop_id=?&&good_id=?", inventory-o.Quantity, o.SellerId, o.GoodId)
	checkError(err, "Can't update inventory")

	if inventory != 0 {
		log.Println("The remaining quantity is", inventory-o.Quantity)
	}else {
		log.Println("Invalid parameters or zero inventory")
	}
}
