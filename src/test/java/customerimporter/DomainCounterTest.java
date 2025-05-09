package customerimporter;

import org.junit.jupiter.api.Test;

import java.util.*;

import static org.junit.jupiter.api.Assertions.*;

class DomainCounterTest {
    @Test
    void testCountDomains() {
        List<String> emails = List.of("john.doe@example.com", "jane.smith@example.com", "invalid-email");
        Map<String, Integer> domainCounts = DomainCounter.countDomains(emails);

        assertEquals(1, domainCounts.get("example.com"));
        assertFalse(domainCounts.containsKey("invalid-email"));
    }

    @Test
    void testSortDomains() {
        Map<String, Integer> domainCounts = Map.of("example.com", 2, "another.com", 1);
        List<String> sortedDomains = DomainCounter.sortDomains(domainCounts);

        assertEquals("another.com: 1", sortedDomains.get(0));
        assertEquals("example.com: 2", sortedDomains.get(1));
    }
}