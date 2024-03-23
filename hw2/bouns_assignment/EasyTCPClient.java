import java.io.*;
import java.net.*;
import java.util.*;

public class EasyTCPClient {

    public static void main(String[] args) {
        final String SERVER_ADDRESS = "localhost";
        final int PORT = 9999;

        try (
                Socket socket = new Socket(SERVER_ADDRESS, PORT);
                PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
                BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()))
        ) {
            Scanner sc = new Scanner(System.in);
            while (true) {
                String strType = sc.nextLine().replaceFirst("\n", "");;
                String message = "Message " + strType;
                out.println(message);
                System.out.println("Sent to server: " + message);
                String response = in.readLine();
                System.out.println("Server response: " + response);
                Thread.sleep(1000); // Optional delay
            }
        } catch (IOException | InterruptedException e) {
            e.printStackTrace();
        }
    }
}
