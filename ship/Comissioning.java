package ship;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.PublicKey;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;

public class Comissioning {

    public static String[] readSkis() {
        // Create file if not exists
        Path skisPath = Paths.get("skis.txt");
        if (Files.notExists(skisPath)) {
            return new String[0];
        }

        try (BufferedReader reader = new BufferedReader(new FileReader("skis.txt"))) {
            return reader.lines()
                    .map(line -> line.split(","))
                    .filter(line -> line.length > 1)
                    .map(line -> line[0])
                    .toArray(String[]::new);
        } catch (IOException e) {
            e.printStackTrace();
            return new String[0];
        }
    }

    public static void writeSkis(String[] newSkis, String[] newDevices) {
        try (BufferedWriter writer = new BufferedWriter(new FileWriter("skis.txt"))) {
            for (int i = 0; i < newSkis.length; i++) {
                if (!newSkis[i].isEmpty()) {
                    writer.write(newSkis[i] + "," + newDevices[i]);
                    if (i < newSkis.length - 1) {
                        writer.newLine();
                    }
                }
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public String getSki(String certName) {
        try {
            byte[] fileBytes = Files.readAllBytes(Paths.get(certName + ".crt"));
            String crt = new String(fileBytes);

            CertificateFactory certFactory = CertificateFactory.getInstance("X.509");
            X509Certificate x509Certificate = (X509Certificate) certFactory.generateCertificate(
                    new java.io.ByteArrayInputStream(crt.getBytes()));

            PublicKey publicKey = x509Certificate.getPublicKey();
            byte[] publicKeyBytes = publicKey.getEncoded();

            MessageDigest sha1 = MessageDigest.getInstance("SHA-1");
            byte[] digest = sha1.digest(publicKeyBytes);

            StringBuilder stringBuilder = new StringBuilder();
            for (byte b : digest) {
                stringBuilder.append(String.format("%02X", b));
            }

            return stringBuilder.toString();
        } catch (IOException | CertificateException | NoSuchAlgorithmException e) {
            e.printStackTrace();
            return "";
        }
    }
}
