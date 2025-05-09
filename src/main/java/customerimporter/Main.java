package customerimporter;

import java.io.IOException;

public class Main {
    public static void main(String[] args) {
        if (args.length < 2) {
            System.out.println("Usage: java Main <inputFile> <outputFile>");
            return;
        }

        String inputFile = args[0];
        String outputFile = args[1];

        try {
            // Process the CSV file and count email domains
            var records = CSVProcessor.readCSV(inputFile);
            var domainCounts = DomainCounter.countDomains(records);

            // Write the sorted domains to the output file
            FileWriterUtil.writeOutput(domainCounts, outputFile);

            System.out.println("Processing completed. Results written to " + outputFile);
        } catch (IOException e) {
            System.err.println("Error: " + e.getMessage());
        }
    }
}