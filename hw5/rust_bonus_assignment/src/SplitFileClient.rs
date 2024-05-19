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
	let file_name = &args[2];

	if command_name == "put" {
		send_file(file_name, SERVER_IP1, SERVER_PORT1, 0);
		send_file(file_name, SERVER_IP2, SERVER_PORT2, 1);
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