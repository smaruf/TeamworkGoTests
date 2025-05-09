package customerimporter;

import java.util.*;
import java.util.regex.Pattern;

public class DomainCounter {
    private static final Pattern EMAIL_PATTERN = Pattern.compile("^[^@]+@([^@]+)$");

    public static Map<String, Integer> countDomains(List<String> emails) {
        Map<String, Integer> domainCounts = new HashMap<>();

        for (String email : emails) {
            var matcher = EMAIL_PATTERN.matcher(email);
            if (matcher.matches()) {
                String domain = matcher.group(1);
                domainCounts.put(domain, domainCounts.getOrDefault(domain, 0) + 1);
            }
        }

        return domainCounts;
    }

    public static List<String> sortDomains(Map<String, Integer> domainCounts) {
        List<String> sortedDomains = new ArrayList<>();
        domainCounts.entrySet().stream()
                .sorted(Map.Entry.comparingByKey()) // Sort alphabetically by domain
                .forEach(entry -> sortedDomains.add(entry.getKey() + ": " + entry.getValue()));
        return sortedDomains;
    }
}