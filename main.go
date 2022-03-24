package main

import (
	"errors"
	"fmt"
	"time"
)

func login() bool {
	var email, password string
	var choice int

	fmt.Println("_< Welcome to banking portal terminal app >_")
	fmt.Println("0: Login\n1: Create Account")

	fmt.Print("Enter choice: ")
	fmt.Scanln(&choice)

	if choice == 0 {
		fmt.Print("Enter Email: ")
		fmt.Scanln(&email)

		fmt.Print("Enter Password: ")
		fmt.Scanln(&password)

		if ValidateUser(email, password) {
			fmt.Println("Account verified!")
			fmt.Println("Redirecting to your profile...")
			userPage(email)
		} else {
			fmt.Println("Account does not exist!")
			login()
		}
	} else if choice == 1 {
		register()
	} else {
		fmt.Println("Invalid option!")
		login()
		return false
	}
	return true
}

func register() error {
	var name, email, password, passwordAgain string

	fmt.Println("Register a new user.")

	fmt.Print("Enter Name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter Email: ")
	fmt.Scanln(&email)

	// Validating if the email already exists in the database
	for _, user := range Users {
		if user.Email == email {
			return errors.New("account with the email exists")
		}
	}

	fmt.Print("Enter Password: ")
	fmt.Scanln(&password)

	fmt.Print("Enter Password Again: ")
	fmt.Scanln(&passwordAgain)

	// A while loop if passwords(password and passwordAgain) did not match
	for password != passwordAgain {
		fmt.Println("Passwords failed to match.")
		fmt.Println("Enter Password Again: ")
		fmt.Scanln(&passwordAgain)
	}

	var person Person = Person{Name: name, User: User{Email: email, Password: password}}

	accountHash := GenerateAccountNumber(person)
	person.AccountHash = accountHash
	var bankDetail BankDetail = BankDetail{AccountHash: accountHash, CurrentBalance: 0}

	Persons = append(Persons, person)
	Users = append(Users, person.User)
	BankDetails = append(BankDetails, bankDetail)

	fmt.Printf("Account created for %v \n", name)
	fmt.Println("Redirecting to login...")

	time.Sleep(1 * time.Second)

	login()

	return nil
}

func userPage(email string) {
	fmt.Printf("Welcome %v\n", email)
	fmt.Println("0: Check balance\n1: Withdraw money\n2: Deposit amount\n3: Transfer amount\n4: Logout")

	var choice int
	fmt.Scanln(&choice)

	var currentPerson Person

	for _, person := range Persons {
		if person.User.Email == email {
			currentPerson = person
		}
	}

	switch choice {
	case 0:
		bankDetails, _, _, _ := GetData(currentPerson, "bankDetails")
		fmt.Printf("Current bank balance: %v\n", bankDetails.CurrentBalance)
		userPage(email)
	case 1:
		/**
		Logic for withdrawing money
		- Deduct the amount from current balance
		- if current balance is 0 prompt to increase the balance first
		**/
		var amount, previousBalance, currentBalance float64
		fmt.Println("Enter amount to withdraw: ")
		fmt.Scanln(&amount)

		for i := range BankDetails {
			currentIter := &BankDetails[i]
			if currentIter.AccountHash == currentPerson.AccountHash && currentIter.CurrentBalance != 0 {
				if amount > 0 {
					previousBalance = currentIter.CurrentBalance
					currentIter.CurrentBalance -= amount
					currentBalance = currentIter.CurrentBalance
				} else {
					fmt.Println("Amount can't be <= 0 !")
					userPage(email)
				}
			} else if currentIter.AccountHash == currentPerson.AccountHash && currentIter.CurrentBalance <= 0 {
				fmt.Printf("Your balance seems to be low by: %v", previousBalance-amount)
				userPage(email)
			}
		}
		fmt.Printf("Previous Balance: %v\nWithdrawn amount: %v\nCurrent Balance: %v\n", previousBalance, amount, currentBalance)
		userPage(email)
	case 2:
		/**
		Logic for deposit money
		- Add the amount to current balance
		**/
		var amount, previousBalance, currentBalance float64
		fmt.Println("Enter amount to deposit: ")
		fmt.Scanln(&amount)

		// if len(BankDetails) == 0 {
		// 	BankDetails = append(BankDetails, BankDetail{AccountHash: currentPerson.AccountHash, CurrentBalance: amount})
		// 	if amount > 0 {
		// 		previousBalance = 0
		// 		currentBalance = amount
		// 		fmt.Printf("Previous Balance: %v\nDeposited amount: %v\nCurrent Balance: %v\n", previousBalance, amount, currentBalance)
		// 		userPage(email)
		// 	} else {
		// 		fmt.Println("Amount can't be <= 0 !")
		// 		userPage(email)
		// 	}
		// } else {
		for i := range BankDetails {
			currentIter := &BankDetails[i]
			if currentIter.AccountHash == currentPerson.AccountHash && currentIter.CurrentBalance >= 0 {
				if amount > 0 {
					previousBalance = currentIter.CurrentBalance
					currentIter.CurrentBalance += amount
					currentBalance = currentIter.CurrentBalance
					fmt.Printf("Previous Balance: %v\nDeposited amount: %v\nCurrent Balance: %v\n", previousBalance, amount, currentBalance)
					userPage(email)
				} else {
					fmt.Println("Amount can't be <= 0 !")
					userPage(email)
				}
			}
		}
		BankDetails = append(BankDetails, BankDetail{AccountHash: currentPerson.AccountHash, CurrentBalance: amount})
		// }
	case 3:
		/**
		Logic for transferring money
		- Deduct the amount of the sender's current balance
		- Increse the amount of the reciever's current balance
		- The sending and the recieving amount can't be 0
		**/
		var amount float64
		var recieverEmail string

		fmt.Println("Enter the reciever's email and amount: ")
		fmt.Scanln(&recieverEmail, &amount)

		for i := range Persons {
			if Persons[i].User.Email == recieverEmail {
				reciever := &Persons[i]
				for i := range BankDetails {
					currentIter := &BankDetails[i]
					if currentIter.AccountHash == reciever.AccountHash {
						currentIter.CurrentBalance += amount
						fmt.Println(currentIter)
					} else if currentIter.AccountHash == currentPerson.AccountHash {
						currentIter.CurrentBalance -= amount
						fmt.Println(currentIter)
					}
				}
				TransferTransactions = append(TransferTransactions, TransferTransaction{Amount: amount, From: currentPerson.AccountHash, To: reciever.AccountHash})
				fmt.Println(TransferTransactions)
				fmt.Printf("Amount send: %v\nTo: %v\n", amount, reciever.AccountHash)
			}
		}

		userPage(email)
	case 4:
		fmt.Println("Logging out...")
		login()
	}
}

func main() {
	// var user = Person{Name: "sarvesh", AccountHash: "1234", User: User{Email: "Sarvesh", Password: "Password"}}
	login()
}
