mkdir -p $HOME/.containerum/src/json_templates
DIRECTORY=$(cd install $0 && pwd)
cp $DIRECTORY/CONFIG.json $HOME/.containerum/CONFIG.json
chmod 777 -R $HOME/.containerum/