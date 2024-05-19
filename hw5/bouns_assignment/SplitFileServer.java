
import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
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
				OutputStream os = clientSocket.getOutputStream();
				int read;
				// get command from client
				byte[] commandBuffer = new byte[1024];
				read = is.read(commandBuffer);
				String commandName = new String(commandBuffer, 0, read);
				os.write("ok".getBytes());
				
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
						String fileName = in.readLine();
						File srcFile = new File(fileName);
						if (srcFile.exists() && srcFile.isFile()){
							long size = srcFile.length();
							out.println(Long.toString(size));
						} else {
							out.println("fail to get file Size");
							continue;
						}
						System.out.println("Request from client to send: " + fileName);

						String statusRes = in.readLine();
						if (!statusRes.equals("ok")){
							continue;
						}
						
						try (FileInputStream fis = new FileInputStream(srcFile)){
							byte[] b = new byte[1];
							while (fis.read(b) > 0){
								os.write(b);
							}
							System.out.println(fileName + " send successful");
						} catch (Exception e) {
							System.out.println("read Error!");
							continue;
						}
						
					}
					default -> System.out.println("Invalid command");
				}
			}
		} catch (Exception e) {
			System.out.println("Error!");
		}
	}
}