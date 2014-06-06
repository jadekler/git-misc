package main

import (
    "fmt"
    "os"
    "time"
    "math/rand"
)

func main() {
    orders := getOrders()
    fmt.Println("I'll let you know as soon as those orders are ready!")
    
    // Spawn goroutine for each order
    // (later) spawn goroutines capped at max amount of barristas
    completeOrders(orders)

    // Give customer their order

    // Calculate total cost and give customer receipt
    fmt.Println("Have a good day!")
    var input string
    fmt.Scanln(&input)
}

func completeOrders(orders map[string]int) () {
    for product, amt := range orders {
        for i := 0; i < amt; i++ {
            go func(product string) {
                amt := time.Duration(rand.Intn(250))
                time.Sleep(time.Millisecond * amt)
                fmt.Printf("%v completed!", product)
            }(product)
        }
    }
}

func getOrders() (map[string]int) {
    var input = ""

    orders := map[string]int{
        "coffees": 0,
        "teas": 0,
        "hotchocs": 0,
    }

    fmt.Println("Welcome! What can I get you? We have coffee (c), hot chocolate (h), and tea (t). If you'd like to quit, press (q).")

    for input != "d" && input != "done" {
        _, err := fmt.Scanf("%s", &input)
        if (err != nil) {
            
        }

        switch input {
            case "c", "coffee":
                fmt.Println("One coffee, got it.")
                orders["coffees"]++
            case "h", "hot chocolate":
                fmt.Println("One hot chocolate, got it.")
                orders["hotchocs"]++
            case "t", "tea":
                fmt.Println("One tea, got it.")
                orders["teas"]++
            case "q", "quit":
                os.Exit(1)
            case "d", "done":
            default:
                fmt.Println("Huh?")
        }

        if input != "d" && input != "done" {
            fmt.Println("What else can I get you?  We have coffee (c), hot chocolate (h), and tea (t). Or, are you done (d)?")
        }
    }

    fmt.Println("That's it? Ok!")

    return orders
}

// Orders amountOfCoffees coffees and returns the price
// when the coffees are ready
func OrderCoffee(amountOfCoffees int) (price int) {
    return amountOfCoffees*3
}