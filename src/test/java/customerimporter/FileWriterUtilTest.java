package customerimporter;

import org.junit.jupiter.api.Test;

import java.io.*;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class FileWriterUtilTest {
    @Test
    void testWriteOutput() throws IOException {
        Map<String, Integer> domainCounts = Map.of("example.com", 2, "another.com", 1);
        File tempFile = File.createTempFile("output", ".txt");

        FileWriterUtil.writeOutput(domainCounts, tempFile.getAbsolutePath());

        try (BufferedReader reader = new BufferedReader(new FileReader(tempFile))) {
            assertEquals("another.com: 1", reader.readLine());
            assertEquals("example.com: 2", reader.readLine());
        }

        tempFile.delete();
    }
}