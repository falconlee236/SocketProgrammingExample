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
    let mut client_socket = TcpStream::connect(server_address).expect("stream failed to connect");
    client_socket.set_nonblocking(true).expect("failed to initialize non-blocking");

    let (tx, rx) = mpsc::channel::<String>();
    thread::spawn(move || loop {
        // Read message:
        let mut buff = vec![0; MSG_SIZE];
        match client_socket.read_exact(&mut buff) {
            Ok(_) => {
                let msg_byte_vec = buff.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
                let msg = String::from_utf8(msg_byte_vec).expect("invalid utf8 message");
                println!("---------------------------------------------------------");
                println!("{}", msg);
                println!("---------------------------------------------------------");
            },
            Err(ref err) if err.kind() == ErrorKind::WouldBlock => (),
            Err(_) => {
                println!("connection with server was served");
                break;
            }
        }

        // Receive message from channel and Write message to the server
        match rx.try_recv() {
            Ok(msg) => {
                let mut buff = msg.clone().into_bytes();
                buff.resize(MSG_SIZE, 0);
                client_socket.write_all(&buff).expect("writing to socket failed");
            },
            Err(TryRecvError::Empty) => (),
            Err(TryRecvError::Disconnected) => break
        }

        thread::sleep(Duration::from_millis(100));
    });

    // println!("\nwhat is your name?");
    // let mut name_buff = String::new();
    // io::stdin().read_line(&mut name_buff).expect("Reading from stdin failed");
    // let name = name_buff.trim().to_string();
    // println!("\nPlease enter a message to send");
    // send nickname to server
    if tx.send(nickname.to_string()).is_err() {exit(1)};


    loop {
        let mut buff = String::new();
        io::stdin().read_line(&mut buff).expect("Reading from stdin failed");
        let msg = format!("{}{}{:?}{}", nickname, &String::from("님이 "), &buff.trim().to_string(), &String::from("을(를) 입력하셨습니다."));
        if tx.send(msg).is_err() { break }
    }

}