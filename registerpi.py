import socket
import urllib
import urllib2

# EDIT these and replace with your app instance and pi name
appstore_url = "http://my-app.appspot.com/pi"
pi_name = "dev-pi"

# Find a public-facing ip adress
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM); 
s.connect(("8.8.8.8", 80));
ip = s.getsockname()[0];
s.close();
print("Going with ip " + ip); 

# POST a registration request witht the pi name and ip.
data = urllib.urlencode({'name' : pi_name, 'ip'  : ip })
content = urllib2.urlopen(url = appstore_url, data = data).read()
print (content)
