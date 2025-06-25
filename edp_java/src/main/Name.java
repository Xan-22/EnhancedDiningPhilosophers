package main;

public class Name {
    protected static final Name[] PHILOSOPHER_NAMES = new Name[] {
            new Name("Susan Haack"),
            new Name("Zhaozhou"),
            new Name("David Hume"),
            new Name("Omar Khayyám"),
            new Name("Kaṇāda")
    };

    protected static final Name[] COOK_NAMES = new Name[] {
            new Name("Eren"),
            new Name("Mikasa"),
            new Name("Armin")
    };

    protected static final Name[] WAITER_NAMES = new Name[] {
            new Name("Miria"),
            new Name("Isaac")
    };

    private final String value;

    public Name(String value) {
        validateName(value);
        this.value = value;
    }

    private static void validateName(String name) {
        if (name == null || name.trim().isEmpty()) {
            throw new IllegalArgumentException("Name cannot be null or empty.");
        }
    }

    @Override
    public String toString() {
        return value;
    }
}
