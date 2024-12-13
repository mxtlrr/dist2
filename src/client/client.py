import http.client, subprocess
from dist2math import MathFunc


threads = int(subprocess.check_output(['nproc']).decode('utf-8').replace("\n",""))

# Stuff to send to server to update status
READY = 0
BUSY  = 1
accuracy = 200

status = READY

class Client:
  def __init__(self):
    pass
  def sendReq(connect: http.client.HTTPConnection, method: str, string: str) -> str:
    connect.request(method,string)
    return connect.getresponse().read().decode("utf-8")
  def SetStatus(connect: http.client.HTTPConnection, v2: int, client_id: int) -> None:
    status = v2
    connect.request("GET", f"/setstatus?val={v2}&client_id={client_id}")
    r = connection.getresponse()
    r.read()
    r.close()


connection = http.client.HTTPConnection("127.0.0.1", 8080, timeout=10)
print("Connected to dist2 server...")


# register ourselves.
v = Client.sendReq(connection, "GET", f"/register?threads={threads}").split(" ")
if v[0] != "OK":
  print("Failure registering to dist2 server!")
  connection.close()

client_id = int(v[1])-1
print(f"Registered as client {client_id}")

# Set ourselves as ready to recieve
Client.SetStatus(connection, READY, client_id)

while True:
  val = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&type=request").split(" ")
  instruction = val[0]
  # Set ourselves as busy
  Client.SetStatus(connection, BUSY, client_id)

  print(instruction)
  match instruction:
    case "COMP":
      dig_count = int(val[1])
      offset    = int(val[3])
      print(f"digit count: {dig_count} | offset: {offset}")

      value = MathFunc.GetOffset(MathFunc.CompSqrt2(accuracy), offset, dig_count)
      print(value)
      accuracy += 1

      zz = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&data={str(value)}&type=data")
      Client.SetStatus(connection, READY, client_id)
    case "CHECK":
      offset     = int(val[3])
      submitted  = val[5]

      accuracy *= 2 # Increase accuracy to actually get a potentially
      # better result

      value = MathFunc.GetOffset(MathFunc.CompSqrt2(accuracy), offset, 3)
      if submitted != value:
        # Let's send the full 20 digits.
        value = MathFunc.GetOffset(MathFunc.CompSqrt2(accuracy), offset, 20)
        zz = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&type=check&return={value}")
      else:
        zz = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&type=check&return=OK")
      Client.SetStatus(connection, READY, client_id)

      # Reset accuracy.
      accuracy /= 2

    # Handle issues from server
    case "Sorry,":
      print("ERR! Something failed with validating ourselves with the server")
      print("Either retry or if this is reoccuring, open an issue on GitHub")
      print("at <https://github.com/mxtlrr/dist2/issues>.\n\nExiting.")
      break

  __import__("time").sleep(0.5)
  print(f"Status: {status}")


connection.close()