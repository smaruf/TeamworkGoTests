package customerimporter;

import org.junit.jupiter.api.Test;

import java.io.*;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class CSVProcessorTest {
    @Test
    void testReadCSV() throws IOException {
        String testCSV = "FirstName,LastName,Email\nJohn,Doe,john.doe@example.com\nJane,Smith,jane.smith@example.com";
        File tempFile = File.createTempFile("test", ".csv");
        try (BufferedWriter writer = new BufferedWriter(new FileWriter(tempFile))) {
            writer.write(testCSV);
        }

        List<String> emails = CSVProcessor.readCSV(tempFile.getAbsolutePath());
        assertEquals(2, emails.size());
        assertTrue(emails.contains("john.doe@example.com"));
        assertTrue(emails.contains("jane.smith@example.com"));

        tempFile.delete();
    }
}