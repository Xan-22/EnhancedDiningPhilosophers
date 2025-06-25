package main;

public class Utility {
    private Utility() {
        // Prevent instantiation
    }

    public static void validatePrice(float price) {
        if (price <= 0) {
            throw new IllegalArgumentException("Food price must be positive.");
        }
    }

    public static void validateTime(long time) {
        if (time < 0) {
            throw new IllegalArgumentException("Time cannot be negative.");
        }
    }

    public static void waitFor(long milliseconds) {
        if (milliseconds < 0) {
            throw new IllegalArgumentException("Wait time cannot be negative.");
        }
        try {
            Thread.sleep(milliseconds);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt(); // Restore interrupted status
        }
    }

}
