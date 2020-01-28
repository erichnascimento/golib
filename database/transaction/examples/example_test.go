package examples_test

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/erichnascimento/golib/database/transaction/sql"
)

func Example() {
	db, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sql.RegisterConnector("main", db)
	cs := &CustomerService{
		storage:        &CustomerStorage{
			connection: sql.Connect("main"),
		},
		addressService: &AddressService{
			storage: &AddressStorage{
				connection: sql.Connect("main"),
			},
		},
	}

	ctx := context.Background()
	err = cs.Add(ctx, "erich", "5th avenue, 1")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// Output:
}



