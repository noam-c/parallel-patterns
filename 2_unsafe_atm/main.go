package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const startingBalance = 500

// A simplistic bank account. Allows you to withdraw, deposit, or
// retrieve the account balance.
type account struct {
	balance int
}

func (a *account) updateBalance(balance int) {
	// Sleep for a millisecond to simulate taking some time to update the
	// account balance.
	time.Sleep(1 * time.Millisecond)

	a.balance = balance
}

// Withdraw removes money from the account, but only if there is
// enough money to cover the withdrawal.
func (a *account) Withdraw(amt int) int {
	if a.balance >= amt {
		a.updateBalance(a.balance - amt)
	}

	return a.balance
}

// Deposit adds money to the account.
func (a *account) Deposit(amt int) int {
	a.updateBalance(a.balance + amt)
	return a.balance
}

// Balance retrieves the account's balance.
func (a *account) Balance() int {
	return a.balance
}

// Log helper to also print the current time.
func log(messages ...interface{}) {
	var args []interface{}
	args = append(args, time.Now().Format(time.RFC3339Nano))
	args = append(args, messages...)
	fmt.Println(args...)
}

func main() {
	// Ensure that the app runs on more than one core
	runtime.GOMAXPROCS(2)

	// Create a new account
	a := account{balance: startingBalance}
	wg := sync.WaitGroup{}

	// Several account users will add and remove money from the account
	// repeatedly. Since they always withdraw exactly what they deposit, there
	// should be no change to the account balance in the end. If there is, then
	// we have a race condition and the account is NOT thread safe.
	numThreads := 16
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func(id int) {
			// Each thread deposits and withdraws several times.
			for j := 0; j < 50; j++ {
				amt := 50
				log("Thread", id, ": Depositing $", amt)
				log("Thread", id, ": New balance is $", a.Deposit(amt))

				log("Thread", id, ": Withdrawing $", amt)
				log("Thread", id, ": New balance is $", a.Withdraw(amt))
			}

			wg.Done()
		}(i)
	}

	// Wait for all the account users to finish
	wg.Wait()

	finalBalance := a.Balance()
	if finalBalance != startingBalance {
		fmt.Println("WOAH! Balance changed! It is now:", finalBalance)
	} else {
		fmt.Println("Finished. Balance is just fine at:", finalBalance)
	}
}
