/*
ChatServer.rs
20190532 sangyun lee
*/
use std::io::{Read, Write};
use std::process::exit;
use std::sync::{Arc, Mutex};
use std::net::TcpListener;
use std::net::TcpStream;
use std::thread;
use std::collections::HashMap;
use ctrlc::set_handler;

const MSG_SIZE: usize = 1024;
const SERVER_PORT: usize = 20532;

fn main() {
    //sigint handler
    set_handler(|| {
        println!("\ngg~\n");
        exit(127);
    }).expect("Error setting Ctrl+C handler");

    // get server address
    let server_address = format!("0.0.0.0{}{}", &String::from(":"), SERVER_PORT);

    // create server socket
    let server_socket = TcpListener::bind(server_address).expect("Lister failed to bind");

    println!("\nWaiting for client connection..");

    // create client info map, total number with Mutex
    // in rust every shared resource must wrap in Mutex and Arc type
    let client_map : Arc<Mutex<HashMap<String, TcpStream>>>  = Arc::new(Mutex::new(HashMap::new()));
    let total_client_num = Arc::new(Mutex::new(0));
    loop {
        // set clinet connection
        if let Ok((mut socket, addr)) = server_socket.accept() {
            //get client's nickname from client
            let mut nickname_buff = vec![0; MSG_SIZE];
            // read client name 
            match socket.read(&mut nickname_buff) {
                Ok(_) => {
                    // convert vec to string
                    let msg = nickname_buff.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
                    let nickname = String::from_utf8(msg).expect("invalid utf8 message");
                    let mut status_code : usize = 200;
                    let send_msg = {
                        // lock mutex, get data
                        let mut client_map = client_map.lock().unwrap();
                        let mut total_client_num = total_client_num.lock().unwrap();
                        // chat room is full
                        if *total_client_num == 8{
                            status_code = 404;
                            format!("{}\n[chatting room full. cannot connect.]\n", status_code)
                        } // already has nickname
                        else if client_map.contains_key(&nickname) {
                            status_code = 404;
                            format!("{}\n[nickname already used by another user. cannot connect.]\n", status_code)
                        } else { // otherwise - enter the chat room
                            // add client number
                            *total_client_num += 1;
                            // insert map to client socket
                            client_map.insert(nickname.clone(), socket.try_clone().expect("failed to clone client"));
                            format!("{}\n[welcome {} to CAU net-class chat room at {}.]\n[There are {} users in the room.]", status_code, &nickname, server_socket.local_addr().expect("failed to get address"), total_client_num)
                        }
                    };

                    // write msg to client
                    if socket.write(send_msg.as_bytes()).is_err() {};
                    // client enter the chat room case
                    if status_code == 200 {
                        println!("[{} has joined from {}.]\n[There are {} users in room.]", &nickname, addr, total_client_num.lock().unwrap());
                        // thread want to use Arc shared resoucre, clone shared data
                        // when rust create thread, outter scope's ownership move to thread
                        let client_map = client_map.clone();
                        let total_client_num = total_client_num.clone();
                        // aka. TCPClientHandler / create client handle thread
                        thread::spawn(move || loop {
                            // Read message:
                            let mut msg_res = vec![0; MSG_SIZE];
                            match socket.read(&mut msg_res) {
                                // read success
                                Ok(n) => {
                                    // read message case
                                    if n > 0 { 
                                        let msg_byte_vec = msg_res.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
                                        let command_type = msg_byte_vec[0];
                                        if command_type > 0 && command_type < 6 {
                                            // ls command
                                            if command_type == 1 {
                                                let mut msg = String::new();
                                                for (other_nickname, stream) in client_map.lock().unwrap().iter_mut(){
                                                    // get remove addr
                                                     let remote_addr = stream.peer_addr().expect("Failed to get peer address");
                                                    msg.push_str(&format!("<{}, {}, {}>\n", other_nickname, remote_addr.ip(), remote_addr.port()));
                                                }
                                                if socket.write(msg.as_bytes()).is_err() {}
                                            } else if command_type == 2 || command_type == 3 { // secret, except case
                                                let command_str = if command_type == 2 {
                                                    "\\secret"
                                                } else {
                                                    "\\except"
                                                };
                                                // get not command part in 3 string
                                                let msg_arr: Vec<&str> = std::str::from_utf8(&msg_byte_vec).expect("invalid utf8 message").splitn(3, " ").collect();
                                                // command parameter error
                                                if msg_arr.len() != 3 {
                                                    println!("Invalid command: {}", command_str);
                                                    continue;
                                                }

                                                // get target nickname, msg
                                                let target_nickname = msg_arr[1];
                                                let target_msg = msg_arr[2];

                                                // secret command
                                                if command_type == 2 {
                                                    // get target connetion info
                                                    match client_map.lock().unwrap().get(target_nickname) {
                                                        Some(mut stream) => {
                                                            let msg = format!("from: {}> {}\n", &nickname, target_msg);
                                                            if stream.write(msg.as_bytes()).is_err() {}
                                                        },
                                                        // cannot find target nickname in client map
                                                        None => {
                                                            println!("Invalid command: \\secret");
                                                            let msg = format!("[{} is not in the chat room.]\n", target_nickname);
                                                            if socket.write(msg.as_bytes()).is_err() {}
                                                            continue;
                                                        },
                                                    }
                                                }

                                                // except command
                                                if command_type == 3 {
                                                    match client_map.lock().unwrap().get(target_nickname) {
                                                        // cannot find target nickname in client map
                                                        Some(_) => {
                                                            let msg = format!("from: {}> {}\n", &nickname, target_msg);
                                                            for (other_nickname, stream) in client_map.lock().unwrap().iter_mut(){
                                                                // send msg excpet target nickname
                                                                if other_nickname != target_nickname {
                                                                    if stream.write(msg.as_bytes()).is_err() {}
                                                                }
                                                            }
                                                        }
                                                        None => {
                                                            println!("Invalid command: \\except");
                                                            let msg = format!("[{} is not in the chat room.]\n", target_nickname);
                                                            if socket.write(msg.as_bytes()).is_err() {}
                                                            continue;
                                                        },
                                                    }
                                                }

                                                // text filtering
                                                if target_msg.to_lowercase().contains("i hate professor") {
                                                    // subtract total_client_num
                                                    *total_client_num.lock().unwrap() -= 1;
                                                    let msg = format!("[{} is disconnected.]\n[There are {} users in the chat room.]\n", &nickname, *total_client_num.lock().unwrap());
                                                    // send to every client info of disconnecting
                                                    for (_, stream) in client_map.lock().unwrap().iter_mut(){
                                                        if stream.write(msg.as_bytes()).is_err() {}
                                                    }
                                                    // disconnect client
                                                    client_map.lock().unwrap().remove(&nickname);
                                                    println!("{}", msg);
                                                    break;
                                                }

                                            } else if command_type == 4 { // ping command
                                                if socket.write(&[4]).is_err() {}
                                            } else if command_type == 5 { //quit command
                                                // subtract client number
                                                *total_client_num.lock().unwrap() -= 1;
                                                // remove client info
                                                client_map.lock().unwrap().remove(&nickname);
                                                let msg = format!("\n[{} left the room.]\n[There are {} users now.]\n\n", &nickname, *total_client_num.lock().unwrap());
                                                // send msg to other client
                                                for (other_nickname, stream) in client_map.lock().unwrap().iter_mut(){
                                                    if &nickname != other_nickname {
                                                        if stream.write(msg.as_bytes()).is_err() {}
                                                    }
                                                }
                                                // print msg to server
                                                println!("{}", msg);
                                                break;
                                            }
                                        } else { // otherwise, get client message
                                            let msg = String::from_utf8(msg_byte_vec).expect("invalid utf8 message");
                                            let msg = format!("{}> {}\n", &nickname, msg);
                                            // send message to other client
                                            for (other_nickname, stream) in client_map.lock().unwrap().iter_mut(){
                                                if &nickname != other_nickname {
                                                    if stream.write(msg.as_bytes()).is_err() {}
                                                }
                                            }
                                            // text filtering
                                            if msg.to_lowercase().contains("i hate professor") {
                                                *total_client_num.lock().unwrap() -= 1;
                                                let msg = format!("[{} is disconnected.]\n[There are {} users in the chat room.]\n", &nickname, *total_client_num.lock().unwrap());
                                                for (_, stream) in client_map.lock().unwrap().iter_mut(){
                                                    if stream.write(msg.as_bytes()).is_err() {}
                                                }
                                                client_map.lock().unwrap().remove(&nickname);
                                                println!("{}", msg);
                                                break;
                                            }
                                        }
                                    } else { // client connection closed
                                        *total_client_num.lock().unwrap() -= 1;
                                        client_map.lock().unwrap().remove(&nickname);
                                        println!("\n[{} left the room.]\n[There are {} users now.]\n\n", &nickname, *total_client_num.lock().unwrap());
                                        break;    
                                    }
                                },
                                Err(_) => {}
                            }
                        });
                    }
                },
                Err(e) => {
                    println!("Error occurred: {:?}", e);
                }
            }
        }
    }
}
