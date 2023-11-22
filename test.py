import websocket
import threading
import time

def on_message(ws, message):
    print(f"Received: {message}")


def on_error(ws, error):
    print(f"Error: {error}")


def on_close(ws, c ,a):
    print("### Closed ###")


def on_open(ws):
    print("### Connected ###")

def run_forever():
    ws.run_forever()

if __name__ == "__main__":
    token = "robotics"  # Replace with your authentication token
    header = {"Authorization": f"Bearer {token}"}
    ws = websocket.WebSocketApp("ws://localhost:8080/api/v1/device/hardware-status",
                              header=header,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close)
    ws.on_open = on_open
    run_forever()
    # threading.Thread(target=run_forever)
    # try:
    #     while True:
    #         time.sleep(0.1)
    # except KeyboardInterrupt:
    #     print("Keyboard interrupt detected. Closing connection.")
    #     ws.close()
