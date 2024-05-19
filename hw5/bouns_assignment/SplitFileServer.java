
import java.io.BufferedReader;
import java.io.FileOutputStream;
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
		if (args.length != 2){
			System.out.println("Invalid argument.");
			System.exit(1);
		}
		// server Port
		String serverPort = args[1];
		try (
			ServerSocket serverSocket = new ServerSocket(Integer.parseInt(serverPort))
		){
			while (true) { 
				Socket clientSocket = serverSocket.accept();
				BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));
				PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true);

				// get command from client
				String commandName = in.readLine();
				out.print("ok");

				if (commandName.equals("put")){
					String fileName = in.readLine();
					out.print("ok");
					System.out.println("Received file: " + fileName);

					String fileSizeBuffer = in.readLine();
					out.print("ok");
					int fileSize = Integer.parseInt(fileSizeBuffer);
					try (
						FileOutputStream fos = new FileOutputStream(fileName)
					){
						long receivedBytes = 0;
						while (receivedBytes >= fileSize) { 
							String n = in.readLine();
							receivedBytes += n.length();	
							fos.write(n.getBytes());
						}
					} catch (Exception e) {
						System.out.println("fail to create file");
						throw new Exception();
					}
					System.out.println(fileName + " file store sucessful!");
				}
			}
		} catch (Exception e) {}
	}
}