use std::io::{ErrorKind, Read, Write};
use std::net::TcpListener;
use std::net::TcpStream;
use std::sync::mpsc;
use std::thread;
use std::time::Duration;
use std::collections::HashMap;

const MSG_SIZE: usize = 1024;
const SERVER_PORT: usize = 20532;

fn main() {
    // get server address
    let server_address = format!("0.0.0.0{}{}", &String::from(":"), SERVER_PORT);

    let server_socket = TcpListener::bind(server_address).expect("Lister failed to bind");
    server_socket.set_nonblocking(true).expect("failed to initialize non blocking listener");

    println!("\nWaiting for client connection..");

    let mut clients = vec![];
    let mut client_map : HashMap<String, TcpStream>  = HashMap::new();
    let mut total_client_num = 0;
    let (tx, rx) = mpsc::channel::<String>();
    loop {
        if let Ok((mut socket, addr)) = server_socket.accept() {
            //get client's nickname from client
            let mut nickname_buff = vec![0; MSG_SIZE];
            let mut nickname_socket = socket.try_clone().expect("failed to clone client");
            match nickname_socket.read_exact(&mut nickname_buff) {
                Ok(_) => {
                    let tx = tx.clone();
                    let msg = nickname_buff.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
                    let nickname = String::from_utf8(msg).expect("invalid utf8 message");
                    let mut status_code : usize = 200;
                    let send_msg = if total_client_num == 8{
                        status_code = 404;
                        format!("{}\n[chatting room full. cannot connect.]\n", status_code)
                    }
                    else if client_map.contains_key(&nickname) {
                        status_code = 404;
                        format!("{}\n[nickname already used by another user. cannot connect.]\n", status_code)
                    } else {
                        total_client_num += 1;
                        client_map.insert(nickname.clone(), socket.try_clone().expect("failed to clone client"));
                        format!("{}\n[welcome {} to CAU net-class chat room at {}.]
                            \n[There are {} users in the room.]", status_code, &nickname, server_socket.local_addr().expect("failed to get address"), total_client_num)
                    };

                    println!("server recived {}", send_msg);
                    // tx.send(msg).expect("failed to send message to rx");
                },
                Err(_) => {}
            }
            println!("client {} connected", addr);
            let tx = tx.clone();
            clients.push(socket.try_clone().expect("failed to clone client"));
            let msg = format!("※ [{}]님이 입장하셨습니다. ※",addr);
            tx.send(msg).expect("failed to send message to rx");
            thread::spawn(move || loop {
                let mut buff = vec![0; MSG_SIZE];
                match socket.read_exact(&mut buff) {
                    Ok(_) => {
                        let msg = buff.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
                        let msg = String::from_utf8(msg).expect("invalid utf8 message");

                        println!("[{}]{}", addr, msg);

                        let msg = format!("{}{}{}{}", &String::from("["), addr, &String::from("]"), &msg);
                        tx.send(msg).expect("failed to send message to rx");
                    },
                    Err(ref err) if err.kind() == ErrorKind::WouldBlock => (),
                    Err(_) => {
                        println!("Closing connection with [{}]", addr);
                        let msg = format!("※ [{}]님이 퇴장하셨습니다. ※",addr);
                        tx.send(msg).expect("failed to send message to rx");
                        break;
                    }
                }

                sleep();
            });
        }

        if let Ok(msg) = rx.try_recv() {
            clients = clients.into_iter().filter_map(|mut client| {
                let mut buff = msg.clone().into_bytes();
                buff.resize(MSG_SIZE, 0);
                client.write_all(&buff).map(|_| client).ok()
            })
                .collect::<Vec<_>>();
        }
        sleep();
    }
}

fn sleep() {
    thread::sleep(Duration::from_millis(100));
}
