# import websocket
# import threading
# import time

# def on_message(ws, message):
#     print(f"Received: {message}")

# def on_error(ws, error):
#     print(f"Error: {error}")


# def on_close(ws, c ,a):
#     print("### Closed ###")


# def on_open(ws):
#     print("### Connected ###")

# def run_forever():
#     ws.run_forever()


# if __name__ == "__main__":
#     token = "robotics"  # Replace with your authentication token
#     header = {"Authorization": f"Bearer {token}"}
#     ws = websocket.WebSocketApp("ws://localhost:8080/api/v1/watchtower/log-stream?container=watchtower",
#                               header=header,
#                               on_message=on_message,
#                               on_error=on_error,
#                               on_close=on_close)
#     ws.on_open = on_open
#     run_forever()
#     # threading.Thread(target=run_forever)
#     # try:
#     #     while True:
#     #         time.sleep(0.1)
#     # except KeyboardInterrupt:
#     #     print("Keyboard interrupt detected. Closing connection.")
#     #     ws.close()

import asyncio
import websockets


async def on_ping(websocket, ping_data):
    print(f"Received ping: {ping_data}")
    await websocket.pong(ping_data)
    print("Sent pong in response to ping.")


async def on_message(websocket, message):
    print(f"Received message: {message}")


async def check_ping_messages():
    # Replace with your WebSocket server URI
    uri = "ws://localhost:8080/api/v1/supervisor/log-stream?container=robotics_supervisor"

    async with websockets.connect(uri) as websocket:
        while True:
            message = await websocket.recv()
            if message.startswith("ping"):
                # Extract ping data (if any)
                ping_data = message.split(" ")[1] if len(
                    message.split(" ")) > 1 else None
                await on_ping(websocket, ping_data)
            else:
                await on_message(websocket, message)

if __name__ == "__main__":
    asyncio.run(check_ping_messages())
