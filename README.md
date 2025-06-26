# Enhanced Dining Philosophers

**Michail Mavromatis**

A multi-language implementation of the classic Dining Philosophers problem with additional features including a restaurant simulation with waiters, cooks, and a payment system.

This project is based on [The Enhanced Dining Philosophers Problem](https://cs.lmu.edu/~ray/notes/edp/) by Ray Toal.

## Problem Description

The Dining Philosophers problem is a classic synchronization problem that illustrates challenges in resource allocation and deadlock prevention. In this enhanced version, we simulate a restaurant where:

- **Philosophers** represent customers who need to eat
- **Chopsticks** represent shared resources (forks in the original problem)
- **Seats** represent dining positions at a circular table
- **Waiters** handle customer orders and food delivery
- **Cooks** prepare meals in a kitchen
- **Payment System** tracks customer spending and restaurant revenue

### Key Features

- **Deadlock Prevention**: Implements resource ordering to prevent deadlocks
- **Concurrent Operations**: Multiple philosophers, waiters, and cooks operate simultaneously
- **Resource Management**: Seats, chopsticks, and kitchen capacity are managed efficiently
- **Timeout Handling**: Philosophers can leave if service takes too long
- **Financial Tracking**: Customers have budgets and pay for meals
- **Coupon System**: Compensation for service delays

## Technologies Used

### Java Version (`edp_java/`)
- **Language**: Java 8+
- **Concurrency**: `synchronized` blocks, `wait()`/`notify()` mechanisms
- **Threading**: `Thread` class and `Runnable` interface
- **Collections**: `ArrayList`, `HashMap`, `Queue`
- **Build**: Standard Java compilation with `javac`

### Go Version (`edp_go/`)
- **Language**: Go 1.19+
- **Concurrency**: Goroutines, channels, `sync.Mutex`, `sync.WaitGroup`
- **Context**: `context.Context` for cancellation and timeouts
- **Modules**: Go modules for dependency management
- **Build**: Standard Go toolchain with `go build`

## Project Structure

```
EnhancedDiningPhilosophers/
├── edp_java/                 # Java implementation
│   ├── src/main/            # Source code
│   │   ├── Cook.java        # Kitchen staff
│   │   ├── Counter.java     # Order management
│   │   ├── EnhancedDiningPhilosophers.java  # Main program
│   │   ├── Name.java        # Name utilities
│   │   ├── Order.java       # Order management
│   │   ├── Philosopher.java # Customer simulation
│   │   ├── Seat.java        # Dining seat management
│   │   ├── Utility.java     # Utility functions
│   │   └── Waiter.java      # Service staff
│   └── bin/                 # Compiled classes
├── edp_go/                  # Go implementation
│   ├── types.go             # Type definitions and constants
│   ├── order.go             # Order management
│   ├── seat.go              # Seat management
│   ├── counter.go           # Order counter
│   ├── philosopher.go       # Philosopher simulation
│   ├── waiter.go            # Waiter management
│   ├── cook.go              # Cook management
│   ├── EnhancedDiningPhilosophers.go  # Main program
│   └── go.mod               # Go module file
├── LICENSE                  # Project license
└── README.md               # This file
```

## Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/EnhancedDiningPhilosophers.git
cd EnhancedDiningPhilosophers
```

### 2. Running the Java Version

#### Compile the Java Code
```bash
cd edp_java
javac -d bin src/main/*.java
```

#### Run the Java Program
```bash
java -cp bin EnhancedDiningPhilosophers
```

### 3. Running the Go Version

#### Navigate to Go Directory
```bash
cd edp_go
```

#### Build the Go Program
```bash
go build -o EnhancedDiningPhilosophers.exe .
```

#### Run the Go Program
```bash
./EnhancedDiningPhilosophers.exe
```

**Note**: On Unix-like systems (Linux/macOS), use:
```bash
go build -o EnhancedDiningPhilosophers .
./EnhancedDiningPhilosophers
```

## Program Behavior

When you run either version, you'll see output similar to:

```
The restaurant is now open for business.
Waiter Alice is ready to take orders.
Waiter Bob is ready to take orders.
Cook Charlie is ready to cook.
Cook David is ready to cook.
Philosopher Aristotle is being seated in chair 0.
Philosopher Plato is being seated in chair 1.
Philosopher Socrates is being seated in chair 2.
Philosopher Aristotle is about to call for a waiter.
Philosopher Plato is about to call for a waiter.
Philosopher Socrates is about to call for a waiter.
Waiter Alice has taken order for Pizza from philosopher Aristotle.
Waiter Bob has taken order for Burger from philosopher Plato.
Waiter Alice placed order for Aristotle.
Cook Charlie is cooking Pizza for Aristotle.
Philosopher Aristotle got waiter, waiting for food.
Waiter Bob placed order for Plato.
Cook David is cooking Burger for Plato.
Philosopher Plato got waiter, waiting for food.
Cook Charlie finished cooking Pizza for Aristotle.
Waiter Alice is serving philosopher Aristotle Pizza.
Philosopher Aristotle got food, about to eat.
Philosopher Aristotle has paid $12.50 and left the restaurant.
...
Restaurant status: 2 active philosophers, 2 seats taken (1, 2), 2 chopsticks taken (1, 2), 1 orders on counter
...
The restaurant has closed down.
```

## Key Components

### Philosophers (Customers)
- Start with $200.00 budget
- Think, then attempt to dine
- Wait for available seats
- Order food through waiters
- Eat using two chopsticks
- Pay for meals and leave

### Waiters (Service Staff)
- Take customer orders
- Manage order queue
- Deliver completed meals
- Handle service timeouts with coupons

### Cooks (Kitchen Staff)
- Prepare meals from order queue
- Limited kitchen capacity
- Notify when meals are ready

### Seats and Chopsticks
- Finite number of dining seats
- Chopsticks shared between adjacent philosophers
- Resource ordering prevents deadlocks

## Configuration

Both implementations use similar configuration constants that can be modified:

- **Philosopher names and starting money**
- **Waiter and cook names**
- **Food types and prices**
- **Timing parameters** (thinking, eating, waiting times)
- **Restaurant capacity** (seats, chopsticks)

### Performance Notes

- The Go version may run faster due to lightweight goroutines
- Both versions implement the same logic and should produce similar output
- The simulation runs until all philosophers leave the restaurant

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
