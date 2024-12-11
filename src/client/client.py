import http.client, subprocess
from dist2math import MathFunc


threads = int(subprocess.check_output(['nproc']).decode('utf-8').replace("\n",""))

# Stuff to send to server to update status
READY = 0
BUSY  = 1
accuracy = 200

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

for i in range(1):
  val = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&type=request").split(" ")
  instruction = val[0]
  match instruction:
    case "COMP":
      dig_count = int(val[1])
      offset    = int(val[3])
      print(f"Digit count: {dig_count} | Offset: {offset}")
      # Set ourselves as busy
      Client.SetStatus(connection, BUSY, client_id)


      # Each iteration improves the approximation, roughly doubling the number of correct digits...
      # it can be modeled with f(x) = 1.2084*(1.8146^x)-1.124, R^2 = 0.9995, so at
      # accuracy = 25, we'd get ~3.56 million correct digits, though this is just a model to fit the
      # data, and may not be correct.
      # To keep up with this, we'll increment one after each iteration
      value = MathFunc.GetOffset(MathFunc.CompSqrt2(accuracy), offset, dig_count)
      accuracy += 1 # f(26) = ~6,467,132

      zz = Client.sendReq(connection, "GET", f"/data?client_id={client_id}&data={str(value)}&type=data")
      Client.SetStatus(connection, READY, client_id)


connection.close()