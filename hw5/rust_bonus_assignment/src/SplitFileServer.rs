/*
SplitFileServer.rs
20190532 sang yun lee
*/

use std::env::args;
use std::fs::File;
use std::io::{Read, Write};
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

				// file open
				match File::open(file_name) {
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
			}
		}
	}
}