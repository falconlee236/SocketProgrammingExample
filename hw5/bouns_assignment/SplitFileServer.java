/*
SplitFileServer.java
20190532 sang yun lee
*/

class SplitFileServer {
	public static void main(String[] args){
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
		String secondServerName = "127.0.0.1";
		String secondServerPort = "50532";
	}
}