package customerimporter;

import java.io.*;
import java.util.List;

public class FileWriterUtil {
    public static void writeOutput(Map<String, Integer> domainCounts, String outputFile) throws IOException {
        List<String> sortedDomains = DomainCounter.sortDomains(domainCounts);

        try (BufferedWriter writer = new BufferedWriter(new FileWriter(outputFile))) {
            for (String domain : sortedDomains) {
                writer.write(domain);
                writer.newLine();
            }
        }
    }
}