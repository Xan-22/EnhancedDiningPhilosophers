package main;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.Semaphore;

public class Cook extends Thread {

    // Constants (Time in milliseconds)
    private static final int COOKING_TIME = 6000;
    private static final int COFFEE_BREAK_TIME = 3000;

    // Static
    private static final Cook[] CHEFS = java.util.Arrays.stream(Name.COOK_NAMES)
            .map(n -> new Cook(n.toString())).toArray(Cook[]::new);

    public static Cook[] list() {
        return CHEFS;
    }

    protected static final BlockingQueue<Order> ORDER_QUEUE = new java.util.concurrent.ArrayBlockingQueue<>(1);

    private static final Semaphore cookSemaphore = Waiter.cookSemaphore;

    // Instance
    private final String name;
    private int mealsPrepared = 0;
    private boolean isOnBreak = false;

    public Cook(String name) {
        this.name = name;
    }

    public String name() {
        return name;
    }

    public boolean isOnBreak() {
        return isOnBreak;
    }

    @Override
    public void run() {
        while (!Thread.currentThread().isInterrupted()) {
            try {
                System.out.println(String.format("Chef %s is waiting for an order.", name));
                Order order = ORDER_QUEUE.take();
                if (order != null) {
                    cook(order);
                    mealsPrepared++;
                    if (mealsPrepared % 4 == 0) {
                        takeCoffeeBreak();
                    }
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break;
            }
        }
    }

    private void cook(Order order) throws InterruptedException {
        System.out.println(String.format("Chef %s is cooking the %s for Philosopher %s.", name, order.mealString(),
                order.philosopher().name()));
        Thread.sleep(COOKING_TIME);
        Counter.placeCompletedMeal(order);
        System.out.println(String.format("Chef %s has finished cooking the %s for Philosopher %s.", name,
                order.mealString(), order.philosopher().name()));
        cookSemaphore.release();
    }

    private void takeCoffeeBreak() throws InterruptedException {
        isOnBreak = true;
        System.out.println(String.format("Chef %s has returned from a coffee break.", name));
        Thread.sleep(COFFEE_BREAK_TIME);
        isOnBreak = false;
    }
}