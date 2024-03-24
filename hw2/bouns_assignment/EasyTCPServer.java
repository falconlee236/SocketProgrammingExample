/*
EasyTCPServer.java
20190532 Sangyun Lee
*/
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.ServerSocket;
import java.net.Socket;
import java.time.Duration;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;

public class EasyTCPServer {
    public static void main(String[] args) {
        LocalTime serverStartTime = LocalTime.now();
        DateTimeFormatter dtf = DateTimeFormatter.ofPattern("HH:mm:ss");
        final int PORT = 30532;

        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("\nBye bye~");
        }));
        try {
            while (true){
                System.out.printf("The server is ready to receive on port %d\n", PORT);
                ServerSocket serverSocket = new ServerSocket(PORT);
                Socket clientSocket = serverSocket.accept();

                BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));
                PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true);

                String inputLine;
                int reqNum = 0;
                while ((inputLine = in.readLine()) != null) {
                    System.out.println("Command " + inputLine);
                    String result = "";
                    if (inputLine.equals("1")){
                        result = in.readLine();
                        result = result.toUpperCase();
                    } else if (inputLine.equals("2")){
                        result = String.format("client IP = %s, port = %d",
                                clientSocket.getInetAddress(), clientSocket.getPort());
                    } else if (inputLine.equals("3")){
                        result = String.format("requests served = %d", reqNum);
                    } else if (inputLine.equals("4")){
                        LocalTime currentServerTime = LocalTime.now();
                        Duration duration = Duration.between(serverStartTime, currentServerTime);
                        result = String.format("run time = %s", formatDuration(duration));
                    }
                    out.println(result);
                    reqNum++;
                }
                clientSocket.close();
                serverSocket.close();
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