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
                System.out.print("Input option: ");
                String strType = sc.nextLine().replaceFirst("\n", "");;
                long startTime = System.nanoTime();
                out.println(strType);

                if (strType.equals("1")){
                    System.out.print("Input sentence: ");
                    String text = sc.nextLine().replaceFirst("\n", "");
                    startTime = System.nanoTime();
                    out.println(text.toUpperCase());
                } else if (strType.equals("5")){
                    System.out.println("Bye bye~");
                    throw new InterruptedException();
                }
                String response = in.readLine();
                long endTime = System.nanoTime();
                System.out.println("Reply from server: " + response);
                System.out.printf("RTT = %fms\n", (endTime - startTime) / 1e+6);
                Thread.sleep(1000); // Optional delay
            }
        } catch (IOException | InterruptedException e) {
            System.out.println(e.getClass().getName());
            if (e.getClass().getName().equals("java.net.SocketException")){
                System.out.println("Server Disconnected");
            }
        }
    }
}
