## pippin ("pi ping")

An google app engine apllication for trakcing raspberry pis.

Your Raspberry Pi on the local network, but where? The DHCP lease might have expired and your Pi will find itself with a new IP adress.

A "ping 192.168.0.255" could reveal the location, but now you can automatically track your Raspberrys from the comfort of your own home:

#### App engine setup

The "pippin-app" subdirectory contains the app egine app. Deploy with
    
    goapp deploy -oauth -application app-id pippin-app

Local test:
     
     goapp serve -clear_datastore pippin-app

#### Raspberry pi setup

The Raspberry Pi registers itself with a POST request:

    curl --data "name=test-pi-1&ip=192.168.0.1" $APP_URL_/pi

The "registerpi.py" script can be used to auto-detect and register the ip.

#### Usage

See all registered Raspberries:

    curl $APP_URL/

Get infor for one Raspberry as a JSON structure:

    curl $URL/pi/raspberry-name

Get the IP for a raspberry as plain text (useful for scripting):

    curl $URL/pi/raspberry-name/ip

