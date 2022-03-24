package main

import (
	"fmt"
	"time"
)

func login() bool {
	var email, password string
	var choice int

	fmt.Println("_< Welcome to banking portal terminal app >_")
	fmt.Println("0: Login\n1: Create Account")

	fmt.Print(">> ")
	fmt.Scanln(&choice)

	if choice == 0 {
		fmt.Print("Enter Email: ")
		fmt.Scanln(&email)

		fmt.Print("Enter Password: ")
		fmt.Scanln(&password)

		if ValidateUser(email, password) {
			fmt.Println("- Account verified!")
			fmt.Print("- Redirecting to your profile...\n\n")
			userPage(email)
		} else {
			fmt.Print("! > Account does not exist!\n\n")
			login()
		}
	} else if choice == 1 {
		register()
	} else {
		fmt.Println("! > Invalid option!")
		login()
		return false
	}
	return true
}

func register() {
	var name, email, password, passwordAgain string

	fmt.Println("~ Register a new user ~")

	fmt.Print("Enter Name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter Email: ")
	fmt.Scanln(&email)

	// Validating if the email already exists in the database
	for _, user := range Users {
		if user.Email == email {
			fmt.Println("! > Email already exists")
			register()
		}
	}

	fmt.Print("Enter Password: ")
	fmt.Scanln(&password)

	fmt.Print("Enter Password Again: ")
	fmt.Scanln(&passwordAgain)

	// A while loop if passwords(password and passwordAgain) did not match
	for password != passwordAgain {
		fmt.Println("! > Passwords failed to match.")
		fmt.Println(">> Enter Password Again: ")
		fmt.Scanln(&passwordAgain)
	}

	var person Person = Person{Name: name, User: User{Email: email, Password: password}}

	accountHash := GenerateAccountNumber(person)
	person.AccountHash = accountHash
	var bankDetail BankDetail = BankDetail{AccountHash: accountHash, CurrentBalance: 0}

	Persons = append(Persons, person)
	Users = append(Users, person.User)
	BankDetails = append(BankDetails, bankDetail)

	fmt.Printf("- Account created for %v \n", name)
	fmt.Print("- Redirecting to login...\n\n")

	time.Sleep(1 * time.Second)

	login()
}

func userPage(email string) {
	fmt.Printf("~ Welcome %v ~\n", email)
	fmt.Println("0: Check balance\n1: Withdraw money\n2: Deposit amount\n3: Transfer amount\n4: Logout")

	var choice int
	fmt.Print(">> ")
	fmt.Scanln(&choice)

	var currentPerson Person

	for _, person := range Persons {
		if person.User.Email == email {
			currentPerson = person
		}
	}

	switch choice {
	case 0:
		/**
		Logic for checking user's bank balance
		- Fetches bank details with helper method GetData
		**/
		bankDetails, _, _, _ := GetData(currentPerson, "bankDetails")
		fmt.Printf("$ > Current bank balance: %v\n\n", bankDetails.CurrentBalance)
		userPage(email)
	case 1:
		/**
		Logic for withdrawing money
		- Deduct the amount from current balance
		- if current balance is 0 prompt to increase the balance first
		**/
		var amount, previousBalance, currentBalance float64
		fmt.Print(">> Enter amount to withdraw: ")
		fmt.Scanln(&amount)

		for i := range BankDetails {
			currentIter := &BankDetails[i]
			if currentIter.AccountHash == currentPerson.AccountHash && currentIter.CurrentBalance != 0 {
				if amount > 0 {
					previousBalance = currentIter.CurrentBalance
					currentIter.CurrentBalance -= amount
					currentBalance = currentIter.CurrentBalance
				} else {
					fmt.Println("! > Amount can't be <= 0")
					userPage(email)
				}
			} else if currentIter.AccountHash == currentPerson.AccountHash && currentIter.CurrentBalance <= 0 {
				fmt.Printf("! > Your balance seems to be low by: %v\n", previousBalance-amount)
				userPage(email)
			}
		}
		fmt.Printf("$ > Previous Balance: %v\n$ > Withdrawn amount: %v\n$ > Current Balance: %v\n\n", previousBalance, amount, currentBalance)
		userPage(email)
	case 2:
		/**
		Logic for deposit money
		- Add the amount to current balance
		**/
		var amount, previousBalance, currentBalance float64
		fmt.Print(">> Enter amount to deposit: ")
		fmt.Scanln(&amount)

		for i := range BankDetails {
			currentIter := &BankDetails[i]
			if currentIter.AccountHash == currentPerson.AccountHash && currentIter.CurrentBalance >= 0 {
				if amount > 0 {
					previousBalance = currentIter.CurrentBalance
					currentIter.CurrentBalance += amount
					currentBalance = currentIter.CurrentBalance
					fmt.Printf("$ > Previous Balance: %v\n$ > Deposit amount: %v\n$ > Current Balance: %v\n\n", previousBalance, amount, currentBalance)
					userPage(email)
				} else {
					fmt.Println("! > Amount can't be <= 0 !")
					userPage(email)
				}
			}
		}
		BankDetails = append(BankDetails, BankDetail{AccountHash: currentPerson.AccountHash, CurrentBalance: amount})
	case 3:
		/**
		Logic for transferring money
		- Deduct the amount of the sender's current balance
		- Increse the amount of the reciever's current balance
		- The sending and the recieving amount can't be 0
		**/
		var amount float64
		var recieverEmail string

		fmt.Print(">> Enter the reciever's email and amount: ")
		fmt.Scanln(&recieverEmail, &amount)

		// A while loop which validates if the amount or the reciever email entered is non empty
		for amount <= 0 || recieverEmail == "" {
			fmt.Println("! > Amount can't be 0 and reciever email can't be empty")
			fmt.Print(">> Enter the reciever's email and amount: ")
			fmt.Scanln(&recieverEmail, &amount)
		}

		for i := range Persons {
			if Persons[i].User.Email == recieverEmail {
				reciever := &Persons[i]
				for i := range BankDetails {
					currentIter := &BankDetails[i]
					if currentIter.AccountHash == reciever.AccountHash {
						currentIter.CurrentBalance += amount
					} else if currentIter.AccountHash == currentPerson.AccountHash {
						currentIter.CurrentBalance -= amount
					}
				}
				TransferTransactions = append(TransferTransactions, TransferTransaction{Amount: amount, From: currentPerson.AccountHash, To: reciever.AccountHash})
				fmt.Printf("$ > Amount send: %v\n$ > To: %v\n\n", amount, reciever.AccountHash)
			}
		}

		userPage(email)
	case 4:
		fmt.Print("# > Logging out...\n\n")
		login()
	}
}

func main() {
	login()
}
