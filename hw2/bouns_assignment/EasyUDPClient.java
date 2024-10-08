/*
 EasyUDPClient.java
 20190532 Sangyun Lee
 */
import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.Scanner;

public class EasyUDPClient {
    public static void main(String[] args) {
        // server info
        final String SERVER_ADDRESS = "127.0.0.1";
        final int PORT = 30532;
        final int TIMEOUT = 5000;

        // add SIGINT HANDLER
        Runtime.getRuntime().addShutdownHook(new Thread(() -> System.out.println("\nBye bye~")));

        // create UDP client socket
        try (DatagramSocket clientSocket = new DatagramSocket()) {
            // set timeout to socket
            clientSocket.setSoTimeout(TIMEOUT);
            InetAddress serverAddress = InetAddress.getByName(SERVER_ADDRESS);
            Scanner sc = new Scanner(System.in);
            System.out.printf("The client is running on port %d\n", clientSocket.getLocalPort());

            while (true) {
                System.out.println("<Menu>");
                System.out.println("1) convert text to UPPER-case");
                System.out.println("2) get my IP address and port number");
                System.out.println("3) get server request count");
                System.out.println("4) get server running time");
                System.out.println("5) exit");
                System.out.print("Input option: ");
                String strType = sc.nextLine().replaceFirst("\n", "");
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
                byte [] strTypeBytes = strType.getBytes();

                // wrap data with UDP Datagram packet
                DatagramPacket sendStrTypePacket =
                        new DatagramPacket(strTypeBytes, strTypeBytes.length, serverAddress, PORT);
                // send to specified server
                clientSocket.send(sendStrTypePacket);

                if (strType.equals("1")){
                    System.out.print("Input sentence: ");
                    String text = sc.nextLine().replaceFirst("\n", "");
                    startTime = System.nanoTime();

                    byte[] textBytes = text.getBytes();
                    DatagramPacket sendTextPacket =
                            new DatagramPacket(textBytes, textBytes.length, serverAddress, PORT);
                    clientSocket.send(sendTextPacket);
                } else if (strType.equals("5")){
                    throw new InterruptedException();
                }

                // create receive UDP datagram packet
                byte[] responseBytes = new byte[4096];
                DatagramPacket receiveResponsePacket =
                        new DatagramPacket(responseBytes, responseBytes.length);
                // receive data to Datagram
                clientSocket.receive(receiveResponsePacket);

                // get Data from Datagram
                String response = new String(receiveResponsePacket.getData());
                long endTime = System.nanoTime();
                System.out.println("Reply from server: " + response);
                System.out.printf("RTT = %fms\n", (endTime - startTime) / 1e+6);
                Thread.sleep(1000); // Optional delay
            }
        } catch (IOException | InterruptedException  e) {
            if (e.getClass().getName().equals("java.net.SocketTimeoutException")){
                System.out.println("Server Disconnected");
            }
        }
    }
}
