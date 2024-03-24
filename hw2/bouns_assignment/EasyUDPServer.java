import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.*;
import java.time.Duration;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;

public class EasyUDPServer {
    public static void main(String[] args) {
        LocalTime serverStartTime = LocalTime.now();
        DateTimeFormatter dtf = DateTimeFormatter.ofPattern("HH:mm:ss");
        final int PORT = 30532;

        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("\nBye bye~");
        }));
        try (DatagramSocket serverSocket = new DatagramSocket(PORT)){
            while (true){
                System.out.printf("The server is ready to receive on port %d\n", PORT);

                byte[] inputLineBytes = new byte[1000];
                DatagramPacket receiveInputLinePacket = new DatagramPacket(inputLineBytes, inputLineBytes.length);
                serverSocket.receive(receiveInputLinePacket);

                String inputLine = new String(receiveInputLinePacket.getData());
                int reqNum = 0;
                System.out.println("Command " + inputLine);

                String result = "";
                if (inputLine.equals("1")){
                    byte[] textBytes = new byte[2048];
                    DatagramPacket receiveTextBytesPacket = new DatagramPacket(textBytes, textBytes.length);
                    serverSocket.receive(receiveTextBytesPacket);

                    result = new String(receiveTextBytesPacket.getData()).toUpperCase();
                } else if (inputLine.equals("2")){
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

                byte[] resultBytes = new byte[4096];
                InetAddress clientAddress = receiveInputLinePacket.getAddress();
                int clientPort = receiveInputLinePacket.getPort();
                DatagramPacket sendResultPacket =
                        new DatagramPacket(resultBytes, resultBytes.length, clientAddress, clientPort);
                serverSocket.send(sendResultPacket);
                reqNum++;
            }
        } catch (IOException e) {
//            e.printStackTrace();
        }
    }

    private static String formatDuration(Duration duration) {
        long hours = duration.toHours();
        long minutes = duration.toMinutes() % 60;
        long seconds = duration.getSeconds() % 60;

        return String.format("%02d:%02d:%02d", hours, minutes, seconds);
    }
}
