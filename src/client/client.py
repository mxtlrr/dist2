import http.client, subprocess, platform
from dist2math import MathFunc

threads = 0
match platform.system():
  case "Windows":
    threads = int(__import__("os").environ["NUMBER_OF_PROCESSORS"])
  case "Linux":
    threads = int(subprocess.check_output(['nproc']).decode('utf-8').replace("\n",""))

# Stuff to send to server to update status
READY = 0
BUSY  = 1
accuracy = 1

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

from time import time

while True:
  try:
    val = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&type=request").split(" ")
  except http.client.RemoteDisconnected:
    print("Server stopped runnig, either digits needed are done,")
    print("or there was some form of power failure. Goodbye.")
    break
  instruction = val[0]
  # Set ourselves as busy
  Client.SetStatus(connection, BUSY, client_id)

  print(instruction)
  match instruction:
    case "COMP":
      dig_count = int(val[1])
      offset    = int(val[3])
      print(f"digit count: {dig_count} | offset: {offset}")

      start = time()

      value = MathFunc.GetActual(offset, dig_count) #MathFunc.GetOffset(MathFunc.CompSqrt2(accuracy), offset, dig_count)
      
      print(f"Elapsed time: {time()-start}")
      # print(value)
      accuracy += 1

      try:
        zz = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&data={str(value)}&type=data")
      except http.client.RemoteDisconnected:
          print("Server terminated. Goodbye")
          break
      try:
        Client.SetStatus(connection, READY, client_id)
      except ConnectionResetError:
        print("Server terminated. All digits calculated successfully. Goodbye.")
        break

    # Handle issues from server
    case "Sorry,":
      print("ERR! Something failed with validating ourselves with the server")
      print("Either retry or if this is reoccuring, open an issue on GitHub")
      print("at <https://github.com/mxtlrr/dist2/issues>.\n\nExiting.")
      break

  print(f"Status: {status}")


connection.close()
