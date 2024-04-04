/*
EasyTCPClient.java
20190532 Sangyun Lee
*/
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.Socket;
import java.util.Scanner;

public class EasyTCPClient {
    public static void main(String[] args) {
        // set server info
        final String SERVER_ADDRESS = "localhost";
        final int PORT = 20532;

        // set SIGINT handler
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("\nBye bye~");
        }));

        try ( // create TCP Socket
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
                // get user input
                System.out.print("Input option: ");
                String strType = sc.nextLine().replaceFirst("\n", "");;
                try{
                    if (Integer.parseInt(strType) < 1 || Integer.parseInt(strType) > 5){
                        System.out.println("invalid argument");
                        continue;
                    }
                } catch (NumberFormatException e){
                    System.out.println("invalid argument");
                    continue;
                }
                long startTime = System.nanoTime();
                out.println(strType);

                if (strType.equals("1")){
                    System.out.print("Input sentence: ");
                    String text = sc.nextLine().replaceFirst("\n", "");
                    startTime = System.nanoTime();
                    // send to Server
                    out.println(text);
                } else if (strType.equals("5")){
                    System.out.println("Bye bye~");
                    throw new InterruptedException();
                }
                // get from Server
                String response = in.readLine();
                long endTime = System.nanoTime();
                System.out.println("Reply from server: " + response);
                System.out.printf("RTT = %fms\n", (endTime - startTime) / 1e+6);
                Thread.sleep(1000); // Optional delay
            }
        } catch (IOException | InterruptedException e) {
            // is TCP Server Disconnected
            if (e.getClass().getName().equals("java.net.SocketException")){
                System.out.println("Server Disconnected");
            }
        }
    }
}
