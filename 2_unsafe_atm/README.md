# The Unsafe ATM

A Go program that maintains a bank account with multiple users simultaneously withdrawing and depositing money.

This program is currently not thread-safe, which means that simultaneous transactions with the bank account cause money to either duplicate or disappear mysteriously!
The program kicks off multiple threads that each repeatedly deposit and withdraw $50. Because they always withdraw and deposit the same amount, the account balance should *NOT* change. Running the program a few times, though, demonstrates that the balance does not remain the same in the end. This is because the `account` object does not have any kind of locking mechanism, which means that the account can easily be corrupted by parallel actors using it at once. This kind of problem is called a [Race Condition](https://en.wikipedia.org/wiki/Race_condition).

## Project Instructions
All work will be done in `main.go`.

1. Build and run the code without any modifications. After a few runs, you should see some output indicating that the balance changed from the `startingBalance` when it shouldn't have. This is the race condition that we will fix.
2. Add a lock (`sync.Mutex`) to the `account` struct. This will be the lock that we use to coordinate access to the account.
3. At the start of the `Deposit` function, call `Lock()` on your new lock object. Follow it immediately with a [deferred](https://tour.golang.org/flowcontrol/12) `Unlock()` call.
4. Repeat step 2 for the `Withdraw` function.
5. Build and run the code! The balance should now always be the same as the `startingBalance` variable, meaning that all deposits and withdrawals worked as expected.