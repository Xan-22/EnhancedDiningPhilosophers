package main;

public class EnhancedDiningPhilosophers {

    static final Object[] CHOPSTICKS;
    static {
        Object[] chopsticks = new Object[Philosopher.list().length];
        java.util.Arrays.setAll(chopsticks, i -> new Object());
        CHOPSTICKS = chopsticks;
    }

    public static void main(String[] args) {
        System.out.println("The restaurant is now open for business.");

        for (Cook cook : Cook.list()) {
            cook.setDaemon(true);
            cook.start();
        }

        for (Waiter waiter : Waiter.list()) {
            waiter.setDaemon(true);
            waiter.start();
        }

        for (Philosopher philosopher : Philosopher.list()) {
            philosopher.start();
        }

        monitorRestaurant();
    }

    private static void monitorRestaurant() {
        boolean done = false;
        while (!done) {
            try {
                Thread.sleep(1000);

                if (allPhilosophersLeft()) {
                    System.out.println("The restaurant has closed down.");
                    done = true;
                } else {
                    logRestaurantStatus();
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                done = true;
            }
        }
    }

    private static boolean allPhilosophersLeft() {
        for (Philosopher philosopher : Philosopher.list()) {
            if (philosopher.isAlive()) {
                return false;
            }
        }
        return true;
    }

    private static void logRestaurantStatus() {
        int activePhilosophers = 0;
        for (Philosopher philosopher : Philosopher.list()) {
            if (philosopher.isAlive()) {
                activePhilosophers++;
            }
        }

        int occupiedSeats = 0;
        java.util.List<Integer> occupiedSeatNumbers = new java.util.ArrayList<>();
        for (Seat seat : Seat.seats()) {
            if (seat.isOccupied()) {
                occupiedSeats++;
                occupiedSeatNumbers.add(seat.number());
            }
        }

        java.util.Set<Integer> heldChopsticks = Philosopher.getHeldChopsticks();
        java.util.List<Integer> chopstickNumbers = new java.util.ArrayList<>(heldChopsticks);
        java.util.Collections.sort(chopstickNumbers);

        if (activePhilosophers > 0) {
            String seatInfo = occupiedSeats + " seats taken (" +
                    occupiedSeatNumbers.stream()
                            .map(String::valueOf)
                            .collect(java.util.stream.Collectors.joining(", "))
                    + ")";

            String chopstickInfo = heldChopsticks.size() + " chopsticks taken (" +
                    chopstickNumbers.stream()
                            .map(String::valueOf)
                            .collect(java.util.stream.Collectors.joining(", "))
                    + ")";

            System.out.println(
                    String.format("%nRestaurant status: %d active philosophers, %s, %s, %d orders on counter%n",
                            activePhilosophers, seatInfo, chopstickInfo, Counter.orderCount()));
        }
    }
}