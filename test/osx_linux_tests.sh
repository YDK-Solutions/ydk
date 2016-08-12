PY_GENERATE="python2.7"
PY_TEST="python2.7"

while getopts g:t: option 
do case "${option}" in 
g) PY_GENERATE=${OPTARG};; 
t) PY_TEST=${OPTARG};; 
esac 
done

printf "Generating api in $PY_GENERATE"
printf "Testing api in $PY_TEST"

virtualenv -p $PY_GENERATE gen_env
source gen_env/bin/activate
pip install coverage
pip install -r requirements.txt

python generate.py --profile profiles/test/ydktest.json
deactivate

virtualenv -p $PY_TEST test_env
source test_env/bin/activate
pip install coverage
pip install -r requirements.txt
pip install gen-api/python/dist/*
export PYTHONPATH=gen-api/python:$PYTHONPATH
python gen-api/python/tests/test_sanity_codec.py
deactivate

rm -rf gen_env
rm -rf test_env