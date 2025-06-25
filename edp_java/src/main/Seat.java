package main;

public class Seat {
    private static final Seat[] SEATS = new Seat[Philosopher.list().length];
    static {
        for (int i = 0; i < SEATS.length; i++)
            SEATS[i] = new Seat();
    }

    public static Seat[] seats() {
        return SEATS.clone();
    }

    public static synchronized Seat availableSeat() {
        for (Seat seat : SEATS)
            if (!seat.occupied)
                return seat;
        return null;
    }

    private static int nextNumber = 0;
    private final int number;
    private boolean occupied = false;

    public Seat() {
        this.number = nextNumber++;
    }

    public int number() {
        return number;
    }

    public synchronized boolean isOccupied() {
        return occupied;
    }

    public synchronized boolean attemptToOccupy() {
        if (!occupied) {
            occupy();
            return true;
        }
        return false;
    }

    private synchronized void occupy() {
        occupied = true;
    }

    public synchronized void vacate() {
        System.out.println(String.format("Seat %d is vacated.", number));
        occupied = false;
    }
}
