/*
SplitFileServer.rs
20190532 sang yun lee
*/

use std::env::args;
use std::fs::File;
use std::io::{Read, Write, BufReader};
use std::net::TcpListener;
use std::process::exit;

// message buffer size
const MSG_SIZE: usize = 1024;

fn main() {
	// get system argument
	let args: Vec<String> = args().collect();
	if args.len() != 2 {
		println!("Invalid argument");
		exit(1);
	}
	// serer port from system argument
	let server_port = &args[1];

	// create server_socket
	let server_socket = TcpListener::bind(format!("0.0.0.0:{}", server_port)).expect("Listener failed to bind");
	println!("Waiting for connections...");

	loop {
		// set client connection
		if let Ok((mut socket, _)) = server_socket.accept() {
			// read from client
			let mut command_buffer = vec![0; MSG_SIZE];
			socket.read(&mut command_buffer).unwrap();
			let command_name = command_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
			let command_name = String::from_utf8(command_name).expect("invalid utf8 message");
			// write to client
			socket.write("ok".as_bytes()).unwrap();

			// put case
			if command_name == "put"{
				//get file Name from client
				// read from client
				let mut filename_buffer = vec![0; MSG_SIZE];
				socket.read(&mut filename_buffer).unwrap();
				let file_name = filename_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
				let file_name = String::from_utf8(file_name).expect("invalid utf8 message");
				// write to client
				socket.write("ok".as_bytes()).unwrap();

				// get file size from client
				// read from client
				let mut filesize_buffer = vec![0; MSG_SIZE];
				socket.read(&mut filesize_buffer).unwrap();
				let filesize_string = filesize_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
				let filesize_string = String::from_utf8(filesize_string).expect("invalid utf8 message");
				// write to client
				socket.write("ok".as_bytes()).unwrap();
				let file_size = match filesize_string.parse() {
					Ok(num) => num,
					Err(err) => {
						eprintln!("convert failed : {}", err);
						0
					}
				};

				// file create
				match File::create(&file_name) {
					Ok(mut file) => {
						let mut received_bytes = 0;
						let mut buffer = vec![0; MSG_SIZE];
						while received_bytes < file_size {
							match socket.read(&mut buffer) {
								Ok(n) => {
									if n == 0 {
										break;
									}
									file.write_all(&buffer[..n]).unwrap();
									received_bytes += n;
								},
								Err(e) => {
									println!("Error occurred: {:?}", e);
									continue;
								}
							}
						}
					},
					Err(error) => {
						eprintln!("File open Error: {}", error);
						exit(1);
					}
				};
				println!("{} save successful!", file_name);
			} else if command_name == "get" {
				//get file Name from client
				// read from client
				let mut filename_buffer = vec![0; MSG_SIZE];
				socket.read(&mut filename_buffer).unwrap();
				let file_name = filename_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
				let file_name = String::from_utf8(file_name).expect("invalid utf8 message");

				// try to open file
				let src_file = match File::open(&file_name) {
					Ok(file) => {
						// try to get file metadata
						match file.metadata() {
							Ok(metadata) => {
								// send file size to client
								socket.write(metadata.len().to_string().as_bytes()).unwrap();
								file
							},
							Err(err) => {
								eprintln!("fail to get file size: {}", err);
								socket.write("error!".as_bytes()).unwrap();		
								continue;
							}
						}
					},
					Err(err) => {
						eprintln!("fail to open file: {}", err);
						socket.write("error!".as_bytes()).unwrap();
						continue;
					}
				};

				println!("Request from client to send : {}", &file_name);
				// read from client
				let mut status_buffer = vec![0; MSG_SIZE];
				socket.read(&mut status_buffer).unwrap();
				let status_res = status_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
				let status_res = String::from_utf8(status_res).expect("invalid utf8 message");
				// client status check
				if status_res != "ok" {
					eprintln!("fail to receive file size");
					continue;
				}

				// set Reader Buffer in file
				let mut reader = BufReader::new(src_file);
				let mut buffer = [0; 1];
				loop {
					// read 1 byte from file
					match reader.read(&mut buffer) {
						// EOF
						Ok(0) => {
							println!("{} send successful!", file_name);
							break
						},
						Ok(_) => {
							let byte = buffer[0];
							socket.write(&[byte]).unwrap();
						}
						Err(error) => {
							eprintln!("read failed: {}", error);
							break;
						}
					}
				}
			} else {
				println!("Unknown command!: {}", command_name);
			}
		}
	}
}