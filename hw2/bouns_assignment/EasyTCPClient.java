import java.io.*;
import java.net.*;
import java.util.*;

public class EasyTCPClient {
    public static void main(String[] args) {
        final String SERVER_ADDRESS = "localhost";
        final int PORT = 30532;

        try (
                Socket socket = new Socket(SERVER_ADDRESS, PORT);
                PrintWriter out = new PrintWriter(
                        socket.getOutputStream(), true);
                BufferedReader in = new BufferedReader(
                        new InputStreamReader(socket.getInputStream()))
        ) {
            Scanner sc = new Scanner(System.in);
            System.out.printf("The client is running on port %d\n", socket.getLocalPort());
            while (true) {

                System.out.println("<Menu>");
                System.out.println("1) convert text to UPPER-case");
                System.out.println("2) get my IP address and port number");
                System.out.println("3) get server request count");
                System.out.println("4) get server running time");
                System.out.println("5) exit");
                System.out.println("Input option: ");
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
