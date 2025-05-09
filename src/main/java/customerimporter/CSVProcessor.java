package customerimporter;

import java.io.*;
import java.util.ArrayList;
import java.util.List;

public class CSVProcessor {
    public static List<String> readCSV(String filePath) throws IOException {
        List<String> emails = new ArrayList<>();
        try (BufferedReader br = new BufferedReader(new FileReader(filePath))) {
            String line;
            boolean isHeader = true;
            while ((line = br.readLine()) != null) {
                if (isHeader) {
                    isHeader = false; // Skip the header row
                    continue;
                }
                String[] fields = line.split(",");
                if (fields.length > 2) {
                    emails.add(fields[2].trim()); // Assuming email is the 3rd column
                }
            }
        }
        return emails;
    }
}