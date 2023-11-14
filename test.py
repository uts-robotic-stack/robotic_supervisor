import websocket


def on_message(ws, message):
    print(f"Received: {message}")


def on_error(ws, error):
    print(f"Error: {error}")


def on_close(ws):
    print("### Closed ###")


def on_open(ws):
    print("### Connected ###")

if __name__ == "__main__":
    token = "robotics"  # Replace with your authentication token
    header = {"Authorization": f"Bearer {token}"}
    ws = websocket.WebSocketApp("ws://localhost:8080/api/v1/watchtower/logs?container_name=watchtower",
                              header=header,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close)
    ws.on_open = on_open
    ws.run_forever()

    try:
        while True:
            pass
    except KeyboardInterrupt:
        print("Keyboard interrupt detected. Closing connection.")
        ws.close()