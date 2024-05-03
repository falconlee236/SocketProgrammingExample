use std::io::{self, ErrorKind, Read, Write};
use std::net::TcpStream;
use std::sync::mpsc::{self, TryRecvError};
use std::thread;
use std::time::Duration;
use std::env::args;
use std::process::exit;
use std::collections::HashMap;

const MSG_SIZE: usize = 1024;
// server info
const SERVER_IP: &str = "127.0.0.1";
const SERVER_PORT: usize = 20532;

fn main() {
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
}