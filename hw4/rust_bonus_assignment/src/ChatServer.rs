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
    let server_address = format!("127.0.0.1{}{}", &String::from(":"), SERVER_PORT);

    let server_socket = TcpListener::bind(server_address).expect("Lister failed to bind");

    println!("\nWaiting for client connection..");

    let client_map : Arc<Mutex<HashMap<String, TcpStream>>>  = Arc::new(Mutex::new(HashMap::new()));
    let total_client_num = Arc::new(Mutex::new(0));
    loop {
        if let Ok((mut socket, addr)) = server_socket.accept() {
            //get client's nickname from client
            let mut nickname_buff = vec![0; MSG_SIZE];
            match socket.read(&mut nickname_buff) {
                Ok(_) => {
                    let msg = nickname_buff.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
                    let nickname = String::from_utf8(msg).expect("invalid utf8 message");
                    let mut status_code : usize = 200;
                    let send_msg = {
                        let mut client_map = client_map.lock().unwrap();
                        let mut total_client_num = total_client_num.lock().unwrap();
                        if *total_client_num == 8{
                            status_code = 404;
                            format!("{}\n[chatting room full. cannot connect.]\n", status_code)
                        }
                        else if client_map.contains_key(&nickname) {
                            status_code = 404;
                            format!("{}\n[nickname already used by another user. cannot connect.]\n", status_code)
                        } else {
                            *total_client_num += 1;
                            client_map.insert(nickname.clone(), socket.try_clone().expect("failed to clone client"));
                            format!("{}\n[welcome {} to CAU net-class chat room at {}.]\n[There are {} users in the room.]", status_code, &nickname, server_socket.local_addr().expect("failed to get address"), total_client_num)
                        }
                    };

                    if socket.write(send_msg.as_bytes()).is_err() {};
                    if status_code == 200 {
                        println!("[{} has joined from {}.]\n[There are {} users in room.]", &nickname, addr, total_client_num.lock().unwrap());
                        let client_map = client_map.clone();
                        let total_client_num = total_client_num.clone();
                        // aka. TCPClientHandler
                        thread::spawn(move || loop {
                            // Read message:
                            let mut msg_res = vec![0; MSG_SIZE];
                            match socket.read(&mut msg_res) {
                                // read success
                                Ok(n) => {
                                    // read message case
                                    if n > 0 { 
                                        let msg_byte_vec = msg_res.into_iter().take_while(|&x| x != 0).collect::<Vec<_>>();
                                        if msg_byte_vec[0] > 0 && msg_byte_vec[0] < 6 {
                                            // ls command
                                            if msg_byte_vec[0] == 1 {
                                                let mut msg = String::new();
                                                for (other_nickname, stream) in client_map.lock().unwrap().iter_mut(){
                                                    // get remove addr
                                                     let remote_addr = stream.peer_addr().expect("Failed to get peer address");
                                                    msg.push_str(&format!("<{}, {}, {}>\n", other_nickname, remote_addr.ip(), remote_addr.port()));
                                                }
                                                if socket.write(msg.as_bytes()).is_err() {}
                                            } else if msg_byte_vec[0] == 4 { // ping command
                                                if socket.write(&[4]).is_err() {}
                                            } else if msg_byte_vec[0] == 5 { //quit command
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
                                            // let msg = String::from_utf8(msg_byte_vec[2..].to_vec()).expect("invalid utf8 message");       
                                        } else {
                                            let msg = String::from_utf8(msg_byte_vec).expect("invalid utf8 message");
                                            let msg = format!("{}> {}\n", &nickname, msg);
                                            for (other_nickname, stream) in client_map.lock().unwrap().iter_mut(){
                                                if &nickname != other_nickname {
                                                    if stream.write(msg.as_bytes()).is_err() {}
                                                }
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
