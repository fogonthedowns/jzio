package orderprocessor

import (
	"context"
	"fmt"
	"log"
	"os"

	"jzio/src/internal/driver"
	"jzio/src/internal/handler"

	bx "github.com/fogonthedowns/orderbook"
)

func ExecuteOrders(ctx context.Context) (err error) {
	actions := make(chan *bx.Action)
	done := make(chan bool)

	// establish db connection
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	connection, err := driver.ConnectSQL(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		os.Exit(-1)
	}

	// start a transaction
	tx, err := connection.SQL.BeginTx(ctx, nil)

	// find Open orders
	orders, err := handler.LoadOpenOrders(ctx, tx)

	fmt.Printf("orders %+v \n", orders)

	// commit a transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// create a new order book
	ob := bx.NewOrderBook(actions)
	//go bx.ConsoleActionHandler(actions, done)

	filled := make([]bx.Action, 0)
	log2 := make(map[string]bx.ActionType)
	// new log inerts
	go func() {
		for {
			action := <-actions
			if action.ActionType == bx.AT_FILLED {
				filled = append(filled, *action)
			}

			log2[action.OrderId] = action.ActionType
			if action.ActionType == bx.AT_DONE {
				fmt.Printf("results%+v \n", filled)
				done <- true
				return
			}
		}
	}()

	// add Order to OrderBook
	for _, order := range orders {
		//fmt.Printf("order id %v\n", order.OrderPublicId)
		ob.AddOrder(bx.NewOrder(order.OrderPublicId, order.Side.IsBuy(), order.Cents, order.GoldTenthMilliGrams))
	}
	ob.Done()
	<-done
	select {
	case foo := <-done:
		fmt.Printf("does this get called %v", foo)
		for _, x := range filled {
			fmt.Printf("%+v \n", x)
		}
	default:
		fmt.Println("no activity")
	}

	return err
}
