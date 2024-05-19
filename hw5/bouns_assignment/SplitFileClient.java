/*
SplitFileClient.java
20190532 sang yun lee
*/

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.PrintWriter;
import java.net.Socket;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

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
		final String threadFileName = fileName;

		// server Info hardcoding
		String firstServerName = "127.0.0.1";
		String firstServerPort = "40532";
		String secondServerName = "127.0.0.1";
		String secondServerPort = "50532";

		ExecutorService executor = Executors.newFixedThreadPool(2);
		CountDownLatch latch = new CountDownLatch(2);
		// put command case
		if (commandName.equals("put")){
			Runnable task1 = () -> {
				try {
					sendFile(threadFileName, firstServerName, firstServerPort, 0);		
				} finally {
					latch.countDown();
				}
			};
			Runnable task2 = () -> {
				try {
					sendFile(threadFileName, secondServerName, secondServerPort, 1);
				} finally {
					latch.countDown();
				}
			};
			executor.submit(task1);
			executor.submit(task2);
			// Wait for all tasks to complete
            try {
                latch.await();  // Wait until the count reaches 0
                executor.shutdown();  // Shutdown the executor
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                System.out.println("Thread interrupted");
            }
		}// get command case
		else if (commandName.equals("get")){
			Runnable task1 = () -> {
				try {
					recieveFile(threadFileName, firstServerName, firstServerPort, 0);		
				} finally {
					latch.countDown();
				}
			};
			Runnable task2 = () -> {
				try {
					recieveFile(threadFileName, secondServerName, secondServerPort, 1);
				} finally {
					latch.countDown();
				}
			};
			executor.submit(task1);
			executor.submit(task2);
			// Wait for all tasks to complete
            try {
                latch.await();  // Wait until the count reaches 0
                executor.shutdown();  // Shutdown the executor
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                System.out.println("Thread interrupted");
            }

			int idx = fileName.lastIndexOf('.');
			String fileExtension = fileName.substring(idx);
			fileName = fileName.substring(0, idx);
			String outputFileName = String.format("%s-merged%s", fileName, fileExtension);
			String tmpFileName1 = String.format("%s-part%d%stmp%s", fileName, 1, fileExtension, fileExtension);
			String tmpFileName2 = String.format("%s-part%d%stmp%s", fileName, 2, fileExtension, fileExtension);
			File tmpFile1 = new File(tmpFileName1);
			File tmpFile2 = new File(tmpFileName2);
			try (
				FileOutputStream fos = new FileOutputStream(outputFileName);
				FileInputStream fis1 = new FileInputStream(tmpFile1);
				FileInputStream fis2 = new FileInputStream(tmpFile2)
			){
				long byteCnt = 0;
				int finishCnt = 0;
				while (true){
					if (finishCnt == 2)
						break;
					if (byteCnt % 2 == 0){
						byte[] b = new byte[1];
						if (fis1.read(b) <= 0){
							finishCnt++;
							continue;
						}
						fos.write(b);
					} else {
						byte[] b = new byte[1];
						if (fis2.read(b) <= 0){
							finishCnt++;
							continue;
						}
						fos.write(b);
					}
					byteCnt++;
				}
				System.out.println(threadFileName + "file merge successful!");
			} catch (Exception e) {
				System.out.println("File Open Error");
				System.exit(1);
			}
			if (tmpFile1.exists() && tmpFile1.delete()){
				System.out.println("delete Tempfile1 successful");
			} else {
				System.out.println("delete TempFile1 failed");
			}
			if (tmpFile2.exists() && tmpFile2.delete()){
				System.out.println("delete Tempfile2 successful");
			} else {
				System.out.println("delete TempFile2 failed");
			}
		}
	}

	// put file to server
	private static void sendFile(String fileName, String serverName, String serverPort, int part){
		try (
			Socket socket = new Socket(serverName, Integer.parseInt(serverPort));
			PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
			BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()));
			OutputStream os = socket.getOutputStream();
			InputStream is = socket.getInputStream();
		) {
			int read;

			// prepare command string to send command
			String commandStr = "put";
			// send to server
			os.write(commandStr.getBytes());
			// get from server
			byte[] commandBuffer = new byte[1024];
			read = is.read(commandBuffer);
			String commandRes = new String(commandBuffer, 0, read);
			if (!commandRes.equals("ok")){
				System.out.println("fail to receive fileName");
				throw new IOException();
			}
			System.out.println("Request to server to put :" + fileName);
			// get file Object
			File srcFile = new File(fileName);

			int idx = fileName.lastIndexOf('.');
			String fileExtension = fileName.substring(idx);
			fileName = fileName.substring(0, idx);
			fileName = String.format("%s-part%d%s", fileName, part + 1, fileExtension);
			os.write(fileName.getBytes());
			// get from server
			byte[] fileNameBuffer = new byte[1024];
			read = is.read(fileNameBuffer);
			String fileNameRes = new String(fileNameBuffer, 0, read);
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

	private static void recieveFile(String fileName, String serverName, String serverPort, int part){
		System.out.println("Request to server to get : " + fileName);
		try (
			Socket socket = new Socket(serverName, Integer.parseInt(serverPort));
			PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
			BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()));
			InputStream is = socket.getInputStream();
		) {
			String commandStr = "get";
			// send to server
			out.println(commandStr);
			// get from server
			String commandRes = in.readLine();
			if (!commandRes.equals("ok")){
				System.out.println("fail to receive fileName");
				throw new IOException();
			}
			
			int idx = fileName.lastIndexOf('.');
			String fileExtension = fileName.substring(idx);
			fileName = fileName.substring(0, idx);
			fileName = String.format("%s-part%d%s", fileName, part + 1, fileExtension);
			out.println(fileName);
			// get file Size from server
			String fileSizeBuffer = in.readLine();
			long fileSize = Long.parseLong(fileSizeBuffer);
			out.println("ok");
			
			
			// get file Object
			File tmpFile = new File(fileName + "tmp" + fileExtension);
			try (FileOutputStream fos = new FileOutputStream(tmpFile)){
				long receivedBytes = 0;
				byte[] buffer = new byte[1024];
				int bytesRead;
				while (receivedBytes < fileSize && (bytesRead = is.read(buffer)) != -1) { 
					fos.write(buffer, 0, bytesRead);
					receivedBytes += bytesRead;	
				}
				System.out.println(fileName + " store successful");
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
