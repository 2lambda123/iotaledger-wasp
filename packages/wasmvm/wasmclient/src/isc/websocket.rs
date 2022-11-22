use crate::*;
// use tungstenite::*;

#[derive(Clone)]
struct Client {
    url: String,
}

impl Client {
    pub fn connect(url: &str) -> errors::Result<Self> {
        let client = Client {
            url: url.to_owned(),
        };

        match tungstenite::connect(&client.url) {
            Ok(_v) => return Ok(client),
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
    use std::{net::TcpListener, thread::spawn};
    use tungstenite::{
        accept, accept_hdr,
        handshake::server::{Request, Response},
    };

    #[test]
    fn client_connect() {
        let url = "localhost:3012";
        spawn(move || mock_server(url));
        let client = Client::connect(&format!("ws://{}", url)).unwrap();
        assert!(client.url == format!("ws://{}", url));
    }

    fn mock_server(url: &str) {
        let server = TcpListener::bind(url).unwrap();
        for stream in server.incoming() {
            spawn(move || {
                accept(stream.unwrap()).unwrap();
            });
        }
    }
}
