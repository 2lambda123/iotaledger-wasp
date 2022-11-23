use crate::*;
use std::net;
use tungstenite::*;

pub type SocketType =
    tungstenite::WebSocket<tungstenite::stream::MaybeTlsStream<std::net::TcpStream>>;

struct Client {
    url: String,
    socket: SocketType,
}

impl Client {
    pub fn connect(url: &str) -> errors::Result<Self> {
        match tungstenite::connect(url) {
            Ok((socket, _res)) => {
                return Ok(Client {
                    url: url.to_owned(),
                    socket: socket,
                })
            }
            Err(e) => return Err(e.to_string()),
        }
    }
    pub fn subscribe() -> errors::Result<()> {
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use crate::websocket::Client;
    use std::{clone, net::TcpListener, ops::Index, thread::spawn};
    use tungstenite::{
        accept,
        handshake::server::{Request, Response},
    };

    #[test]
    fn client_connect() {
        let url = "ws://localhost:3012";
        mock_server(url, None);
        let client = Client::connect(&format!("ws://{}", url)).unwrap();
        assert!(client.url == format!("ws://{}", url));
    }

    #[test]
    fn client_subscribe() {
        let url = "ws://localhost:3012";
        mock_server(url, Some("hi"));
        let client = Client::connect(&format!("ws://{}", url)).unwrap();
    }

    fn mock_server(input_url: &str, response_msg: Option<&str>) {
        let ws_prefix = "ws://";
        let url = match input_url.to_string().strip_prefix(ws_prefix) {
            Some(u) => u.to_string(),
            None => input_url.to_string(),
        };
        let msg = match response_msg {
            Some(m) => m.to_string(),
            None => "".to_string(),
        };
        spawn(move || {
            let server = TcpListener::bind(url).unwrap();
            for stream in server.incoming() {
                let mut socket = accept(stream.unwrap()).unwrap();
                if !msg.is_empty() {
                    socket.write_message(msg.clone().into()).unwrap();
                }
            }
        });
    }
}
