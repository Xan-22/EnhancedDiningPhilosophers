package main;

import java.util.Optional;
import java.util.concurrent.Semaphore;

public class Waiter extends Thread {
    // Constants (Time in milliseconds)
    private static final int CHECK_ORDERS_INTERVAL = 500;

    // Static
    private static final Waiter[] waiters = java.util.Arrays.stream(Name.WAITER_NAMES)
            .map(n -> new Waiter(n.toString())).toArray(Waiter[]::new);

    private static final int NUM_COOKS = Cook.list().length;
    public static final Semaphore cookSemaphore = new Semaphore(NUM_COOKS, true);

    // Synchronized queue for philosophers calling waiters
    private static final java.util.concurrent.BlockingQueue<Philosopher> WAITER_CALLING_QUEUE = new java.util.concurrent.LinkedBlockingQueue<>();

    public static Waiter[] list() {
        return waiters;
    }

    public static final Waiter availableWaiter() {
        for (Waiter waiter : waiters) {
            if (!waiter.isProcessingOrder()) {
                return waiter;
            }
        }
        return null; // No available waiter found
    }

    public static void addPhilosopherToQueue(Philosopher philosopher) {
        try {
            WAITER_CALLING_QUEUE.put(philosopher);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    public static void removePhilosopherFromQueue(Philosopher philosopher) {
        WAITER_CALLING_QUEUE.remove(philosopher);
    }

    // Instance
    private Name name;
    private Optional<Order> order = Optional.empty();
    private Optional<Philosopher> patron = Optional.empty();
    private boolean isProcessingOrder = false;

    public Waiter(String name) {
        this.name = new Name(name);
    }

    public String name() {
        return name.toString();
    }

    private boolean isProcessingOrder() {
        return isProcessingOrder;
    }

    public void takeOrder(Philosopher philosopher, Order order) {
        this.isProcessingOrder = true;
        this.order = Optional.of(order);
        this.patron = Optional.of(philosopher);
    }

    @Override
    public void run() {
        System.out.println(String.format("Waiter %s is ready to take orders.", name));
        while (!Thread.currentThread().isInterrupted()) {
            try {
                if (isProcessingOrder && order.isPresent() && patron.isPresent()) {
                    processOrder();
                } else {
                    // Check for completed meals to deliver
                    if (Counter.hasCompletedMeals()) {
                        Order completedOrder = Counter.pollCompletedMeal();
                        if (completedOrder != null) {
                            deliverOrder(completedOrder);
                        }
                    } else {
                        // No completed meals available, check for philosophers who need service
                        checkForPhilosophersNeedingService();
                    }
                }

                Thread.sleep(CHECK_ORDERS_INTERVAL);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break;
            }
        }
    }

    private void processOrder() {
        if (!order.isPresent() || !patron.isPresent()) {
            return;
        }
        Order currentOrder = order.get();
        Philosopher currentPatron = patron.get();
        boolean cookAcquired = cookSemaphore.tryAcquire();
        if (cookAcquired) {
            boolean orderPlaced = false;
            try {
                Cook.ORDER_QUEUE.put(currentOrder);
                orderPlaced = true;
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
            if (orderPlaced) {
                System.out.println(
                        String.format("Waiter %s placed order for %s.", name, currentPatron.name()));
            }
        } else {
            System.out
                    .println(String.format("Waiter %s cannot place order for %s - all chefs busy. Giving $5.00 coupon.",
                            name, currentPatron.name()));
            currentPatron.setShouldReceiveCoupon(true);
            currentPatron.notifyFoodServed(); // Countdown the latch to wake up the philosopher
            // Clear the philosopher's order so they can try again
            currentPatron.clearOrder();
            // Vacate the seat and print message
            currentPatron.vacateSeat();
            System.out.println(String.format("Philosopher %s has left the restaurant without being served.",
                    currentPatron.name()));
        }
        order = Optional.empty();
        patron = Optional.empty();
        isProcessingOrder = false;
    }

    private void deliverOrder(Order completedOrder) {
        Philosopher targetPhilosopher = completedOrder.philosopher();
        targetPhilosopher.notifyFoodServed();
        System.out.println(String.format("Waiter %s is serving philosopher %s %s.", name,
                targetPhilosopher.name(), completedOrder.mealString()));
        // Make waiter available for new orders
        order = Optional.empty();
        patron = Optional.empty();
        isProcessingOrder = false;
    }

    private void checkForPhilosophersNeedingService() {
        if (isProcessingOrder) {
            return; // Already processing an order
        }

        // Check the synchronized queue for philosophers calling waiters
        Philosopher callingPhilosopher = WAITER_CALLING_QUEUE.poll();

        if (callingPhilosopher != null) {
            // Philosopher needs service
            Order newOrder = new Order(callingPhilosopher);
            callingPhilosopher.setOrder(newOrder);
            takeOrder(callingPhilosopher, newOrder);

            String mealDesc = newOrder.mealString();
            System.out.println(String.format("Waiter %s has taken order for %s from philosopher %s.",
                    name, mealDesc, callingPhilosopher.name()));
        }
    }

    // class implementation
    public void serve(Philosopher philosopher) {
        if (order.isPresent()) {
            System.out.println(String.format("Waiter %s is serving philosopher %s %s.", name,
                    philosopher.name(), order.get().mealString()));
        }
    }
}
