/*
SplitFileClient.java
20190532 sang yun lee
*/

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.Socket;

class SplitFileClient {
	public static void main(String[] args) {
		// argument handling
		if (args.length != 3){
			System.out.println("Invalid argument.");
			System.exit(1);
		}
		// get argument from system args
		String commandName = args[1];
		String fileName = args[2];

		// server Info hardcoding
		String firstServerName = "127.0.0.1";
		String firstServerPort = "40532";
		String secondServerPort = "50532";

		// put command case
		if (commandName.equals("put")){
			sendFile(fileName, firstServerName, firstServerPort, 0);
		}
	}

	// put file to server
	private static void sendFile(String fileName, String serverName, String serverPort, int part){
		try (
			Socket socket = new Socket(serverName, Integer.parseInt(serverPort));
			PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
			BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()))
		) {

		} catch(IOException e){
			System.out.println("failt to connect server");
			System.exit(1);
		}
	}
}

/*!SECTION
 * PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
			BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream())
 */