import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.concurrent.atomic.AtomicInteger;

public class MyTCPServer {

    private static final AtomicInteger reqNum = new AtomicInteger(0);
    private static final long startTime = System.currentTimeMillis();

    public static void main(String[] args) {
        final int serverPort = 20532;

        try (ServerSocket serverSocket = new ServerSocket(serverPort)) {
            Runtime.getRuntime().addShutdownHook(new Thread(() -> System.out.println("\nBye bye~")));

            System.out.println("Server is ready to receive on port " + serverPort);

            while (true) {
                Socket clientSocket = serverSocket.accept();
                System.out.println("Connection request from " + clientSocket.getRemoteSocketAddress());

                new Thread(new TCPClientHandler(clientSocket)).start();
            }
        } catch (IOException e) {
            System.err.println("Error occurred while setting up server socket");
            e.printStackTrace();
        }
    }

    private static class TCPClientHandler implements Runnable {
        private final Socket clientSocket;

        public TCPClientHandler(Socket clientSocket) {
            this.clientSocket = clientSocket;
        }

        @Override
        public void run() {
            try (BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));
                 PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true)) {

                byte[] typeBuffer = new byte[1024];
                byte[] buffer = new byte[1024];

                while (true) {
                    int t = in.read(typeBuffer);
                    if (t == -1) {
                        break;
                    }
                    String typeStr = new String(typeBuffer, 0, t - 1);

                    System.out.println("Command " + typeStr);

                    switch (typeStr) {
                        case "1":
                            int count = in.read(buffer);
                            String input = new String(buffer, 0, count);
                            out.println(input.toUpperCase());
                            break;
                        case "2":
                            out.println(clientSocket.getRemoteSocketAddress());
                            break;
                        case "3":
                            out.println(reqNum.get());
                            break;
                        case "4":
                            long duration = System.currentTimeMillis() - startTime;
                            int hour = (int) (duration / 3600000);
                            int minute = (int) ((duration / 60000) % 60);
                            int second = (int) ((duration / 1000) % 60);
                            String totalRuntime = String.format("%02d:%02d:%02d\n", hour, minute, second);
                            out.println(totalRuntime);
                            break;
                        case "5":
                            out.println("-1");
                            break;
                    }
                    reqNum.incrementAndGet();
                }
            } catch (IOException e) {
                System.err.println("Error occurred while handling client request");
                e.printStackTrace();
            }
        }
    }
}