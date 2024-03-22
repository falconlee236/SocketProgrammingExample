import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.Socket;
import java.net.UnknownHostException;
import java.util.Scanner;

public class EasyTCPClient {

    public static void main(String[] args) {
        final String serverName = "nsl2.cau.ac.kr";
        final int serverPort = 20532;

        try (Socket socket = new Socket(serverName, serverPort);
             PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
             BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()));
             Scanner scanner = new Scanner(System.in)) {

            Runtime.getRuntime().addShutdownHook(new Thread(() -> {
                System.out.println("\nBye bye~");
                System.exit(1);
            }));

            System.out.println("Client is running on port " + socket.getLocalPort());

            while (true) {
                System.out.println("<Menu>");
                System.out.println("1) convert text to UPPER-case");
                System.out.println("2) get my IP address and port number");
                System.out.println("3) get server request count");
                System.out.println("4) get server running time");
                System.out.println("5) exit");
                System.out.print("Input option: ");
                String inputOption = scanner.nextLine().trim();

                if (!inputOption.matches("[1-5]")) {
                    System.out.println("Invalid option");
                    continue;
                }

                long start = System.currentTimeMillis();
                out.println(inputOption);

                if (inputOption.equals("1")) {
                    System.out.print("Input lowercase sentence: ");
                    String input = scanner.nextLine();
                    out.println(input);
                }

                String response = in.readLine();
                long duration = System.currentTimeMillis() - start;

                switch (inputOption) {
                    case "1":
                        System.out.println("\nReply from server: " + response);
                        break;
                    case "2":
                        String[] recInfo = response.split(":");
                        System.out.printf("\nReply from server: client IP = %s PORT = %s\n", recInfo[0], recInfo[1]);
                        break;
                    case "3":
                        System.out.println("\nReply from server: requests served = " + response);
                        break;
                    case "4":
                        System.out.println("\nReply from server: run time = " + response);
                        break;
                    case "5":
                        System.out.println("Bye bye~");
                        return;
                }
                System.out.printf("RTT = %dms\n", duration);
            }
        } catch (UnknownHostException e) {
            System.err.println("Unknown host: " + serverName);
            e.printStackTrace();
        } catch (IOException e) {
            System.err.println("I/O error occurred while connecting to server");
            e.printStackTrace();
        }
    }
}
