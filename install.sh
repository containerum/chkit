sudo pip3 install -r requirements.txt
sudo mkdir -p $HOME/.containerium/src
sudo cp -r ./*  $HOME/.containerium/src/
sudo chmod +x client.py
sudo ln -srf $HOME/.containerium/src/client.py /usr/bin/client
sudo cp $HOME/.containerium/src/CONFIG.json $HOME/.containerium/CONFIG.json
sudo chmod 777 $HOME/.containerium/CONFIG.json
sudo chmod 777 $HOME/.containerium/src/json_templates/run.json