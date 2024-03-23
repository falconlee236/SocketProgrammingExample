import java.io.*;
import java.net.*;
import java.util.*;

public class EasyTCPServer {
    public static void main(String[] args) {
        final int PORT = 30532;

        System.out.printf("The server is ready to receive on port %d\n", PORT);
        try {
            ServerSocket serverSocket = new ServerSocket(PORT);
            Socket clientSocket = serverSocket.accept();

            BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));
            PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true);

            String inputLine;
            while ((inputLine = in.readLine()) != null) {
                System.out.println("Client connected: " + clientSocket);
                System.out.println("Command " + inputLine);
                String result = "";
                if (inputLine.equals("1")){
                    result = in.readLine();
                }
                System.out.println("Server received result " + result);
                out.println("Server received: " + inputLine);
            }

            System.out.println("Client disconnected: " + clientSocket);
            clientSocket.close();
            serverSocket.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}