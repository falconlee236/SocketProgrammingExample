
import java.io.BufferedReader;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.ServerSocket;
import java.net.Socket;

/*
SplitFileServer.java
20190532 sang yun lee
*/

class SplitFileServer {
	public static void main(String[] args){
		// argument handling
		if (args.length != 1){
			System.out.println("Invalid argument.");
			System.exit(1);
		}
		// server Port
		String serverPort = args[0];
		try (
			ServerSocket serverSocket = new ServerSocket(Integer.parseInt(serverPort))
		){
			System.out.println("wait for client request");
			while (true) { 
				Socket clientSocket = serverSocket.accept();
				BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));
				PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true);
				InputStream is = clientSocket.getInputStream();
				
				// get command from client
				String commandName = in.readLine();
				out.println("ok");
				
				switch (commandName){
					case "put" -> {
						String fileName = in.readLine();
						out.println("ok");
										
						String fileSizeBuffer = in.readLine();
						out.println("ok");
						try (
							FileOutputStream fos = new FileOutputStream(fileName)
						){
							int fileSize = Integer.parseInt(fileSizeBuffer);
							long receivedBytes = 0;
							byte[] buffer = new byte[1024];
							int bytesRead;
							while (receivedBytes < fileSize && (bytesRead = is.read(buffer)) != -1) {
								fos.write(buffer, 0, bytesRead);
								receivedBytes += bytesRead;
							}
						} catch (Exception e) {
							System.out.println("fail to create file");
							continue;
						}
						System.out.println(fileName + " file store sucessful!");
					}
					case "get" -> {
						
					}
					default -> System.out.println("Invalid command");
				}
			}
		} catch (Exception e) {
			System.out.println("Error!");
		}
	}
}