import http.client, subprocess, platform
from dist2math import MathFunc

# See https://stackoverflow.com/a/44793537, https://docs.python.org/3/library/threading.html (first paragraph)
from multiprocessing import Process, Pool

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

# How many digits until we split up the workload to multiple threads?
# TODO: put this in a configuration file.
DIGITS_BEFORE_SEP = 1000


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
    print("Server stopped running, either digits needed are done,")
    print("or there was some form of power failure. Goodbye.")
    break
  instruction = val[0]
  # Set ourselves as busy
  Client.SetStatus(connection, BUSY, client_id)

  print(instruction)
  match instruction:
    case "COMP":
      dig_count  = int(val[1])
      offset     = int(val[3])
      max_digits = int(val[5])
      print(f"digit count: {dig_count} | offset: {offset}")

      start = time()
      if offset > max_digits:
        break
      
      value = 0
      if dig_count > DIGITS_BEFORE_SEP:
        # How many times do we want to seperate into?
        nT = 0
        if threads > 2:
          nT = threads//3
        else:
          nT = 2
        print(f"Splitting into {nT} different workers.")

        task = [(offset+(i * dig_count//nT), dig_count//nT) for i in range(nT)]
        print(task)
        with Pool(processes=nT) as p:
          r = p.starmap(MathFunc.GetActual, task)
          # Concatenate the two, then send it back
          value = ''.join(r)
      else:
        value = MathFunc.GetActual(offset, dig_count)
      
      print(f"Elapsed time: {time()-start}")
      accuracy += 1

      try:
        zz = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&data={str(value)}&type=data&timing={time()-start}")
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
