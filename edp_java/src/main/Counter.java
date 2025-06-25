package main;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;

public class Counter {
    private static final BlockingQueue<Order> orders = new LinkedBlockingQueue<>();
    private static final BlockingQueue<Order> completedMeals = new LinkedBlockingQueue<>();

    private Counter() {
    }

    public static void placeOrder(Order order) {
        try {
            orders.put(order);
            System.out.println(String.format("Order for %s placed on counter.", order.philosopher().name()));
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    public static Order takeOrder() {
        try {
            Order order = orders.take();
            System.out
                    .println(String.format("Order for %s picked up from counter.", order.philosopher().name()));
            return order;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            return null;
        }
    }

    public static boolean hasOrders() {
        return !orders.isEmpty();
    }

    public static int orderCount() {
        return orders.size();
    }

    public static void placeCompletedMeal(Order order) {
        try {
            completedMeals.put(order);
            System.out.println(String.format("Order for %s placed on counter.", order.philosopher().name()));
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    public static boolean hasCompletedMeals() {
        return !completedMeals.isEmpty();
    }

    public static Order pollCompletedMeal() {
        Order order = completedMeals.poll();
        if (order != null) {
            System.out
                    .println(String.format("Order for %s picked up from counter.", order.philosopher().name()));
        }
        return order;
    }
}
