
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
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
		try ( // create server Socket
			ServerSocket serverSocket = new ServerSocket(Integer.parseInt(serverPort))
		){
			System.out.println("wait for client request");
			while (true) { 
				// accept client request
				Socket clientSocket = serverSocket.accept();
				InputStream is = clientSocket.getInputStream();
				OutputStream os = clientSocket.getOutputStream();
				int read;
				// get command from client
				byte[] commandBuffer = new byte[1024];
				read = is.read(commandBuffer);
				String commandName = new String(commandBuffer, 0, read);
				// send ok to client
				os.write("ok".getBytes());
				
				// commandName case divide
				switch (commandName){
					case "put" -> {
						// get fileName from client
						byte[] fileNameBuffer = new byte[1024];
						read = is.read(fileNameBuffer);
						String fileName = new String(fileNameBuffer, 0, read);
						os.write("ok".getBytes());
						System.out.println("Recieved file: " + fileName);
						// get fileSize from client
						byte[] fileSizeBuffer = new byte[1024];
						read = is.read(fileSizeBuffer);
						String fileSizeString = new String(fileSizeBuffer, 0, read);
						os.write("ok".getBytes());
						// create FileOutputStream object write byte to File
						try (
							FileOutputStream fos = new FileOutputStream(fileName)
						){
							// convert String to int 
							int fileSize = Integer.parseInt(fileSizeString);
							long receivedBytes = 0;
							byte[] buffer = new byte[1024];
							int bytesRead;
							// read byte from client
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
						byte[] fileNameBuffer = new byte[1024];
						read = is.read(fileNameBuffer);
						String fileName = new String(fileNameBuffer, 0, read);
						File srcFile = new File(fileName);
						if (srcFile.exists() && srcFile.isFile()){
							long size = srcFile.length();
							os.write(Long.toString(size).getBytes());
						} else {
							os.write("fail to get file Size".getBytes());
							continue;
						}
						System.out.println("Request from client to send: " + fileName);

						byte[] statusResBuffer = new byte[1024];
						read = is.read(statusResBuffer);
						String statusRes = new String(statusResBuffer, 0, read);
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