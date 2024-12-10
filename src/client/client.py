import http.client, subprocess
threads = int(subprocess.check_output(['nproc']).decode('utf-8').replace("\n",""))


# Stuff to send to server to update status

READY = 0
BUSY  = 1


class Client:
  def __init__(self):
    pass
  def sendReq(connect: http.client.HTTPConnection, method: str, string: str) -> str:
    connect.request(method,string)
    return connect.getresponse().read().decode("utf-8")
  def SetStatus(connect: http.client.HTTPConnection, v2: int, client_id: int) -> None:
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
print(f"We are client {client_id}")

# Set ourselves as ready to recieve
Client.SetStatus(connection, READY, client_id)

for i in range(5):
  val = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&type=request").split(" ")
  instruction = val[0]
  match instruction:
    case "COMP":
      # Set ourselves as busy
      Client.SetStatus(connection, BUSY, client_id)
      # Go compute something
      print("I am computing.")


      zz = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&data=VERYSPECIFICTHING&type=data")
      Client.SetStatus(connection, READY, client_id)


connection.close()