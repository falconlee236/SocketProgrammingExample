/*
SplitFileClient.java
20190532 sang yun lee
*/

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.PrintWriter;
import java.net.Socket;

class SplitFileClient {
	public static void main(String[] args) {
		// argument handling
		if (args.length != 2){
			System.out.println("Invalid argument.");
			System.exit(1);
		}
		// get argument from system args
		String commandName = args[0];
		String fileName = args[1];

		// server Info hardcoding
		String firstServerName = "127.0.0.1";
		String firstServerPort = "40532";
		String secondServerName = "127.0.0.1";
		String secondServerPort = "50532";

		// put command case
		if (commandName.equals("put")){
			sendFile(fileName, firstServerName, firstServerPort, 0);
			sendFile(fileName, secondServerName, secondServerPort, 1);
		}
	}

	// put file to server
	private static void sendFile(String fileName, String serverName, String serverPort, int part){
		try (
			Socket socket = new Socket(serverName, Integer.parseInt(serverPort));
			PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
			BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()));
			OutputStream os = socket.getOutputStream();
		) {
			String commandStr = "put";
			// send to server
			out.println(commandStr);
			// get from server
			String commandRes = in.readLine();
			System.out.println(commandRes);
			if (!commandRes.equals("ok")){
				System.out.println("fail to receive fileName");
				throw new IOException();
			}
			System.out.println("Request to server to put :" + fileName);
			// get file Object
			File srcFile = new File(fileName);

			int idx = fileName.lastIndexOf('.');
			String fileExtension = fileName.substring(idx+1);
			fileName = fileName.substring(0, idx);
			fileName = String.format("%s-part%d.%s", fileName, part + 1, fileExtension);
			out.println(fileName);
			// get from server
			String fileNameRes = in.readLine();
			if (!fileNameRes.equals("ok")){
				System.out.println("fail to receive fileName");
				throw new IOException();
			}

			if (srcFile.exists() && srcFile.isFile()){
				long size = srcFile.length();
				out.println(Long.toString(size));
				String fileSizeRes = in.readLine();
				if (!fileSizeRes.equals("ok")){
					System.out.println("fail to receive fileSize");
					throw new IOException();
				}
			}

			try (FileInputStream fis = new FileInputStream(srcFile)){
				int byteCount = 0;
				byte[] b = new byte[1];
				while (fis.read(b) > 0){
					if (byteCount % 2 == part){
						os.write(b);
					}
					byteCount++;
				}
				System.out.println(fileName + " send successful");
			} catch (Exception e) {
				System.out.println("File Error");
				System.exit(1);
			}

		} catch(IOException e){
			System.out.println("fail to connect server");
			System.exit(1);
		}
	}
}
