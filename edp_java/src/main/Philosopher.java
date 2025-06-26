package main;

import java.util.Optional;
import java.util.concurrent.CountDownLatch;

public class Philosopher extends Thread {

    // Constants (Time in milliseconds)
    private static final long THINKING_TIME = 1000;
    private static final long EATING_TIME = 2500;
    private static final long WAITING_TIME = 100;
    private static final long TIMEOUT = 2500;
    private static final float STARTING_MONEY = 200.0f;
    private static final float COUPON_VALUE = 5.0f;

    private static final Philosopher[] PHILOSOPHERS = java.util.Arrays.stream(Name.PHILOSOPHER_NAMES)
            .map(n -> new Philosopher(n.toString())).toArray(Philosopher[]::new);

    // Track which chopsticks are currently being held
    private static final java.util.Set<Integer> heldChopsticks = java.util.concurrent.ConcurrentHashMap.newKeySet();

    public static Philosopher[] list() {
        return PHILOSOPHERS;
    }

    public static java.util.Set<Integer> getHeldChopsticks() {
        return new java.util.HashSet<>(heldChopsticks);
    }

    private final Name name;
    private float money = STARTING_MONEY;
    private Optional<Seat> seat = Optional.empty();
    private Optional<Order> order = Optional.empty();
    private CountDownLatch foodServedLatch;
    private volatile boolean shouldReceiveCoupon = false;

    public Name name() {
        return name;
    }

    public Philosopher(String name) {
        this.name = new Name(name);
        Utility.validateTime(TIMEOUT);
    }

    public Optional<Order> order() {
        return order;
    }

    public void setOrder(Order order) {
        this.order = Optional.of(order);
    }

    public void clearOrder() {
        this.order = Optional.empty();
    }

    private void think() {
        try {
            Thread.sleep(THINKING_TIME);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    private void eat() {
        int seatNumber = seat.get().number();
        int leftChopstick = seatNumber;
        int rightChopstick = (seatNumber + 1) % PHILOSOPHERS.length;

        int firstChopstick;
        int secondChopstick;
        // Last philosopher picks up right chopstick first, then left
        if (seatNumber == PHILOSOPHERS.length - 1) {
            firstChopstick = rightChopstick;
            secondChopstick = leftChopstick;
        } else {
            firstChopstick = leftChopstick;
            secondChopstick = rightChopstick;
        }

        synchronized (EnhancedDiningPhilosophers.CHOPSTICKS[firstChopstick]) {
            heldChopsticks.add(firstChopstick);
            synchronized (EnhancedDiningPhilosophers.CHOPSTICKS[secondChopstick]) {
                heldChopsticks.add(secondChopstick);
                Utility.waitFor(EATING_TIME);
                heldChopsticks.remove(secondChopstick);
            }
            heldChopsticks.remove(firstChopstick);
        }
    }

    private void pay() {
        float mealCost = order.get().cost();
        money -= mealCost;
        if (money < 0) {
            System.out.println(String.format(
                    "Philosopher %s cannot afford the meal ($%.2f) and is leaving for good. Balance: $%.2f",
                    name, mealCost, money));
            money = 0;
        } else {
            System.out.println(
                    String.format("Philosopher %s has paid $%.2f and left the restaurant.", name, mealCost));
        }
        vacateSeat();
    }

    private boolean waitForWaiter() {
        System.out.println(String.format("Philosopher %s is waiting for a waiter.", name));

        // Add to the blocking queue
        Waiter.addPhilosopherToQueue(this);

        // Wait for a waiter to take the order within TIMEOUT
        long startTime = System.currentTimeMillis();
        while (System.currentTimeMillis() - startTime <= TIMEOUT) {
            // Check if we have an order (meaning a waiter took our order)
            if (order.isPresent()) {
                System.out.println(String.format("Philosopher %s got an order from waiter.", name));
                return true;
            }
            Utility.waitFor(WAITING_TIME);
        }

        // Timeout reached - remove from queue if still there
        Waiter.removePhilosopherFromQueue(this);
        System.out.println(String.format("Philosopher %s gave up waiting for a waiter.", name));
        return false;
    }

    public void notifyFoodServed() {
        if (foodServedLatch != null)
            foodServedLatch.countDown();
    }

    public void setShouldReceiveCoupon(boolean shouldReceive) {
        this.shouldReceiveCoupon = shouldReceive;
    }

    private void giveCoupon(float amount) {
        money += amount;
        System.out.println(String.format("Philosopher %s received a $%.2f coupon. New balance: $%.2f", name,
                amount, money));
    }

    @Override
    public void run() {
        while (money > 0) {
            attemptToDine();
            think();
        }
        System.out.println(String.format("Philosopher %s has left the restaurant for good.", name));
    }

    private void attemptToDine() {
        seat = Optional.ofNullable(Seat.availableSeat());
        if (seat.isPresent() && seat.get().attemptToOccupy()) {
            System.out.println(
                    String.format("Philosopher %s is being seated in chair %d.", name, seat.get().number()));
            think();
            System.out.println(String.format("Philosopher %s is about to call for a waiter.", name));
            boolean hasWaiter = waitForWaiter();
            if (!hasWaiter) {
                System.out.println(
                        String.format("Philosopher %s has left the restaurant without being served.", name));
                seat.get().vacate();
                return;
            }
            System.out.println(String.format("Philosopher %s got waiter, waiting for food.", name));
            shouldReceiveCoupon = false;
            foodServedLatch = new CountDownLatch(1);
            try {
                foodServedLatch.await();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
            if (shouldReceiveCoupon) {
                giveCoupon(COUPON_VALUE);
                return; // Leave without eating
            }
            System.out.println(String.format("Philosopher %s got food, about to eat.", name));
            eat();
            pay();
            clearOrder();
        } else {
            System.out.println(String.format("Philosopher %s could not get a seat.", name));
        }
    }

    public void vacateSeat() {
        if (seat.isPresent()) {
            seat.get().vacate();
            seat = Optional.empty(); // Clear the seat reference so philosopher can get a new seat
        }
    }
}
