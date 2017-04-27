sudo pip3 install -r requirements.txt
sudo mkdir -p /var/lib/containerium/src
sudo cp -r ./*  /var/lib/containerium/src/
sudo chmod +x client.py
sudo ln -srf /var/lib/containerium/src/client.py /usr/bin/client
sudo cp /var/lib/containerium/src/CONFIG.json /var/lib/containerium/CONFIG.json
sudo chmod 777 /var/lib/containerium/CONFIG.json
sudo chmod 777 /var/lib/containerium/src/json_templates/run.json