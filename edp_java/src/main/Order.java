package main;

public class Order {
    private final Philosopher philosopher;
    private final Food[] meal;
    private final float cost;

    public Order(Philosopher philosopher) {
        this.philosopher = philosopher;
        this.meal = Food.randomMeal();
        this.cost = calculateCost();
    }

    private float calculateCost() {
        float total = 0;
        for (Food food : meal)
            total += food.price();
        return total;
    }

    public Philosopher philosopher() {
        return philosopher;
    }

    public float cost() {
        return cost;
    }

    public Food[] meal() {
        return meal;
    }

    public String mealString() {
        return java.util.Arrays.stream(meal())
                .filter(f -> !(f.type() == Food.Type.SOUP && f.name().equals("No Soup")))
                .filter(f -> !(f.type() == Food.Type.DESSERT && f.name().equals("No Dessert")))
                .map(Food::toString)
                .reduce((a, b) -> a + " and " + b)
                .orElse("");
    }

    public static class Food {
        public enum Type {
            ENTREE, SOUP, DESSERT
        }

        private final Type type;
        private final Name name;
        private final float price;
        private static final java.util.Random rand = new java.util.Random();

        public Food(Type type, String name, float price) {
            this.type = type;
            this.name = new Name(name);
            if (price < 0)
                throw new IllegalArgumentException("Food price cannot be negative.");
            this.price = price;
        }

        public String name() {
            return name.toString();
        }

        public float price() {
            return price;
        }

        public Type type() {
            return type;
        }

        @Override
        public String toString() {
            return name.toString();
        }

        private static final Food[] entrees = new Food[] {
                new Food(Type.ENTREE, "Paella", 13.25f),
                new Food(Type.ENTREE, "Wu Hsiang Chi", 10.00f),
                new Food(Type.ENTREE, "Bogrács Gulyás", 11.25f),
                new Food(Type.ENTREE, "Spanakopita", 6.50f),
                new Food(Type.ENTREE, "Moui Nagden", 12.95f),
                new Food(Type.ENTREE, "Sambal Goreng Udang", 14.95f)
        };
        private static final Food[] soups = new Food[] {
                new Food(Type.SOUP, "No Soup", 0.00f),
                new Food(Type.SOUP, "Albóndigas", 3.00f)
        };
        private static final Food[] desserts = new Food[] {
                new Food(Type.DESSERT, "No Dessert", 0.00f),
                new Food(Type.DESSERT, "Berog", 3.50f)
        };

        public static Food[] menu() {
            return entrees.clone();
        }

        public static Food[] soups() {
            return soups.clone();
        }

        public static Food[] desserts() {
            return desserts.clone();
        }

        public static Food[] randomMeal() {
            return new Food[] {
                    entrees[rand.nextInt(entrees.length)],
                    soups[rand.nextInt(soups.length)],
                    desserts[rand.nextInt(desserts.length)]
            };
        }
    }
}
