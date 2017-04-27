sudo pip3 install -r requirements.txt
sudo mkdir -p $HOME/.containerum/src
sudo cp -r ./*  $HOME/.containerum/src/
sudo chmod +x client.py
sudo ln -srf $HOME/.containerum/src/client.py /usr/bin/client
sudo cp $HOME/.containerum/src/CONFIG.json $HOME/.containerum/CONFIG.json
sudo chmod 777 $HOME/.containerum/CONFIG.json
sudo chmod 777 $HOME/.containerum/src/json_templates/run.json