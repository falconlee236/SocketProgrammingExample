/*
SplitFileClient.rs
20190532 sang yun lee
*/

use std::env::args;
use std::io::{Read, Write, BufReader};
use std::net::TcpStream;
use std::path::Path;
use std::process::exit;
use std::fs::File;
use std::fs::remove_file;
use std::thread;

// message buffer size
const MSG_SIZE: usize = 1024;
// server info
const SERVER_IP1: &str = "127.0.0.1";
const SERVER_PORT1: &str = "40532";
const SERVER_IP2: &str = "127.0.0.1";
const SERVER_PORT2: &str = "50532";

fn main() {
	// get system argument
	let args: Vec<String> = args().collect();
	if args.len() != 3 {
		println!("Invalid argument");
		exit(1);
	}
	// get argument from system args
	let command_name = &args[1];
	let file_name = args[2].clone();

	if command_name == "put" {
		let file_name1 = file_name.clone();
		let file_name2 = file_name.clone();
		let handle1 = thread::spawn(move || {
			send_file(&file_name1, SERVER_IP1, SERVER_PORT1, 0);
		});
		let handle2 = thread::spawn(move || {
			send_file(&file_name2, SERVER_IP2, SERVER_PORT2, 1);
		});
		handle1.join().unwrap();
        handle2.join().unwrap();
	} else if command_name == "get" {
		let file_name1 = file_name.clone();
		let file_name2 = file_name.clone();
		let handle1 = thread::spawn(move || {
			receive_file(&file_name1, SERVER_IP1, SERVER_PORT1, 0);
		});
		let handle2 = thread::spawn(move || {
			receive_file(&file_name2, SERVER_IP2, SERVER_PORT2, 1);
		});
		handle1.join().unwrap();
        handle2.join().unwrap();

		// get file extension from fileName
		let file_extension = match Path::new(&file_name).extension().and_then(|s| s.to_str()) {
			Some(ext) => ".".to_string() + ext,
			None => "".to_string()
		};
		let file_name = &file_name[0..(file_name.len() - file_extension.len())];
		// set tmp file name
		let tmp_file_name1 = format!("{}-part{}{}tmp{}", &file_name, 1, file_extension, file_extension);
		let tmp_file_name2 = format!("{}-part{}{}tmp{}", &file_name, 2, file_extension, file_extension);
		// make result file name
		let file_name = format!("{}-merged{}", file_name, file_extension);
		// file open
		let tmp_file1 = match File::open(&tmp_file_name1) {
			Ok(file) => file,
			Err(error) => {
				eprintln!("{} File open Error: {}", &tmp_file_name1, error);
				exit(1);
			}
		};
		let tmp_file2 = match File::open(&tmp_file_name2) {
			Ok(file) => file,
			Err(error) => {
				eprintln!("{} File open Error: {}", &tmp_file_name2, error);
				exit(1);
			}
		};
		// dst file create
		match File::create(&file_name) {
			Ok(mut dst_file) => {
				let mut reader1 = BufReader::new(tmp_file1);
				let mut reader2 = BufReader::new(tmp_file2);

				let mut byte_cnt: i64 = 0;
				let mut finish_cnt: i32 = 0;
				while finish_cnt != 2 {
					let mut buffer = [0; 1];
					if byte_cnt % 2 == 0 {
						// read 1 byte from file
						match reader1.read(&mut buffer) {
							Ok(0) => { // EOF
								finish_cnt += 1;
								continue;
							},
							Ok(_) => {
								let byte = buffer[0];
								dst_file.write(&[byte]).unwrap();
							}
							Err(error) => {
								eprintln!("tmp file 1read failed: {}", error);
								exit(1);
							}
						}
					} else {
						// read 1 byte from file
						match reader2.read(&mut buffer) {
							Ok(0) => { // EOF
								finish_cnt += 1;
								continue;
							},
							Ok(_) => {
								let byte = buffer[0];
								dst_file.write(&[byte]).unwrap();
							}
							Err(error) => {
								eprintln!("tmp file 1read failed: {}", error);
								exit(1);
							}
						}
					}
					byte_cnt += 1;
				}
				println!("{}{} file merge sucessful!", file_name, file_extension);
				match remove_file(tmp_file_name1) {
					Ok(_) => println!("tmp1 file remove sucess!"),
					Err(_) => {
						eprintln!("tmp1 file remove failed!");
						exit(1);
					}
				};
				match remove_file(tmp_file_name2) {
					Ok(_) => println!("tmp2 file remove sucess!"),
					Err(_) => {
						eprintln!("tmp2 file remove failed!");
						exit(1);
					}
				};
			},
			Err(e) => {
				eprintln!("file creation error: {}", e);
				exit(1);
			}
		}
	} else {
		println!("Invalid argument");
	}
}

fn send_file(file_name: &str, server_name: &str, server_port: &str, part: i32){
	// create client socket
	let mut client_socket = TcpStream::connect(format!("{}:{}", server_name, server_port)).expect("stream failed to connect");
	// prepare command string to send command
	let command_str = "put";
	// write to server
	client_socket.write(command_str.as_bytes()).unwrap();
	// read from server
	let mut command_buffer = vec![0; MSG_SIZE];
	client_socket.read(&mut command_buffer).unwrap();
	let command_res = command_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
	let command_res = String::from_utf8(command_res).expect("invalid utf8 message");
	// server response is not ok
	if command_res != "ok" {
		println!("fail to receive command name");
		exit(1);
	}
	println!("Request to server to put : {}", file_name);

	// file open
	let src_file = match File::open(file_name) {
        Ok(file) => file,
        Err(error) => {
            eprintln!("File open Error: {}", error);
            exit(1);
        }
    };

	// get file extension from fileName
	let file_extension = match Path::new(file_name).extension().and_then(|s| s.to_str()) {
		Some(ext) => ".".to_string() + ext,
		None => "".to_string()
	};
	let file_name = &file_name[0..(file_name.len() - file_extension.len())];
	let file_name = format!("{}-part{}{}", file_name, part + 1, file_extension);
	// write to server
	client_socket.write(file_name.as_bytes()).unwrap();
	// read from server
	let mut filename_buffer = vec![0; MSG_SIZE];
	client_socket.read(&mut filename_buffer).unwrap();
	let filename_res = filename_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
	let filename_res = String::from_utf8(filename_res).expect("invalid utf8 message");
	// server response is not ok
	if filename_res != "ok" {
		println!("fail to receive file name");
		exit(1);
	}
	
	// get file size from File object metadata
	let file_size = match src_file.metadata() {
		Ok(metadata) => metadata.len(),
		Err(_) => 0
	}.to_string();
	client_socket.write(file_size.as_bytes()).unwrap();
	// read from server
	let mut filesize_buffer = vec![0; MSG_SIZE];
	client_socket.read(&mut filesize_buffer).unwrap();
	let filesize_res = filesize_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
	let filesize_res = String::from_utf8(filesize_res).expect("invalid utf8 message");
	// server response is not ok
	if filesize_res != "ok" {
		println!("fail to receive file name");
		exit(1);
	}
	
	let mut reader = BufReader::new(src_file);
    let mut buffer = [0; 1];
	let mut byte_cnt = 0;
    loop {
		// read 1 byte from file
        match reader.read(&mut buffer) {
            Ok(0) => break, // EOF
            Ok(_) => {
                let byte = buffer[0];
                if byte_cnt % 2 == part {
					client_socket.write(&[byte]).unwrap();
				}
            }
            Err(error) => {
                eprintln!("read failed: {}", error);
				exit(1);
            }
        }
		byte_cnt += 1;
    }
	println!("{} send sucessful!", file_name);
}

fn receive_file(file_name: &str, server_name: &str, server_port: &str, part: i32){
	println!("Request to server to get : {}", file_name);

	// create client socket
	let mut client_socket = TcpStream::connect(format!("{}:{}", server_name, server_port)).expect("stream failed to connect");
	// prepare command string to send command
	let command_str = "get";
	client_socket.write(command_str.as_bytes()).unwrap();
	// read from server
	let mut command_buffer = vec![0; MSG_SIZE];
	client_socket.read(&mut command_buffer).unwrap();
	let command_res = command_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
	let command_res = String::from_utf8(command_res).expect("invalid utf8 message");
	// server response is not ok
	if command_res != "ok" {
		println!("fail to receive command name");
		exit(1);
	}

	// get file extension from fileName
	let file_extension = match Path::new(file_name).extension().and_then(|s| s.to_str()) {
		Some(ext) => ".".to_string() + ext,
		None => "".to_string()
	};
	let file_name = &file_name[0..(file_name.len() - file_extension.len())];
	let file_name = format!("{}-part{}{}", file_name, part + 1, file_extension);
	// send file name to server
	client_socket.write(file_name.as_bytes()).unwrap();

	// try to get file size
	let mut filesize_buffer = vec![0; MSG_SIZE];
	client_socket.read(&mut filesize_buffer).unwrap();
	let filesize_res = filesize_buffer.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
	let filesize = String::from_utf8(filesize_res).expect("invalid utf8 message");
	// parse file size string to long
	let file_size: u64 = match filesize.parse() {
		Ok(n) => {
			client_socket.write("ok".as_bytes()).unwrap();
			n
		}
		Err(e) => {
            eprintln!("convert error : {}", e);
            exit(1);
        }
	};

	// file create
	match File::create(format!("{}tmp{}", &file_name, file_extension)) {
		Ok(mut file) => {
			let mut received_bytes = 0;
			let mut buffer = vec![0; MSG_SIZE];
			while received_bytes < file_size {
				match client_socket.read(&mut buffer) {
					Ok(n) => {
						if n == 0 {
							break;
						}
						file.write_all(&buffer[..n]).unwrap();
						received_bytes += n as u64;
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
	println!("{} file store successful!", file_name);

}