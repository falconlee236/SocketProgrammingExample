use std::io::{self, Read, Write};
use std::net::TcpStream;
use std::thread;
use std::env::args;
use std::process::exit;
use std::collections::HashMap;
use std::time::Instant;
use std::sync::{Arc, Mutex};

use ctrlc::set_handler;

const MSG_SIZE: usize = 1024;
// server info
const SERVER_IP: &str = "127.0.0.1";
const SERVER_PORT: usize = 20532;

fn main() {
    // start time
    let start = Arc::new(Mutex::new(Instant::now()));
    // get system argument 
    let args: Vec<String> = args().collect();
    // system args must have 2
    if args.len() != 2 {
        println!("This program must be run a one name argument");
        exit(1);
    }
    // get nickname from system args
    let nickname = &args[1];
    // define command Map
    let mut command_map : HashMap<&str, u8>  = HashMap::new();
    command_map.insert("ls", 0x01);
    command_map.insert("secret", 0x02);
    command_map.insert("except", 0x03);
    command_map.insert("ping", 0x04);
    command_map.insert("quit", 0x05);
    
    // server address ip
    let server_address = format!("{}:{}", SERVER_IP, SERVER_PORT);
    // make client socket;
    let mut client_socket = TcpStream::connect(server_address).expect("stream failed to connect");
    // write to server
    if client_socket.write(nickname.as_bytes()).is_err() {};

    // read from server
    let mut access_buff = vec![0; MSG_SIZE];
    if client_socket.read(&mut access_buff).is_err() {};
    let access_msg = access_buff.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
    let access_msg = String::from_utf8(access_msg).expect("invalid utf8 message");
    // msg to split by new line delimeter
    let access_res: Vec<&str> = access_msg.splitn(2, "\n").collect();
    // print message
    println!("{}", access_res[1]);
    // if status is error
    if access_res[0] == "404" {
        exit(1);
    }

    let mut sigint_client_socket = client_socket.try_clone().expect("failed clone");
    set_handler(move || {
        println!("\ngg~\n");
        if sigint_client_socket.write(&[5]).is_err() {};
        exit(1);
    }).expect("Error setting Ctrl+C handler");

    let mut read_client_socket = client_socket.try_clone().expect("failed clone");
    let start_clone = Arc::clone(&start);
    thread::spawn(move || loop{
        // Read message:
        let mut msg_res = vec![0; MSG_SIZE];
        match read_client_socket.read(&mut msg_res) {
            // read success
            Ok(n) => {
                // ping command case
                if n == 1 || msg_res[0] == 4 {
                    // get duration since start time
                    let start = start_clone.lock().unwrap();
                    let since = Instant::now().duration_since(*start);
                    let nanosecond = since.as_nanos();
                    println!("RTT = {}ms", nanosecond as f64 / 1e+6);
                } else if n > 0{ // read message case
                    let msg_byte_vec = msg_res.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
                    let msg = String::from_utf8(msg_byte_vec).expect("invalid utf8 message");
                    println!("{}", msg);
                } else {
                    println!("Server connection closed");
                    exit(1);
                }
            },
            Err(_) => {
                println!("Server connection closed");
                exit(1);
            }
        }
    });

    loop {
        // input msg from user
        let mut msg_buff = String::new();
        io::stdin().read_line(&mut msg_buff).expect("Reading from stdin failed");
        let msg_input = msg_buff.trim_end_matches("\n").to_string();
        println!();
        // // find command index
        match msg_input.find("\\") {
            // if command
            Some(index) => {
                if index == 0 {
                    let msg_arr: Vec<&str> = msg_input.splitn(2, " ").collect();
                    let command = &msg_arr[0][1..];
                    // encoding command to byte
                    match command_map.get(command) {
                        Some(&encoding) => {
                            // command argument error handling
                            if ((encoding == 1 || encoding == 4 || encoding == 5) && msg_arr.len() != 1) || 
                                ((encoding == 2 || encoding == 3) && msg_arr.len() != 2) {
                                    println!("Invalid command");
                                    continue;
                            }
                            if msg_arr.len() > 1 {
                                let msg = format!("{} {}", encoding, msg_arr[1]);
                                if client_socket.write(msg.as_bytes()).is_err() {
                                    println!("Server connection closed");
                                    exit(1);
                                }
                            } else {
                                // start calculate start time
                                let mut start = start.lock().unwrap();
                                *start = Instant::now();
                                if client_socket.write(&[encoding]).is_err() {
                                    println!("Server connection closed");
                                    exit(1);
                                }
                            }
                        },
                        // cannot find command table
                        None => println!("Invalid command")
                    }
                } else {
                    if client_socket.write(msg_input.as_bytes()).is_err() {
                        println!("Server connection closed");
                        exit(1);
                    }
                }
            },
            None => {
                if client_socket.write(msg_input.as_bytes()).is_err() {
                    println!("Server connection closed");
                    exit(1);
                }
            }
        }
    }
}