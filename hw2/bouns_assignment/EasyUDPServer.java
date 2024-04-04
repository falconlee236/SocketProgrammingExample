/*
* EasyUDPServer.java
* 20190532 Sangyun Lee
* */
import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.time.Duration;
import java.time.LocalTime;

public class EasyUDPServer {
    public static void main(String[] args) {
        LocalTime serverStartTime = LocalTime.now();
        final int PORT = 30532;

        // set SIGINT handler
        Runtime.getRuntime().addShutdownHook(new Thread(() -> System.out.println("\nBye bye~")));

        // create UDP socket
        try (DatagramSocket serverSocket = new DatagramSocket(PORT)){
            int reqNum = 0;
            while (true){
                System.out.printf("The server is ready to receive on port %d\n", PORT);

                // create receive datagram
                byte[] inputLineBytes = new byte[1000];
                DatagramPacket receiveInputLinePacket = new DatagramPacket(inputLineBytes, inputLineBytes.length);
                // receive command type data from client
                serverSocket.receive(receiveInputLinePacket);

                // get data from packet
                String inputLine = new String(receiveInputLinePacket.getData()).trim();
                System.out.println("Command " + inputLine);

                String result = "";
                if (inputLine.equals("1")){
                    // get Text data from client
                    byte[] textBytes = new byte[2048];
                    DatagramPacket receiveTextBytesPacket = new DatagramPacket(textBytes, textBytes.length);
                    serverSocket.receive(receiveTextBytesPacket);

                    result = new String(receiveTextBytesPacket.getData()).trim().toUpperCase();
                } else if (inputLine.equals("2")){
                    // get client ip, client port
                    result = String.format("client IP = %s, port = %d",
                            receiveInputLinePacket.getAddress().getHostAddress(),
                            receiveInputLinePacket.getPort());
                } else if (inputLine.equals("3")){
                    result = String.format("requests served = %d", reqNum);
                } else if (inputLine.equals("4")){
                    LocalTime currentServerTime = LocalTime.now();
                    Duration duration = Duration.between(serverStartTime, currentServerTime);
                    result = String.format("run time = %s", formatDuration(duration));
                }

                byte[] resultBytes = result.getBytes();
                // get dst address and port from recently received packet
                InetAddress clientAddress = receiveInputLinePacket.getAddress();
                int clientPort = receiveInputLinePacket.getPort();
                // create UDP Packet
                DatagramPacket sendResultPacket =
                        new DatagramPacket(resultBytes, resultBytes.length, clientAddress, clientPort);
                // send packet to client
                serverSocket.send(sendResultPacket);
                reqNum++;
            }
        } catch (IOException e) {
//            e.printStackTrace();
        }
    }
    // change to format
    private static String formatDuration(Duration duration) {
        long hours = duration.toHours();
        long minutes = duration.toMinutes() % 60;
        long seconds = duration.getSeconds() % 60;

        return String.format("%02d:%02d:%02d", hours, minutes, seconds);
    }
}
